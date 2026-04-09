//go:build !cgo && !windows

package mbtime

import (
	"time"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerMilliSecs(m *Module, reg runtime.Registrar) {
	reg.Register("MilliSecs", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return value.FromFloat(time.Since(m.start).Seconds() * 1000.0), nil
	})
}
