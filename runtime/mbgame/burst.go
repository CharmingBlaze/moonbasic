package mbgame

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerBurst(r runtime.Registrar) {
	_ = m
	r.Register("GAME.BURSTSPAWN", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		_ = args
		return value.Nil, fmt.Errorf("GAME.BURSTSPAWN is not implemented yet (particle burst bridge)")
	})
}
