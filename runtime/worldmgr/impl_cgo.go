//go:build cgo || (windows && !cgo)

package worldmgr

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/runtime/mbentity"
	"moonbasic/runtime/mbmodel3d"
	scat "moonbasic/runtime/scatter"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func syncWorldFogGPU(m *Module) {
	mbmodel3d.SyncSceneFogWorld(m.FogMode, m.FogColor[0], m.FogColor[1], m.FogColor[2], m.FogDensity)
}

func (m *Module) worldSetReflection(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WORLD.SETREFLECTION expects (entity)")
	}
	_, err := rt.ArgInt(args, 0)
	if err != nil {
		return value.Nil, err
	}
	return value.Nil, fmt.Errorf("WORLD.SETREFLECTION: reflection probe capture not implemented yet")
}

func (m *Module) worldFogMode(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("FOGMODE expects mode")
	}
	mode, _ := rt.ArgInt(args, 0)
	m.FogMode = int(mode)
	syncWorldFogGPU(m)
	return value.Nil, nil
}

func (m *Module) worldFogColor(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("FOGCOLOR expects (r, g, b)")
	}
	r, _ := rt.ArgInt(args, 0)
	g, _ := rt.ArgInt(args, 1)
	b, _ := rt.ArgInt(args, 2)
	m.FogColor = [4]uint8{uint8(r), uint8(g), uint8(b), 255}
	syncWorldFogGPU(m)
	return value.Nil, nil
}

func (m *Module) worldFogDensity(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("FOGDENSITY expects density")
	}
	d, _ := rt.ArgFloat(args, 0)
	m.FogDensity = float32(d)
	syncWorldFogGPU(m)
	return value.Nil, nil
}

func worldSetCenter(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("WORLD.SETCENTER expects x, z")
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

func worldSetCenterEntity(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WORLD.SETCENTERENTITY expects entity")
	}
	id, err := rt.ArgInt(args, 0)
	if err != nil {
		return value.Nil, err
	}
	x, z, err := mbentity.EntityWorldXZ(m.h, id)
	if err != nil {
		return value.Nil, err
	}
	m.terr.SetCenter(x, z)
	return value.Nil, nil
}

func worldUpdate(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WORLD.UPDATE expects dt")
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
		return value.Nil, fmt.Errorf("WORLD.STREAMENABLE expects enabled")
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

func worldSetVegetation(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("WORLD.SETVEGETATION: heap not bound")
	}
	if m.scat == nil {
		return value.Nil, fmt.Errorf("WORLD.SETVEGETATION: scatter module not wired (internal)")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("WORLD.SETVEGETATION expects (terrain, billboard, density)")
	}
	ht, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	if _, err := rt.ArgHandle(args, 1); err != nil {
		return value.Nil, err
	}
	den, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	if m.vegScatter == 0 {
		id, err := m.h.Alloc(&scat.ScatterObject{Name: "world_vegetation", Seed: 0x5EED})
		if err != nil {
			return value.Nil, err
		}
		m.vegScatter = id
	}
	if err := m.scat.ApplyToTerrain(m.vegScatter, heap.Handle(ht), den); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}
