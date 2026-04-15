//go:build (!linux && !windows) || !cgo

package mbphysics3d

// MatrixBufferForEntitySync is only populated on Linux+CGO+Jolt builds.
func MatrixBufferForEntitySync() []float32 { return nil }

// MatrixBufferPrevForEntitySync is only populated on native Jolt builds.
func MatrixBufferPrevForEntitySync() []float32 { return nil }

// PhysicsMatrixInterpAlpha is 1 on stub builds.
func PhysicsMatrixInterpAlpha() float64 { return 1 }
