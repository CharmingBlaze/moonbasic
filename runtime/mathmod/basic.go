package mathmod

import (
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerBasic(r runtime.Registrar) {
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
	zeroFloat := func(c float64) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			if len(args) != 0 {
				return value.Nil, errNArgs(0, len(args))
			}
			return value.FromFloat(c), nil
		}
	}

	regFlat("ABS", "MATH.ABS", oneFloat(math.Abs))
	regFlat("SQRT", "MATH.SQRT", oneFloat(math.Sqrt))
	regFlat("SQR", "MATH.SQR", oneFloat(math.Sqrt))
	regFlat("EXP", "MATH.EXP", oneFloat(math.Exp))
	regFlat("LOG", "MATH.LOG", oneFloat(math.Log))
	regFlat("LOG2", "MATH.LOG2", oneFloat(math.Log2))
	regFlat("LOG10", "MATH.LOG10", oneFloat(math.Log10))
	regFlat("ASIN", "MATH.ASIN", oneFloat(math.Asin))
	regFlat("ACOS", "MATH.ACOS", oneFloat(math.Acos))
	regFlat("ATAN", "MATH.ATAN", oneFloat(math.Atan))
	regFlat("POW", "MATH.POW", twoFloat(math.Pow))
	regFlat("FLOOR", "MATH.FLOOR", oneFloat(math.Floor))
	regFlat("CEIL", "MATH.CEIL", oneFloat(math.Ceil))

	regFlat("PI", "MATH.PI", zeroFloat(math.Pi))
	regFlat("TAU", "MATH.TAU", zeroFloat(2*math.Pi))
	regFlat("E", "MATH.E", zeroFloat(math.E))

	regFlat("SGN", "MATH.SGN", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, errNArgs(1, len(args))
		}
		x, _ := args[0].ToFloat()
		switch {
		case x > 0:
			return value.FromInt(1), nil
		case x < 0:
			return value.FromInt(-1), nil
		default:
			return value.FromInt(0), nil
		}
	})
	regFlat("FIX", "MATH.FIX", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, errNArgs(1, len(args))
		}
		x, _ := args[0].ToFloat()
		return value.FromFloat(math.Trunc(x)), nil
	})
	regFlat("MIN", "MATH.MIN", twoFloat(math.Min))
	regFlat("MAX", "MATH.MAX", twoFloat(math.Max))

	regFlat("ROUND", "MATH.ROUND", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) < 1 || len(args) > 2 {
			return value.Nil, errNArgsRange("1 or 2", len(args))
		}
		x, _ := args[0].ToFloat()
		if len(args) == 1 {
			return value.FromFloat(math.Round(x)), nil
		}
		dec, _ := args[1].ToFloat()
		p := math.Pow(10, dec)
		if math.IsInf(p, 0) || p == 0 {
			return value.FromFloat(x), nil
		}
		return value.FromFloat(math.Round(x*p) / p), nil
	})
}
