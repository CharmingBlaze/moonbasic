//go:build cgo || (windows && !cgo)

package mbdraw

import (
	"fmt"

	"moonbasic/hal"
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerShapeCmds(m *Module, r runtime.Registrar) {
	r.Register("DRAW.RECTANGLE", "draw", m.rectangle)
	r.Register("DRAW.CIRCLE", "draw", m.circle)
	r.Register("DRAW.TEXT", "draw", m.drawText)
	r.Register("DRAW.RECTLINES", "draw", m.rectLines)
	r.Register("DRAW.CIRCLELINES", "draw", m.circleLines)
	r.Register("DRAW.TRIANGLE", "draw", m.triangle)
	r.Register("DRAW.POLY", "draw", m.poly)
}

func (m *Module) rectLines(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 9 {
		return value.Nil, fmt.Errorf("DRAW.RECTLINES expects 9 arguments (x,y,w,h, thick, r,g,b,a)")
	}
	x, ok1 := argInt(args[0])
	y, ok2 := argInt(args[1])
	w, ok3 := argInt(args[2])
	h, ok4 := argInt(args[3])
	thick, ok5 := argFloat(args[4])
	r, ok6 := argInt(args[5])
	g, ok7 := argInt(args[6])
	b, ok8 := argInt(args[7])
	a, ok9 := argInt(args[8])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 || !ok9 {
		return value.Nil, fmt.Errorf("DRAW.RECTLINES: numeric arguments required")
	}
	rt.Driver.Video.DrawRectangleLines(x, y, w, h, thick, hal.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
	return value.Nil, nil
}

func (m *Module) poly(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 9 {
		return value.Nil, fmt.Errorf("DRAW.POLY expects 9 arguments (cx,cy, sides, radius, rotation, r,g,b,a)")
	}
	cx, ok1 := argFloat(args[0])
	cy, ok2 := argFloat(args[1])
	sides, ok3 := argInt(args[2])
	radius, ok4 := argFloat(args[3])
	rotation, ok5 := argFloat(args[4])
	r, ok6 := argInt(args[5])
	g, ok7 := argInt(args[6])
	b, ok8 := argInt(args[7])
	a, ok9 := argInt(args[8])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 || !ok9 {
		return value.Nil, fmt.Errorf("DRAW.POLY: numeric arguments required")
	}
	rt.Driver.Video.DrawPoly(hal.V2{X: cx, Y: cy}, sides, radius, rotation, hal.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
	return value.Nil, nil
}

func (m *Module) triangle(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 10 {
		return value.Nil, fmt.Errorf("DRAW.TRIANGLE expects 10 arguments (x1,y1, x2,y2, x3,y3, r,g,b,a)")
	}
	v1x, ok1 := argFloat(args[0])
	v1y, ok2 := argFloat(args[1])
	v2x, ok3 := argFloat(args[2])
	v2y, ok4 := argFloat(args[3])
	v3x, ok5 := argFloat(args[4])
	v3y, ok6 := argFloat(args[5])
	r, ok7 := argInt(args[6])
	g, ok8 := argInt(args[7])
	b, ok9 := argInt(args[8])
	a, ok10 := argInt(args[9])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 || !ok9 || !ok10 {
		return value.Nil, fmt.Errorf("DRAW.TRIANGLE: numeric arguments required")
	}
	rt.Driver.Video.DrawTriangle(
		hal.V2{X: v1x, Y: v1y},
		hal.V2{X: v2x, Y: v2y},
		hal.V2{X: v3x, Y: v3y},
		hal.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)},
	)
	return value.Nil, nil
}

func (m *Module) circle(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
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
	rt.Driver.Video.DrawCircle(cx, cy, radius, hal.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
	return value.Nil, nil
}

func (m *Module) circleLines(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
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
	rt.Driver.Video.DrawCircleLines(cx, cy, radius, hal.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
	return value.Nil, nil
}

func (m *Module) rectangle(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 8 {
		return value.Nil, fmt.Errorf("DRAW.RECTANGLE expects 8 arguments (x,y,w,h, r,g,b,a)")
	}
	x, ok1 := argInt(args[0])
	y, ok2 := argInt(args[1])
	w, ok3 := argInt(args[2])
	h, ok4 := argInt(args[3])
	r, ok5 := argInt(args[4])
	g, ok6 := argInt(args[5])
	b, ok7 := argInt(args[6])
	a, ok8 := argInt(args[7])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 {
		return value.Nil, fmt.Errorf("DRAW.RECTANGLE: arguments must be numeric")
	}
	rt.Driver.Video.DrawRectangle(x, y, w, h, hal.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
	return value.Nil, nil
}

func (m *Module) drawText(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 8 {
		return value.Nil, fmt.Errorf("DRAW.TEXT expects 8 arguments (text, x, y, size, r, g, b, a)")
	}
	text, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argInt(args[1])
	y, ok2 := argInt(args[2])
	size, ok3 := argInt(args[3])
	r, ok4 := argInt(args[4])
	g, ok5 := argInt(args[5])
	b, ok6 := argInt(args[6])
	a, ok7 := argInt(args[7])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 {
		return value.Nil, fmt.Errorf("DRAW.TEXT: numeric arguments required")
	}
	rt.Driver.Video.DrawText(text, x, y, size, hal.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
	return value.Nil, nil
}

// Remaining shapes like gradients and ellipses can be added to hal.VideoDevice as needed.
