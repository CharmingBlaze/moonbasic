// Package mbgrid implements GRID.* tactical cell helpers (flat array + optional terrain Y).
package mbgrid

import (
	"fmt"
	"sync"

	"moonbasic/runtime"
	"moonbasic/runtime/mbentity"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// Module registers GRID.* commands.
type Module struct {
	mu  sync.Mutex
	h   *heap.Store
	ent *mbentity.Module
}

func NewModule() *Module { return &Module{} }

func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// BindEntity wires the entity module for GRID.SNAP and placement helpers.
func (m *Module) BindEntity(mod runtime.Module) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if e, ok := mod.(*mbentity.Module); ok {
		m.ent = e
	}
}

func (m *Module) requireHeap(rt *runtime.Runtime) (*heap.Store, error) {
	if rt != nil && rt.Heap != nil {
		return rt.Heap, nil
	}
	if m.h != nil {
		return m.h, nil
	}
	return nil, fmt.Errorf("grid: heap not bound")
}

func argF64(v value.Value) (float64, bool) {
	if f, ok := v.ToFloat(); ok {
		return f, true
	}
	if i, ok := v.ToInt(); ok {
		return float64(i), true
	}
	return 0, false
}
