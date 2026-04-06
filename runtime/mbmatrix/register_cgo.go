//go:build cgo || (windows && !cgo)

package mbmatrix

import "moonbasic/runtime"

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	m.registerMat4(reg)
	m.registerVec3(reg)
	m.registerVec3Extras(reg)
	m.registerVec2(reg)
	m.registerVec2Extras(reg)
	m.registerQuat(reg)
	m.registerQuatExtras(reg)
	m.registerColor(reg)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
