//go:build !cgo

package mbmodel3d

// SetSceneAmbient is a no-op without CGO (no PBR draw path).
func SetSceneAmbient(r, g, b float32) {}

// SetSceneAmbientScaled is a no-op without CGO.
func SetSceneAmbientScaled(r, g, b, scale float32) {}
