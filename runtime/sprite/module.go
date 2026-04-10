// Package mbsprite registers SPRITE.* (Raylib when CGO enabled).
package mbsprite

import "moonbasic/vm/heap"

// Module holds sprite/animation natives.
type Module struct {
	h *heap.Store
}

// NewModule creates a sprite module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

func (m *Module) Reset() {}


