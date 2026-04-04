package texture

import "moonbasic/vm/heap"

type Module struct {
	h *heap.Store
}

func NewModule() *Module { return &Module{} }

func (m *Module) BindHeap(h *heap.Store) { m.h = h }
