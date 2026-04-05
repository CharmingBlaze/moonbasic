package mbgame

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerTileSprite(r runtime.Registrar) {
	_ = m
	r.Register("GAME.SPRITETILEBRIDGE", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		_ = args
		return value.Nil, fmt.Errorf("GAME.SPRITETILEBRIDGE is not implemented yet (use collision math + sprite bounds in script)")
	})
}
