//go:build !cgo && !windows

package mbentity

// SetKinematicCharacterLookup is a no-op without CGO (see kinematic_lookup_cgo.go).
func SetKinematicCharacterLookup(fn func(int64) bool) { _ = fn }

// SetCharacterGroundNormalResolver is a no-op without CGO (see entity_gameplay_intel_cgo.go).
func SetCharacterGroundNormalResolver(fn func(int64) (float64, float64, float64, bool)) { _ = fn }
