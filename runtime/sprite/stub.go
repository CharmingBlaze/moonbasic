//go:build !cgo

package mbsprite

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const hint = "SPRITE.* requires CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild"

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			return value.Nil, fmt.Errorf("%s: %s", name, hint)
		}
	}
	reg.Register("SPRITE.LOAD", "sprite", stub("SPRITE.LOAD"))
	reg.Register("SPRITE.DRAW", "sprite", stub("SPRITE.DRAW"))
	reg.Register("SPRITE.SETPOS", "sprite", stub("SPRITE.SETPOS"))
	reg.Register("SPRITE.SETPOSITION", "sprite", stub("SPRITE.SETPOSITION"))
	reg.Register("SPRITE.DEFANIM", "sprite", stub("SPRITE.DEFANIM"))
	reg.Register("SPRITE.PLAYANIM", "sprite", stub("SPRITE.PLAYANIM"))
	reg.Register("SPRITE.UPDATEANIM", "sprite", stub("SPRITE.UPDATEANIM"))
	reg.Register("SPRITE.HIT", "sprite", stub("SPRITE.HIT"))
	reg.Register("ATLAS.LOAD", "sprite", stub("ATLAS.LOAD"))
	reg.Register("ATLAS.FREE", "sprite", stub("ATLAS.FREE"))
	reg.Register("ATLAS.GETSPRITE", "sprite", stub("ATLAS.GETSPRITE"))
	reg.Register("ANIM.DEFINE", "sprite", stub("ANIM.DEFINE"))
	reg.Register("ANIM.ADDTRANSITION", "sprite", stub("ANIM.ADDTRANSITION"))
	reg.Register("ANIM.UPDATE", "sprite", stub("ANIM.UPDATE"))
	reg.Register("ANIM.SETPARAM", "sprite", stub("ANIM.SETPARAM"))
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
