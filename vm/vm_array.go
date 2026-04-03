package vm

import (
	"fmt"

	"moonbasic/vm/heap"
	"moonbasic/vm/opcode"
	"moonbasic/vm/value"
)

func arrayKindFromFlags(f uint8) heap.ArrayKind {
	switch f {
	case 1:
		return heap.ArrayKindString
	case 2:
		return heap.ArrayKindBool
	default:
		return heap.ArrayKindFloat
	}
}

func (v *VM) doArrayMake(i opcode.Instruction) error {
	nd := int(i.Operand)
	if nd < 1 {
		return v.runtimeError("ARRAYMAKE: need at least 1 dimension")
	}
	if len(v.Stack) < nd {
		return v.runtimeError("ARRAYMAKE: stack underflow")
	}
	dims := make([]int64, nd)
	for j := nd - 1; j >= 0; j-- {
		dims[j] = v.popInt64ForDim()
	}
	kind := arrayKindFromFlags(i.Flags)
	emptyStr := int32(0)
	if kind == heap.ArrayKindString {
		emptyStr = v.Heap.Intern("")
	}
	arr, err := heap.NewArrayOfKind(dims, kind, emptyStr)
	if err != nil {
		return v.runtimeError(err.Error())
	}
	h, err := v.Heap.Alloc(arr)
	if err != nil {
		return v.runtimeError(err.Error())
	}
	v.push(value.FromHandle(h))
	return nil
}

func (v *VM) popInt64ForDim() int64 {
	x := v.pop()
	if i, ok := x.ToInt(); ok {
		return i
	}
	if f, ok := x.ToFloat(); ok {
		return int64(f)
	}
	return 0
}

func (v *VM) doArrayGet(i opcode.Instruction) error {
	nd := int(i.Operand)
	if nd < 1 {
		return v.runtimeError("ARRAYGET: bad dimension count")
	}
	need := nd + 1
	if len(v.Stack) < need {
		return v.runtimeError("ARRAYGET: stack underflow")
	}
	indices := make([]int64, nd)
	for j := nd - 1; j >= 0; j-- {
		indices[j] = v.popInt64ForDim()
	}
	hv := v.pop()
	if hv.Kind != value.KindHandle {
		return v.runtimeError("ARRAYGET: not a handle")
	}
	obj, ok := v.Heap.Get(int32(hv.IVal))
	if !ok {
		return v.runtimeError("ARRAYGET: invalid handle")
	}
	arr, ok := obj.(*heap.Array)
	if !ok {
		return v.runtimeError(fmt.Sprintf("ARRAYGET: expected array, got %s", obj.TypeName()))
	}
	switch arr.Kind {
	case heap.ArrayKindString:
		si, err := arr.GetStringIndex(indices)
		if err != nil {
			return v.runtimeError(err.Error())
		}
		v.push(value.FromStringIndex(si))
	case heap.ArrayKindBool:
		f, err := arr.GetFloat(indices)
		if err != nil {
			return v.runtimeError(err.Error())
		}
		v.push(value.FromBool(f != 0))
	default:
		f, err := arr.GetFloat(indices)
		if err != nil {
			return v.runtimeError(err.Error())
		}
		v.push(value.FromFloat(f))
	}
	return nil
}

func (v *VM) doArraySet(i opcode.Instruction) error {
	nd := int(i.Operand)
	if nd < 1 {
		return v.runtimeError("ARRAYSET: bad dimension count")
	}
	need := nd + 2
	if len(v.Stack) < need {
		return v.runtimeError("ARRAYSET: stack underflow")
	}
	val := v.pop()
	indices := make([]int64, nd)
	for j := nd - 1; j >= 0; j-- {
		indices[j] = v.popInt64ForDim()
	}
	hv := v.pop()
	if hv.Kind != value.KindHandle {
		return v.runtimeError("ARRAYSET: not a handle")
	}
	obj, ok := v.Heap.Get(int32(hv.IVal))
	if !ok {
		return v.runtimeError("ARRAYSET: invalid handle")
	}
	arr, ok := obj.(*heap.Array)
	if !ok {
		return v.runtimeError(fmt.Sprintf("ARRAYSET: expected array, got %s", obj.TypeName()))
	}
	switch arr.Kind {
	case heap.ArrayKindString:
		if val.Kind != value.KindString {
			return v.runtimeError("ARRAYSET: string array expects string value")
		}
		if err := arr.SetStringIndex(indices, int32(val.IVal)); err != nil {
			return v.runtimeError(err.Error())
		}
	case heap.ArrayKindBool:
		b := value.Truthy(val, v.Program.StringTable, v.Heap)
		f := 0.0
		if b {
			f = 1
		}
		if err := arr.SetFloat(indices, f); err != nil {
			return v.runtimeError(err.Error())
		}
	default:
		f, okf := val.ToFloat()
		if !okf {
			if val.Kind == value.KindInt {
				f = float64(val.IVal)
				okf = true
			}
		}
		if !okf {
			return v.runtimeError("ARRAYSET: numeric value expected")
		}
		if err := arr.SetFloat(indices, f); err != nil {
			return v.runtimeError(err.Error())
		}
	}
	return nil
}

func (v *VM) doArrayRedim(i opcode.Instruction) error {
	nd := int(i.Operand)
	if nd < 1 {
		return v.runtimeError("ARRAYREDIM: bad dimension count")
	}
	if len(v.Stack) < nd+1 {
		return v.runtimeError("ARRAYREDIM: stack underflow")
	}
	dims := make([]int64, nd)
	for j := nd - 1; j >= 0; j-- {
		dims[j] = v.popInt64ForDim()
	}
	hv := v.pop()
	if hv.Kind != value.KindHandle {
		return v.runtimeError("ARRAYREDIM: not a handle")
	}
	obj, ok := v.Heap.Get(int32(hv.IVal))
	if !ok {
		return v.runtimeError("ARRAYREDIM: invalid handle")
	}
	arr, ok := obj.(*heap.Array)
	if !ok {
		return v.runtimeError(fmt.Sprintf("ARRAYREDIM: expected array, got %s", obj.TypeName()))
	}
	preserve := i.Flags != 0
	if err := arr.Redim(dims, preserve); err != nil {
		return v.runtimeError(err.Error())
	}
	return nil
}
