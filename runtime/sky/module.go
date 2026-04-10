package sky

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

type Module struct {
	h *heap.Store
}

func NewModule() *Module { return &Module{} }
func (m *Module) BindHeap(h *heap.Store) { m.h = h }
func (m *Module) Register(r runtime.Registrar) { registerSky(m, r) }
func (m *Module) Shutdown() {}

func (m *Module) Reset() {}

