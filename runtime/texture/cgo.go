//go:build cgo || (windows && !cgo)

package texture

import "moonbasic/runtime"

func (m *Module) Register(r runtime.Registrar) {
	registerTextureLoadCmds(m, r)
	registerTexturePropCmds(m, r)
	registerTextureGenCmds(m, r)
	registerRenderTargetCmds(m, r)
	registerTextureBlitzCmds(m, r)
}

func (m *Module) Shutdown() {}
