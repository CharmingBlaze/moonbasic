// Package mbcollision implements RAY.*, BBOX.*, BSPHERE.* (Raylib geometric queries) when CGO is enabled.
package mbcollision

import "moonbasic/vm/heap"

// Module registers collision / ray builtins.
type Module struct {
	h *heap.Store
}

// NewModule creates the module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }
