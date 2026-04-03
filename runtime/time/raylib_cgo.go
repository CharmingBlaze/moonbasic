//go:build cgo

package mbtime

import (
	"moonbasic/runtime"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerRaylibTiming(reg runtime.Registrar) {
	reg.Register("TIME.DELTA", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return value.FromFloat(float64(rl.GetFrameTime())), nil
	})
	reg.Register("TIME.GETFPS", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return value.FromFloat(float64(rl.GetFPS())), nil
	})
}
