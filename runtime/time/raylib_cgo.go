//go:build cgo || (windows && !cgo)

package mbtime

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerRaylibTiming(reg runtime.Registrar) {
	reg.Register("TIME.DELTA", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		switch len(args) {
		case 0:
			return value.FromFloat(DeltaSeconds(rt)), nil
		case 2:
			minF, ok1 := args[0].ToFloat()
			maxF, ok2 := args[1].ToFloat()
			if !ok1 || !ok2 {
				return value.Nil, fmt.Errorf("TIME.DELTA(min#, max#): min and max must be numeric")
			}
			if maxF < minF {
				minF, maxF = maxF, minF
			}
			if rt != nil && rt.GamePaused {
				return value.FromFloat(0), nil
			}
			dt := float64(rl.GetFrameTime())
			if dt <= 0 {
				dt = minF
			} else if dt < minF {
				dt = minF
			}
			if dt > maxF {
				dt = maxF
			}
			if rt != nil {
				s := rt.TimeScale
				if s != 0 && s != 1 {
					dt *= s
				}
			}
			return value.FromFloat(dt), nil
		default:
			return value.Nil, fmt.Errorf("TIME.DELTA expects 0 or 2 arguments, got %d", len(args))
		}
	})
	reg.Register("TIME.GETFPS", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return value.FromFloat(float64(rl.GetFPS())), nil
	})
}
