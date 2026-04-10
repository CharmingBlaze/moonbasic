//go:build !linux || !cgo

package mbphysics3d

// SetPickLayerLookup is a no-op without Jolt.
func SetPickLayerLookup(fn func(int64) (uint8, bool)) { _ = fn }

func resetPickState() {}

// PickCastEntityID is unavailable without Jolt.
func PickCastEntityID(ox, oy, oz, dx, dy, dz, maxDist float64) int64 {
	_, _, _, _, _, _, _ = ox, oy, oz, dx, dy, dz, maxDist
	return 0
}

// RaycastDownNormal is unavailable without Jolt.
func RaycastDownNormal(ox, oy, oz, maxDown float64) (nx, ny, nz float64, ok bool) {
	_, _, _, _ = ox, oy, oz, maxDown
	return 0, 1, 0, false
}
