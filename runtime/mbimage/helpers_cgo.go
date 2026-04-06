//go:build cgo || (windows && !cgo)

package mbimage

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/convert"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) requireHeap() error {
	if m.h == nil {
		return runtime.Errorf("IMAGE.*: heap not bound")
	}
	return nil
}

func argInt(v value.Value) (int32, bool) {
	if i, ok := v.ToInt(); ok {
		return int32(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return int32(f), true
	}
	return 0, false
}

func argFloat(v value.Value) (float32, bool) {
	if f, ok := v.ToFloat(); ok {
		return float32(f), true
	}
	if i, ok := v.ToInt(); ok {
		return float32(i), true
	}
	return 0, false
}

func rgbaFromArgs(r, g, b, a value.Value) (color.RGBA, error) {
	ri, ok1 := argInt(r)
	gi, ok2 := argInt(g)
	bi, ok3 := argInt(b)
	ai, ok4 := argInt(a)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return color.RGBA{}, fmt.Errorf("expected numeric RGBA components")
	}
	c := convert.NewColor4(ri, gi, bi, ai)
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}, nil
}

func (m *Module) getImage(args []value.Value, ix int, op string) (*rl.Image, error) {
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: argument %d must be image handle", op, ix+1)
	}
	o, err := heap.Cast[*imageObj](m.h, heap.Handle(args[ix].IVal))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if o.img == nil {
		return nil, fmt.Errorf("%s: %v", op, errInvalidImage)
	}
	return o.img, nil
}

func (m *Module) allocImage(img *rl.Image, op string) (value.Value, error) {
	if img == nil || !rl.IsImageValid(img) {
		return value.Nil, fmt.Errorf("%s: invalid image", op)
	}
	id, err := m.h.Alloc(&imageObj{img: img})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

// ImageFromArgs retrieves a Raylib image from a heap handle.
func ImageFromArgs(s *heap.Store, h heap.Handle) (*rl.Image, error) {
	o, err := heap.Cast[*imageObj](s, h)
	if err != nil {
		return nil, err
	}
	return o.img, nil
}
