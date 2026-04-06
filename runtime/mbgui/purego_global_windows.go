//go:build !cgo && windows

package mbgui

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerPuregoGlobal(m *Module, reg runtime.Registrar) {
	reg.Register("GUI.ENABLE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.ENABLE expects 0 arguments")
		}
		pg.enabled = true
		return value.Nil, nil
	})
	reg.Register("GUI.DISABLE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.DISABLE expects 0 arguments")
		}
		pg.disabled = true
		return value.Nil, nil
	})
	reg.Register("GUI.LOCK", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.LOCK expects 0 arguments")
		}
		pg.locked = true
		return value.Nil, nil
	})
	reg.Register("GUI.UNLOCK", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.UNLOCK expects 0 arguments")
		}
		pg.locked = false
		return value.Nil, nil
	})
	reg.Register("GUI.ISLOCKED", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.ISLOCKED expects 0 arguments")
		}
		return rt.RetBool(pg.locked), nil
	})
	reg.Register("GUI.SETALPHA", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.SETALPHA expects (alpha#)")
		}
		a, ok := argF32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.SETALPHA: alpha must be numeric")
		}
		if a < 0 {
			a = 0
		}
		if a > 1 {
			a = 1
		}
		pg.alpha = a
		return value.Nil, nil
	})
	reg.Register("GUI.SETSTATE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.SETSTATE expects (state)")
		}
		s, ok := argI32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.SETSTATE: state must be numeric")
		}
		pg.guiState = s
		return value.Nil, nil
	})
	reg.Register("GUI.GETSTATE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.GETSTATE expects 0 arguments")
		}
		return rt.RetInt(int64(pg.guiState)), nil
	})
	reg.Register("GUI.SETSTYLE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("GUI.SETSTYLE expects (control, property, value)")
		}
		c, ok1 := argI32(args[0])
		p, ok2 := argI32(args[1])
		v, ok3 := argI32(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("GUI.SETSTYLE: arguments must be numeric")
		}
		pg.styleInt[styleKey(c, p)] = int64(v)
		return value.Nil, nil
	})
	reg.Register("GUI.GETSTYLE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("GUI.GETSTYLE expects (control, property)")
		}
		c, ok1 := argI32(args[0])
		p, ok2 := argI32(args[1])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("GUI.GETSTYLE: arguments must be numeric")
		}
		v := pg.styleInt[styleKey(c, p)]
		return rt.RetInt(v), nil
	})
	reg.Register("GUI.GETCOLOR", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("GUI.GETCOLOR expects (control, property)")
		}
		c, ok1 := argI32(args[0])
		p, ok2 := argI32(args[1])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("GUI.GETCOLOR: arguments must be numeric")
		}
		k := styleKey(c, p)
		col, ok := pg.styleColor[k]
		if !ok {
			col = puregoBaseTextColor()
		}
		return allocRGBA(m, col)
	})
	reg.Register("GUI.SETFONT", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return value.Nil, fmt.Errorf("GUI.SETFONT: custom fonts require CGO raygui; using default font in purego GUI")
	})
}
