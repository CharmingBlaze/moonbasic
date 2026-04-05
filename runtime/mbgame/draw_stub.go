//go:build !cgo

package mbgame

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func stubDraw(name string) runtime.BuiltinFn {
	return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		_ = args
		return value.Nil, fmt.Errorf("%s requires CGO_ENABLED=1 (raylib)", name)
	}
}

func (m *Module) registerDrawHelpers(r runtime.Registrar) {
	_ = m
	r.Register("GAME.DRAWSCREENFLASH", "game", stubDraw("GAME.DRAWSCREENFLASH"))
	r.Register("GAME.SCREENFLASH", "game", stubDraw("GAME.SCREENFLASH"))
}
