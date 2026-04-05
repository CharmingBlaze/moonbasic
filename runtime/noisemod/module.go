// Package noisemod implements NOISE.* stateful procedural noise (pure Go; shared core in runtime/procnoise).
package noisemod

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

// Module registers NOISE.* builtins.
type Module struct {
	h *heap.Store
}

// NewModule creates the noise module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) { registerNoise(m, r) }

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
