//go:build !cgo && windows

package mbgui

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerPuregoTooltipDraw(m *Module, reg runtime.Registrar) {
	reg.Register("GUI.ENABLETOOLTIP", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.ENABLETOOLTIP expects 0 arguments")
		}
		pg.tooltipOn = true
		return value.Nil, nil
	})
	reg.Register("GUI.DISABLETOOLTIP", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.DISABLETOOLTIP expects 0 arguments")
		}
		pg.tooltipOn = false
		return value.Nil, nil
	})
	reg.Register("GUI.SETTOOLTIP", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.SETTOOLTIP expects (text$)")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		pg.tooltipText = s
		return value.Nil, nil
	})
	reg.Register("GUI.ICONTEXT", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("GUI.ICONTEXT expects (iconId, text$)")
		}
		id, ok := argI32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("iconId must be numeric")
		}
		s, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Nil, err
		}
		return rt.RetString(fmt.Sprintf("#%d %s", id, s)), nil
	})
	reg.Register("GUI.DRAWICON", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 8 {
			return value.Nil, fmt.Errorf("GUI.DRAWICON expects (iconId, x, y, pixelSize, r,g,b,a)")
		}
		id, ok := argI32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("iconId must be numeric")
		}
		x, ok1 := argF32(args[1])
		y, ok2 := argF32(args[2])
		ps, ok3 := argF32(args[3])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("position/size must be numeric")
		}
		col, err := colorArgs(args, 4)
		if err != nil {
			return value.Nil, err
		}
		rgba := color.RGBA{R: col.R, G: col.G, B: col.B, A: col.A}
		s := fmt.Sprintf("%d", id)
		font := rl.GetFontDefault()
		rl.DrawTextEx(font, s, rl.Vector2{X: x, Y: y}, ps, 1, rgba)
		return value.Nil, nil
	})
	reg.Register("GUI.SETICONSCALE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.SETICONSCALE expects (scale)")
		}
		s, ok := argI32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("scale must be numeric")
		}
		pg.iconScale = s
		return value.Nil, nil
	})
	reg.Register("GUI.GETTEXTWIDTH", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.GETTEXTWIDTH expects (text$)")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		font := rl.GetFontDefault()
		m := rl.MeasureTextEx(font, s, pg.textSize, pg.textSpacing)
		return rt.RetInt(int64(m.X)), nil
	})
	reg.Register("GUI.FADE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("GUI.FADE expects (r,g,b,a, alpha#)")
		}
		col, err := colorArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		a, ok := argF32(args[4])
		if !ok {
			return value.Nil, fmt.Errorf("alpha must be numeric")
		}
		out := puregoMulAlpha(col, a)
		return allocRGBA(m, out)
	})
	reg.Register("GUI.DRAWRECTANGLE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 13 {
			return value.Nil, fmt.Errorf("GUI.DRAWRECTANGLE expects (x,y,w,h, borderW, br,bg,bb,ba, fr,fg,fb,fa)")
		}
		b, err := rectArgs(args, 0)
		if err != nil {
			return value.Nil, err
		}
		bw, ok := argI32(args[4])
		if !ok {
			return value.Nil, fmt.Errorf("border width must be numeric")
		}
		bc, err := colorArgs(args, 5)
		if err != nil {
			return value.Nil, err
		}
		fc, err := colorArgs(args, 9)
		if err != nil {
			return value.Nil, err
		}
		rl.DrawRectangleRec(b, fc)
		rl.DrawRectangleLinesEx(b, float32(bw), bc)
		return value.Nil, nil
	})
	reg.Register("GUI.DRAWTEXT", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 9 {
			return value.Nil, fmt.Errorf("GUI.DRAWTEXT expects (text$, x,y,w,h, align, r,g,b,a)")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		rb, err := rectArgs(args, 1)
		if err != nil {
			return value.Nil, err
		}
		al, ok := argI32(args[5])
		if !ok {
			return value.Nil, fmt.Errorf("align must be numeric")
		}
		col, err := colorArgs(args, 6)
		if err != nil {
			return value.Nil, err
		}
		saved := pg.alignH
		pg.alignH = al
		puregoDrawLabelText(s, rb, col)
		pg.alignH = saved
		return value.Nil, nil
	})
	reg.Register("GUI.GETTEXTBOUNDS", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("GUI.GETTEXTBOUNDS expects (control, x,y,w,h)")
		}
		_, ok := argI32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("control must be numeric")
		}
		b, err := rectArgs(args, 1)
		if err != nil {
			return value.Nil, err
		}
		font := rl.GetFontDefault()
		ms := rl.MeasureTextEx(font, "Mg", pg.textSize, pg.textSpacing)
		return allocRect(m, rl.Rectangle{X: b.X + 2, Y: b.Y + (b.Height-ms.Y)/2, Width: ms.X, Height: ms.Y})
	})
}
