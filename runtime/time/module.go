// Package mbtime implements TIME.* (clock + optional Raylib frame timing).
package mbtime

import (
	"time"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// Module tracks monotonic clock origin.
type Module struct {
	start time.Time
}

// NewModule creates TIME builtins.
func NewModule() *Module {
	return &Module{start: time.Now()}
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

	registerWallClock(reg)
	registerDeltaCapCommands(reg)
	registerRaylibTiming(reg)
	registerMilliSecs(m, reg)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

func errArgs(want, got int) error {
	return runtime.Errorf("expects %d argument(s), got %d", want, got)
}
