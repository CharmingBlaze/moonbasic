//go:build cgo || (windows && !cgo)

package mbdraw

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/convert"
	"moonbasic/vm/value"
)

func registerDrawHelperCmds(m *Module, r runtime.Registrar) {
	r.Register("DRAW.PROGRESSBAR", "draw", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.drawProgressBar(rt, args)
	})
	r.Register("DRAW.HEALTHBAR", "draw", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.drawHealthBar(rt, args)
	})
	r.Register("DRAW.CENTERTEXT", "draw", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.drawCenterText(rt, args)
	})
	r.Register("DRAW.RIGHTTEXT", "draw", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.drawRightText(rt, args)
	})
	r.Register("DRAW.SHADOWTEXT", "draw", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.drawShadowText(rt, args)
	})
	r.Register("DRAW.OUTLINETEXT", "draw", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.drawOutlineText(rt, args)
	})
	r.Register("DRAW.CROSSHAIR", "draw", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.drawCrosshair(rt, args)
	})
	r.Register("DRAW.RECTGRID", "draw", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.drawRectGrid(rt, args)
	})
}

func rgbaAt(args []value.Value, offset int) (color.RGBA, bool) {
	if len(args) < offset+4 {
		return color.RGBA{}, false
	}
	var rgb [4]int32
	for i := 0; i < 4; i++ {
		v, ok := argInt(args[offset+i])
		if !ok {
			return color.RGBA{}, false
		}
		rgb[i] = v
	}
	c := convert.NewColor4(rgb[0], rgb[1], rgb[2], rgb[3])
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}, true
}

func (m *Module) drawProgressBar(rt *runtime.Runtime, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 9 {
		return value.Nil, fmt.Errorf("DRAW.PROGRESSBAR expects (x, y, w, h, progress#, r, g, b, a)")
	}
	x, ok1 := argInt(args[0])
	y, ok2 := argInt(args[1])
	w, ok3 := argInt(args[2])
	h, ok4 := argInt(args[3])
	t, ok5 := args[4].ToFloat()
	if !ok5 {
		if ti, ok := args[4].ToInt(); ok {
			t = float64(ti)
			ok5 = true
		}
	}
	fg, ok6 := rgbaAt(args, 5)
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("DRAW.PROGRESSBAR: invalid arguments")
	}
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	ix, iy := int32(x), int32(y)
	iw, ih := int32(w), int32(h)
	bg := color.RGBA{R: 40, G: 40, B: 48, A: 255}
	rl.DrawRectangle(ix, iy, iw, ih, bg)
	innerW := int32(float64(iw) * t)
	if innerW > 0 {
		rl.DrawRectangle(ix, iy, innerW, ih, fg)
	}
	rl.DrawRectangleLines(ix, iy, iw, ih, rl.Fade(rl.White, 0.35))
	return value.Nil, nil
}

func (m *Module) drawHealthBar(rt *runtime.Runtime, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 10 {
		return value.Nil, fmt.Errorf("DRAW.HEALTHBAR expects (x, y, w, h, current#, max#, r, g, b, a)")
	}
	x, ok1 := argInt(args[0])
	y, ok2 := argInt(args[1])
	w, ok3 := argInt(args[2])
	h, ok4 := argInt(args[3])
	cur, ok5 := args[4].ToFloat()
	max, ok6 := args[5].ToFloat()
	if !ok5 {
		if ci, ok := args[4].ToInt(); ok {
			cur = float64(ci)
			ok5 = true
		}
	}
	if !ok6 {
		if mi, ok := args[5].ToInt(); ok {
			max = float64(mi)
			ok6 = true
		}
	}
	fg, ok7 := rgbaAt(args, 6)
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 {
		return value.Nil, fmt.Errorf("DRAW.HEALTHBAR: invalid arguments")
	}
	if max <= 0 {
		max = 1
	}
	t := cur / max
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	ix, iy := int32(x), int32(y)
	iw, ih := int32(w), int32(h)
	bg := color.RGBA{R: 40, G: 40, B: 48, A: 255}
	rl.DrawRectangle(ix, iy, iw, ih, bg)
	innerW := int32(float64(iw) * t)
	if innerW > 0 {
		rl.DrawRectangle(ix, iy, innerW, ih, fg)
	}
	rl.DrawRectangleLines(ix, iy, iw, ih, rl.Fade(rl.White, 0.35))
	return value.Nil, nil
}

func (m *Module) drawCenterText(rt *runtime.Runtime, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("DRAW.CENTERTEXT expects (text$, y, size, r, g, b, a)")
	}
	text := stringFromRT(rt, args[0])
	y, ok2 := argInt(args[1])
	size, ok3 := argInt(args[2])
	c, ok4 := rgbaAt(args, 3)
	if !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAW.CENTERTEXT: invalid arguments")
	}
	sw := int32(rl.GetScreenWidth())
	tw := rl.MeasureText(text, size)
	x := (sw - tw) / 2
	if x < 0 {
		x = 0
	}
	rl.DrawText(text, x, int32(y), size, c)
	return value.Nil, nil
}

func (m *Module) drawRightText(rt *runtime.Runtime, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("DRAW.RIGHTTEXT expects (text$, marginRight, y, size, r, g, b, a)")
	}
	text := stringFromRT(rt, args[0])
	marg, ok2 := argInt(args[1])
	y, ok3 := argInt(args[2])
	size, ok4 := argInt(args[3])
	c, ok5 := rgbaAt(args, 4)
	if !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("DRAW.RIGHTTEXT: invalid arguments")
	}
	sw := int32(rl.GetScreenWidth())
	tw := rl.MeasureText(text, size)
	right := sw - int32(marg)
	x := right - tw
	if x < 0 {
		x = 0
	}
	rl.DrawText(text, x, int32(y), size, c)
	return value.Nil, nil
}

func (m *Module) drawShadowText(rt *runtime.Runtime, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("DRAW.SHADOWTEXT expects (text$, x, y, size, r, g, b, a)")
	}
	text := stringFromRT(rt, args[0])
	x, ok2 := argInt(args[1])
	y, ok3 := argInt(args[2])
	size, ok4 := argInt(args[3])
	c, ok5 := rgbaAt(args, 4)
	if !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("DRAW.SHADOWTEXT: invalid arguments")
	}
	ix, iy := int32(x), int32(y)
	sh := color.RGBA{R: 0, G: 0, B: 0, A: 200}
	rl.DrawText(text, ix+2, iy+2, size, sh)
	rl.DrawText(text, ix, iy, size, c)
	return value.Nil, nil
}

func (m *Module) drawOutlineText(rt *runtime.Runtime, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("DRAW.OUTLINETEXT expects (text$, x, y, size, r, g, b, a)")
	}
	text := stringFromRT(rt, args[0])
	x, ok2 := argInt(args[1])
	y, ok3 := argInt(args[2])
	size, ok4 := argInt(args[3])
	c, ok5 := rgbaAt(args, 4)
	if !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("DRAW.OUTLINETEXT: invalid arguments")
	}
	ix, iy := int32(x), int32(y)
	out := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			rl.DrawText(text, ix+int32(dx), iy+int32(dy), size, out)
		}
	}
	rl.DrawText(text, ix, iy, size, c)
	return value.Nil, nil
}

func (m *Module) drawCrosshair(rt *runtime.Runtime, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("DRAW.CROSSHAIR expects (cx, cy, radius, r, g, b, a)")
	}
	cx, ok1 := argInt(args[0])
	cy, ok2 := argInt(args[1])
	rad, ok3 := argInt(args[2])
	c, ok4 := rgbaAt(args, 3)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAW.CROSSHAIR: invalid arguments")
	}
	ix, iy := int32(cx), int32(cy)
	r32 := int32(rad)
	rl.DrawLine(ix-r32, iy, ix+r32, iy, c)
	rl.DrawLine(ix, iy-r32, ix, iy+r32, c)
	return value.Nil, nil
}

func (m *Module) drawRectGrid(rt *runtime.Runtime, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 14 {
		return value.Nil, fmt.Errorf("DRAW.RECTGRID expects (x, y, cellW, cellH, cols, rows, line r,g,b,a, fill r,g,b,a)")
	}
	x, ok1 := argInt(args[0])
	y, ok2 := argInt(args[1])
	cw, ok3 := argInt(args[2])
	ch, ok4 := argInt(args[3])
	cols, ok5 := argInt(args[4])
	rows, ok6 := argInt(args[5])
	lineC, ok7 := rgbaAt(args, 6)
	fillC, ok8 := rgbaAt(args, 10)
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 {
		return value.Nil, fmt.Errorf("DRAW.RECTGRID: invalid arguments")
	}
	if cw <= 0 || ch <= 0 || cols <= 0 || rows <= 0 {
		return value.Nil, nil
	}
	ix, iy := int32(x), int32(y)
	icw, ich := int32(cw), int32(ch)
	for row := int32(0); row < int32(rows); row++ {
		for col := int32(0); col < int32(cols); col++ {
			rx := ix + col*icw
			ry := iy + row*ich
			rl.DrawRectangle(rx, ry, icw, ich, fillC)
			rl.DrawRectangleLines(rx, ry, icw, ich, lineC)
		}
	}
	return value.Nil, nil
}
