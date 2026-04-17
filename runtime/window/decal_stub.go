//go:build !cgo && !windows

package window

import "moonbasic/runtime"

func (m *Module) registerDecalCommands(r runtime.Registrar) {
	stub := stubFn(stubHint)
	for _, name := range []string{
		"DECAL.CREATE", "DECAL.MAKE", "DECAL.FREE", "DECAL.SETPOS", "DECAL.GETPOS",
		"DECAL.SETSIZE", "DECAL.GETSIZE", "DECAL.SETLIFETIME", "DECAL.GETLIFETIME",
		"DECAL.GETROT", "DECAL.DRAW",
	} {
		r.Register(name, "decal", stub(name))
	}
}
