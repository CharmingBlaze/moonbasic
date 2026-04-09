//go:build cgo || (windows && !cgo)

package window

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerWindowMetricsCommands(reg runtime.Registrar) {
	reg.Register("WINDOW.WIDTH", "window", runtime.AdaptLegacy(m.wWidth))
	reg.Register("WINDOW.HEIGHT", "window", runtime.AdaptLegacy(m.wHeight))
	reg.Register("WINDOW.DPISCALE", "window", runtime.AdaptLegacy(m.wDPIScale))
	reg.Register("WINDOW.GETFPS", "window", runtime.AdaptLegacy(m.wGetFPS))
	reg.Register("WINDOW.ISFULLSCREEN", "window", runtime.AdaptLegacy(m.wIsFullscreen))
	reg.Register("WINDOW.TOGGLEFULLSCREEN", "window", runtime.AdaptLegacy(m.wToggleFullscreen))
	reg.Register("WINDOW.ISRESIZED", "window", runtime.AdaptLegacy(m.wIsResized))
	reg.Register("WINDOW.SETTITLE", "window", m.wSetTitle)
	reg.Register("RENDER.WIDTH", "render", runtime.AdaptLegacy(m.rRenderWidth))
	reg.Register("RENDER.HEIGHT", "render", runtime.AdaptLegacy(m.rRenderHeight))

	// Global shorthands (Easy Mode)
	reg.Register("SCREENWIDTH", "window", runtime.AdaptLegacy(m.wWidth))
	reg.Register("SCREENHEIGHT", "window", runtime.AdaptLegacy(m.wHeight))
}

func (m *Module) wWidth(args []value.Value) (value.Value, error) {
	if err := m.requireOpen("WINDOW.WIDTH"); err != nil {
		return value.Nil, err
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.WIDTH expects 0 arguments")
	}
	return value.FromFloat(float64(rl.GetScreenWidth())), nil
}

func (m *Module) wDPIScale(args []value.Value) (value.Value, error) {
	if err := m.requireOpen("WINDOW.DPISCALE"); err != nil {
		return value.Nil, err
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.DPISCALE expects 0 arguments")
	}
	v := rl.GetWindowScaleDPI()
	return value.FromFloat(float64(v.X)), nil
}

func (m *Module) wHeight(args []value.Value) (value.Value, error) {
	if err := m.requireOpen("WINDOW.HEIGHT"); err != nil {
		return value.Nil, err
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.HEIGHT expects 0 arguments")
	}
	return value.FromFloat(float64(rl.GetScreenHeight())), nil
}

func (m *Module) rRenderWidth(args []value.Value) (value.Value, error) {
	if err := m.requireOpen("RENDER.WIDTH"); err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(rl.GetRenderWidth())), nil
}

func (m *Module) rRenderHeight(args []value.Value) (value.Value, error) {
	if err := m.requireOpen("RENDER.HEIGHT"); err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(rl.GetRenderHeight())), nil
}

func (m *Module) wGetFPS(args []value.Value) (value.Value, error) {
	if err := m.requireOpen("WINDOW.GETFPS"); err != nil {
		return value.Nil, err
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.GETFPS expects 0 arguments")
	}
	return value.FromInt(int64(rl.GetFPS())), nil
}

func (m *Module) wIsFullscreen(args []value.Value) (value.Value, error) {
	if err := m.requireOpen("WINDOW.ISFULLSCREEN"); err != nil {
		return value.Nil, err
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.ISFULLSCREEN expects 0 arguments")
	}
	return value.FromBool(rl.IsWindowFullscreen()), nil
}

func (m *Module) wToggleFullscreen(args []value.Value) (value.Value, error) {
	if err := m.requireOpen("WINDOW.TOGGLEFULLSCREEN"); err != nil {
		return value.Nil, err
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.TOGGLEFULLSCREEN expects 0 arguments")
	}
	rl.ToggleFullscreen()
	return value.Nil, nil
}

func (m *Module) wIsResized(args []value.Value) (value.Value, error) {
	if err := m.requireOpen("WINDOW.ISRESIZED"); err != nil {
		return value.Nil, err
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.ISRESIZED expects 0 arguments")
	}
	return value.FromBool(rl.IsWindowResized()), nil
}

func (m *Module) wSetTitle(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := m.requireOpen("WINDOW.SETTITLE"); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("WINDOW.SETTITLE expects title$")
	}
	title, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	rl.SetWindowTitle(title)
	return value.Nil, nil
}
