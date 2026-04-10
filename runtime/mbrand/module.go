// Package mbrand implements RAND.* handle-based PRNGs (math/rand/v2).
package mbrand

import "moonbasic/vm/heap"

// Module registers RNG builtins.
type Module struct {
	h *heap.Store
}

// NewModule creates the module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

func (m *Module) Reset() {}


