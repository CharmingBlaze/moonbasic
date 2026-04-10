// Package mbgui registers GUI.* immediate-mode widgets: full raygui when CGO is enabled,
// or a minimal Raylib-drawn subset on Windows when CGO is disabled.
package mbgui

import (
	"moonbasic/runtime"
	mbcamera "moonbasic/runtime/camera"
	"moonbasic/vm/heap"
)

// Module implements GUI bindings for moonBASIC.
type Module struct {
	h   *heap.Store
	cam *mbcamera.Module
}

// NewModule constructs the GUI builtin module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// BindCamera implements runtime.CameraAware.
func (m *Module) BindCamera(c runtime.Module) {
	if cam, ok := c.(*mbcamera.Module); ok {
		m.cam = cam
	}
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

// Reset implements runtime.Module.
func (m *Module) Reset() {}

