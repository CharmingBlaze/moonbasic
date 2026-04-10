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

	r.Register("HELP", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("HELP expects 1 argument (commandName$)")
		}
		cmd, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		sig, ok := helpMap[strings.ToUpper(cmd)]
		if !ok {
			fmt.Fprintf(rt.DiagOut, "HELP: No documentation found for '%s'\n", cmd)
			return value.Nil, nil
		}
		fmt.Fprintf(rt.DiagOut, "\x1b[1;36m>> %s \x1b[1;33m%s\x1b[0m\n", strings.ToUpper(cmd), sig)
		return value.Nil, nil
	})

	r.Register("COLORPRINT", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) < 4 {
			return value.Nil, fmt.Errorf("COLORPRINT expects (r, g, b, text...): too few arguments")
		}
		r, _ := rt.ArgInt(args, 0)
		g, _ := rt.ArgInt(args, 1)
		b, _ := rt.ArgInt(args, 2)
		out := rt.formatPrintParts(args[3:])
		// ANSI TRUECOLOR
		fmt.Fprintf(rt.DiagOut, "\x1b[38;2;%d;%d;%dm%s\x1b[0m\n", r, g, b, out)
		return value.Nil, nil
	})
}

var helpMap = map[string]string{
	"POSENT":         "(entity, x#, y#, z#)  - Position an entity in 3D space.",
	"ROTENT":         "(entity, p#, y#, r#)  - Rotate an entity (Euler angles).",
	"SCALENT":        "(entity, x#, y#, z#)  - Scale an entity.",
	"ENTHIT":         "(entity, type#)      - Returns handle of entity hit after last ENTITY.UPDATE.",
	"HITCOUNT":       "(entity)              - Returns number of active collisions.",
	"HITENT":         "(entity, index#)      - Get hit entity handle by index.",
	"UPDW":           "(dt#)                 - Alias intent: use ENTITY.UPDATE(dt) (Blitz UpdateWorld replaced).",
	"ENTRAD":         "(entity, radius#)     - Set sphere collision radius.",
	"ENTTYPE":        "(entity, type#)       - Set collision group (1-32).",
	"CREATECAMERA":   "()                    - Create a new 3D camera entity.",
	"LOADMESH":       "(path$)               - Load a static 3D model.",
	"LOADSPRITE":     "(path$)               - Create a 3D billboard sprite.",
	"ENTITY.UPDATE":  "(dt#)                 - Step entities, collisions, particles (replaces Blitz UpdateWorld).",
	"POSITIONENTITY": "(entity, x#, y#, z#)  - Full name for POSENT.",
	"ROTATEENTITY":   "(entity, p#, y#, r#)  - Full name for ROTENT.",
	"ENTITYTYPE":     "(entity, type#)       - Full name for ENTTYPE.",
	"SKYCOLOR":       "(r, g, b)             - Set clear color.",
	"CREATESPHERE":   "(radius#, segs#)      - Create a sphere entity.",
	"CREATECUBE":     "(w#, h#, d#)          - Create a box entity.",
	"MOUSEHIT":       "(button#)             - 1 if button was pressed this frame.",
	"MOUSEX":         "()                    - Current screen mouse X.",
	"MOUSEY":         "()                    - Current screen mouse Y.",
	"STR$":           "(value)               - Convert value to string.",
	"Window.Open":    "(w, h, title$)       - Initialize display.",
}
