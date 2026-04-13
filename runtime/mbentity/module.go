// Package mbentity registers Blitz3D-style ENTITY.* helpers (lightweight transforms + simple physics).
package mbentity

import (
	"moonbasic/runtime"
	mbcamera "moonbasic/runtime/camera"
	"moonbasic/runtime/texture"
	"moonbasic/vm/heap"
)

// ModulesByStore maps VM heaps to their entity module (for cross-package helpers).
var ModulesByStore = make(map[*heap.Store]*Module)

// Module holds entity state for one registry.
type Module struct {
	h   *heap.Store
	reg runtime.Registrar
	tex *texture.Module // set by compiler registry (LEVEL.LOADSKYBOX, etc.)
	cam *mbcamera.Module
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
}

