//go:build cgo

package mbgame

import (
	"fmt"
	"image/color"

	"moonbasic/runtime"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (m *Module) registerDrawHelpers(r runtime.Registrar) {
	// Outlined boxes: use DRAW.RECTLINES (same as a classic "draw box ex" helper).
	r.Register("GAME.DRAWSCREENFLASH", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GAME.DRAWSCREENFLASH expects 1 argument (dt#)")
		}
		dt, ok := argF(args[0])
		if !ok || dt < 0 {
			return value.Nil, fmt.Errorf("GAME.DRAWSCREENFLASH: dt must be a non-negative number")
		}
		if m.flash == nil || m.flash.tRem <= 0 {
			return value.Nil, nil
		}
		w := rl.GetScreenWidth()
		h := rl.GetScreenHeight()
		var aOut float64 = float64(m.flash.a)
		if m.flash.dur > 0 {
			aOut *= m.flash.tRem / m.flash.dur
		}
		if aOut > 0 {
			col := color.RGBA{R: uint8(m.flash.r), G: uint8(m.flash.g), B: uint8(m.flash.b), A: uint8(aOut)}
			rl.DrawRectangle(0, 0, int32(w), int32(h), col)
		}
		m.flash.tRem -= dt
		if m.flash.tRem <= 0 {
			m.flash = nil
		}
		return value.Nil, nil
	}))
	r.Register("GAME.SCREENFLASH", "game", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("GAME.SCREENFLASH expects 5 arguments (r,g,b,a, seconds#)")
		}
		r0, ok1 := argI(args[0])
		g0, ok2 := argI(args[1])
		b0, ok3 := argI(args[2])
		a0, ok4 := argI(args[3])
		sec, ok5 := argF(args[4])
		if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || sec <= 0 {
			return value.Nil, fmt.Errorf("GAME.SCREENFLASH: invalid arguments")
		}
		m.flash = &screenFlashState{r: int(r0), g: int(g0), b: int(b0), a: int(a0), tRem: sec, dur: sec}
		return value.Nil, nil
	}))
}
