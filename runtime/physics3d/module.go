// Package mbphysics3d registers PHYSICS3D.* and BODY3D.* (Jolt on linux+cgo; stubs elsewhere).
package mbphysics3d

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// UserInvokerSetter allows the pipeline to wire VM.CallUserFunction for collision callbacks.
type UserInvokerSetter interface {
	SetUserInvoker(fn func(string, []value.Value) (value.Value, error))
}

// Module is the 3D physics orchestrator.
type Module struct {
	h      *heap.Store
	invoke func(string, []value.Value) (value.Value, error)
}

// NewModule creates the physics3d module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// SetUserInvoker wires callback dispatch from PHYSICS3D.PROCESSCOLLISIONS.
func (m *Module) SetUserInvoker(fn func(string, []value.Value) (value.Value, error)) {
	m.invoke = fn
}

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	registerPhysics3DCommands(m, reg)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {
	shutdownPhysics3D(m)
}
