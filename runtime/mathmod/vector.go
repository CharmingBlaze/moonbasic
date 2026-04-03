package mathmod

import (
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerVector(r runtime.Registrar) {
	at2 := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, errNArgs(2, len(args))
		}
		y, _ := args[0].ToFloat()
		x, _ := args[1].ToFloat()
		return value.FromFloat(math.Atan2(y, x)), nil
	}
	r.Register("ATAN2", "math", at2)
	r.Register("MATH.ATAN2", "math", at2)
}
