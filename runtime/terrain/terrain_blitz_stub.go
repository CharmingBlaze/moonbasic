//go:build (!cgo && !windows) || (windows && gopls_stub)

package terrain

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerTerrainBlitzAliases(m *Module, r runtime.Registrar) {
	_ = m
	hint := func(name string) func(*runtime.Runtime, ...value.Value) (value.Value, error) {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			_ = args
			return value.Nil, fmt.Errorf("%s requires CGO and Raylib (set CGO_ENABLED=1)", name)
		}
	}
	for _, k := range []string{
		"TerrainHeight", "ModifyTerrain", "TerrainX", "TerrainZ", "TerrainSize",
		"LoadTerrain", "TerrainDetail", "TerrainShading",
	} {
		r.Register(k, "terrain", hint(k))
	}
}
