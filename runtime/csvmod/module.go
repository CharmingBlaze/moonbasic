// Package mbcsv implements CSV.* heap-backed tables (encoding/csv).
package mbcsv

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

// Module registers CSV.* commands.
type Module struct {
	h *heap.Store
}

// NewModule creates the CSV module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	registerCSVCommands(m, r)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
