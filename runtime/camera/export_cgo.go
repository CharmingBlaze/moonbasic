//go:build cgo

package mbcamera

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

// RayCamera3D returns the Raylib camera for a CAMERA.MAKE heap handle.
func RayCamera3D(store *heap.Store, h heap.Handle) (rl.Camera3D, error) {
	if store == nil {
		return rl.Camera3D{}, fmt.Errorf("heap is nil")
	}
	o, err := heap.Cast[*camObj](store, h)
	if err != nil {
		return rl.Camera3D{}, err
	}
	return o.cam, nil
}
