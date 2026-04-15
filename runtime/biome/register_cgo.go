//go:build cgo || (windows && !cgo)

package biome

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerBiome(m *Module, r runtime.Registrar) {
	r.Register("BIOME.CREATE", "biome", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return bMake(m, rt, args...) })
	r.Register("BIOME.MAKE", "biome", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return bMake(m, rt, args...) })
	r.Register("BIOME.FREE", "biome", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return bFree(m, rt, args...) })
	r.Register("BIOME.SETTEMP", "biome", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return bSetTemp(m, rt, args...) })
	r.Register("BIOME.SETHUMIDITY", "biome", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return bSetHumidity(m, rt, args...)
	})
}

func bMake(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil || len(args) != 1 {
		return value.Nil, fmt.Errorf("BIOME.MAKE expects name")
	}
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	id, err := m.h.Alloc(&BiomeObject{Name: s})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(int32(id)), nil
}

func bFree(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("BIOME.FREE expects handle")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	m.h.Free(heap.Handle(h))
	return value.Nil, nil
}

func bSetTemp(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("BIOME.SETTEMP expects biome, celsius")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	t, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	o, err := heap.Cast[*BiomeObject](m.h, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	o.TempC = float32(t)
	return value.Nil, nil
}

func bSetHumidity(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("BIOME.SETHUMIDITY expects biome, amount")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	v, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	o, err := heap.Cast[*BiomeObject](m.h, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	o.Humidity = float32(v)
	return value.Nil, nil
}
