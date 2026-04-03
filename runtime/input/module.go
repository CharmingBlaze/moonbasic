// Package input implements INPUT.* builtins (keyboard, etc.) backed by Raylib when CGO is enabled.
package input

import "moonbasic/vm/heap"

// Module registers INPUT.* handlers into the runtime Registry command map.
type Module struct {
	h *heap.Store
}

// NewModule returns a new input module.
func NewModule() *Module {
	return &Module{}
}

// BindHeap implements runtime.HeapAware (INPUT.GETMOUSEWORLDPOS allocates numeric arrays).
func (m *Module) BindHeap(h *heap.Store) { m.h = h }
