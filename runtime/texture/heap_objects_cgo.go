//go:build cgo

package texture

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

// TextureObject owns a GPU texture; no separate child handles.
// Borrowed is true for views created by RENDERTARGET.TEXTURE — do not unload the GPU texture here.
type TextureObject struct {
	Tex      rl.Texture2D
	Borrowed bool
	release  heap.ReleaseOnce
}

func (t *TextureObject) TypeName() string { return "Texture" }
func (t *TextureObject) TypeTag() uint16  { return heap.TagTexture }
func (t *TextureObject) Free() {
	if t.Borrowed {
		return
	}
	t.release.Do(func() { rl.UnloadTexture(t.Tex) })
}

// RenderTargetObject owns a Raylib render target (FBO + color/depth attachments).
type RenderTargetObject struct {
	RT rl.RenderTexture2D
}

func (r *RenderTargetObject) TypeName() string { return "RenderTexture" }
func (r *RenderTargetObject) TypeTag() uint16   { return heap.TagRenderTexture }
func (r *RenderTargetObject) Free() {
	rl.UnloadRenderTexture(r.RT)
}
