//go:build !linux || !cgo

package mbcharcontroller

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const stubHint = "CHARCONTROLLER requires Linux with CGO and Jolt (same as PHYSICS3D)."

func registerCharControllerCommands(m *Module, reg runtime.Registrar) {
	_ = m
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			return value.Nil, fmt.Errorf("%s: %s", name, stubHint)
		}
	}
	for _, k := range []string{
		"CHARCONTROLLER.MAKE", "CHARCONTROLLER.SETPOS", "CHARCONTROLLER.SETPOSITION", "CHARCONTROLLER.GETPOS",
		"CHARCONTROLLER.MOVE", "CHARCONTROLLER.ISGROUNDED",
		"CHARCONTROLLER.X", "CHARCONTROLLER.Y", "CHARCONTROLLER.Z",
		"CHARCONTROLLER.FREE",
	} {
		reg.Register(k, "charcontroller", stub(k))
	}
}

func shutdownCharController(m *Module) { _ = m }
