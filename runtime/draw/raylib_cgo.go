//go:build cgo

package mbdraw

import (
	"fmt"
	"image/color"

	"moonbasic/runtime"
	"moonbasic/runtime/convert"
	"moonbasic/runtime/mbimage"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type rlTex struct {
	t rl.Texture2D
}

func (r *rlTex) TypeName() string { return "Texture" }

func (r *rlTex) TypeTag() uint16 { return heap.TagTexture }

func (r *rlTex) Free() {
	rl.UnloadTexture(r.t)
}

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	r.Register("DRAW.RECTANGLE", "draw", runtime.AdaptLegacy(m.rectangle))
	r.Register("DRAW.RECTANGLE_ROUNDED", "draw", runtime.AdaptLegacy(m.rectangleRounded))
	r.Register("DRAW.TEXTURE", "draw", runtime.AdaptLegacy(m.drawTexture))
	m.registerTextureNPatch(r)
	r.Register("TEXTURE.LOAD", "draw", m.texLoad)
	r.Register("TEXTURE.FROMIMAGE", "draw", runtime.AdaptLegacy(m.texFromImage))
	r.Register("TEXTURE.FREE", "draw", runtime.AdaptLegacy(m.texFree))
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

func argInt(v value.Value) (int32, bool) {
	if i, ok := v.ToInt(); ok {
		return int32(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return int32(f), true
	}
	return 0, false
}

func (m *Module) rectangle(args []value.Value) (value.Value, error) {
	if len(args) != 8 {
		return value.Nil, fmt.Errorf("DRAW.RECTANGLE expects 8 arguments (x,y,w,h, r,g,b,a)")
	}
	var xywh [4]int32
	var rgb [4]int32
	for i := 0; i < 4; i++ {
		v, ok := argInt(args[i])
		if !ok {
			return value.Nil, fmt.Errorf("DRAW.RECTANGLE: non-numeric argument %d", i+1)
		}
		xywh[i] = v
	}
	for i := 0; i < 4; i++ {
		v, ok := argInt(args[4+i])
		if !ok {
			return value.Nil, fmt.Errorf("DRAW.RECTANGLE: non-numeric color argument %d", i+1)
		}
		rgb[i] = v
	}
	c := convert.NewColor4(rgb[0], rgb[1], rgb[2], rgb[3])
	col := color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
	rl.DrawRectangle(xywh[0], xywh[1], xywh[2], xywh[3], col)
	return value.Nil, nil
}

func (m *Module) rectangleRounded(args []value.Value) (value.Value, error) {
	if len(args) != 9 {
		return value.Nil, fmt.Errorf("DRAW.RECTANGLE_ROUNDED expects 9 arguments (x,y,w,h, radius, r,g,b,a)")
	}
	var xywh [4]int32
	for i := 0; i < 4; i++ {
		v, ok := argInt(args[i])
		if !ok {
			return value.Nil, fmt.Errorf("DRAW.RECTANGLE_ROUNDED: non-numeric argument %d", i+1)
		}
		xywh[i] = v
	}
	rad, ok := argInt(args[4])
	if !ok {
		return value.Nil, fmt.Errorf("DRAW.RECTANGLE_ROUNDED: radius must be numeric")
	}
	var rgb [4]int32
	for i := 0; i < 4; i++ {
		v, ok := argInt(args[5+i])
		if !ok {
			return value.Nil, fmt.Errorf("DRAW.RECTANGLE_ROUNDED: non-numeric color argument %d", i+1)
		}
		rgb[i] = v
	}
	c := convert.NewColor4(rgb[0], rgb[1], rgb[2], rgb[3])
	col := color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
	rl.DrawRectangleRounded(
		rl.Rectangle{X: float32(xywh[0]), Y: float32(xywh[1]), Width: float32(xywh[2]), Height: float32(xywh[3])},
		float32(rad),
		8,
		col,
	)
	return value.Nil, nil
}

func (m *Module) texLoad(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TEXTURE.LOAD: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("TEXTURE.LOAD expects 1 string path")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	t := rl.LoadTexture(path)
	id, err := m.h.Alloc(&rlTex{t: t})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) texFromImage(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TEXTURE.FROMIMAGE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TEXTURE.FROMIMAGE expects 1 image handle")
	}
	img, err := mbimage.RayImageForTexture(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("TEXTURE.FROMIMAGE: %w", err)
	}
	t := rl.LoadTextureFromImage(img)
	id, err := m.h.Alloc(&rlTex{t: t})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) texFree(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TEXTURE.FREE expects handle")
	}
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func (m *Module) drawTexture(args []value.Value) (value.Value, error) {
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("DRAW.TEXTURE expects 7 arguments (handle, x, y, r, g, b, a)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAW.TEXTURE: first argument must be texture handle")
	}
	obj, err := heap.Cast[*rlTex](m.h, heap.Handle(args[0].IVal))
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
	rl.DrawTexture(obj.t, x, y, tint)
	return value.Nil, nil
}

// TextureForBinding returns the Raylib texture for a TagTexture heap handle (materials, model setup).
func TextureForBinding(store *heap.Store, h heap.Handle) (rl.Texture2D, error) {
	o, err := heap.Cast[*rlTex](store, h)
	if err != nil {
		return rl.Texture2D{}, err
	}
	return o.t, nil
}
