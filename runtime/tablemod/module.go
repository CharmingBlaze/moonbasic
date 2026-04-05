// Package mbtable implements TABLE.* in-memory tables with JSON/CSV bridges.
package mbtable

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

// Module registers TABLE.* commands.
type Module struct {
	h *heap.Store
}

// NewModule creates the table module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	registerTableCommands(m, r)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
