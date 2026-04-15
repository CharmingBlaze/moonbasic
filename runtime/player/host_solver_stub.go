//go:build (!linux && !windows) || !cgo

package player

// Process is a no-op without native Jolt + CGO (stub builds).
func (m *Module) Process(dt float64) {
	_ = m
	_ = dt
}
