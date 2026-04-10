// Package terrain implements TERRAIN.* and CHUNK.* heightfield commands (Raylib CGO).
package terrain

import (
	mbentity "moonbasic/runtime/mbentity"
	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

// Module registers terrain builtins.
type Module struct {
	h      *heap.Store
	active heap.Handle
	ent    *mbentity.Module
}

// NewModule creates a terrain module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// BindEntity wires the entity module for TERRAIN.APPLYTILES (ENTITY.COPY + position).
func (m *Module) BindEntity(mod *mbentity.Module) { m.ent = mod }

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	registerTerrain(m, r)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

// ActiveHandle returns the most recently created terrain handle (for world streaming).
func (m *Module) ActiveHandle() heap.Handle { return m.active }

func (m *Module) Reset() {}

