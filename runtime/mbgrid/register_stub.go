//go:build !cgo && !windows

package mbgrid

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const hint = "GRID.* requires CGO (Raylib)"

func (m *Module) Register(reg runtime.Registrar) {
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			_ = args
			return value.Nil, fmt.Errorf("%s: %s", name, hint)
		}
	}
	for _, n := range []string{
		"GRID.CREATE", "GRID.MAKE", "GRID.FREE", "GRID.SETCELL", "GRID.GETCELL", "GRID.WORLDTOCELL",
		"GRID.DRAW", "GRID.SNAP", "GRID.GETPATH", "GRID.FOLLOWTERRAIN", "GRID.PLACEENTITY",
		"GRID.RAYCAST", "GRID.GETNEIGHBORS",
	} {
		reg.Register(n, "grid", stub(n))
	}
}
