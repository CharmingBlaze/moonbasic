//go:build !cgo && !windows

package input

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerBlitzMouseKeys(r runtime.Registrar) {
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			_ = args
			return value.Nil, fmt.Errorf("%s requires CGO (raylib)", name)
		}
	}
	for _, n := range []string{
		"MOUSEDOWN", "MOUSEHIT", "MOUSEXSPEED", "MOUSEYSPEED",
		"FlushMouse", "WaitMouse", "MoveMouse", "HidePointer", "ShowPointer",
		"GetKey", "WaitKey", "FlushKeys",
	} {
		r.Register(n, "input", stub(n))
	}
}
