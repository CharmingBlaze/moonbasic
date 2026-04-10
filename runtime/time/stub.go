//go:build !cgo && !windows

package mbtime

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerRaylibTiming(r runtime.Registrar) {
	r.Register("TIME.DELTA", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		switch len(args) {
		case 0:
			return value.FromFloat(0), nil
		case 2:
			minF, ok1 := args[0].ToFloat()
			maxF, ok2 := args[1].ToFloat()
			if !ok1 || !ok2 {
				return value.Nil, fmt.Errorf("TIME.DELTA(min#, max#): min and max must be numeric")
			}
			if maxF < minF {
				minF, maxF = maxF, minF
			}
			return value.FromFloat(minF), nil
		default:
			return value.Nil, fmt.Errorf("TIME.DELTA expects 0 or 2 arguments, got %d", len(args))
		}
	})
	r.Register("TIME.GETFPS", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return value.FromFloat(0), nil
	})
}
