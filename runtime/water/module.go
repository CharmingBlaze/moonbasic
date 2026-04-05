// Package water implements WATER.* (Raylib CGO).
package water

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

type Module struct {
	h      *heap.Store
	waters []heap.Handle
}

func NewModule() *Module { return &Module{} }

func (m *Module) BindHeap(h *heap.Store) { m.h = h }

func (m *Module) Register(r runtime.Registrar) { registerWater(m, r) }

func (m *Module) Shutdown() {}
