//go:build !cgo && windows

package mbgui

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerPuregoSliders(m *Module, reg runtime.Registrar) {
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
			return value.Nil, fmt.Errorf("value/min/max must be numeric")
		}
		_ = ls
		_ = rs
		track := b
		track.Height = 24
		out := puregoDragValue(track, val, minV, maxV)
		puregoDrawButtonChrome(track, false)
		s := fmt.Sprintf("%.2f", out)
		puregoDrawLabelText(s, track, puregoBaseTextColor())
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
		_, err = rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		_, err = rt.ArgString(args, 5)
		if err != nil {
			return value.Nil, err
		}
		val, ok1 := argF32(args[6])
		minV, ok2 := argF32(args[7])
		maxV, ok3 := argF32(args[8])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("value/min/max must be numeric")
		}
		out := puregoDragValue(b, val, minV, maxV)
		fillW := b.Width * ((out - minV) / (maxV - minV + 1e-9))
		if fillW < 0 {
			fillW = 0
		}
		if fillW > b.Width {
			fillW = b.Width
		}
		rl.DrawRectangleRec(rl.Rectangle{X: b.X, Y: b.Y, Width: fillW, Height: b.Height}, puregoMulAlpha(puregoPanelColor(), pg.alpha))
		rl.DrawRectangleLinesEx(b, 1, puregoBaseTextColor())
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
		_, err = rt.ArgString(args, 4)
		if err != nil {
			return value.Nil, err
		}
		_, err = rt.ArgString(args, 5)
		if err != nil {
			return value.Nil, err
		}
		val, ok1 := argF32(args[6])
		minV, ok2 := argF32(args[7])
		maxV, ok3 := argF32(args[8])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("value/min/max must be numeric")
		}
		t := (val - minV) / (maxV - minV + 1e-9)
		if t < 0 {
			t = 0
		}
		if t > 1 {
			t = 1
		}
		rl.DrawRectangleRec(rl.Rectangle{X: b.X, Y: b.Y, Width: b.Width * t, Height: b.Height}, puregoMulAlpha(puregoBaseTextColor(), 0.35*pg.alpha))
		rl.DrawRectangleLinesEx(b, 1, puregoBaseTextColor())
		return rt.RetFloat(float64(val)), nil
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
			return value.Nil, fmt.Errorf("arguments must be numeric")
		}
		vf := float32(val)
		out := puregoDragValue(b, vf, float32(minV), float32(maxV))
		rl.DrawRectangleLinesEx(b, 1, puregoBaseTextColor())
		return rt.RetInt(int64(out)), nil
	})
	_ = m
}
