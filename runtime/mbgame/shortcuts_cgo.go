//go:build cgo || (windows && !cgo)

package mbgame

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/runtime/input"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (m *Module) registerShortcuts(r runtime.Registrar) {
	m.registerPauseFrame(r)

	regLegacy2 := func(short, dotted string, impl func([]value.Value) (value.Value, error)) {
		w := runtime.AdaptLegacy(impl)
		r.Register(short, "game", w)
		r.Register(dotted, "game", w)
	}
	regRT0 := func(short, dotted string, impl func(*runtime.Runtime, ...value.Value) (value.Value, error)) {
		r.Register(short, "game", impl)
		r.Register(dotted, "game", impl)
	}

	regLegacy2("SCREENW", "GAME.SCREENW", func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("SCREENW expects 0 arguments")
		}
		return value.FromInt(int64(rl.GetScreenWidth())), nil
	})
	regLegacy2("SCREENH", "GAME.SCREENH", func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("SCREENH expects 0 arguments")
		}
		return value.FromInt(int64(rl.GetScreenHeight())), nil
	})
	regLegacy2("SCREENCX", "GAME.SCREENCX", func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("SCREENCX expects 0 arguments")
		}
		return value.FromFloat(float64(rl.GetScreenWidth()) * 0.5), nil
	})
	regLegacy2("SCREENCY", "GAME.SCREENCY", func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("SCREENCY expects 0 arguments")
		}
		return value.FromFloat(float64(rl.GetScreenHeight()) * 0.5), nil
	})

	regRT0("DT", "GAME.DT", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("DT expects 0 arguments")
		}
		if rt != nil && rt.GamePaused {
			return value.FromFloat(0), nil
		}
		return value.FromFloat(float64(rl.GetFrameTime())), nil
	})
	regLegacy2("FPS", "GAME.FPS", func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("FPS expects 0 arguments")
		}
		return value.FromInt(int64(rl.GetFPS())), nil
	})

	regRT0("ENDGAME", "GAME.ENDGAME", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("ENDGAME expects 0 arguments")
		}
		if rt != nil && rt.TerminateVM != nil {
			rt.TerminateVM()
		}
		return value.Nil, nil
	})

	regLegacy2("MX", "GAME.MX", func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("MX expects 0 arguments")
		}
		return value.FromInt(int64(rl.GetMouseX())), nil
	})
	regLegacy2("MY", "GAME.MY", func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("MY expects 0 arguments")
		}
		return value.FromInt(int64(rl.GetMouseY())), nil
	})
	regLegacy2("MOUSEX", "GAME.MOUSEX", func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("MOUSEX expects 0 arguments")
		}
		return value.FromInt(int64(rl.GetMouseX())), nil
	})
	regLegacy2("MOUSEY", "GAME.MOUSEY", func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("MOUSEY expects 0 arguments")
		}
		return value.FromInt(int64(rl.GetMouseY())), nil
	})
	regLegacy2("MDX", "GAME.MDX", func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("MDX expects 0 arguments")
		}
		d := rl.GetMouseDelta()
		return value.FromFloat(float64(d.X)), nil
	})
	regLegacy2("MDY", "GAME.MDY", func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("MDY expects 0 arguments")
		}
		d := rl.GetMouseDelta()
		return value.FromFloat(float64(d.Y)), nil
	})
	regLegacy2("MWHEEL", "GAME.MWHEEL", func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("MWHEEL expects 0 arguments")
		}
		return value.FromFloat(float64(rl.GetMouseWheelMove())), nil
	})
	regLegacy2("MLEFT", "GAME.MLEFT", func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("MLEFT expects 0 arguments")
		}
		return value.FromBool(rl.IsMouseButtonDown(rl.MouseLeftButton)), nil
	})
	regLegacy2("MRIGHT", "GAME.MRIGHT", func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("MRIGHT expects 0 arguments")
		}
		return value.FromBool(rl.IsMouseButtonDown(rl.MouseRightButton)), nil
	})
	regLegacy2("MMIDDLE", "GAME.MMIDDLE", func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("MMIDDLE expects 0 arguments")
		}
		return value.FromBool(rl.IsMouseButtonDown(rl.MouseMiddleButton)), nil
	})
	regLegacy2("MLEFTPRESSED", "GAME.MLEFTPRESSED", func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("MLEFTPRESSED expects 0 arguments")
		}
		return value.FromBool(rl.IsMouseButtonPressed(rl.MouseLeftButton)), nil
	})
	regLegacy2("MRIGHTPRESSED", "GAME.MRIGHTPRESSED", func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("MRIGHTPRESSED expects 0 arguments")
		}
		return value.FromBool(rl.IsMouseButtonPressed(rl.MouseRightButton)), nil
	})

	regLegacy2("KEYDOWN", "GAME.KEYDOWN", func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("KEYDOWN expects 1 argument")
		}
		kc, err := input.KeyCodeFromValue(args[0])
		if err != nil {
			return value.Nil, err
		}
		return value.FromBool(rl.IsKeyDown(kc)), nil
	})
	regLegacy2("KEYPRESSED", "GAME.KEYPRESSED", func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("KEYPRESSED expects 1 argument")
		}
		kc, err := input.KeyCodeFromValue(args[0])
		if err != nil {
			return value.Nil, err
		}
		return value.FromBool(rl.IsKeyPressed(kc)), nil
	})
	regLegacy2("KEYRELEASED", "GAME.KEYRELEASED", func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("KEYRELEASED expects 1 argument")
		}
		kc, err := input.KeyCodeFromValue(args[0])
		if err != nil {
			return value.Nil, err
		}
		return value.FromBool(rl.IsKeyReleased(kc)), nil
	})
	regLegacy2("KEYCHAR", "GAME.KEYCHAR", func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("KEYCHAR expects 0 arguments")
		}
		c := rl.GetCharPressed()
		return value.FromInt(int64(c)), nil
	})
	regLegacy2("ANYKEY", "GAME.ANYKEY", func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("ANYKEY expects 0 arguments")
		}
		for k := int32(32); k <= 348; k++ {
			if rl.IsKeyPressed(k) {
				return value.FromBool(true), nil
			}
		}
		return value.FromBool(false), nil
	})
}
