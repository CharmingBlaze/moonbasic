//go:build cgo

package mbdraw

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerCircleExtraCmds(m *Module, r runtime.Registrar) {
	r.Register("DRAW.CIRCLESECTOR", "draw", runtime.AdaptLegacy(m.circleSector))
	r.Register("DRAW.CIRCLEGRADIENT", "draw", runtime.AdaptLegacy(m.circleGradient))
}

func (m *Module) circleSector(args []value.Value) (value.Value, error) {
	if len(args) != 10 {
		return value.Nil, fmt.Errorf("DRAW.CIRCLESECTOR expects 10 arguments (cx, cy, radius, startAngle, endAngle, segments, r, g, b, a)")
	}
	cx, ok0 := argInt(args[0])
	cy, ok1 := argInt(args[1])
	rad, ok2 := argFloat(args[2])
	start, ok3 := argFloat(args[3])
	end, ok4 := argFloat(args[4])
	segs, ok5 := argInt(args[5])
	if !ok0 || !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("DRAW.CIRCLESECTOR: geometry must be numeric")
	}
	cr, ok6 := argInt(args[6])
	cg, ok7 := argInt(args[7])
	cb, ok8 := argInt(args[8])
	ca, ok9 := argInt(args[9])
	if !ok6 || !ok7 || !ok8 || !ok9 {
		return value.Nil, fmt.Errorf("DRAW.CIRCLESECTOR: color must be numeric")
	}
	col := color.RGBA{R: uint8(cr), G: uint8(cg), B: uint8(cb), A: uint8(ca)}
	rl.DrawCircleSector(rl.Vector2{X: float32(cx), Y: float32(cy)}, rad, start, end, int32(segs), col)
	return value.Nil, nil
}

func (m *Module) circleGradient(args []value.Value) (value.Value, error) {
	if len(args) != 11 {
		return value.Nil, fmt.Errorf("DRAW.CIRCLEGRADIENT expects 11 arguments (cx, cy, radius, inner r,g,b,a, outer r,g,b,a)")
	}
	cx, ok0 := argInt(args[0])
	cy, ok1 := argInt(args[1])
	rad, ok2 := argFloat(args[2])
	ir, ok3 := argInt(args[3])
	ig, ok4 := argInt(args[4])
	ib, ok5 := argInt(args[5])
	ia, ok6 := argInt(args[6])
	or_, ok7 := argInt(args[7])
	og, ok8 := argInt(args[8])
	ob, ok9 := argInt(args[9])
	oa, ok10 := argInt(args[10])
	if !ok0 || !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 || !ok9 || !ok10 {
		return value.Nil, fmt.Errorf("DRAW.CIRCLEGRADIENT: arguments must be numeric")
	}
	inner := color.RGBA{R: uint8(ir), G: uint8(ig), B: uint8(ib), A: uint8(ia)}
	outer := color.RGBA{R: uint8(or_), G: uint8(og), B: uint8(ob), A: uint8(oa)}
	rl.DrawCircleGradient(int32(cx), int32(cy), rad, inner, outer)
	return value.Nil, nil
}
