//go:build !cgo && windows

package mbgui

import (
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func stringFromRT(rt *runtime.Runtime, v value.Value) string {
	var pool []string
	var hg value.StringGetter
	if rt != nil {
		if rt.Prog != nil {
			pool = rt.Prog.StringTable
		}
		if rt.Heap != nil {
			hg = rt.Heap
		}
	}
	return value.StringAt(v, pool, hg)
}

func argF32(v value.Value) (float32, bool) {
	if f, ok := v.ToFloat(); ok {
		return float32(f), true
	}
	if i, ok := v.ToInt(); ok {
		return float32(i), true
	}
	return 0, false
}

func argI32(v value.Value) (int32, bool) {
	if i, ok := v.ToInt(); ok {
		return int32(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return int32(f), true
	}
	return 0, false
}

func rectArgs(args []value.Value, o int) (rl.Rectangle, error) {
	if o+3 >= len(args) {
		return rl.Rectangle{}, fmt.Errorf("need rectangle (x,y,w,h) at arg %d", o+1)
	}
	x, okx := argF32(args[o])
	y, oky := argF32(args[o+1])
	w, okw := argF32(args[o+2])
	h, okh := argF32(args[o+3])
	if !okx || !oky || !okw || !okh {
		return rl.Rectangle{}, fmt.Errorf("rectangle components must be numeric")
	}
	return rl.Rectangle{X: x, Y: y, Width: w, Height: h}, nil
}

func colorArgs(args []value.Value, o int) (rl.Color, error) {
	if o+3 >= len(args) {
		return rl.Color{}, fmt.Errorf("need RGBA at arg %d", o+1)
	}
	r, okr := argI32(args[o])
	g, okg := argI32(args[o+1])
	b, okb := argI32(args[o+2])
	a, oka := argI32(args[o+3])
	if !okr || !okg || !okb || !oka {
		return rl.Color{}, fmt.Errorf("color components must be numeric")
	}
	return rl.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}, nil
}

func splitItems(s string) []string {
	parts := strings.Split(s, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return []string{""}
	}
	return out
}

func allocRGBA(m *Module, c rl.Color) (value.Value, error) {
	arr, err := heap.NewArray([]int64{4})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, float64(c.R))
	_ = arr.Set([]int64{1}, float64(c.G))
	_ = arr.Set([]int64{2}, float64(c.B))
	_ = arr.Set([]int64{3}, float64(c.A))
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func allocRect(m *Module, r rl.Rectangle) (value.Value, error) {
	arr, err := heap.NewArray([]int64{4})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, float64(r.X))
	_ = arr.Set([]int64{1}, float64(r.Y))
	_ = arr.Set([]int64{2}, float64(r.Width))
	_ = arr.Set([]int64{3}, float64(r.Height))
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func readFloat2(store *heap.Store, h heap.Handle) (float64, float64, error) {
	a, err := heap.Cast[*heap.Array](store, h)
	if err != nil {
		return 0, 0, err
	}
	if a.TotalElements() < 2 {
		return 0, 0, fmt.Errorf("array needs at least 2 elements")
	}
	x, e1 := a.Get([]int64{0})
	y, e2 := a.Get([]int64{1})
	if e1 != nil {
		return 0, 0, e1
	}
	if e2 != nil {
		return 0, 0, e2
	}
	return x, y, nil
}

func writeFloat2(store *heap.Store, h heap.Handle, x, y float64) error {
	a, err := heap.Cast[*heap.Array](store, h)
	if err != nil {
		return err
	}
	if a.TotalElements() < 2 {
		return fmt.Errorf("array needs at least 2 elements")
	}
	if e := a.Set([]int64{0}, x); e != nil {
		return e
	}
	return a.Set([]int64{1}, y)
}

func boolAsFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
