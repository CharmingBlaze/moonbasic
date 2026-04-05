package mbgame

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerEaseNoiseRandBuiltins(r runtime.Registrar) {
	r.Register("EASEIN", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("EASEIN expects 1 argument")
		}
		t, ok := argF(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("EASEIN: numeric required")
		}
		return value.FromFloat(easeInQuad(t)), nil
	}))
	r.Register("EASEOUT", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("EASEOUT expects 1 argument")
		}
		t, ok := argF(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("EASEOUT: numeric required")
		}
		return value.FromFloat(easeOutQuad(t)), nil
	}))
	r.Register("EASEINOUT", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("EASEINOUT expects 1 argument")
		}
		t, ok := argF(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("EASEINOUT: numeric required")
		}
		return value.FromFloat(easeInOutQuad(t)), nil
	}))
	r.Register("EASEIN3", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("EASEIN3 expects 1 argument")
		}
		t, ok := argF(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("EASEIN3: numeric required")
		}
		return value.FromFloat(easeInCubic(t)), nil
	}))
	r.Register("EASEOUT3", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("EASEOUT3 expects 1 argument")
		}
		t, ok := argF(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("EASEOUT3: numeric required")
		}
		return value.FromFloat(easeOutCubic(t)), nil
	}))
	r.Register("EASEINOUT3", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("EASEINOUT3 expects 1 argument")
		}
		t, ok := argF(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("EASEINOUT3: numeric required")
		}
		return value.FromFloat(easeInOutCubic(t)), nil
	}))
	r.Register("EASEINSINE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("EASEINSINE expects 1 argument")
		}
		t, ok := argF(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("EASEINSINE: numeric required")
		}
		return value.FromFloat(easeInSine(t)), nil
	}))
	r.Register("EASEOUTSINE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("EASEOUTSINE expects 1 argument")
		}
		t, ok := argF(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("EASEOUTSINE: numeric required")
		}
		return value.FromFloat(easeOutSine(t)), nil
	}))
	r.Register("EASEINOUTSINE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("EASEINOUTSINE expects 1 argument")
		}
		t, ok := argF(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("EASEINOUTSINE: numeric required")
		}
		return value.FromFloat(easeInOutSine(t)), nil
	}))
	r.Register("EASEINBACK", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("EASEINBACK expects 1 argument")
		}
		t, ok := argF(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("EASEINBACK: numeric required")
		}
		return value.FromFloat(easeInBack(t)), nil
	}))
	r.Register("EASEOUTBACK", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("EASEOUTBACK expects 1 argument")
		}
		t, ok := argF(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("EASEOUTBACK: numeric required")
		}
		return value.FromFloat(easeOutBack(t)), nil
	}))
	r.Register("EASEINBOUNCE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("EASEINBOUNCE expects 1 argument")
		}
		t, ok := argF(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("EASEINBOUNCE: numeric required")
		}
		return value.FromFloat(easeInBounce(t)), nil
	}))
	r.Register("EASEOUTBOUNCE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("EASEOUTBOUNCE expects 1 argument")
		}
		t, ok := argF(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("EASEOUTBOUNCE: numeric required")
		}
		return value.FromFloat(easeOutBounce(t)), nil
	}))
	r.Register("EASEINELASTIC", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("EASEINELASTIC expects 1 argument")
		}
		t, ok := argF(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("EASEINELASTIC: numeric required")
		}
		return value.FromFloat(easeInElastic(t)), nil
	}))
	r.Register("EASEOUTELASTIC", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("EASEOUTELASTIC expects 1 argument")
		}
		t, ok := argF(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("EASEOUTELASTIC: numeric required")
		}
		return value.FromFloat(easeOutElastic(t)), nil
	}))
	r.Register("EASELERP", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 4 || args[3].Kind != value.KindString {
			return value.Nil, fmt.Errorf("EASELERP expects (a, b, t, easing$)")
		}
		a, ok1 := argF(args[0])
		b, ok2 := argF(args[1])
		t, ok3 := argF(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("EASELERP: numeric a,b,t required")
		}
		name, err := rt.ArgString(args, 3)
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(easeLerp(a, b, t, name)), nil
	})

	r.Register("PERLIN", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		switch len(args) {
		case 2:
			x, ok1 := argF(args[0])
			y, ok2 := argF(args[1])
			if !ok1 || !ok2 {
				return value.Nil, fmt.Errorf("PERLIN: numeric required")
			}
			return value.FromFloat(perlin2(x, y)), nil
		case 3:
			// either (x,y,z) or (x,y, octaves) — disambiguate: if arg2 is integer-like small, treat as octaves for FBM on x,y
			x, ok1 := argF(args[0])
			y, ok2 := argF(args[1])
			z, ok3 := argF(args[2])
			if !ok1 || !ok2 || !ok3 {
				return value.Nil, fmt.Errorf("PERLIN: numeric required")
			}
			return value.FromFloat(perlin3(x, y, z)), nil
		case 5:
			x, ok1 := argF(args[0])
			y, ok2 := argF(args[1])
			oct, ok3 := argI(args[2])
			freq, ok4 := argF(args[3])
			amp, ok5 := argF(args[4])
			if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
				return value.Nil, fmt.Errorf("PERLIN: invalid arguments")
			}
			return value.FromFloat(fbm2(x*freq, y*freq, int(oct))*amp), nil
		default:
			return value.Nil, fmt.Errorf("PERLIN: 2, 3, or 5 arguments")
		}
	}))
	r.Register("SIMPLEX", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) == 2 {
			x, ok1 := argF(args[0])
			y, ok2 := argF(args[1])
			if !ok1 || !ok2 {
				return value.Nil, fmt.Errorf("SIMPLEX: numeric required")
			}
			return value.FromFloat(simplex2(x, y)), nil
		}
		if len(args) == 3 {
			x, ok1 := argF(args[0])
			y, ok2 := argF(args[1])
			z, ok3 := argF(args[2])
			if !ok1 || !ok2 || !ok3 {
				return value.Nil, fmt.Errorf("SIMPLEX: numeric required")
			}
			return value.FromFloat(simplex3(x, y, z)), nil
		}
		return value.Nil, fmt.Errorf("SIMPLEX expects 2 or 3 arguments")
	}))
	r.Register("VORONOI", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("VORONOI expects 2 arguments")
		}
		x, ok1 := argF(args[0])
		y, ok2 := argF(args[1])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("VORONOI: numeric required")
		}
		return value.FromFloat(voronoi2(x, y)), nil
	}))
	r.Register("FBMNOISE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("FBMNOISE expects (x, y, octaves)")
		}
		x, ok1 := argF(args[0])
		y, ok2 := argF(args[1])
		o, ok3 := argI(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("FBMNOISE: invalid arguments")
		}
		return value.FromFloat(fbm2(x, y, int(o))), nil
	}))

	r.Register("HASHINT", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) == 2 {
			x, ok1 := argI(args[0])
			y, ok2 := argI(args[1])
			if !ok1 || !ok2 {
				return value.Nil, fmt.Errorf("HASHINT: integer arguments required")
			}
			return value.FromInt(int64(hashInt2(int32(x), int32(y)))), nil
		}
		if len(args) == 3 {
			x, ok1 := argI(args[0])
			y, ok2 := argI(args[1])
			z, ok3 := argI(args[2])
			if !ok1 || !ok2 || !ok3 {
				return value.Nil, fmt.Errorf("HASHINT: integer arguments required")
			}
			return value.FromInt(int64(hashInt3(int32(x), int32(y), int32(z)))), nil
		}
		return value.Nil, fmt.Errorf("HASHINT expects 2 or 3 arguments")
	}))
	r.Register("HASHFLOAT", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("HASHFLOAT expects 2 arguments")
		}
		x, ok1 := argF(args[0])
		y, ok2 := argF(args[1])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("HASHFLOAT: numeric required")
		}
		return value.FromFloat(hashFloat2(x, y)), nil
	}))

	r.Register("RNDRANGE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("RNDRANGE expects 2 arguments (min, max)")
		}
		a, ok1 := argI(args[0])
		b, ok2 := argI(args[1])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("RNDRANGE: integer bounds required")
		}
		return value.FromInt(m.rndRange(a, b)), nil
	}))
	r.Register("WEIGHTEDRND", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if m.h == nil {
			return value.Nil, fmt.Errorf("WEIGHTEDRND: heap not bound")
		}
		if len(args) != 1 || args[0].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("WEIGHTEDRND expects 1 array handle")
		}
		ix, err := m.weightedRndIndex(m.h, heap.Handle(args[0].IVal))
		if err != nil {
			return value.Nil, err
		}
		return value.FromInt(int64(ix)), nil
	}))
	r.Register("SHUFFLE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if m.h == nil {
			return value.Nil, fmt.Errorf("SHUFFLE: heap not bound")
		}
		if len(args) != 1 || args[0].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("SHUFFLE expects 1 array handle")
		}
		return value.Nil, m.shuffleArray(m.h, heap.Handle(args[0].IVal))
	}))
	r.Register("RANDOMELEMENT", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if m.h == nil {
			return value.Nil, fmt.Errorf("RANDOMELEMENT: heap not bound")
		}
		if len(args) != 1 || args[0].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("RANDOMELEMENT expects 1 array handle")
		}
		return m.randomElement(m.h, heap.Handle(args[0].IVal))
	}))
}
