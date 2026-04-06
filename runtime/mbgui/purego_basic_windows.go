//go:build !cgo && windows

package mbgui

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerPuregoBasic(m *Module, reg runtime.Registrar) {
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
		puregoDrawLabelText(text, b, puregoBaseTextColor())
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
		pressed := puregoClickIn(b)
		puregoDrawButtonChrome(b, pressed)
		puregoDrawLabelText(text, b, puregoBaseTextColor())
		return rt.RetBool(pressed), nil
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
		hov := rl.CheckCollisionPointRec(rl.GetMousePosition(), b)
		puregoDrawButtonChrome(b, hov && rl.IsMouseButtonDown(rl.MouseLeftButton))
		puregoDrawLabelText(text, b, puregoBaseTextColor())
		return rt.RetBool(puregoClickIn(b)), nil
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
		if puregoClickIn(b) {
			active = !active
		}
		box := rl.Rectangle{X: b.X, Y: b.Y + (b.Height-18)/2, Width: 18, Height: 18}
		rl.DrawRectangleLinesEx(box, 1, puregoBaseTextColor())
		if active {
			rl.DrawRectangleRec(rl.Rectangle{X: box.X + 4, Y: box.Y + 4, Width: 10, Height: 10}, puregoBaseTextColor())
		}
		tb := rl.Rectangle{X: b.X + 24, Y: b.Y, Width: b.Width - 24, Height: b.Height}
		puregoDrawLabelText(text, tb, puregoBaseTextColor())
		return rt.RetBool(active), nil
	})
	reg.Register("GUI.TOGGLEGROUP", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("GUI.TOGGLEGROUP expects (x,y,w,h, items$)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		items := splitItems(text)
		if len(items) == 0 {
			return rt.RetInt(0), nil
		}
		k := makePuregoRectKey(b)
		active := pg.widgetInt[k]
		if puregoClickIn(b) {
			active = (active + 1) % int32(len(items))
			pg.widgetInt[k] = active
			puregoNoteWidgetWrite()
		}
		puregoDrawButtonChrome(b, false)
		puregoDrawLabelText(items[active], b, puregoBaseTextColor())
		return rt.RetInt(int64(active)), nil
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
			return value.Nil, fmt.Errorf("active must be numeric")
		}
		items := splitItems(text)
		if len(items) == 0 {
			return rt.RetInt(0), nil
		}
		active := a % int32(len(items))
		if active < 0 {
			active += int32(len(items))
		}
		if puregoClickIn(b) {
			active = (active + 1) % int32(len(items))
		}
		puregoDrawButtonChrome(b, false)
		puregoDrawLabelText(items[active], b, puregoBaseTextColor())
		return rt.RetInt(int64(active)), nil
	})
	reg.Register("GUI.TOGGLESLIDER", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("GUI.TOGGLESLIDER: requires CGO raygui")
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
		if puregoClickIn(b) {
			checked = !checked
		}
		box := rl.Rectangle{X: b.X, Y: b.Y + (b.Height-18)/2, Width: 18, Height: 18}
		rl.DrawRectangleLinesEx(box, 1, puregoBaseTextColor())
		if checked {
			rl.DrawRectangleRec(rl.Rectangle{X: box.X + 3, Y: box.Y + 3, Width: 12, Height: 12}, puregoBaseTextColor())
		}
		tb := rl.Rectangle{X: b.X + 24, Y: b.Y, Width: b.Width - 24, Height: b.Height}
		puregoDrawLabelText(text, tb, puregoBaseTextColor())
		return rt.RetBool(checked), nil
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
			return value.Nil, fmt.Errorf("active must be numeric")
		}
		items := splitItems(text)
		if len(items) == 0 {
			return rt.RetInt(0), nil
		}
		active := a % int32(len(items))
		if active < 0 {
			active += int32(len(items))
		}
		if puregoClickIn(b) {
			active = (active + 1) % int32(len(items))
		}
		puregoDrawButtonChrome(b, false)
		puregoDrawLabelText(items[active], b, puregoBaseTextColor())
		return rt.RetInt(int64(active)), nil
	})
	reg.Register("GUI.DROPDOWNBOX", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("GUI.DROPDOWNBOX expects (x,y,w,h, items$, stateHandle)")
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
			return value.Nil, fmt.Errorf("stateHandle must be handle")
		}
		h := heap.Handle(args[5].IVal)
		af, bf, err := readFloat2(m.h, h)
		if err != nil {
			return value.Nil, err
		}
		active := int32(af)
		edit := bf != 0
		items := splitItems(text)
		if len(items) == 0 {
			return rt.RetBool(false), nil
		}
		if active < 0 || active >= int32(len(items)) {
			active = 0
		}
		open := edit
		if puregoClickIn(b) {
			open = !open
			if !edit {
				active = (active + 1) % int32(len(items))
			}
		}
		puregoDrawButtonChrome(b, open)
		puregoDrawLabelText(items[active], b, puregoBaseTextColor())
		if err := writeFloat2(m.h, h, float64(active), boolAsFloat(open)); err != nil {
			return value.Nil, err
		}
		return rt.RetBool(open), nil
	})
	_ = m
}
