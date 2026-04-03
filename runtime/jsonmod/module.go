// Package mbjson implements JSON.* heap-backed flat JSON objects (encoding/json).
package mbjson

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

// Module registers JSON.* commands.
type Module struct {
	h *heap.Store
}

// NewModule creates the json module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	registerJSONCommands(m, r)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
