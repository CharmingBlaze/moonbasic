// Package mbnet implements NET.*, PEER.*, and EVENT.* using ENet (CGO) or stubs.
package mbnet

import (
	"fmt"
	"sync"

	mbruntime "moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// Module registers networking commands.
type Module struct {
	h *heap.Store

	invokeMu sync.Mutex
	invoke   func(string, []value.Value) (value.Value, error)
}

// NewModule creates the net module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// SetUserInvoker wires VM.CallUserFunction for SERVER/CLIENT/RPC/LOBBY callbacks.
func (m *Module) SetUserInvoker(fn func(string, []value.Value) (value.Value, error)) {
	m.invokeMu.Lock()
	defer m.invokeMu.Unlock()
	m.invoke = fn
}

func (m *Module) callUser(name string, args []value.Value) (value.Value, error) {
	m.invokeMu.Lock()
	fn := m.invoke
	m.invokeMu.Unlock()
	if fn == nil {
		return value.Nil, fmt.Errorf("multiplayer callbacks require SetUserInvoker (host must wire VM)")
	}
	return fn(name, args)
}

// Register implements runtime.Module.
func (m *Module) Register(reg mbruntime.Registrar) {
	registerNetCommands(m, reg)
}

// Shutdown releases ENet and any open hosts (CGO path).
func (m *Module) Shutdown() {
	shutdownNet(m)
}
