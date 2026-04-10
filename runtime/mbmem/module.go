// Package mbmem implements MEM.* heap-backed raw byte buffers.
package mbmem

import "moonbasic/vm/heap"

// Module registers memory block builtins.
type Module struct {
	h *heap.Store
}

// NewModule creates the module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

func (m *Module) Reset() {}


