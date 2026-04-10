// Package mbpool implements POOL.* — reusable handle pools with user FUNCTION factory/reset.
package mbpool

import (
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// Module registers pool builtins.
type Module struct {
	h      *heap.Store
	invoke func(string, []value.Value) (value.Value, error)
}

// NewModule creates the module.
func NewModule() *Module { return &Module{} }

// BindHeap binds the VM heap before Register.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// SetUserInvoker wires VM.CallUserFunction for POOL.SETFACTORY / GET / RETURN / PREWARM.
func (m *Module) SetUserInvoker(fn func(string, []value.Value) (value.Value, error)) {
	m.invoke = fn
}

func (m *Module) Reset() {}


