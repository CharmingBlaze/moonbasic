//go:build !cgo && !windows

package biome

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerBiome(m *Module, r runtime.Registrar) {
	h := func(n string) func(*runtime.Runtime, ...value.Value) (value.Value, error) {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			return value.Nil, fmt.Errorf("%s requires CGO", n)
		}
	}
	r.Register("BIOME.CREATE", "biome", h("BIOME.CREATE"))
	r.Register("BIOME.MAKE", "biome", h("BIOME.MAKE"))
	r.Register("BIOME.FREE", "biome", h("BIOME.FREE"))
	r.Register("BIOME.SETTEMP", "biome", h("BIOME.SETTEMP"))
	r.Register("BIOME.SETHUMIDITY", "biome", h("BIOME.SETHUMIDITY"))
}
