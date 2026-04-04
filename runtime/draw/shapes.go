//go:build cgo

package mbdraw

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/convert"
	"moonbasic/vm/value"
)

func registerShapeCmds(m *Module, r runtime.Registrar) {
	r.Register("DRAW.RECTANGLE", "draw", runtime.AdaptLegacy(m.rectangle))
	r.Register("DRAW.RECTANGLE_ROUNDED", "draw", runtime.AdaptLegacy(m.rectangleRounded))
	r.Register("DRAW.CIRCLE", "draw", runtime.AdaptLegacy(m.circle))
	r.Register("DRAW.CIRCLELINES", "draw", runtime.AdaptLegacy(m.circleLines))
	r.Register("DRAW.ELLIPSE", "draw", runtime.AdaptLegacy(m.ellipse))
	r.Register("DRAW.ELLIPSELINES", "draw", runtime.AdaptLegacy(m.ellipseLines))
	r.Register("DRAW.RING", "draw", runtime.AdaptLegacy(m.ring))
	r.Register("DRAW.RINGLINES", "draw", runtime.AdaptLegacy(m.ringLines))
	r.Register("DRAW.TRIANGLE", "draw", runtime.AdaptLegacy(m.triangle))
	r.Register("DRAW.TRIANGLELINES", "draw", runtime.AdaptLegacy(m.triangleLines))
	r.Register("DRAW.POLY", "draw", runtime.AdaptLegacy(m.poly))
	r.Register("DRAW.POLYLINES", "draw", runtime.AdaptLegacy(m.polyLines))
	r.Register("DRAW.RECTLINES", "draw", runtime.AdaptLegacy(m.rectLines))
	r.Register("DRAW.RECTPRO", "draw", runtime.AdaptLegacy(m.rectPro))
	r.Register("DRAW.RECTGRADV", "draw", runtime.AdaptLegacy(m.rectGradV))
	r.Register("DRAW.RECTGRADH", "draw", runtime.AdaptLegacy(m.rectGradH))
	r.Register("DRAW.RECTGRAD", "draw", runtime.AdaptLegacy(m.rectGrad))
}

func (m *Module) rectGradV(args []value.Value) (value.Value, error) {
	if len(args) != 12 {
		return value.Nil, fmt.Errorf("DRAW.RECTGRADV expects 12 arguments (x,y,w,h, r1,g1,b1,a1, r2,g2,b2,a2)")
	}
	x, ok1 := argInt(args[0])
	y, ok2 := argInt(args[1])
	w, ok3 := argInt(args[2])
	h, ok4 := argInt(args[3])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAW.RECTGRADV: geometry arguments must be numeric")
	}
	r1, ok5 := argInt(args[4])
	g1, ok6 := argInt(args[5])
	b1, ok7 := argInt(args[6])
	a1, ok8 := argInt(args[7])
	r2, ok9 := argInt(args[8])
	g2, ok10 := argInt(args[9])
	b2, ok11 := argInt(args[10])
	a2, ok12 := argInt(args[11])
	if !ok5 || !ok6 || !ok7 || !ok8 || !ok9 || !ok10 || !ok11 || !ok12 {
		return value.Nil, fmt.Errorf("DRAW.RECTGRADV: color components must be numeric")
	}
	col1 := color.RGBA{R: uint8(r1), G: uint8(g1), B: uint8(b1), A: uint8(a1)}
	col2 := color.RGBA{R: uint8(r2), G: uint8(g2), B: uint8(b2), A: uint8(a2)}
	rl.DrawRectangleGradientV(x, y, w, h, col1, col2)
	return value.Nil, nil
}

func (m *Module) rectGradH(args []value.Value) (value.Value, error) {
	if len(args) != 12 {
		return value.Nil, fmt.Errorf("DRAW.RECTGRADH expects 12 arguments (x,y,w,h, r1,g1,b1,a1, r2,g2,b2,a2)")
	}
	x, ok1 := argInt(args[0])
	y, ok2 := argInt(args[1])
	w, ok3 := argInt(args[2])
	h, ok4 := argInt(args[3])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAW.RECTGRADH: geometry arguments must be numeric")
	}
	r1, ok5 := argInt(args[4])
	g1, ok6 := argInt(args[5])
	b1, ok7 := argInt(args[6])
	a1, ok8 := argInt(args[7])
	r2, ok9 := argInt(args[8])
	g2, ok10 := argInt(args[9])
	b2, ok11 := argInt(args[10])
	a2, ok12 := argInt(args[11])
	if !ok5 || !ok6 || !ok7 || !ok8 || !ok9 || !ok10 || !ok11 || !ok12 {
		return value.Nil, fmt.Errorf("DRAW.RECTGRADH: color components must be numeric")
	}
	col1 := color.RGBA{R: uint8(r1), G: uint8(g1), B: uint8(b1), A: uint8(a1)}
	col2 := color.RGBA{R: uint8(r2), G: uint8(g2), B: uint8(b2), A: uint8(a2)}
	rl.DrawRectangleGradientH(x, y, w, h, col1, col2)
	return value.Nil, nil
}

func (m *Module) rectGrad(args []value.Value) (value.Value, error) {
	if len(args) != 20 {
		return value.Nil, fmt.Errorf("DRAW.RECTGRAD expects 20 arguments (x,y,w,h, r1,g1,b1,a1, r2,g2,b2,a2, r3,g3,b3,a3, r4,g4,b4,a4)")
	}
	x, ok1 := argFloat(args[0])
	y, ok2 := argFloat(args[1])
	w, ok3 := argFloat(args[2])
	h, ok4 := argFloat(args[3])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAW.RECTGRAD: geometry arguments must be numeric")
	}
	r1, ok5 := argInt(args[4])
	g1, ok6 := argInt(args[5])
	b1, ok7 := argInt(args[6])
	a1, ok8 := argInt(args[7])
	r2, ok9 := argInt(args[8])
	g2, ok10 := argInt(args[9])
	b2, ok11 := argInt(args[10])
	a2, ok12 := argInt(args[11])
	r3, ok13 := argInt(args[12])
	g3, ok14 := argInt(args[13])
	b3, ok15 := argInt(args[14])
	a3, ok16 := argInt(args[15])
	r4, ok17 := argInt(args[16])
	g4, ok18 := argInt(args[17])
	b4, ok19 := argInt(args[18])
	a4, ok20 := argInt(args[19])
	if !ok5 || !ok6 || !ok7 || !ok8 || !ok9 || !ok10 || !ok11 || !ok12 || !ok13 || !ok14 || !ok15 || !ok16 || !ok17 || !ok18 || !ok19 || !ok20 {
		return value.Nil, fmt.Errorf("DRAW.RECTGRAD: color components must be numeric")
	}
	col1 := color.RGBA{R: uint8(r1), G: uint8(g1), B: uint8(b1), A: uint8(a1)}
	col2 := color.RGBA{R: uint8(r2), G: uint8(g2), B: uint8(b2), A: uint8(a2)}
	col3 := color.RGBA{R: uint8(r3), G: uint8(g3), B: uint8(b3), A: uint8(a3)}
	col4 := color.RGBA{R: uint8(r4), G: uint8(g4), B: uint8(b4), A: uint8(a4)}
	rl.DrawRectangleGradientEx(rl.Rectangle{X: x, Y: y, Width: w, Height: h}, col1, col2, col3, col4)
	return value.Nil, nil
}

func (m *Module) rectLines(args []value.Value) (value.Value, error) {
	if len(args) != 9 {
		return value.Nil, fmt.Errorf("DRAW.RECTLINES expects 9 arguments (x,y,w,h, thick, r,g,b,a)")
	}
	x, ok1 := argInt(args[0])
	y, ok2 := argInt(args[1])
	w, ok3 := argInt(args[2])
	h, ok4 := argInt(args[3])
	thick, ok5 := argFloat(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("DRAW.RECTLINES: geometry arguments must be numeric")
	}
	r, ok6 := argInt(args[5])
	g, ok7 := argInt(args[6])
	b, ok8 := argInt(args[7])
	a, ok9 := argInt(args[8])
	if !ok6 || !ok7 || !ok8 || !ok9 {
		return value.Nil, fmt.Errorf("DRAW.RECTLINES: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawRectangleLinesEx(rl.Rectangle{X: float32(x), Y: float32(y), Width: float32(w), Height: float32(h)}, thick, col)
	return value.Nil, nil
}

func (m *Module) rectPro(args []value.Value) (value.Value, error) {
	if len(args) != 11 {
		return value.Nil, fmt.Errorf("DRAW.RECTPRO expects 11 arguments (x,y,w,h, ox,oy, rot, r,g,b,a)")
	}
	x, ok1 := argFloat(args[0])
	y, ok2 := argFloat(args[1])
	w, ok3 := argFloat(args[2])
	h, ok4 := argFloat(args[3])
	ox, ok5 := argFloat(args[4])
	oy, ok6 := argFloat(args[5])
	rot, ok7 := argFloat(args[6])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 {
		return value.Nil, fmt.Errorf("DRAW.RECTPRO: geometry arguments must be numeric")
	}
	r, ok8 := argInt(args[7])
	g, ok9 := argInt(args[8])
	b, ok10 := argInt(args[9])
	a, ok11 := argInt(args[10])
	if !ok8 || !ok9 || !ok10 || !ok11 {
		return value.Nil, fmt.Errorf("DRAW.RECTPRO: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawRectanglePro(rl.Rectangle{X: x, Y: y, Width: w, Height: h}, rl.Vector2{X: ox, Y: oy}, rot, col)
	return value.Nil, nil
}

func (m *Module) poly(args []value.Value) (value.Value, error) {
	if len(args) != 9 {
		return value.Nil, fmt.Errorf("DRAW.POLY expects 9 arguments (cx,cy, sides, radius, rotation, r,g,b,a)")
	}
	cx, ok1 := argFloat(args[0])
	cy, ok2 := argFloat(args[1])
	sides, ok3 := argInt(args[2])
	radius, ok4 := argFloat(args[3])
	rotation, ok5 := argFloat(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("DRAW.POLY: geometry arguments must be numeric")
	}
	r, ok6 := argInt(args[5])
	g, ok7 := argInt(args[6])
	b, ok8 := argInt(args[7])
	a, ok9 := argInt(args[8])
	if !ok6 || !ok7 || !ok8 || !ok9 {
		return value.Nil, fmt.Errorf("DRAW.POLY: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawPoly(rl.Vector2{X: cx, Y: cy}, int32(sides), radius, rotation, col)
	return value.Nil, nil
}

func (m *Module) polyLines(args []value.Value) (value.Value, error) {
	if len(args) != 10 {
		return value.Nil, fmt.Errorf("DRAW.POLYLINES expects 10 arguments (cx,cy, sides, radius, rotation, thick, r,g,b,a)")
	}
	cx, ok1 := argFloat(args[0])
	cy, ok2 := argFloat(args[1])
	sides, ok3 := argInt(args[2])
	radius, ok4 := argFloat(args[3])
	rotation, ok5 := argFloat(args[4])
	thick, ok6 := argFloat(args[5])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("DRAW.POLYLINES: geometry arguments must be numeric")
	}
	r, ok7 := argInt(args[6])
	g, ok8 := argInt(args[7])
	b, ok9 := argInt(args[8])
	a, ok10 := argInt(args[9])
	if !ok7 || !ok8 || !ok9 || !ok10 {
		return value.Nil, fmt.Errorf("DRAW.POLYLINES: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawPolyLinesEx(rl.Vector2{X: cx, Y: cy}, int32(sides), radius, rotation, thick, col)
	return value.Nil, nil
}

func (m *Module) triangle(args []value.Value) (value.Value, error) {
	if len(args) != 10 {
		return value.Nil, fmt.Errorf("DRAW.TRIANGLE expects 10 arguments (x1,y1, x2,y2, x3,y3, r,g,b,a)")
	}
	v1x, ok1 := argFloat(args[0])
	v1y, ok2 := argFloat(args[1])
	v2x, ok3 := argFloat(args[2])
	v2y, ok4 := argFloat(args[3])
	v3x, ok5 := argFloat(args[4])
	v3y, ok6 := argFloat(args[5])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("DRAW.TRIANGLE: vertex coordinates must be numeric")
	}
	r, ok7 := argInt(args[6])
	g, ok8 := argInt(args[7])
	b, ok9 := argInt(args[8])
	a, ok10 := argInt(args[9])
	if !ok7 || !ok8 || !ok9 || !ok10 {
		return value.Nil, fmt.Errorf("DRAW.TRIANGLE: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawTriangle(rl.Vector2{X: v1x, Y: v1y}, rl.Vector2{X: v2x, Y: v2y}, rl.Vector2{X: v3x, Y: v3y}, col)
	return value.Nil, nil
}

func (m *Module) triangleLines(args []value.Value) (value.Value, error) {
	if len(args) != 10 {
		return value.Nil, fmt.Errorf("DRAW.TRIANGLELINES expects 10 arguments (x1,y1, x2,y2, x3,y3, r,g,b,a)")
	}
	v1x, ok1 := argFloat(args[0])
	v1y, ok2 := argFloat(args[1])
	v2x, ok3 := argFloat(args[2])
	v2y, ok4 := argFloat(args[3])
	v3x, ok5 := argFloat(args[4])
	v3y, ok6 := argFloat(args[5])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("DRAW.TRIANGLELINES: vertex coordinates must be numeric")
	}
	r, ok7 := argInt(args[6])
	g, ok8 := argInt(args[7])
	b, ok9 := argInt(args[8])
	a, ok10 := argInt(args[9])
	if !ok7 || !ok8 || !ok9 || !ok10 {
		return value.Nil, fmt.Errorf("DRAW.TRIANGLELINES: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawTriangleLines(rl.Vector2{X: v1x, Y: v1y}, rl.Vector2{X: v2x, Y: v2y}, rl.Vector2{X: v3x, Y: v3y}, col)
	return value.Nil, nil
}

func (m *Module) ring(args []value.Value) (value.Value, error) {
	if len(args) != 11 {
		return value.Nil, fmt.Errorf("DRAW.RING expects 11 arguments (cx,cy, innerR, outerR, start,end, segs, r,g,b,a)")
	}
	cx, ok1 := argFloat(args[0])
	cy, ok2 := argFloat(args[1])
	innerR, ok3 := argFloat(args[2])
	outerR, ok4 := argFloat(args[3])
	start, ok5 := argFloat(args[4])
	end, ok6 := argFloat(args[5])
	segs, ok7 := argInt(args[6])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 {
		return value.Nil, fmt.Errorf("DRAW.RING: numeric arguments required for geometry")
	}
	r, ok8 := argInt(args[7])
	g, ok9 := argInt(args[8])
	b, ok10 := argInt(args[9])
	a, ok11 := argInt(args[10])
	if !ok8 || !ok9 || !ok10 || !ok11 {
		return value.Nil, fmt.Errorf("DRAW.RING: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawRing(rl.Vector2{X: cx, Y: cy}, innerR, outerR, start, end, int32(segs), col)
	return value.Nil, nil
}

func (m *Module) ringLines(args []value.Value) (value.Value, error) {
	if len(args) != 11 {
		return value.Nil, fmt.Errorf("DRAW.RINGLINES expects 11 arguments (cx,cy, innerR, outerR, start,end, segs, r,g,b,a)")
	}
	cx, ok1 := argFloat(args[0])
	cy, ok2 := argFloat(args[1])
	innerR, ok3 := argFloat(args[2])
	outerR, ok4 := argFloat(args[3])
	start, ok5 := argFloat(args[4])
	end, ok6 := argFloat(args[5])
	segs, ok7 := argInt(args[6])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 {
		return value.Nil, fmt.Errorf("DRAW.RINGLINES: numeric arguments required for geometry")
	}
	r, ok8 := argInt(args[7])
	g, ok9 := argInt(args[8])
	b, ok10 := argInt(args[9])
	a, ok11 := argInt(args[10])
	if !ok8 || !ok9 || !ok10 || !ok11 {
		return value.Nil, fmt.Errorf("DRAW.RINGLINES: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawRingLines(rl.Vector2{X: cx, Y: cy}, innerR, outerR, start, end, int32(segs), col)
	return value.Nil, nil
}

func (m *Module) ellipse(args []value.Value) (value.Value, error) {
	if len(args) != 8 {
		return value.Nil, fmt.Errorf("DRAW.ELLIPSE expects 8 arguments (cx, cy, rx, ry, r, g, b, a)")
	}
	cx, ok1 := argInt(args[0])
	cy, ok2 := argInt(args[1])
	rx, ok3 := argFloat(args[2])
	ry, ok4 := argFloat(args[3])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAW.ELLIPSE: cx, cy, rx, ry must be numeric")
	}
	r, ok5 := argInt(args[4])
	g, ok6 := argInt(args[5])
	b, ok7 := argInt(args[6])
	a, ok8 := argInt(args[7])
	if !ok5 || !ok6 || !ok7 || !ok8 {
		return value.Nil, fmt.Errorf("DRAW.ELLIPSE: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawEllipse(cx, cy, rx, ry, col)
	return value.Nil, nil
}

func (m *Module) ellipseLines(args []value.Value) (value.Value, error) {
	if len(args) != 8 {
		return value.Nil, fmt.Errorf("DRAW.ELLIPSELINES expects 8 arguments (cx, cy, rx, ry, r, g, b, a)")
	}
	cx, ok1 := argInt(args[0])
	cy, ok2 := argInt(args[1])
	rx, ok3 := argFloat(args[2])
	ry, ok4 := argFloat(args[3])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAW.ELLIPSELINES: cx, cy, rx, ry must be numeric")
	}
	r, ok5 := argInt(args[4])
	g, ok6 := argInt(args[5])
	b, ok7 := argInt(args[6])
	a, ok8 := argInt(args[7])
	if !ok5 || !ok6 || !ok7 || !ok8 {
		return value.Nil, fmt.Errorf("DRAW.ELLIPSELINES: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawEllipseLines(cx, cy, rx, ry, col)
	return value.Nil, nil
}

func (m *Module) circle(args []value.Value) (value.Value, error) {
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("DRAW.CIRCLE expects 7 arguments (cx, cy, radius, r, g, b, a)")
	}
	cx, ok1 := argInt(args[0])
	cy, ok2 := argInt(args[1])
	radius, ok3 := argFloat(args[2])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("DRAW.CIRCLE: cx, cy, radius must be numeric")
	}
	r, ok4 := argInt(args[3])
	g, ok5 := argInt(args[4])
	b, ok6 := argInt(args[5])
	a, ok7 := argInt(args[6])
	if !ok4 || !ok5 || !ok6 || !ok7 {
		return value.Nil, fmt.Errorf("DRAW.CIRCLE: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawCircle(cx, cy, radius, col)
	return value.Nil, nil
}

func (m *Module) circleLines(args []value.Value) (value.Value, error) {
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("DRAW.CIRCLELINES expects 7 arguments (cx, cy, radius, r, g, b, a)")
	}
	cx, ok1 := argInt(args[0])
	cy, ok2 := argInt(args[1])
	radius, ok3 := argFloat(args[2])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("DRAW.CIRCLELINES: cx, cy, radius must be numeric")
	}
	r, ok4 := argInt(args[3])
	g, ok5 := argInt(args[4])
	b, ok6 := argInt(args[5])
	a, ok7 := argInt(args[6])
	if !ok4 || !ok5 || !ok6 || !ok7 {
		return value.Nil, fmt.Errorf("DRAW.CIRCLELINES: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawCircleLines(cx, cy, radius, col)
	return value.Nil, nil
}

func (m *Module) rectangle(args []value.Value) (value.Value, error) {
	if len(args) != 8 {
		return value.Nil, fmt.Errorf("DRAW.RECTANGLE expects 8 arguments (x,y,w,h, r,g,b,a)")
	}
	var xywh [4]int32
	var rgb [4]int32
	for i := 0; i < 4; i++ {
		v, ok := argInt(args[i])
		if !ok {
			return value.Nil, fmt.Errorf("DRAW.RECTANGLE: non-numeric argument %d", i+1)
		}
		xywh[i] = v
	}
	for i := 0; i < 4; i++ {
		v, ok := argInt(args[4+i])
		if !ok {
			return value.Nil, fmt.Errorf("DRAW.RECTANGLE: non-numeric color argument %d", i+1)
		}
		rgb[i] = v
	}
	c := convert.NewColor4(rgb[0], rgb[1], rgb[2], rgb[3])
	col := color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
	rl.DrawRectangle(xywh[0], xywh[1], xywh[2], xywh[3], col)
	return value.Nil, nil
}

func (m *Module) rectangleRounded(args []value.Value) (value.Value, error) {
	if len(args) != 9 {
		return value.Nil, fmt.Errorf("DRAW.RECTANGLE_ROUNDED expects 9 arguments (x,y,w,h, radius, r,g,b,a)")
	}
	var xywh [4]int32
	for i := 0; i < 4; i++ {
		v, ok := argInt(args[i])
		if !ok {
			return value.Nil, fmt.Errorf("DRAW.RECTANGLE_ROUNDED: non-numeric argument %d", i+1)
		}
		xywh[i] = v
	}
	rad, ok := argInt(args[4])
	if !ok {
		return value.Nil, fmt.Errorf("DRAW.RECTANGLE_ROUNDED: radius must be numeric")
	}
	var rgb [4]int32
	for i := 0; i < 4; i++ {
		v, ok := argInt(args[5+i])
		if !ok {
			return value.Nil, fmt.Errorf("DRAW.RECTANGLE_ROUNDED: non-numeric color argument %d", i+1)
		}
		rgb[i] = v
	}
	c := convert.NewColor4(rgb[0], rgb[1], rgb[2], rgb[3])
	col := color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
	rl.DrawRectangleRounded(
		rl.Rectangle{X: float32(xywh[0]), Y: float32(xywh[1]), Width: float32(xywh[2]), Height: float32(xywh[3])},
		float32(rad),
		8,
		col,
	)
	return value.Nil, nil
}
