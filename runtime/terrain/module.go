// Package terrain implements TERRAIN.* and CHUNK.* heightfield commands (Raylib CGO).
package terrain

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

// Module registers terrain builtins.
type Module struct {
	h      *heap.Store
	active heap.Handle
}

// NewModule creates the terrain module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	registerTerrain(m, r)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

// ActiveHandle returns the most recently created terrain handle (for world streaming).
func (m *Module) ActiveHandle() heap.Handle { return m.active }
