//go:build !cgo && windows

package mbgui

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerPuregoLayout(m *Module, reg runtime.Registrar) {
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
		head := float32(24)
		if b.Height < head+4 {
			head = b.Height / 2
		}
		headRec := rl.Rectangle{X: b.X, Y: b.Y, Width: b.Width, Height: head}
		bodyRec := rl.Rectangle{X: b.X, Y: b.Y + head, Width: b.Width, Height: b.Height - head}
		rl.DrawRectangleRec(headRec, puregoMulAlpha(puregoPanelColor(), pg.alpha))
		rl.DrawRectangleLinesEx(b, 1, puregoBaseTextColor())
		puregoDrawLabelText(title, headRec, puregoBaseTextColor())
		closeRec := rl.Rectangle{X: b.X + b.Width - head, Y: b.Y, Width: head, Height: head}
		rl.DrawRectangleLinesEx(closeRec, 1, puregoBaseTextColor())
		puregoDrawLabelText("x", closeRec, puregoBaseTextColor())
		rl.DrawRectangleRec(bodyRec, puregoMulAlpha(rl.Color{R: 28, G: 30, B: 34, A: 255}, 0.9*pg.alpha))
		closed := puregoClickIn(closeRec)
		return rt.RetBool(closed), nil
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
		rl.DrawRectangleLinesEx(b, 1, puregoBaseTextColor())
		puregoDrawLabelText(text, rl.Rectangle{X: b.X + 6, Y: b.Y - 10, Width: b.Width - 12, Height: 20}, puregoBaseTextColor())
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
		mid := b.Y + b.Height/2
		rl.DrawLineEx(rl.Vector2{X: b.X, Y: mid}, rl.Vector2{X: b.X + b.Width, Y: mid}, 1, puregoBaseTextColor())
		puregoDrawLabelText(text, b, puregoBaseTextColor())
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
		rl.DrawRectangleRec(b, puregoMulAlpha(puregoPanelColor(), 0.85*pg.alpha))
		rl.DrawRectangleLinesEx(b, 1, puregoBaseTextColor())
		puregoDrawLabelText(text, b, puregoBaseTextColor())
		return value.Nil, nil
	})
	reg.Register("GUI.TABBAR", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("GUI.TABBAR: requires CGO raygui")
	})
	reg.Register("GUI.SCROLLPANEL", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("GUI.SCROLLPANEL: requires CGO raygui")
	})
	_ = m
}
