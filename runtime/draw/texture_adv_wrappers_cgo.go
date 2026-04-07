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

// drawTexRecObj — state for DRAW.TEXTUREREC (sub-rect + position + tint).
type drawTexRecObj struct {
	release        heap.ReleaseOnce
	tex            heap.Handle
	srcX, srcY     float32
	srcW, srcH     float32
	posX, posY     float32
	cr, cg, cb, ca int32
}

func (o *drawTexRecObj) Free() { o.release.Do(func() {}) }

func (o *drawTexRecObj) TypeTag() uint16 { return heap.TagTextureDraw }

func (o *drawTexRecObj) TypeName() string { return "DRAWTEXREC" }

// drawTexProObj — state for DRAW.TEXTUREPRO (full pro quad).
type drawTexProObj struct {
	release        heap.ReleaseOnce
	tex            heap.Handle
	srcX, srcY     float32
	srcW, srcH     float32
	dstX, dstY     float32
	dstW, dstH     float32
	ox, oy         float32
	rot            float32
	cr, cg, cb, ca int32
}

func (o *drawTexProObj) Free() { o.release.Do(func() {}) }

func (o *drawTexProObj) TypeTag() uint16 { return heap.TagTextureDraw }

func (o *drawTexProObj) TypeName() string { return "DRAWTEXPRO" }

func registerTextureAdvWrappers(m *Module, r runtime.Registrar) {
	r.Register("DRAWTEXREC.SRC", "draw", m.texRecSrc)
	r.Register("DRAWTEXREC.POS", "draw", m.texRecPos)
	r.Register("DRAWTEXREC.COLOR", "draw", m.texRecColor)
	r.Register("DRAWTEXREC.COL", "draw", m.texRecColor)
	r.Register("DRAWTEXREC.SETTEXTURE", "draw", m.texRecSetTexture)
	r.Register("DRAWTEXREC.DRAW", "draw", m.texRecDraw)
	r.Register("DRAWTEXREC.FREE", "draw", m.texRecFree)
	r.Register("DRAWTEXREC", "draw", runtime.AdaptLegacy(m.makeDrawTexRec))

	r.Register("DRAWTEXPRO.SRC", "draw", m.texProSrc)
	r.Register("DRAWTEXPRO.DST", "draw", m.texProDst)
	r.Register("DRAWTEXPRO.ORIGIN", "draw", m.texProOrigin)
	r.Register("DRAWTEXPRO.ROT", "draw", m.texProRot)
	r.Register("DRAWTEXPRO.COLOR", "draw", m.texProColor)
	r.Register("DRAWTEXPRO.COL", "draw", m.texProColor)
	r.Register("DRAWTEXPRO.SETTEXTURE", "draw", m.texProSetTexture)
	r.Register("DRAWTEXPRO.DRAW", "draw", m.texProDraw)
	r.Register("DRAWTEXPRO.FREE", "draw", m.texProFree)
	r.Register("DRAWTEXPRO", "draw", runtime.AdaptLegacy(m.makeDrawTexPro))
}

func castTexRec(h *heap.Store, v value.Value) (*drawTexRecObj, error) {
	if v.Kind != value.KindHandle {
		return nil, fmt.Errorf("expected handle")
	}
	return heap.Cast[*drawTexRecObj](h, heap.Handle(v.IVal))
}

func castTexPro(h *heap.Store, v value.Value) (*drawTexProObj, error) {
	if v.Kind != value.KindHandle {
		return nil, fmt.Errorf("expected handle")
	}
	return heap.Cast[*drawTexProObj](h, heap.Handle(v.IVal))
}

func (m *Module) makeDrawTexRec(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("DRAWTEXREC expects 1 argument (textureHandle)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAWTEXREC: texture handle")
	}
	o := &drawTexRecObj{
		tex: heap.Handle(args[0].IVal),
		srcW: 1, srcH: 1,
		cr: 255, cg: 255, cb: 255, ca: 255,
	}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) makeDrawTexPro(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("DRAWTEXPRO expects 1 argument (textureHandle)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAWTEXPRO: texture handle")
	}
	o := &drawTexProObj{
		tex: heap.Handle(args[0].IVal),
		srcW: 1, srcH: 1, dstW: 1, dstH: 1,
		cr: 255, cg: 255, cb: 255, ca: 255,
	}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) texRecSrc(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("DRAWTEXREC.SRC expects (handle, srcX#, srcY#, srcW#, srcH#)")
	}
	o, err := castTexRec(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	sx, ok1 := argFloat(args[1])
	sy, ok2 := argFloat(args[2])
	sw, ok3 := argFloat(args[3])
	sh, ok4 := argFloat(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAWTEXREC.SRC: numeric")
	}
	o.srcX, o.srcY, o.srcW, o.srcH = sx, sy, sw, sh
	return value.Nil, nil
}

func (m *Module) texRecPos(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("DRAWTEXREC.POS expects (handle, x#, y#)")
	}
	o, err := castTexRec(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argFloat(args[1])
	y, ok2 := argFloat(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("DRAWTEXREC.POS: numeric")
	}
	o.posX, o.posY = x, y
	return value.Nil, nil
}

func (m *Module) texRecColor(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("DRAWTEXREC.COLOR expects (handle, r,g,b,a)")
	}
	o, err := castTexRec(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	r, ok1 := argInt(args[1])
	g, ok2 := argInt(args[2])
	b, ok3 := argInt(args[3])
	a, ok4 := argInt(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAWTEXREC.COLOR: numeric")
	}
	o.cr, o.cg, o.cb, o.ca = int32(r), int32(g), int32(b), int32(a)
	return value.Nil, nil
}

func (m *Module) texRecSetTexture(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAWTEXREC.SETTEXTURE expects (handle, texHandle)")
	}
	o, err := castTexRec(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	if args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAWTEXREC.SETTEXTURE: handle")
	}
	o.tex = heap.Handle(args[1].IVal)
	return value.Nil, nil
}

func (m *Module) texRecDraw(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("DRAWTEXREC.DRAW expects (handle)")
	}
	o, err := castTexRec(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	tex, err := texture.ForBinding(m.h, o.tex)
	if err != nil {
		return value.Nil, err
	}
	c := convert.NewColor4(o.cr, o.cg, o.cb, o.ca)
	tint := color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
	sourceRec := rl.Rectangle{X: o.srcX, Y: o.srcY, Width: o.srcW, Height: o.srcH}
	pos := rl.Vector2{X: o.posX, Y: o.posY}
	rl.DrawTextureRec(tex, sourceRec, pos, tint)
	return value.Nil, nil
}

func (m *Module) texRecFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAWTEXREC.FREE expects handle")
	}
	if _, err := castTexRec(m.h, args[0]); err != nil {
		return value.Nil, err
	}
	return value.Nil, rt.Heap.Free(heap.Handle(args[0].IVal))
}

func (m *Module) texProSrc(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("DRAWTEXPRO.SRC expects (handle, srcX#, srcY#, srcW#, srcH#)")
	}
	o, err := castTexPro(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	sx, ok1 := argFloat(args[1])
	sy, ok2 := argFloat(args[2])
	sw, ok3 := argFloat(args[3])
	sh, ok4 := argFloat(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAWTEXPRO.SRC: numeric")
	}
	o.srcX, o.srcY, o.srcW, o.srcH = sx, sy, sw, sh
	return value.Nil, nil
}

func (m *Module) texProDst(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("DRAWTEXPRO.DST expects (handle, dstX#, dstY#, dstW#, dstH#)")
	}
	o, err := castTexPro(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	dx, ok1 := argFloat(args[1])
	dy, ok2 := argFloat(args[2])
	dw, ok3 := argFloat(args[3])
	dh, ok4 := argFloat(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAWTEXPRO.DST: numeric")
	}
	o.dstX, o.dstY, o.dstW, o.dstH = dx, dy, dw, dh
	return value.Nil, nil
}

func (m *Module) texProOrigin(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("DRAWTEXPRO.ORIGIN expects (handle, ox#, oy#)")
	}
	o, err := castTexPro(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	ox, ok1 := argFloat(args[1])
	oy, ok2 := argFloat(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("DRAWTEXPRO.ORIGIN: numeric")
	}
	o.ox, o.oy = ox, oy
	return value.Nil, nil
}

func (m *Module) texProRot(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAWTEXPRO.ROT expects (handle, radians#)")
	}
	o, err := castTexPro(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	rad, ok := argFloat(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("DRAWTEXPRO.ROT: numeric")
	}
	o.rot = rad
	return value.Nil, nil
}

func (m *Module) texProColor(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("DRAWTEXPRO.COLOR expects (handle, r,g,b,a)")
	}
	o, err := castTexPro(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	r, ok1 := argInt(args[1])
	g, ok2 := argInt(args[2])
	b, ok3 := argInt(args[3])
	a, ok4 := argInt(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAWTEXPRO.COLOR: numeric")
	}
	o.cr, o.cg, o.cb, o.ca = int32(r), int32(g), int32(b), int32(a)
	return value.Nil, nil
}

func (m *Module) texProSetTexture(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAWTEXPRO.SETTEXTURE expects (handle, texHandle)")
	}
	o, err := castTexPro(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	if args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAWTEXPRO.SETTEXTURE: handle")
	}
	o.tex = heap.Handle(args[1].IVal)
	return value.Nil, nil
}

func (m *Module) texProDraw(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("DRAWTEXPRO.DRAW expects (handle)")
	}
	o, err := castTexPro(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	tex, err := texture.ForBinding(m.h, o.tex)
	if err != nil {
		return value.Nil, err
	}
	c := convert.NewColor4(o.cr, o.cg, o.cb, o.ca)
	tint := color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
	sourceRec := rl.Rectangle{X: o.srcX, Y: o.srcY, Width: o.srcW, Height: o.srcH}
	destRec := rl.Rectangle{X: o.dstX, Y: o.dstY, Width: o.dstW, Height: o.dstH}
	origin := rl.Vector2{X: o.ox, Y: o.oy}
	rl.DrawTexturePro(tex, sourceRec, destRec, origin, o.rot, tint)
	return value.Nil, nil
}

func (m *Module) texProFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAWTEXPRO.FREE expects handle")
	}
	if _, err := castTexPro(m.h, args[0]); err != nil {
		return value.Nil, err
	}
	return value.Nil, rt.Heap.Free(heap.Handle(args[0].IVal))
}
