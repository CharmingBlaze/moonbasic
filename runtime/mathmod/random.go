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
		if len(args) > 1 {
			return value.Nil, fmt.Errorf("[moonBASIC] Runtime Error: RND expects at most 1 argument, got %d", len(args))
		}
		if len(args) == 0 {
			return value.FromFloat(m.rng.Float64()), nil
		}
		nf, _ := args[0].ToFloat()
		n := int64(math.Floor(nf))
		if n < 1 {
			return value.FromInt(0), nil
		}
		return value.FromInt(int64(m.rng.Intn(int(n)))), nil
	}
	r.Register("RND", "math", rndFn)
	r.Register("MATH.RND", "math", rndFn)

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
