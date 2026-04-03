package runtime

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerStringsCheck(r Registrar) {
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

	r.Register("TYPEOF", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("TYPEOF expects 1 argument")
		}
		return rt.RetString(typeOfString(args[0])), nil
	})
	r.Register("ISNULL", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("ISNULL expects 1 argument")
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
	r.Register("ISHANDLE", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("ISHANDLE expects 1 argument")
		}
		return value.FromBool(args[0].Kind == value.KindHandle && args[0].IVal != 0), nil
	})
	r.Register("ISTYPE", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Value{}, Errorf("ISTYPE expects 2 arguments (value, typename$)")
		}
		want, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		want = strings.TrimSpace(strings.ToUpper(want))
		v := args[0]
		if v.Kind == value.KindHandle && v.IVal != 0 {
			if rt.Heap == nil {
				return value.FromBool(false), nil
			}
			obj, ok := rt.Heap.Get(heap.Handle(v.IVal))
			if !ok || obj == nil {
				return value.FromBool(false), nil
			}
			return value.FromBool(strings.EqualFold(obj.TypeName(), want)), nil
		}
		return value.FromBool(typeOfString(v) == want), nil
	})

	r.Register("DUMP", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, Errorf("DUMP expects 1 argument")
		}
		line := dumpLine(rt, args[0])
		if rt.DiagOut != nil {
			fmt.Fprintln(rt.DiagOut, line)
		} else {
			fmt.Println(line)
		}
		return value.Value{}, nil
	})
}

func typeOfString(v value.Value) string {
	switch v.Kind {
	case value.KindNil:
		return "NIL"
	case value.KindInt:
		return "INT"
	case value.KindFloat:
		return "FLOAT"
	case value.KindString:
		return "STRING"
	case value.KindBool:
		return "BOOL"
	case value.KindHandle:
		return "HANDLE"
	default:
		return "UNKNOWN"
	}
}

func dumpLine(rt *Runtime, v value.Value) string {
	switch v.Kind {
	case value.KindNil:
		return "[NIL]"
	case value.KindInt:
		return fmt.Sprintf("[INT %d]", v.IVal)
	case value.KindFloat:
		return fmt.Sprintf("[FLOAT %s]", strconv.FormatFloat(v.FVal, 'g', -1, 64))
	case value.KindBool:
		if v.IVal != 0 {
			return "[BOOL TRUE]"
		}
		return "[BOOL FALSE]"
	case value.KindString:
		s, _ := rt.ArgString([]value.Value{v}, 0)
		return fmt.Sprintf("[STRING %q]", s)
	case value.KindHandle:
		if v.IVal == 0 || rt.Heap == nil {
			return "[HANDLE 0]"
		}
		obj, ok := rt.Heap.Get(heap.Handle(v.IVal))
		if !ok || obj == nil {
			return fmt.Sprintf("[HANDLE ? %d]", v.IVal)
		}
		switch o := obj.(type) {
		case *heap.Array:
			return dumpArrayLine(rt, o)
		case *heap.StringList:
			return fmt.Sprintf("[STRINGLIST len=%d %v]", len(o.Items), o.Items)
		default:
			return fmt.Sprintf("[HANDLE:%s %d]", obj.TypeName(), v.IVal)
		}
	default:
		return fmt.Sprintf("[%s]", v.String())
	}
}

func dumpArrayLine(rt *Runtime, a *heap.Array) string {
	n := a.TotalElements()
	max := n
	if max > 32 {
		max = 32
	}
	var parts []string
	pool := rt.Prog.StringTable
	switch a.Kind {
	case heap.ArrayKindFloat:
		for i := 0; i < max; i++ {
			parts = append(parts, strconv.FormatFloat(a.Floats[i], 'g', -1, 64))
		}
	case heap.ArrayKindBool:
		for i := 0; i < max; i++ {
			if a.Floats[i] != 0 {
				parts = append(parts, "1")
			} else {
				parts = append(parts, "0")
			}
		}
	case heap.ArrayKindString:
		for i := 0; i < max; i++ {
			idx := a.Strings[i]
			s := value.StringAt(value.FromStringIndex(idx), pool)
			if s == "" && rt.Heap != nil {
				if hs, ok := rt.Heap.GetString(idx); ok {
					s = hs
				}
			}
			parts = append(parts, strconv.Quote(s))
		}
	}
	suffix := ""
	if n > max {
		suffix = ",..."
	}
	return fmt.Sprintf("[ARRAY len=%d [%s%s]]", n, strings.Join(parts, ","), suffix)
}
