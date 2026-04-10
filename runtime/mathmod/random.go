package mathmod

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerRandom(r runtime.Registrar) {
	rndFn := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		switch len(args) {
		case 0:
			return value.FromFloat(m.rng.Float64()), nil
		case 1:
			nf, _ := args[0].ToFloat()
			n := int64(math.Floor(nf))
			if n < 1 {
				return value.FromInt(0), nil
			}
			return value.FromInt(int64(m.rng.Intn(int(n)))), nil
		case 2:
			// Inclusive integer range [min, max] — classic Blitz-style Rand(low, high).
			lo, ok1 := args[0].ToInt()
			hi, ok2 := args[1].ToInt()
			if !ok1 || !ok2 {
				lf, okf1 := args[0].ToFloat()
				hf, okf2 := args[1].ToFloat()
				if !okf1 || !okf2 {
					return value.Nil, fmt.Errorf("RND(min, max): min and max must be numeric")
				}
				lo = int64(math.Floor(lf))
				hi = int64(math.Floor(hf))
			}
			if hi < lo {
				lo, hi = hi, lo
			}
			span := hi - lo + 1
			if span <= 0 {
				return value.FromInt(lo), nil
			}
			return value.FromInt(lo + int64(m.rng.Intn(int(span)))), nil
		default:
			return value.Nil, fmt.Errorf("[moonBASIC] Runtime Error: RND expects 0, 1, or 2 arguments, got %d", len(args))
		}
	}
	r.Register("RND", "math", rndFn)
	r.Register("MATH.RND", "math", rndFn)

	randAlias := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("RAND expects (min, max) inclusive integers")
		}
		return rndFn(rt, args[0], args[1])
	}
	r.Register("RAND", "math", randAlias)
	r.Register("MATH.RAND", "math", randAlias)

	rndfFn := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, errNArgs(2, len(args))
		}
		lo, _ := args[0].ToFloat()
		hi, _ := args[1].ToFloat()
		if hi < lo {
			lo, hi = hi, lo
		}
		return value.FromFloat(lo + (hi-lo)*m.rng.Float64()), nil
	}
	r.Register("RNDF", "math", rndfFn)
	r.Register("MATH.RNDF", "math", rndfFn)

	seedFn := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, errNArgs(1, len(args))
		}
		s, _ := args[0].ToFloat()
		m.reseed(int64(s))
		return value.Nil, nil
	}
	r.Register("RNDSEED", "math", seedFn)
	r.Register("MATH.RNDSEED", "math", seedFn)

	randomizeFn := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) > 1 {
			return value.Nil, fmt.Errorf("[moonBASIC] Runtime Error: RANDOMIZE expects 0 or 1 arguments, got %d", len(args))
		}
		if len(args) == 0 {
			m.reseed(time.Now().UnixNano())
			return value.Nil, nil
		}
		s, _ := args[0].ToFloat()
		m.reseed(int64(s))
		return value.Nil, nil
	}
	r.Register("RANDOMIZE", "math", randomizeFn)
}

func (m *Module) reseed(seed int64) {
	m.rng = rand.New(rand.NewSource(seed))
}
