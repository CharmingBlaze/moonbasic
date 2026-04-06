//go:build !cgo && !windows

package window

import "moonbasic/runtime"

func (m *Module) registerWindowStateCommands(r runtime.Registrar) {
	stub := stubFn(stubHint)
	names := []string{
		"WINDOW.SETFLAG", "WINDOW.CLEARFLAG", "WINDOW.CHECKFLAG", "WINDOW.SETSTATE",
		"WINDOW.SETMINSIZE", "WINDOW.SETMAXSIZE",
		"WINDOW.GETPOSITIONX", "WINDOW.GETPOSITIONY",
		"WINDOW.SETMONITOR", "WINDOW.GETMONITORCOUNT", "WINDOW.GETMONITORNAME",
		"WINDOW.GETMONITORWIDTH", "WINDOW.GETMONITORHEIGHT", "WINDOW.GETMONITORREFRESHRATE",
		"WINDOW.GETSCALEDPIX", "WINDOW.GETSCALEDPIY", "WINDOW.SETICON", "WINDOW.SETOPACITY",
	}
	for _, n := range names {
		r.Register(n, "window", stub(n))
	}
}
