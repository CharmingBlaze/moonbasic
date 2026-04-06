//go:build cgo || (windows && !cgo)

package mbdraw

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/convert"
	"moonbasic/runtime/texture"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) textureFromArg(v value.Value) (rl.Texture2D, error) {
	if v.Kind != value.KindHandle {
		return rl.Texture2D{}, fmt.Errorf("texture handle required")
	}
	return texture.ForBinding(m.h, heap.Handle(v.IVal))
}

// TextureForBinding resolves a texture handle for Raylib use (decals, particles, etc.).
func TextureForBinding(store *heap.Store, h heap.Handle) (rl.Texture2D, error) {
	return texture.ForBinding(store, h)
}

func registerTextureCmds(m *Module, r runtime.Registrar) {
	r.Register("DRAW.TEXTURE", "draw", runtime.AdaptLegacy(m.drawTexture))
	r.Register("DRAW.TEXTUREV", "draw", runtime.AdaptLegacy(m.drawTextureV))
	r.Register("DRAW.TEXTUREEX", "draw", runtime.AdaptLegacy(m.drawTextureEx))
	r.Register("DRAW.TEXTUREREC", "draw", runtime.AdaptLegacy(m.drawTextureRec))
	r.Register("DRAW.TEXTUREPRO", "draw", runtime.AdaptLegacy(m.drawTexturePro))
	r.Register("DRAW.TEXTUREFULL", "draw", runtime.AdaptLegacy(m.drawTextureFull))
	r.Register("DRAW.TEXTUREFLIPPED", "draw", runtime.AdaptLegacy(m.drawTextureFlipped))
	r.Register("DRAW.TEXTURETILED", "draw", runtime.AdaptLegacy(m.drawTextureTiled))
	r.Register("DRAW.TEXTURENPATCH", "draw", runtime.AdaptLegacy(m.drawTextureNPatch))
}

func (m *Module) drawTexture(args []value.Value) (value.Value, error) {
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("DRAW.TEXTURE expects 7 arguments (handle, x, y, r, g, b, a)")
	}
	tex, err := m.textureFromArg(args[0])
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argInt(args[1])
	y, ok2 := argInt(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("DRAW.TEXTURE: x,y must be numeric")
	}
	var rgb [4]int32
	for i := 0; i < 4; i++ {
		v, ok := argInt(args[3+i])
		if !ok {
			return value.Nil, fmt.Errorf("DRAW.TEXTURE: tint must be numeric")
		}
		rgb[i] = v
	}
	c := convert.NewColor4(rgb[0], rgb[1], rgb[2], rgb[3])
	tint := color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
	rl.DrawTexture(tex, int32(x), int32(y), tint)
	return value.Nil, nil
}

func (m *Module) drawTextureTiled(args []value.Value) (value.Value, error) {
	if len(args) != 17 {
		return value.Nil, fmt.Errorf("DRAW.TEXTURETILED expects 17 arguments (handle, srcx,srcy,srcw,srch, dstx,dsty,dstw,dsth, ox,oy, rot, scale, r,g,b,a)")
	}
	tex, err := m.textureFromArg(args[0])
	if err != nil {
		return value.Nil, err
	}
	srcx, ok1 := argFloat(args[1])
	srcy, ok2 := argFloat(args[2])
	srcw, ok3 := argFloat(args[3])
	srch, ok4 := argFloat(args[4])
	dstx, ok5 := argFloat(args[5])
	dsty, ok6 := argFloat(args[6])
	dstw, ok7 := argFloat(args[7])
	dsth, ok8 := argFloat(args[8])
	ox, ok9 := argFloat(args[9])
	oy, ok10 := argFloat(args[10])
	rot, ok11 := argFloat(args[11])
	scale, ok12 := argFloat(args[12])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 || !ok9 || !ok10 || !ok11 || !ok12 {
		return value.Nil, fmt.Errorf("DRAW.TEXTURETILED: geometry arguments must be numeric")
	}
	var rgb [4]int32
	for i := 0; i < 4; i++ {
		v, ok := argInt(args[13+i])
		if !ok {
			return value.Nil, fmt.Errorf("DRAW.TEXTURETILED: tint must be numeric")
		}
		rgb[i] = v
	}
	c := convert.NewColor4(rgb[0], rgb[1], rgb[2], rgb[3])
	tint := color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
	sourceRec := rl.Rectangle{X: srcx, Y: srcy, Width: srcw, Height: srch}
	destRec := rl.Rectangle{X: dstx, Y: dsty, Width: dstw, Height: dsth}
	origin := rl.Vector2{X: ox, Y: oy}
	drawTextureTiledRaylib(tex, sourceRec, destRec, origin, rot, scale, tint)
	return value.Nil, nil
}

// drawTextureTiledRaylib matches raylib's DrawTextureTiled using DrawTexturePro.
// raylib-go does not expose DrawTextureTiled; rotation/origin non-zero falls back to a single DrawTexturePro over destRec (stretched, not tiled).
func drawTextureTiledRaylib(tex rl.Texture2D, sourceRec, destRec rl.Rectangle, origin rl.Vector2, rotation, scale float32, tint color.RGBA) {
	if scale <= 0 || sourceRec.Width <= 0 || sourceRec.Height <= 0 {
		return
	}
	if rotation != 0 || origin.X != 0 || origin.Y != 0 {
		rl.DrawTexturePro(tex, sourceRec, destRec, origin, rotation, tint)
		return
	}
	tileW := sourceRec.Width * scale
	tileH := sourceRec.Height * scale
	if tileW <= 0 || tileH <= 0 {
		return
	}
	destRight := destRec.X + destRec.Width
	destBottom := destRec.Y + destRec.Height
	for tileTop := destRec.Y; tileTop < destBottom; tileTop += tileH {
		for tileLeft := destRec.X; tileLeft < destRight; tileLeft += tileW {
			visL := maxFloat32(tileLeft, destRec.X)
			visR := minFloat32(tileLeft+tileW, destRight)
			visT := maxFloat32(tileTop, destRec.Y)
			visB := minFloat32(tileTop+tileH, destBottom)
			if visL >= visR || visT >= visB {
				continue
			}
			u := visL - tileLeft
			v := visT - tileTop
			srcX := sourceRec.X + u/scale
			srcY := sourceRec.Y + v/scale
			srcW := (visR - visL) / scale
			srcH := (visB - visT) / scale
			rl.DrawTexturePro(tex,
				rl.Rectangle{X: srcX, Y: srcY, Width: srcW, Height: srcH},
				rl.Rectangle{X: visL, Y: visT, Width: visR - visL, Height: visB - visT},
				rl.Vector2{}, 0, tint)
		}
	}
}

func maxFloat32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func minFloat32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func (m *Module) drawTextureFull(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("DRAW.TEXTUREFULL expects 1 argument (handle)")
	}
	tex, err := m.textureFromArg(args[0])
	if err != nil {
		return value.Nil, err
	}
	screenW := float32(rl.GetScreenWidth())
	screenH := float32(rl.GetScreenHeight())
	sourceRec := rl.Rectangle{X: 0, Y: 0, Width: float32(tex.Width), Height: float32(tex.Height)}
	destRec := rl.Rectangle{X: 0, Y: 0, Width: screenW, Height: screenH}
	rl.DrawTexturePro(tex, sourceRec, destRec, rl.Vector2{}, 0, rl.White)
	return value.Nil, nil
}

func (m *Module) drawTextureFlipped(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("DRAW.TEXTUREFLIPPED expects 1 argument (render texture handle)")
	}
	tex, err := m.textureFromArg(args[0])
	if err != nil {
		return value.Nil, err
	}
	screenW := float32(rl.GetScreenWidth())
	screenH := float32(rl.GetScreenHeight())
	sourceRec := rl.Rectangle{X: 0, Y: 0, Width: float32(tex.Width), Height: -float32(tex.Height)}
	destRec := rl.Rectangle{X: 0, Y: 0, Width: screenW, Height: screenH}
	rl.DrawTexturePro(tex, sourceRec, destRec, rl.Vector2{}, 0, rl.White)
	return value.Nil, nil
}

func (m *Module) drawTextureRec(args []value.Value) (value.Value, error) {
	if len(args) != 11 {
		return value.Nil, fmt.Errorf("DRAW.TEXTUREREC expects 11 arguments (handle, srcx,srcy,srcw,srch, x,y, r,g,b,a)")
	}
	tex, err := m.textureFromArg(args[0])
	if err != nil {
		return value.Nil, err
	}
	srcx, ok1 := argFloat(args[1])
	srcy, ok2 := argFloat(args[2])
	srcw, ok3 := argFloat(args[3])
	srch, ok4 := argFloat(args[4])
	x, ok5 := argFloat(args[5])
	y, ok6 := argFloat(args[6])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("DRAW.TEXTUREREC: geometry arguments must be numeric")
	}
	var rgb [4]int32
	for i := 0; i < 4; i++ {
		v, ok := argInt(args[7+i])
		if !ok {
			return value.Nil, fmt.Errorf("DRAW.TEXTUREREC: tint must be numeric")
		}
		rgb[i] = v
	}
	c := convert.NewColor4(rgb[0], rgb[1], rgb[2], rgb[3])
	tint := color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
	sourceRec := rl.Rectangle{X: srcx, Y: srcy, Width: srcw, Height: srch}
	pos := rl.Vector2{X: x, Y: y}
	rl.DrawTextureRec(tex, sourceRec, pos, tint)
	return value.Nil, nil
}

func (m *Module) drawTexturePro(args []value.Value) (value.Value, error) {
	if len(args) != 16 {
		return value.Nil, fmt.Errorf("DRAW.TEXTUREPRO expects 16 arguments (handle, srcx,srcy,srcw,srch, dstx,dsty,dstw,dsth, ox,oy, rot, r,g,b,a)")
	}
	tex, err := m.textureFromArg(args[0])
	if err != nil {
		return value.Nil, err
	}
	srcx, ok1 := argFloat(args[1])
	srcy, ok2 := argFloat(args[2])
	srcw, ok3 := argFloat(args[3])
	srch, ok4 := argFloat(args[4])
	dstx, ok5 := argFloat(args[5])
	dsty, ok6 := argFloat(args[6])
	dstw, ok7 := argFloat(args[7])
	dsth, ok8 := argFloat(args[8])
	ox, ok9 := argFloat(args[9])
	oy, ok10 := argFloat(args[10])
	rot, ok11 := argFloat(args[11])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 || !ok9 || !ok10 || !ok11 {
		return value.Nil, fmt.Errorf("DRAW.TEXTUREPRO: geometry arguments must be numeric")
	}
	var rgb [4]int32
	for i := 0; i < 4; i++ {
		v, ok := argInt(args[12+i])
		if !ok {
			return value.Nil, fmt.Errorf("DRAW.TEXTUREPRO: tint must be numeric")
		}
		rgb[i] = v
	}
	c := convert.NewColor4(rgb[0], rgb[1], rgb[2], rgb[3])
	tint := color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
	sourceRec := rl.Rectangle{X: srcx, Y: srcy, Width: srcw, Height: srch}
	destRec := rl.Rectangle{X: dstx, Y: dsty, Width: dstw, Height: dsth}
	origin := rl.Vector2{X: ox, Y: oy}
	rl.DrawTexturePro(tex, sourceRec, destRec, origin, rot, tint)
	return value.Nil, nil
}

func (m *Module) drawTextureV(args []value.Value) (value.Value, error) {
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("DRAW.TEXTUREV expects 7 arguments (handle, x, y, r, g, b, a)")
	}
	tex, err := m.textureFromArg(args[0])
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argFloat(args[1])
	y, ok2 := argFloat(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("DRAW.TEXTUREV: x,y must be numeric")
	}
	var rgb [4]int32
	for i := 0; i < 4; i++ {
		v, ok := argInt(args[3+i])
		if !ok {
			return value.Nil, fmt.Errorf("DRAW.TEXTUREV: tint must be numeric")
		}
		rgb[i] = v
	}
	c := convert.NewColor4(rgb[0], rgb[1], rgb[2], rgb[3])
	tint := color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
	rl.DrawTextureV(tex, rl.Vector2{X: x, Y: y}, tint)
	return value.Nil, nil
}

func (m *Module) drawTextureEx(args []value.Value) (value.Value, error) {
	if len(args) != 9 {
		return value.Nil, fmt.Errorf("DRAW.TEXTUREEX expects 9 arguments (handle, x, y, rot, scale, r, g, b, a)")
	}
	tex, err := m.textureFromArg(args[0])
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argFloat(args[1])
	y, ok2 := argFloat(args[2])
	rot, ok3 := argFloat(args[3])
	scale, ok4 := argFloat(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAW.TEXTUREEX: position, rotation, and scale must be numeric")
	}
	var rgb [4]int32
	for i := 0; i < 4; i++ {
		v, ok := argInt(args[5+i])
		if !ok {
			return value.Nil, fmt.Errorf("DRAW.TEXTUREEX: tint must be numeric")
		}
		rgb[i] = v
	}
	c := convert.NewColor4(rgb[0], rgb[1], rgb[2], rgb[3])
	tint := color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
	rl.DrawTextureEx(tex, rl.Vector2{X: x, Y: y}, rot, scale, tint)
	return value.Nil, nil
}

func (m *Module) drawTextureNPatch(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("DRAW.TEXTURENPATCH: heap not bound")
	}
	if len(args) != 13 {
		return value.Nil, fmt.Errorf("DRAW.TEXTURENPATCH expects 13 arguments (tex, L,T,R,B, x,y,w,h, r,g,b,a)")
	}
	tex, err := m.textureFromArg(args[0])
	if err != nil {
		return value.Nil, err
	}
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
