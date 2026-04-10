package mbgame

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerTimeScale(r runtime.Registrar) {
	r.Register("GAME.SETTIMESCALE", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GAME.SETTIMESCALE expects scale#")
		}
		var s float64
		if f, ok := args[0].ToFloat(); ok {
			s = f
		} else if i, ok := args[0].ToInt(); ok {
			s = float64(i)
		} else {
			return value.Nil, fmt.Errorf("GAME.SETTIMESCALE: scale must be numeric")
		}
		if rt == nil {
			return value.Nil, fmt.Errorf("GAME.SETTIMESCALE: runtime not available")
		}
		rt.TimeScale = s
		return value.Nil, nil
	})
	r.Register("GAME.GETTIMESCALE", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GAME.GETTIMESCALE expects 0 arguments")
		}
		if rt == nil {
			return value.FromFloat(1), nil
		}
		s := rt.TimeScale
		if s == 0 {
			s = 1
		}
		return value.FromFloat(s), nil
	})
}
