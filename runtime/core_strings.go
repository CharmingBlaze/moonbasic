package runtime

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerStringBuiltins(r Registrar) {
	r.Register("STR$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("STR$ expects 1 argument")
		}
		if args[0].Kind == value.KindString {
			return args[0], nil
		}
		return rt.RetString(args[0].String()), nil
	})
	r.Register("INT", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("INT expects 1 argument")
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
	r.Register("FLOAT", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("FLOAT expects 1 argument")
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
	r.Register("LEN", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("LEN expects 1 argument")
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
	r.Register("LEFT$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, Errorf("LEFT$ expects 2 arguments")
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
	r.Register("RIGHT$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, Errorf("RIGHT$ expects 2 arguments")
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
	r.Register("MID$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return value.Value{}, Errorf("MID$ expects 2 or 3 arguments")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		start, err := rt.ArgInt(args, 1) // 1-based
		if err != nil {
			return value.Value{}, err
		}
		var ln int64 = int64(len([]rune(s)))
		if len(args) == 3 {
			ln, err = rt.ArgInt(args, 2)
			if err != nil {
				return value.Value{}, err
			}
		}
		runes := []rune(s)
		if start < 1 {
			start = 1
		}
		i := start - 1
		if i >= int64(len(runes)) {
			return rt.RetString(""), nil
		}
		end := i + ln
		if end > int64(len(runes)) {
			end = int64(len(runes))
		}
		return rt.RetString(string(runes[i:end])), nil
	})
	r.Register("UPPER$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("UPPER$ expects 1 argument")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetString(strings.ToUpper(s)), nil
	})
	r.Register("LOWER$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("LOWER$ expects 1 argument")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetString(strings.ToLower(s)), nil
	})
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
		start := 0
		if len(args) == 3 {
			st, err := rt.ArgInt(args, 2)
			if err != nil {
				return value.Value{}, err
			}
			if st > 1 {
				start = int(st) - 1
			}
		}
		if start < 0 || start > len(hay) {
			return value.FromInt(0), nil
		}
		i := strings.Index(hay[start:], needle)
		if i < 0 {
			return value.FromInt(0), nil
		}
		return value.FromInt(int64(i + start + 1)), nil
	})
	r.Register("CHR$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("CHR$ expects 1 argument")
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
	r.Register("ASC", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("ASC expects 1 argument")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		if s == "" {
			return value.FromInt(0), nil
		}
		r, _ := utf8.DecodeRuneInString(s)
		return value.FromInt(int64(r)), nil
	})
	r.Register("BIN$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("BIN$ expects 1 argument")
		}
		n, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetString(strconv.FormatInt(n, 2)), nil
	})
	r.Register("HEX$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("HEX$ expects 1 argument")
		}
		n, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetString(strings.ToUpper(strconv.FormatInt(n, 16))), nil
	})
	r.Register("OCT$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("OCT$ expects 1 argument")
		}
		n, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetString(strconv.FormatInt(n, 8)), nil
	})
	r.Register("VAL", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 || args[0].Kind != value.KindString {
			return value.Value{}, Errorf("VAL expects 1 string argument")
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
	r.Register("BOOL", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("BOOL expects 1 argument")
		}
		return rt.RetBool(value.Truthy(args[0], rt.Prog.StringTable)), nil
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
	r.Register("SPACE$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("SPACE$ expects 1 argument")
		}
		n, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		if n < 0 {
			n = 0
		}
		if n > 1<<20 {
			return value.Value{}, Errorf("SPACE$: n too large")
		}
		return rt.RetString(strings.Repeat(" ", int(n))), nil
	})
	r.Register("STRING$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, Errorf("STRING$ expects 2 arguments (n, c$)")
		}
		n, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		if n < 0 {
			n = 0
		}
		if n > 1<<20 {
			return value.Value{}, Errorf("STRING$: n too large")
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
	r.Register("LSET$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, Errorf("LSET$ expects 2 arguments (s$, width)")
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
	r.Register("RSET$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, Errorf("RSET$ expects 2 arguments (s$, width)")
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
	r.Register("SPLIT$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if rt.Heap == nil {
			return value.Value{}, Errorf("SPLIT$: runtime heap not available")
		}
		if len(args) != 2 {
			return value.Value{}, Errorf("SPLIT$ expects 2 arguments (s$, delim$)")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		d, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		var parts []string
		if d == "" {
			for _, ch := range s {
				parts = append(parts, string(ch))
			}
		} else {
			parts = strings.Split(s, d)
		}
		id, err := rt.Heap.Alloc(&heap.StringList{Items: parts})
		if err != nil {
			return value.Value{}, err
		}
		return value.FromHandle(id), nil
	})
	r.Register("JOIN$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if rt.Heap == nil {
			return value.Value{}, Errorf("JOIN$: runtime heap not available")
		}
		if len(args) != 2 {
			return value.Value{}, Errorf("JOIN$ expects 2 arguments (list, delim$)")
		}
		if args[0].Kind != value.KindHandle {
			return value.Value{}, Errorf("JOIN$: first argument must be a handle from SPLIT$")
		}
		sl, err := heap.Cast[*heap.StringList](rt.Heap, heap.Handle(args[0].IVal))
		if err != nil {
			return value.Value{}, Errorf("JOIN$: %w", err)
		}
		delim, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetString(strings.Join(sl.Items, delim)), nil
	})
	r.Register("REVERSE$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("REVERSE$ expects 1 argument")
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
	r.Register("REPEAT$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, Errorf("REPEAT$ expects 2 arguments (s$, n)")
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
			return value.Value{}, Errorf("REPEAT$: n too large")
		}
		return rt.RetString(strings.Repeat(s, int(n))), nil
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
	r.Register("ISALPHA", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("ISALPHA expects 1 argument")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		if s == "" {
			return value.FromBool(false), nil
		}
		for _, r0 := range s {
			if !unicode.IsLetter(r0) {
				return value.FromBool(false), nil
			}
		}
		return value.FromBool(true), nil
	})
	r.Register("ISNUMERIC", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("ISNUMERIC expects 1 argument")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		s = strings.TrimSpace(s)
		if s == "" {
			return value.FromBool(false), nil
		}
		_, err = strconv.ParseFloat(s, 64)
		return value.FromBool(err == nil), nil
	})
	r.Register("ISALPHANUM", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("ISALPHANUM expects 1 argument")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		if s == "" {
			return value.FromBool(false), nil
		}
		for _, r0 := range s {
			if !unicode.IsLetter(r0) && !unicode.IsDigit(r0) {
				return value.FromBool(false), nil
			}
		}
		return value.FromBool(true), nil
	})
	r.Register("FORMAT$", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 || args[1].Kind != value.KindString {
			return value.Value{}, Errorf("FORMAT$ expects (v, pattern$)")
		}
		pat, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		v := args[0]
		switch v.Kind {
		case value.KindInt:
			return rt.RetString(fmt.Sprintf(pat, v.IVal)), nil
		case value.KindFloat:
			return rt.RetString(fmt.Sprintf(pat, v.FVal)), nil
		case value.KindString:
			s, err := rt.ArgString(args, 0)
			if err != nil {
				return value.Value{}, err
			}
			return rt.RetString(fmt.Sprintf(pat, s)), nil
		case value.KindBool:
			b := v.IVal != 0
			return rt.RetString(fmt.Sprintf(pat, b)), nil
		case value.KindHandle:
			return rt.RetString(fmt.Sprintf(pat, v.IVal)), nil
		default:
			return rt.RetString(fmt.Sprintf(pat, v.String())), nil
		}
	})
}
