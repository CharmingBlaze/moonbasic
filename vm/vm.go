// Package vm implements the moonBASIC bytecode interpreter.
package vm

import (
	"fmt"
	"io"
	"os"

	"moonbasic/runtime"
	"moonbasic/vm/callstack"
	"moonbasic/vm/heap"
	"moonbasic/vm/opcode"
	"moonbasic/vm/value"
)

// VM is the moonBASIC virtual machine instance.
type VM struct {
	Stack     []value.Value
	CallStack *callstack.Stack
	Registry  *runtime.Registry
	Globals   map[string]value.Value
	Heap      *heap.Store

	Program *opcode.Program
	Halted  bool

	// Debug / Trace
	Trace    bool      // If true, dump machine state after each instruction
	TraceOut io.Writer // Destination for trace output (default os.Stderr)

	// StackHygieneDebug when true (pipeline Options.Debug / CLI --info): if execution
	// finishes with values left on the operand stack, return an error (ARCHITECTURE §6).
	StackHygieneDebug bool

	// Profiler when non-nil receives one tick per instruction executed (by source line).
	Profiler *ProfileRecorder
}

// New creates a new VM instance with a linked registry and heap.
func New(reg *runtime.Registry, h *heap.Store) *VM {
	return &VM{
		Stack:     make([]value.Value, 0, 1024),
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

// Push a value onto the operand stack.
func (v *VM) push(val value.Value) {
	v.Stack = append(v.Stack, val)
}

// Pop a value from the operand stack.
func (v *VM) pop() value.Value {
	if len(v.Stack) == 0 {
		return value.Nil
	}
	val := v.Stack[len(v.Stack)-1]
	v.Stack = v.Stack[:len(v.Stack)-1]
	return val
}

// Peek at the stack top.
func (v *VM) peek() value.Value {
	if len(v.Stack) == 0 {
		return value.Nil
	}
	return v.Stack[len(v.Stack)-1]
}

// Execute runs the given program from its main entry point.
func (v *VM) Execute(prog *opcode.Program) error {
	v.Program = prog
	v.Halted = false
	v.Registry.Prog = prog
	v.Heap.SeedProgramStrings(prog.StringTable)
	defer func() { v.Registry.Prog = nil }()

	v.Registry.StackTraceFn = func() string { return v.FormatCallStack() }
	defer func() { v.Registry.StackTraceFn = nil }()

	v.Registry.TerminateVM = func() { v.Halted = true }
	defer func() { v.Registry.TerminateVM = nil }()

	// Push the <MAIN> chunk as our first call frame.
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

	if v.StackHygieneDebug && len(v.Stack) != 0 && !v.Halted {
		return fmt.Errorf("[moonBASIC] stack hygiene: program finished with %d value(s) on operand stack (expected 0 after statements; see ARCHITECTURE §6)", len(v.Stack))
	}

	return nil
}

// step processes one instruction.
// Sub-implementations are in vm_arith.go and vm_control.go.
func (v *VM) step(i opcode.Instruction) error {
	switch i.Op {
	case opcode.OpPushInt:
		frame := v.CallStack.Top()
		val := frame.Chunk.IntConsts[i.Operand]
		v.push(value.FromInt(val))

	case opcode.OpPushFloat:
		frame := v.CallStack.Top()
		val := frame.Chunk.FloatConsts[i.Operand]
		v.push(value.FromFloat(val))

	case opcode.OpPushString:
		idx := int(i.Operand)
		if idx < 0 || idx >= len(v.Program.StringTable) {
			return v.runtimeError("PUSH_STRING: string pool index out of range")
		}
		v.push(value.FromStringIndex(int32(idx)))

	case opcode.OpPushBool:
		v.push(value.FromBool(i.Operand != 0))

	case opcode.OpPushNull:
		v.push(value.Nil)

	case opcode.OpPop:
		v.pop()

	case opcode.OpLoadGlobal:
		frame := v.CallStack.Top()
		name := frame.Chunk.Names[i.Operand]
		if val, ok := v.Globals[name]; ok {
			v.push(val)
		} else {
			v.push(value.Nil)
		}

	case opcode.OpStoreGlobal:
		frame := v.CallStack.Top()
		name := frame.Chunk.Names[i.Operand]
		v.Globals[name] = v.peek()

	case opcode.OpLoadLocal:
		frame := v.CallStack.Top()
		slot := int(i.Operand)
		idx := frame.StackBase + slot
		if idx >= 0 && idx < len(v.Stack) {
			v.push(v.Stack[idx])
		} else {
			v.push(value.Nil)
		}

	case opcode.OpStoreLocal:
		frame := v.CallStack.Top()
		slot := int(i.Operand)
		idx := frame.StackBase + slot
		if idx >= 0 && idx < len(v.Stack) {
			v.Stack[idx] = v.peek()
		}

	case opcode.OpHalt:
		v.Halted = true

	case opcode.OpSwap:
		if len(v.Stack) < 2 {
			return v.runtimeError("SWAP: stack underflow")
		}
		n := len(v.Stack)
		v.Stack[n-1], v.Stack[n-2] = v.Stack[n-2], v.Stack[n-1]

	default:
		// Delegate to specialized handlers for complex logic
		return v.dispatchComplex(i)
	}

	return nil
}

// trace dumps IP, opcode, stack depth, and stack contents after each step.
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
	fmt.Fprintf(out, "[trace] %s L%d IP=%d %s | depth=%d stack=%v\n",
		chunk, line, ip, instr.String(), len(v.Stack), v.Stack)
}
