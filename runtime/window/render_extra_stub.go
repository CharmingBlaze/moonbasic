//go:build !cgo && !windows

package window

import "moonbasic/runtime"

func (m *Module) registerRenderAdvanced(r runtime.Registrar) {
	stub := stubFn(stubHint)
	for _, name := range []string{
		"RENDER.SETBLEND", "RENDER.SETBLENDMODE", "RENDER.SETDEPTHWRITE", "RENDER.SETDEPTHMASK",
		"RENDER.SETDEPTHTEST", "RENDER.SETSCISSOR", "RENDER.CLEARSCISSOR", "RENDER.SETWIREFRAME",
		"RENDER.SCREENSHOT", "RENDER.SETMSAA", "RENDER.SETSHADOWMAPSIZE", "RENDER.SETAMBIENT", "RENDER.SETMODE",
	} {
		r.Register(name, "window", stub(name))
	}
}
