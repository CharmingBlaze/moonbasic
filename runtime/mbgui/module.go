// Package mbgui registers GUI.* immediate-mode widgets: full raygui when CGO is enabled,
// or a minimal Raylib-drawn subset on Windows when CGO is disabled.
package mbgui

import "moonbasic/vm/heap"

// Module implements GUI bindings for moonBASIC.
type Module struct {
	h *heap.Store
}

// NewModule constructs the GUI builtin module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }
