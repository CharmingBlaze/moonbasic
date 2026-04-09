package mathmod

import (
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const degToRad = math.Pi / 180.0

func sinDeg(x float64) float64 { return math.Sin(x * degToRad) }
func cosDeg(x float64) float64 { return math.Cos(x * degToRad) }
func tanDeg(x float64) float64 { return math.Tan(x * degToRad) }

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
	// Primary SIN/COS/TAN: degrees (Blitz / MoonBasic 3D spec).
	regFlat("SIN", "MATH.SIN", oneFloat(sinDeg))
	regFlat("COS", "MATH.COS", oneFloat(cosDeg))
	regFlat("TAN", "MATH.TAN", oneFloat(tanDeg))
	regFlat("ATN", "MATH.ATN", oneFloat(math.Atan)) // legacy name; radians (Go math)

	// Radian trig for physics / legacy scripts.
	regFlat("SINRAD", "MATH.SINRAD", oneFloat(math.Sin))
	regFlat("COSRAD", "MATH.COSRAD", oneFloat(math.Cos))
	regFlat("TANRAD", "MATH.TANRAD", oneFloat(math.Tan))

	regFlat("ATAN2", "MATH.ATAN2", twoFloat(math.Atan2))
}
