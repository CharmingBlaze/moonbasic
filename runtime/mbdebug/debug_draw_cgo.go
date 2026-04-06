//go:build cgo || (windows && !cgo)

package mbdebug

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerDebugDraw3D(r runtime.Registrar) {
	r.Register("DEBUG.DRAWLINE", "debug", runtime.AdaptLegacy(m.debugDrawLine))
	r.Register("DEBUG.DRAWBOX", "debug", runtime.AdaptLegacy(m.debugDrawBox))
}

func argF(v value.Value) (float32, bool) {
	if f, ok := v.ToFloat(); ok {
		return float32(f), true
	}
	if i, ok := v.ToInt(); ok {
		return float32(i), true
	}
	return 0, false
}

func (m *Module) debugDrawLine(args []value.Value) (value.Value, error) {
	if len(args) != 9 {
		return value.Nil, fmt.Errorf("DEBUG.DRAWLINE expects (x1#, y1#, z1#, x2#, y2#, z2#, r, g, b)")
	}
	x1, ok1 := argF(args[0])
	y1, ok2 := argF(args[1])
	z1, ok3 := argF(args[2])
	x2, ok4 := argF(args[3])
	y2, ok5 := argF(args[4])
	z2, ok6 := argF(args[5])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("DEBUG.DRAWLINE: positions must be numeric")
	}
	ri, ok7 := args[6].ToInt()
	gi, ok8 := args[7].ToInt()
	bi, ok9 := args[8].ToInt()
	if !ok7 || !ok8 || !ok9 {
		return value.Nil, fmt.Errorf("DEBUG.DRAWLINE: r,g,b must be integers 0–255")
	}
	col := rl.Color{R: uint8(ri), G: uint8(gi), B: uint8(bi), A: 255}
	rl.DrawLine3D(rl.Vector3{X: x1, Y: y1, Z: z1}, rl.Vector3{X: x2, Y: y2, Z: z2}, col)
	return value.Nil, nil
}

func (m *Module) debugDrawBox(args []value.Value) (value.Value, error) {
	if len(args) != 9 {
		return value.Nil, fmt.Errorf("DEBUG.DRAWBOX expects (cx#, cy#, cz#, w#, h#, d#, r, g, b)")
	}
	cx, ok1 := argF(args[0])
	cy, ok2 := argF(args[1])
	cz, ok3 := argF(args[2])
	w, ok4 := argF(args[3])
	h, ok5 := argF(args[4])
	d, ok6 := argF(args[5])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("DEBUG.DRAWBOX: numeric box required")
	}
	ri, ok7 := args[6].ToInt()
	gi, ok8 := args[7].ToInt()
	bi, ok9 := args[8].ToInt()
	if !ok7 || !ok8 || !ok9 {
		return value.Nil, fmt.Errorf("DEBUG.DRAWBOX: r,g,b must be integers 0–255")
	}
	col := rl.Color{R: uint8(ri), G: uint8(gi), B: uint8(bi), A: 255}
	hx, hy, hz := w*0.5, h*0.5, d*0.5
	bb := rl.BoundingBox{
		Min: rl.Vector3{X: cx - hx, Y: cy - hy, Z: cz - hz},
		Max: rl.Vector3{X: cx + hx, Y: cy + hy, Z: cz + hz},
	}
	rl.DrawBoundingBox(bb, col)
	return value.Nil, nil
}
