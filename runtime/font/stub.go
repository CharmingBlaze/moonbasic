//go:build !cgo

package mbfont

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const hint = "FONT.* requires CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild"

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			return value.Nil, fmt.Errorf("%s: %s", name, hint)
		}
	}
	r.Register("FONT.LOAD", "font", stub("FONT.LOAD"))
	r.Register("FONT.LOADBDF", "font", stub("FONT.LOADBDF"))
	r.Register("FONT.FREE", "font", stub("FONT.FREE"))
	r.Register("FONT.DRAWDEFAULT", "font", stub("FONT.DRAWDEFAULT"))
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
