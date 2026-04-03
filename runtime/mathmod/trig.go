package mathmod

import (
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerTrig(r runtime.Registrar) {
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
	regFlat("SIN", "MATH.SIN", oneFloat(math.Sin))
	regFlat("COS", "MATH.COS", oneFloat(math.Cos))
	regFlat("TAN", "MATH.TAN", oneFloat(math.Tan))
	regFlat("ATN", "MATH.ATN", oneFloat(math.Atan))
}
