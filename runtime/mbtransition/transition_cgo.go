//go:build cgo || (windows && !cgo)

package mbtransition

import (
	"fmt"
	"image/color"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/window"
	"moonbasic/vm/value"
)

// RegisterFrameHook registers the transition overlay (CGO only).
func RegisterFrameHook(w *window.Module) {
	if w == nil {
		return
	}
	w.AppendFrameDrawHook(transitionDraw)
}

func (m *Module) Register(reg runtime.Registrar) {
	reg.Register("TRANSITION.FADEOUT", "transition", runtime.AdaptLegacy(transFadeOut))
	reg.Register("TRANSITION.FADEIN", "transition", runtime.AdaptLegacy(transFadeIn))
	reg.Register("TRANSITION.ISDONE", "transition", runtime.AdaptLegacy(transIsDone))
	reg.Register("TRANSITION.WIPE", "transition", m.transWipe)
	reg.Register("TRANSITION.SETCOLOR", "transition", runtime.AdaptLegacy(transSetColor))
}

func (m *Module) Shutdown() {}

func argF(v value.Value) (float32, bool) {
	if f, ok := v.ToFloat(); ok {
		return float32(f), true
	}
	if i, ok := v.ToInt(); ok {
		return float32(i), true
	}
	return 0, false
}

func argInt32(v value.Value) (int32, bool) {
	if i, ok := v.ToInt(); ok {
		return int32(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return int32(f), true
	}
	return 0, false
}

func transFadeOut(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("TRANSITION.FADEOUT expects (seconds)")
	}
	d, ok := argF(args[0])
	if !ok || d <= 0 {
		return value.Nil, fmt.Errorf("TRANSITION.FADEOUT: duration must be positive")
	}
	trMu.Lock()
	defer trMu.Unlock()
	trMode = trFadeOut
	trDuration = d
	trElapsed = 0
	trDone = false
	return value.Nil, nil
}

func transFadeIn(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("TRANSITION.FADEIN expects (seconds)")
	}
	d, ok := argF(args[0])
	if !ok || d <= 0 {
		return value.Nil, fmt.Errorf("TRANSITION.FADEIN: duration must be positive")
	}
	trMu.Lock()
	defer trMu.Unlock()
	trMode = trFadeIn
	trDuration = d
	trElapsed = 0
	trDone = false
	return value.Nil, nil
}

func transIsDone(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("TRANSITION.ISDONE expects 0 arguments")
	}
	trMu.Lock()
	d := trDone
	trMu.Unlock()
	return value.FromBool(d), nil
}

func (m *Module) transWipe(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = m
	if len(args) != 2 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("TRANSITION.WIPE expects (direction$, seconds)")
	}
	dir, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	d, ok := argF(args[1])
	if !ok || d <= 0 {
		return value.Nil, fmt.Errorf("TRANSITION.WIPE: duration must be positive")
	}
	trMu.Lock()
	defer trMu.Unlock()
	trMode = trWipe
	trWipeDir = strings.TrimSpace(strings.ToLower(dir))
	trDuration = d
	trElapsed = 0
	trDone = false
	return value.Nil, nil
}

func transSetColor(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("TRANSITION.SETCOLOR expects (r, g, b, a)")
	}
	r0, _ := argInt32(args[0])
	g0, _ := argInt32(args[1])
	b0, _ := argInt32(args[2])
	a0, _ := argInt32(args[3])
	trMu.Lock()
	trColor = color.RGBA{R: clampU8(r0), G: clampU8(g0), B: clampU8(b0), A: clampU8(a0)}
	trMu.Unlock()
	return value.Nil, nil
}

func clampU8(n int32) uint8 {
	if n < 0 {
		return 0
	}
	if n > 255 {
		return 255
	}
	return uint8(n)
}

func transitionDraw() {
	trMu.Lock()
	mode := trMode
	if mode == trIdle {
		trMu.Unlock()
		return
	}
	dt := rl.GetFrameTime()
	trElapsed += dt
	t := float32(1)
	if trDuration > 0 {
		t = trElapsed / trDuration
	}
	fin := t >= 1
	if fin {
		trDone = true
		t = 1
		if mode == trFadeIn {
			trMode = trIdle
		}
	}
	col := trColor
	dir := trWipeDir
	trMu.Unlock()

	w := float32(rl.GetRenderWidth())
	h := float32(rl.GetRenderHeight())
	if w < 1 || h < 1 {
		return
	}

	switch mode {
	case trFadeOut:
		a := uint8(float32(col.A) * clamp01(t))
		if fin {
			a = col.A
		}
		rl.DrawRectangle(0, 0, int32(w), int32(h), color.RGBA{R: col.R, G: col.G, B: col.B, A: a})
	case trFadeIn:
		a := uint8(float32(col.A) * (1 - clamp01(t)))
		rl.DrawRectangle(0, 0, int32(w), int32(h), color.RGBA{R: col.R, G: col.G, B: col.B, A: a})
	case trWipe:
		p := clamp01(t)
		switch dir {
		case "left":
			rw := w * p
			rl.DrawRectangle(0, 0, int32(rw), int32(h), col)
		case "right":
			x0 := w * (1 - p)
			rl.DrawRectangle(int32(x0), 0, int32(w-x0), int32(h), col)
		case "up", "top":
			rh := h * p
			rl.DrawRectangle(0, 0, int32(w), int32(rh), col)
		case "down", "bottom":
			y0 := h * (1 - p)
			rl.DrawRectangle(0, int32(y0), int32(w), int32(h-y0), col)
		default:
			rl.DrawRectangle(0, 0, int32(w), int32(h), col)
		}
	}
}
