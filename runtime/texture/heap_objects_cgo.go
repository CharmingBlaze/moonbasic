//go:build cgo

package texture

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

// TextureObject owns a GPU texture; no separate child handles.
type TextureObject struct {
	Tex     rl.Texture2D
	release heap.ReleaseOnce
}

func (t *TextureObject) TypeName() string { return "Texture" }
func (t *TextureObject) TypeTag() uint16  { return heap.TagTexture }
func (t *TextureObject) Free() {
	t.release.Do(func() { rl.UnloadTexture(t.Tex) })
}
