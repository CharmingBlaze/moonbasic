//go:build !cgo && !windows

package window

import "moonbasic/runtime"

func (m *Module) registerWindowPlacementCommands(reg runtime.Registrar) {
	stub := stubFn(stubHint)
	for _, k := range []string{
		"WINDOW.SETPOS", "WINDOW.SETPOSITION", "WINDOW.SETSIZE",
		"WINDOW.MINIMIZE", "WINDOW.MAXIMIZE", "WINDOW.RESTORE",
		"WINDOW.SETTARGETFPS",
	} {
		reg.Register(k, "window", stub(k))
	}
}
