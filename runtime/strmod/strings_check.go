package strmod

import (
	"strings"
	"unicode"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerStringsCheck(r runtime.Registrar) {
	r.Register("ISALPHA", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 || args[0].Kind != value.KindString {
			return value.Value{}, runtime.Errorf("ISALPHA expects 1 string argument")
		}
		s, _ := rt.ArgString(args, 0)
		if s == "" {
			return value.FromBool(false), nil
		}
		for _, r := range s {
			if !unicode.IsLetter(r) {
				return value.FromBool(false), nil
			}
		}
		return value.FromBool(true), nil
	})
	r.Register("ISALPHANUM", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 || args[0].Kind != value.KindString {
			return value.Value{}, runtime.Errorf("ISALPHANUM expects 1 string argument")
		}
		s, _ := rt.ArgString(args, 0)
		if s == "" {
			return value.FromBool(false), nil
		}
		for _, r := range s {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				return value.FromBool(false), nil
			}
		}
		return value.FromBool(true), nil
	})
	r.Register("ISNUMERIC", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 || args[0].Kind != value.KindString {
			return value.Value{}, runtime.Errorf("ISNUMERIC expects 1 string argument")
		}
		s, _ := rt.ArgString(args, 0)
		if s == "" {
			return value.FromBool(false), nil
		}
		for _, r := range s {
			if !unicode.IsDigit(r) && r != '.' && r != '-' {
				return value.FromBool(false), nil
			}
		}
		return value.FromBool(true), nil
	})
	r.Register("CONTAINS", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, runtime.Errorf("CONTAINS expects 2 arguments")
		}
		s, _ := rt.ArgString(args, 0)
		sub, _ := rt.ArgString(args, 1)
		return value.FromBool(strings.Contains(s, sub)), nil
	})
	r.Register("STARTSWITH", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, runtime.Errorf("STARTSWITH expects 2 arguments")
		}
		s, _ := rt.ArgString(args, 0)
		sub, _ := rt.ArgString(args, 1)
		return value.FromBool(strings.HasPrefix(s, sub)), nil
	})
	r.Register("ENDSWITH", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, runtime.Errorf("ENDSWITH expects 2 arguments")
		}
		s, _ := rt.ArgString(args, 0)
		sub, _ := rt.ArgString(args, 1)
		return value.FromBool(strings.HasSuffix(s, sub)), nil
	})

	r.Register("TYPEOF", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("TYPEOF expects 1 argument")
		}
		return rt.RetString(typeOfString(args[0])), nil
	})
	r.Register("ISNULL", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("ISNULL expects 1 argument")
		}
		v := args[0]
		if v.Kind == value.KindNil {
			return value.FromBool(true), nil
		}
		if v.Kind == value.KindHandle && v.IVal == 0 {
			return value.FromBool(true), nil
		}
		return value.FromBool(false), nil
	})
	r.Register("ISHANDLE", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("ISHANDLE expects 1 argument")
		}
		return value.FromBool(args[0].Kind == value.KindHandle && args[0].IVal != 0), nil
	})
	r.Register("ISTYPE", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, runtime.Errorf("ISTYPE expects 2 arguments (value, typename$)")
		}
		want, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		want = strings.TrimSpace(strings.ToUpper(want))
		got := strings.ToUpper(typeOfString(args[0]))
		return value.FromBool(got == want), nil
	})
}

func typeOfString(v value.Value) string {
	switch v.Kind {
	case value.KindNil:
		return "nil"
	case value.KindInt:
		return "int"
	case value.KindFloat:
		return "float"
	case value.KindString:
		return "string"
	case value.KindBool:
		return "bool"
	case value.KindHandle:
		return "handle"
	default:
		return "unknown"
	}
}
