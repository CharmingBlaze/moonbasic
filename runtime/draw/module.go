// Package mbdraw registers DRAW.* and (when CGO) TEXTURE.LOAD / TEXTURE.FREE / DRAW.TEXTURE.
package mbdraw

import "moonbasic/vm/heap"

// Module is the draw/texture builtin module.
type Module struct {
	h *heap.Store
}

// NewModule creates DRAW.* registration.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }
