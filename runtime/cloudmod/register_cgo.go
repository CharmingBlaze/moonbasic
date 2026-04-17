//go:build cgo || (windows && !cgo)

package cloudmod

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerCloud(m *Module, r runtime.Registrar) {
	r.Register("CLOUD.CREATE", "cloud", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return cMake(m, rt, args...) })
	r.Register("CLOUD.MAKE", "cloud", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return cMake(m, rt, args...) })
	r.Register("CLOUD.FREE", "cloud", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return cFree(m, rt, args...) })
	r.Register("CLOUD.UPDATE", "cloud", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return cUpdate(m, rt, args...) })
	r.Register("CLOUD.DRAW", "cloud", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return cDraw(m, rt, args...) })
	r.Register("CLOUD.SETCOVERAGE", "cloud", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return cSetCoverage(m, rt, args...) })
	r.Register("CLOUD.GETCOVERAGE", "cloud", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return cGetCoverage(m, rt, args...) })
}

func cMake(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil || len(args) != 0 {
		return value.Nil, fmt.Errorf("CLOUD.MAKE expects no args")
	}
	id, err := m.h.Alloc(&CloudObject{Coverage: 0.3})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(int32(id)), nil
}

func cFree(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CLOUD.FREE expects handle")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	m.h.Free(heap.Handle(h))
	return value.Nil, nil
}

func cUpdate(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CLOUD.UPDATE expects cloud, dt#")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	if _, err := rt.ArgFloat(args, 1); err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func cDraw(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CLOUD.DRAW expects cloud")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	// Coverage reserved for future volumetric rendering.
	return value.FromHandle(h), nil
}

func cSetCoverage(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CLOUD.SETCOVERAGE expects cloud, amount#")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	v, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	o, err := heap.Cast[*CloudObject](m.h, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	o.Coverage = float32(v)
	return args[0], nil
}

func cGetCoverage(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CLOUD.GETCOVERAGE expects cloud handle")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	o, err := heap.Cast[*CloudObject](m.h, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(o.Coverage)), nil
}
