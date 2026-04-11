//go:build cgo || (windows && !cgo)

package input

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerCursor(r runtime.Registrar) {
	r.Register("INPUT.LOCKMOUSE", "input", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("INPUT.LOCKMOUSE expects (lock#)")
		}
		var lock bool
		switch args[0].Kind {
		case value.KindBool:
			lock = args[0].IVal != 0
		case value.KindInt:
			lock = args[0].IVal != 0
		default:
			f, ok := args[0].ToFloat()
			if !ok {
				return value.Nil, fmt.Errorf("INPUT.LOCKMOUSE: lock must be TRUE/FALSE or numeric")
			}
			lock = f != 0
		}
		if lock {
			rl.DisableCursor()
		} else {
			rl.EnableCursor()
		}
		return value.Nil, nil
	}))
	r.Register("CURSOR.SHOW", "input", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("CURSOR.SHOW expects no arguments")
		}
		rl.ShowCursor()
		return value.Nil, nil
	}))
	r.Register("CURSOR.HIDE", "input", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("CURSOR.HIDE expects no arguments")
		}
		rl.HideCursor()
		return value.Nil, nil
	}))
	r.Register("CURSOR.ISHIDDEN", "input", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("CURSOR.ISHIDDEN expects no arguments")
		}
		return value.FromBool(rl.IsCursorHidden()), nil
	}))
	r.Register("CURSOR.ISONSCREEN", "input", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("CURSOR.ISONSCREEN expects no arguments")
		}
		return value.FromBool(rl.IsCursorOnScreen()), nil
	}))
	r.Register("CURSOR.ENABLE", "input", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("CURSOR.ENABLE expects no arguments")
		}
		rl.EnableCursor()
		return value.Nil, nil
	}))
	r.Register("CURSOR.DISABLE", "input", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("CURSOR.DISABLE expects no arguments")
		}
		rl.DisableCursor()
		return value.Nil, nil
	}))
	r.Register("CURSOR.SET", "input", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("CURSOR.SET expects 1 numeric cursor id (Raylib MouseCursor)")
		}
		c, ok := cursorIntArg(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("CURSOR.SET: cursor id must be numeric")
		}
		rl.SetMouseCursor(c)
		return value.Nil, nil
	}))
}

func cursorIntArg(v value.Value) (int32, bool) {
	if i, ok := v.ToInt(); ok {
		return int32(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return int32(f), true
	}
	return 0, false
}
