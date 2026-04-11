// Package mbcamera registers CAMERA.* and CAMERA2D.* (Raylib when CGO enabled).
package mbcamera

import "moonbasic/vm/heap"

// Module holds 3D/2D camera natives.
type Module struct {
	h *heap.Store

	// lastActive3D is set by CAMERA.BEGIN for CAMERA.GETACTIVE / shadow helpers.
	lastActive3D heap.Handle

	// entityWorldPos resolves an EntityRef handle to world XYZ (set by mbentity.BindCamera).
	entityWorldPos func(*heap.Store, heap.Handle) (x, y, z float32, ok bool)
}

// NewModule creates a camera module.
func NewModule() *Module { return &Module{} }

// SetEntityWorldPosHook is called from mbentity.BindCamera so CAMERA.ORBIT(cam, entity, dist) can track targets.
func (m *Module) SetEntityWorldPosHook(fn func(*heap.Store, heap.Handle) (float32, float32, float32, bool)) {
	m.entityWorldPos = fn
}

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }
