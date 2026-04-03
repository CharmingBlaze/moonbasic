// Package mbfont registers FONT.* (Raylib when CGO enabled).
// Handle 0 is reserved for FONT.DRAWDEFAULT (Raylib built-in default font).
package mbfont

import "moonbasic/vm/heap"

// Module registers font natives.
type Module struct {
	h *heap.Store
}

// NewModule creates a font module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }
