//go:build cgo

package mbdraw

import (
	"fmt"
	"image/color"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

var vecPool = sync.Pool{
	New: func() interface{} {
		s := make([]rl.Vector2, 0, 64)
		return s
	},
}

func registerLineCmds(m *Module, r runtime.Registrar) {
	r.Register("DRAW.LINE", "draw", runtime.AdaptLegacy(m.line))
	r.Register("DRAW.LINEEX", "draw", runtime.AdaptLegacy(m.lineEx))
	r.Register("DRAW.LINEBEZIER", "draw", runtime.AdaptLegacy(m.lineBezier))
	r.Register("DRAW.LINEBEZIERQUAD", "draw", runtime.AdaptLegacy(m.lineBezierQuad))
	r.Register("DRAW.LINEBEZIERCUBIC", "draw", runtime.AdaptLegacy(m.lineBezierCubic))
	r.Register("DRAW.SPLINELINEAR", "draw", runtime.AdaptLegacy(m.splineLinear))
	r.Register("DRAW.SPLINEBASIS", "draw", runtime.AdaptLegacy(m.splineBasis))
	r.Register("DRAW.SPLINECATMULLROM", "draw", runtime.AdaptLegacy(m.splineCatmullRom))
	r.Register("DRAW.SPLINEBEZIERQUAD", "draw", runtime.AdaptLegacy(m.splineBezierQuad))
	r.Register("DRAW.SPLINEBEZIERCUBIC", "draw", runtime.AdaptLegacy(m.splineBezierCubic))
}

func (m *Module) splineBezierQuad(args []value.Value) (value.Value, error) {
	return m.drawSpline(args, "DRAW.SPLINEBEZIERQUAD", rl.DrawSplineBezierQuadratic)
}

func (m *Module) splineBezierCubic(args []value.Value) (value.Value, error) {
	return m.drawSpline(args, "DRAW.SPLINEBEZIERCUBIC", rl.DrawSplineBezierCubic)
}

func (m *Module) splineLinear(args []value.Value) (value.Value, error) {
	return m.drawSpline(args, "DRAW.SPLINELINEAR", rl.DrawSplineLinear)
}

func (m *Module) splineBasis(args []value.Value) (value.Value, error) {
	return m.drawSpline(args, "DRAW.SPLINEBASIS", rl.DrawSplineBasis)
}

func (m *Module) splineCatmullRom(args []value.Value) (value.Value, error) {
	return m.drawSpline(args, "DRAW.SPLINECATMULLROM", rl.DrawSplineCatmullRom)
}

func (m *Module) drawSpline(args []value.Value, name string, drawFn func([]rl.Vector2, float32, color.RGBA)) (value.Value, error) {
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("%s expects 6 arguments (points_array, thick, r,g,b,a)", name)
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("%s: first argument must be an array handle", name)
	}
	arrH := heap.Handle(args[0].IVal)
	thick, ok1 := argFloat(args[1])
	if !ok1 {
		return value.Nil, fmt.Errorf("%s: thickness must be numeric", name)
	}
	r, ok2 := argInt(args[2])
	g, ok3 := argInt(args[3])
	b, ok4 := argInt(args[4])
	a, ok5 := argInt(args[5])
	if !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("%s: color components must be numeric", name)
	}

	n := m.h.ArrayFlatLen(arrH)
	if n < 0 || n%2 != 0 {
		return value.Nil, fmt.Errorf("%s: points array must have even length (x,y pairs)", name)
	}
	count := n / 2
	pts := vecPool.Get().([]rl.Vector2)
	pts = pts[:0]
	defer vecPool.Put(pts)

	for i := 0; i < count; i++ {
		x, okx := m.h.ArrayGetFloat(arrH, int64(i*2))
		y, oky := m.h.ArrayGetFloat(arrH, int64(i*2+1))
		if !okx || !oky {
			return value.Nil, fmt.Errorf("%s: invalid point data", name)
		}
		pts = append(pts, rl.NewVector2(float32(x), float32(y)))
	}

	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	drawFn(pts, thick, col)
	return value.Nil, nil
}

func (m *Module) lineBezier(args []value.Value) (value.Value, error) {
	if len(args) != 9 {
		return value.Nil, fmt.Errorf("DRAW.LINEBEZIER expects 9 arguments (x1,y1, x2,y2, thick, r,g,b,a)")
	}
	x1, ok1 := argFloat(args[0])
	y1, ok2 := argFloat(args[1])
	x2, ok3 := argFloat(args[2])
	y2, ok4 := argFloat(args[3])
	thick, ok5 := argFloat(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("DRAW.LINEBEZIER: coordinates and thickness must be numeric")
	}
	r, ok6 := argInt(args[5])
	g, ok7 := argInt(args[6])
	b, ok8 := argInt(args[7])
	a, ok9 := argInt(args[8])
	if !ok6 || !ok7 || !ok8 || !ok9 {
		return value.Nil, fmt.Errorf("DRAW.LINEBEZIER: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawLineBezier(rl.Vector2{X: x1, Y: y1}, rl.Vector2{X: x2, Y: y2}, thick, col)
	return value.Nil, nil
}

func (m *Module) lineBezierQuad(args []value.Value) (value.Value, error) {
	if len(args) != 11 {
		return value.Nil, fmt.Errorf("DRAW.LINEBEZIERQUAD expects 11 arguments (x1,y1, cx,cy, x2,y2, thick, r,g,b,a)")
	}
	x1, ok1 := argFloat(args[0])
	y1, ok2 := argFloat(args[1])
	cx, ok3 := argFloat(args[2])
	cy, ok4 := argFloat(args[3])
	x2, ok5 := argFloat(args[4])
	y2, ok6 := argFloat(args[5])
	thick, ok7 := argFloat(args[6])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 {
		return value.Nil, fmt.Errorf("DRAW.LINEBEZIERQUAD: coordinates and thickness must be numeric")
	}
	r, ok8 := argInt(args[7])
	g, ok9 := argInt(args[8])
	b, ok10 := argInt(args[9])
	a, ok11 := argInt(args[10])
	if !ok8 || !ok9 || !ok10 || !ok11 {
		return value.Nil, fmt.Errorf("DRAW.LINEBEZIERQUAD: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawSplineSegmentBezierQuadratic(rl.Vector2{X: x1, Y: y1}, rl.Vector2{X: cx, Y: cy}, rl.Vector2{X: x2, Y: y2}, thick, col)
	return value.Nil, nil
}

func (m *Module) lineBezierCubic(args []value.Value) (value.Value, error) {
	if len(args) != 13 {
		return value.Nil, fmt.Errorf("DRAW.LINEBEZIERCUBIC expects 13 arguments (x1,y1, c1x,c1y, c2x,c2y, x2,y2, thick, r,g,b,a)")
	}
	x1, ok1 := argFloat(args[0])
	y1, ok2 := argFloat(args[1])
	c1x, ok3 := argFloat(args[2])
	c1y, ok4 := argFloat(args[3])
	c2x, ok5 := argFloat(args[4])
	c2y, ok6 := argFloat(args[5])
	x2, ok7 := argFloat(args[6])
	y2, ok8 := argFloat(args[7])
	thick, ok9 := argFloat(args[8])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 || !ok9 {
		return value.Nil, fmt.Errorf("DRAW.LINEBEZIERCUBIC: coordinates and thickness must be numeric")
	}
	r, ok10 := argInt(args[9])
	g, ok11 := argInt(args[10])
	b, ok12 := argInt(args[11])
	a, ok13 := argInt(args[12])
	if !ok10 || !ok11 || !ok12 || !ok13 {
		return value.Nil, fmt.Errorf("DRAW.LINEBEZIERCUBIC: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawSplineSegmentBezierCubic(rl.Vector2{X: x1, Y: y1}, rl.Vector2{X: c1x, Y: c1y}, rl.Vector2{X: c2x, Y: c2y}, rl.Vector2{X: x2, Y: y2}, thick, col)
	return value.Nil, nil
}

func (m *Module) line(args []value.Value) (value.Value, error) {
	if len(args) != 8 {
		return value.Nil, fmt.Errorf("DRAW.LINE expects 8 arguments (x1,y1, x2,y2, r,g,b,a)")
	}
	x1, ok1 := argInt(args[0])
	y1, ok2 := argInt(args[1])
	x2, ok3 := argInt(args[2])
	y2, ok4 := argInt(args[3])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAW.LINE: coordinates must be numeric")
	}
	r, ok5 := argInt(args[4])
	g, ok6 := argInt(args[5])
	b, ok7 := argInt(args[6])
	a, ok8 := argInt(args[7])
	if !ok5 || !ok6 || !ok7 || !ok8 {
		return value.Nil, fmt.Errorf("DRAW.LINE: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawLine(x1, y1, x2, y2, col)
	return value.Nil, nil
}

func (m *Module) lineEx(args []value.Value) (value.Value, error) {
	if len(args) != 9 {
		return value.Nil, fmt.Errorf("DRAW.LINEEX expects 9 arguments (x1,y1, x2,y2, thick, r,g,b,a)")
	}
	x1, ok1 := argFloat(args[0])
	y1, ok2 := argFloat(args[1])
	x2, ok3 := argFloat(args[2])
	y2, ok4 := argFloat(args[3])
	thick, ok5 := argFloat(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("DRAW.LINEEX: coordinates and thickness must be numeric")
	}
	r, ok6 := argInt(args[5])
	g, ok7 := argInt(args[6])
	b, ok8 := argInt(args[7])
	a, ok9 := argInt(args[8])
	if !ok6 || !ok7 || !ok8 || !ok9 {
		return value.Nil, fmt.Errorf("DRAW.LINEEX: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawLineEx(rl.Vector2{X: x1, Y: y1}, rl.Vector2{X: x2, Y: y2}, thick, col)
	return value.Nil, nil
}
