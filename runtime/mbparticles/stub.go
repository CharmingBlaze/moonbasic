//go:build !cgo

package mbparticles

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const hint = "PARTICLE natives require CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild"

var stubNames = []string{
	"PARTICLE.MAKE", "PARTICLE.FREE", "PARTICLE.SETTEXTURE", "PARTICLE.SETEMITRATE",
	"PARTICLE.SETLIFETIME", "PARTICLE.SETVELOCITY", "PARTICLE.SETCOLOR", "PARTICLE.SETCOLOREND",
	"PARTICLE.SETSIZE", "PARTICLE.SETGRAVITY", "PARTICLE.SETPOS", "PARTICLE.PLAY",
	"PARTICLE.UPDATE", "PARTICLE.DRAW",
}

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	for _, name := range stubNames {
		n := name
		reg.Register(n, "particle", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			return value.Nil, fmt.Errorf("%s: %s", n, hint)
		})
	}
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
