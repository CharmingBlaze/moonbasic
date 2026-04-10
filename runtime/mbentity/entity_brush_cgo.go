//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// brushObj is a Blitz-style material bundle (color, optional texture, FX flags).
type brushObj struct {
	release heap.ReleaseOnce
	self    heap.Handle

	r, g, b uint8
	alpha   float32 // if >= 0, multiplied into tint; -1 = ignore (use entity alpha only)
	texH    heap.Handle
	texSlot int32 // UV layer index (0 = primary)
	texFrame int32 // reserved for animated sheets
	texOwned bool  // if true, FreeBrush frees texH
	uScale  float32
	vScale  float32
	fx       int32
	shine    float32
	blendMode int32 // -1 = ignore (use FX/entity); else 0–3 Blitz blend (see BrushBlend)
}

func (b *brushObj) TypeName() string { return "Brush" }

func (b *brushObj) TypeTag() uint16 { return heap.TagBrush }

func (b *brushObj) Free() {
	b.release.Do(func() {
		b.texH = 0
	})
}

func (m *Module) getBrush(h heap.Handle) *brushObj {
	if m.h == nil || h == 0 {
		return nil
	}
	o, ok := m.h.Get(h)
	if !ok {
		return nil
	}
	b, ok := o.(*brushObj)
	if !ok {
		return nil
	}
	return b
}

func (m *Module) entTintResolved(e *ent) rl.Color {
	a := e.alpha
	if a < 0 {
		a = 0
	}
	if a > 1 {
		a = 1
	}
	r, g, b := e.r, e.g, e.b
	if bobj := m.getBrush(e.getExt().brushH); bobj != nil {
		r, g, b = bobj.r, bobj.g, bobj.b
		if bobj.alpha >= 0 {
			ba := bobj.alpha
			if ba > 1 {
				ba = 1
			}
			a *= ba
		}
	}
	if a < 0 {
		a = 0
	}
	if a > 1 {
		a = 1
	}
	col := rl.Color{R: r, G: g, B: b, A: uint8(a * 255)}
	if m.brushFullBright(e) {
		col = rl.ColorBrightness(col, 0.35)
	}
	return col
}

func (m *Module) brushFullBright(e *ent) bool {
	if e.fxFlags&1 != 0 {
		return true
	}
	if b := m.getBrush(e.getExt().brushH); b != nil && b.fx&1 != 0 {
		return true
	}
	return false
}

func blitzBrushBlendToRL(m int32) rl.BlendMode {
	switch m {
	case 0, 1:
		return rl.BlendAlpha
	case 2:
		return rl.BlendMultiplied
	case 3:
		return rl.BlendAdditive
	default:
		return rl.BlendAlpha
	}
}

func (m *Module) entDrawBlendMode(e *ent) (use bool, mode rl.BlendMode) {
	if e.blendMode >= 0 {
		return true, rl.BlendMode(e.blendMode)
	}
	if b := m.getBrush(e.getExt().brushH); b != nil && b.blendMode >= 0 {
		return true, blitzBrushBlendToRL(b.blendMode)
	}
	if e.fxFlags&16 != 0 {
		return true, rl.BlendAdditive
	}
	if b := m.getBrush(e.getExt().brushH); b != nil && b.fx&16 != 0 {
		return true, rl.BlendAdditive
	}
	return false, rl.BlendAlpha
}

func (m *Module) entCreateBrush(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CreateBrush: heap not bound")
	}
	var rf, gf, bf float32 = 255, 255, 255
	switch len(args) {
	case 0:
	case 3:
		var ok1, ok2, ok3 bool
		rf, ok1 = argF32(args[0])
		gf, ok2 = argF32(args[1])
		bf, ok3 = argF32(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("CreateBrush: RGB must be numeric")
		}
	default:
		return value.Nil, fmt.Errorf("CreateBrush expects () or (r#, g#, b#)")
	}
	o := &brushObj{r: f32ToU8(rf), g: f32ToU8(gf), b: f32ToU8(bf), alpha: -1, uScale: 1, vScale: 1, blendMode: -1}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	o.self = id
	return value.FromHandle(id), nil
}

func f32ToU8(v float32) uint8 {
	if v > 1 {
		if v > 255 {
			return 255
		}
		if v < 0 {
			return 0
		}
		return uint8(v)
	}
	if v < 0 {
		return 0
	}
	return uint8(v * 255)
}

func (m *Module) entBrushTexture(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("BrushTexture: heap not bound")
	}
	if len(args) < 2 || len(args) > 4 {
		return value.Nil, fmt.Errorf("BrushTexture expects (brush, texture [, frame#, uvIndex#])")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("brush must be handle")
	}
	b := m.getBrush(h)
	if b == nil {
		return value.Nil, fmt.Errorf("invalid brush")
	}
	if b.texOwned && b.texH != 0 {
		_ = m.h.Free(b.texH)
		b.texOwned = false
	}
	th, ok2 := argHandle(args[1])
	if !ok2 {
		return value.Nil, fmt.Errorf("texture must be handle")
	}
	b.texH = th
	b.texFrame = 0
	b.texSlot = 0
	if len(args) >= 3 {
		if fi, ok := args[2].ToInt(); ok {
			b.texFrame = int32(fi)
		}
	}
	if len(args) >= 4 {
		if si, ok := args[3].ToInt(); ok {
			b.texSlot = int32(si)
		}
	}
	return value.Nil, nil
}

func (m *Module) entBrushFX(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("BrushFX expects (brush, fx#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("brush must be handle")
	}
	b := m.getBrush(h)
	if b == nil {
		return value.Nil, fmt.Errorf("invalid brush")
	}
	fx, _ := args[1].ToInt()
	b.fx = int32(fx)
	return value.Nil, nil
}

func (m *Module) entBrushShininess(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("BrushShininess expects (brush, amount#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("brush must be handle")
	}
	b := m.getBrush(h)
	if b == nil {
		return value.Nil, fmt.Errorf("invalid brush")
	}
	s, ok1 := argF32(args[1])
	if !ok1 {
		return value.Nil, fmt.Errorf("amount must be numeric")
	}
	b.shine = s
	return value.Nil, nil
}

func (m *Module) entEntityShadow(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("EntityShadow expects (entity#, state#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	sv, _ := args[1].ToInt()
	e.getExt().shadowCast = int32(sv)
	return value.Nil, nil
}

func (m *Module) entPaintEntity(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PaintEntity expects (entity#, brush)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	bh, ok := argHandle(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("brush must be handle")
	}
	if m.getBrush(bh) == nil {
		return value.Nil, fmt.Errorf("invalid brush")
	}
	e.getExt().brushH = bh
	if b := m.getBrush(bh); b != nil {
		e.r, e.g, e.b = b.r, b.g, b.b
		e.fxFlags = b.fx
		e.shininess = b.shine
		if b.alpha >= 0 {
			e.alpha = b.alpha
		}
		if b.texH != 0 {
			e.texHandle = b.texH
		}
	}
	return value.Nil, nil
}
