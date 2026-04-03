//go:build cgo

package mbcollision

import "moonbasic/runtime"

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	m.registerRayBuiltins(reg)
	m.registerBBoxBuiltins(reg)
	m.registerBSphereBuiltins(reg)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
