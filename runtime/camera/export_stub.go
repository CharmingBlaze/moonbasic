//go:build !cgo && !windows

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

// CameraXZStrafeBasis is unavailable without CGO.
func CameraXZStrafeBasis(store *heap.Store, ch heap.Handle) (fwd rl.Vector3, right rl.Vector3, err error) {
	return rl.Vector3{}, rl.Vector3{}, fmt.Errorf("CameraXZStrafeBasis requires CGO-enabled build")
}
