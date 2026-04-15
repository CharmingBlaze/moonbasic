//go:build !cgo && !windows

package texture

import (
	"fmt"
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func stub(name string) runtime.BuiltinFn {
	return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("%s requires CGO_ENABLED=1", name)
	}
}

func (m *Module) Register(r runtime.Registrar) {
	r.Register("TEXTURE.LOAD", "texture", stub("TEXTURE.LOAD"))
	r.Register("LOADTEXTURE", "texture", stub("LOADTEXTURE"))
	r.Register("LoadTexture", "texture", stub("LoadTexture"))
	r.Register("TEXTURE.FREE", "texture", stub("TEXTURE.FREE"))
	r.Register("FREETEXTURE", "texture", stub("FREETEXTURE"))
	r.Register("TEXTURE.FROMIMAGE", "texture", stub("TEXTURE.FROMIMAGE"))
	r.Register("IMAGE.TOTEXTURE", "texture", stub("IMAGE.TOTEXTURE"))
	for _, k := range []string{
		"TEXTURE.SETGRID", "TEXTURE.SETFRAME", "TEXTURE.LOADANIM", "TEXTURE.PLAY", "TEXTURE.STOPANIM", "TEXTURE.TICKALL",
		"TEXTURE.SETUVSCROLL", "TEXTURE.SETDISTORTION",
		"TEXTURE.WIDTH", "TEXTURE.HEIGHT", "TEXTUREWIDTH", "TEXTUREHEIGHT", "TEXTURE.SETFILTER", "TEXTURE.SETDEFAULTFILTER", "MATERIAL.AUTOFILTER", "TEXTURE.SETWRAP", "TEXTURE.UPDATE", "TEXTURE.RELOAD",
		"LEVEL.PRELOAD", "RENDER.CLEARCACHE",
		"TEXTURE.GENWHITENOISE", "TEXTURE.GENCHECKED", "TEXTURE.GENGRADIENTV", "TEXTURE.GENGRADIENTH", "TEXTURE.GENCOLOR",
		"RENDERTARGET.CREATE", "RENDERTARGET.MAKE", "RENDERTARGET.FREE", "RENDERTARGET.BEGIN", "RENDERTARGET.END", "RENDERTARGET.TEXTURE",
		"CreateTexture", "LoadAnimTexture", "TextureWidth", "TextureHeight", "TextureName",
		"SetCubeFace", "SetCubeMode", "TextureCoords", "ScaleTexture", "RotateTexture", "PositionTexture",
	} {
		r.Register(k, "texture", stub(k))
	}
}

func (m *Module) Shutdown() {}
