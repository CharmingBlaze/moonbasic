//go:build (!linux && !windows) || !cgo

package mbphysics3d

// MatrixBufferForEntitySync is only populated on Linux+CGO+Jolt builds.
func MatrixBufferForEntitySync() []float32 { return nil }
