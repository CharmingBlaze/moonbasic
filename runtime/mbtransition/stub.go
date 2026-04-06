//go:build !cgo && !windows

package mbtransition

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/runtime/window"
	"moonbasic/vm/value"
)

const hint = "TRANSITION.* requires CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild"

// RegisterFrameHook is a no-op without CGO.
func RegisterFrameHook(w *window.Module) { _ = w }

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	for _, n := range []string{
		"TRANSITION.FADEOUT", "TRANSITION.FADEIN", "TRANSITION.ISDONE",
		"TRANSITION.WIPE", "TRANSITION.SETCOLOR",
	} {
		name := n
		reg.Register(name, "transition", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			return value.Nil, fmt.Errorf("%s: %s", name, hint)
		})
	}
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
