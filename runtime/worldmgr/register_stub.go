//go:build !cgo && !windows

package worldmgr

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerWorld(m *Module, r runtime.Registrar) {
	hint := func(name string) func(*runtime.Runtime, ...value.Value) (value.Value, error) {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			return value.Nil, fmt.Errorf("%s requires CGO", name)
		}
	}
	r.Register("WORLD.SETCENTER", "world", hint("WORLD.SETCENTER"))
	r.Register("WORLD.SETCENTERENTITY", "world", hint("WORLD.SETCENTERENTITY"))
	r.Register("WORLD.UPDATE", "world", hint("WORLD.UPDATE"))
	r.Register("WORLD.STREAMENABLE", "world", hint("WORLD.STREAMENABLE"))
	r.Register("WORLD.PRELOAD", "world", hint("WORLD.PRELOAD"))
	r.Register("WORLD.STATUS", "world", hint("WORLD.STATUS"))
	r.Register("WORLD.ISREADY", "world", hint("WORLD.ISREADY"))
	r.Register("WORLD.SETREFLECTION", "world", hint("WORLD.SETREFLECTION"))
	r.Register("WORLD.SETVEGETATION", "world", hint("WORLD.SETVEGETATION"))
}
