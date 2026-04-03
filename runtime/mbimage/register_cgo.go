//go:build cgo

package mbimage

import "moonbasic/runtime"

// Register wires IMAGE.* Raylib builtins.
func (m *Module) Register(reg runtime.Registrar) {
	registerImageLoad(m, reg)
	registerImageTransform(m, reg)
	registerImageDraw(m, reg)
	registerImageQuery(m, reg)
	registerClipboardImage(m, reg)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
