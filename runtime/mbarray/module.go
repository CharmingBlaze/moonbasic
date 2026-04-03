package mbarray

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

// Module implements ARRAY.*-style flat builtins (ARRAYLEN, ARRAYPUSH, ...).
type Module struct {
	h *heap.Store
}

// NewModule creates an array helper module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	registerAll(m, r)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
