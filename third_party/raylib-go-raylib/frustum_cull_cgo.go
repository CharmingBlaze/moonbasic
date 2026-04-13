//go:build cgo

package rl

// FrustumCullDistances returns the rlgl near/far cull distances used with BeginMode3D projection
// (same source as GetCullDistanceNear/Far). Used by moonBASIC CPU frustum extraction.
func FrustumCullDistances() (near, far float32) {
	return float32(GetCullDistanceNear()), float32(GetCullDistanceFar())
}
