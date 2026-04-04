// Package mbparticles implements PARTICLE.* billboard emitters (Raylib) when CGO is enabled.
package mbparticles

import "moonbasic/vm/heap"

// Module registers particle builtins.
type Module struct {
	h *heap.Store
}

// NewModule creates the module.
func NewModule() *Module { return &Module{} }

// BindHeap binds the VM heap before Register (implements runtime.HeapAware).
func (m *Module) BindHeap(h *heap.Store) { m.h = h }
