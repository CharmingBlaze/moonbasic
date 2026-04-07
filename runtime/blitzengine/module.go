// Package blitzengine registers Blitz3D-style flat command names (APPTITLE, GRAPHICS, PLOT, …)
// that forward to the existing NAMESPACE.NAME surface. See docs/reference/moonbasic-command-set/.
package blitzengine

import (
	"moonbasic/runtime"
)

// Module holds optional 2D pen state for SETCOLOR / SETALPHA / SETORIGIN / SETVIEWPORT.
type Module struct {
	pen pen2D
}

// NewModule creates the Blitz-style facade module.
func NewModule() *Module {
	m := &Module{}
	m.pen.setAlpha(1)
	m.pen.setColor(255, 255, 255)
	return m
}

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	registerBlitzAPI(m, reg)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
