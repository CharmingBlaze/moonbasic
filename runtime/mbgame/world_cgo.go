//go:build cgo || (windows && !cgo)

package mbgame

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (m *Module) registerWorldCamera(r runtime.Registrar) {
	_ = m
	r.Register("GAME.ISCURSORONSCREEN", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GAME.ISCURSORONSCREEN expects 0 arguments")
		}
		return value.FromBool(rl.IsCursorOnScreen()), nil
	}))
}
