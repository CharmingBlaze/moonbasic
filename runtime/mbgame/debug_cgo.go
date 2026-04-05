//go:build cgo

package mbgame

import (
	"fmt"
	"image/color"

	"moonbasic/runtime"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (m *Module) registerDebugDraw(r runtime.Registrar) {
	_ = m
	r.Register("GAME.DEBUGRECT", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if rt != nil && !rt.DebugMode {
			return value.Nil, nil
		}
		if len(args) != 8 {
			return value.Nil, fmt.Errorf("GAME.DEBUGRECT expects 8 arguments (x,y,w,h, r,g,b,a)")
		}
		x, ok1 := argI(args[0])
		y, ok2 := argI(args[1])
		w, ok3 := argI(args[2])
		h, ok4 := argI(args[3])
		cr, ok5 := argI(args[4])
		cg, ok6 := argI(args[5])
		cb, ok7 := argI(args[6])
		ca, ok8 := argI(args[7])
		if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 {
			return value.Nil, fmt.Errorf("GAME.DEBUGRECT: arguments must be numeric")
		}
		col := color.RGBA{R: uint8(cr), G: uint8(cg), B: uint8(cb), A: uint8(ca)}
		rl.DrawRectangleLinesEx(rl.Rectangle{X: float32(x), Y: float32(y), Width: float32(w), Height: float32(h)}, 1, col)
		return value.Nil, nil
	})
}
