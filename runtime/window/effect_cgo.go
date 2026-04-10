//go:build cgo || (windows && !cgo)

package window

import (
	"fmt"
	"strings"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (m *Module) registerEffectCommands(r runtime.Registrar) {
	r.Register("WORLD.FLASH", "post", runtime.AdaptLegacy(m.worldFlash))
	r.Register("EFFECT.SSAO", "post", m.effectSSAO)
	r.Register("EFFECT.SSR", "post", m.effectSSR)
	r.Register("EFFECT.MOTIONBLUR", "post", m.effectMotionBlur)
	r.Register("EFFECT.DEPTHOFFIELD", "post", m.effectDOF)
	r.Register("EFFECT.BLOOM", "post", m.effectBloom)
	r.Register("EFFECT.TONEMAPPING", "post", m.effectTonemapping)
	r.Register("EFFECT.SHARPEN", "post", m.effectSharpen)
	r.Register("EFFECT.GRAIN", "post", m.effectGrain)
	r.Register("EFFECT.VIGNETTE", "post", m.effectVignette)
	r.Register("EFFECT.CHROMATICABERRATION", "post", m.effectChromatic)
	r.Register("EFFECT.FXAA", "post", m.effectFXAA)
}

func effectGuard() error {
	postMu.Lock()
	custom := postCustomOn
	postMu.Unlock()
	if custom {
		return fmt.Errorf("EFFECT.* is unavailable while POST.ADDSHADER custom pass is active")
	}
	return nil
}

func effectEnableBasics() {
	postMu.Lock()
	postActive = true
	postCustomOn = false
	ensureBuiltInPostShader()
	postMu.Unlock()
}

func argBool01(v value.Value) bool {
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

func argFloat32(v value.Value) (float32, bool) {
	if f, ok := v.ToFloat(); ok {
		return float32(f), true
	}
	if i, ok := v.ToInt(); ok {
		return float32(i), true
	}
	return 0, false
}

func (m *Module) effectSSAO(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if err := effectGuard(); err != nil {
		return value.Nil, err
	}
	if len(args) < 1 {
		return value.Nil, fmt.Errorf("EFFECT.SSAO expects (enable, radius?, intensity?)")
	}
	on := argBool01(args[0])
	effectEnableBasics()
	postMu.Lock()
	postSSAO = on
	if len(args) >= 2 {
		if f, ok := argFloat32(args[1]); ok {
			postKV["ssao.radius"] = f
		}
	}
	if len(args) >= 3 {
		if f, ok := argFloat32(args[2]); ok {
			postKV["ssao.intensity"] = f
		}
	}
	postMu.Unlock()
	return value.Nil, nil
}

func (m *Module) effectSSR(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if err := effectGuard(); err != nil {
		return value.Nil, err
	}
	if len(args) < 1 {
		return value.Nil, fmt.Errorf("EFFECT.SSR expects (enable, steps?, stride?)")
	}
	on := argBool01(args[0])
	effectEnableBasics()
	postMu.Lock()
	postSSR = on
	if len(args) >= 2 {
		if f, ok := argFloat32(args[1]); ok {
			postKV["ssr.steps"] = f
		}
	}
	if len(args) >= 3 {
		if f, ok := argFloat32(args[2]); ok {
			postKV["ssr.stride"] = f
		}
	}
	postMu.Unlock()
	return value.Nil, nil
}

func (m *Module) effectMotionBlur(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if err := effectGuard(); err != nil {
		return value.Nil, err
	}
	if len(args) < 1 {
		return value.Nil, fmt.Errorf("EFFECT.MOTIONBLUR expects (enable, strength?)")
	}
	on := argBool01(args[0])
	effectEnableBasics()
	postMu.Lock()
	postMotionBlur = on
	if len(args) >= 2 {
		if f, ok := argFloat32(args[1]); ok {
			postKV["motionblur.strength"] = f
		}
	}
	postMu.Unlock()
	return value.Nil, nil
}

func (m *Module) effectDOF(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if err := effectGuard(); err != nil {
		return value.Nil, err
	}
	if len(args) < 1 {
		return value.Nil, fmt.Errorf("EFFECT.DEPTHOFFIELD expects (enable, focus?, range?)")
	}
	on := argBool01(args[0])
	effectEnableBasics()
	postMu.Lock()
	postDOF = on
	if len(args) >= 2 {
		if f, ok := argFloat32(args[1]); ok {
			postKV["dof.focus"] = f
		}
	}
	if len(args) >= 3 {
		if f, ok := argFloat32(args[2]); ok {
			postKV["dof.range"] = f
		}
	}
	postMu.Unlock()
	return value.Nil, nil
}

func (m *Module) effectBloom(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if err := effectGuard(); err != nil {
		return value.Nil, err
	}
	if len(args) < 1 {
		return value.Nil, fmt.Errorf("EFFECT.BLOOM expects (enable, threshold?, intensity?)")
	}
	on := argBool01(args[0])
	effectEnableBasics()
	postMu.Lock()
	postBloom = on
	if len(args) >= 2 {
		if f, ok := argFloat32(args[1]); ok {
			postKV["bloom.threshold"] = f
		}
	}
	if len(args) >= 3 {
		if f, ok := argFloat32(args[2]); ok {
			postKV["bloom.intensity"] = f
		}
	}
	postMu.Unlock()
	return value.Nil, nil
}

func (m *Module) effectTonemapping(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if err := effectGuard(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("EFFECT.TONEMAPPING expects 1 string (reinhard, filmic, aces, none)")
	}
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	effectEnableBasics()
	var mode int32
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "none", "off", "linear":
		mode = 0
	case "reinhard":
		mode = 1
	case "filmic":
		mode = 2
	case "aces":
		mode = 3
	default:
		return value.Nil, fmt.Errorf("EFFECT.TONEMAPPING: unknown mode %q", s)
	}
	postMu.Lock()
	postTonemapMode = mode
	postMu.Unlock()
	return value.Nil, nil
}

func (m *Module) effectSharpen(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if err := effectGuard(); err != nil {
		return value.Nil, err
	}
	if len(args) < 1 {
		return value.Nil, fmt.Errorf("EFFECT.SHARPEN expects (enable, amount?)")
	}
	on := argBool01(args[0])
	effectEnableBasics()
	postMu.Lock()
	postSharpen = on
	if len(args) >= 2 {
		if f, ok := argFloat32(args[1]); ok {
			postKV["sharpen.amount"] = f
		}
	}
	postMu.Unlock()
	return value.Nil, nil
}

func (m *Module) effectGrain(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if err := effectGuard(); err != nil {
		return value.Nil, err
	}
	if len(args) < 1 {
		return value.Nil, fmt.Errorf("EFFECT.GRAIN expects (enable, amount?)")
	}
	on := argBool01(args[0])
	effectEnableBasics()
	postMu.Lock()
	postGrain = on
	if len(args) >= 2 {
		if f, ok := argFloat32(args[1]); ok {
			postKV["grain.amount"] = f
		}
	}
	postMu.Unlock()
	return value.Nil, nil
}

func (m *Module) effectVignette(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if err := effectGuard(); err != nil {
		return value.Nil, err
	}
	if len(args) < 1 {
		return value.Nil, fmt.Errorf("EFFECT.VIGNETTE expects (enable, amount?)")
	}
	on := argBool01(args[0])
	effectEnableBasics()
	postMu.Lock()
	postVignette = on
	if len(args) >= 2 {
		if f, ok := argFloat32(args[1]); ok {
			postKV["vignette.strength"] = f
		}
	}
	postMu.Unlock()
	return value.Nil, nil
}

func (m *Module) effectChromatic(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if err := effectGuard(); err != nil {
		return value.Nil, err
	}
	if len(args) < 1 {
		return value.Nil, fmt.Errorf("EFFECT.CHROMATICABERRATION expects (enable, amount?)")
	}
	on := argBool01(args[0])
	effectEnableBasics()
	postMu.Lock()
	postChromatic = on
	if len(args) >= 2 {
		if f, ok := argFloat32(args[1]); ok {
			postKV["chromatic.offset"] = f
		}
	}
	postMu.Unlock()
	return value.Nil, nil
}
func (m *Module) effectFXAA(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if err := effectGuard(); err != nil {
		return value.Nil, err
	}
	if len(args) < 1 {
		return value.Nil, fmt.Errorf("EFFECT.FXAA expects (enable)")
	}
	on := argBool01(args[0])
	effectEnableBasics()
	postMu.Lock()
	postFXAA = on
	postMu.Unlock()
	return value.Nil, nil
}
func (m *Module) worldFlash(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("WORLD.FLASH expects (color#, duration#)")
	}
	colorH := heap.Handle(args[0].IVal)
	dur, ok := argFloat32(args[1])
	if !ok || dur < 0 {
		return value.Nil, fmt.Errorf("WORLD.FLASH: duration must be numeric and non-negative")
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Default to white if invalid color handle
	m.flashColor = rl.White
	if colorH != 0 && m.h != nil {
		rgba, err := mbmatrix.HeapColorRGBA(m.h, colorH)
		if err == nil {
			m.flashColor = rl.Color{R: rgba.R, G: rgba.G, B: rgba.B, A: rgba.A}
		}
	}
	m.flashDur = dur
	m.flashElapsed = 0
	return value.Nil, nil
}
