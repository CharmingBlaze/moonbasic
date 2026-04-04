//go:build !cgo

package window

import "moonbasic/runtime"

func (m *Module) registerDecalCommands(r runtime.Registrar) {
	stub := stubFn(stubHint)
	for _, name := range []string{
		"DECAL.MAKE", "DECAL.FREE", "DECAL.SETPOS", "DECAL.SETSIZE",
		"DECAL.SETLIFETIME", "DECAL.DRAW",
	} {
		r.Register(name, "decal", stub(name))
	}
}
