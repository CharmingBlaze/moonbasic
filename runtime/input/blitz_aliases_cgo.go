//go:build cgo || (windows && !cgo)

package input

import (
	"fmt"

	mbruntime "moonbasic/runtime"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (m *Module) registerBlitzAliases(r mbruntime.Registrar) {
	r.Register("INPUT.KEYHIT", "input", mbruntime.AdaptLegacy(m.inKeyHit))
	r.Register("INPUT.MOUSEXSPEED", "input", mbruntime.AdaptLegacy(m.inMouseXSpeed))
	r.Register("INPUT.MOUSEYSPEED", "input", mbruntime.AdaptLegacy(m.inMouseYSpeed))
	r.Register("INPUT.JOYX", "input", mbruntime.AdaptLegacy(m.inJoyX))
	r.Register("INPUT.JOYY", "input", mbruntime.AdaptLegacy(m.inJoyY))
	r.Register("INPUT.JOYBUTTON", "input", mbruntime.AdaptLegacy(m.inJoyButton))
	r.Register("INPUT.JOYDOWN", "input", mbruntime.AdaptLegacy(m.inJoyButton))

	// Easy Mode Global Aliases (and namespaces)
	r.Register("MOUSEDX", "input", mbruntime.AdaptLegacy(m.inMouseXSpeed))
	r.Register("INPUT.MOUSEDX", "input", mbruntime.AdaptLegacy(m.inMouseXSpeed))
	r.Register("MOUSEDY", "input", mbruntime.AdaptLegacy(m.inMouseYSpeed))
	r.Register("INPUT.MOUSEDY", "input", mbruntime.AdaptLegacy(m.inMouseYSpeed))
	r.Register("MOUSEWHEEL", "input", mbruntime.AdaptLegacy(m.inMouseWheel))
	r.Register("INPUT.MOUSEWHEEL", "input", mbruntime.AdaptLegacy(m.inMouseWheel))
	r.Register("MOUSEX", "input", mbruntime.AdaptLegacy(m.inMouseX))
	r.Register("MOUSEY", "input", mbruntime.AdaptLegacy(m.inMouseY))
	r.Register("MOUSEZ", "input", mbruntime.AdaptLegacy(m.inMouseWheel))
	r.Register("KEYHIT", "input", mbruntime.AdaptLegacy(m.inKeyHit))
	r.Register("KEYDOWN", "input", mbruntime.AdaptLegacy(m.inKeyDown))
	r.Register("INPUT.KEYDOWN", "input", mbruntime.AdaptLegacy(m.inKeyDown))
	r.Register("KeyDown", "input", mbruntime.AdaptLegacy(m.inKeyDown))
	r.Register("KEYUP", "input", mbruntime.AdaptLegacy(m.inKeyUp))
	r.Register("INPUT.KEYUP", "input", mbruntime.AdaptLegacy(m.inKeyUp))
	r.Register("AXIS", "input", mbruntime.AdaptLegacy(m.inAxis))
}

func (m *Module) inKeyHit(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("INPUT.KEYHIT expects 1 argument (key)")
	}
	kc, err := KeyCodeFromValue(args[0])
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(rl.IsKeyPressed(kc)), nil
}

func (m *Module) inMouseXSpeed(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("INPUT.MOUSEXSPEED expects 0 arguments")
	}
	d := rl.GetMouseDelta()
	return value.FromFloat(float64(d.X)), nil
}

func (m *Module) inMouseYSpeed(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("INPUT.MOUSEYSPEED expects 0 arguments")
	}
	d := rl.GetMouseDelta()
	return value.FromFloat(float64(d.Y)), nil
}

func (m *Module) inJoyX(args []value.Value) (value.Value, error) {
	gp, ax, err := joyParseArgs(args, rl.GamepadAxisLeftX)
	if err != nil {
		return value.Nil, err
	}
	if !rl.IsGamepadAvailable(gp) {
		return value.FromFloat(0), nil
	}
	v := rl.GetGamepadAxisMovement(gp, ax)
	return value.FromFloat(float64(v)), nil
}

func (m *Module) inJoyY(args []value.Value) (value.Value, error) {
	gp, ax, err := joyParseArgs(args, rl.GamepadAxisLeftY)
	if err != nil {
		return value.Nil, err
	}
	if !rl.IsGamepadAvailable(gp) {
		return value.FromFloat(0), nil
	}
	v := rl.GetGamepadAxisMovement(gp, ax)
	return value.FromFloat(float64(v)), nil
}

func joyParseArgs(args []value.Value, defaultAxis int32) (gp int32, ax int32, err error) {
	switch len(args) {
	case 0:
		return 0, defaultAxis, nil
	case 1:
		i, ok := args[0].ToInt()
		if !ok || i < 0 {
			return 0, 0, fmt.Errorf("INPUT.JOYX/JOYY: gamepad index must be non-negative int")
		}
		return int32(i), defaultAxis, nil
	case 2:
		i, ok := args[0].ToInt()
		if !ok || i < 0 {
			return 0, 0, fmt.Errorf("INPUT.JOYX/JOYY: gamepad index must be non-negative int")
		}
		a, ok2 := args[1].ToInt()
		if !ok2 || a < 0 {
			return 0, 0, fmt.Errorf("INPUT.JOYX/JOYY: axis index must be non-negative int")
		}
		return int32(i), int32(a), nil
	default:
		return 0, 0, fmt.Errorf("INPUT.JOYX/JOYY: expected 0–2 arguments (gamepad#, axis#)")
	}
}

func (m *Module) inJoyButton(args []value.Value) (value.Value, error) {
	var gp, btn int32
	switch len(args) {
	case 1:
		b, ok := args[0].ToInt()
		if !ok || b < 0 {
			return value.Nil, fmt.Errorf("INPUT.JOYBUTTON: button index must be non-negative int")
		}
		gp, btn = 0, int32(b)
	case 2:
		g, ok1 := args[0].ToInt()
		b, ok2 := args[1].ToInt()
		if !ok1 || g < 0 || !ok2 || b < 0 {
			return value.Nil, fmt.Errorf("INPUT.JOYBUTTON: (gamepad#, button#) must be non-negative ints")
		}
		gp, btn = int32(g), int32(b)
	default:
		return value.Nil, fmt.Errorf("INPUT.JOYBUTTON expects 1 or 2 arguments (button#) or (gamepad#, button#)")
	}
	if !rl.IsGamepadAvailable(gp) {
		return value.FromBool(false), nil
	}
	return value.FromBool(rl.IsGamepadButtonDown(gp, btn)), nil
}
func (m *Module) inMouseWheel(args []value.Value) (value.Value, error) {
	return value.FromFloat(float64(rl.GetMouseWheelMove())), nil
}

func (m *Module) inKeyDown(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("KEYDOWN expects 1 argument (key)")
	}
	kc, err := KeyCodeFromValue(args[0])
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(rl.IsKeyDown(kc)), nil
}

func (m *Module) inKeyUp(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("KEYUP expects 1 argument (key)")
	}
	kc, err := KeyCodeFromValue(args[0])
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(rl.IsKeyUp(kc)), nil
}

func (m *Module) inAxis(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("AXIS expects 2 arguments (negKey, posKey)")
	}
	k1, err1 := KeyCodeFromValue(args[0])
	k2, err2 := KeyCodeFromValue(args[1])
	if err1 != nil || err2 != nil {
		return value.Nil, fmt.Errorf("AXIS: invalid key codes")
	}
	v := 0.0
	if rl.IsKeyDown(k1) {
		v -= 1.0
	}
	if rl.IsKeyDown(k2) {
		v += 1.0
	}
	return value.FromFloat(v), nil
}
