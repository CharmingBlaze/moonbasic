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

// CameraXZWalkBasis returns horizontal forward/right for camera-relative walking.
// When orbit-follow is active (CAMERA.ORBIT cam, entity, dist), forward matches internal orbit yaw
// (same as CAMERA.YAW / SetRot(0, cam.Yaw(), 0)) so WASD does not drift from player facing at
// non-zero pitch. Otherwise uses CameraXZStrafeBasis.
func CameraXZWalkBasis(store *heap.Store, ch heap.Handle) (fwd rl.Vector3, right rl.Vector3, err error) {
	if store == nil {
		return rl.Vector3{}, rl.Vector3{}, fmt.Errorf("heap is nil")
	}
	o, err := heap.Cast[*camObj](store, ch)
	if err != nil {
		return CameraXZStrafeBasis(store, ch)
	}
	if o.orbitInited {
		sy := math.Sin(float64(o.orbitYaw))
		cy := math.Cos(float64(o.orbitYaw))
		// Same horizontal “into the scene” axis as orbit_follow (target - camera) on XZ → (-sin,-cos).
		fwd = rl.Vector3{X: float32(-sy), Y: 0, Z: float32(-cy)}
		h := math.Hypot(float64(fwd.X), float64(fwd.Z))
		if h < 1e-5 {
			return CameraXZStrafeBasis(store, ch)
		}
		fwd.X = float32(float64(fwd.X) / h)
		fwd.Z = float32(float64(fwd.Z) / h)
		right = rl.Vector3{X: -fwd.Z, Y: 0, Z: fwd.X}
		return fwd, right, nil
	}
	return CameraXZStrafeBasis(store, ch)
}
