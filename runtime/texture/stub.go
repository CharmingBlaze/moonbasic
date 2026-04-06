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
	r.Register("TEXTURE.FREE", "texture", stub("TEXTURE.FREE"))
	r.Register("TEXTURE.FROMIMAGE", "texture", stub("TEXTURE.FROMIMAGE"))
	for _, k := range []string{
		"TEXTURE.WIDTH", "TEXTURE.HEIGHT", "TEXTURE.SETFILTER", "TEXTURE.SETWRAP", "TEXTURE.UPDATE",
		"TEXTURE.GENWHITENOISE", "TEXTURE.GENCHECKED", "TEXTURE.GENGRADIENTV", "TEXTURE.GENGRADIENTH", "TEXTURE.GENCOLOR",
		"RENDERTARGET.MAKE", "RENDERTARGET.FREE", "RENDERTARGET.BEGIN", "RENDERTARGET.END", "RENDERTARGET.TEXTURE",
	} {
		r.Register(k, "texture", stub(k))
	}
}

func (m *Module) Shutdown() {}
