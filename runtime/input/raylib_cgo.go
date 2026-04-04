//go:build cgo

package input

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func keyCodeArg(v value.Value) (int32, error) {
	if i, ok := v.ToInt(); ok {
		return int32(i), nil
	}
	if f, ok := v.ToFloat(); ok {
		return int32(f), nil
	}
	return 0, fmt.Errorf("expected numeric key code (use KEY_ESCAPE etc.)")
}

// Register wires Raylib keyboard builtins.
func (m *Module) Register(r runtime.Registrar) {
	registerCursor(r)
	registerGesture(r)
	m.registerInputAdvanced(r)
	m.registerMouseExtra(r)
	r.Register("INPUT.KEYDOWN", "input", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("INPUT.KEYDOWN expects 1 argument")
		}
		kc, err := keyCodeArg(args[0])
		if err != nil {
			return value.Nil, err
		}
		return value.FromBool(rl.IsKeyDown(kc)), nil
	}))
	r.Register("INPUT.KEYPRESSED", "input", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("INPUT.KEYPRESSED expects 1 argument")
		}
		kc, err := keyCodeArg(args[0])
		if err != nil {
			return value.Nil, err
		}
		return value.FromBool(rl.IsKeyPressed(kc)), nil
	}))
	r.Register("INPUT.KEYUP", "input", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("INPUT.KEYUP expects 1 argument")
		}
		kc, err := keyCodeArg(args[0])
		if err != nil {
			return value.Nil, err
		}
		return value.FromBool(rl.IsKeyReleased(kc)), nil
	}))
	r.Register("INPUT.GETKEYNAME", "input", m.inGetKeyName)
	m.registerActionMapping(r)
}

func (m *Module) Shutdown() {}
