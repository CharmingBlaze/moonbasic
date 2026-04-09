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
func (v *VM) CallUserFunction(name string, args []value.Value) (ret value.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = v.runtimeError(fmt.Sprintf("VM Panic Recovery: %v", r))
			if v.Registry != nil && err != nil {
				v.Registry.SetLastScriptError(err)
			}
		}
	}()

	if v.Program == nil {
		return value.Nil, fmt.Errorf("no program loaded")
	}
	key := strings.ToUpper(strings.TrimSpace(name))
	chunk, ok := v.Program.Functions[key]
	if !ok {
		err = fmt.Errorf("undefined function: %s", key)
		if v.Registry != nil {
			v.Registry.SetLastScriptError(err)
		}
		return value.Nil, err
	}

	baseDepth := v.CallStack.Depth()
	
	// We'll use register 255 of the CURRENT top frame to receive the return value.
	// If there are no frames, we have a problem. 
	// Usually CallUserFunction is called during execution or we can push a dummy frame.
	var retReg uint8 = 255
	
	// Push the new frame
	v.CallStack.Push(chunk, 0, retReg)
	newFrame := v.CallStack.Top()
	
	// Pass arguments in R0, R1...
	for idx, a := range args {
		if idx < 256 {
			newFrame.Registers[idx] = a
		}
	}

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
		if stepErr := v.step(instr); stepErr != nil {
			for v.CallStack.Depth() > baseDepth {
				v.CallStack.Pop()
			}
			err = stepErr
			if v.Registry != nil {
				v.Registry.SetLastScriptError(stepErr)
			}
			return value.Nil, stepErr
		}
	}

	// After the function returns, the result should be in the parent's retReg.
	ret = value.Nil
	parent := v.CallStack.Top()
	if parent != nil {
		ret = parent.Registers[retReg]
		// Clear it just in case
		parent.Registers[retReg] = value.Nil
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
