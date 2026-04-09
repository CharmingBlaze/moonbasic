//go:build cgo || (windows && !cgo)

package window

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerBlitzDisplayQueries(reg runtime.Registrar) {
	reg.Register("WindowWidth", "window", runtime.AdaptLegacy(m.blitzWindowWidth))
	reg.Register("WindowHeight", "window", runtime.AdaptLegacy(m.blitzWindowHeight))
	reg.Register("ScreenWidth", "window", runtime.AdaptLegacy(m.blitzScreenW))
	reg.Register("ScreenHeight", "window", runtime.AdaptLegacy(m.blitzScreenH))
	reg.Register("GraphicsWidth", "window", runtime.AdaptLegacy(m.blitzGraphicsW))
	reg.Register("GraphicsHeight", "window", runtime.AdaptLegacy(m.blitzGraphicsH))
	reg.Register("GraphicsDepth", "window", runtime.AdaptLegacy(m.blitzGraphicsDepth))
	reg.Register("AvailVidMem", "window", runtime.AdaptLegacy(m.blitzAvailVidMem))
	reg.Register("TotalVidMem", "window", runtime.AdaptLegacy(m.blitzTotalVidMem))
}

func (m *Module) blitzWindowWidth(args []value.Value) (value.Value, error) {
	return m.wWidth(args)
}

func (m *Module) blitzWindowHeight(args []value.Value) (value.Value, error) {
	return m.wHeight(args)
}

func (m *Module) blitzScreenW(args []value.Value) (value.Value, error) {
	return m.wWidth(args)
}

func (m *Module) blitzScreenH(args []value.Value) (value.Value, error) {
	return m.wHeight(args)
}

func (m *Module) blitzGraphicsW(args []value.Value) (value.Value, error) {
	return m.rRenderWidth(args)
}

func (m *Module) blitzGraphicsH(args []value.Value) (value.Value, error) {
	return m.rRenderHeight(args)
}

func (m *Module) blitzGraphicsDepth(args []value.Value) (value.Value, error) {
	if err := m.requireOpen("GraphicsDepth"); err != nil {
		return value.Nil, err
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("GraphicsDepth expects 0 arguments")
	}
	// Raylib does not expose drawable depth bits per target; common default for desktop GL.
	return value.FromInt(32), nil
}

func (m *Module) blitzAvailVidMem(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("AvailVidMem expects 0 arguments")
	}
	return value.FromInt(-1), nil
}

func (m *Module) blitzTotalVidMem(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("TotalVidMem expects 0 arguments")
	}
	return value.FromInt(-1), nil
}
