//go:build !cgo && !windows

package mbgame

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerWorldCamera(r runtime.Registrar) {
	_ = m
	r.Register("GAME.ISCURSORONSCREEN", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		_ = args
		return value.Nil, fmt.Errorf("GAME.ISCURSORONSCREEN requires CGO_ENABLED=1")
	})
}
