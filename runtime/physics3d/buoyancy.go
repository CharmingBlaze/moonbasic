package mbphysics3d

import (
	"fmt"
	"sync"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

var (
	buoyMu     sync.Mutex
	buoyByEnt  map[int64]float64
)

func registerBuoyancyCommands(m *Module, reg runtime.Registrar) {
	_ = m
	reg.Register("PHYSICS.SETBUOYANCY", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phSetBuoyancy(a) }))
	reg.Register("PHYSICS.GETBUOYANCY", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phGetBuoyancy(a) }))
	reg.Register("ENTITY.SETBUOYANCY", "entity", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phSetBuoyancy(a) }))
	reg.Register("ENTITY.GETBUOYANCY", "entity", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phGetBuoyancy(a) }))
}

func phSetBuoyancy(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PHYSICS.SETBUOYANCY expects (entity#, density#)")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PHYSICS.SETBUOYANCY: invalid entity")
	}
	d, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("PHYSICS.SETBUOYANCY: density must be numeric")
	}
	buoyMu.Lock()
	defer buoyMu.Unlock()
	if buoyByEnt == nil {
		buoyByEnt = make(map[int64]float64)
	}
	buoyByEnt[id] = d
	return value.Nil, nil
}

func phGetBuoyancy(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PHYSICS.GETBUOYANCY expects (entity#)")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PHYSICS.GETBUOYANCY: invalid entity")
	}
	buoyMu.Lock()
	defer buoyMu.Unlock()
	v := buoyByEnt[id]
	return value.FromFloat(v), nil
}
