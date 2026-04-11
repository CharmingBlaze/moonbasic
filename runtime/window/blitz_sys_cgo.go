//go:build cgo || (windows && !cgo)

package window

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (m *Module) registerBlitzSysCommands(reg runtime.Registrar) {
	reg.Register("AppTitle", "window", m.appTitleBlitz)
	reg.Register("Graphics3D", "window", m.graphics3DBlitz)
	reg.Register("ActiveShader", "window", m.activeShaderBlitz)
	reg.Register("SetMSAA", "window", m.wSetMSAA)
	reg.Register("SetSSAO", "window", m.setSSAOBlitz)
	reg.Register("SetPostProcess", "window", m.postAddShader)
	reg.Register("SetBloom", "window", m.setBloomBlitz)
}

func (m *Module) setBloomBlitz(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.effectBloom(rt, args...)
}

func (m *Module) setSSAOBlitz(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.effectSSAO(rt, args...)
}

func (m *Module) appTitleBlitz(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.wSetTitle(rt, args...)
}

func (m *Module) activeShaderBlitz(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.postAddShader(rt, args...)
}

// Graphics3D resizes the window; depth is reserved (z-buffer is configured at WINDOW.OPEN).
// mode: bit 0 = enable FLAG_WINDOW_HIGHDPI via SetWindowState (best-effort after open).
// Two-argument form Graphics3D(w, h) uses depth=24 (reserved) and mode=1 (HighDPI).
func (m *Module) graphics3DBlitz(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := m.requireOpen("Graphics3D"); err != nil {
		return value.Nil, err
	}
	var w, h int64
	var mode int64 = 0
	switch len(args) {
	case 2:
		var okw, okh bool
		w, okw = argInt(args[0])
		h, okh = argInt(args[1])
		if !okw || !okh || w < 1 || h < 1 {
			return value.Nil, fmt.Errorf("Graphics3D: width and height must be positive")
		}
	case 4:
		var okw, okh bool
		w, okw = argInt(args[0])
		h, okh = argInt(args[1])
		if !okw || !okh || w < 1 || h < 1 {
			return value.Nil, fmt.Errorf("Graphics3D: width and height must be positive")
		}
		_, okd := argInt(args[2])
		if !okd {
			return value.Nil, fmt.Errorf("Graphics3D: depth must be numeric (reserved)")
		}
		var okm bool
		mode, okm = args[3].ToInt()
		if !okm {
			return value.Nil, fmt.Errorf("Graphics3D: mode must be numeric")
		}
	default:
		return value.Nil, fmt.Errorf("Graphics3D expects (width, height) or (width, height, depth, mode)")
	}
	rl.SetWindowSize(int(w), int(h))
	if mode&1 != 0 {
		rl.SetWindowState(rl.FlagWindowHighdpi)
	}
	_ = rt
	return value.Nil, nil
}
