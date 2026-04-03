// Package mbcharcontroller registers CHARCONTROLLER.* (Jolt CharacterVirtual on Linux+CGO).
package mbcharcontroller

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

// Module is the character controller orchestrator.
type Module struct {
	h *heap.Store
}

// NewModule creates the charcontroller module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	registerCharControllerCommands(m, reg)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {
	shutdownCharController(m)
}
