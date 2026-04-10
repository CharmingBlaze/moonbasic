// Package mbimage implements IMAGE.* builtins (Raylib CPU images) when CGO is enabled.
package mbimage

import "moonbasic/vm/heap"

// Module registers IMAGE.* handlers.
type Module struct {
	h *heap.Store
}

// NewModule creates an image builtin module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

func (m *Module) Reset() {}


