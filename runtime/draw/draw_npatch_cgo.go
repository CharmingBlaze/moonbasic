//go:build cgo

package mbdraw

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/convert"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerTextureNPatch(r runtime.Registrar) {
	r.Register("DRAW.TEXTURENPATCH", "draw", runtime.AdaptLegacy(m.drawTextureNPatch))
}

// Args: texture, left, top, right, bottom, dest x,y,w,h, tint r,g,b,a — Source uses full texture; rotation 0.
func (m *Module) drawTextureNPatch(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("DRAW.TEXTURENPATCH: heap not bound")
	}
	if len(args) != 13 {
		return value.Nil, fmt.Errorf("DRAW.TEXTURENPATCH expects 13 arguments (tex, L,T,R,B, x,y,w,h, r,g,b,a)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAW.TEXTURENPATCH: first argument must be texture handle")
	}
	obj, err := heap.Cast[*rlTex](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	tex := obj.t
	var border [4]int32
	for i := 0; i < 4; i++ {
		v, ok := argInt(args[1+i])
		if !ok {
			return value.Nil, fmt.Errorf("DRAW.TEXTURENPATCH: border %d must be numeric", i+1)
		}
		border[i] = v
	}
	left, top, right, bottom := border[0], border[1], border[2], border[3]
	var xywh [4]int32
	for i := 0; i < 4; i++ {
		v, ok := argInt(args[5+i])
		if !ok {
			return value.Nil, fmt.Errorf("DRAW.TEXTURENPATCH: dest must be numeric")
		}
		xywh[i] = v
	}
	var rgb [4]int32
	for i := 0; i < 4; i++ {
		v, ok := argInt(args[9+i])
		if !ok {
			return value.Nil, fmt.Errorf("DRAW.TEXTURENPATCH: tint must be numeric")
		}
		rgb[i] = v
	}
	c := convert.NewColor4(rgb[0], rgb[1], rgb[2], rgb[3])
	tint := color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
	np := rl.NPatchInfo{
		Source: rl.Rectangle{X: 0, Y: 0, Width: float32(tex.Width), Height: float32(tex.Height)},
		Left:   left, Top: top, Right: right, Bottom: bottom,
		Layout: rl.NPatchNinePatch,
	}
	dest := rl.Rectangle{X: float32(xywh[0]), Y: float32(xywh[1]), Width: float32(xywh[2]), Height: float32(xywh[3])}
	rl.DrawTextureNPatch(tex, np, dest, rl.Vector2{}, 0, tint)
	return value.Nil, nil
}
