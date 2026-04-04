//go:build cgo

package texture

import (
	"fmt"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/mbimage"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerTexturePropCmds(m *Module, r runtime.Registrar) {
	r.Register("TEXTURE.WIDTH", "texture", runtime.AdaptLegacy(m.texWidth))
	r.Register("TEXTURE.HEIGHT", "texture", runtime.AdaptLegacy(m.texHeight))
	r.Register("TEXTURE.SETFILTER", "texture", runtime.AdaptLegacy(m.texSetFilter))
	r.Register("TEXTURE.SETWRAP", "texture", runtime.AdaptLegacy(m.texSetWrap))
	r.Register("TEXTURE.UPDATE", "texture", runtime.AdaptLegacy(m.texUpdate))
}

func (m *Module) texWidth(args []value.Value) (value.Value, error) {
	t, err := m.tex2D(args, "TEXTURE.WIDTH")
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(int64(t.Width)), nil
}

func (m *Module) texHeight(args []value.Value) (value.Value, error) {
	t, err := m.tex2D(args, "TEXTURE.HEIGHT")
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(int64(t.Height)), nil
}

func (m *Module) texSetFilter(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TEXTURE.SETFILTER expects (texture, filterMode)")
	}
	t, err := ForBinding(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	f, ok := argTexFilter(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("TEXTURE.SETFILTER: filter must be numeric (use FILTER_POINT etc.)")
	}
	rl.SetTextureFilter(t, f)
	return value.Nil, nil
}

func (m *Module) texSetWrap(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TEXTURE.SETWRAP expects (texture, wrapMode)")
	}
	t, err := ForBinding(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	w, ok := argTexWrap(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("TEXTURE.SETWRAP: wrap must be numeric (use WRAP_REPEAT etc.)")
	}
	rl.SetTextureWrap(t, w)
	return value.Nil, nil
}

func (m *Module) texUpdate(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TEXTURE.UPDATE expects (texture, image)")
	}
	t, err := ForBinding(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	img, err := mbimage.RayImageForTexture(m.h, heap.Handle(args[1].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("TEXTURE.UPDATE: %w", err)
	}
	if img.Data == nil || img.Width <= 0 || img.Height <= 0 {
		return value.Nil, fmt.Errorf("TEXTURE.UPDATE: image has no pixel data")
	}
	sz := rl.GetPixelDataSize(img.Width, img.Height, int32(img.Format))
	if sz <= 0 {
		return value.Nil, fmt.Errorf("TEXTURE.UPDATE: unsupported image format")
	}
	pix := unsafe.Slice((*byte)(img.Data), sz)
	rl.UpdateTexture(t, pix)
	return value.Nil, nil
}

func (m *Module) tex2D(args []value.Value, op string) (rl.Texture2D, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return rl.Texture2D{}, fmt.Errorf("%s expects texture handle", op)
	}
	return ForBinding(m.h, heap.Handle(args[0].IVal))
}

func argTexFilter(v value.Value) (rl.TextureFilterMode, bool) {
	if i, ok := v.ToInt(); ok {
		return rl.TextureFilterMode(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return rl.TextureFilterMode(f), true
	}
	return 0, false
}

func argTexWrap(v value.Value) (rl.TextureWrapMode, bool) {
	if i, ok := v.ToInt(); ok {
		return rl.TextureWrapMode(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return rl.TextureWrapMode(f), true
	}
	return 0, false
}
