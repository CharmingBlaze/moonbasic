//go:build cgo

package window

import (
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmodel3d"
	"moonbasic/vm/value"
)

func (m *Module) registerRenderAdvanced(r runtime.Registrar) {
	r.Register("RENDER.SETBLEND", "window", runtime.AdaptLegacy(m.rSetBlend))
	r.Register("RENDER.SETBLENDMODE", "window", runtime.AdaptLegacy(m.rSetBlend))
	r.Register("RENDER.SETDEPTHWRITE", "window", runtime.AdaptLegacy(m.rSetDepthWrite))
	r.Register("RENDER.SETDEPTHMASK", "window", runtime.AdaptLegacy(m.rSetDepthWrite))
	r.Register("RENDER.SETDEPTHTEST", "window", runtime.AdaptLegacy(m.rSetDepthTest))
	r.Register("RENDER.SETSCISSOR", "window", runtime.AdaptLegacy(m.rSetScissor))
	r.Register("RENDER.CLEARSCISSOR", "window", runtime.AdaptLegacy(m.rClearScissor))
	r.Register("RENDER.SETWIREFRAME", "window", runtime.AdaptLegacy(m.rSetWireframe))
	r.Register("RENDER.SCREENSHOT", "window", m.rScreenshot)
	r.Register("RENDER.SETMSAA", "window", runtime.AdaptLegacy(m.rSetMSAA))
	r.Register("RENDER.SETSHADOWMAPSIZE", "render", runtime.AdaptLegacy(m.rSetShadowMapSize))
	r.Register("RENDER.SETMODE", "render", m.rSetRenderMode)
}

func argTruth(v value.Value) bool {
	if v.Kind == value.KindBool {
		return v.IVal != 0
	}
	if i, ok := v.ToInt(); ok {
		return i != 0
	}
	if f, ok := v.ToFloat(); ok {
		return f != 0
	}
	return false
}

func (m *Module) rSetBlend(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("RENDER.SETBLEND expects 1 argument (mode int, use BLEND_*)")
	}
	mode, ok := args[0].ToInt()
	if !ok {
		if f, okf := args[0].ToFloat(); okf {
			mode = int64(f)
		} else {
			return value.Nil, fmt.Errorf("RENDER.SETBLEND: mode must be numeric")
		}
	}
	rl.SetBlendMode(rl.BlendMode(mode))
	return value.Nil, nil
}

func (m *Module) rSetDepthWrite(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("RENDER.SETDEPTHWRITE expects 1 bool argument")
	}
	if argTruth(args[0]) {
		rl.EnableDepthMask()
	} else {
		rl.DisableDepthMask()
	}
	return value.Nil, nil
}

func (m *Module) rSetDepthTest(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("RENDER.SETDEPTHTEST expects 1 bool argument")
	}
	if argTruth(args[0]) {
		rl.EnableDepthTest()
	} else {
		rl.DisableDepthTest()
	}
	return value.Nil, nil
}

func (m *Module) rSetScissor(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("RENDER.SETSCISSOR expects 4 arguments (x, y, w, h)")
	}
	var xywh [4]int32
	for i := 0; i < 4; i++ {
		v, ok := args[i].ToInt()
		if !ok {
			if f, okf := args[i].ToFloat(); okf {
				v = int64(f)
			} else {
				return value.Nil, fmt.Errorf("RENDER.SETSCISSOR: arguments must be numeric")
			}
		}
		xywh[i] = int32(v)
	}
	rl.EnableScissorTest()
	rl.Scissor(xywh[0], xywh[1], xywh[2], xywh[3])
	return value.Nil, nil
}

func (m *Module) rClearScissor(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("RENDER.CLEARSCISSOR expects 0 arguments")
	}
	rl.DisableScissorTest()
	return value.Nil, nil
}

func (m *Module) rSetWireframe(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("RENDER.SETWIREFRAME expects 1 bool argument")
	}
	if argTruth(args[0]) {
		rl.EnableWireMode()
	} else {
		rl.DisableWireMode()
	}
	return value.Nil, nil
}

func (m *Module) rScreenshot(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("RENDER.SCREENSHOT expects 1 string argument (path$)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	rl.TakeScreenshot(path)
	return value.Nil, nil
}

func (m *Module) rSetMSAA(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("RENDER.SETMSAA expects 1 bool argument")
	}
	if argTruth(args[0]) {
		rl.SetWindowState(rl.FlagMsaa4xHint)
	} else {
		rl.ClearWindowState(rl.FlagMsaa4xHint)
	}
	return value.Nil, nil
}

func (m *Module) rSetRenderMode(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("RENDER.SETMODE expects 1 string (forward, deferred)")
	}
	modeStr, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	mode := strings.ToLower(strings.TrimSpace(modeStr))
	switch mode {
	case "forward", "deferred":
		setRenderPipelineMode(mode)
		return value.Nil, nil
	default:
		return value.Nil, fmt.Errorf("RENDER.SETMODE: unknown mode %q (use forward or deferred)", mode)
	}
}

func (m *Module) rSetShadowMapSize(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("RENDER.SETSHADOWMAPSIZE expects 1 argument (size in pixels)")
	}
	var sz int64
	if i, ok := args[0].ToInt(); ok {
		sz = i
	} else if f, ok := args[0].ToFloat(); ok {
		sz = int64(f)
	} else {
		return value.Nil, fmt.Errorf("RENDER.SETSHADOWMAPSIZE: size must be numeric")
	}
	mbmodel3d.SetShadowMapResolution(int32(sz))
	return value.Nil, nil
}
