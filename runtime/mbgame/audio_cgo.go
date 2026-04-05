//go:build cgo

package mbgame

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (m *Module) registerAudioHelpers(r runtime.Registrar) {
	_ = m
	r.Register("GAME.SETMASTERVOLUME", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GAME.SETMASTERVOLUME expects 1 argument (volume# 0..1)")
		}
		v, ok := argF(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("GAME.SETMASTERVOLUME: numeric volume")
		}
		if v < 0 {
			v = 0
		}
		if v > 1 {
			v = 1
		}
		rl.SetMasterVolume(float32(v))
		return value.Nil, nil
	}))
	r.Register("GAME.GETMASTERVOLUME", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GAME.GETMASTERVOLUME expects 0 arguments")
		}
		return value.FromFloat(float64(rl.GetMasterVolume())), nil
	}))
}
