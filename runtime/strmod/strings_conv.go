package strmod

import (
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerStringsConv(r runtime.Registrar) {
	r.Register("STR$", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("STR$ expects 1 argument")
		}
		if args[0].Kind == value.KindString {
			return args[0], nil
		}
		return rt.RetString(args[0].String()), nil
	})
	r.Register("INT", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("INT expects 1 argument")
		}
		v := args[0]
		if v.Kind == value.KindString {
			s, err := rt.ArgString(args, 0)
			if err != nil {
				return value.Value{}, err
			}
			s = strings.TrimSpace(s)
			if s == "" {
				return value.FromInt(0), nil
			}
			if f, err := strconv.ParseFloat(s, 64); err == nil {
				return value.FromInt(int64(math.Floor(f))), nil
			}
			if i, err := strconv.ParseInt(s, 10, 64); err == nil {
				return value.FromInt(i), nil
			}
			return value.FromInt(0), nil
		}
		if f, ok := v.ToFloat(); ok {
			return value.FromInt(int64(math.Floor(f))), nil
		}
		if i, ok := v.ToInt(); ok {
			return value.FromInt(i), nil
		}
		return value.FromInt(0), nil
	})
	r.Register("FLOAT", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("FLOAT expects 1 argument")
		}
		f, ok := args[0].ToFloat()
		if !ok {
			if args[0].Kind == value.KindString {
				s, err := rt.ArgString(args, 0)
				if err != nil {
					return value.Value{}, err
				}
				v, _ := strconv.ParseFloat(s, 64)
				return value.FromFloat(v), nil
			}
			return value.FromFloat(0.0), nil
		}
		return value.FromFloat(f), nil
	})
	r.Register("LEN", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("LEN expects 1 argument")
		}
		if args[0].Kind != value.KindString {
			return value.FromInt(0), nil
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		return value.FromInt(int64(utf8.RuneCountInString(s))), nil
	})
	r.Register("CHR$", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("CHR$ expects 1 argument")
		}
		c, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		if c < 0 || c > 0x10FFFF {
			return rt.RetString(""), nil
		}
		return rt.RetString(string(rune(c))), nil
	})
	r.Register("ASC", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("ASC expects 1 argument")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		if s == "" {
			return value.FromInt(0), nil
		}
		r0, _ := utf8.DecodeRuneInString(s)
		return value.FromInt(int64(r0)), nil
	})
	r.Register("VAL", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 || args[0].Kind != value.KindString {
			return value.Value{}, runtime.Errorf("VAL expects 1 string argument")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		s = strings.TrimSpace(s)
		if s == "" {
			return value.FromFloat(0), nil
		}
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return value.FromFloat(0), nil
		}
		return value.FromFloat(v), nil
	})
	r.Register("BOOL", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("BOOL expects 1 argument")
		}
		return rt.RetBool(value.Truthy(args[0], rt.Prog.StringTable, rt.Heap)), nil
	})
}
