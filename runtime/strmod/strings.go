package strmod

import (
	"fmt"
	"moonbasic/runtime"
	"moonbasic/vm/value"
	"strings"
	"unicode"
)

func registerStrings(r runtime.Registrar) {
	r.Register("strmod.LSET$", "strmod", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		width, err := rt.ArgInt(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		if len(s) > int(width) {
			return rt.RetString(s[:width]), nil
		}
		return rt.RetString(s + strings.Repeat(" ", int(width)-len(s))), nil
	})
	r.Register("strmod.RSET$", "strmod", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		width, err := rt.ArgInt(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		if len(s) > int(width) {
			return rt.RetString(s[:width]), nil
		}
		return rt.RetString(strings.Repeat(" ", int(width)-len(s)) + s), nil
	})
	r.Register("strmod.SPACE$", "strmod", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		n, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetString(strings.Repeat(" ", int(n))), nil
	})
	r.Register("strmod.STRING$", "strmod", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		n, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		c, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetString(strings.Repeat(c, int(n))), nil
	})
	r.Register("strmod.REVERSE$", "strmod", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		runes := []rune(s)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return rt.RetString(string(runes)), nil
	})
	r.Register("strmod.REPEAT$", "strmod", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		n, err := rt.ArgInt(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetString(strings.Repeat(s, int(n))), nil
	})
	r.Register("strmod.COUNT$", "strmod", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		src, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		find, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetInt(int64(strings.Count(src, find))), nil
	})
	r.Register("strmod.ISALPHA", "strmod", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		for _, r := range s {
			if !unicode.IsLetter(r) {
				return rt.RetBool(false), nil
			}
		}
		return rt.RetBool(len(s) > 0), nil
	})
	r.Register("strmod.ISNUMERIC", "strmod", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		for _, r := range s {
			if !unicode.IsDigit(r) {
				return rt.RetBool(false), nil
			}
		}
		return rt.RetBool(len(s) > 0), nil
	})
	r.Register("strmod.ISALPHANUM", "strmod", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		for _, r := range s {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				return rt.RetBool(false), nil
			}
		}
		return rt.RetBool(len(s) > 0), nil
	})
	r.Register("strmod.FORMAT$", "strmod", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		format, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		// This is a simplified version. A full implementation would need to handle multiple args.
		val := args[0]
		var goVal interface{}
		switch val.Kind {
		case value.KindInt:
			goVal = val.IVal
		case value.KindFloat:
			goVal = val.FVal
		case value.KindString:
			goVal, _ = rt.Heap.GetString(int32(val.IVal))
		case value.KindBool:
			goVal = val.IVal != 0
		default:
			goVal = ""
		}
		return rt.RetString(fmt.Sprintf(format, goVal)), nil
	})
}
