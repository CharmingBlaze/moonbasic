package vm

import (
	"fmt"
	"strings"

	"moonbasic/vm/value"
)

// CallUserFunction invokes a user-defined FUNCTION by name with the given arguments.
// It saves and restores the operand stack and call stack depth so it can be used from
// native code while the VM is executing (e.g. physics collision processing).
// The function name is matched case-insensitively against Program.Functions keys.
func (v *VM) CallUserFunction(name string, args []value.Value) (value.Value, error) {
	if v.Program == nil {
		return value.Nil, fmt.Errorf("no program loaded")
	}
	key := strings.ToUpper(strings.TrimSpace(name))
	chunk, ok := v.Program.Functions[key]
	if !ok {
		return value.Nil, fmt.Errorf("undefined function: %s", key)
	}

	baseDepth := v.CallStack.Depth()
	baseLen := len(v.Stack)

	for _, a := range args {
		v.push(a)
	}
	v.CallStack.Push(chunk, 0, len(v.Stack)-len(args))
	targetDepth := baseDepth + 1

	for v.CallStack.Depth() >= targetDepth {
		frame := v.CallStack.Top()
		if frame == nil {
			break
		}
		if frame.IP >= len(frame.Chunk.Instructions) {
			v.CallStack.Pop()
			continue
		}
		instr := frame.Chunk.Instructions[frame.IP]
		frame.IP++
		if err := v.step(instr); err != nil {
			v.Stack = v.Stack[:baseLen]
			for v.CallStack.Depth() > baseDepth {
				v.CallStack.Pop()
			}
			return value.Nil, err
		}
	}

	var ret value.Value
	if len(v.Stack) > baseLen {
		ret = v.Stack[len(v.Stack)-1]
		v.Stack = v.Stack[:baseLen]
	} else {
		ret = value.Nil
	}
	for v.CallStack.Depth() > baseDepth {
		v.CallStack.Pop()
	}
	return ret, nil
}

// ProgramLoaded reports whether Execute has been called (or Program assigned).
func (v *VM) ProgramLoaded() bool {
	return v.Program != nil
}
