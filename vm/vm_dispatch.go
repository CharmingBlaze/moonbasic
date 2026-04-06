// Package vm implements the moonBASIC bytecode interpreter.
package vm

import (
	"fmt"
	"strings"

	"moonbasic/vm/heap"
	"moonbasic/vm/opcode"
	"moonbasic/vm/value"
)

// dispatchComplex handles opcodes that require logic beyond simple stack push/pop.
func (v *VM) dispatchComplex(i opcode.Instruction) error {
	switch i.Op {
	// Binary Operations
	case opcode.OpAdd, opcode.OpSub, opcode.OpMul, opcode.OpDiv, opcode.OpMod, opcode.OpPow:
		return v.doArithmetic(i.Op)
	case opcode.OpNeg:
		return v.doNegation()

	// Comparisons
	case opcode.OpEq, opcode.OpNeq, opcode.OpLt, opcode.OpGt, opcode.OpLte, opcode.OpGte:
		return v.doComparison(i.Op)

	// Logical Operations
	case opcode.OpAnd, opcode.OpOr, opcode.OpNot, opcode.OpXor:
		return v.doLogic(i.Op)

	// String Concat
	case opcode.OpConcat:
		return v.doConcat()

	// Control Flow (Jumps)
	case opcode.OpJump, opcode.OpJumpIfFalse, opcode.OpJumpIfTrue:
		return v.doJump(i)

	// Function Calls
	case opcode.OpCallBuiltin:
		return v.doCallBuiltin(i)
	case opcode.OpCallUser:
		return v.doCallUser(i)
	case opcode.OpReturn, opcode.OpReturnVoid:
		return v.doReturn(i)

	// Handles
	case opcode.OpNew:
		return v.doNew(i)
	case opcode.OpFieldGet:
		return v.doFieldGet(i)
	case opcode.OpFieldSet:
		return v.doFieldSet(i)
	case opcode.OpCallHandle:
		return v.doCallHandle(i)
	case opcode.OpDelete:
		h := v.pop()
		if h.Kind == value.KindHandle {
			_ = v.Heap.Free(heap.Handle(h.IVal))
		}

	case opcode.OpArrayMake:
		return v.doArrayMake(i)
	case opcode.OpArrayGet:
		return v.doArrayGet(i)
	case opcode.OpArraySet:
		return v.doArraySet(i)
	case opcode.OpArrayRedim:
		return v.doArrayRedim(i)
	case opcode.OpArrayMakeTyped:
		return v.doArrayMakeTyped(i)
	case opcode.OpNewFilled:
		return v.doNewFilled(i)

	default:
		return v.runtimeError(fmt.Sprintf("unknown or unimplemented opcode: %s", i.Op.String()))
	}
	return nil
}

func (v *VM) doArithmetic(op opcode.OpCode) error {
	right := v.pop()
	left := v.pop()
	pool := v.Program.StringTable
	h := v.Heap
	if left.Kind == value.KindString || right.Kind == value.KindString {
		if op != opcode.OpAdd {
			return v.runtimeError("only addition is defined for strings")
		}
		s := value.StringAt(left, pool, h) + value.StringAt(right, pool, h)
		v.push(value.FromStringIndex(v.Heap.Intern(s)))
		return nil
	}
	var res value.Value
	var err error

	switch op {
	case opcode.OpAdd:
		res, err = value.Add(left, right)
	case opcode.OpSub:
		res, err = value.Sub(left, right)
	case opcode.OpMul:
		res, err = value.Mul(left, right)
	case opcode.OpDiv:
		res, err = value.Div(left, right)
	case opcode.OpMod:
		res, err = value.Mod(left, right)
	case opcode.OpPow:
		res, err = value.Pow(left, right)
	}

	if err != nil {
		return v.runtimeError(err.Error())
	}
	v.push(res)
	return nil
}

func (v *VM) doNegation() error {
	val := v.pop()
	res, err := value.Neg(val)
	if err != nil {
		return v.runtimeError(err.Error())
	}
	v.push(res)
	return nil
}

func (v *VM) doComparison(op opcode.OpCode) error {
	right := v.pop()
	left := v.pop()
	pool := v.Program.StringTable
	h := v.Heap

	switch op {
	case opcode.OpEq:
		var eq bool
		if left.Kind == value.KindString && right.Kind == value.KindString {
			eq = value.EqualStringValue(left, right, pool, h)
		} else {
			eq = value.Equal(left, right)
		}
		v.push(value.FromBool(eq))
	case opcode.OpNeq:
		var eq bool
		if left.Kind == value.KindString && right.Kind == value.KindString {
			eq = value.EqualStringValue(left, right, pool, h)
		} else {
			eq = value.Equal(left, right)
		}
		v.push(value.FromBool(!eq))
	case opcode.OpLt:
		res, err := value.Less(left, right, pool, h)
		if err != nil {
			return v.runtimeError(err.Error())
		}
		v.push(value.FromBool(res))
	case opcode.OpGt:
		res, err := value.Less(right, left, pool, h)
		if err != nil {
			return v.runtimeError(err.Error())
		}
		v.push(value.FromBool(res))
	case opcode.OpLte:
		res, err := value.Less(right, left, pool, h)
		if err != nil {
			return v.runtimeError(err.Error())
		}
		v.push(value.FromBool(!res))
	case opcode.OpGte:
		res, err := value.Less(left, right, pool, h)
		if err != nil {
			return v.runtimeError(err.Error())
		}
		v.push(value.FromBool(!res))
	}
	return nil
}

func (v *VM) doLogic(op opcode.OpCode) error {
	pool := v.Program.StringTable
	h := v.Heap
	switch op {
	case opcode.OpNot:
		v.push(value.FromBool(!value.Truthy(v.pop(), pool, h)))
	case opcode.OpAnd:
		r, l := v.pop(), v.pop()
		v.push(value.FromBool(value.Truthy(l, pool, h) && value.Truthy(r, pool, h)))
	case opcode.OpOr:
		r, l := v.pop(), v.pop()
		v.push(value.FromBool(value.Truthy(l, pool, h) || value.Truthy(r, pool, h)))
	case opcode.OpXor:
		r, l := v.pop(), v.pop()
		v.push(value.FromBool(value.Truthy(l, pool, h) != value.Truthy(r, pool, h)))
	}
	return nil
}

func (v *VM) doConcat() error {
	r, l := v.pop(), v.pop()
	pool := v.Program.StringTable
	h := v.Heap
	s := value.StringAt(l, pool, h) + value.StringAt(r, pool, h)
	v.push(value.FromStringIndex(v.Heap.Intern(s)))
	return nil
}

// Milestone 6: Heap & Type logic

func (v *VM) doNew(i opcode.Instruction) error {
	frame := v.CallStack.Top()
	typeName := frame.Chunk.Names[i.Operand]

	// Check if type exists in program metadata
	if _, ok := v.Program.Types[typeName]; !ok {
		return v.runtimeError(fmt.Sprintf("NEW: unknown type %s", typeName))
	}

	// Create type instance
	obj := heap.NewInstance(typeName)
	h, err := v.Heap.Alloc(obj)
	if err != nil {
		return v.runtimeError(err.Error())
	}
	v.push(value.FromHandle(h))
	return nil
}

func (v *VM) doFieldGet(i opcode.Instruction) error {
	frame := v.CallStack.Top()
	fieldName := frame.Chunk.Names[i.Operand]

	hVal := v.pop()
	if hVal.Kind != value.KindHandle {
		return v.runtimeError(fmt.Sprintf("attempted to access field %s on %s (not a handle)", fieldName, hVal.TypeName()))
	}

	obj, ok := v.Heap.Get(heap.Handle(hVal.IVal))
	if !ok {
		return v.runtimeError(fmt.Sprintf("invalid handle %d for field %s", hVal.IVal, fieldName))
	}

	inst, ok := obj.(*heap.Instance)
	if !ok {
		return v.runtimeError(fmt.Sprintf("field %s exists only on user types, got %s", fieldName, obj.TypeName()))
	}

	v.push(inst.GetField(fieldName))
	return nil
}

func (v *VM) doFieldSet(i opcode.Instruction) error {
	frame := v.CallStack.Top()
	fieldName := frame.Chunk.Names[i.Operand]

	val := v.pop()
	hVal := v.pop()

	if hVal.Kind != value.KindHandle {
		return v.runtimeError(fmt.Sprintf("attempted to set field %s on %s (not a handle)", fieldName, hVal.TypeName()))
	}

	obj, ok := v.Heap.Get(heap.Handle(hVal.IVal))
	if !ok {
		return v.runtimeError(fmt.Sprintf("invalid handle %d for field set %s", hVal.IVal, fieldName))
	}

	inst, ok := obj.(*heap.Instance)
	if !ok {
		return v.runtimeError(fmt.Sprintf("cannot set field %s on engine type %s", fieldName, obj.TypeName()))
	}

	inst.SetField(fieldName, val)

	// Return the value (assignments as expressions return the value)
	v.push(val)
	return nil
}

func (v *VM) doCallHandle(i opcode.Instruction) error {
	frame := v.CallStack.Top()
	methodName := frame.Chunk.Names[i.Operand]
	argCount := int(i.Flags)

	// [handle] [arg1] [arg2] ... [argN]
	idx := len(v.Stack) - argCount - 1
	if idx < 0 {
		return v.runtimeError(fmt.Sprintf("stack underflow during handle call %s", methodName))
	}
	hVal := v.Stack[idx]

	if hVal.Kind != value.KindHandle {
		return v.runtimeError(fmt.Sprintf("cannot call method %s on %s (not a handle)", methodName, hVal.TypeName()))
	}

	// Extract args
	args := make([]value.Value, argCount)
	copy(args, v.Stack[len(v.Stack)-argCount:])
	v.Stack = v.Stack[:len(v.Stack)-argCount]

	// Pop the handle itself
	v.Stack = v.Stack[:len(v.Stack)-1]

	hid := heap.Handle(hVal.IVal)
	obj, ok := v.Heap.Get(hid)
	if !ok {
		if hid == 0 {
			return v.runtimeError("method " + methodName + " called on null handle (0)\n  Hint: Initialize the handle from MAKE/LOAD before calling methods.")
		}
		return v.runtimeError(fmt.Sprintf("method %s called with invalid handle %d\n  Hint: Object may have been freed; do not use handles after FREE.", methodName, hVal.IVal))
	}

	typeName := obj.TypeName()
	methUp := strings.ToUpper(strings.TrimSpace(methodName))
	key, prepend, mapped := handleCallBuiltin(typeName, methUp)
	var callArgs []value.Value
	var callKey string
	if mapped {
		callKey = key
		if prepend {
			callArgs = make([]value.Value, 0, len(args)+1)
			callArgs = append(callArgs, hVal)
			callArgs = append(callArgs, args...)
		} else {
			callArgs = args
		}
	} else {
		callKey = handleCallRegistryPrefix(typeName) + strings.ToUpper(strings.TrimSpace(methodName))
		callArgs = args
	}
	res, err := v.Registry.Call(callKey, callArgs)
	if err != nil {
		return v.runtimeError(v.formatHandleCallError(typeName, methodName, callKey, mapped, err))
	}

	v.push(res)
	return nil
}
