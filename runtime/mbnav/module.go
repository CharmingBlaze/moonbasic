// Package mbnav implements NAV.*, PATH.*, NAVAGENT.*, STEER.*, and BTREE.* (grid pathfinding,
// simple agents, steering forces as VEC3 handles, and VM-callback behaviour trees).
package mbnav

import (
	"fmt"
	"sync"

	"moonbasic/runtime"
	"moonbasic/runtime/mbentity"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// Module registers navigation and AI-related builtins.
type Module struct {
	mu     sync.Mutex
	h      *heap.Store
	ent    *mbentity.Module
	invoke func(string, []value.Value) (value.Value, error)

	// terrainNav maps a baked terrain heap handle to its NAV grid handle (see NAV.BAKE).
	terrainNav map[heap.Handle]heap.Handle
	// enemyFollow tracks waypoint indices for ENEMY.FOLLOWPATH(entity#, path#, speed#).
	enemyFollow map[int64]enemyFollowState
}

type enemyFollowState struct {
	pathH heap.Handle
	idx   int
}

// NewModule creates the module.
func NewModule() *Module {
	return &Module{
		terrainNav:  make(map[heap.Handle]heap.Handle),
		enemyFollow: make(map[int64]enemyFollowState),
	}
}

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// BindEntity binds the entity module for world scanning.
func (m *Module) BindEntity(mod runtime.Module) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if e, ok := mod.(*mbentity.Module); ok {
		m.ent = e
	}
}

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

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

// Reset implements runtime.Module.
func (m *Module) Reset() {}

