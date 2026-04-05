package mbgame

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerColorFormatBuiltins(r runtime.Registrar) {
	_ = m
	r.Register("RGB", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("RGB expects 3 arguments")
		}
		a := make([]int64, 3)
		for i := 0; i < 3; i++ {
			x, ok := argI(args[i])
			if !ok {
				return value.Nil, fmt.Errorf("RGB: integer components required")
			}
			a[i] = x
		}
		return value.FromInt(packRGB(a[0], a[1], a[2])), nil
	}))
	r.Register("ARGB", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("ARGB expects 4 arguments")
		}
		a := make([]int64, 4)
		for i := 0; i < 4; i++ {
			x, ok := argI(args[i])
			if !ok {
				return value.Nil, fmt.Errorf("ARGB: integer components required")
			}
			a[i] = x
		}
		return value.FromInt(packARGB(a[0], a[1], a[2], a[3])), nil
	}))
	r.Register("RGBR", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("RGBR expects 1 argument")
		}
		c, ok := argI(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("RGBR: integer color required")
		}
		return value.FromInt(unpackR(c)), nil
	}))
	r.Register("RGBG", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("RGBG expects 1 argument")
		}
		c, ok := argI(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("RGBG: integer color required")
		}
		return value.FromInt(unpackG(c)), nil
	}))
	r.Register("RGBB", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("RGBB expects 1 argument")
		}
		c, ok := argI(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("RGBB: integer color required")
		}
		return value.FromInt(unpackB(c)), nil
	}))
	r.Register("RGBA", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("RGBA expects 1 argument (alpha channel)")
		}
		c, ok := argI(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("RGBA: integer color required")
		}
		return value.FromInt(unpackA(c)), nil
	}))
	r.Register("RGBMIX", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("RGBMIX expects 3 arguments (c1, c2, t)")
		}
		c1, ok1 := argI(args[0])
		c2, ok2 := argI(args[1])
		t, ok3 := argF(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("RGBMIX: invalid arguments")
		}
		return value.FromInt(rgbMix(c1, c2, t)), nil
	}))
	r.Register("RGBBRIGHTEN", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("RGBBRIGHTEN expects 2 arguments")
		}
		c, ok1 := argI(args[0])
		amt, ok2 := argF(args[1])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("RGBBRIGHTEN: invalid arguments")
		}
		return value.FromInt(rgbBrighten(c, amt)), nil
	}))
	r.Register("RGBDARKEN", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("RGBDARKEN expects 2 arguments")
		}
		c, ok1 := argI(args[0])
		amt, ok2 := argF(args[1])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("RGBDARKEN: invalid arguments")
		}
		return value.FromInt(rgbDarken(c, amt)), nil
	}))
	r.Register("RGBFADE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("RGBFADE expects 2 arguments")
		}
		c, ok1 := argI(args[0])
		a, ok2 := argF(args[1])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("RGBFADE: invalid arguments")
		}
		return value.FromInt(rgbFade(c, a)), nil
	}))

	// Named palette — zero-arg functions (no CONST keyword required in BASIC).
	r.Register("COL_WHITE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("COL_WHITE expects 0 arguments")
		}
		return value.FromInt(colWhite()), nil
	}))
	r.Register("COL_BLACK", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("COL_BLACK expects 0 arguments")
		}
		return value.FromInt(colBlack()), nil
	}))
	r.Register("COL_RED", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("COL_RED expects 0 arguments")
		}
		return value.FromInt(colRed()), nil
	}))
	r.Register("COL_GREEN", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("COL_GREEN expects 0 arguments")
		}
		return value.FromInt(colGreen()), nil
	}))
	r.Register("COL_BLUE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("COL_BLUE expects 0 arguments")
		}
		return value.FromInt(colBlue()), nil
	}))
	r.Register("COL_YELLOW", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("COL_YELLOW expects 0 arguments")
		}
		return value.FromInt(colYellow()), nil
	}))
	r.Register("COL_CYAN", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("COL_CYAN expects 0 arguments")
		}
		return value.FromInt(colCyan()), nil
	}))
	r.Register("COL_MAGENTA", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("COL_MAGENTA expects 0 arguments")
		}
		return value.FromInt(colMagenta()), nil
	}))
	r.Register("COL_ORANGE", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("COL_ORANGE expects 0 arguments")
		}
		return value.FromInt(colOrange()), nil
	}))
	r.Register("COL_GRAY", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("COL_GRAY expects 0 arguments")
		}
		return value.FromInt(colGray()), nil
	}))
	r.Register("COL_DARKGRAY", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("COL_DARKGRAY expects 0 arguments")
		}
		return value.FromInt(colDarkGray()), nil
	}))
	r.Register("COL_LIGHTGRAY", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("COL_LIGHTGRAY expects 0 arguments")
		}
		return value.FromInt(colLightGray()), nil
	}))
	r.Register("COL_TRANSPARENT", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("COL_TRANSPARENT expects 0 arguments")
		}
		return value.FromInt(colTransparent()), nil
	}))

	r.Register("FORMATINT", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("FORMATINT expects 2 arguments")
		}
		n, ok1 := argI(args[0])
		d, ok2 := argI(args[1])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("FORMATINT: integer arguments required")
		}
		return rt.RetString(formatInt(int(n), int(d))), nil
	})
	r.Register("FORMATSCORE", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("FORMATSCORE expects 1 argument")
		}
		n, ok := argI(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("FORMATSCORE: integer required")
		}
		return rt.RetString(formatScore(n)), nil
	})
	r.Register("FORMATTIME", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("FORMATTIME expects 1 argument (seconds)")
		}
		n, ok := argI(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("FORMATTIME: integer seconds required")
		}
		return rt.RetString(formatTime(int(n))), nil
	})
	r.Register("FORMATTIME2", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("FORMATTIME2 expects 1 argument (seconds)")
		}
		n, ok := argI(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("FORMATTIME2: integer seconds required")
		}
		return rt.RetString(formatTime2(int(n))), nil
	})
}
