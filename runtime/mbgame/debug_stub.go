//go:build !cgo

package mbgame

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerDebugDraw(r runtime.Registrar) {
	_ = m
	r.Register("GAME.DEBUGRECT", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		_ = args
		return value.Nil, fmt.Errorf("GAME.DEBUGRECT requires CGO_ENABLED=1")
	})
}
