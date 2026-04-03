//go:build !cgo

package mbdraw

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const hint = "DRAW/TEXTURE natives require CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild"

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			return value.Nil, fmt.Errorf("%s: %s", name, hint)
		}
	}
	r.Register("DRAW.RECTANGLE", "draw", stub("DRAW.RECTANGLE"))
	r.Register("DRAW.RECTANGLE_ROUNDED", "draw", stub("DRAW.RECTANGLE_ROUNDED"))
	r.Register("DRAW.TEXTURE", "draw", stub("DRAW.TEXTURE"))
	r.Register("DRAW.TEXTURENPATCH", "draw", stub("DRAW.TEXTURENPATCH"))
	r.Register("TEXTURE.LOAD", "draw", stub("TEXTURE.LOAD"))
	r.Register("TEXTURE.FROMIMAGE", "draw", stub("TEXTURE.FROMIMAGE"))
	r.Register("TEXTURE.FREE", "draw", stub("TEXTURE.FREE"))
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
