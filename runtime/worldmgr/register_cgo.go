//go:build cgo || (windows && !cgo)

package worldmgr

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerWorld(m *Module, r runtime.Registrar) {
	r.Register("WORLD.SETCENTER", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return worldSetCenter(m, rt, args...) })
	r.Register("WORLD.UPDATE", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return worldUpdate(m, rt, args...) })
	r.Register("WORLD.STREAMENABLE", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return worldStreamEnable(m, rt, args...) })
	r.Register("WORLD.PRELOAD", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return worldPreload(m, rt, args...) })
	r.Register("WORLD.STATUS", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return worldStatus(m, rt, args...) })
	r.Register("WORLD.ISREADY", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return worldIsReady(m, rt, args...) })
}

func worldSetCenter(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("WORLD.SETCENTER expects x#, z#")
	}
	x, err := rt.ArgFloat(args, 0)
	if err != nil {
		return value.Nil, err
	}
	z, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	m.terr.SetCenter(float32(x), float32(z))
	return value.Nil, nil
}

func worldUpdate(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WORLD.UPDATE expects dt#")
	}
	_, err := rt.ArgFloat(args, 0)
	if err != nil {
		return value.Nil, err
	}
	m.terr.TickStreaming(rt)
	return value.Nil, nil
}

func worldStreamEnable(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WORLD.STREAMENABLE expects enabled?")
	}
	b, err := rt.ArgBool(args, 0)
	if err != nil {
		return value.Nil, err
	}
	m.terr.SetStreamEnabled(b)
	return value.Nil, nil
}

func worldPreload(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("WORLD.PRELOAD expects terrain, radius")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	rad, err := rt.ArgInt(args, 1)
	if err != nil {
		return value.Nil, err
	}
	err = m.terr.PreloadTerrain(heap.Handle(h), int(rad))
	if err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func worldStatus(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WORLD.STATUS expects no arguments")
	}
	return rt.RetString(m.terr.StatusString()), nil
}

func worldIsReady(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WORLD.ISREADY expects terrain")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(m.terr.IsReadyTerrain(heap.Handle(h))), nil
}
