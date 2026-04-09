package mathmod

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerBlitzSurface(r runtime.Registrar) {
	oneFloat := func(f func(float64) float64) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			if len(args) != 1 {
				return value.Nil, errNArgs(1, len(args))
			}
			x, _ := args[0].ToFloat()
			return value.FromFloat(f(x)), nil
		}
	}
	twoFloat := func(f func(float64, float64) float64) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			if len(args) != 2 {
				return value.Nil, errNArgs(2, len(args))
			}
			a, _ := args[0].ToFloat()
			b, _ := args[1].ToFloat()
			return value.FromFloat(f(a, b)), nil
		}
	}
	// Friendly names (same behavior as uppercase MATH.* / flat).
	r.Register("Sin", "math", oneFloat(sinDeg))
	r.Register("Cos", "math", oneFloat(cosDeg))
	r.Register("Tan", "math", oneFloat(tanDeg))
	r.Register("Sqrt", "math", oneFloat(math.Sqrt))
	r.Register("Abs", "math", oneFloat(math.Abs))

	// Rnd(min, max) float — same as RNDF.
	r.Register("Rnd", "math", twoFloat(func(a, b float64) float64 {
		if b < a {
			a, b = b, a
		}
		return a + (b-a)*m.rng.Float64()
	}))

	// Rand(min, max) inclusive integer.
	r.Register("Rand", "math", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, errNArgs(2, len(args))
		}
		lo, ok1 := args[0].ToInt()
		if !ok1 {
			if f, ok := args[0].ToFloat(); ok {
				lo = int64(math.Floor(f))
				ok1 = true
			}
		}
		hi, ok2 := args[1].ToInt()
		if !ok2 {
			if f, ok := args[1].ToFloat(); ok {
				hi = int64(math.Floor(f))
				ok2 = true
			}
		}
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("Rand: min and max must be numeric")
		}
		if hi < lo {
			lo, hi = hi, lo
		}
		span := hi - lo + 1
		if span <= 0 {
			return value.FromInt(lo), nil
		}
		return value.FromInt(lo + int64(m.rng.Intn(int(span)))), nil
	})

	r.Register("SeedRnd", "math", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, errNArgs(1, len(args))
		}
		s, _ := args[0].ToFloat()
		m.reseed(int64(s))
		return value.Nil, nil
	})

	// PascalCase aliases (same Float64 behavior as MATH.* / flat uppercase).
	r.Register("ASin", "math", oneFloat(math.Asin))
	r.Register("ACos", "math", oneFloat(math.Acos))
	r.Register("ATan", "math", oneFloat(math.Atan))
	r.Register("ATan2", "math", twoFloat(math.Atan2))
	r.Register("Floor", "math", oneFloat(math.Floor))
	r.Register("Ceil", "math", oneFloat(math.Ceil))
	// Round: use flat ROUND / MATH.ROUND (supports 1- or 2-argument form); do not register "Round" here — it would overwrite ROUND in the registry.
	r.Register("Exp", "math", oneFloat(math.Exp))
	r.Register("Log", "math", oneFloat(math.Log))
	r.Register("Log10", "math", oneFloat(math.Log10))
}
