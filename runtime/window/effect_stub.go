//go:build !cgo

package window

import "moonbasic/runtime"

func (m *Module) registerEffectCommands(r runtime.Registrar) {
	stub := stubFn(stubHint)
	for _, name := range []string{
		"EFFECT.SSAO", "EFFECT.SSR", "EFFECT.MOTIONBLUR", "EFFECT.DEPTHOFFIELD",
		"EFFECT.BLOOM", "EFFECT.TONEMAPPING", "EFFECT.SHARPEN", "EFFECT.GRAIN",
		"EFFECT.VIGNETTE", "EFFECT.CHROMATICABERRATION",
	} {
		r.Register(name, "post", stub(name))
	}
}
