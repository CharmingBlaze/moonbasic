//go:build cgo || (windows && !cgo)

package mbgame

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (m *Module) registerGamepad(r runtime.Registrar) {
	_ = m
	r.Register("GAME.ISGAMEPADAVAILABLE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GAME.ISGAMEPADAVAILABLE expects 1 argument (gamepad#)")
		}
		gp, ok := argI(args[0])
		if !ok || gp < 0 {
			return value.Nil, fmt.Errorf("GAME.ISGAMEPADAVAILABLE: gamepad index must be non-negative")
		}
		return value.FromBool(rl.IsGamepadAvailable(int32(gp))), nil
	}))
	r.Register("GAME.GETGAMEPADNAME$", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GAME.GETGAMEPADNAME$ expects 1 argument (gamepad#)")
		}
		gp, ok := argI(args[0])
		if !ok || gp < 0 {
			return value.Nil, fmt.Errorf("GAME.GETGAMEPADNAME$: gamepad index must be non-negative")
		}
		if !rl.IsGamepadAvailable(int32(gp)) {
			return rt.RetString(""), nil
		}
		return rt.RetString(rl.GetGamepadName(int32(gp))), nil
	})
}
