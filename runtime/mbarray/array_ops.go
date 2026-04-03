package mbarray

import (
	"fmt"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerAll(m *Module, r runtime.Registrar) {
	r.Register("ARRAYLEN", "array", m.arrayLen)
	r.Register("ARRAYFILL", "array", m.arrayFill)
	r.Register("ARRAYCOPY", "array", m.arrayCopy)
	r.Register("ARRAYSORT", "array", m.arraySort)
	r.Register("ARRAYREVERSE", "array", m.arrayReverse)
	r.Register("ARRAYFIND", "array", m.arrayFind)
	r.Register("ARRAYCONTAINS", "array", m.arrayContains)
	r.Register("ARRAYPUSH", "array", m.arrayPush)
	r.Register("ARRAYPOP", "array", m.arrayPop)
	r.Register("ARRAYSHIFT", "array", m.arrayShift)
	r.Register("ARRAYUNSHIFT", "array", m.arrayUnshift)
	r.Register("ARRAYSPLICE", "array", m.arraySplice)
	r.Register("ARRAYSLICE", "array", m.arraySlice)
	r.Register("ARRAYJOINS", "array", m.arrayJoins)
	r.Register("ARRAYJOINS$", "array", m.arrayJoins)
	r.Register("ARRAYFREE", "array", m.arrayFree)
	r.Register("ERASE", "array", m.arrayFree)
}

func (m *Module) getArr(v value.Value, op string) (*heap.Array, error) {
	if m.h == nil {
		return nil, runtime.Errorf("%s: heap not bound", op)
	}
	if v.Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: expected array handle", op)
	}
	return heap.Cast[*heap.Array](m.h, heap.Handle(v.IVal))
}

func (m *Module) arrayLen(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return value.Nil, fmt.Errorf("ARRAYLEN expects 1 or 2 arguments")
	}
	a, err := m.getArr(args[0], "ARRAYLEN")
	if err != nil {
		return value.Nil, err
	}
	if len(args) == 1 {
		return value.FromInt(int64(a.TotalElements())), nil
	}
	d1, ok := args[1].ToInt()
	if !ok {
		if f, okf := args[1].ToFloat(); okf {
			d1 = int64(f)
			ok = true
		}
	}
	if !ok {
		return value.Nil, fmt.Errorf("ARRAYLEN: dimension must be numeric")
	}
	sz, err := a.DimSize(int(d1))
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(sz), nil
}

func (m *Module) arrayFill(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ARRAYFILL expects (array, value)")
	}
	a, err := m.getArr(args[0], "ARRAYFILL")
	if err != nil {
		return value.Nil, err
	}
	switch a.Kind {
	case heap.ArrayKindString:
		if args[1].Kind != value.KindString {
			return value.Nil, fmt.Errorf("ARRAYFILL: string array expects string value")
		}
		if err := a.FillStringIndex(int32(args[1].IVal)); err != nil {
			return value.Nil, err
		}
	case heap.ArrayKindBool:
		b := value.Truthy(args[1], rt.Prog.StringTable)
		f := 0.0
		if b {
			f = 1
		}
		if err := a.FillScalar(f); err != nil {
			return value.Nil, err
		}
	default:
		f, ok := args[1].ToFloat()
		if !ok {
			if args[1].Kind == value.KindInt {
				f = float64(args[1].IVal)
				ok = true
			}
		}
		if !ok {
			return value.Nil, fmt.Errorf("ARRAYFILL: numeric value expected")
		}
		if err := a.FillScalar(f); err != nil {
			return value.Nil, err
		}
	}
	return value.Nil, nil
}

func (m *Module) arrayCopy(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ARRAYCOPY expects (src, dst)")
	}
	src, err := m.getArr(args[0], "ARRAYCOPY")
	if err != nil {
		return value.Nil, err
	}
	dst, err := m.getArr(args[1], "ARRAYCOPY")
	if err != nil {
		return value.Nil, err
	}
	if err := dst.CopyFrom(src); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) arraySort(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return value.Nil, fmt.Errorf("ARRAYSORT expects 1 or 2 arguments")
	}
	a, err := m.getArr(args[0], "ARRAYSORT")
	if err != nil {
		return value.Nil, err
	}
	desc := false
	if len(args) == 2 && args[1].Kind == value.KindString {
		s := strings.ToUpper(rt.Prog.StringTable[args[1].IVal])
		desc = s == "DESC" || s == "DESCENDING"
	}
	pool := rt.Prog.StringTable
	if a.Kind == heap.ArrayKindString {
		less := func(i, j int32) bool {
			si := ""
			sj := ""
			if int(i) >= 0 && int(i) < len(pool) {
				si = pool[i]
			}
			if int(j) >= 0 && int(j) < len(pool) {
				sj = pool[j]
			}
			return si < sj
		}
		if err := a.Sort1D(desc, less); err != nil {
			return value.Nil, err
		}
	} else {
		if err := a.Sort1D(desc, nil); err != nil {
			return value.Nil, err
		}
	}
	return value.Nil, nil
}

func (m *Module) arrayReverse(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ARRAYREVERSE expects array handle")
	}
	a, err := m.getArr(args[0], "ARRAYREVERSE")
	if err != nil {
		return value.Nil, err
	}
	a.Reverse()
	return value.Nil, nil
}

func (m *Module) arrayFind(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ARRAYFIND expects (array, value)")
	}
	a, err := m.getArr(args[0], "ARRAYFIND")
	if err != nil {
		return value.Nil, err
	}
	switch a.Kind {
	case heap.ArrayKindString:
		if args[1].Kind != value.KindString {
			return value.FromInt(-1), nil
		}
		want := int32(args[1].IVal)
		fi := a.FindStringIndex(want)
		return value.FromInt(int64(fi)), nil
	case heap.ArrayKindBool:
		want := value.Truthy(args[1], rt.Prog.StringTable)
		wf := 0.0
		if want {
			wf = 1
		}
		fi := a.FindFlat(wf)
		return value.FromInt(int64(fi)), nil
	default:
		f, ok := args[1].ToFloat()
		if !ok && args[1].Kind == value.KindInt {
			f = float64(args[1].IVal)
			ok = true
		}
		if !ok {
			return value.FromInt(-1), nil
		}
		fi := a.FindFlat(f)
		return value.FromInt(int64(fi)), nil
	}
}

func (m *Module) arrayContains(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	v, err := m.arrayFind(rt, args...)
	if err != nil {
		return value.Nil, err
	}
	i, _ := v.ToInt()
	return value.FromBool(i >= 0), nil
}

func (m *Module) arrayPush(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ARRAYPUSH expects (array, value)")
	}
	a, err := m.getArr(args[0], "ARRAYPUSH")
	if err != nil {
		return value.Nil, err
	}
	switch a.Kind {
	case heap.ArrayKindString:
		if args[1].Kind != value.KindString {
			return value.Nil, fmt.Errorf("ARRAYPUSH: string array expects string")
		}
		return value.Nil, a.Push1D(0, int32(args[1].IVal), false)
	case heap.ArrayKindBool:
		b := value.Truthy(args[1], rt.Prog.StringTable)
		f := 0.0
		if b {
			f = 1
		}
		return value.Nil, a.Push1D(f, 0, true)
	default:
		f, ok := args[1].ToFloat()
		if !ok && args[1].Kind == value.KindInt {
			f = float64(args[1].IVal)
			ok = true
		}
		if !ok {
			return value.Nil, fmt.Errorf("ARRAYPUSH: numeric value expected")
		}
		return value.Nil, a.Push1D(f, 0, false)
	}
}

func (m *Module) arrayPop(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ARRAYPOP expects array handle")
	}
	a, err := m.getArr(args[0], "ARRAYPOP")
	if err != nil {
		return value.Nil, err
	}
	f, si, ok := a.Pop1D()
	if !ok {
		return value.Nil, fmt.Errorf("ARRAYPOP: empty array")
	}
	switch a.Kind {
	case heap.ArrayKindString:
		return value.FromStringIndex(si), nil
	case heap.ArrayKindBool:
		return value.FromBool(f != 0), nil
	default:
		return value.FromFloat(f), nil
	}
}

func (m *Module) arrayShift(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ARRAYSHIFT expects array handle")
	}
	a, err := m.getArr(args[0], "ARRAYSHIFT")
	if err != nil {
		return value.Nil, err
	}
	f, si, ok := a.Shift1D()
	if !ok {
		return value.Nil, fmt.Errorf("ARRAYSHIFT: empty array")
	}
	switch a.Kind {
	case heap.ArrayKindString:
		return value.FromStringIndex(si), nil
	case heap.ArrayKindBool:
		return value.FromBool(f != 0), nil
	default:
		return value.FromFloat(f), nil
	}
}

func (m *Module) arrayUnshift(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ARRAYUNSHIFT expects (array, value)")
	}
	a, err := m.getArr(args[0], "ARRAYUNSHIFT")
	if err != nil {
		return value.Nil, err
	}
	switch a.Kind {
	case heap.ArrayKindString:
		if args[1].Kind != value.KindString {
			return value.Nil, fmt.Errorf("ARRAYUNSHIFT: string array expects string")
		}
		return value.Nil, a.Unshift1D(0, int32(args[1].IVal))
	case heap.ArrayKindBool:
		b := value.Truthy(args[1], rt.Prog.StringTable)
		f := 0.0
		if b {
			f = 1
		}
		return value.Nil, a.Unshift1D(f, 0)
	default:
		f, ok := args[1].ToFloat()
		if !ok && args[1].Kind == value.KindInt {
			f = float64(args[1].IVal)
			ok = true
		}
		if !ok {
			return value.Nil, fmt.Errorf("ARRAYUNSHIFT: numeric value expected")
		}
		return value.Nil, a.Unshift1D(f, 0)
	}
}

func (m *Module) arraySplice(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ARRAYSPLICE expects (array, pos, count)")
	}
	a, err := m.getArr(args[0], "ARRAYSPLICE")
	if err != nil {
		return value.Nil, err
	}
	pos, ok := args[1].ToInt()
	if !ok {
		if f, okf := args[1].ToFloat(); okf {
			pos = int64(f)
			ok = true
		}
	}
	if !ok {
		return value.Nil, fmt.Errorf("ARRAYSPLICE: pos must be numeric")
	}
	cnt, ok := args[2].ToInt()
	if !ok {
		if f, okf := args[2].ToFloat(); okf {
			cnt = int64(f)
			ok = true
		}
	}
	if !ok {
		return value.Nil, fmt.Errorf("ARRAYSPLICE: count must be numeric")
	}
	return value.Nil, a.Splice1D(pos, cnt)
}

func (m *Module) arraySlice(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ARRAYSLICE expects (array, start, end)")
	}
	a, err := m.getArr(args[0], "ARRAYSLICE")
	if err != nil {
		return value.Nil, err
	}
	start, ok := args[1].ToInt()
	if !ok {
		if f, okf := args[1].ToFloat(); okf {
			start = int64(f)
			ok = true
		}
	}
	end, ok2 := args[2].ToInt()
	if !ok2 {
		if f, okf := args[2].ToFloat(); okf {
			end = int64(f)
			ok2 = true
		}
	}
	if !ok || !ok2 {
		return value.Nil, fmt.Errorf("ARRAYSLICE: start and end must be numeric")
	}
	empty := rt.Prog.InternString("")
	out, err := a.Slice1D(start, end, empty)
	if err != nil {
		return value.Nil, err
	}
	h, err := m.h.Alloc(out)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func (m *Module) arrayJoins(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ARRAYJOINS expects (array, delimiter$)")
	}
	a, err := m.getArr(args[0], "ARRAYJOINS")
	if err != nil {
		return value.Nil, err
	}
	if a.Kind != heap.ArrayKindString {
		return value.Nil, fmt.Errorf("ARRAYJOINS: string array expected")
	}
	delim, _ := rt.ArgString(args, 1)
	s := a.JoinStrings(rt.Prog.StringTable, delim)
	return rt.RetString(s), nil
}

func (m *Module) arrayFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ARRAYFREE expects array handle")
	}
	if m.h == nil {
		return value.Nil, runtime.Errorf("ARRAYFREE: heap not bound")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("ARRAYFREE: expected handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}
