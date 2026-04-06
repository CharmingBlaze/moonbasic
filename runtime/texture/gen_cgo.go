//go:build cgo || (windows && !cgo)

package texture

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerTextureGenCmds(m *Module, r runtime.Registrar) {
	r.Register("TEXTURE.GENWHITENOISE", "texture", runtime.AdaptLegacy(m.texGenWhiteNoise))
	r.Register("TEXTURE.GENCHECKED", "texture", runtime.AdaptLegacy(m.texGenChecked))
	r.Register("TEXTURE.GENGRADIENTV", "texture", runtime.AdaptLegacy(m.texGenGradientV))
	r.Register("TEXTURE.GENGRADIENTH", "texture", runtime.AdaptLegacy(m.texGenGradientH))
	r.Register("TEXTURE.GENCOLOR", "texture", runtime.AdaptLegacy(m.texGenColor))
}

func (m *Module) texGenWhiteNoise(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("TEXTURE.GENWHITENOISE: heap not bound")
	}
	if len(args) != 2 && len(args) != 3 {
		return value.Nil, fmt.Errorf("TEXTURE.GENWHITENOISE expects (w, h) or (w, h, factor#)")
	}
	w, ok0 := argI32(args[0])
	h, ok1 := argI32(args[1])
	if !ok0 || !ok1 {
		return value.Nil, fmt.Errorf("TEXTURE.GENWHITENOISE: w,h must be numeric")
	}
	factor := float32(1)
	if len(args) == 3 {
		if f, ok := argF32(args[2]); ok {
			factor = f
		} else {
			return value.Nil, fmt.Errorf("TEXTURE.GENWHITENOISE: factor must be numeric")
		}
	}
	im := rl.GenImageWhiteNoise(int(w), int(h), factor)
	defer rl.UnloadImage(im)
	tex := rl.LoadTextureFromImage(im)
	id, err := m.h.Alloc(&TextureObject{Tex: tex})
	if err != nil {
		rl.UnloadTexture(tex)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func argF32(v value.Value) (float32, bool) {
	if f, ok := v.ToFloat(); ok {
		return float32(f), true
	}
	if i, ok := v.ToInt(); ok {
		return float32(i), true
	}
	return 0, false
}

func (m *Module) texGenChecked(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("TEXTURE.GENCHECKED: heap not bound")
	}
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("TEXTURE.GENCHECKED expects (w, h, tileW, tileH, color1, color2)")
	}
	w, ok0 := argI32(args[0])
	h, ok1 := argI32(args[1])
	tw, ok2 := argI32(args[2])
	th, ok3 := argI32(args[3])
	if !ok0 || !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("TEXTURE.GENCHECKED: dimensions must be numeric")
	}
	if args[4].Kind != value.KindHandle || args[5].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TEXTURE.GENCHECKED: colors must be COLOR handles")
	}
	c1, err := mbmatrix.HeapColorRGBA(m.h, heap.Handle(args[4].IVal))
	if err != nil {
		return value.Nil, err
	}
	c2, err := mbmatrix.HeapColorRGBA(m.h, heap.Handle(args[5].IVal))
	if err != nil {
		return value.Nil, err
	}
	im := rl.GenImageChecked(int(w), int(h), int(tw), int(th), c1, c2)
	defer rl.UnloadImage(im)
	tex := rl.LoadTextureFromImage(im)
	id, err := m.h.Alloc(&TextureObject{Tex: tex})
	if err != nil {
		rl.UnloadTexture(tex)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) texGenGradientV(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("TEXTURE.GENGRADIENTV: heap not bound")
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("TEXTURE.GENGRADIENTV expects (w, h, topColor, bottomColor)")
	}
	w, ok0 := argI32(args[0])
	h, ok1 := argI32(args[1])
	if !ok0 || !ok1 {
		return value.Nil, fmt.Errorf("TEXTURE.GENGRADIENTV: w,h must be numeric")
	}
	if args[2].Kind != value.KindHandle || args[3].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TEXTURE.GENGRADIENTV: colors must be COLOR handles")
	}
	top, err := mbmatrix.HeapColorRGBA(m.h, heap.Handle(args[2].IVal))
	if err != nil {
		return value.Nil, err
	}
	bot, err := mbmatrix.HeapColorRGBA(m.h, heap.Handle(args[3].IVal))
	if err != nil {
		return value.Nil, err
	}
	im := rl.GenImageGradientLinear(int(w), int(h), 0, top, bot)
	defer rl.UnloadImage(im)
	tex := rl.LoadTextureFromImage(im)
	id, err := m.h.Alloc(&TextureObject{Tex: tex})
	if err != nil {
		rl.UnloadTexture(tex)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) texGenGradientH(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("TEXTURE.GENGRADIENTH: heap not bound")
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("TEXTURE.GENGRADIENTH expects (w, h, leftColor, rightColor)")
	}
	w, ok0 := argI32(args[0])
	h, ok1 := argI32(args[1])
	if !ok0 || !ok1 {
		return value.Nil, fmt.Errorf("TEXTURE.GENGRADIENTH: w,h must be numeric")
	}
	if args[2].Kind != value.KindHandle || args[3].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TEXTURE.GENGRADIENTH: colors must be COLOR handles")
	}
	left, err := mbmatrix.HeapColorRGBA(m.h, heap.Handle(args[2].IVal))
	if err != nil {
		return value.Nil, err
	}
	right, err := mbmatrix.HeapColorRGBA(m.h, heap.Handle(args[3].IVal))
	if err != nil {
		return value.Nil, err
	}
	im := rl.GenImageGradientLinear(int(w), int(h), 90, left, right)
	defer rl.UnloadImage(im)
	tex := rl.LoadTextureFromImage(im)
	id, err := m.h.Alloc(&TextureObject{Tex: tex})
	if err != nil {
		rl.UnloadTexture(tex)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) texGenColor(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("TEXTURE.GENCOLOR: heap not bound")
	}
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("TEXTURE.GENCOLOR expects (w, h, r, g, b, a)")
	}
	w, ok0 := argI32(args[0])
	h, ok1 := argI32(args[1])
	if !ok0 || !ok1 {
		return value.Nil, fmt.Errorf("TEXTURE.GENCOLOR: w,h must be numeric")
	}
	r_, ok2 := argU8(args[2])
	g_, ok3 := argU8(args[3])
	b_, ok4 := argU8(args[4])
	a_, ok5 := argU8(args[5])
	if !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("TEXTURE.GENCOLOR: rgba must be numeric")
	}
	col := color.RGBA{R: r_, G: g_, B: b_, A: a_}
	im := rl.GenImageColor(int(w), int(h), col)
	defer rl.UnloadImage(im)
	tex := rl.LoadTextureFromImage(im)
	id, err := m.h.Alloc(&TextureObject{Tex: tex})
	if err != nil {
		rl.UnloadTexture(tex)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func argI32(v value.Value) (int32, bool) {
	if i, ok := v.ToInt(); ok {
		return int32(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return int32(f), true
	}
	return 0, false
}

func argU8(v value.Value) (uint8, bool) {
	if i, ok := v.ToInt(); ok {
		if i < 0 {
			return 0, false
		}
		if i > 255 {
			return 255, true
		}
		return uint8(i), true
	}
	if f, ok := v.ToFloat(); ok {
		if f < 0 {
			return 0, false
		}
		if f > 255 {
			return 255, true
		}
		return uint8(f), true
	}
	return 0, false
}
