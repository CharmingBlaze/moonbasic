//go:build cgo

package mbimage

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerImageProcess(m *Module, reg runtime.Registrar) {
	reg.Register("IMAGE.DRAWIMAGE", "image", runtime.AdaptLegacy(m.imageDrawImage))
	reg.Register("IMAGE.DITHER", "image", runtime.AdaptLegacy(m.imageDither))
	reg.Register("IMAGE.MIPMAPS", "image", runtime.AdaptLegacy(m.imageMipmaps))
	reg.Register("IMAGE.FORMAT", "image", runtime.AdaptLegacy(m.imageFormat))
	reg.Register("IMAGE.DRAWRECTLINES", "image", runtime.AdaptLegacy(m.imageDrawRectLines))
	reg.Register("IMAGE.ALPHACROP", "image", runtime.AdaptLegacy(m.imageAlphaCrop))
	reg.Register("IMAGE.ALPHACLEAR", "image", runtime.AdaptLegacy(m.imageAlphaClear))
}

// IMAGE.DRAWIMAGE: blit src into dest with source and dest rectangles and tint (Raylib ImageDraw).
func (m *Module) imageDrawImage(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 14 {
		return value.Nil, fmt.Errorf("IMAGE.DRAWIMAGE expects (dest, src, sx,sy,sw,sh, dx,dy,dw,dh, r,g,b,a)")
	}
	dst, err := m.getImage(args, 0, "IMAGE.DRAWIMAGE")
	if err != nil {
		return value.Nil, err
	}
	src, err := m.getImage(args, 1, "IMAGE.DRAWIMAGE")
	if err != nil {
		return value.Nil, err
	}
	sx, ok1 := argFloat(args[2])
	sy, ok2 := argFloat(args[3])
	sw, ok3 := argFloat(args[4])
	sh, ok4 := argFloat(args[5])
	dx, ok5 := argFloat(args[6])
	dy, ok6 := argFloat(args[7])
	dw, ok7 := argFloat(args[8])
	dh, ok8 := argFloat(args[9])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 {
		return value.Nil, fmt.Errorf("IMAGE.DRAWIMAGE: rectangle components must be numeric")
	}
	tint, err := rgbaFromArgs(args[10], args[11], args[12], args[13])
	if err != nil {
		return value.Nil, fmt.Errorf("IMAGE.DRAWIMAGE: %w", err)
	}
	srcRec := rl.Rectangle{X: sx, Y: sy, Width: sw, Height: sh}
	dstRec := rl.Rectangle{X: dx, Y: dy, Width: dw, Height: dh}
	rl.ImageDraw(dst, src, srcRec, dstRec, tint)
	return value.Nil, nil
}

func (m *Module) imageDither(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("IMAGE.DITHER expects (handle, rBpp, gBpp, bBpp, aBpp)")
	}
	img, err := m.getImage(args, 0, "IMAGE.DITHER")
	if err != nil {
		return value.Nil, err
	}
	var bpp [4]int32
	for i := 0; i < 4; i++ {
		v, ok := argInt(args[1+i])
		if !ok {
			return value.Nil, fmt.Errorf("IMAGE.DITHER: bpp values must be numeric")
		}
		bpp[i] = v
	}
	rl.ImageDither(img, bpp[0], bpp[1], bpp[2], bpp[3])
	return value.Nil, nil
}

func (m *Module) imageMipmaps(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("IMAGE.MIPMAPS expects image handle")
	}
	img, err := m.getImage(args, 0, "IMAGE.MIPMAPS")
	if err != nil {
		return value.Nil, err
	}
	rl.ImageMipmaps(img)
	return value.Nil, nil
}

func (m *Module) imageFormat(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("IMAGE.FORMAT expects (handle, newFormat) Raylib PixelFormat int")
	}
	img, err := m.getImage(args, 0, "IMAGE.FORMAT")
	if err != nil {
		return value.Nil, err
	}
	f, ok := argInt(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("IMAGE.FORMAT: format must be numeric")
	}
	rl.ImageFormat(img, rl.PixelFormat(f))
	return value.Nil, nil
}

func (m *Module) imageDrawRectLines(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 10 {
		return value.Nil, fmt.Errorf("IMAGE.DRAWRECTLINES expects (handle, x, y, w, h, thick, r, g, b, a)")
	}
	img, err := m.getImage(args, 0, "IMAGE.DRAWRECTLINES")
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argFloat(args[1])
	y, ok2 := argFloat(args[2])
	w, ok3 := argFloat(args[3])
	h, ok4 := argFloat(args[4])
	thick, ok5 := argInt(args[5])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("IMAGE.DRAWRECTLINES: geometry must be numeric")
	}
	tint, err := rgbaFromArgs(args[6], args[7], args[8], args[9])
	if err != nil {
		return value.Nil, fmt.Errorf("IMAGE.DRAWRECTLINES: %w", err)
	}
	rec := rl.Rectangle{X: x, Y: y, Width: w, Height: h}
	rl.ImageDrawRectangleLines(img, rec, int(thick), tint)
	return value.Nil, nil
}

func (m *Module) imageAlphaCrop(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("IMAGE.ALPHACROP expects (handle, threshold#)")
	}
	img, err := m.getImage(args, 0, "IMAGE.ALPHACROP")
	if err != nil {
		return value.Nil, err
	}
	th, ok := argFloat(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("IMAGE.ALPHACROP: threshold must be numeric")
	}
	rl.ImageAlphaCrop(img, th)
	return value.Nil, nil
}

func (m *Module) imageAlphaClear(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("IMAGE.ALPHACLEAR expects (handle, r, g, b, a, threshold#)")
	}
	img, err := m.getImage(args, 0, "IMAGE.ALPHACLEAR")
	if err != nil {
		return value.Nil, err
	}
	c, err := rgbaFromArgs(args[1], args[2], args[3], args[4])
	if err != nil {
		return value.Nil, fmt.Errorf("IMAGE.ALPHACLEAR: %w", err)
	}
	th, ok := argFloat(args[5])
	if !ok {
		return value.Nil, fmt.Errorf("IMAGE.ALPHACLEAR: threshold must be numeric")
	}
	rl.ImageAlphaClear(img, c, th)
	return value.Nil, nil
}
