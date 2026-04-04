//go:build !cgo

package mbtilemap

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const hint = "TILEMAP.* requires CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild"

var stubNames = []string{
	"TILEMAP.LOAD", "TILEMAP.FREE", "TILEMAP.SETTILESIZE", "TILEMAP.DRAW",
	"TILEMAP.GETTILE", "TILEMAP.SETTILE", "TILEMAP.ISSOLID", "TILEMAP.WIDTH", "TILEMAP.HEIGHT",
	"TILEMAP.LAYERCOUNT", "TILEMAP.LAYERNAME", "TILEMAP.DRAWLAYER",
	"TILEMAP.COLLISIONAT", "TILEMAP.SETCOLLISION", "TILEMAP.MERGECOLLISIONLAYER", "TILEMAP.ISSOLIDCATEGORY",
}

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	for _, name := range stubNames {
		n := name
		reg.Register(n, "tilemap", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			return value.Nil, fmt.Errorf("%s: %s", n, hint)
		})
	}
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
