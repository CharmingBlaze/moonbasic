package mbgame

import (
	"fmt"
	"time"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerPauseFrame(r runtime.Registrar) {
	r.Register("PAUSEGAME", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("PAUSEGAME expects 0 arguments")
		}
		if rt != nil {
			rt.GamePaused = true
		}
		return value.Nil, nil
	})
	r.Register("RESUMEGAME", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("RESUMEGAME expects 0 arguments")
		}
		if rt != nil {
			rt.GamePaused = false
		}
		return value.Nil, nil
	})
	r.Register("GAMEPAUSE", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GAMEPAUSE expects 0 arguments")
		}
		if rt == nil {
			return value.FromBool(false), nil
		}
		return value.FromBool(rt.GamePaused), nil
	})
	r.Register("FRAMECOUNT", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("FRAMECOUNT expects 0 arguments")
		}
		if rt == nil {
			return value.FromInt(0), nil
		}
		return value.FromInt(int64(rt.FrameCount)), nil
	})
	r.Register("ELAPSED", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("ELAPSED expects 0 arguments")
		}
		return value.FromFloat(time.Since(m.t0).Seconds()), nil
	})
}
