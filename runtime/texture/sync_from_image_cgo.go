//go:build cgo || (windows && !cgo)

package texture

import (
	"fmt"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// SyncTextureFromImage replaces or updates the GPU texture backing a TextureObject from a CPU image.
// If dimensions differ, the old texture is unloaded and a new one is loaded (same flags as load path uses 1).
func SyncTextureFromImage(to *TextureObject, im *rl.Image) error {
	if to == nil {
		return fmt.Errorf("SyncTextureFromImage: nil texture object")
	}
	if im == nil || im.Data == nil || im.Width <= 0 || im.Height <= 0 {
		return fmt.Errorf("SyncTextureFromImage: invalid image")
	}
	to.mu.Lock()
	defer to.mu.Unlock()

	if int32(im.Width) != to.Tex.Width || int32(im.Height) != to.Tex.Height {
		if to.Borrowed {
			return fmt.Errorf("SyncTextureFromImage: cannot resize borrowed/render-target texture")
		}
		rl.UnloadTexture(to.Tex)
		to.Tex = rl.LoadTextureFromImage(im)
		texApplyLoadFlags(&to.Tex, to.Flags)
		return nil
	}
	sz := rl.GetPixelDataSize(im.Width, im.Height, int32(im.Format))
	if sz <= 0 {
		return fmt.Errorf("SyncTextureFromImage: unsupported image format")
	}
	pix := unsafe.Slice((*byte)(im.Data), sz)
	rl.UpdateTexture(to.Tex, pix)
	return nil
}
