//go:build cgo || (windows && !cgo)

package mbimage

import "moonbasic/runtime"

// Register wires IMAGE.* Raylib builtins.
func (m *Module) Register(reg runtime.Registrar) {
	registerImageLoad(m, reg)
	registerImageSequence(m, reg)
	registerImageTransform(m, reg)
	registerImageDraw(m, reg)
	registerImageProcess(m, reg)
	registerImageQuery(m, reg)
	registerPixelFilterCmds(m, reg)
	registerClipboardImage(m, reg)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
