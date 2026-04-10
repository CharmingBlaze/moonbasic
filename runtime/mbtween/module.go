package mbtween

import (
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// Module implements TWEEN.* for global float animation.
type Module struct {
	h *heap.Store

	invoke func(string, []value.Value) (value.Value, error)
	getG   func(string) (value.Value, bool)
	setG   func(string, value.Value)
}

// NewModule creates the module.
func NewModule() *Module { return &Module{} }

// BindHeap binds the VM heap before Register.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// SetUserInvoker wires VM.CallUserFunction for TWEEN.ONCOMPLETE.
func (m *Module) SetUserInvoker(fn func(string, []value.Value) (value.Value, error)) {
	m.invoke = fn
}

// SetGlobalAccessor wires read/write of script globals by name (uppercased like the VM).
func (m *Module) SetGlobalAccessor(get func(string) (value.Value, bool), set func(string, value.Value)) {
	m.getG, m.setG = get, set
}

func (m *Module) Reset() {}


