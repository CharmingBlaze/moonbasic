package strmod

import (
	"strings"
	"unicode/utf8"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// blitzRuneStart converts a 1-based Blitz/MoonBasic string index to a 0-based rune offset.
func blitzRuneStart(pos int64) int {
	if pos < 1 {
		return 0
	}
	return int(pos - 1)
}

func registerStringsSlice(r runtime.Registrar) {
	leftFn := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, runtime.Errorf("LEFT expects 2 arguments")
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
	}
	r.Register("LEFT", "core", leftFn)
	r.Register("LEFT$", "core", leftFn)

	rightFn := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, runtime.Errorf("RIGHT expects 2 arguments")
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
	}
	r.Register("RIGHT", "core", rightFn)
	r.Register("RIGHT$", "core", rightFn)

	midFn := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return value.Value{}, runtime.Errorf("MID expects 2 or 3 arguments")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		// start is 1-based character index (UTF-8 runes); floats coerce via ArgInt (truncates toward zero).
		startBlitz, err := rt.ArgInt(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		runes := []rune(s)
		i := int64(blitzRuneStart(startBlitz))
		ln := int64(len(runes)) - i
		if len(args) == 3 {
			ln, err = rt.ArgInt(args, 2)
			if err != nil {
				return value.Value{}, err
			}
		}
		if i < 0 {
			i = 0
		}
		if i >= int64(len(runes)) {
			return rt.RetString(""), nil
		}
		end := i + ln
		if end > int64(len(runes)) {
			end = int64(len(runes))
		}
		if ln < 0 {
			return rt.RetString(""), nil
		}
		return rt.RetString(string(runes[i:end])), nil
	}
	r.Register("MID", "core", midFn)
	r.Register("MID$", "core", midFn)
	r.Register("SPACE", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("SPACE expects 1 argument")
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
	r.Register("STRING", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, runtime.Errorf("STRING expects 2 arguments (n, c) or (char, n)")
		}
		// MoonBasic / Blitz order: (char$, n). Legacy: (n, char$) when first arg is numeric.
		var n int64
		var chStr string
		var err error
		if args[0].Kind == value.KindString {
			chStr, err = rt.ArgString(args, 0)
			if err != nil {
				return value.Value{}, err
			}
			n, err = rt.ArgInt(args, 1)
			if err != nil {
				return value.Value{}, err
			}
		} else {
			n, err = rt.ArgInt(args, 0)
			if err != nil {
				return value.Value{}, err
			}
			chStr, err = rt.ArgString(args, 1)
			if err != nil {
				return value.Value{}, err
			}
		}
		if n < 0 {
			n = 0
		}
		if n > 1<<20 {
			return value.Value{}, runtime.Errorf("STRING$: n too large")
		}
		ch := " "
		if chStr != "" {
			r0, _ := utf8.DecodeRuneInString(chStr)
			if r0 != utf8.RuneError {
				ch = string(r0)
			}
		}
		return rt.RetString(strings.Repeat(ch, int(n))), nil
	})
	lsetFn := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, runtime.Errorf("LSET expects 2 arguments (s, width)")
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
	}
	r.Register("LSET", "core", lsetFn)
	r.Register("LSET$", "core", lsetFn)

	rsetFn := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, runtime.Errorf("RSET expects 2 arguments (s, width)")
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
	}
	r.Register("RSET", "core", rsetFn)
	r.Register("RSET$", "core", rsetFn)

	r.Register("REVERSE", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("REVERSE expects 1 argument")
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
	r.Register("REPEAT", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, runtime.Errorf("REPEAT expects 2 arguments (s, n)")
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
