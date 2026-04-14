// Package water implements WATER.* (Raylib CGO).
package water

import (
	"moonbasic/vm/heap"
)

type Module struct {
	h      *heap.Store
	waters []heap.Handle
}

func NewModule() *Module { return &Module{} }

func (m *Module) BindHeap(h *heap.Store) { m.h = h }

func (m *Module) Shutdown() {}

func (m *Module) Reset() {}

