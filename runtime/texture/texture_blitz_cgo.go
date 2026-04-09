//go:build cgo || (windows && !cgo)

package texture

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerTextureBlitzCmds(m *Module, r runtime.Registrar) {
	r.Register("CreateTexture", "texture", m.texCreateBlank)
	r.Register("LoadAnimTexture", "texture", m.texLoadAnimStrip)
	r.Register("TextureWidth", "texture", runtime.AdaptLegacy(m.texWidth))
	r.Register("TextureHeight", "texture", runtime.AdaptLegacy(m.texHeight))
	r.Register("TextureName$", "texture", m.texNameStr)
	r.Register("SetCubeFace", "texture", runtime.AdaptLegacy(m.texSetCubeFace))
	r.Register("SetCubeMode", "texture", runtime.AdaptLegacy(m.texSetCubeMode))
	r.Register("TextureCoords", "texture", runtime.AdaptLegacy(m.texCoordsMode))
	r.Register("ScaleTexture", "texture", runtime.AdaptLegacy(m.texScaleUV))
	r.Register("RotateTexture", "texture", runtime.AdaptLegacy(m.texRotateUV))
	r.Register("PositionTexture", "texture", runtime.AdaptLegacy(m.texPositionUV))
}

func (m *Module) texObj1(args []value.Value, op string) (*TextureObject, error) {
	if m.h == nil {
		return nil, fmt.Errorf("%s: heap not bound", op)
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s expects texture handle", op)
	}
	o, err := heap.Cast[*TextureObject](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return nil, err
	}
	return o, nil
}

// CreateTexture(width#, height# [, flags#]) — blank RGBA texture; flags passed to loader preset (default 1).
func (m *Module) texCreateBlank(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CreateTexture: heap not bound")
	}
	if len(args) < 2 || len(args) > 3 {
		return value.Nil, fmt.Errorf("CreateTexture expects (width#, height# [, flags#])")
	}
	w, ok1 := argDim(args[0])
	h, ok2 := argDim(args[1])
	if !ok1 || !ok2 || w < 1 || h < 1 {
		return value.Nil, fmt.Errorf("CreateTexture: width/height must be positive integers")
	}
	flags := int32(1)
	if len(args) == 3 {
		if fi, ok := args[2].ToInt(); ok {
			flags = int32(fi)
		}
	}
	im := rl.GenImageColor(int(w), int(h), rl.Blank)
	defer rl.UnloadImage(im)
	t := rl.LoadTextureFromImage(im)
	obj := &TextureObject{Tex: t, loaded: true, Flags: flags, UScl: 1, VScl: 1}
	obj.setFinalizer()
	texApplyLoadFlags(&obj.Tex, flags)
	id, err := m.h.Alloc(obj)
	if err != nil {
		rl.UnloadTexture(t)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func argDim(v value.Value) (int, bool) {
	if i, ok := v.ToInt(); ok && i > 0 {
		return int(i), true
	}
	if f, ok := v.ToFloat(); ok && f > 0 {
		return int(f), true
	}
	return 0, false
}

// LoadAnimTexture(path$, flags#, cellW#, cellH#, firstFrame#, frameCount#) — loads one cell from a horizontal strip as a standalone texture.
func (m *Module) texLoadAnimStrip(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LoadAnimTexture: heap not bound")
	}
	if len(args) != 6 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("LoadAnimTexture expects (path$, flags#, cellW#, cellH#, firstFrame#, frameCount#)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	flags, _ := args[1].ToInt()
	cw, okw := argDim(args[2])
	ch, okh := argDim(args[3])
	first, okf := args[4].ToInt()
	_, okc := args[5].ToInt()
	if !okw || !okh || !okf || !okc {
		return value.Nil, fmt.Errorf("LoadAnimTexture: numeric cell size / frames required")
	}
	if cw < 1 || ch < 1 {
		return value.Nil, fmt.Errorf("LoadAnimTexture: cell size must be >= 1")
	}
	src := rl.LoadImage(path)
	if src == nil || src.Data == nil {
		return value.Nil, fmt.Errorf("LoadAnimTexture: failed to load %q", path)
	}
	defer rl.UnloadImage(src)
	x0 := int(first) * cw
	if x0+cw > int(src.Width) || ch > int(src.Height) {
		return value.Nil, fmt.Errorf("LoadAnimTexture: frame rect outside image")
	}
	cp := rl.ImageCopy(src)
	defer rl.UnloadImage(cp)
	rl.ImageCrop(cp, rl.Rectangle{X: float32(x0), Y: 0, Width: float32(cw), Height: float32(ch)})
	t := rl.LoadTextureFromImage(cp)
	obj := &TextureObject{
		Tex: t, loaded: true, SourcePath: path, Flags: int32(flags), UScl: 1, VScl: 1,
	}
	obj.setFinalizer()
	texApplyLoadFlags(&obj.Tex, int32(flags))
	id, err2 := m.h.Alloc(obj)
	if err2 != nil {
		rl.UnloadTexture(t)
		return value.Nil, err2
	}
	return value.FromHandle(id), nil
}

func (m *Module) texNameStr(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := m.texObj1(args, "TextureName$")
	if err != nil {
		return value.Nil, err
	}
	o.mu.RLock()
	p := o.SourcePath
	o.mu.RUnlock()
	if p == "" {
		return rt.RetString(""), nil
	}
	return rt.RetString(p), nil
}

func (m *Module) texSetCubeFace(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SetCubeFace expects (texture, face#)")
	}
	o, err := m.texObj1([]value.Value{args[0]}, "SetCubeFace")
	if err != nil {
		return value.Nil, err
	}
	f, ok := args[1].ToInt()
	if !ok {
		return value.Nil, fmt.Errorf("face must be numeric")
	}
	o.mu.Lock()
	o.CubeFace = int32(f)
	o.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) texSetCubeMode(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SetCubeMode expects (texture, mode#)")
	}
	o, err := m.texObj1([]value.Value{args[0]}, "SetCubeMode")
	if err != nil {
		return value.Nil, err
	}
	md, ok := args[1].ToInt()
	if !ok {
		return value.Nil, fmt.Errorf("mode must be numeric")
	}
	o.mu.Lock()
	o.CubeMode = int32(md)
	o.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) texCoordsMode(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("TextureCoords expects (texture, coords#)")
	}
	o, err := m.texObj1([]value.Value{args[0]}, "TextureCoords")
	if err != nil {
		return value.Nil, err
	}
	c, ok := args[1].ToInt()
	if !ok {
		return value.Nil, fmt.Errorf("coords must be numeric")
	}
	o.mu.Lock()
	o.CoordsMode = int32(c)
	o.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) texScaleUV(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ScaleTexture expects (texture, uScale#, vScale#)")
	}
	o, err := m.texObj1([]value.Value{args[0]}, "ScaleTexture")
	if err != nil {
		return value.Nil, err
	}
	u, ok1 := argF32(args[1])
	v, ok2 := argF32(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("scales must be numeric")
	}
	o.mu.Lock()
	o.UScl, o.VScl = u, v
	o.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) texRotateUV(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("RotateTexture expects (texture, degrees#)")
	}
	o, err := m.texObj1([]value.Value{args[0]}, "RotateTexture")
	if err != nil {
		return value.Nil, err
	}
	a, ok := argF32(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("angle must be numeric")
	}
	o.mu.Lock()
	o.RotDeg = a
	o.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) texPositionUV(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("PositionTexture expects (texture, uPos#, vPos#)")
	}
	o, err := m.texObj1([]value.Value{args[0]}, "PositionTexture")
	if err != nil {
		return value.Nil, err
	}
	u, ok1 := argF32(args[1])
	v, ok2 := argF32(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("offsets must be numeric")
	}
	o.mu.Lock()
	o.UPos, o.VPos = u, v
	o.mu.Unlock()
	return value.Nil, nil
}
