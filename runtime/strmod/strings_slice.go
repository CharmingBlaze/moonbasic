package strmod

import (
	"strings"
	"unicode/utf8"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerStringsSlice(r runtime.Registrar) {
	r.Register("LEFT$", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, runtime.Errorf("LEFT$ expects 2 arguments")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		n, err := rt.ArgInt(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		if n <= 0 {
			return rt.RetString(""), nil
		}
		runes := []rune(s)
		if int64(len(runes)) < n {
			n = int64(len(runes))
		}
		return rt.RetString(string(runes[:n])), nil
	})
	r.Register("RIGHT$", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, runtime.Errorf("RIGHT$ expects 2 arguments")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		n, err := rt.ArgInt(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		runes := []rune(s)
		if n <= 0 {
			return rt.RetString(""), nil
		}
		if int64(len(runes)) < n {
			return rt.RetString(s), nil
		}
		return rt.RetString(string(runes[len(runes)-int(n):])), nil
	})
	r.Register("MID$", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return value.Value{}, runtime.Errorf("MID$ expects 2 or 3 arguments")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		pos, err := rt.ArgInt(args, 1) // 0-based rune index
		if err != nil {
			return value.Value{}, err
		}
		runes := []rune(s)
		ln := int64(len(runes))
		if len(args) == 3 {
			ln, err = rt.ArgInt(args, 2)
			if err != nil {
				return value.Value{}, err
			}
		}
		if pos < 0 {
			pos = 0
		}
		if pos >= int64(len(runes)) {
			return rt.RetString(""), nil
		}
		i := pos
		end := i + ln
		if end > int64(len(runes)) {
			end = int64(len(runes))
		}
		if ln < 0 {
			return rt.RetString(""), nil
		}
		return rt.RetString(string(runes[i:end])), nil
	})
	r.Register("SPACE$", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("SPACE$ expects 1 argument")
		}
		n, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		if n < 0 {
			n = 0
		}
		if n > 1<<20 {
			return value.Value{}, runtime.Errorf("SPACE$: n too large")
		}
		return rt.RetString(strings.Repeat(" ", int(n))), nil
	})
	r.Register("STRING$", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, runtime.Errorf("STRING$ expects 2 arguments (n, c$)")
		}
		n, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		if n < 0 {
			n = 0
		}
		if n > 1<<20 {
			return value.Value{}, runtime.Errorf("STRING$: n too large")
		}
		ch := " "
		if args[1].Kind == value.KindString {
			s, err := rt.ArgString(args, 1)
			if err != nil {
				return value.Value{}, err
			}
			if s != "" {
				r0, _ := utf8.DecodeRuneInString(s)
				if r0 != utf8.RuneError {
					ch = string(r0)
				}
			}
		}
		return rt.RetString(strings.Repeat(ch, int(n))), nil
	})
	r.Register("LSET$", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, runtime.Errorf("LSET$ expects 2 arguments (s$, width)")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		w, err := rt.ArgInt(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		if w < 0 {
			w = 0
		}
		runes := []rune(s)
		if int64(len(runes)) > w {
			return rt.RetString(string(runes[:w])), nil
		}
		pad := int(w) - len(runes)
		if pad <= 0 {
			return rt.RetString(s), nil
		}
		return rt.RetString(s + strings.Repeat(" ", pad)), nil
	})
	r.Register("RSET$", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, runtime.Errorf("RSET$ expects 2 arguments (s$, width)")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		w, err := rt.ArgInt(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		if w < 0 {
			w = 0
		}
		runes := []rune(s)
		if int64(len(runes)) > w {
			return rt.RetString(string(runes[int64(len(runes))-w:])), nil
		}
		pad := int(w) - len(runes)
		if pad <= 0 {
			return rt.RetString(s), nil
		}
		return rt.RetString(strings.Repeat(" ", pad) + s), nil
	})
	r.Register("REVERSE$", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("REVERSE$ expects 1 argument")
		}
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
	r.Register("REPEAT$", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, runtime.Errorf("REPEAT$ expects 2 arguments (s$, n)")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		n, err := rt.ArgInt(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		if n <= 0 {
			return rt.RetString(""), nil
		}
		if n > 1<<20 {
			return value.Value{}, runtime.Errorf("REPEAT$: n too large")
		}
		return rt.RetString(strings.Repeat(s, int(n))), nil
	})
}
