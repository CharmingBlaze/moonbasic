//go:build !linux || !cgo

package mbphysics3d

// SetPickLayerLookup is a no-op without Jolt.
func SetPickLayerLookup(fn func(int64) (uint8, bool)) { _ = fn }

func resetPickState() {}
