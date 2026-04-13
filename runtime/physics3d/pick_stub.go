//go:build !linux || !cgo
//
// Non-Jolt stub: same API as pick_linux.go; raycasts return no hit. Keeps Windows + linux-nocgo
// fullruntime builds working. See AGENTS.md “Physics sync & Jolt”.
//
package mbphysics3d

// SetPickLayerLookup is a no-op without Jolt.
func SetPickLayerLookup(fn func(int64) (uint8, bool)) { _ = fn }

func resetPickState() {}

// PickCastEntityID is unavailable without Jolt.
func PickCastEntityID(ox, oy, oz, dx, dy, dz, maxDist float64) int64 {
	_, _, _, _, _, _, _ = ox, oy, oz, dx, dy, dz, maxDist
	return 0
}

var groundRaycastHook func(ox, oy, oz, maxDown float64) (nx, ny, nz, hitY float64, ok bool)

func SetRaycastHook(fn func(ox, oy, oz, maxDown float64) (nx, ny, nz, hitY float64, ok bool)) {
	groundRaycastHook = fn
}

// RaycastDownGroundProbe is unavailable without Jolt (unless hook is wired).
func RaycastDownGroundProbe(ox, oy, oz, maxDown float64) (nx, ny, nz, hitY float64, ok bool) {
	if groundRaycastHook != nil {
		return groundRaycastHook(ox, oy, oz, maxDown)
	}
	return 0, 1, 0, 0, false
}

// RaycastDownNormal is unavailable without Jolt.
func RaycastDownNormal(ox, oy, oz, maxDown float64) (nx, ny, nz float64, ok bool) {
	nx, ny, nz, _, ok = RaycastDownGroundProbe(ox, oy, oz, maxDown)
	return nx, ny, nz, ok
}
