// Package mblight2d implements LIGHT2D.* and hooks a simple 2D light overlay each frame (CGO).
package mblight2d

import "moonbasic/vm/heap"

// Module registers LIGHT2D.* (and RENDER.SET2DAMBIENT when CGO enabled).
type Module struct {
	h *heap.Store
}

// NewModule creates the module.
func NewModule() *Module { return &Module{} }

// BindHeap binds the VM heap before Register.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

func (m *Module) Reset() {}
