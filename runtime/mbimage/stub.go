//go:build !cgo

package mbimage

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const hint = "IMAGE.* natives require CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild"

var imageStubNames = []string{
	"IMAGE.LOAD", "IMAGE.LOADRAW", "IMAGE.MAKEBLANK", "IMAGE.MAKE", "IMAGE.MAKECOPY", "IMAGE.MAKETEXT",
	"IMAGE.EXPORT", "IMAGE.WIDTH", "IMAGE.HEIGHT", "IMAGE.FREE",
	"IMAGE.CROP", "IMAGE.RESIZE", "IMAGE.RESIZENN", "IMAGE.FLIPH", "IMAGE.FLIPV",
	"IMAGE.ROTATE", "IMAGE.ROTATECW", "IMAGE.ROTATECCW",
	"IMAGE.COLORTINT", "IMAGE.COLORINVERT", "IMAGE.COLORGRAYSCALE", "IMAGE.COLORCONTRAST",
	"IMAGE.COLORBRIGHTNESS", "IMAGE.COLORREPLACE", "IMAGE.CLEARBACKGROUND",
	"IMAGE.DRAWPIXEL", "IMAGE.DRAWRECT", "IMAGE.DRAWLINE", "IMAGE.DRAWCIRCLE", "IMAGE.DRAWTEXT",
	"IMAGE.DRAWIMAGE", "IMAGE.DITHER", "IMAGE.MIPMAPS", "IMAGE.FORMAT", "IMAGE.DRAWRECTLINES",
	"IMAGE.ALPHACROP", "IMAGE.ALPHACLEAR",
	"IMAGE.GETCOLORR", "IMAGE.GETCOLORG", "IMAGE.GETCOLORB", "IMAGE.GETCOLORA",
	"IMAGE.GETBBOXX", "IMAGE.GETBBOXY", "IMAGE.GETBBOXW", "IMAGE.GETBBOXH",
	"CLIPBOARD.GETIMAGE",
}

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	for _, name := range imageStubNames {
		n := name
		reg.Register(n, "image", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			return value.Nil, fmt.Errorf("%s: %s", n, hint)
		})
	}
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
