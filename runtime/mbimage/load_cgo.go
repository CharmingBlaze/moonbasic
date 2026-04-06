//go:build cgo || (windows && !cgo)

package mbimage

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerImageLoad(m *Module, reg runtime.Registrar) {
	reg.Register("IMAGE.LOAD", "image", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 || args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("IMAGE.LOAD expects 1 string path")
		}
		path, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		img := rl.LoadImage(path)
		return m.allocImage(img, "IMAGE.LOAD")
	})

	reg.Register("IMAGE.LOADRAW", "image", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 5 || args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("IMAGE.LOADRAW expects path$, width, height, format, headerSize (Raylib PixelFormat int)")
		}
		w, ok1 := argInt(args[1])
		h, ok2 := argInt(args[2])
		fmtVal, ok3 := argInt(args[3])
		hdr, ok4 := argInt(args[4])
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return value.Nil, fmt.Errorf("IMAGE.LOADRAW: width, height, format, headerSize must be numeric")
		}
		path, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		img := rl.LoadImageRaw(path, w, h, rl.PixelFormat(fmtVal), hdr)
		return m.allocImage(img, "IMAGE.LOADRAW")
	})

	// makeBlank accepts (w, h) → transparent black RGBA(0,0,0,0), or (w,h,r,g,b,a) → filled color.
	makeBlank := func(op string) func([]value.Value) (value.Value, error) {
		return func(args []value.Value) (value.Value, error) {
			if err := m.requireHeap(); err != nil {
				return value.Nil, err
			}
			if len(args) != 2 && len(args) != 6 {
				return value.Nil, fmt.Errorf("%s expects (width, height) or (width, height, r, g, b, a)", op)
			}
			w, ok1 := argInt(args[0])
			h, ok2 := argInt(args[1])
			if !ok1 || !ok2 {
				return value.Nil, fmt.Errorf("%s: width and height must be numeric", op)
			}
			var col rl.Color
			if len(args) == 2 {
				col = rl.Color{R: 0, G: 0, B: 0, A: 0}
			} else {
				c, err := rgbaFromArgs(args[2], args[3], args[4], args[5])
				if err != nil {
					return value.Nil, fmt.Errorf("%s: %w", op, err)
				}
				col = rl.Color{R: c.R, G: c.G, B: c.B, A: c.A}
			}
			img := rl.GenImageColor(int(w), int(h), col)
			return m.allocImage(img, op)
		}
	}
	reg.Register("IMAGE.MAKEBLANK", "image", runtime.AdaptLegacy(makeBlank("IMAGE.MAKEBLANK")))
	reg.Register("IMAGE.MAKE", "image", runtime.AdaptLegacy(makeBlank("IMAGE.MAKE")))

	makeCopy := func(op string) func([]value.Value) (value.Value, error) {
		return func(args []value.Value) (value.Value, error) {
			if err := m.requireHeap(); err != nil {
				return value.Nil, err
			}
			if len(args) != 1 {
				return value.Nil, fmt.Errorf("%s expects image handle", op)
			}
			src, err := m.getImage(args, 0, op)
			if err != nil {
				return value.Nil, err
			}
			cp := rl.ImageCopy(src)
			return m.allocImage(cp, op)
		}
	}
	reg.Register("IMAGE.MAKECOPY", "image", runtime.AdaptLegacy(makeCopy("IMAGE.MAKECOPY")))
	reg.Register("IMAGE.COPY", "image", runtime.AdaptLegacy(makeCopy("IMAGE.COPY")))

	// Fixed canvas: GenImageColor + ImageDrawText (not raw GenImageText).
	reg.Register("IMAGE.MAKETEXT", "image", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 6 || args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("IMAGE.MAKETEXT expects text$, fontSize, r, g, b, a")
		}
		fs, ok := argInt(args[1])
		if !ok || fs < 1 {
			return value.Nil, fmt.Errorf("IMAGE.MAKETEXT: fontSize must be a positive number")
		}
		fg, err := rgbaFromArgs(args[2], args[3], args[4], args[5])
		if err != nil {
			return value.Nil, fmt.Errorf("IMAGE.MAKETEXT: %w", err)
		}
		text, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		// Canvas: width from text length × font size (capped), height fontSize + padding.
		w := int32(len(text)) * fs
		if w < 64 {
			w = 64
		}
		if w > 2048 {
			w = 2048
		}
		h := fs + 16
		if h > 512 {
			h = 512
		}
		bg := color.RGBA{R: 0, G: 0, B: 0, A: 0}
		img := rl.GenImageColor(int(w), int(h), bg)
		if img == nil || !rl.IsImageValid(img) {
			return value.Nil, fmt.Errorf("IMAGE.MAKETEXT: could not allocate canvas")
		}
		rl.ImageDrawText(img, 4, 4, text, fs, fg)
		return m.allocImage(img, "IMAGE.MAKETEXT")
	})

	reg.Register("IMAGE.EXPORT", "image", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 || args[1].Kind != value.KindString {
			return value.Nil, fmt.Errorf("IMAGE.EXPORT expects image handle and path$")
		}
		img, err := m.getImage(args, 0, "IMAGE.EXPORT")
		if err != nil {
			return value.Nil, err
		}
		path, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Nil, err
		}
		ok := rl.ExportImage(*img, path)
		return value.FromBool(ok), nil
	})

	reg.Register("IMAGE.WIDTH", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("IMAGE.WIDTH expects image handle")
		}
		img, err := m.getImage(args, 0, "IMAGE.WIDTH")
		if err != nil {
			return value.Nil, err
		}
		return value.FromInt(int64(img.Width)), nil
	}))

	reg.Register("IMAGE.HEIGHT", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("IMAGE.HEIGHT expects image handle")
		}
		img, err := m.getImage(args, 0, "IMAGE.HEIGHT")
		if err != nil {
			return value.Nil, err
		}
		return value.FromInt(int64(img.Height)), nil
	}))

	reg.Register("IMAGE.FREE", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 || args[0].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("IMAGE.FREE expects image handle")
		}
		if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
			return value.Nil, err
		}
		return value.Nil, nil
	}))
}
