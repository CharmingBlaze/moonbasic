//go:build cgo

package mbgui

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raygui"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerTextControls(m *Module, reg runtime.Registrar) {
	reg.Register("GUI.TEXTBOX", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 7 {
			return value.Nil, fmt.Errorf("GUI.TEXTBOX expects (x,y,w,h, text$, maxLen, editMode?)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		s := stringFromRT(rt, args[4])
		maxLen, ok := argI32(args[5])
		if !ok || maxLen < 1 {
			return value.Nil, fmt.Errorf("GUI.TEXTBOX: maxLen must be positive")
		}
		edit, err := rt.ArgBool(args, 6)
		if err != nil {
			return value.Nil, err
		}
		raygui.TextBox(b, &s, int(maxLen), edit)
		return rt.RetString(s), nil
	})
	reg.Register("GUI.SPINNER", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 9 {
			return value.Nil, fmt.Errorf("GUI.SPINNER expects (x,y,w,h, text$, value, min, max, editMode?)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		val, ok := argI32(args[5])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.SPINNER: value must be numeric")
		}
		minV, ok1 := argI32(args[6])
		maxV, ok2 := argI32(args[7])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("GUI.SPINNER: min/max must be numeric")
		}
		edit, err := rt.ArgBool(args, 8)
		if err != nil {
			return value.Nil, err
		}
		v := val
		raygui.Spinner(b, text, &v, int(minV), int(maxV), edit)
		return rt.RetInt(int64(v)), nil
	})
	reg.Register("GUI.VALUEBOX", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 9 {
			return value.Nil, fmt.Errorf("GUI.VALUEBOX expects (x,y,w,h, text$, value, min, max, editMode?)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		val, ok := argI32(args[5])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.VALUEBOX: value must be numeric")
		}
		minV, ok1 := argI32(args[6])
		maxV, ok2 := argI32(args[7])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("GUI.VALUEBOX: min/max must be numeric")
		}
		edit, err := rt.ArgBool(args, 8)
		if err != nil {
			return value.Nil, err
		}
		v := val
		raygui.ValueBox(b, text, &v, int(minV), int(maxV), edit)
		return rt.RetInt(int64(v)), nil
	})
	reg.Register("GUI.VALUEBOXFLOAT", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 8 {
			return value.Nil, fmt.Errorf("GUI.VALUEBOXFLOAT expects (x,y,w,h, label$, value#, textBuf$, editMode?)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		label, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		vf, ok := argF32(args[5])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.VALUEBOXFLOAT: value must be numeric")
		}
		buf := stringFromRT(rt, args[6])
		edit, err := rt.ArgBool(args, 7)
		if err != nil {
			return value.Nil, err
		}
		v := vf
		raygui.ValueBoxFloat(b, label, &buf, &v, edit)
		lastValueBoxFloat = float64(v)
		lastValueBoxFloatS = buf
		return rt.RetFloat(float64(v)), nil
	})
	reg.Register("GUI.VALUEBOXFLOATTEXT$", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.VALUEBOXFLOATTEXT$ expects 0 arguments (call after GUI.VALUEBOXFLOAT in the same frame)")
		}
		return rt.RetString(lastValueBoxFloatS), nil
	})
}
