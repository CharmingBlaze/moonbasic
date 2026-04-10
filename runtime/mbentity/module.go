// Package mbentity registers Blitz3D-style ENTITY.* helpers (lightweight transforms + simple physics).
package mbentity

import (
	"moonbasic/runtime"
	mbcamera "moonbasic/runtime/camera"
	"moonbasic/runtime/texture"
	"moonbasic/vm/heap"

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
}

// NewModule constructs the entity module.
func NewModule() *Module { return &Module{} }

// BindTextureModule wires the texture module for commands that load GPU textures from paths.
func (m *Module) BindTextureModule(t *texture.Module) { m.tex = t }

// BindCamera wires the camera module for world-space projection and shake helpers.
func (m *Module) BindCamera(c runtime.Module) {
	if cam, ok := c.(*mbcamera.Module); ok {
		m.cam = cam
	}
}

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) {
	m.h = h
	ModulesByStore[h] = m
}

// GetWorldPosByID resolves the world position of an entity handle securely for cross-module helpers.
func (m *Module) GetWorldPosByID(id int) (rl.Vector3, bool) {
	st := m.store() // Ensure we have a store() helper or access ents.
	if id < 1 || id >= len(st.ents) { return rl.Vector3{}, false }
	e := st.ents[int64(id)]
	if e == nil { return rl.Vector3{}, false }
	return m.worldPos(e), true
}


