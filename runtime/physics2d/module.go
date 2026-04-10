// Package mbphysics2d registers PHYSICS2D.* and BODY2D.* (stubs until Box2D is integrated).
package mbphysics2d

import (
	"moonbasic/vm/heap"
)

// Module is the 2D physics orchestrator.
type Module struct {
	h *heap.Store
}

// NewModule creates the physics2d module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }


// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

func (m *Module) Reset() {}

