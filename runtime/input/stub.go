//go:build !cgo && !windows

package input

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) Register(r runtime.Registrar) {
	registerCursor(r)
	registerGesture(r)
	m.registerInputAdvanced(r)
	m.registerMouseExtra(r)
	registerAxis(r)
	registerMovement2D(m, r)
	// Without CGO, keys are never "down" so WHILE NOT Input.KeyDown(...) keeps looping
	// until the user closes the window (window package stub still errors on OPEN).
	r.Register("INPUT.KEYDOWN", "input", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("INPUT.KEYDOWN expects 1 argument")
		}
		return value.FromBool(false), nil
	}))
	r.Register("INPUT.KEYPRESSED", "input", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("INPUT.KEYPRESSED expects 1 argument")
		}
		return value.FromBool(false), nil
	}))
	r.Register("INPUT.KEYUP", "input", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("INPUT.KEYUP expects 1 argument")
		}
		return value.FromBool(false), nil
	}))
	r.Register("INPUT.GETKEYNAME", "input", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("INPUT.GETKEYNAME expects 1 argument")
		}
		return rt.RetString(""), nil
	})
	m.registerBlitzAliases(r)
	m.registerActionMapping(r)
	m.registerInputFacadeStub(r)
}

func (m *Module) registerInputFacadeStub(r runtime.Registrar) {
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			_ = args
			return value.Nil, fmt.Errorf("%s requires CGO_ENABLED=1", name)
		}
	}
	for _, name := range []string{
		"MOUSE", "MOUSE.DX", "MOUSE.DY", "MOUSE.WHEEL", "MOUSE.DOWN", "MOUSE.PRESSED", "MOUSE.RELEASED",
		"KEY", "KEY.DOWN", "KEY.HIT", "KEY.UP",
		"GAMEPAD", "GAMEPAD.AXIS", "GAMEPAD.BUTTON",
	} {
		r.Register(name, "input", stub(name))
	}
}

func (m *Module) Shutdown() {}
