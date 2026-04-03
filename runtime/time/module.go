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
	registerRaylibTiming(reg)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

func errArgs(want, got int) error {
	return runtime.Errorf("expects %d argument(s), got %d", want, got)
}
