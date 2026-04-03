package runtime

import (
	"strings"

	"moonbasic/vm/value"
)

// runeIndex returns the 0-based rune index of needle in hay, searching from startRune inclusive, or -1.
func runeIndex(hay, needle string, startRune int) int {
	hr := []rune(hay)
	if startRune < 0 || startRune > len(hr) {
		return -1
	}
	if needle == "" {
		return startRune
	}
	nr := []rune(needle)
	if len(nr) == 0 {
		return startRune
	}
outer:
	for i := startRune; i <= len(hr)-len(nr); i++ {
		for j := 0; j < len(nr); j++ {
			if hr[i+j] != nr[j] {
				continue outer
			}
		}
		return i
	}
	return -1
}

func registerStringsSearch(r Registrar) {
	r.Register("TRIM$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("TRIM$ expects 1 argument")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetString(strings.TrimSpace(s)), nil
	})
	r.Register("LTRIM$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("LTRIM$ expects 1 argument")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetString(strings.TrimLeft(s, " \t\r\n")), nil
	})
	r.Register("RTRIM$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("RTRIM$ expects 1 argument")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetString(strings.TrimRight(s, " \t\r\n")), nil
	})
	r.Register("CONTAINS", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, Errorf("CONTAINS expects 2 arguments")
		}
		s1, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		s2, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		return value.FromBool(strings.Contains(s1, s2)), nil
	})
	r.Register("REPLACE$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Value{}, Errorf("REPLACE$ expects 3 arguments")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		from, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		to, err := rt.ArgString(args, 2)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetString(strings.ReplaceAll(s, from, to)), nil
	})
	r.Register("INSTR", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return value.Value{}, Errorf("INSTR expects 2 or 3 arguments")
		}
		hay, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		needle, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		startRune := 0
		if len(args) == 3 {
			st, err := rt.ArgInt(args, 2)
			if err != nil {
				return value.Value{}, err
			}
			startRune = int(st)
		}
		return value.FromInt(int64(runeIndex(hay, needle, startRune))), nil
	})
	r.Register("STARTSWITH", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, Errorf("STARTSWITH expects 2 arguments")
		}
		s1, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		s2, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		return value.FromBool(strings.HasPrefix(s1, s2)), nil
	})
	r.Register("ENDSWITH", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, Errorf("ENDSWITH expects 2 arguments")
		}
		s1, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		s2, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		return value.FromBool(strings.HasSuffix(s1, s2)), nil
	})
	r.Register("COUNT$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, Errorf("COUNT$ expects 2 arguments (src$, find$)")
		}
		s1, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		s2, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		return value.FromInt(int64(strings.Count(s1, s2))), nil
	})
}
