//go:build cgo || (windows && !cgo)

package texture

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"moonbasic/vm/heap"
)

// ForBinding returns the Raylib texture for a TagTexture heap handle.
func ForBinding(store *heap.Store, h heap.Handle) (rl.Texture2D, error) {
	o, err := heap.Cast[*TextureObject](store, h)
	if err != nil {
		return rl.Texture2D{}, fmt.Errorf("not a texture handle: %w", err)
	}
	return o.Tex, nil
}
