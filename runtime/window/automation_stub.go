//go:build !cgo && !windows

package window

import "moonbasic/runtime"

func (m *Module) registerAutomationCommands(r runtime.Registrar) {
	stub := stubFn(stubHint)
	names := []string{
		"EVENT.LISTMAKE", "EVENT.LISTLOAD", "EVENT.LISTEXPORT", "EVENT.SETACTIVELIST",
		"EVENT.RECSTART", "EVENT.RECSTOP", "EVENT.REPLAY", "EVENT.RECPLAYING", "EVENT.ISPLAYING",
		"EVENT.LISTCLEAR", "EVENT.LISTCOUNT", "EVENT.LISTFREE",
	}
	for _, n := range names {
		r.Register(n, "window", stub(n))
	}
}

func (m *Module) shutdownAutomation() {}
