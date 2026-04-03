//go:build !cgo

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
}

func (m *Module) Shutdown() {}
