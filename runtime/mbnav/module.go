// Package mbnav implements NAV.*, PATH.*, NAVAGENT.*, STEER.*, and BTREE.* (grid pathfinding,
// simple agents, steering forces as VEC3 handles, and VM-callback behaviour trees).
package mbnav

import (
	"fmt"
	"sync"

	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// Module registers navigation and AI-related builtins.
type Module struct {
	mu     sync.Mutex
	h      *heap.Store
	invoke func(string, []value.Value) (value.Value, error)
}

// NewModule creates the module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// SetUserInvoker wires VM.CallUserFunction for BTREE.Run condition/action callbacks.
func (m *Module) SetUserInvoker(fn func(string, []value.Value) (value.Value, error)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.invoke = fn
}

func (m *Module) callUser(name string, args []value.Value) (value.Value, error) {
	m.mu.Lock()
	fn := m.invoke
	m.mu.Unlock()
	if fn == nil {
		return value.Nil, fmt.Errorf("BTREE.*: SetUserInvoker not configured (host must wire VM)")
	}
	return fn(name, args)
}
