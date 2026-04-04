package mbevent

import "moonbasic/vm/value"

// Module implements EVENT.* global bus.
type Module struct {
	invoke func(string, []value.Value) (value.Value, error)
}

// NewModule creates the module.
func NewModule() *Module { return &Module{} }

// SetUserInvoker wires VM.CallUserFunction for listeners.
func (m *Module) SetUserInvoker(fn func(string, []value.Value) (value.Value, error)) {
	m.invoke = fn
}
