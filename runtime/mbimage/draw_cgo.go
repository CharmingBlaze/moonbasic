//go:build cgo

package mbimage

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerImageDraw(m *Module, reg runtime.Registrar) {
	reg.Register("IMAGE.DRAWPIXEL", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 7 {
			return value.Nil, fmt.Errorf("IMAGE.DRAWPIXEL expects handle, x, y, r, g, b, a")
		}
		img, err := m.getImage(args, 0, "IMAGE.DRAWPIXEL")
		if err != nil {
			return value.Nil, err
		}
		x, ok1 := argInt(args[1])
		y, ok2 := argInt(args[2])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("IMAGE.DRAWPIXEL: x, y must be numeric")
		}
		col, err := rgbaFromArgs(args[3], args[4], args[5], args[6])
		if err != nil {
			return value.Nil, fmt.Errorf("IMAGE.DRAWPIXEL: %w", err)
		}
		rl.ImageDrawPixel(img, x, y, col)
		return value.Nil, nil
	}))

	reg.Register("IMAGE.DRAWRECT", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 9 {
			return value.Nil, fmt.Errorf("IMAGE.DRAWRECT expects handle, x, y, w, h, r, g, b, a")
		}
		img, err := m.getImage(args, 0, "IMAGE.DRAWRECT")
		if err != nil {
			return value.Nil, err
		}
		x, ok1 := argInt(args[1])
		y, ok2 := argInt(args[2])
		w, ok3 := argInt(args[3])
		h, ok4 := argInt(args[4])
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return value.Nil, fmt.Errorf("IMAGE.DRAWRECT: x, y, w, h must be numeric")
		}
		col, err := rgbaFromArgs(args[5], args[6], args[7], args[8])
		if err != nil {
			return value.Nil, fmt.Errorf("IMAGE.DRAWRECT: %w", err)
		}
		rl.ImageDrawRectangle(img, x, y, w, h, col)
		return value.Nil, nil
	}))

	reg.Register("IMAGE.DRAWLINE", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 9 {
			return value.Nil, fmt.Errorf("IMAGE.DRAWLINE expects handle, x1, y1, x2, y2, r, g, b, a")
		}
		img, err := m.getImage(args, 0, "IMAGE.DRAWLINE")
		if err != nil {
			return value.Nil, err
		}
		x1, ok1 := argInt(args[1])
		y1, ok2 := argInt(args[2])
		x2, ok3 := argInt(args[3])
		y2, ok4 := argInt(args[4])
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return value.Nil, fmt.Errorf("IMAGE.DRAWLINE: coordinates must be numeric")
		}
		col, err := rgbaFromArgs(args[5], args[6], args[7], args[8])
		if err != nil {
			return value.Nil, fmt.Errorf("IMAGE.DRAWLINE: %w", err)
		}
		rl.ImageDrawLine(img, x1, y1, x2, y2, col)
		return value.Nil, nil
	}))

	reg.Register("IMAGE.DRAWCIRCLE", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 8 {
			return value.Nil, fmt.Errorf("IMAGE.DRAWCIRCLE expects handle, cx, cy, radius, r, g, b, a")
		}
		img, err := m.getImage(args, 0, "IMAGE.DRAWCIRCLE")
		if err != nil {
			return value.Nil, err
		}
		cx, ok1 := argInt(args[1])
		cy, ok2 := argInt(args[2])
		rad, ok3 := argInt(args[3])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("IMAGE.DRAWCIRCLE: center and radius must be numeric")
		}
		col, err := rgbaFromArgs(args[4], args[5], args[6], args[7])
		if err != nil {
			return value.Nil, fmt.Errorf("IMAGE.DRAWCIRCLE: %w", err)
		}
		rl.ImageDrawCircle(img, cx, cy, rad, col)
		return value.Nil, nil
	}))

	reg.Register("IMAGE.DRAWTEXT", "image", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 9 || args[3].Kind != value.KindString {
			return value.Nil, fmt.Errorf("IMAGE.DRAWTEXT expects handle, x, y, text$, fontSize, r, g, b, a")
		}
		img, err := m.getImage(args, 0, "IMAGE.DRAWTEXT")
		if err != nil {
			return value.Nil, err
		}
		x, ok1 := argInt(args[1])
		y, ok2 := argInt(args[2])
		fs, ok3 := argInt(args[4])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("IMAGE.DRAWTEXT: x, y, fontSize must be numeric")
		}
		col, err := rgbaFromArgs(args[5], args[6], args[7], args[8])
		if err != nil {
			return value.Nil, fmt.Errorf("IMAGE.DRAWTEXT: %w", err)
		}
		text, err := rt.ArgString(args, 3)
		if err != nil {
			return value.Nil, err
		}
		rl.ImageDrawText(img, x, y, text, fs, col)
		return value.Nil, nil
	})
}
