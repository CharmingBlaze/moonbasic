//go:build !cgo && !windows

package mblight2d

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/runtime/window"
	"moonbasic/vm/value"
)

// RegisterFrameHook is a no-op without CGO.
func RegisterFrameHook(w *window.Module) { _ = w }

const hint = "LIGHT2D / RENDER.SET2DAmbIENT require CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild"

var stubNames = []string{
	"LIGHT2D.CREATE", "LIGHT2D.MAKE", "LIGHT2D.FREE", "LIGHT2D.SETPOS", "LIGHT2D.SETPOSITION", "LIGHT2D.GETPOS", "LIGHT2D.GETCOLOR", "LIGHT2D.SETCOLOR",
	"LIGHT2D.SETRADIUS", "LIGHT2D.SETINTENSITY", "RENDER.SET2DAmbIENT",
}

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	for _, name := range stubNames {
		n := name
		reg.Register(n, "light2d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			return value.Nil, fmt.Errorf("%s: %s", n, hint)
		})
	}
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
