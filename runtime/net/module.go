// Package mbnet implements NET.*, PEER.*, and EVENT.* using ENet (CGO) or stubs.
package mbnet

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

// Module registers networking commands.
type Module struct {
	h *heap.Store
}

// NewModule creates the net module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	registerNetCommands(m, reg)
}

// Shutdown releases ENet and any open hosts (CGO path).
func (m *Module) Shutdown() {
	shutdownNet(m)
}
