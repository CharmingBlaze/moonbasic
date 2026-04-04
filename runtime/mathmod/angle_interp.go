package mathmod

import (
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerAngleInterp(r runtime.Registrar) {
	regFlat := func(short, long string, fn runtime.BuiltinFn) {
		r.Register(short, "math", fn)
		r.Register(long, "math", fn)
	}
	oneFloat := func(f func(float64) float64) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			if len(args) != 1 {
				return value.Nil, errNArgs(1, len(args))
			}
			x, _ := args[0].ToFloat()
			return value.FromFloat(f(x)), nil
		}
	}
	threeFloat := func(f func(float64, float64, float64) float64) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			if len(args) != 3 {
				return value.Nil, errNArgs(3, len(args))
			}
			a, _ := args[0].ToFloat()
			b, _ := args[1].ToFloat()
			c, _ := args[2].ToFloat()
			return value.FromFloat(f(a, b, c)), nil
		}
	}

	regFlat("DEG2RAD", "MATH.DEG2RAD", oneFloat(func(d float64) float64 { return d * math.Pi / 180 }))
	regFlat("RAD2DEG", "MATH.RAD2DEG", oneFloat(func(r float64) float64 { return r * 180 / math.Pi }))
	regFlat("WRAPANGLE", "MATH.WRAPANGLE", oneFloat(wrapAngle360))
	regFlat("WRAPANGLE180", "MATH.WRAPANGLE180", oneFloat(wrapAngle180))
	regFlat("ANGLEDIFF", "MATH.ANGLEDIFF", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, errNArgs(2, len(args))
		}
		a, _ := args[0].ToFloat()
		b, _ := args[1].ToFloat()
		return value.FromFloat(angleDiffDeg(a, b)), nil
	})

	regFlat("LERP", "MATH.LERP", threeFloat(lerpRL55))
	regFlat("SMOOTHSTEP", "MATH.SMOOTHSTEP", threeFloat(smoothstep))
	regFlat("CLAMP", "MATH.CLAMP", threeFloat(func(v, lo, hi float64) float64 {
		if lo > hi {
			lo, hi = hi, lo
		}
		return clampRL55(v, lo, hi)
	}))
	regFlat("PINGPONG", "MATH.PINGPONG", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, errNArgs(2, len(args))
		}
		t, _ := args[0].ToFloat()
		L, _ := args[1].ToFloat()
		return value.FromFloat(pingPong(t, L)), nil
	})
	regFlat("WRAP", "MATH.WRAP", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Nil, errNArgs(3, len(args))
		}
		v, _ := args[0].ToFloat()
		min, _ := args[1].ToFloat()
		max, _ := args[2].ToFloat()
		if max < min {
			min, max = max, min
		}
		return value.FromFloat(wrapRL55(v, min, max)), nil
	})
}

func wrapAngle360(a float64) float64 {
	a = math.Mod(a, 360)
	if a < 0 {
		a += 360
	}
	return a
}

func wrapAngle180(a float64) float64 {
	return math.Mod(a+540, 360) - 180
}

func angleDiffDeg(a, b float64) float64 {
	return math.Mod(b-a+540, 360) - 180
}

func smoothstep(lo, hi, x float64) float64 {
	if hi == lo {
		return 0
	}
	t := (x - lo) / (hi - lo)
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	return t * t * (3 - 2*t)
}

func pingPong(t, length float64) float64 {
	if length <= 0 {
		return 0
	}
	t = math.Mod(t, 2*length)
	if t < 0 {
		t += 2 * length
	}
	if t > length {
		return 2*length - t
	}
	return t
}
