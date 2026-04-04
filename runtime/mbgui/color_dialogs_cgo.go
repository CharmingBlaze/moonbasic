//go:build cgo

package mbgui

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerColorAndDialogs(m *Module, reg runtime.Registrar) {
	reg.Register("GUI.COLORPANEL", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 9 {
			return value.Nil, fmt.Errorf("GUI.COLORPANEL expects (x,y,w,h, text$, r,g,b,a)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		col, err := colorArgs(args, 5)
		if err != nil {
			return value.Nil, err
		}
		out := raygui.ColorPanel(b, text, col)
		return allocRGBA(m, out)
	})
	reg.Register("GUI.COLORBARALPHA", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("GUI.COLORBARALPHA expects (x,y,w,h, text$, alpha#)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		a, ok := argF32(args[5])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.COLORBARALPHA: alpha must be numeric")
		}
		out := raygui.ColorBarAlpha(b, text, a)
		return rt.RetFloat(float64(out)), nil
	})
	reg.Register("GUI.COLORBARHUE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("GUI.COLORBARHUE expects (x,y,w,h, text$, value#)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		v, ok := argF32(args[5])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.COLORBARHUE: value must be numeric")
		}
		out := raygui.ColorBarHue(b, text, v)
		return rt.RetFloat(float64(out)), nil
	})
	reg.Register("GUI.COLORPICKER", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 9 {
			return value.Nil, fmt.Errorf("GUI.COLORPICKER expects (x,y,w,h, text$, r,g,b,a)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		col, err := colorArgs(args, 5)
		if err != nil {
			return value.Nil, err
		}
		out := raygui.ColorPicker(b, text, col)
		return allocRGBA(m, out)
	})
	reg.Register("GUI.COLORPICKERHSV", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("GUI.COLORPICKERHSV expects (x,y,w,h, text$, hsvHandle) 3 floats XYZ=HSV")
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
			return value.Nil, fmt.Errorf("GUI.COLORPICKERHSV: hsvHandle must be 3-element numeric array")
		}
		h := heap.Handle(args[5].IVal)
		a, err := heap.Cast[*heap.Array](m.h, h)
		if err != nil {
			return value.Nil, err
		}
		if a.TotalElements() < 3 {
			return value.Nil, fmt.Errorf("GUI.COLORPICKERHSV: need 3 elements")
		}
		x, e1 := a.Get([]int64{0})
		y, e2 := a.Get([]int64{1})
		z, e3 := a.Get([]int64{2})
		if err := firstErr(e1, e2, e3); err != nil {
			return value.Nil, err
		}
		v := rl.Vector3{X: float32(x), Y: float32(y), Z: float32(z)}
		res := raygui.ColorPickerHSV(b, text, &v)
		_ = a.Set([]int64{0}, float64(v.X))
		_ = a.Set([]int64{1}, float64(v.Y))
		_ = a.Set([]int64{2}, float64(v.Z))
		return rt.RetInt(int64(res)), nil
	})
	reg.Register("GUI.COLORPANELHSV", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("GUI.COLORPANELHSV expects (x,y,w,h, text$, hsvHandle)")
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
			return value.Nil, fmt.Errorf("GUI.COLORPANELHSV: hsvHandle must be 3-element numeric array")
		}
		h := heap.Handle(args[5].IVal)
		a, err := heap.Cast[*heap.Array](m.h, h)
		if err != nil {
			return value.Nil, err
		}
		if a.TotalElements() < 3 {
			return value.Nil, fmt.Errorf("GUI.COLORPANELHSV: need 3 elements")
		}
		x, e1 := a.Get([]int64{0})
		y, e2 := a.Get([]int64{1})
		z, e3 := a.Get([]int64{2})
		if err := firstErr(e1, e2, e3); err != nil {
			return value.Nil, err
		}
		v := rl.Vector3{X: float32(x), Y: float32(y), Z: float32(z)}
		res := raygui.ColorPanelHSV(b, text, &v)
		_ = a.Set([]int64{0}, float64(v.X))
		_ = a.Set([]int64{1}, float64(v.Y))
		_ = a.Set([]int64{2}, float64(v.Z))
		return rt.RetInt(int64(res)), nil
	})
	reg.Register("GUI.MESSAGEBOX", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 7 {
			return value.Nil, fmt.Errorf("GUI.MESSAGEBOX expects (x,y,w,h, title$, message$, buttons$)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		title, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		msg, err := rt.ArgString(args, 5)
		if err != nil {
			return value.Nil, err
		}
		btns, err := rt.ArgString(args, 6)
		if err != nil {
			return value.Nil, err
		}
		res := raygui.MessageBox(b, title, msg, btns)
		return rt.RetInt(int64(res)), nil
	})
	reg.Register("GUI.TEXTINPUTBOX", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 10 {
			return value.Nil, fmt.Errorf("GUI.TEXTINPUTBOX expects (x,y,w,h, title$, message$, buttons$, text$, maxLen, secretHandle)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		title, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		msg, err := rt.ArgString(args, 5)
		if err != nil {
			return value.Nil, err
		}
		btns, err := rt.ArgString(args, 6)
		if err != nil {
			return value.Nil, err
		}
		buf := stringFromRT(rt, args[7])
		maxLen, ok := argI32(args[8])
		if !ok || maxLen < 1 {
			return value.Nil, fmt.Errorf("GUI.TEXTINPUTBOX: maxLen must be positive")
		}
		if args[9].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("GUI.TEXTINPUTBOX: secretHandle must be 1-element numeric array (0/1)")
		}
		h := heap.Handle(args[9].IVal)
		sf, err := readFloat1(m.h, h)
		if err != nil {
			return value.Nil, err
		}
		secret := sf != 0
		btn := raygui.TextInputBox(b, title, msg, btns, &buf, int32(maxLen), &secret)
		lastTextInputBuf = buf
		lastTextInputBtn = btn
		if err := writeFloat1(m.h, h, boolAsFloat(secret)); err != nil {
			return value.Nil, err
		}
		return rt.RetInt(int64(btn)), nil
	})
	reg.Register("GUI.TEXTINPUTLAST$", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.TEXTINPUTLAST$ expects 0 arguments (after GUI.TEXTINPUTBOX)")
		}
		return rt.RetString(lastTextInputBuf), nil
	})
	reg.Register("GUI.GRID", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 8 {
			return value.Nil, fmt.Errorf("GUI.GRID expects (x,y,w,h, text$, spacing#, subdivs, cellHandle)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		sp, ok1 := argF32(args[5])
		sub, ok2 := argI32(args[6])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("GUI.GRID: spacing and subdivs must be numeric")
		}
		if args[7].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("GUI.GRID: cellHandle must be 2-element numeric array (mouse cell)")
		}
		h := heap.Handle(args[7].IVal)
		cx, cy, err := readFloat2(m.h, h)
		if err != nil {
			return value.Nil, err
		}
		cell := rl.Vector2{X: float32(cx), Y: float32(cy)}
		res := raygui.Grid(b, text, sp, sub, &cell)
		if err := writeFloat2(m.h, h, float64(cell.X), float64(cell.Y)); err != nil {
			return value.Nil, err
		}
		return rt.RetInt(int64(res)), nil
	})
}
