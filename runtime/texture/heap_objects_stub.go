//go:build !cgo

package texture

import "moonbasic/vm/heap"

// TextureObject is a stub when CGO is off; no texture heap objects are created (see load.go / gen_cgo.go tags).
type TextureObject struct{}

func (t *TextureObject) TypeName() string { return "Texture" }
func (t *TextureObject) TypeTag() uint16  { return heap.TagTexture }
func (t *TextureObject) Free()            {}

// RenderTargetObject is a stub when CGO is off.
type RenderTargetObject struct{}

func (r *RenderTargetObject) TypeName() string { return "RenderTexture" }
func (r *RenderTargetObject) TypeTag() uint16  { return heap.TagRenderTexture }
func (r *RenderTargetObject) Free()            {}
