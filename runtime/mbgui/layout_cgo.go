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

func registerLayout(m *Module, reg runtime.Registrar) {
	reg.Register("GUI.WINDOWBOX", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("GUI.WINDOWBOX expects (x,y,w,h, title$)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		title, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		return rt.RetBool(raygui.WindowBox(b, title)), nil
	})
	reg.Register("GUI.GROUPBOX", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("GUI.GROUPBOX expects (x,y,w,h, text$)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		raygui.GroupBox(b, text)
		return value.Nil, nil
	})
	reg.Register("GUI.LINE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("GUI.LINE expects (x,y,w,h, text$)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		raygui.Line(b, text)
		return value.Nil, nil
	})
	reg.Register("GUI.PANEL", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("GUI.PANEL expects (x,y,w,h, text$)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		text, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		raygui.Panel(b, text)
		return value.Nil, nil
	})
	reg.Register("GUI.TABBAR", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("GUI.TABBAR expects (x,y,w,h, tabs$, stateHandle)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		tabs, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		if args[5].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("GUI.TABBAR: stateHandle must be a handle to a 1-element numeric array (active tab)")
		}
		h := heap.Handle(args[5].IVal)
		activeF, err := readFloat1(m.h, h)
		if err != nil {
			return value.Nil, err
		}
		active := int32(activeF)
		closeReq := raygui.TabBar(b, splitItems(tabs), &active)
		if err := writeFloat1(m.h, h, float64(active)); err != nil {
			return value.Nil, err
		}
		return rt.RetInt(int64(closeReq)), nil
	})
	reg.Register("GUI.SCROLLPANEL", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 10 {
			return value.Nil, fmt.Errorf("GUI.SCROLLPANEL expects (px,py,pw,ph, title$, cx,cy,cw,ch, stateHandle)")
		}
		pb, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		title, err := rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		cb, err := rectArgs(args, 5)
		if err != nil {
			return value.Nil, err
		}
		if args[9].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("GUI.SCROLLPANEL: stateHandle must be a 6-element array [sx,sy,vx,vy,vw,vh]")
		}
		h := heap.Handle(args[9].IVal)
		v6, err := readFloat6(m.h, h)
		if err != nil {
			return value.Nil, err
		}
		scroll := rl.Vector2{X: float32(v6[0]), Y: float32(v6[1])}
		view := rl.Rectangle{X: float32(v6[2]), Y: float32(v6[3]), Width: float32(v6[4]), Height: float32(v6[5])}
		raygui.ScrollPanel(pb, title, cb, &scroll, &view)
		v6[0] = float64(scroll.X)
		v6[1] = float64(scroll.Y)
		v6[2] = float64(view.X)
		v6[3] = float64(view.Y)
		v6[4] = float64(view.Width)
		v6[5] = float64(view.Height)
		if err := writeFloat6(m.h, h, v6); err != nil {
			return value.Nil, err
		}
		return value.Nil, nil
	})
}
