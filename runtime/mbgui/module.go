// Package mbgui registers GUI.* raygui (Raylib immediate-mode widgets) when CGO is enabled.
package mbgui

import "moonbasic/vm/heap"

// Module implements raygui bindings for moonBASIC.
type Module struct {
	h *heap.Store
}

// NewModule constructs the GUI builtin module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }
