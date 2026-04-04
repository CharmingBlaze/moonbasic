//go:build !cgo

package input

import "moonbasic/runtime"

func (m *Module) registerMouseExtra(reg runtime.Registrar) {
	hint := "INPUT mouse/touch/gamepad requires CGO (Raylib)"
	for _, k := range []string{
		"INPUT.MOUSEPRESSED", "INPUT.MOUSERELEASED", "INPUT.MOUSEWHEELMOVE",
		"INPUT.MOUSEDELTAX", "INPUT.MOUSEDELTAY", "INPUT.SETMOUSEPOS", "INPUT.CHARPRESSED",
		"INPUT.ISGAMEPADAVAILABLE", "INPUT.GETGAMEPADAXISVALUE",
	} {
		reg.Register(k, "input", stubErr(hint, k))
	}
}
