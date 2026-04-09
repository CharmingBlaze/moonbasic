package driver

import (
	"errors"
	"fmt"

	"moonbasic/internal/raylibpurego"
)

// ErrRaylibUnavailable is returned when no Raylib backend could be resolved (no linked CGO and no loadable sidecar DLL/SO).
var ErrRaylibUnavailable = errors.New("raylib unavailable")

// CheckWindow returns nil if a window backend can be used, or [ErrRaylibUnavailable] with a hint.
func CheckWindow(sel Selection) error {
	if sel.Kind == KindUnavailable {
		return fmt.Errorf("%w: %s (place %s next to the executable, set %s=purego, or rebuild with CGO and a C toolchain)",
			ErrRaylibUnavailable, sel.Detail, raylibpurego.LibPath(), EnvDriver)
	}
	return nil
}
