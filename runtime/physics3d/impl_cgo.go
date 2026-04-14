//go:build (linux || windows) && cgo

package mbphysics3d


func shutdownPhysics3D(m *Module) {
	_, _ = m.phStop(nil)
}
