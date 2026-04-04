// Package mbmatrix implements TRANSFORM.* / MAT4.*, VEC3.*, and VEC2.* heap-backed math (Raylib raymath) when CGO is enabled.
package mbmatrix

import "moonbasic/vm/heap"

// Module registers matrix and vector builtins.
type Module struct {
	h *heap.Store
}

// NewModule creates the module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }
