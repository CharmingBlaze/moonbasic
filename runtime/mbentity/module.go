// Package mbentity registers Blitz3D-style ENTITY.* helpers (lightweight transforms + simple physics).
package mbentity

import (
	"moonbasic/runtime/texture"
	"moonbasic/vm/heap"
)

// modulesByHeap maps VM heaps to their entity module (for cross-package helpers).
var modulesByHeap = make(map[*heap.Store]*Module)

// Module holds entity state for one registry.
type Module struct {
	h *heap.Store
	tex *texture.Module // set by compiler registry (LEVEL.LOADSKYBOX, etc.)
}

// NewModule constructs the entity module.
func NewModule() *Module { return &Module{} }

// BindTextureModule wires the texture module for commands that load GPU textures from paths.
func (m *Module) BindTextureModule(t *texture.Module) { m.tex = t }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) {
	m.h = h
	modulesByHeap[h] = m
}
