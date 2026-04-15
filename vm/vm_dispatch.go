// Package vm implements the moonBASIC bytecode interpreter.
package vm

import (
	"fmt"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/opcode"
	"moonbasic/vm/value"
)

// dispatchComplex handles opcodes that require logic beyond simple stack push/pop.
func (v *VM) dispatchComplex(i opcode.Instruction) error {
	switch i.Op {
	// Binary Operations
	case opcode.OpAdd, opcode.OpSub, opcode.OpMul, opcode.OpDiv, opcode.OpMod, opcode.OpPow:
		return v.doArithmetic(i)
	case opcode.OpNeg:
		return v.doNegation(i)

	// Comparisons
	case opcode.OpEq, opcode.OpNeq, opcode.OpLt, opcode.OpGt, opcode.OpLte, opcode.OpGte:
		return v.doComparison(i)

	// Logical Operations
	case opcode.OpAnd, opcode.OpOr, opcode.OpNot, opcode.OpXor:
		return v.doLogic(i)

	// String Concat
	case opcode.OpConcat:
		return v.doConcat(i)

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
		h := v.reg(i.SrcA)
		if h.Kind == value.KindHandle {
			_ = v.Heap.Free(heap.Handle(h.IVal))
		}

	case opcode.OpEraseAll:
		return v.EraseAllHandles()

	case opcode.OpSyncPhysics:
		if err := v.RunSyncPhysicsOpcode(i); err != nil {
			return v.runtimeError(err.Error())
		}

	case opcode.OpArrayLen:
		return v.doArrayLen(i)

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
	
	case opcode.OpEntityPropGet:
		return v.doEntityPropGet(i)
	case opcode.OpEntityPropSet:
		return v.doEntityPropSet(i)

	default:
		return v.runtimeError(fmt.Sprintf("unknown or unimplemented opcode: %s", i.Op.String()))
	}
	return nil
}

func (v *VM) doArithmetic(i opcode.Instruction) error {
	right := v.reg(i.SrcB)
	left := v.reg(i.SrcA)
	pool := v.Program.StringTable
	h := v.Heap
	if left.Kind == value.KindString || right.Kind == value.KindString {
		if i.Op != opcode.OpAdd {
			return v.runtimeError("only addition is defined for strings")
		}
		s := value.StringAt(left, pool, h) + value.StringAt(right, pool, h)
		v.setReg(i.Dst, value.FromStringIndex(v.Heap.Intern(s)))
		return nil
	}
	var res value.Value
	var err error

	switch i.Op {
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
	v.setReg(i.Dst, res)
	return nil
}

func (v *VM) doNegation(i opcode.Instruction) error {
	val := v.reg(i.SrcA)
	res, err := value.Neg(val)
	if err != nil {
		return v.runtimeError(err.Error())
	}
	v.setReg(i.Dst, res)
	return nil
}

func (v *VM) doComparison(i opcode.Instruction) error {
	right := v.reg(i.SrcB)
	left := v.reg(i.SrcA)
	pool := v.Program.StringTable
	h := v.Heap

	switch i.Op {
	case opcode.OpEq:
		var eq bool
		if left.Kind == value.KindString && right.Kind == value.KindString {
			eq = value.EqualStringValue(left, right, pool, h)
		} else {
			eq = value.Equal(left, right)
		}
		v.setReg(i.Dst, value.FromBool(eq))
	case opcode.OpNeq:
		var eq bool
		if left.Kind == value.KindString && right.Kind == value.KindString {
			eq = value.EqualStringValue(left, right, pool, h)
		} else {
			eq = value.Equal(left, right)
		}
		v.setReg(i.Dst, value.FromBool(!eq))
	case opcode.OpLt:
		res, err := value.Less(left, right, pool, h)
		if err != nil {
			return v.runtimeError(err.Error())
		}
		v.setReg(i.Dst, value.FromBool(res))
	case opcode.OpGt:
		res, err := value.Less(right, left, pool, h)
		if err != nil {
			return v.runtimeError(err.Error())
		}
		v.setReg(i.Dst, value.FromBool(res))
	case opcode.OpLte:
		res, err := value.Less(right, left, pool, h)
		if err != nil {
			return v.runtimeError(err.Error())
		}
		v.setReg(i.Dst, value.FromBool(!res))
	case opcode.OpGte:
		res, err := value.Less(left, right, pool, h)
		if err != nil {
			return v.runtimeError(err.Error())
		}
		v.setReg(i.Dst, value.FromBool(!res))
	}
	return nil
}

func (v *VM) doLogic(i opcode.Instruction) error {
	pool := v.Program.StringTable
	h := v.Heap
	switch i.Op {
	case opcode.OpNot:
		v.setReg(i.Dst, value.FromBool(!value.Truthy(v.reg(i.SrcA), pool, h)))
	case opcode.OpAnd:
		l, r := v.reg(i.SrcA), v.reg(i.SrcB)
		v.setReg(i.Dst, value.FromBool(value.Truthy(l, pool, h) && value.Truthy(r, pool, h)))
	case opcode.OpOr:
		l, r := v.reg(i.SrcA), v.reg(i.SrcB)
		v.setReg(i.Dst, value.FromBool(value.Truthy(l, pool, h) || value.Truthy(r, pool, h)))
	case opcode.OpXor:
		l, r := v.reg(i.SrcA), v.reg(i.SrcB)
		v.setReg(i.Dst, value.FromBool(value.Truthy(l, pool, h) != value.Truthy(r, pool, h)))
	}
	return nil
}

func (v *VM) doConcat(i opcode.Instruction) error {
	l, r := v.reg(i.SrcA), v.reg(i.SrcB)
	pool := v.Program.StringTable
	h := v.Heap
	s := value.StringAt(l, pool, h) + value.StringAt(r, pool, h)
	v.setReg(i.Dst, value.FromStringIndex(v.Heap.Intern(s)))
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
	v.setReg(i.Dst, value.FromHandle(h))
	return nil
}

func (v *VM) doFieldGet(i opcode.Instruction) error {
	frame := v.CallStack.Top()
	fieldName := frame.Chunk.Names[i.Operand]

	hVal := v.reg(i.SrcA)
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

	if td, ok2 := v.Program.Types[strings.ToUpper(inst.Type)]; ok2 {
		found := false
		for _, f := range td.Fields {
			if f == fieldName {
				found = true
				break
			}
		}
		if !found {
			hint := ""
			if sug := suggestFieldName(fieldName, td.Fields); sug != "" {
				hint = fmt.Sprintf("\n  Hint: did you mean %q?", sug)
			}
			return v.runtimeError(fmt.Sprintf("%s has no field %q%s", inst.Type, fieldName, hint))
		}
	}

	v.setReg(i.Dst, inst.GetField(fieldName))
	return nil
}

func (v *VM) doFieldSet(i opcode.Instruction) error {
	frame := v.CallStack.Top()
	fieldName := frame.Chunk.Names[i.Operand]

	val := v.reg(i.SrcB)
	hVal := v.reg(i.SrcA)

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

	if td, ok2 := v.Program.Types[strings.ToUpper(inst.Type)]; ok2 {
		found := false
		for _, f := range td.Fields {
			if f == fieldName {
				found = true
				break
			}
		}
		if !found {
			hint := ""
			if sug := suggestFieldName(fieldName, td.Fields); sug != "" {
				hint = fmt.Sprintf("\n  Hint: did you mean %q?", sug)
			}
			return v.runtimeError(fmt.Sprintf("%s has no field %q%s", inst.Type, fieldName, hint))
		}
	}

	inst.SetField(fieldName, val)

	// Return the value (assignments as expressions return the value)
	v.setReg(i.Dst, val)
	return nil
}

func (v *VM) doCallHandle(i opcode.Instruction) error {
	frame := v.CallStack.Top()
	// Operand encodes (ArgCount << 24 | NameIdx)
	argCount := int(uint32(i.Operand) >> 24)
	nameIdx := int32(i.Operand & 0x00FFFFFF)
	methodName := frame.Chunk.Names[nameIdx]

	hVal := v.reg(i.SrcA)
	argStart := i.SrcB

	if hVal.Kind != value.KindHandle {
		return v.runtimeError(fmt.Sprintf("cannot call method %s on %s (not a handle)", methodName, hVal.TypeName()))
	}

	// Extract args from registers starting at ArgStart
	args := make([]value.Value, argCount)
	for idx := 0; idx < argCount; idx++ {
		args[idx] = v.reg(argStart + uint8(idx))
	}

	hid := heap.Handle(hVal.IVal)
	obj, ok := v.Heap.Get(hid)
	if !ok {
		if hid == 0 {
			return v.runtimeError("method " + methodName + " called on null handle (0)\n  Hint: Initialize the handle from MAKE/LOAD before calling methods.")
		}
		return v.runtimeError(fmt.Sprintf("method %s called with invalid handle %d\n  Hint: Object may have been freed; do not use handles after FREE.", methodName, hVal.IVal))
	}

	// ENTITYREF wraps an integer entity id; ENTITY.* builtins expect that id as the first argument, not a heap handle.
	firstReceiver := hVal
	if er, ok := obj.(*heap.EntityRef); ok {
		firstReceiver = value.FromInt(er.ID)
	}

	tag := obj.TypeTag()
	typeName := obj.TypeName()
	key, prepend, mapped := handleCallDispatch(tag, methodName, argCount)
	var callArgs []value.Value
	var callKey string
	if mapped {
		callKey = key
		if prepend {
			callArgs = make([]value.Value, 0, len(args)+1)
			callArgs = append(callArgs, firstReceiver)
			callArgs = append(callArgs, args...)
		} else {
			callArgs = args
		}
	} else {
		prefix := handleCallRegistryPrefix(tag)
		if prefix == "" {
			prefix = strings.ToUpper(strings.TrimSpace(typeName)) + "."
		}
		callKey = prefix + strings.ToUpper(strings.TrimSpace(methodName))
		callArgs = args
	}
	res, err := v.Registry.Call(callKey, callArgs)
	if err != nil {
		return v.runtimeError(v.formatHandleCallError(tag, typeName, methodName, callKey, mapped, err))
	}

	v.setReg(i.Dst, res)
	return nil
}

func suggestFieldName(want string, fields []string) string {
	var best string
	bestD := 100
	for _, f := range fields {
		d := levenshtein(want, f)
		if d < bestD {
			bestD, best = d, f
		}
	}
	if bestD <= 2 && best != "" {
		return best
	}
	return ""
}

func levenshtein(a, b string) int {
	if len(a) < len(b) {
		a, b = b, a
	}
	row := make([]int, len(b)+1)
	for j := 0; j <= len(b); j++ {
		row[j] = j
	}
	for i := 1; i <= len(a); i++ {
		prev := row[0]
		row[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			tmp := row[j]
			row[j] = min3(prev+cost, row[j]+1, row[j-1]+1)
			prev = tmp
		}
	}
	return row[len(b)]
}

func min3(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= c {
		return b
	}
	return c
}

func (v *VM) validateEntityMacroID(id int64) error {
	if id < 0 {
		return v.runtimeError(fmt.Sprintf("ENTITY: entity id %d is invalid (negative)", id))
	}
	if id >= runtime.MaxEntitySpatialIndex {
		return v.runtimeError(fmt.Sprintf("ENTITY: entity id %d exceeds maximum %d", id, runtime.MaxEntitySpatialIndex-1))
	}
	return nil
}

func (v *VM) doEntityPropGet(i opcode.Instruction) error {
	idVal := v.reg(i.SrcA)
	var id int64
	if idVal.Kind == value.KindHandle {
		id = idVal.IVal
		if obj, ok := v.Heap.Get(heap.Handle(idVal.IVal)); ok {
			if er, ok := obj.(*heap.EntityRef); ok {
				id = er.ID
			}
		}
	} else {
		// INTEGER arrays store as KindFloat; entity# lives in FVal, not IVal.
		var ok bool
		id, ok = idVal.ToInt()
		if !ok {
			return v.runtimeError("ENTITY_PROP_GET: entity id must be numeric")
		}
	}
	if err := v.validateEntityMacroID(id); err != nil {
		return err
	}

	// 1. Zero-Copy SoA Path
	sp := v.Registry.Spatial
	if sp != nil && id >= 0 && id < int64(len(sp.X)) {
		if v.Registry.EntityIDActive != nil && !v.Registry.EntityIDActive(id) {
			return v.runtimeError(fmt.Sprintf("ENTITY: no active entity with id %d (spatial read)", id))
		}
		var f float32
		switch i.Operand {
		case 0: f = sp.X[id]
		case 1: f = sp.Y[id]
		case 2: f = sp.Z[id]
		case 3: f = sp.P[id]
		case 4: f = sp.W[id]
		case 5: f = sp.R[id]
		default: goto fallback
		}
		v.setReg(i.Dst, value.FromFloat(float64(f)))
		return nil
	}

fallback:
	// 2. Legacy/Complex Property Path
	if v.Registry.FastEntityPropGet == nil {
		return v.runtimeError("ENTITY_PROP_GET: optimized spatial access not supported by active runtime")
	}
	res, err := v.Registry.FastEntityPropGet(id, int(i.Operand))
	if err != nil {
		return v.runtimeError(err.Error())
	}
	v.setReg(i.Dst, res)
	return nil
}

func (v *VM) doEntityPropSet(i opcode.Instruction) error {
	idVal := v.reg(i.SrcA)
	var id int64
	if idVal.Kind == value.KindHandle {
		id = idVal.IVal
		if obj, ok := v.Heap.Get(heap.Handle(idVal.IVal)); ok {
			if er, ok := obj.(*heap.EntityRef); ok {
				id = er.ID
			}
		}
	} else {
		var ok bool
		id, ok = idVal.ToInt()
		if !ok {
			return v.runtimeError("ENTITY_PROP_SET: entity id must be numeric")
		}
	}
	if err := v.validateEntityMacroID(id); err != nil {
		return err
	}
	val := v.reg(i.SrcB)

	// 1. Zero-Copy SoA Path
	sp := v.Registry.Spatial
	if sp != nil && id >= 0 && id < int64(len(sp.X)) {
		if v.Registry.EntityIDActive != nil && !v.Registry.EntityIDActive(id) {
			return v.runtimeError(fmt.Sprintf("ENTITY: no active entity with id %d (spatial write)", id))
		}
		fV, _ := val.ToFloat()
		f := float32(fV)
		switch i.Operand {
		case 0: sp.X[id] = f
		case 1: sp.Y[id] = f
		case 2: sp.Z[id] = f
		case 3: sp.P[id] = f
		case 4: sp.W[id] = f
		case 5: sp.R[id] = f
		default: goto fallback
		}
		return nil
	}

fallback:
	// 2. Legacy/Complex Property Path
	if v.Registry.FastEntityPropSet == nil {
		return v.runtimeError("ENTITY_PROP_SET: optimized spatial access not supported by active runtime")
	}
	err := v.Registry.FastEntityPropSet(id, int(i.Operand), val)
	if err != nil {
		return v.runtimeError(err.Error())
	}
	return nil
}
