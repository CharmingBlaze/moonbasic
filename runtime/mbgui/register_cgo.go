//go:build cgo

package mbgui

import "moonbasic/runtime"

// Register wires all GUI.* raygui builtins.
func (m *Module) Register(reg runtime.Registrar) {
	registerGlobalAndStyle(m, reg)
	registerThemeCommands(m, reg)
	registerLayout(m, reg)
	registerBasicControls(m, reg)
	registerTextControls(m, reg)
	registerSlidersAndLists(m, reg)
	registerColorAndDialogs(m, reg)
	registerTooltipIconsDraw(m, reg)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
