// Package mbdb implements DB.* SQLite handles: mattn driver when CGO is on,
// modernc.org/sqlite when built with -tags modernc_sqlite, or stubs otherwise.
package mbdb

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

// Module registers database commands.
type Module struct {
	h *heap.Store
}

// NewModule creates the db module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	registerDBCommands(m, r)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
