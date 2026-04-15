//go:build !cgo && !windows

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
	reg.Register("SPRITE.GETPOS", "sprite", stub("SPRITE.GETPOS"))
	reg.Register("SPRITE.DEFANIM", "sprite", stub("SPRITE.DEFANIM"))
	reg.Register("SPRITE.PLAYANIM", "sprite", stub("SPRITE.PLAYANIM"))
	reg.Register("SPRITE.UPDATEANIM", "sprite", stub("SPRITE.UPDATEANIM"))
	reg.Register("SPRITE.SETFRAME", "sprite", stub("SPRITE.SETFRAME"))
	reg.Register("SPRITE.PLAY", "sprite", stub("SPRITE.PLAY"))
	reg.Register("SPRITE.SETORIGIN", "sprite", stub("SPRITE.SETORIGIN"))
	reg.Register("SPRITE.GETSCALE", "sprite", stub("SPRITE.GETSCALE"))
	reg.Register("SPRITE.SETSCALE", "sprite", stub("SPRITE.SETSCALE"))
	reg.Register("SPRITE.GETROT", "sprite", stub("SPRITE.GETROT"))
	reg.Register("SPRITE.SETROT", "sprite", stub("SPRITE.SETROT"))
	reg.Register("SPRITE.GETCOLOR", "sprite", stub("SPRITE.GETCOLOR"))
	reg.Register("SPRITE.SETCOLOR", "sprite", stub("SPRITE.SETCOLOR"))
	reg.Register("SPRITE.GETALPHA", "sprite", stub("SPRITE.GETALPHA"))
	reg.Register("SPRITE.SETALPHA", "sprite", stub("SPRITE.SETALPHA"))
	reg.Register("SPRITE.HIT", "sprite", stub("SPRITE.HIT"))
	reg.Register("SPRITECOLLIDE", "sprite", stub("SPRITECOLLIDE"))
	reg.Register("SPRITE.POINTHIT", "sprite", stub("SPRITE.POINTHIT"))
	reg.Register("SPRITE.FREE", "sprite", stub("SPRITE.FREE"))
	reg.Register("ATLAS.LOAD", "sprite", stub("ATLAS.LOAD"))
	reg.Register("ATLAS.FREE", "sprite", stub("ATLAS.FREE"))
	reg.Register("ATLAS.GETSPRITE", "sprite", stub("ATLAS.GETSPRITE"))
	reg.Register("ANIM.DEFINE", "sprite", stub("ANIM.DEFINE"))
	reg.Register("ANIM.ADDTRANSITION", "sprite", stub("ANIM.ADDTRANSITION"))
	reg.Register("ANIM.UPDATE", "sprite", stub("ANIM.UPDATE"))
	reg.Register("ANIM.SETPARAM", "sprite", stub("ANIM.SETPARAM"))

	extras := []string{
		"SPRITEGROUP.CREATE", "SPRITEGROUP.MAKE", "SPRITEGROUP.ADD", "SPRITEGROUP.REMOVE", "SPRITEGROUP.CLEAR", "SPRITEGROUP.DRAW", "SPRITEGROUP.FREE",
		"SPRITELAYER.MAKE", "SPRITELAYER.CREATE", "SPRITELAYER.ADD", "SPRITELAYER.CLEAR", "SPRITELAYER.SETZ", "SPRITELAYER.DRAW", "SPRITELAYER.FREE",
		"SPRITEBATCH.MAKE", "SPRITEBATCH.CREATE", "SPRITEBATCH.ADD", "SPRITEBATCH.CLEAR", "SPRITEBATCH.DRAW", "SPRITEBATCH.FREE",
		"SPRITEUI.CREATE", "SPRITEUI.MAKE", "SPRITEUI.DRAW", "SPRITEUI.FREE",
		"PARTICLE2D.CREATE", "PARTICLE2D.MAKE", "PARTICLE2D.EMIT", "PARTICLE2D.UPDATE", "PARTICLE2D.DRAW", "PARTICLE2D.FREE",
	}
	for _, name := range extras {
		n := name
		reg.Register(n, "sprite", stub(n))
	}
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
