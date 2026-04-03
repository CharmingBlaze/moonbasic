//go:build !cgo

package mbtime

import (
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerRaylibTiming(r runtime.Registrar) {
	r.Register("TIME.DELTA", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return value.FromFloat(0), nil
	})
	r.Register("TIME.GETFPS", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return value.FromFloat(0), nil
	})
}
