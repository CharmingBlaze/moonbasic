//go:build !cgo && !windows

package mbdebug

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerDebugDraw3D(r runtime.Registrar) {
	_ = m
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			_ = args
			return value.Nil, fmt.Errorf("%s requires CGO (Raylib)", name)
		}
	}
	r.Register("DEBUG.DRAWLINE", "debug", stub("DEBUG.DRAWLINE"))
	r.Register("DEBUG.DRAWBOX", "debug", stub("DEBUG.DRAWBOX"))
}
