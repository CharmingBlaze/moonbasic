package mathmod

import (
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerRemap(r runtime.Registrar) {
	regFlat := func(short, long string, fn runtime.BuiltinFn) {
		r.Register(short, "math", fn)
		r.Register(long, "math", fn)
	}

	regFlat("REMAP", "MATH.REMAP", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 5 {
			return value.Nil, errNArgs(5, len(args))
		}
		v, _ := args[0].ToFloat()
		inMin, _ := args[1].ToFloat()
		inMax, _ := args[2].ToFloat()
		outMin, _ := args[3].ToFloat()
		outMax, _ := args[4].ToFloat()
		if inMin == inMax {
			return value.FromFloat(outMin), nil
		}
		t := (v - inMin) / (inMax - inMin)
		return value.FromFloat(outMin + t*(outMax-outMin)), nil
	})

	regFlat("INVERSE_LERP", "MATH.INVERSE_LERP", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 3 {
			return value.Nil, errNArgs(3, len(args))
		}
		a, _ := args[0].ToFloat()
		b, _ := args[1].ToFloat()
		x, _ := args[2].ToFloat()
		if a == b {
			return value.FromFloat(0), nil
		}
		return value.FromFloat((x - a) / (b - a)), nil
	})

	regFlat("SATURATE", "MATH.SATURATE", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 1 {
			return value.Nil, errNArgs(1, len(args))
		}
		x, _ := args[0].ToFloat()
		return value.FromFloat(clampRL55(x, 0, 1)), nil
	})
}
