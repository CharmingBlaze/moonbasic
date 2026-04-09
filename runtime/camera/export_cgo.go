//go:build cgo || (windows && !cgo)

package mbcamera

import (
	"fmt"
	"math"

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

// CameraXZStrafeBasis returns unit forward (camera→target flattened to XZ) and right for camera-relative movement.
func CameraXZStrafeBasis(store *heap.Store, ch heap.Handle) (fwd rl.Vector3, right rl.Vector3, err error) {
	cam, err := RayCamera3D(store, ch)
	if err != nil {
		return rl.Vector3{}, rl.Vector3{}, err
	}
	dx := float64(cam.Target.X - cam.Position.X)
	dz := float64(cam.Target.Z - cam.Position.Z)
	h := math.Hypot(dx, dz)
	if h < 1e-5 {
		return rl.Vector3{X: 0, Y: 0, Z: 1}, rl.Vector3{X: 1, Y: 0, Z: 0}, nil
	}
	fwd = rl.Vector3{X: float32(dx / h), Y: 0, Z: float32(dz / h)}
	// Horizontal right = negated (Y × forward) so strafe matches SetOrbit yaw and main.mb-style WASD.
	right = rl.Vector3{X: -fwd.Z, Y: 0, Z: fwd.X}
	return fwd, right, nil
}
