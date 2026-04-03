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
		if !value.Truthy(v.pop(), pool) {
			frame.IP = int(i.Operand)
		}
	case opcode.OpJumpIfTrue:
		pool := v.Program.StringTable
		if value.Truthy(v.pop(), pool) {
			frame.IP = int(i.Operand)
		}
	}
	return nil
}

// Function Calls (User and Built-in)

func (v *VM) doCallBuiltin(i opcode.Instruction) error {
	frame := v.CallStack.Top()
	name := frame.Chunk.Names[i.Operand]
	argCount := int(i.Flags)

	if len(v.Stack) < argCount {
		return v.runtimeError(fmt.Sprintf("not enough arguments for %s", name))
	}

	// Extract args from stack
	args := make([]value.Value, argCount)
	copy(args, v.Stack[len(v.Stack)-argCount:])
	v.Stack = v.Stack[:len(v.Stack)-argCount]

	// Call the native registry
	res, err := v.Registry.Call(name, args)
	if err != nil {
		return v.runtimeError(err.Error())
	}

	v.push(res)
	return nil
}

func (v *VM) doCallUser(i opcode.Instruction) error {
	frame := v.CallStack.Top()
	name := frame.Chunk.Names[i.Operand]
	argCount := int(i.Flags)

	targetChunk, ok := v.Program.Functions[name]
	if !ok {
		return v.runtimeError(fmt.Sprintf("undefined function: %s", name))
	}

	// The frame's stack base is the first argument
	newBase := len(v.Stack) - argCount
	v.CallStack.Push(targetChunk, 0, newBase)
	
	return nil
}

func (v *VM) doReturn(i opcode.Instruction) error {
	hasValue := i.Op == opcode.OpReturn

	var res value.Value
	if hasValue {
		res = v.pop()
	}

	// Exit the current frame
	oldFrame := v.CallStack.Pop()
	
	// Truncate the stack back to where this frame began (dropping locals/params)
	v.Stack = v.Stack[:oldFrame.StackBase]

	if hasValue {
		v.push(res)
	} else {
		v.push(value.Nil)
	}

	return nil
}
