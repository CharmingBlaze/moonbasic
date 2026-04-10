//go:build cgo || (windows && !cgo)

package mbimage

import (
	"errors"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

var errInvalidImage = errors.New("image handle has no backing image")

// imageObj owns a CPU-side Raylib image (heap-backed); no child handles.
type imageObj struct {
	img     *rl.Image
	release heap.ReleaseOnce
	// filterMode is Raylib TextureFilter (0 = unset). Applied when uploading via TEXTURE.FROMIMAGE / IMAGE.TOTEXTURE.
	filterMode int32
}

func (o *imageObj) TypeName() string { return "Image" }

func (o *imageObj) TypeTag() uint16 { return heap.TagImage }

func (o *imageObj) Free() {
	o.release.Do(func() {
		if o.img != nil {
			rl.UnloadImage(o.img)
			o.img = nil
		}
	})
}

// TextureFilterForHeapImage returns a Raylib TextureFilter to apply after LoadTextureFromImage, or 0 if unset.
func TextureFilterForHeapImage(s *heap.Store, h heap.Handle) int32 {
	if s == nil {
		return 0
	}
	o, err := heap.Cast[*imageObj](s, h)
	if err != nil {
		return 0
	}
	return o.filterMode
}

// RayImageForTexture returns the Raylib image for a TagImage heap handle (TEXTURE.FROMIMAGE).
func RayImageForTexture(s *heap.Store, h heap.Handle) (*rl.Image, error) {
	o, err := heap.Cast[*imageObj](s, h)
	if err != nil {
		return nil, err
	}
	if o.img == nil {
		return nil, errInvalidImage
	}
	return o.img, nil
}
