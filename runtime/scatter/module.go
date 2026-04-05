package scatter

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

type Module struct {
	h     *heap.Store
	props []heap.Handle
}

func NewModule() *Module { return &Module{} }
func (m *Module) BindHeap(h *heap.Store) { m.h = h }
func (m *Module) Register(r runtime.Registrar) { registerScatter(m, r) }
func (m *Module) Shutdown() {}
