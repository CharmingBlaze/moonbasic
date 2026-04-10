// Package mbmodel3d implements MESH.*, MATERIAL.*, MODEL.*, and SHADER.LOAD (Raylib 3D) when CGO is enabled.
package mbmodel3d

import "moonbasic/vm/heap"

// Module registers 3D mesh/material/model/shader builtins.
type Module struct {
	h              *heap.Store
	enqueueCleanup func(func())
	// scratch buffers for shader uniforms (zero-alloc)
	u1 []float32
	u2 []float32
	u3 []float32
	u4 []float32
}

// NewModule creates the module.
func NewModule() *Module {
	return &Module{
		u1: make([]float32, 1),
		u2: make([]float32, 2),
		u3: make([]float32, 3),
		u4: make([]float32, 4),
	}
}

// BindHeap binds the VM heap before Register (implements runtime.HeapAware).
func (m *Module) BindHeap(h *heap.Store) {
	m.h = h
	SetGlobalHeapGetter(func() *heap.Store { return m.h })
}

// BindCleanup receives the main-thread enqueuer (implements runtime.MainThreadCleanupAware).
func (m *Module) BindCleanup(enqueuer func(func())) {
	m.enqueueCleanup = enqueuer
	setGlobalCleanupEnqueuer(enqueuer)
}

var globalEnqueuer func(func())

func setGlobalCleanupEnqueuer(fn func(func())) {
	globalEnqueuer = fn
}

func enqueueOnMainThread(fn func()) {
	if globalEnqueuer != nil {
		globalEnqueuer(fn)
	} else {
		// No window/cleanup provider? Raylib calls might crash.
		// Silently drop or log? Drop for safety in a finalizer context.
	}
}

func (m *Module) Reset() {}

