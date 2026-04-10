// Package mbtime implements TIME.* (clock + optional Raylib frame timing).
package mbtime

import (
	"time"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// Module tracks monotonic clock origin.
type Module struct {
	start     time.Time
	timeScale float64
	paused    bool
	lerpActive  bool
	targetScale float64
	lerpDur     float64
	lerpElapsed float64
}

var (
	GlobalTimeScale = 1.0
	GlobalPaused    = false
)

// NewModule creates TIME builtins.
func NewModule() *Module {
	return &Module{
		start:     time.Now(),
		timeScale: 1.0,
		paused:    false,
	}
}

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	reg.Register("TIME.GET", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		sec := time.Since(m.start).Seconds()
		return value.FromFloat(sec), nil
	})
	// Flat manifest names (same monotonic origin as TIME.GET).
	reg.Register("TIMER", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return value.FromFloat(time.Since(m.start).Seconds()), nil
	})
	reg.Register("TICKCOUNT", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return value.FromInt(time.Since(m.start).Milliseconds()), nil
	})

	delayFn := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 1 {
			return value.Nil, errArgs(1, len(args))
		}
		var ms float64
		if i, ok := args[0].ToInt(); ok {
			ms = float64(i)
		} else if f, ok := args[0].ToFloat(); ok {
			ms = f
		} else {
			return value.Nil, runtime.Errorf("DELAY: milliseconds must be numeric")
		}
		if ms > 0 {
			time.Sleep(time.Duration(ms * float64(time.Millisecond)))
		}
		return value.Nil, nil
	}
	reg.Register("DELAY", "time", delayFn)
	reg.Register("Delay", "time", delayFn)

	reg.Register("TIME.UPDATE", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 { return value.Nil, errArgs(1, len(args)) }
		dt, _ := args[0].ToFloat()
		if m.lerpActive {
			m.lerpElapsed += dt
			if m.lerpElapsed >= m.lerpDur {
				m.timeScale = m.targetScale
				m.lerpActive = false
			} else {
				u := m.lerpElapsed / m.lerpDur
				// Linear lerp back to 1.0 or target
				m.timeScale = m.timeScale + (m.targetScale-m.timeScale)*u
			}
			GlobalTimeScale = m.timeScale
			if rt != nil { rt.TimeScale = GlobalTimeScale }
		}
		return value.Nil, nil
	})

	registerWallClock(reg)
	registerDeltaCapCommands(reg)
	registerRaylibTiming(reg)
	registerMilliSecs(m, reg)

	reg.Register("GAME.SETPAUSE", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, errArgs(1, len(args))
		}
		b, _ := rt.ArgBool(args, 0)
		m.paused = b
		GlobalPaused = m.paused
		if rt != nil { rt.GamePaused = b }
		return value.Nil, nil
	})

	reg.Register("GAME.SLOWMOTION", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, errArgs(2, len(args))
		}
		factor, ok1 := args[0].ToFloat()
		dur, ok2 := args[1].ToFloat()
		if ok1 && ok2 {
			m.timeScale = factor
			m.targetScale = 1.0
			m.lerpDur = dur
			m.lerpElapsed = 0
			m.lerpActive = true
			GlobalTimeScale = factor
			if rt != nil { rt.TimeScale = GlobalTimeScale }
		}
		return value.Nil, nil
	})
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

func errArgs(want, got int) error {
	return runtime.Errorf("expects %d argument(s), got %d", want, got)
}

func (m *Module) Reset() {}

