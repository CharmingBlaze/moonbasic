//go:build !cgo

package mbcamera

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"moonbasic/vm/heap"
)

// RayCamera3D is unavailable without CGO (camera module is stubbed).
func RayCamera3D(store *heap.Store, h heap.Handle) (rl.Camera3D, error) {
	return rl.Camera3D{}, fmt.Errorf("RayCamera3D requires CGO-enabled build")
}
