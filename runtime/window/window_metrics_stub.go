//go:build !cgo && !windows

package window

import "moonbasic/runtime"

func (m *Module) registerWindowMetricsCommands(reg runtime.Registrar) {
	stub := stubFn(stubHint)
	for _, k := range []string{
		"WINDOW.WIDTH", "WINDOW.HEIGHT", "WINDOW.GETFPS",
		"WINDOW.ISFULLSCREEN", "WINDOW.TOGGLEFULLSCREEN", "WINDOW.ISRESIZED", "WINDOW.SETTITLE",
	} {
		reg.Register(k, "window", stub(k))
	}
}
