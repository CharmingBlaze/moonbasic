//go:build cgo

package mbgui

import (
	"fmt"
	"os"

	"github.com/gen2brain/raylib-go/raygui"

	"moonbasic/runtime"
	mbfont "moonbasic/runtime/font"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerGlobalAndStyle(m *Module, reg runtime.Registrar) {
	reg.Register("GUI.ENABLE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.ENABLE expects 0 arguments")
		}
		raygui.Enable()
		return value.Nil, nil
	})
	reg.Register("GUI.DISABLE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.DISABLE expects 0 arguments")
		}
		raygui.Disable()
		return value.Nil, nil
	})
	reg.Register("GUI.LOCK", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.LOCK expects 0 arguments")
		}
		raygui.Lock()
		return value.Nil, nil
	})
	reg.Register("GUI.UNLOCK", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.UNLOCK expects 0 arguments")
		}
		raygui.Unlock()
		return value.Nil, nil
	})
	reg.Register("GUI.ISLOCKED", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.ISLOCKED expects 0 arguments")
		}
		return rt.RetBool(raygui.IsLocked()), nil
	})
	reg.Register("GUI.SETALPHA", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.SETALPHA expects (alpha#)")
		}
		a, ok := argF32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.SETALPHA: alpha must be numeric")
		}
		raygui.SetAlpha(a)
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
		raygui.SetState(raygui.PropertyValue(s))
		return value.Nil, nil
	})
	reg.Register("GUI.GETSTATE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.GETSTATE expects 0 arguments")
		}
		return rt.RetInt(int64(raygui.GetState())), nil
	})
	reg.Register("GUI.SETFONT", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 || args[0].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("GUI.SETFONT expects (fontHandle)")
		}
		f, err := mbfont.FontForHandle(m.h, heap.Handle(args[0].IVal))
		if err != nil {
			return value.Nil, err
		}
		raygui.SetFont(f)
		return value.Nil, nil
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
		raygui.SetStyle(raygui.ControlID(c), raygui.PropertyID(p), raygui.PropertyValue(v))
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
		v := raygui.GetStyle(raygui.ControlID(c), raygui.PropertyID(p))
		return rt.RetInt(int64(v)), nil
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
		col := raygui.GetColor(raygui.ControlID(c), raygui.PropertyID(p))
		return allocRGBA(m, col)
	})
	reg.Register("GUI.LOADSTYLE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.LOADSTYLE expects (path$)")
		}
		path, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		raygui.LoadStyle(path)
		return value.Nil, nil
	})
	reg.Register("GUI.LOADDEFAULTSTYLE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.LOADDEFAULTSTYLE expects 0 arguments")
		}
		raygui.LoadStyleDefault()
		return value.Nil, nil
	})
	reg.Register("GUI.LOADSTYLEMEM", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.LOADSTYLEMEM expects (path$) to a binary .rgs file")
		}
		path, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return value.Nil, err
		}
		if len(data) == 0 {
			return value.Nil, fmt.Errorf("GUI.LOADSTYLEMEM: empty file")
		}
		raygui.LoadStyleFromMemory(data)
		return value.Nil, nil
	})
	reg.Register("GUI.LOADICONS", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("GUI.LOADICONS expects (path$, loadNames?)")
		}
		path, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		names, err := rt.ArgBool(args, 1)
		if err != nil {
			return value.Nil, err
		}
		raygui.LoadIcons(path, names)
		return value.Nil, nil
	})
	reg.Register("GUI.LOADICONSMEM", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("GUI.LOADICONSMEM expects (path$, loadNames?)")
		}
		path, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		names, err := rt.ArgBool(args, 1)
		if err != nil {
			return value.Nil, err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return value.Nil, err
		}
		if len(data) == 0 {
			return value.Nil, fmt.Errorf("GUI.LOADICONSMEM: empty file")
		}
		raygui.LoadIconsFromMemory(data, names)
		return value.Nil, nil
	})
}
