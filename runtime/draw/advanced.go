//go:build cgo || (windows && !cgo)

package mbdraw

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerAdvancedCmds(m *Module, r runtime.Registrar) {
	r.Register("DRAW.ARC", "draw", runtime.AdaptLegacy(m.arc))
	r.Register("DRAW.DOT", "draw", runtime.AdaptLegacy(m.dot))
	r.Register("DRAW.PIXEL", "draw", runtime.AdaptLegacy(m.pixel))
	r.Register("DRAW.PIXELV", "draw", runtime.AdaptLegacy(m.pixelV))
	r.Register("DRAW.SETPIXELCOLOR", "draw", runtime.AdaptLegacy(m.pixel)) // Alias
	r.Register("DRAW.GRID2D", "draw", runtime.AdaptLegacy(m.grid2D))
	r.Register("DRAW.GETPIXELCOLOR", "draw", runtime.AdaptLegacy(m.getPixelColor))
}

func (m *Module) grid2D(args []value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("DRAW.GRID2D expects 5 arguments (spacing, r,g,b,a)")
	}
	spacing, ok1 := argInt(args[0])
	if !ok1 {
		return value.Nil, fmt.Errorf("DRAW.GRID2D: spacing must be numeric")
	}
	r, ok2 := argInt(args[1])
	g, ok3 := argInt(args[2])
	b, ok4 := argInt(args[3])
	a, ok5 := argInt(args[4])
	if !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("DRAW.GRID2D: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}

	screenWidth := rl.GetScreenWidth()
	screenHeight := rl.GetScreenHeight()

	for i := int32(0); i < int32(screenWidth); i += int32(spacing) {
		rl.DrawLine(i, 0, i, int32(screenHeight), col)
	}
	for i := int32(0); i < int32(screenHeight); i += int32(spacing) {
		rl.DrawLine(0, i, int32(screenWidth), i, col)
	}

	return value.Nil, nil
}

func (m *Module) getPixelColor(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAW.GETPIXELCOLOR expects 2 arguments (x, y)")
	}
	x, ok1 := argInt(args[0])
	y, ok2 := argInt(args[1])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("DRAW.GETPIXELCOLOR: x and y must be numeric")
	}
	img := rl.LoadImageFromScreen()
	defer rl.UnloadImage(img)
	col := rl.GetImageColor(*img, x, y)
	arr, err := heap.NewArray([]int64{4})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, float64(col.R))
	_ = arr.Set([]int64{1}, float64(col.G))
	_ = arr.Set([]int64{2}, float64(col.B))
	_ = arr.Set([]int64{3}, float64(col.A))
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) pixel(args []value.Value) (value.Value, error) {
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("DRAW.PIXEL expects 6 arguments (x, y, r,g,b,a)")
	}
	x, ok1 := argInt(args[0])
	y, ok2 := argInt(args[1])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("DRAW.PIXEL: x and y must be numeric")
	}
	r, ok3 := argInt(args[2])
	g, ok4 := argInt(args[3])
	b, ok5 := argInt(args[4])
	a, ok6 := argInt(args[5])
	if !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("DRAW.PIXEL: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawPixel(x, y, col)
	return value.Nil, nil
}

func (m *Module) pixelV(args []value.Value) (value.Value, error) {
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("DRAW.PIXELV expects 6 arguments (x, y, r,g,b,a)")
	}
	x, ok1 := argFloat(args[0])
	y, ok2 := argFloat(args[1])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("DRAW.PIXELV: x and y must be numeric")
	}
	r, ok3 := argInt(args[2])
	g, ok4 := argInt(args[3])
	b, ok5 := argInt(args[4])
	a, ok6 := argInt(args[5])
	if !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("DRAW.PIXELV: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawPixelV(rl.Vector2{X: x, Y: y}, col)
	return value.Nil, nil
}

func (m *Module) arc(args []value.Value) (value.Value, error) {
	if len(args) != 10 {
		return value.Nil, fmt.Errorf("DRAW.ARC expects 10 arguments (cx,cy, r, start, end, thick, r,g,b,a)")
	}
	cx, ok1 := argFloat(args[0])
	cy, ok2 := argFloat(args[1])
	radius, ok3 := argFloat(args[2])
	start, ok4 := argFloat(args[3])
	end, ok5 := argFloat(args[4])
	thick, ok6 := argFloat(args[5])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("DRAW.ARC: geometry arguments must be numeric")
	}
	if thick > radius {
		thick = radius
	}
	r, ok7 := argInt(args[6])
	g, ok8 := argInt(args[7])
	b, ok9 := argInt(args[8])
	a, ok10 := argInt(args[9])
	if !ok7 || !ok8 || !ok9 || !ok10 {
		return value.Nil, fmt.Errorf("DRAW.ARC: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawRing(rl.Vector2{X: cx, Y: cy}, radius-thick, radius, start, end, 32, col)
	return value.Nil, nil
}

func (m *Module) dot(args []value.Value) (value.Value, error) {
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("DRAW.DOT expects 7 arguments (x, y, size, r, g, b, a)")
	}
	x, ok1 := argFloat(args[0])
	y, ok2 := argFloat(args[1])
	size, ok3 := argFloat(args[2])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("DRAW.DOT: x, y, and size must be numeric")
	}
	r, ok4 := argInt(args[3])
	g, ok5 := argInt(args[4])
	b, ok6 := argInt(args[5])
	a, ok7 := argInt(args[6])
	if !ok4 || !ok5 || !ok6 || !ok7 {
		return value.Nil, fmt.Errorf("DRAW.DOT: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawCircleV(rl.Vector2{X: x, Y: y}, size, col)
	return value.Nil, nil
}
