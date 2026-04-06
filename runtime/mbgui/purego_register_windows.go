//go:build !cgo && windows

package mbgui

import "moonbasic/runtime"

// Register wires a minimal Raylib-drawn GUI for Windows when CGO is disabled (purego Raylib).
// Full raygui behavior requires CGO_ENABLED=1.
func (m *Module) Register(reg runtime.Registrar) {
	registerPuregoGlobal(m, reg)
	registerPuregoTheme(m, reg)
	registerPuregoBasic(m, reg)
	registerPuregoSliders(m, reg)
	registerPuregoLayout(m, reg)
	registerPuregoTooltipDraw(m, reg)
	registerPuregoUnimplemented(reg)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
