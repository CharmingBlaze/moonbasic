//go:build cgo || (windows && !cgo)

package input

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Register wires Raylib keyboard builtins.
func (m *Module) Register(r runtime.Registrar) {
	registerCursor(r)
	registerGesture(r)
	m.registerInputAdvanced(r)
	m.registerMouseExtra(r)
	m.registerInputFacade(r)
	registerAxis(r)
	registerMovement2D(m, r)
	r.Register("INPUT.KEYDOWN", "input", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("INPUT.KEYDOWN expects 1 argument")
		}
		kc, err := KeyCodeFromValue(args[0])
		if err != nil {
			return value.Nil, err
		}
		return value.FromBool(rl.IsKeyDown(kc)), nil
	}))
	r.Register("INPUT.KEYPRESSED", "input", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("INPUT.KEYPRESSED expects 1 argument")
		}
		kc, err := KeyCodeFromValue(args[0])
		if err != nil {
			return value.Nil, err
		}
		return value.FromBool(rl.IsKeyPressed(kc)), nil
	}))
	r.Register("INPUT.KEYUP", "input", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("INPUT.KEYUP expects 1 argument")
		}
		kc, err := KeyCodeFromValue(args[0])
		if err != nil {
			return value.Nil, err
		}
		return value.FromBool(rl.IsKeyReleased(kc)), nil
	}))
	r.Register("INPUT.GETKEYNAME", "input", m.inGetKeyName)
	m.registerBlitzAliases(r)
	m.registerActionMapping(r)
}

func (m *Module) Shutdown() {}
