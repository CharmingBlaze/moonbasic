//go:build cgo

package input

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerMouseExtra(reg runtime.Registrar) {
	reg.Register("INPUT.MOUSEPRESSED", "input", runtime.AdaptLegacy(m.inMousePressed))
	reg.Register("INPUT.MOUSERELEASED", "input", runtime.AdaptLegacy(m.inMouseReleased))
	reg.Register("INPUT.MOUSEWHEELMOVE", "input", runtime.AdaptLegacy(m.inMouseWheelMove))
	reg.Register("INPUT.MOUSEDELTAX", "input", runtime.AdaptLegacy(m.inMouseDeltaX))
	reg.Register("INPUT.MOUSEDELTAY", "input", runtime.AdaptLegacy(m.inMouseDeltaY))
	reg.Register("INPUT.SETMOUSEPOS", "input", runtime.AdaptLegacy(m.inSetMousePos))
	reg.Register("INPUT.CHARPRESSED", "input", runtime.AdaptLegacy(m.inCharPressed))
	reg.Register("INPUT.ISGAMEPADAVAILABLE", "input", runtime.AdaptLegacy(m.inIsGamepadAvailable))
	reg.Register("INPUT.GETGAMEPADAXISVALUE", "input", runtime.AdaptLegacy(m.inGetGamepadAxisValue))
}

func (m *Module) inMousePressed(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("INPUT.MOUSEPRESSED expects 1 argument (button)")
	}
	b, err := argMouseButton(args[0])
	if err != nil {
		return value.Nil, fmt.Errorf("INPUT.MOUSEPRESSED: %w", err)
	}
	return value.FromBool(rl.IsMouseButtonPressed(b)), nil
}

func (m *Module) inMouseReleased(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("INPUT.MOUSERELEASED expects 1 argument (button)")
	}
	b, err := argMouseButton(args[0])
	if err != nil {
		return value.Nil, fmt.Errorf("INPUT.MOUSERELEASED: %w", err)
	}
	return value.FromBool(rl.IsMouseButtonReleased(b)), nil
}

func (m *Module) inMouseWheelMove(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("INPUT.MOUSEWHEELMOVE expects 0 arguments")
	}
	return value.FromFloat(float64(rl.GetMouseWheelMove())), nil
}

func (m *Module) inMouseDeltaX(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("INPUT.MOUSEDELTAX expects 0 arguments")
	}
	d := rl.GetMouseDelta()
	return value.FromFloat(float64(d.X)), nil
}

func (m *Module) inMouseDeltaY(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("INPUT.MOUSEDELTAY expects 0 arguments")
	}
	d := rl.GetMouseDelta()
	return value.FromFloat(float64(d.Y)), nil
}

func (m *Module) inSetMousePos(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("INPUT.SETMOUSEPOS expects (x, y)")
	}
	x, ok1 := args[0].ToInt()
	if !ok1 {
		if xf, ok := args[0].ToFloat(); ok {
			x = int64(xf)
			ok1 = true
		}
	}
	y, ok2 := args[1].ToInt()
	if !ok2 {
		if yf, ok := args[1].ToFloat(); ok {
			y = int64(yf)
			ok2 = true
		}
	}
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("INPUT.SETMOUSEPOS: x,y must be numeric")
	}
	rl.SetMousePosition(int(x), int(y))
	return value.Nil, nil
}

func (m *Module) inCharPressed(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("INPUT.CHARPRESSED expects 0 arguments")
	}
	c := rl.GetCharPressed()
	if c == 0 {
		return value.FromInt(0), nil
	}
	return value.FromInt(int64(c)), nil
}

func (m *Module) inIsGamepadAvailable(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("INPUT.ISGAMEPADAVAILABLE expects gamepad index")
	}
	gp, ok := args[0].ToInt()
	if !ok {
		if f, okf := args[0].ToFloat(); okf {
			gp = int64(f)
			ok = true
		}
	}
	if !ok || gp < 0 {
		return value.Nil, fmt.Errorf("INPUT.ISGAMEPADAVAILABLE: index must be numeric")
	}
	return value.FromBool(rl.IsGamepadAvailable(int32(gp))), nil
}

func (m *Module) inGetGamepadAxisValue(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("INPUT.GETGAMEPADAXISVALUE expects (gamepad, axis)")
	}
	gp, ok1 := args[0].ToInt()
	if !ok1 {
		if f, okf := args[0].ToFloat(); okf {
			gp = int64(f)
			ok1 = true
		}
	}
	ax, ok2 := args[1].ToInt()
	if !ok2 {
		if f, okf := args[1].ToFloat(); okf {
			ax = int64(f)
			ok2 = true
		}
	}
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("INPUT.GETGAMEPADAXISVALUE: gamepad and axis must be numeric")
	}
	v := rl.GetGamepadAxisMovement(int32(gp), int32(ax))
	return value.FromFloat(float64(v)), nil
}
