//go:build cgo

package mbimage

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerImageTransform(m *Module, reg runtime.Registrar) {
	reg.Register("IMAGE.CROP", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("IMAGE.CROP expects handle, x, y, width, height")
		}
		img, err := m.getImage(args, 0, "IMAGE.CROP")
		if err != nil {
			return value.Nil, err
		}
		x, ok1 := argInt(args[1])
		y, ok2 := argInt(args[2])
		w, ok3 := argInt(args[3])
		h, ok4 := argInt(args[4])
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return value.Nil, fmt.Errorf("IMAGE.CROP: x, y, width, height must be numeric")
		}
		rl.ImageCrop(img, rl.Rectangle{X: float32(x), Y: float32(y), Width: float32(w), Height: float32(h)})
		return value.Nil, nil
	}))

	reg.Register("IMAGE.RESIZE", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("IMAGE.RESIZE expects handle, newWidth, newHeight")
		}
		img, err := m.getImage(args, 0, "IMAGE.RESIZE")
		if err != nil {
			return value.Nil, err
		}
		w, ok1 := argInt(args[1])
		h, ok2 := argInt(args[2])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("IMAGE.RESIZE: dimensions must be numeric")
		}
		rl.ImageResize(img, w, h)
		return value.Nil, nil
	}))

	reg.Register("IMAGE.RESIZENN", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("IMAGE.RESIZENN expects handle, newWidth, newHeight")
		}
		img, err := m.getImage(args, 0, "IMAGE.RESIZENN")
		if err != nil {
			return value.Nil, err
		}
		w, ok1 := argInt(args[1])
		h, ok2 := argInt(args[2])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("IMAGE.RESIZENN: dimensions must be numeric")
		}
		rl.ImageResizeNN(img, w, h)
		return value.Nil, nil
	}))

	reg.Register("IMAGE.FLIPH", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("IMAGE.FLIPH expects image handle")
		}
		img, err := m.getImage(args, 0, "IMAGE.FLIPH")
		if err != nil {
			return value.Nil, err
		}
		rl.ImageFlipHorizontal(img)
		return value.Nil, nil
	}))

	reg.Register("IMAGE.FLIPV", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("IMAGE.FLIPV expects image handle")
		}
		img, err := m.getImage(args, 0, "IMAGE.FLIPV")
		if err != nil {
			return value.Nil, err
		}
		rl.ImageFlipVertical(img)
		return value.Nil, nil
	}))

	reg.Register("IMAGE.ROTATE", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("IMAGE.ROTATE expects handle, degrees")
		}
		img, err := m.getImage(args, 0, "IMAGE.ROTATE")
		if err != nil {
			return value.Nil, err
		}
		deg, ok := argInt(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("IMAGE.ROTATE: degrees must be numeric")
		}
		rl.ImageRotate(img, deg)
		return value.Nil, nil
	}))

	reg.Register("IMAGE.ROTATECW", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("IMAGE.ROTATECW expects image handle")
		}
		img, err := m.getImage(args, 0, "IMAGE.ROTATECW")
		if err != nil {
			return value.Nil, err
		}
		rl.ImageRotateCW(img)
		return value.Nil, nil
	}))

	reg.Register("IMAGE.ROTATECCW", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("IMAGE.ROTATECCW expects image handle")
		}
		img, err := m.getImage(args, 0, "IMAGE.ROTATECCW")
		if err != nil {
			return value.Nil, err
		}
		rl.ImageRotateCCW(img)
		return value.Nil, nil
	}))

	reg.Register("IMAGE.COLORTINT", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("IMAGE.COLORTINT expects handle, r, g, b, a")
		}
		img, err := m.getImage(args, 0, "IMAGE.COLORTINT")
		if err != nil {
			return value.Nil, err
		}
		col, err := rgbaFromArgs(args[1], args[2], args[3], args[4])
		if err != nil {
			return value.Nil, fmt.Errorf("IMAGE.COLORTINT: %w", err)
		}
		rl.ImageColorTint(img, col)
		return value.Nil, nil
	}))

	reg.Register("IMAGE.COLORINVERT", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("IMAGE.COLORINVERT expects image handle")
		}
		img, err := m.getImage(args, 0, "IMAGE.COLORINVERT")
		if err != nil {
			return value.Nil, err
		}
		rl.ImageColorInvert(img)
		return value.Nil, nil
	}))

	reg.Register("IMAGE.COLORGRAYSCALE", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("IMAGE.COLORGRAYSCALE expects image handle")
		}
		img, err := m.getImage(args, 0, "IMAGE.COLORGRAYSCALE")
		if err != nil {
			return value.Nil, err
		}
		rl.ImageColorGrayscale(img)
		return value.Nil, nil
	}))

	reg.Register("IMAGE.COLORCONTRAST", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("IMAGE.COLORCONTRAST expects handle, contrast")
		}
		img, err := m.getImage(args, 0, "IMAGE.COLORCONTRAST")
		if err != nil {
			return value.Nil, err
		}
		c, ok := argFloat(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("IMAGE.COLORCONTRAST: contrast must be numeric")
		}
		rl.ImageColorContrast(img, c)
		return value.Nil, nil
	}))

	reg.Register("IMAGE.COLORBRIGHTNESS", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("IMAGE.COLORBRIGHTNESS expects handle, brightness")
		}
		img, err := m.getImage(args, 0, "IMAGE.COLORBRIGHTNESS")
		if err != nil {
			return value.Nil, err
		}
		b, ok := argInt(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("IMAGE.COLORBRIGHTNESS: brightness must be numeric")
		}
		rl.ImageColorBrightness(img, b)
		return value.Nil, nil
	}))

	reg.Register("IMAGE.COLORREPLACE", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 9 {
			return value.Nil, fmt.Errorf("IMAGE.COLORREPLACE expects handle, fromR,fromG,fromB,fromA, toR,toG,toB,toA")
		}
		img, err := m.getImage(args, 0, "IMAGE.COLORREPLACE")
		if err != nil {
			return value.Nil, err
		}
		from, err := rgbaFromArgs(args[1], args[2], args[3], args[4])
		if err != nil {
			return value.Nil, fmt.Errorf("IMAGE.COLORREPLACE (from): %w", err)
		}
		to, err := rgbaFromArgs(args[5], args[6], args[7], args[8])
		if err != nil {
			return value.Nil, fmt.Errorf("IMAGE.COLORREPLACE (to): %w", err)
		}
		rl.ImageColorReplace(img, from, to)
		return value.Nil, nil
	}))

	imageClearBackground := func(op string) func([]value.Value) (value.Value, error) {
		return func(args []value.Value) (value.Value, error) {
			if err := m.requireHeap(); err != nil {
				return value.Nil, err
			}
			if len(args) != 5 {
				return value.Nil, fmt.Errorf("%s expects handle, r, g, b, a", op)
			}
			img, err := m.getImage(args, 0, op)
			if err != nil {
				return value.Nil, err
			}
			col, err := rgbaFromArgs(args[1], args[2], args[3], args[4])
			if err != nil {
				return value.Nil, fmt.Errorf("%s: %w", op, err)
			}
			rl.ImageClearBackground(img, col)
			return value.Nil, nil
		}
	}
	reg.Register("IMAGE.CLEARBACKGROUND", "image", runtime.AdaptLegacy(imageClearBackground("IMAGE.CLEARBACKGROUND")))
	reg.Register("IMAGE.CLEAR", "image", runtime.AdaptLegacy(imageClearBackground("IMAGE.CLEAR")))
}
