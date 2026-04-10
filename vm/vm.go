// Package vm implements the moonBASIC bytecode interpreter.
package vm

import (
	"fmt"
	"io"
	"os"

	"moonbasic/lineprof"
	"moonbasic/runtime"
	"moonbasic/vm/callstack"
	"moonbasic/vm/heap"
	"moonbasic/vm/opcode"
	"moonbasic/vm/value"
)

// VM is the moonBASIC virtual machine instance.
type VM struct {
	CallStack *callstack.Stack
	Registry  *runtime.Registry
	Globals   map[string]value.Value
	Heap      *heap.Store

	Program *opcode.Program
	Halted  bool

	// Debug / Trace
	Trace    bool      // If true, dump machine state after each instruction
	TraceOut io.Writer // Destination for trace output (default os.Stderr)

	// Profiler when non-nil receives one tick per instruction executed (by source line).
	Profiler lineprof.LineProfiler

	// PhysicsScratch holds the last Jolt/WASM SoA float payload (host-filled before [OpSyncPhysics]).
	PhysicsScratch []float64
}

// New creates a new VM instance with a linked registry and heap.
func New(reg *runtime.Registry, h *heap.Store) *VM {
	return &VM{
		CallStack: callstack.New(),
		Registry:  reg,
		Globals:   make(map[string]value.Value),
		Heap:      h,
		Halted:    false,
	}
}

// Runtime Error Helper
func (v *VM) runtimeError(msg string) error {
	frame := v.CallStack.Top()
	line := -1
	if frame != nil && frame.IP < len(frame.Chunk.SourceLines) {
		line = int(frame.Chunk.SourceLines[frame.IP])
	}
	where := "unknown source"
	if v.Program != nil && v.Program.SourcePath != "" {
		where = v.Program.SourcePath
	} else if frame != nil && frame.Chunk != nil && frame.Chunk.Name != "" {
		where = frame.Chunk.Name
	}
	if line >= 1 {
		return fmt.Errorf("[moonBASIC] Error in %s line %d:\n  %s", where, line, msg)
	}
	return fmt.Errorf("[moonBASIC] Error in %s:\n  %s", where, msg)
}

// reg returns a register from the current frame.
func (v *VM) reg(idx uint8) value.Value {
	frame := v.CallStack.Top()
	if frame == nil {
		return value.Nil
	}
	return frame.Registers[idx]
}

// setReg sets a register in the current frame.
func (v *VM) setReg(idx uint8, val value.Value) {
	frame := v.CallStack.Top()
	if frame != nil {
		frame.Registers[idx] = val
	}
}

// SyncPhysicsFromFloat32View resizes [PhysicsScratch] and copies guest floats for use by [opcode.OpSyncPhysics].
func (v *VM) SyncPhysicsFromFloat32View(src []float32) {
	if len(src) == 0 {
		v.PhysicsScratch = v.PhysicsScratch[:0]
		return
	}
	if cap(v.PhysicsScratch) < len(src) {
		v.PhysicsScratch = make([]float64, len(src))
	} else {
		v.PhysicsScratch = v.PhysicsScratch[:len(src)]
	}
	for i := range src {
		v.PhysicsScratch[i] = float64(src[i])
	}
}

// Execute runs the given program from its main entry point.
func (v *VM) Execute(prog *opcode.Program) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = v.runtimeError(fmt.Sprintf("VM Panic Recovery: %v", r))
		}
		if v.Registry != nil {
			if err != nil {
				v.Registry.SetLastScriptError(err)
			} else {
				v.Registry.ClearLastScriptError()
			}
		}
	}()
	v.Program = prog
	v.Halted = false
	v.Registry.Prog = prog
	v.Heap.SeedProgramStrings(prog.StringTable)
	defer func() { v.Registry.Prog = nil }()

	v.Registry.StackTraceFn = func() string { return v.FormatCallStack() }
	defer func() { v.Registry.StackTraceFn = nil }()

	v.Registry.TerminateVM = func() { v.Halted = true }
	defer func() { v.Registry.TerminateVM = nil }()

	v.Registry.EraseAllHandlesFn = v.EraseAllHandles
	defer func() { v.Registry.EraseAllHandlesFn = nil }()

	// Push the <MAIN> chunk as our first call frame (R0 as return reg, irrelevant for MAIN).
	v.CallStack.Push(prog.Main, 0, 0)

	for !v.Halted && v.CallStack.Depth() > 0 {
		frame := v.CallStack.Top()
		if frame.IP >= len(frame.Chunk.Instructions) {
			// Implicit return at end of chunk
			v.CallStack.Pop()
			continue
		}

		instr := frame.Chunk.Instructions[frame.IP]
		// Increment IP *before* execution so OpJump etc can override it.
		frame.IP++

		if err := v.step(instr); err != nil {
			return err
		}

		if v.Profiler != nil {
			ip := frame.IP - 1
			if ip >= 0 && ip < len(frame.Chunk.SourceLines) {
				if ln := int(frame.Chunk.SourceLines[ip]); ln > 0 {
					v.Profiler.RecordLine(ln)
				}
			}
		}

		if v.Trace {
			v.trace(instr)
		}
	}

	return nil
}

// EraseAllHandles runs [heap.Store.FreeAll], then sets every [value.KindHandle] in [VM.Globals]
// and on the operand stack to nil. Intended for ERASE ALL / FREE.ALL when tearing down resources;
// do not use in the middle of an expression that still expects a handle on the stack.
func (v *VM) EraseAllHandles() error {
	if v.Heap == nil {
		return fmt.Errorf("EraseAllHandles: heap not bound")
	}
	v.Heap.FreeAll()
	if v.Registry != nil {
		v.Registry.ResetModules()
	}
	for name, val := range v.Globals {
		if val.Kind == value.KindHandle {
			v.Globals[name] = value.Nil
		}
	}
	// Iterate through all active call frames and clear their registers
	v.CallStack.Range(func(f *callstack.Frame) {
		for i := range f.Registers {
			if f.Registers[i].Kind == value.KindHandle {
				f.Registers[i] = value.Nil
			}
		}
	})
	return nil
}

// step processes one instruction.
// Sub-implementations are in vm_arith.go and vm_control.go.
func (v *VM) step(i opcode.Instruction) error {
	switch i.Op {
	case opcode.OpPushInt:
		frame := v.CallStack.Top()
		val := frame.Chunk.IntConsts[i.Operand]
		v.setReg(i.Dst, value.FromInt(val))

	case opcode.OpPushFloat:
		frame := v.CallStack.Top()
		val := frame.Chunk.FloatConsts[i.Operand]
		v.setReg(i.Dst, value.FromFloat(val))

	case opcode.OpPushString:
		idx := int(i.Operand)
		if idx < 0 || idx >= len(v.Program.StringTable) {
			return v.runtimeError("PUSH_STRING: string pool index out of range")
		}
		v.setReg(i.Dst, value.FromStringIndex(int32(idx)))

	case opcode.OpPushBool:
		v.setReg(i.Dst, value.FromBool(i.Operand != 0))

	case opcode.OpPushNull:
		v.setReg(i.Dst, value.Nil)

	case opcode.OpPop:
		// No-op in register VM.

	case opcode.OpLoadGlobal:
		frame := v.CallStack.Top()
		name := frame.Chunk.Names[i.Operand]
		if val, ok := v.Globals[name]; ok {
			v.setReg(i.Dst, val)
		} else {
			v.setReg(i.Dst, value.Nil)
		}

	case opcode.OpStoreGlobal:
		frame := v.CallStack.Top()
		name := frame.Chunk.Names[i.Operand]
		v.Globals[name] = v.reg(i.SrcA)

	case opcode.OpMove:
		v.setReg(i.Dst, v.reg(i.SrcA))

	case opcode.OpLoadLocal:
		return v.runtimeError("OpLoadLocal: use register moves instead")

	case opcode.OpStoreLocal:
		return v.runtimeError("OpStoreLocal: use register moves instead")

	case opcode.OpHalt:
		v.Halted = true

	case opcode.OpSwap:
		a, b := v.reg(i.SrcA), v.reg(i.SrcB)
		v.setReg(i.SrcA, b)
		v.setReg(i.SrcB, a)

	default:
		// Delegate to specialized handlers for complex logic
		return v.dispatchComplex(i)
	}

	return nil
}

// trace dumps IP, opcode, and register state after each step.
func (v *VM) trace(instr opcode.Instruction) {
	out := v.TraceOut
	if out == nil {
		out = os.Stderr
	}
	frame := v.CallStack.Top()
	ip := -1
	line := -1
	chunk := "?"
	if frame != nil {
		ip = frame.IP - 1
		if ip < 0 {
			ip = 0
		}
		chunk = frame.Chunk.Name
		if ip < len(frame.Chunk.SourceLines) {
			line = int(frame.Chunk.SourceLines[ip])
		}
	}
	// Trace now dumps the first few registers of the current frame
	regs := "[]"
	if frame != nil {
		regs = fmt.Sprintf("%v", frame.Registers[:8]) // dump R0-R7
	}
	fmt.Fprintf(out, "[trace] %s L%d IP=%d %s | depth=%d regs=%v\n",
		chunk, line, ip, instr.String(), v.CallStack.Depth(), regs)
}
