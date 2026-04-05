//go:build !cgo

package mbgame

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerShortcuts(r runtime.Registrar) {
	m.registerPauseFrame(r)

	end := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("ENDGAME expects 0 arguments")
		}
		if rt != nil && rt.TerminateVM != nil {
			rt.TerminateVM()
		}
		return value.Nil, nil
	}
	r.Register("ENDGAME", "game", end)
	r.Register("GAME.ENDGAME", "game", end)

	err := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			_ = args
			return value.Nil, fmt.Errorf("%s requires CGO (raylib)", name)
		}
	}
	names := []string{
		"SCREENW", "SCREENH", "SCREENCX", "SCREENCY", "DT", "FPS",
		"MX", "MY", "MOUSEX", "MOUSEY", "MDX", "MDY", "MWHEEL",
		"MLEFT", "MRIGHT", "MMIDDLE", "MLEFTPRESSED", "MRIGHTPRESSED",
		"KEYDOWN", "KEYPRESSED", "KEYRELEASED", "KEYCHAR", "ANYKEY",
	}
	for _, n := range names {
		r.Register(n, "game", err(n))
		r.Register("GAME."+n, "game", err("GAME."+n))
	}
}
