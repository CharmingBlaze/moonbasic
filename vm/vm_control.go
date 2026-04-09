// Package vm implements the moonBASIC bytecode interpreter.
package vm

import (
	"fmt"
	"moonbasic/vm/opcode"
	"moonbasic/vm/value"
)

// Control flow (Jumps and Branching)

func (v *VM) doJump(i opcode.Instruction) error {
	frame := v.CallStack.Top()
	switch i.Op {
	case opcode.OpJump:
		frame.IP = int(i.Operand)
	case opcode.OpJumpIfFalse:
		pool := v.Program.StringTable
		if !value.Truthy(v.reg(i.SrcA), pool, v.Heap) {
			frame.IP = int(i.Operand)
		}
	case opcode.OpJumpIfTrue:
		pool := v.Program.StringTable
		if value.Truthy(v.reg(i.SrcA), pool, v.Heap) {
			frame.IP = int(i.Operand)
		}
	}
	return nil
}

// Function Calls (User and Built-in)

func (v *VM) doCallBuiltin(i opcode.Instruction) error {
	frame := v.CallStack.Top()
	// Operand encodes (ArgCount << 24 | NameIdx)
	argCount := int(uint32(i.Operand) >> 24)
	nameIdx := int32(i.Operand & 0x00FFFFFF)
	name := frame.Chunk.Names[nameIdx]

	argStart := i.SrcB

	// Extract args from registers starting at argStart
	args := make([]value.Value, argCount)
	for idx := 0; idx < argCount; idx++ {
		args[idx] = v.reg(argStart + uint8(idx))
	}

	// Call the native registry
	res, err := v.Registry.Call(name, args)
	if err != nil {
		return v.runtimeError(err.Error())
	}

	v.setReg(i.Dst, res)
	return nil
}

func (v *VM) doCallUser(i opcode.Instruction) error {
	frame := v.CallStack.Top()
	// Operand encodes (ArgCount << 24 | NameIdx)
	argCount := int(uint32(i.Operand) >> 24)
	nameIdx := int32(i.Operand & 0x00FFFFFF)
	name := frame.Chunk.Names[nameIdx]

	targetChunk, ok := v.Program.Functions[name]
	if !ok {
		return v.runtimeError(fmt.Sprintf("undefined function: %s", name))
	}

	// Arguments are in R[SrcB]...R[SrcB+argCount-1] in CURRENT frame
	args := make([]value.Value, argCount)
	argStart := i.SrcB
	for idx := 0; idx < argCount; idx++ {
		args[idx] = v.reg(argStart + uint8(idx))
	}

	// Push the new frame, saving i.Dst as the register to receive the result
	v.CallStack.Push(targetChunk, 0, i.Dst)
	newFrame := v.CallStack.Top()

	// Parameters go into R0, R1, ... in the NEW frame
	for idx := 0; idx < argCount; idx++ {
		newFrame.Registers[idx] = args[idx]
	}
	
	return nil
}

func (v *VM) doReturn(i opcode.Instruction) error {
	hasValue := i.Op == opcode.OpReturn

	var res value.Value
	if hasValue {
		res = v.reg(i.SrcA)
	}

	// Exit the current frame
	oldFrame := v.CallStack.Pop()
	
	// If there's a caller frame, set the return register
	parent := v.CallStack.Top()
	if parent != nil {
		if hasValue {
			parent.Registers[oldFrame.ReturnReg] = res
		} else {
			parent.Registers[oldFrame.ReturnReg] = value.Nil
		}
	}

	return nil
}
