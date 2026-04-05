package mbgame

import (
	"fmt"
	"time"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerMathBuiltins(r runtime.Registrar) {
	r.Register("NEWXVALUE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("NEWXVALUE expects 3 arguments")
		}
		x, ok1 := argF(args[0])
		a, ok2 := argF(args[1])
		d, ok3 := argF(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("NEWXVALUE: numeric arguments required")
		}
		return value.FromFloat(newXValue(x, a, d)), nil
	}))
	r.Register("NEWYVALUE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("NEWYVALUE expects 3 arguments")
		}
		y, ok1 := argF(args[0])
		a, ok2 := argF(args[1])
		d, ok3 := argF(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("NEWYVALUE: numeric arguments required")
		}
		return value.FromFloat(newYValue(y, a, d)), nil
	}))
	r.Register("NEWZVALUE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("NEWZVALUE expects 4 arguments")
		}
		z, ok1 := argF(args[0])
		ax, ok2 := argF(args[1])
		ay, ok3 := argF(args[2])
		d, ok4 := argF(args[3])
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return value.Nil, fmt.Errorf("NEWZVALUE: numeric arguments required")
		}
		return value.FromFloat(newZValue(z, ax, ay, d)), nil
	}))
	r.Register("POINTDIR2D", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("POINTDIR2D expects 4 arguments")
		}
		fs := make([]float64, 4)
		for i := 0; i < 4; i++ {
			f, ok := argF(args[i])
			if !ok {
				return value.Nil, fmt.Errorf("POINTDIR2D: numeric arguments required")
			}
			fs[i] = f
		}
		return value.FromFloat(pointDir2D(fs[0], fs[1], fs[2], fs[3])), nil
	}))
	r.Register("POINTDIR3D", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 7 || args[6].Kind != value.KindString {
			return value.Nil, fmt.Errorf("POINTDIR3D expects (x1,y1,z1,x2,y2,z2, axis$)")
		}
		axis, err := rt.ArgString(args, 6)
		if err != nil {
			return value.Nil, err
		}
		fs := make([]float64, 6)
		for i := 0; i < 6; i++ {
			f, ok := argF(args[i])
			if !ok {
				return value.Nil, fmt.Errorf("POINTDIR3D: numeric coordinates required")
			}
			fs[i] = f
		}
		return value.FromFloat(pointDir3D(fs[0], fs[1], fs[2], fs[3], fs[4], fs[5], axis)), nil
	})
	r.Register("CURVEVALUE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("CURVEVALUE expects 3 arguments")
		}
		d, ok1 := argF(args[0])
		s, ok2 := argF(args[1])
		sp, ok3 := argF(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("CURVEVALUE: numeric arguments required")
		}
		return value.FromFloat(curveValue(d, s, sp)), nil
	}))
	r.Register("CURVEANGLE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("CURVEANGLE expects 3 arguments")
		}
		d, ok1 := argF(args[0])
		s, ok2 := argF(args[1])
		sp, ok3 := argF(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("CURVEANGLE: numeric arguments required")
		}
		return value.FromFloat(curveAngle(d, s, sp)), nil
	}))
	r.Register("OSCILLATE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("OSCILLATE expects 3 arguments (speed, min, max)")
		}
		spd, ok1 := argF(args[0])
		minV, ok2 := argF(args[1])
		maxV, ok3 := argF(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("OSCILLATE: numeric arguments required")
		}
		elapsed := time.Since(m.t0).Seconds()
		return value.FromFloat(oscillate(elapsed, spd, minV, maxV)), nil
	}))
	r.Register("WRAPVALUE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("WRAPVALUE expects 3 arguments")
		}
		v, ok1 := argF(args[0])
		minV, ok2 := argF(args[1])
		maxV, ok3 := argF(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("WRAPVALUE: numeric arguments required")
		}
		return value.FromFloat(wrapValue(v, minV, maxV)), nil
	}))
	r.Register("APPROACH", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("APPROACH expects 3 arguments")
		}
		c, ok1 := argF(args[0])
		t, ok2 := argF(args[1])
		s, ok3 := argF(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("APPROACH: numeric arguments required")
		}
		return value.FromFloat(approach(c, t, s)), nil
	}))
}
