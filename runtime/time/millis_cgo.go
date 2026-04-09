//go:build cgo || (windows && !cgo)

package mbtime

import (
	"moonbasic/runtime"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerMilliSecs(m *Module, reg runtime.Registrar) {
	_ = m
	reg.Register("MilliSecs", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return value.FromFloat(float64(rl.GetTime()) * 1000.0), nil
	})
}
