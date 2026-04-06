//go:build cgo || (windows && !cgo)

package mbimage

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// imageAlphaBBox returns a tight rectangle around pixels with alpha/255 > threshold.
// raylib-go's CGO build does not wrap GetImageAlphaBorder; this matches its 0–1 threshold semantics.
func imageAlphaBBox(img *rl.Image, threshold float32) rl.Rectangle {
	if img == nil || img.Width <= 0 || img.Height <= 0 {
		return rl.Rectangle{}
	}
	w, h := img.Width, img.Height
	cols := rl.LoadImageColors(img)
	defer rl.UnloadImageColors(cols)
	minX, minY := w, h
	maxX, maxY := int32(-1), int32(-1)
	for y := int32(0); y < h; y++ {
		row := y * w
		for x := int32(0); x < w; x++ {
			a := cols[row+x].A
			if float32(a)/255.0 > threshold {
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}
	if maxX < 0 {
		return rl.Rectangle{}
	}
	return rl.Rectangle{
		X:      float32(minX),
		Y:      float32(minY),
		Width:  float32(maxX - minX + 1),
		Height: float32(maxY - minY + 1),
	}
}
