package vm

import (
	"fmt"

	"moonbasic/vm/heap"
	"moonbasic/vm/opcode"
	"moonbasic/vm/value"
)

func dimsFromLinear(dims []int64, li int64) []int64 {
	out := make([]int64, len(dims))
	rem := li
	for i := 0; i < len(dims); i++ {
		stride := int64(1)
		for j := i + 1; j < len(dims); j++ {
			stride *= dims[j]
		}
		out[i] = rem / stride
		rem %= stride
	}
	return out
}

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
	
	dims := make([]int64, nd)
	argStart := i.SrcA
	for j := 0; j < nd; j++ {
		dims[j] = v.regInt64(argStart + uint8(j))
	}
	kind := arrayKindFromFlags(i.SrcB)
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
	v.setReg(i.Dst, value.FromHandle(h))
	return nil
}

func (v *VM) regInt64(reg uint8) int64 {
	x := v.reg(reg)
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
	
	indices := make([]int64, nd)
	dimStart := i.SrcB
	for j := 0; j < nd; j++ {
		indices[j] = v.regInt64(dimStart + uint8(j))
	}

	hv := v.reg(i.SrcA)
	if hv.Kind != value.KindHandle {
		return v.runtimeError("ARRAYGET: not a handle")
	}
	obj, ok := v.Heap.Get(heap.Handle(hv.IVal))
	if !ok {
		return v.runtimeError("ARRAYGET: invalid handle")
	}
	arr, ok := obj.(*heap.Array)
	if !ok {
		return v.runtimeError(fmt.Sprintf("ARRAYGET: expected array, got %s", obj.TypeName()))
	}
	switch arr.Kind {
	case heap.ArrayKindHandle:
		hid, err := arr.GetHandle(indices)
		if err != nil {
			return v.runtimeError(err.Error())
		}
		v.setReg(i.Dst, value.FromHandle(heap.Handle(hid)))
	case heap.ArrayKindString:
		si, err := arr.GetStringIndex(indices)
		if err != nil {
			return v.runtimeError(err.Error())
		}
		v.setReg(i.Dst, value.FromStringIndex(si))
	case heap.ArrayKindBool:
		f, err := arr.GetFloat(indices)
		if err != nil {
			return v.runtimeError(err.Error())
		}
		v.setReg(i.Dst, value.FromBool(f != 0))
	default:
		f, err := arr.GetFloat(indices)
		if err != nil {
			return v.runtimeError(err.Error())
		}
		v.setReg(i.Dst, value.FromFloat(f))
	}
	return nil
}

func (v *VM) doArraySet(i opcode.Instruction) error {
	nd := int(i.Operand)
	if nd < 1 {
		return v.runtimeError("ARRAYSET: bad dimension count")
	}
	
	val := v.reg(i.Dst)
	indices := make([]int64, nd)
	dimStart := i.SrcB
	for j := 0; j < nd; j++ {
		indices[j] = v.regInt64(dimStart + uint8(j))
	}

	hv := v.reg(i.SrcA)
	if hv.Kind != value.KindHandle {
		return v.runtimeError("ARRAYSET: not a handle")
	}
	obj, ok := v.Heap.Get(heap.Handle(hv.IVal))
	if !ok {
		return v.runtimeError("ARRAYSET: invalid handle")
	}
	arr, ok := obj.(*heap.Array)
	if !ok {
		return v.runtimeError(fmt.Sprintf("ARRAYSET: expected array, got %s", obj.TypeName()))
	}
	switch arr.Kind {
	case heap.ArrayKindHandle:
		if val.Kind != value.KindHandle {
			return v.runtimeError("ARRAYSET: handle array expects handle value")
		}
		newH := int32(val.IVal)
		oldH, err := arr.GetHandle(indices)
		if err != nil {
			return v.runtimeError(err.Error())
		}
		if oldH != 0 && oldH != newH {
			_ = v.Heap.Free(heap.Handle(oldH))
		}
		if err := arr.SetHandle(indices, newH); err != nil {
			return v.runtimeError(err.Error())
		}
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
	
	dims := make([]int64, nd)
	dimStart := i.SrcB
	for j := 0; j < nd; j++ {
		dims[j] = v.regInt64(dimStart + uint8(j))
	}

	hv := v.reg(i.SrcA)
	if hv.Kind != value.KindHandle {
		return v.runtimeError("ARRAYREDIM: not a handle")
	}
	obj, ok := v.Heap.Get(heap.Handle(hv.IVal))
	if !ok {
		return v.runtimeError("ARRAYREDIM: invalid handle")
	}
	arr, ok := obj.(*heap.Array)
	if !ok {
		return v.runtimeError(fmt.Sprintf("ARRAYREDIM: expected array, got %s", obj.TypeName()))
	}
	preserve := i.Dst != 0 // use Dst as preserve flag
	if err := arr.Redim(dims, preserve); err != nil {
		return v.runtimeError(err.Error())
	}
	return nil
}

func (v *VM) doArrayMakeTyped(i opcode.Instruction) error {
	ch := v.CallStack.Top().Chunk
	nd := int(i.SrcB)
	typeName := ch.Names[i.Operand]
	if _, ok := v.Program.Types[typeName]; !ok {
		return v.runtimeError("ARRAY_MAKE_TYPED: unknown type " + typeName)
	}
	
	dims := make([]int64, nd)
	dimStart := i.SrcA
	for j := 0; j < nd; j++ {
		dims[j] = v.regInt64(dimStart + uint8(j))
	}
	arr, err := heap.NewArrayOfKind(dims, heap.ArrayKindHandle, 0)
	if err != nil {
		return v.runtimeError(err.Error())
	}
	n := arr.TotalElements()
	for li := 0; li < n; li++ {
		idxs := dimsFromLinear(arr.Dims, int64(li))
		inst := heap.NewInstance(typeName)
		hid, err2 := v.Heap.Alloc(inst)
		if err2 != nil {
			return v.runtimeError(err2.Error())
		}
		if err := arr.SetHandle(idxs, int32(hid)); err != nil {
			return v.runtimeError(err.Error())
		}
	}
	ah, err := v.Heap.Alloc(arr)
	if err != nil {
		return v.runtimeError(err.Error())
	}
	v.setReg(i.Dst, value.FromHandle(ah))
	return nil
}

func (v *VM) doNewFilled(i opcode.Instruction) error {
	ch := v.CallStack.Top().Chunk
	typeName := ch.Names[i.Operand]
	nf := int(i.SrcA) // Use SrcA for field count
	td, ok := v.Program.Types[typeName]
	if !ok {
		return v.runtimeError("NEW_FILLED: unknown type " + typeName)
	}
	if nf != len(td.Fields) {
		return v.runtimeError(fmt.Sprintf("NEW_FILLED: %s needs %d fields, got %d", typeName, len(td.Fields), nf))
	}
	
	vals := make([]value.Value, nf)
	valStart := i.SrcB
	for j := 0; j < nf; j++ {
		vals[j] = v.reg(valStart + uint8(j))
	}
	inst := heap.NewInstance(typeName)
	for j, fn := range td.Fields {
		inst.SetField(fn, vals[j])
	}
	hid, err := v.Heap.Alloc(inst)
	if err != nil {
		return v.runtimeError(err.Error())
	}
	v.setReg(i.Dst, value.FromHandle(hid))
	return nil
}
