//go:build !cgo && !windows

package mbdebug

// DrawFrameOverlay is a no-op without Raylib.
func (m *Module) DrawFrameOverlay() {}
