//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerBrushBlitzAPI(m *Module, r runtime.Registrar) {
	r.Register("LoadBrush", "entity", m.entLoadBrush)
	r.Register("FreeBrush", "entity", runtime.AdaptLegacy(m.entFreeBrush))
	r.Register("BrushColor", "entity", runtime.AdaptLegacy(m.entBrushColor))
	r.Register("BrushAlpha", "entity", runtime.AdaptLegacy(m.entBrushAlpha))
	r.Register("BrushBlend", "entity", runtime.AdaptLegacy(m.entBrushBlend))
	r.Register("GetEntityBrush", "entity", runtime.AdaptLegacy(m.entGetEntityBrush))
	r.Register("PaintSurface", "entity", runtime.AdaptLegacy(m.entPaintSurface))
	r.Register("GetSurfaceBrush", "entity", runtime.AdaptLegacy(m.entGetSurfaceBrush))
}

// entLoadBrush loads a texture from disk, allocates a brush with white base color, and owns the texture handle until FreeBrush.
// LoadBrush(path$ [, flags#, uScale#, vScale#]) — flags reserved for filter/wrap presets (default 1).
func (m *Module) entLoadBrush(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LoadBrush: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("LoadBrush expects (path$ [, flags#, uScale#, vScale#])")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	flags := int64(1)
	uScale := float32(1)
	vScale := float32(1)
	if len(args) >= 2 {
		if fi, ok := args[1].ToInt(); ok {
			flags = fi
		}
	}
	if len(args) >= 3 {
		if v, ok := argF32(args[2]); ok {
			uScale = v
		}
	}
	if len(args) >= 4 {
		if v, ok := argF32(args[3]); ok {
			vScale = v
		}
	}
	_ = flags // reserved (documented); may map to TEXTURE.SETFILTER later
	tex := rl.LoadTexture(path)
	if tex.ID <= 0 {
		return value.Nil, fmt.Errorf("LoadBrush: failed to load %q", path)
	}
	th, err := m.h.Alloc(&textureObj{tex: tex})
	if err != nil {
		rl.UnloadTexture(tex)
		return value.Nil, err
	}
	o := &brushObj{
		r: 255, g: 255, b: 255,
		alpha: -1, uScale: uScale, vScale: vScale,
		texH: th, texOwned: true, blendMode: -1,
	}
	id, err := m.h.Alloc(o)
	if err != nil {
		_ = m.h.Free(th)
		return value.Nil, err
	}
	o.self = id
	return value.FromHandle(id), nil
}

func (m *Module) entFreeBrush(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("FreeBrush: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("FreeBrush expects (brush)")
	}
	bh, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("brush must be handle")
	}
	b := m.getBrush(bh)
	if b == nil {
		return value.Nil, fmt.Errorf("invalid brush")
	}
	if b.texOwned && b.texH != 0 {
		_ = m.h.Free(b.texH)
		b.texH = 0
		b.texOwned = false
	}
	if err := m.h.Free(bh); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) entBrushColor(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("BrushColor expects (brush, r#, g#, b#)")
	}
	bh, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("brush must be handle")
	}
	b := m.getBrush(bh)
	if b == nil {
		return value.Nil, fmt.Errorf("invalid brush")
	}
	rf, ok1 := argF32(args[1])
	gf, ok2 := argF32(args[2])
	bf, ok3 := argF32(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("BrushColor: RGB must be numeric")
	}
	b.r, b.g, b.b = f32ToU8(rf), f32ToU8(gf), f32ToU8(bf)
	return value.Nil, nil
}

func (m *Module) entBrushAlpha(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("BrushAlpha expects (brush, alpha#)")
	}
	bh, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("brush must be handle")
	}
	b := m.getBrush(bh)
	if b == nil {
		return value.Nil, fmt.Errorf("invalid brush")
	}
	a, ok1 := argF32(args[1])
	if !ok1 {
		return value.Nil, fmt.Errorf("alpha must be numeric")
	}
	if a > 1 {
		a /= 255
	}
	if a < 0 {
		a = 0
	}
	if a > 1 {
		a = 1
	}
	b.alpha = a
	return value.Nil, nil
}

func (m *Module) entBrushBlend(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("BrushBlend expects (brush, mode#) — 0=opaque 1=alpha 2=multiply 3=additive")
	}
	bh, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("brush must be handle")
	}
	b := m.getBrush(bh)
	if b == nil {
		return value.Nil, fmt.Errorf("invalid brush")
	}
	mode, okm := args[1].ToInt()
	if !okm {
		return value.Nil, fmt.Errorf("mode must be numeric")
	}
	if mode < 0 || mode > 3 {
		return value.Nil, fmt.Errorf("BrushBlend: mode must be 0..3")
	}
	b.blendMode = int32(mode)
	return value.Nil, nil
}

func (m *Module) entGetEntityBrush(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("GetEntityBrush expects (entity#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	if e.getExt().brushH == 0 {
		return value.FromInt(0), nil
	}
	return value.FromHandle(e.getExt().brushH), nil
}

func (m *Module) entPaintSurface(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PaintSurface expects (surface, brush)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("surface must be mesh builder handle")
	}
	mb, err := castMeshBuilder(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	bh, ok := argHandle(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("brush must be handle")
	}
	if m.getBrush(bh) == nil {
		return value.Nil, fmt.Errorf("invalid brush")
	}
	mb.brushH = bh
	return value.Nil, nil
}

func (m *Module) entGetSurfaceBrush(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("GetSurfaceBrush expects (surface)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("surface must be mesh builder handle")
	}
	mb, err := castMeshBuilder(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if mb.brushH == 0 {
		return value.FromInt(0), nil
	}
	return value.FromHandle(mb.brushH), nil
}
