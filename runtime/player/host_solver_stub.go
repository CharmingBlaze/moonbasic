//go:build !cgo && !windows

package player

// Process is a no-op without entity kinematic helpers (Linux, CGO disabled).
func (m *Module) Process(dt float64) {
	_ = m
	_ = dt
}

func (m *Module) processNav(dt float64) {
	_ = m
	_ = dt
}
