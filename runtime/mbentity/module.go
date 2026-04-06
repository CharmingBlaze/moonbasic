// Package mbentity registers Blitz3D-style ENTITY.* helpers (lightweight transforms + simple physics).
package mbentity

import "moonbasic/vm/heap"

// Module holds entity state for one registry.
type Module struct {
	h *heap.Store
}

// NewModule constructs the entity module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }
