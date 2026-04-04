//go:build cgo

package mbgui

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raygui"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerBasicControls(m *Module, reg runtime.Registrar) {
	reg.Register("GUI.LABEL", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("GUI.LABEL expects (x,y,w,h, text$)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		raygui.Label(b, text)
		return value.Nil, nil
	})
	reg.Register("GUI.BUTTON", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("GUI.BUTTON expects (x,y,w,h, text$)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		return rt.RetBool(raygui.Button(b, text)), nil
	})
	reg.Register("GUI.LABELBUTTON", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("GUI.LABELBUTTON expects (x,y,w,h, text$)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		return rt.RetBool(raygui.LabelButton(b, text)), nil
	})
	reg.Register("GUI.TOGGLE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("GUI.TOGGLE expects (x,y,w,h, text$, active?)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		active, err := rt.ArgBool(args, 5)
		if err != nil {
			return value.Nil, err
		}
		return rt.RetBool(raygui.Toggle(b, text, active)), nil
	})
	reg.Register("GUI.TOGGLEGROUP", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("GUI.TOGGLEGROUP expects (x,y,w,h, items$) semicolon-separated")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		return rt.RetInt(int64(raygui.ToggleGroup(b, text, 0))), nil
	})
	reg.Register("GUI.TOGGLEGROUPAT", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("GUI.TOGGLEGROUPAT expects (x,y,w,h, items$, active)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		a, ok := argI32(args[5])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.TOGGLEGROUPAT: active must be numeric")
		}
		return rt.RetInt(int64(raygui.ToggleGroup(b, text, a))), nil
	})
	reg.Register("GUI.TOGGLESLIDER", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("GUI.TOGGLESLIDER expects (x,y,w,h, text$, active)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		a, ok := argI32(args[5])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.TOGGLESLIDER: active must be numeric")
		}
		return rt.RetInt(int64(raygui.ToggleSlider(b, text, a))), nil
	})
	reg.Register("GUI.CHECKBOX", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("GUI.CHECKBOX expects (x,y,w,h, text$, checked?)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		checked, err := rt.ArgBool(args, 5)
		if err != nil {
			return value.Nil, err
		}
		return rt.RetBool(raygui.CheckBox(b, text, checked)), nil
	})
	reg.Register("GUI.COMBOBOX", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("GUI.COMBOBOX expects (x,y,w,h, items$, active)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		a, ok := argI32(args[5])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.COMBOBOX: active must be numeric")
		}
		return rt.RetInt(int64(raygui.ComboBox(b, text, a))), nil
	})
	reg.Register("GUI.DROPDOWNBOX", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("GUI.DROPDOWNBOX expects (x,y,w,h, items$, stateHandle) 2 floats: active, editMode")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		if args[5].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("GUI.DROPDOWNBOX: stateHandle must be numeric array len>=2")
		}
		h := heap.Handle(args[5].IVal)
		af, bf, err := readFloat2(m.h, h)
		if err != nil {
			return value.Nil, err
		}
		active := int32(af)
		edit := bf != 0
		open := raygui.DropdownBox(b, text, &active, edit)
		if err := writeFloat2(m.h, h, float64(active), boolAsFloat(open)); err != nil {
			return value.Nil, err
		}
		return rt.RetBool(open), nil
	})
}

func boolAsFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
