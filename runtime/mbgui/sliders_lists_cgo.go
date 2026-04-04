//go:build cgo

package mbgui

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raygui"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerSlidersAndLists(m *Module, reg runtime.Registrar) {
	reg.Register("GUI.SLIDER", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 9 {
			return value.Nil, fmt.Errorf("GUI.SLIDER expects (x,y,w,h, left$, right$, value#, min#, max#)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		ls, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		rs, err := rt.ArgString(args, 5)
		if err != nil {
			return value.Nil, err
		}
		val, ok1 := argF32(args[6])
		minV, ok2 := argF32(args[7])
		maxV, ok3 := argF32(args[8])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("GUI.SLIDER: value/min/max must be numeric")
		}
		out := raygui.Slider(b, ls, rs, val, minV, maxV)
		return rt.RetFloat(float64(out)), nil
	})
	reg.Register("GUI.SLIDERBAR", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 9 {
			return value.Nil, fmt.Errorf("GUI.SLIDERBAR expects (x,y,w,h, left$, right$, value#, min#, max#)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		ls, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		rs, err := rt.ArgString(args, 5)
		if err != nil {
			return value.Nil, err
		}
		val, ok1 := argF32(args[6])
		minV, ok2 := argF32(args[7])
		maxV, ok3 := argF32(args[8])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("GUI.SLIDERBAR: value/min/max must be numeric")
		}
		out := raygui.SliderBar(b, ls, rs, val, minV, maxV)
		return rt.RetFloat(float64(out)), nil
	})
	reg.Register("GUI.PROGRESSBAR", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 9 {
			return value.Nil, fmt.Errorf("GUI.PROGRESSBAR expects (x,y,w,h, left$, right$, value#, min#, max#)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		ls, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		rs, err := rt.ArgString(args, 5)
		if err != nil {
			return value.Nil, err
		}
		val, ok1 := argF32(args[6])
		minV, ok2 := argF32(args[7])
		maxV, ok3 := argF32(args[8])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("GUI.PROGRESSBAR: value/min/max must be numeric")
		}
		out := raygui.ProgressBar(b, ls, rs, val, minV, maxV)
		return rt.RetFloat(float64(out)), nil
	})
	reg.Register("GUI.SCROLLBAR", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 7 {
			return value.Nil, fmt.Errorf("GUI.SCROLLBAR expects (x,y,w,h, value, min, max)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		val, ok1 := argI32(args[4])
		minV, ok2 := argI32(args[5])
		maxV, ok3 := argI32(args[6])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("GUI.SCROLLBAR: value/min/max must be numeric")
		}
		out := raygui.ScrollBar(b, val, minV, maxV)
		return rt.RetInt(int64(out)), nil
	})
	reg.Register("GUI.STATUSBAR", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("GUI.STATUSBAR expects (x,y,w,h, text$)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		raygui.StatusBar(b, text)
		return value.Nil, nil
	})
	reg.Register("GUI.DUMMYREC", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("GUI.DUMMYREC expects (x,y,w,h, text$)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		raygui.DummyRec(b, text)
		return value.Nil, nil
	})
	reg.Register("GUI.LISTVIEW", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("GUI.LISTVIEW expects (x,y,w,h, items$, stateHandle) [scroll, active]")
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
			return value.Nil, fmt.Errorf("GUI.LISTVIEW: stateHandle must be numeric array len>=2")
		}
		h := heap.Handle(args[5].IVal)
		sf, af, err := readFloat2(m.h, h)
		if err != nil {
			return value.Nil, err
		}
		scroll := int32(sf)
		active := int32(af)
		sel := raygui.ListView(b, text, &scroll, active)
		if err := writeFloat2(m.h, h, float64(scroll), float64(sel)); err != nil {
			return value.Nil, err
		}
		return rt.RetInt(int64(sel)), nil
	})
	reg.Register("GUI.LISTVIEWEX", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("GUI.LISTVIEWEX expects (x,y,w,h, items$, stateHandle) [focus, scroll, active]")
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
			return value.Nil, fmt.Errorf("GUI.LISTVIEWEX: stateHandle must be numeric array len>=3")
		}
		h := heap.Handle(args[5].IVal)
		a, err := heap.Cast[*heap.Array](m.h, h)
		if err != nil {
			return value.Nil, err
		}
		if a.TotalElements() < 3 {
			return value.Nil, fmt.Errorf("GUI.LISTVIEWEX: state array needs 3 elements")
		}
		f0, e1 := a.Get([]int64{0})
		f1, e2 := a.Get([]int64{1})
		f2, e3 := a.Get([]int64{2})
		if err := firstErr(e1, e2, e3); err != nil {
			return value.Nil, err
		}
		focus := int32(f0)
		scroll := int32(f1)
		active := int32(f2)
		sel := raygui.ListViewEx(b, splitItems(text), &focus, &scroll, active)
		_ = a.Set([]int64{0}, float64(focus))
		_ = a.Set([]int64{1}, float64(scroll))
		_ = a.Set([]int64{2}, float64(sel))
		return rt.RetInt(int64(sel)), nil
	})
}
