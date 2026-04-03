package bitwise

import (
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerBitwise(r runtime.Registrar) {
	r.Register("core.BAND", "bitwise", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		a, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		b, err := rt.ArgInt(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetInt(a & b), nil
	})
	r.Register("core.BOR", "bitwise", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		a, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		b, err := rt.ArgInt(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetInt(a | b), nil
	})
	r.Register("core.BXOR", "bitwise", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		a, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		b, err := rt.ArgInt(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetInt(a ^ b), nil
	})
	r.Register("core.BNOT", "bitwise", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		a, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetInt(^a), nil
	})
	r.Register("core.BLSHIFT", "bitwise", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		a, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		n, err := rt.ArgInt(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetInt(a << n), nil
	})
	r.Register("core.BRSHIFT", "bitwise", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		a, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		n, err := rt.ArgInt(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetInt(a >> n), nil
	})
	r.Register("core.BTEST", "bitwise", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		a, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		bit, err := rt.ArgInt(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetBool((a & (1 << bit)) != 0), nil
	})
	r.Register("core.BSET", "bitwise", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		a, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		bit, err := rt.ArgInt(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetInt(a | (1 << bit)), nil
	})
	r.Register("core.BCLEAR", "bitwise", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		a, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		bit, err := rt.ArgInt(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetInt(a &^ (1 << bit)), nil
	})
	r.Register("core.BTOGGLE", "bitwise", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		a, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		bit, err := rt.ArgInt(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetInt(a ^ (1 << bit)), nil
	})
	r.Register("core.BCOUNT", "bitwise", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		a, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		count := 0
		for a > 0 {
			a &= (a - 1)
			count++
		}
		return rt.RetInt(int64(count)), nil
	})
}
