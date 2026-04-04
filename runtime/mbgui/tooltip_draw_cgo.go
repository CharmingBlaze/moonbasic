//go:build cgo

package mbgui

import (
	"fmt"
	"image/color"

	"github.com/gen2brain/raylib-go/raygui"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerTooltipIconsDraw(m *Module, reg runtime.Registrar) {
	reg.Register("GUI.ENABLETOOLTIP", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.ENABLETOOLTIP expects 0 arguments")
		}
		raygui.EnableTooltip()
		return value.Nil, nil
	})
	reg.Register("GUI.DISABLETOOLTIP", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.DISABLETOOLTIP expects 0 arguments")
		}
		raygui.DisableTooltip()
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
		raygui.SetTooltip(s)
		return value.Nil, nil
	})
	reg.Register("GUI.ICONTEXT", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("GUI.ICONTEXT expects (iconId, text$)")
		}
		id, ok := argI32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.ICONTEXT: iconId must be numeric")
		}
		s, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Nil, err
		}
		return rt.RetString(raygui.IconText(raygui.IconID(id), s)), nil
	})
	reg.Register("GUI.DRAWICON", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 8 {
			return value.Nil, fmt.Errorf("GUI.DRAWICON expects (iconId, x, y, pixelSize, r,g,b,a)")
		}
		id, ok := argI32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.DRAWICON: iconId must be numeric")
		}
		x, ok1 := argI32(args[1])
		y, ok2 := argI32(args[2])
		ps, ok3 := argI32(args[3])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("GUI.DRAWICON: position/size must be numeric")
		}
		col, err := colorArgs(args, 4)
		if err != nil {
			return value.Nil, err
		}
		rgba := color.RGBA{R: col.R, G: col.G, B: col.B, A: col.A}
		raygui.DrawIcon(raygui.IconID(id), x, y, ps, rgba)
		return value.Nil, nil
	})
	reg.Register("GUI.SETICONSCALE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.SETICONSCALE expects (scale)")
		}
		s, ok := argI32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.SETICONSCALE: scale must be numeric")
		}
		raygui.SetIconScale(s)
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
		return rt.RetInt(int64(raygui.GetTextWidth(s))), nil
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
			return value.Nil, fmt.Errorf("GUI.FADE: alpha must be numeric")
		}
		out := raygui.Fade(col, a)
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
			return value.Nil, fmt.Errorf("GUI.DRAWRECTANGLE: border width must be numeric")
		}
		bc, err := colorArgs(args, 5)
		if err != nil {
			return value.Nil, err
		}
		fc, err := colorArgs(args, 9)
		if err != nil {
			return value.Nil, err
		}
		raygui.DrawRectangle(b, bw, bc, fc)
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
			return value.Nil, fmt.Errorf("GUI.DRAWTEXT: align must be numeric")
		}
		col, err := colorArgs(args, 6)
		if err != nil {
			return value.Nil, err
		}
		raygui.DrawText(s, rb, al, col)
		return value.Nil, nil
	})
	reg.Register("GUI.GETTEXTBOUNDS", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("GUI.GETTEXTBOUNDS expects (control, x,y,w,h)")
		}
		cid, ok := argI32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.GETTEXTBOUNDS: control must be numeric (raygui.ControlID)")
		}
		b, err := rectArgs(args, 1)
		if err != nil {
			return value.Nil, err
		}
		r := raygui.GetTextBounds(raygui.ControlID(cid), b)
		return allocRect(m, r)
	})
}
