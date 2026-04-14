// Package mbentity registers Blitz3D-style ENTITY.* helpers (lightweight transforms + simple physics).
package mbentity

import (
	"moonbasic/runtime"
	mbcamera "moonbasic/runtime/camera"
	mbphysics3d "moonbasic/runtime/physics3d"
	"moonbasic/runtime/texture"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
	"unsafe"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ModulesByStore maps VM heaps to their entity module (for cross-package helpers).
var ModulesByStore = make(map[*heap.Store]*Module)

// Module holds entity state for one registry.
type Module struct {
	h   *heap.Store
	reg runtime.Registrar
	tex *texture.Module // set by compiler registry (LEVEL.LOADSKYBOX, etc.)
	cam *mbcamera.Module
	autoBuoyancy bool
}

// NewModule constructs the entity module.
func NewModule() *Module { return &Module{} }

// BindTextureModule wires the texture module for commands that load GPU textures from paths.
func (m *Module) BindTextureModule(t *texture.Module) { m.tex = t }

// BindCamera wires the camera module for world-space projection and shake helpers.
func (m *Module) BindCamera(c runtime.Module) {
	if cam, ok := c.(*mbcamera.Module); ok {
		m.cam = cam
		cam.SetEntityWorldPosHook(func(hs *heap.Store, eh heap.Handle) (float32, float32, float32, bool) {
			if m.h == nil || hs != m.h {
				return 0, 0, 0, false
			}
			p, ok := m.WorldPosFromEntityHandle(eh)
			if !ok {
				return 0, 0, 0, false
			}
			return p.X, p.Y, p.Z, true
		})
	}
}

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) {
	m.h = h
	ModulesByStore[h] = m
	// Wire physics3d mesh lookup for LEVEL.STATIC
	mbphysics3d.SetMeshLookupForHeap(h, m.GetEntityMeshes)
	mbphysics3d.SetVehicleHooksForHeap(h, func(id int64) (rl.Vector3, float32, bool) {
		e := m.store().ents[id]
		if e == nil {
			return rl.Vector3{}, 0, false
		}
		p := m.worldPos(e)
		_, ew, _ := e.getRot()
		return p, ew, true
	}, func(id int64, pos rl.Vector3) {
		e := m.store().ents[id]
		if e == nil {
			return
		}
		m.setLocalFromWorld(e, pos.X, pos.Y, pos.Z)
	})
	mbphysics3d.SetRaycastHook(func(ox, oy, oz, maxDown float64) (nx, ny, nz, hitY float64, ok bool) {
		// Bridge to entity floor query
		y := m.queryFloorYAt(float32(ox), float32(oy), float32(oz))
		if y > float64(oy)-maxDown && y < float64(oy)+1.0 {
			return 0, 1, 0, y, true
		}
		return 0, 1, 0, 0, false
	})
	m.reg.Register("WATER.AUTOPHYSICS", "entity", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("WATER.AUTOPHYSICS expects (toggle)")
		}
		on, _ := rt.ArgBool(args, 0)
		m.autoBuoyancy = on
		return value.Nil, nil
	})
}

// GetEntityMeshes retrieves all meshes associated with an entity (via rlModel if CGO).
func (m *Module) GetEntityMeshes(id int64) []rl.Mesh {
	e := m.store().ents[id]
	if e == nil || !e.hasRLModel {
		return nil
	}
	// Models can have multiple meshes
	mCount := int(e.rlModel.MeshCount)
	if mCount <= 0 {
		return nil
	}
	// Return a slice of meshes
	meshes := (*[1 << 30]rl.Mesh)(unsafe.Pointer(e.rlModel.Meshes))[:mCount:mCount]
	return meshes
}

