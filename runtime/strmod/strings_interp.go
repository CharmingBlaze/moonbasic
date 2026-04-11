package strmod

import (
	"fmt"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerStringsInterp(r runtime.Registrar) {
	interp := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) < 2 {
			return value.Value{}, runtime.Errorf("INTERP expects (template, value0, [value1, ...])")
		}
		tmpl, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		n := len(args) - 1
		if n > 10 {
			return value.Value{}, runtime.Errorf("INTERP$: at most 10 placeholders {0}..{9}")
		}
		pairs := make([]string, 0, n*2)
		for i := 0; i < n; i++ {
			pairs = append(pairs, fmt.Sprintf("{%d}", i))
			pairs = append(pairs, valueStringForInterp(rt, args[i+1]))
		}
		out := strings.NewReplacer(pairs...).Replace(tmpl)
		return rt.RetString(out), nil
	}
	r.Register("INTERP", "core", interp)
	r.Register("STRING.INTERP", "core", interp)
}

func valueStringForInterp(rt *runtime.Runtime, v value.Value) string {
	if v.Kind == value.KindString {
		s, err := rt.ArgString([]value.Value{v}, 0)
		if err != nil {
			return ""
		}
		return s
	}
	return v.String()
}
