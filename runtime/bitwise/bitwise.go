package bitwise

import (
	"math/bits"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerBitwise(r runtime.Registrar) {
	// Manifest uses flat BAND, BOR, …; legacy entries are core.B* — same implementation.
	reg := func(coreKey, flatKey string, fn runtime.BuiltinFn) {
		r.Register(coreKey, "bitwise", fn)
		r.Register(flatKey, "bitwise", fn)
	}

	reg("core.BAND", "BAND", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
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
	reg("core.BOR", "BOR", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
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
	reg("core.BXOR", "BXOR", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
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
	reg("core.BNOT", "BNOT", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		a, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetInt(^a), nil
	})
	reg("core.BLSHIFT", "BLSHIFT", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
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
	reg("core.BRSHIFT", "BRSHIFT", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
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
	reg("core.BTEST", "BTEST", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
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
	reg("core.BSET", "BSET", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
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
	reg("core.BCLEAR", "BCLEAR", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
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
	reg("core.BTOGGLE", "BTOGGLE", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
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
	reg("core.BCOUNT", "BCOUNT", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		a, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		// Population count of the 64-bit two's-complement pattern (negative values count high bits).
		return rt.RetInt(int64(bits.OnesCount64(uint64(a)))), nil
	})
}
