package mbmatrix

import (
	"fmt"
	"image/color"

	"moonbasic/vm/heap"
)

// colorObj is a heap-wrapped RGBA (Raylib Color); used with and without CGO.
type colorObj struct {
	c color.RGBA
}

func (o *colorObj) TypeName() string { return "Color" }

func (o *colorObj) TypeTag() uint16 { return heap.TagColor }

func (o *colorObj) Free() {}

// HeapColorRGBA returns the RGBA for a heap Color handle (TagColor).
func HeapColorRGBA(s *heap.Store, h heap.Handle) (color.RGBA, error) {
	if s == nil || h == 0 {
		return color.RGBA{}, fmt.Errorf("invalid color handle")
	}
	o, err := heap.Cast[*colorObj](s, h)
	if err != nil {
		return color.RGBA{}, err
	}
	return o.c, nil
}
// GetColor is an alias for HeapColorRGBA for cross-module usage.
func GetColor(s *heap.Store, h heap.Handle) (color.RGBA, error) {
	return HeapColorRGBA(s, h)
}
