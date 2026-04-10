//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// GetTerrainHeightCb breaks the import cycle (terrain -> mbentity -> terrain). Assigned by terrain during init.
var GetTerrainHeightCb func(h *heap.Store, th heap.Handle, x, z float32) float32

// EntityWorldXZ returns world-space X and Z for an entity id (for WORLD.SETCENTERENTITY, etc.).
func EntityWorldXZ(h *heap.Store, id int64) (float32, float32, error) {
	mod := ModulesByStore[h]
	if mod == nil {
		return 0, 0, fmt.Errorf("EntityWorldXZ: entity module not bound")
	}
	e := mod.store().ents[id]
	if e == nil {
		return 0, 0, fmt.Errorf("EntityWorldXZ: unknown entity %d", id)
	}
	wp := mod.worldPos(e)
	return wp.X, wp.Z, nil
}

func (m *Module) entGetXZ(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ENTITY.GETXZ: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.GETXZ expects entity#")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.GETXZ: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.GETXZ: unknown entity %d", id)
	}
	p := m.worldPos(e)
	arr, err := heap.NewArrayOfKind([]int64{2}, heap.ArrayKindFloat, 0)
	if err != nil {
		return value.Nil, err
	}
	arr.Floats[0] = float64(p.X)
	arr.Floats[1] = float64(p.Z)
	h, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

// entClampToTerrain implements ENTITY.CLAMPTOTERRAIN — same as TERRAIN.SNAPY with yOffset 0, argument order (entity, terrain).
func (m *Module) entClampToTerrain(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ENTITY.CLAMPTOTERRAIN: heap not bound")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.CLAMPTOTERRAIN expects (entity#, terrain#)")
	}
	return m.entTerrainSnapY([]value.Value{args[1], args[0], value.FromFloat(0)})
}

// entTerrainSnapY implements TERRAIN.SNAPY — samples terrain at entity world XZ and sets Y.
func (m *Module) entTerrainSnapY(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TERRAIN.SNAPY: heap not bound")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("TERRAIN.SNAPY expects (terrain, entity#, yOffset#)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TERRAIN.SNAPY: terrain must be a handle")
	}
	th := heap.Handle(args[0].IVal)
	id, ok := m.entID(args[1])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("TERRAIN.SNAPY: invalid entity")
	}
	off, ok := argF32(args[2])
	if !ok {
		return value.Nil, fmt.Errorf("TERRAIN.SNAPY: yOffset must be numeric")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("TERRAIN.SNAPY: unknown entity %d", id)
	}
	wp := m.worldPos(e)
	var hy float32 = off
	if GetTerrainHeightCb != nil {
		hy = GetTerrainHeightCb(m.h, th, wp.X, wp.Z) + off
	}
	return m.entSetPosition([]value.Value{
		args[1],
		value.FromFloat(float64(wp.X)),
		value.FromFloat(float64(hy)),
		value.FromFloat(float64(wp.Z)),
		value.FromBool(true),
	})
}

// entTerrainPlace sets world XZ, samples terrain height at that XZ, sets Y = height + offset (Position + SnapY in one call).
func (m *Module) entTerrainPlace(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TERRAIN.PLACE: heap not bound")
	}
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("TERRAIN.PLACE expects (terrain, entity#, x#, z#, yOffset#)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TERRAIN.PLACE: terrain must be a handle")
	}
	th := heap.Handle(args[0].IVal)
	id, ok := m.entID(args[1])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("TERRAIN.PLACE: invalid entity")
	}
	x, ok1 := argF32(args[2])
	z, ok2 := argF32(args[3])
	off, ok3 := argF32(args[4])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("TERRAIN.PLACE: x, z, yOffset must be numeric")
	}
	if m.store().ents[id] == nil {
		return value.Nil, fmt.Errorf("TERRAIN.PLACE: unknown entity %d", id)
	}
	var hy float32 = off
	if GetTerrainHeightCb != nil {
		hy = GetTerrainHeightCb(m.h, th, x, z) + off
	}
	return m.entSetPosition([]value.Value{
		args[1],
		value.FromFloat(float64(x)),
		value.FromFloat(float64(hy)),
		value.FromFloat(float64(z)),
		value.FromBool(true),
	})
}
