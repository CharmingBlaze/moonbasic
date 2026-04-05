//go:build cgo

package scatter

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	terr "moonbasic/runtime/terrain"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerScatter(m *Module, r runtime.Registrar) {
	r.Register("SCATTER.CREATE", "scatter", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return sCreate(m, rt, args...) })
	r.Register("SCATTER.FREE", "scatter", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return sFree(m, rt, args...) })
	r.Register("SCATTER.APPLY", "scatter", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return sApply(m, rt, args...) })
	r.Register("SCATTER.DRAWALL", "scatter", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return sDrawAll(m, rt, args...) })
	r.Register("PROP.PLACE", "prop", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return pPlace(m, rt, args...) })
	r.Register("PROP.FREE", "prop", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return pFree(m, rt, args...) })
	r.Register("PROP.DRAWALL", "prop", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return pDrawAll(m, rt, args...) })
}

func sCreate(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil || len(args) != 1 {
		return value.Nil, fmt.Errorf("SCATTER.CREATE expects name$")
	}
	name, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	id, err := m.h.Alloc(&ScatterObject{Name: name, Seed: 1})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(int32(id)), nil
}

func sFree(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SCATTER.FREE expects handle")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	m.h.Free(heap.Handle(h))
	return value.Nil, nil
}

func sApply(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("SCATTER.APPLY expects scatter, terrain, density#")
	}
	hs, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	ht, err := rt.ArgHandle(args, 1)
	if err != nil {
		return value.Nil, err
	}
	den, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	o, err := heap.Cast[*ScatterObject](m.h, heap.Handle(hs))
	if err != nil {
		return value.Nil, err
	}
	n := int(den * 200)
	if n < 10 {
		n = 10
	}
	if n > 800 {
		n = 800
	}
	rng := rand.New(rand.NewSource(o.Seed))
	o.X = make([]float32, n)
	o.Y = make([]float32, n)
	o.Z = make([]float32, n)
	for i := 0; i < n; i++ {
		fx := rng.Float64() * 400
		fz := rng.Float64() * 400
		o.X[i] = float32(fx + 50)
		o.Z[i] = float32(fz + 50)
		o.Y[i] = terr.HeightWorldPublic(m.h, heap.Handle(ht), o.X[i], o.Z[i])
	}
	return value.Nil, nil
}

func sDrawAll(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SCATTER.DRAWALL expects scatter")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	o, err := heap.Cast[*ScatterObject](m.h, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	c2 := rl.NewColor(80, 160, 90, 255)
	for i := range o.X {
		rl.DrawSphere(rl.Vector3{X: o.X[i], Y: o.Y[i] + 0.5, Z: o.Z[i]}, 0.4, c2)
	}
	return value.Nil, nil
}

func pPlace(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("PROP.PLACE expects x#, y#, z# (model param reserved)")
	}
	x, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	y, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	z, err := rt.ArgFloat(args, 3)
	if err != nil {
		return value.Nil, err
	}
	id, err := m.h.Alloc(&PropObject{X: float32(x), Y: float32(y), Z: float32(z)})
	if err != nil {
		return value.Nil, err
	}
	_ = args[0]
	m.props = append(m.props, id)
	return value.FromHandle(int32(id)), nil
}

func pFree(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PROP.FREE expects handle")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	hh := heap.Handle(h)
	for i, id := range m.props {
		if id == hh {
			m.props = append(m.props[:i], m.props[i+1:]...)
			break
		}
	}
	m.h.Free(hh)
	return value.Nil, nil
}

func pDrawAll(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = args
	_ = rt
	col := rl.NewColor(200, 120, 60, 255)
	for _, id := range m.props {
		p, err := heap.Cast[*PropObject](m.h, id)
		if err != nil {
			continue
		}
		rl.DrawCube(rl.Vector3{X: p.X, Y: p.Y, Z: p.Z}, 1, 1.2, 1, col)
	}
	return value.Nil, nil
}
