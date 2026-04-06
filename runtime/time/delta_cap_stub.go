//go:build !cgo && !windows

package mbtime

import (
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func applyDeltaCap(dt float64) float64 { return dt }

func DeltaSeconds(rt *runtime.Runtime) float64 {
	_ = rt
	return 0
}

func registerDeltaCapCommands(reg runtime.Registrar) {
	reg.Register("TIME.SETMAXDELTA", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 1 {
			return value.Nil, errArgs(1, len(args))
		}
		return value.Nil, nil
	})
}
