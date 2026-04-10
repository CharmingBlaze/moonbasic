//go:build cgo || (windows && !cgo)

package mbimage

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerPixelFilterCmds(m *Module, r runtime.Registrar) {
	r.Register("IMAGE.GETPIXEL", "image", runtime.AdaptLegacy(m.imageGetPixel))
	r.Register("IMAGE.SETFILTER", "image", runtime.AdaptLegacy(m.imageSetFilter))
}

// imageGetPixel returns packed RGBA as a single int (0xAARRGGBB in host byte order; use bit ops in BASIC or MATH).
func (m *Module) imageGetPixel(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("IMAGE.GETPIXEL expects handle, x, y")
	}
	img, err := m.getImage(args, 0, "IMAGE.GETPIXEL")
	if err != nil {
		return value.Nil, err
	}
	xi, ok1 := argInt(args[1])
	yi, ok2 := argInt(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("IMAGE.GETPIXEL: x, y must be numeric")
	}
	c := rl.GetImageColor(*img, xi, yi)
	packed := int64(uint32(c.R)<<24 | uint32(c.G)<<16 | uint32(c.B)<<8 | uint32(c.A))
	return value.FromInt(packed), nil
}

func (m *Module) imageSetFilter(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("IMAGE.SETFILTER expects handle, filterMode# (Raylib TextureFilter int, e.g. TEXTURE_FILTER_POINT)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("IMAGE.SETFILTER: image handle required")
	}
	o, err := heap.Cast[*imageObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	fm, ok := args[1].ToInt()
	if !ok {
		if f, okf := args[1].ToFloat(); okf {
			fm = int64(f)
			ok = true
		}
	}
	if !ok {
		return value.Nil, fmt.Errorf("IMAGE.SETFILTER: filterMode must be numeric")
	}
	o.filterMode = int32(fm)
	return value.Nil, nil
}
