package strmod

import (
	"strings"

	"moonbasic/runtime"
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

func registerStringsSearch(r runtime.Registrar) {
	r.Register("TRIM", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("TRIM expects 1 argument")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetString(strings.TrimSpace(s)), nil
	})
	r.Register("LTRIM", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("LTRIM expects 1 argument")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetString(strings.TrimLeft(s, " \t\r\n")), nil
	})
	r.Register("RTRIM", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("RTRIM expects 1 argument")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetString(strings.TrimRight(s, " \t\r\n")), nil
	})
	r.Register("REPLACE", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Value{}, runtime.Errorf("REPLACE expects 3 arguments")
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
	instrImpl := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return value.Value{}, runtime.Errorf("INSTR expects 2 or 3 arguments")
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
			if startRune > 0 {
				startRune-- // 1-based to 0-based
			}
		}
		res := runeIndex(hay, needle, startRune)
		if res == -1 {
			return value.FromInt(0), nil
		}
		return value.FromInt(int64(res + 1)), nil
	}
	r.Register("INSTR", "core", instrImpl)
	r.Register("Instr", "core", instrImpl)
	r.Register("COUNT", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, runtime.Errorf("COUNT expects 2 arguments (src, find)")
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
