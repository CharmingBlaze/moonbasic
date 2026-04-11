//go:build !cgo && windows

package mbgui

import (
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func puregoResetDefaultTheme() {
	pg.styleInt = make(map[int64]int64)
	pg.styleColor = make(map[int64]rl.Color)
	pg.widgetInt = make(map[puregoRectKey]int32)
	// DEFAULT control colors (approximate raygui dark)
	pg.styleColor[styleKey(0, 0)] = rl.Color{R: 72, G: 78, B: 88, A: 255}
	pg.styleColor[styleKey(0, 1)] = rl.Color{R: 48, G: 54, B: 64, A: 255}
	pg.styleColor[styleKey(0, 2)] = rl.Color{R: 228, G: 232, B: 242, A: 255}
	pg.lastTheme = "DEFAULT"
}

func registerPuregoTheme(m *Module, reg runtime.Registrar) {
	reg.Register("GUI.SETCOLOR", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 7 {
			return value.Nil, fmt.Errorf("GUI.SETCOLOR expects (control, property, r, g, b, a)")
		}
		c, ok1 := argI32(args[0])
		p, ok2 := argI32(args[1])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("GUI.SETCOLOR: control and property must be numeric")
		}
		col, err := colorArgs(args, 2)
		if err != nil {
			return value.Nil, err
		}
		pg.styleColor[styleKey(c, p)] = col
		return value.Nil, nil
	})
	reg.Register("GUI.SETTEXTSIZE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.SETTEXTSIZE expects (size)")
		}
		sz, ok := argI32(args[0])
		if !ok || sz < 1 {
			return value.Nil, fmt.Errorf("GUI.SETTEXTSIZE: size must be a positive integer")
		}
		pg.textSize = float32(sz)
		return value.Nil, nil
	})
	reg.Register("GUI.SETTEXTSPACING", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.SETTEXTSPACING expects (spacing)")
		}
		s, ok := argF32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.SETTEXTSPACING: spacing must be numeric")
		}
		pg.textSpacing = s
		return value.Nil, nil
	})
	reg.Register("GUI.SETTEXTLINEHEIGHT", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.SETTEXTLINEHEIGHT expects (extraSpacing)")
		}
		s, ok := argF32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("value must be numeric")
		}
		pg.lineExtra = s
		return value.Nil, nil
	})
	reg.Register("GUI.SETTEXTWRAP", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.SETTEXTWRAP expects (mode)")
		}
		mo, ok := argI32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("mode must be numeric")
		}
		pg.wrapMode = mo
		return value.Nil, nil
	})
	reg.Register("GUI.SETTEXTALIGN", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.SETTEXTALIGN expects (mode)")
		}
		mo, ok := argI32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("mode must be numeric")
		}
		pg.alignH = mo
		return value.Nil, nil
	})
	reg.Register("GUI.SETTEXTALIGNVERT", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.SETTEXTALIGNVERT expects (mode)")
		}
		mo, ok := argI32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("mode must be numeric")
		}
		pg.alignV = mo
		return value.Nil, nil
	})
	reg.Register("GUI.GETTEXTSIZE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.GETTEXTSIZE expects 0 arguments")
		}
		return rt.RetInt(int64(pg.textSize)), nil
	})
	reg.Register("GUI.THEMEAPPLY", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.THEMEAPPLY expects (name)")
		}
		name, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		switch strings.ToUpper(strings.TrimSpace(name)) {
		case "DEFAULT", "RESET", "LIGHT", "BUILTIN_DARK":
			puregoResetDefaultTheme()
			return value.Nil, nil
		default:
			return value.Nil, fmt.Errorf("GUI.THEMEAPPLY: unknown theme %q in purego GUI (use DEFAULT)", name)
		}
	})
	reg.Register("GUI.THEMENAMES", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.THEMENAMES expects 0 arguments")
		}
		return rt.RetString("DEFAULT;LIGHT;BUILTIN_DARK"), nil
	})
	reg.Register("GUI.LOADSTYLE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.LOADSTYLE expects (path)")
		}
		_, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		return value.Nil, fmt.Errorf("GUI.LOADSTYLE: .rgs styles require CGO raygui")
	})
	reg.Register("GUI.LOADDEFAULTSTYLE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.LOADDEFAULTSTYLE expects 0 arguments")
		}
		puregoResetDefaultTheme()
		return value.Nil, nil
	})
	reg.Register("GUI.LOADSTYLEMEM", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("GUI.LOADSTYLEMEM: requires CGO raygui")
	})
	reg.Register("GUI.LOADICONS", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("GUI.LOADICONS: requires CGO raygui")
	})
	reg.Register("GUI.LOADICONSMEM", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("GUI.LOADICONSMEM: requires CGO raygui")
	})
	_ = m
}
