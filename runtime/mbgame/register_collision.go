package mbgame

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerCollisionBuiltins(r runtime.Registrar) {
	_ = m
	r.Register("BOXCOLLIDE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 8 {
			return value.Nil, fmt.Errorf("BOXCOLLIDE expects 8 arguments")
		}
		x := make([]float64, 8)
		for i := 0; i < 8; i++ {
			f, ok := argF(args[i])
			if !ok {
				return value.Nil, fmt.Errorf("BOXCOLLIDE: numeric arguments required")
			}
			x[i] = f
		}
		return value.FromBool(boxCollide2D(x[0], x[1], x[2], x[3], x[4], x[5], x[6], x[7])), nil
	}))
	r.Register("CIRCLECOLLIDE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("CIRCLECOLLIDE expects 6 arguments")
		}
		fs := make([]float64, 6)
		for i := 0; i < 6; i++ {
			f, ok := argF(args[i])
			if !ok {
				return value.Nil, fmt.Errorf("CIRCLECOLLIDE: numeric arguments required")
			}
			fs[i] = f
		}
		return value.FromBool(circleCollide2D(fs[0], fs[1], fs[2], fs[3], fs[4], fs[5])), nil
	}))
	r.Register("POINTINBOX", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("POINTINBOX expects 6 arguments")
		}
		fs := make([]float64, 6)
		for i := 0; i < 6; i++ {
			f, ok := argF(args[i])
			if !ok {
				return value.Nil, fmt.Errorf("POINTINBOX: numeric arguments required")
			}
			fs[i] = f
		}
		return value.FromBool(pointInBox2D(fs[0], fs[1], fs[2], fs[3], fs[4], fs[5])), nil
	}))
	r.Register("POINTINCIRCLE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("POINTINCIRCLE expects 5 arguments")
		}
		fs := make([]float64, 5)
		for i := 0; i < 5; i++ {
			f, ok := argF(args[i])
			if !ok {
				return value.Nil, fmt.Errorf("POINTINCIRCLE: numeric arguments required")
			}
			fs[i] = f
		}
		return value.FromBool(pointInCircle2D(fs[0], fs[1], fs[2], fs[3], fs[4])), nil
	}))
	r.Register("CIRCLEBOXCOLLIDE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 7 {
			return value.Nil, fmt.Errorf("CIRCLEBOXCOLLIDE expects 7 arguments")
		}
		fs := make([]float64, 7)
		for i := 0; i < 7; i++ {
			f, ok := argF(args[i])
			if !ok {
				return value.Nil, fmt.Errorf("CIRCLEBOXCOLLIDE: numeric arguments required")
			}
			fs[i] = f
		}
		return value.FromBool(circleBoxCollide2D(fs[0], fs[1], fs[2], fs[3], fs[4], fs[5], fs[6])), nil
	}))
	r.Register("LINECOLLIDE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 8 {
			return value.Nil, fmt.Errorf("LINECOLLIDE expects 8 arguments")
		}
		fs := make([]float64, 8)
		for i := 0; i < 8; i++ {
			f, ok := argF(args[i])
			if !ok {
				return value.Nil, fmt.Errorf("LINECOLLIDE: numeric arguments required")
			}
			fs[i] = f
		}
		return value.FromBool(lineCollide2D(fs[0], fs[1], fs[2], fs[3], fs[4], fs[5], fs[6], fs[7])), nil
	}))
	r.Register("POINTONLINE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 7 {
			return value.Nil, fmt.Errorf("POINTONLINE expects 7 arguments")
		}
		fs := make([]float64, 7)
		for i := 0; i < 7; i++ {
			f, ok := argF(args[i])
			if !ok {
				return value.Nil, fmt.Errorf("POINTONLINE: numeric arguments required")
			}
			fs[i] = f
		}
		return value.FromBool(pointOnLine2D(fs[0], fs[1], fs[2], fs[3], fs[4], fs[5], fs[6])), nil
	}))
	r.Register("SPHERECOLLIDE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 8 {
			return value.Nil, fmt.Errorf("SPHERECOLLIDE expects 8 arguments (x1,y1,z1,r1,x2,y2,z2,r2)")
		}
		fs := make([]float64, 8)
		for i := 0; i < 8; i++ {
			f, ok := argF(args[i])
			if !ok {
				return value.Nil, fmt.Errorf("SPHERECOLLIDE: numeric arguments required")
			}
			fs[i] = f
		}
		return value.FromBool(sphereCollide3D(fs[0], fs[1], fs[2], fs[3], fs[4], fs[5], fs[6], fs[7])), nil
	}))
	r.Register("AABBCOLLIDE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 12 {
			return value.Nil, fmt.Errorf("AABBCOLLIDE expects 12 arguments")
		}
		fs := make([]float64, 12)
		for i := 0; i < 12; i++ {
			f, ok := argF(args[i])
			if !ok {
				return value.Nil, fmt.Errorf("AABBCOLLIDE: numeric arguments required")
			}
			fs[i] = f
		}
		return value.FromBool(aabbCollide3D(fs[0], fs[1], fs[2], fs[3], fs[4], fs[5], fs[6], fs[7], fs[8], fs[9], fs[10], fs[11])), nil
	}))
	r.Register("SPHEREBOXCOLLIDE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 10 {
			return value.Nil, fmt.Errorf("SPHEREBOXCOLLIDE expects 10 arguments")
		}
		fs := make([]float64, 10)
		for i := 0; i < 10; i++ {
			f, ok := argF(args[i])
			if !ok {
				return value.Nil, fmt.Errorf("SPHEREBOXCOLLIDE: numeric arguments required")
			}
			fs[i] = f
		}
		return value.FromBool(sphereBoxCollide3D(fs[0], fs[1], fs[2], fs[3], fs[4], fs[5], fs[6], fs[7], fs[8], fs[9])), nil
	}))
	r.Register("POINTINAABB", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 9 {
			return value.Nil, fmt.Errorf("POINTINAABB expects 9 arguments")
		}
		fs := make([]float64, 9)
		for i := 0; i < 9; i++ {
			f, ok := argF(args[i])
			if !ok {
				return value.Nil, fmt.Errorf("POINTINAABB: numeric arguments required")
			}
			fs[i] = f
		}
		return value.FromBool(pointInAABB3D(fs[0], fs[1], fs[2], fs[3], fs[4], fs[5], fs[6], fs[7], fs[8])), nil
	}))
	r.Register("DISTANCE2D", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("DISTANCE2D expects 4 arguments")
		}
		fs := make([]float64, 4)
		for i := 0; i < 4; i++ {
			f, ok := argF(args[i])
			if !ok {
				return value.Nil, fmt.Errorf("DISTANCE2D: numeric arguments required")
			}
			fs[i] = f
		}
		return value.FromFloat(distance2D(fs[0], fs[1], fs[2], fs[3])), nil
	}))
	r.Register("DISTANCE3D", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("DISTANCE3D expects 6 arguments")
		}
		fs := make([]float64, 6)
		for i := 0; i < 6; i++ {
			f, ok := argF(args[i])
			if !ok {
				return value.Nil, fmt.Errorf("DISTANCE3D: numeric arguments required")
			}
			fs[i] = f
		}
		return value.FromFloat(distance3D(fs[0], fs[1], fs[2], fs[3], fs[4], fs[5])), nil
	}))
	r.Register("DISTANCESQ2D", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("DISTANCESQ2D expects 4 arguments")
		}
		fs := make([]float64, 4)
		for i := 0; i < 4; i++ {
			f, ok := argF(args[i])
			if !ok {
				return value.Nil, fmt.Errorf("DISTANCESQ2D: numeric arguments required")
			}
			fs[i] = f
		}
		return value.FromFloat(distanceSq2D(fs[0], fs[1], fs[2], fs[3])), nil
	}))
	r.Register("DISTANCESQ3D", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("DISTANCESQ3D expects 6 arguments")
		}
		fs := make([]float64, 6)
		for i := 0; i < 6; i++ {
			f, ok := argF(args[i])
			if !ok {
				return value.Nil, fmt.Errorf("DISTANCESQ3D: numeric arguments required")
			}
			fs[i] = f
		}
		return value.FromFloat(distanceSq3D(fs[0], fs[1], fs[2], fs[3], fs[4], fs[5])), nil
	}))
}
