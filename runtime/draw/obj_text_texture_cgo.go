//go:build cgo || (windows && !cgo)

package mbdraw

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/convert"
	mbfont "moonbasic/runtime/font"
	"moonbasic/runtime/texture"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// textDrawObj — default font DRAW.TEXT style (8-arg) wrapper.
type textDrawObj struct {
	release heap.ReleaseOnce
	textVal value.Value
	x, y    int32
	size    int32
	cr, cg, cb, ca int32
}

func (o *textDrawObj) Free() { o.release.Do(func() {}) }

func (o *textDrawObj) TypeTag() uint16 { return heap.TagTextDraw }

func (o *textDrawObj) TypeName() string { return "TEXTOBJ" }

// textureDrawObj — DRAW.TEXTUREV-style (position + tint).
type textureDrawObj struct {
	release heap.ReleaseOnce
	tex     heap.Handle
	x, y    int32
	cr, cg, cb, ca int32
}

func (o *textureDrawObj) Free() { o.release.Do(func() {}) }

func (o *textureDrawObj) TypeTag() uint16 { return heap.TagTextureDraw }

func (o *textureDrawObj) TypeName() string { return "DRAWTEX2" }

func registerTextTextureObjs(m *Module, r runtime.Registrar) {
	r.Register("TEXTDRAW.POS", "draw", m.textDrawPos)
	r.Register("TEXTDRAW.SIZE", "draw", m.textDrawSize)
	r.Register("TEXTDRAW.COLOR", "draw", m.textDrawColor)
	r.Register("TEXTDRAW.COL", "draw", m.textDrawColor)
	r.Register("TEXTDRAW.SETTEXT", "draw", m.textDrawSetText)
	r.Register("TEXTDRAW.DRAW", "draw", m.textDrawDraw)
	r.Register("TEXTDRAW.FREE", "draw", m.textDrawFree)
	r.Register("TEXTOBJ", "draw", runtime.AdaptLegacy(m.makeTextObj))

	r.Register("DRAWTEX2.POS", "draw", m.texDrawPos)
	r.Register("DRAWTEX2.COLOR", "draw", m.texDrawColor)
	r.Register("DRAWTEX2.COL", "draw", m.texDrawColor)
	r.Register("DRAWTEX2.SETTEXTURE", "draw", m.texDrawSetTexture)
	r.Register("DRAWTEX2.DRAW", "draw", m.texDrawDraw)
	r.Register("DRAWTEX2.FREE", "draw", m.texDrawFree)
	r.Register("DRAWTEX2", "draw", runtime.AdaptLegacy(m.makeTexDraw2))
}

func castTextDraw(h *heap.Store, v value.Value) (*textDrawObj, error) {
	if v.Kind != value.KindHandle {
		return nil, fmt.Errorf("expected handle")
	}
	return heap.Cast[*textDrawObj](h, heap.Handle(v.IVal))
}

func castTexDraw2(h *heap.Store, v value.Value) (*textureDrawObj, error) {
	if v.Kind != value.KindHandle {
		return nil, fmt.Errorf("expected handle")
	}
	return heap.Cast[*textureDrawObj](h, heap.Handle(v.IVal))
}

func (m *Module) makeTextObj(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("TEXTOBJ expects 1 argument (text$)")
	}
	o := &textDrawObj{
		textVal: args[0],
		size:    20,
		cr: 255, cg: 255, cb: 255, ca: 255,
	}
	return m.allocTextDraw(o)
}

func (m *Module) allocTextDraw(o *textDrawObj) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("draw: heap not bound")
	}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) textDrawPos(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("TEXTDRAW.POS expects (handle, x, y)")
	}
	o, err := castTextDraw(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argInt(args[1])
	y, ok2 := argInt(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("TEXTDRAW.POS: int coordinates")
	}
	o.x, o.y = x, y
	return value.Nil, nil
}

func (m *Module) textDrawSize(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("TEXTDRAW.SIZE expects (handle, px)")
	}
	o, err := castTextDraw(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	s, ok := argInt(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("TEXTDRAW.SIZE: int")
	}
	o.size = s
	return value.Nil, nil
}

func (m *Module) textDrawColor(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("TEXTDRAW.COLOR expects (handle, r,g,b,a)")
	}
	o, err := castTextDraw(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	r, ok1 := argInt(args[1])
	g, ok2 := argInt(args[2])
	b, ok3 := argInt(args[3])
	a, ok4 := argInt(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("TEXTDRAW.COLOR: numeric")
	}
	o.cr, o.cg, o.cb, o.ca = int32(r), int32(g), int32(b), int32(a)
	return value.Nil, nil
}

func (m *Module) textDrawSetText(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("TEXTDRAW.SETTEXT expects (handle, text$)")
	}
	o, err := castTextDraw(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	o.textVal = args[1]
	return value.Nil, nil
}

func (m *Module) textDrawDraw(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("TEXTDRAW.DRAW expects (handle)")
	}
	o, err := castTextDraw(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	text := stringFromRT(rt, o.textVal)
	col := color.RGBA{R: uint8(o.cr), G: uint8(o.cg), B: uint8(o.cb), A: uint8(o.ca)}
	rl.DrawText(text, o.x, o.y, o.size, col)
	return value.Nil, nil
}

func (m *Module) textDrawFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TEXTDRAW.FREE expects handle")
	}
	if _, err := castTextDraw(m.h, args[0]); err != nil {
		return value.Nil, err
	}
	return value.Nil, rt.Heap.Free(heap.Handle(args[0].IVal))
}

func (m *Module) makeTexDraw2(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("DRAWTEX2 expects 1 argument (textureHandle)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAWTEX2: texture handle")
	}
	o := &textureDrawObj{tex: heap.Handle(args[0].IVal), cr: 255, cg: 255, cb: 255, ca: 255}
	h, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func (m *Module) texDrawPos(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("DRAWTEX2.POS expects (handle, x, y)")
	}
	o, err := castTexDraw2(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argInt(args[1])
	y, ok2 := argInt(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("DRAWTEX2.POS: int")
	}
	o.x, o.y = x, y
	return value.Nil, nil
}

func (m *Module) texDrawColor(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("DRAWTEX2.COLOR expects (handle, r,g,b,a)")
	}
	o, err := castTexDraw2(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	r, ok1 := argInt(args[1])
	g, ok2 := argInt(args[2])
	b, ok3 := argInt(args[3])
	a, ok4 := argInt(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAWTEX2.COLOR: numeric")
	}
	o.cr, o.cg, o.cb, o.ca = int32(r), int32(g), int32(b), int32(a)
	return value.Nil, nil
}

func (m *Module) texDrawSetTexture(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAWTEX2.SETTEXTURE expects (handle, texHandle)")
	}
	o, err := castTexDraw2(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	if args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAWTEX2.SETTEXTURE: handle")
	}
	o.tex = heap.Handle(args[1].IVal)
	return value.Nil, nil
}

func (m *Module) texDrawDraw(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("DRAWTEX2.DRAW expects (handle)")
	}
	o, err := castTexDraw2(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	tex, err := texture.ForBinding(m.h, o.tex)
	if err != nil {
		return value.Nil, err
	}
	c := convert.NewColor4(o.cr, o.cg, o.cb, o.ca)
	tint := color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
	rl.DrawTexture(tex, o.x, o.y, tint)
	return value.Nil, nil
}

func (m *Module) texDrawFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAWTEX2.FREE expects handle")
	}
	if _, err := castTexDraw2(m.h, args[0]); err != nil {
		return value.Nil, err
	}
	return value.Nil, rt.Heap.Free(heap.Handle(args[0].IVal))
}

// --- Optional TEXTOBJEX (font + DrawTextEx) ---
type textDrawExObj struct {
	release  heap.ReleaseOnce
	fontH    heap.Handle
	textVal  value.Value
	x, y     float32
	size     float32
	spacing  float32
	cr, cg, cb, ca int32
}

func (o *textDrawExObj) Free() { o.release.Do(func() {}) }

func (o *textDrawExObj) TypeTag() uint16 { return heap.TagTextDrawEx }

func (o *textDrawExObj) TypeName() string { return "TEXTOBJEX" }

func registerTextExObj(m *Module, r runtime.Registrar) {
	r.Register("TEXTEXOBJ.POS", "draw", m.textExObjPos)
	r.Register("TEXTEXOBJ.SIZE", "draw", m.textExObjSize)
	r.Register("TEXTEXOBJ.SPACING", "draw", m.textExObjSpacing)
	r.Register("TEXTEXOBJ.COLOR", "draw", m.textExObjColor)
	r.Register("TEXTEXOBJ.SETTEXT", "draw", m.textExObjSetText)
	r.Register("TEXTEXOBJ.DRAW", "draw", m.textExObjDraw)
	r.Register("TEXTEXOBJ.FREE", "draw", m.textExObjFree)
	r.Register("TEXTOBJEX", "draw", runtime.AdaptLegacy(m.makeTextObjEx))
}

func (m *Module) makeTextObjEx(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("TEXTOBJEX expects 2 arguments (fontHandle, text$)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TEXTOBJEX: font handle")
	}
	o := &textDrawExObj{fontH: heap.Handle(args[0].IVal), textVal: args[1], size: 20, spacing: 1, cr: 255, cg: 255, cb: 255, ca: 255}
	h, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func castTextEx(h *heap.Store, v value.Value) (*textDrawExObj, error) {
	return heap.Cast[*textDrawExObj](h, heap.Handle(v.IVal))
}

func (m *Module) textExObjPos(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("TEXTEXOBJ.POS expects (handle, x#, y#)")
	}
	o, err := castTextEx(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argFloat(args[1])
	y, ok2 := argFloat(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("TEXTEXOBJ.POS: float")
	}
	o.x, o.y = x, y
	return value.Nil, nil
}

func (m *Module) textExObjSize(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("TEXTEXOBJ.SIZE expects (handle, px#)")
	}
	o, err := castTextEx(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	s, ok := argFloat(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("TEXTEXOBJ.SIZE: float")
	}
	o.size = s
	return value.Nil, nil
}

func (m *Module) textExObjSpacing(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("TEXTEXOBJ.SPACING expects (handle, spacing#)")
	}
	o, err := castTextEx(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	s, ok := argFloat(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("TEXTEXOBJ.SPACING: float")
	}
	o.spacing = s
	return value.Nil, nil
}

func (m *Module) textExObjColor(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("TEXTEXOBJ.COLOR expects (handle, r,g,b,a)")
	}
	o, err := castTextEx(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	r, ok1 := argInt(args[1])
	g, ok2 := argInt(args[2])
	b, ok3 := argInt(args[3])
	a, ok4 := argInt(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("TEXTEXOBJ.COLOR: numeric")
	}
	o.cr, o.cg, o.cb, o.ca = int32(r), int32(g), int32(b), int32(a)
	return value.Nil, nil
}

func (m *Module) textExObjSetText(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("TEXTEXOBJ.SETTEXT expects (handle, text$)")
	}
	o, err := castTextEx(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	o.textVal = args[1]
	return value.Nil, nil
}

func (m *Module) textExObjDraw(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("TEXTEXOBJ.DRAW expects (handle)")
	}
	o, err := castTextEx(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	font, err := mbfont.FontForHandle(m.h, o.fontH)
	if err != nil {
		return value.Nil, err
	}
	text := stringFromRT(rt, o.textVal)
	col := color.RGBA{R: uint8(o.cr), G: uint8(o.cg), B: uint8(o.cb), A: uint8(o.ca)}
	rl.DrawTextEx(font, text, rl.Vector2{X: o.x, Y: o.y}, o.size, o.spacing, col)
	return value.Nil, nil
}

func (m *Module) textExObjFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TEXTEXOBJ.FREE expects handle")
	}
	if _, err := castTextEx(m.h, args[0]); err != nil {
		return value.Nil, err
	}
	return value.Nil, rt.Heap.Free(heap.Handle(args[0].IVal))
}
