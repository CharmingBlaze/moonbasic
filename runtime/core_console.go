package runtime

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"moonbasic/vm/value"
)

var consoleIn = bufio.NewReader(os.Stdin)

func (rt *Runtime) formatPrintParts(args []value.Value) string {
	if len(args) == 0 {
		return ""
	}
	var sb strings.Builder
	for i, a := range args {
		if i > 0 {
			sb.WriteByte(' ')
		}
		if a.Kind == value.KindString {
			idx := int32(a.IVal)
			s, ok := rt.Heap.GetString(idx)
			if !ok && rt.Prog != nil && idx >= 0 && int(idx) < len(rt.Prog.StringTable) {
				s = rt.Prog.StringTable[idx]
				ok = true
			}
			if ok {
				sb.WriteString(s)
			}
		} else {
			sb.WriteString(a.String())
		}
	}
	return sb.String()
}

func registerConsoleIO(r Registrar) {
	r.Register("PRINT", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		fmt.Fprintln(rt.DiagOut, rt.formatPrintParts(args))
		return value.Value{}, nil
	})
	r.Register("PRINTLN", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		fmt.Fprintln(rt.DiagOut, rt.formatPrintParts(args))
		return value.Value{}, nil
	})
	r.Register("WRITE", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		fmt.Fprint(rt.DiagOut, rt.formatPrintParts(args))
		return value.Value{}, nil
	})
	r.Register("INPUT", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		prompt, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		def := ""
		if len(args) > 1 {
			def, err = rt.ArgString(args, 1)
			if err != nil {
				return value.Value{}, err
			}
		}
		fmt.Fprint(rt.DiagOut, prompt)
		line, err := consoleIn.ReadString('\n')
		if err != nil {
			return value.Value{}, err
		}
		line = strings.TrimRight(line, "\r\n")
		if strings.TrimSpace(line) == "" && def != "" {
			return rt.RetString(def), nil
		}
		return rt.RetString(line), nil
	})
	r.Register("CLS", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		fmt.Fprint(rt.DiagOut, "\x1b[2J\x1b[H")
		return value.Value{}, nil
	})
	r.Register("LOCATE", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		row, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		col, err := rt.ArgInt(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		if row < 1 {
			row = 1
		}
		if col < 1 {
			col = 1
		}
		fmt.Fprintf(rt.DiagOut, "\x1b[%d;%dH", row, col)
		return value.Value{}, nil
	})

	spacesFn := func(name string) BuiltinFn {
		return func(rt *Runtime, args ...value.Value) (value.Value, error) {
			n, err := rt.ArgInt(args, 0)
			if err != nil {
				return value.Value{}, err
			}
			if n < 0 {
				n = 0
			}
			if n > 1<<20 {
				return value.Value{}, Errorf("%s: n too large", name)
			}
			fmt.Fprint(rt.DiagOut, strings.Repeat(" ", int(n)))
			return value.Value{}, nil
		}
	}
	r.Register("TAB", "core", spacesFn("TAB"))
	r.Register("SPC", "core", spacesFn("SPC"))
}
