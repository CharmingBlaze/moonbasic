// Package mbcamera registers CAMERA.* and CAMERA2D.* (Raylib when CGO enabled).
package mbcamera

import "moonbasic/vm/heap"

// Module holds 3D/2D camera natives.
type Module struct {
	h *heap.Store

	// lastActive3D is set by CAMERA.BEGIN for CAMERA.GETACTIVE / shadow helpers.
	lastActive3D heap.Handle
}

// NewModule creates a camera module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }
