//go:build cgo || (windows && !cgo)

package texture

import (
	"runtime"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

// TextureObject owns a GPU texture; no separate child handles.
// Borrowed is true for views created by RENDERTARGET.TEXTURE — do not unload the GPU texture here.
type TextureObject struct {
	Tex      rl.Texture2D
	Borrowed bool

	// Blitz-style metadata (UV/cube — honored by materials that read these fields)
	SourcePath string
	Flags      int32
	UScl       float32
	VScl       float32
	UPos       float32
	VPos       float32
	RotDeg     float32
	CubeFace   int32
	CubeMode   int32
	CoordsMode int32

	// Asynchronous state
	mu        sync.RWMutex
	isLoading bool
	loaded    bool
	loadError string

	release heap.ReleaseOnce
}

func (t *TextureObject) TypeName() string { return "Texture" }
func (t *TextureObject) TypeTag() uint16  { return heap.TagTexture }
func (t *TextureObject) Free() {
	if t.Borrowed {
		return
	}
	t.release.Do(func() { rl.UnloadTexture(t.Tex) })
}

func (t *TextureObject) setFinalizer() {
	runtime.SetFinalizer(t, func(o *TextureObject) {
		enqueueOnMainThread(func() { o.Free() })
	})
}

// RenderTargetObject owns a Raylib render target (FBO + color/depth attachments).
type RenderTargetObject struct {
	RT      rl.RenderTexture2D
	release heap.ReleaseOnce
}

func (r *RenderTargetObject) TypeName() string { return "RenderTexture" }
func (r *RenderTargetObject) TypeTag() uint16   { return heap.TagRenderTexture }
func (r *RenderTargetObject) Free() {
	r.release.Do(func() { rl.UnloadRenderTexture(r.RT) })
}

func (r *RenderTargetObject) setFinalizer() {
	runtime.SetFinalizer(r, func(o *RenderTargetObject) {
		enqueueOnMainThread(func() { o.Free() })
	})
}
