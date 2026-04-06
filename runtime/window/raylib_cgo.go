//go:build cgo || (windows && !cgo)

package window

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func argInt(v value.Value) (int64, bool) {
	if i, ok := v.ToInt(); ok {
		return i, true
	}
	if f, ok := v.ToFloat(); ok {
		return int64(f), true
	}
	return 0, false
}

func clampU8(v int64) uint8 {
	switch {
	case v < 0:
		return 0
	case v > 255:
		return 255
	default:
		return uint8(v)
	}
}

// Register wires Raylib-backed WINDOW.* and minimal RENDER.* handlers.

func (m *Module) Register(reg runtime.Registrar) {
	reg.Register("WINDOW.OPEN", "window", m.wOpen)
	reg.Register("WINDOW.SETFPS", "window", m.wSetFPS)
	reg.Register("WINDOW.CLOSE", "window", m.wClose)
	reg.Register("WINDOW.SHOULDCLOSE", "window", m.wShouldClose)
	reg.Register("RENDER.CLEAR", "render", m.rClear)
	reg.Register("RENDER.FRAME", "render", m.rFrame)
	m.registerRenderAdvanced(reg)
	m.registerPostCommands(reg)
	m.registerEffectCommands(reg)
	m.registerComputeShaderCommands(reg)
	m.registerDecalCommands(reg)
	m.registerWindowStateCommands(reg)
	m.registerWindowMetricsCommands(reg)
	m.registerWindowPlacementCommands(reg)
	m.registerAutomationCommands(reg)
}

func (m *Module) wOpen(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("WINDOW.OPEN expects 3 arguments (width, height, title$)")
	}
	w, okw := argInt(args[0])
	h, okh := argInt(args[1])
	if !okw || !okh {
		return value.Nil, fmt.Errorf("WINDOW.OPEN: width and height must be numeric")
	}
	if args[2].Kind != value.KindString {
		return value.Nil, fmt.Errorf("WINDOW.OPEN: title must be a string")
	}
	title, err := rt.ArgString(args, 2)
	if err != nil {
		return value.Nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.opened {
		m.closeWindowLocked()
	}

	rl.InitWindow(int32(w), int32(h), title)
	rl.SetTargetFPS(60)
	if !rl.IsWindowReady() {
		rl.CloseWindow()
		m.opened = false
		m.inFrame = false
		return value.FromBool(false), nil
	}
	m.opened = true
	m.inFrame = false
	m.fpsTick = 0
	if m.onAudioOpen != nil {
		m.onAudioOpen()
	}
	return value.FromBool(true), nil
}

func (m *Module) wSetFPS(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WINDOW.SETFPS expects 1 argument (fps)")
	}
	fps, ok := argInt(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("WINDOW.SETFPS: fps must be numeric")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.opened {
		return value.Nil, fmt.Errorf("WINDOW.SETFPS: window is not open (call WINDOW.OPEN first)")
	}
	rl.SetTargetFPS(int32(fps))
	return value.Nil, nil
}

func (m *Module) wClose(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.CLOSE expects 0 arguments")
	}
	m.mu.Lock()
	m.closeWindowLocked()
	m.mu.Unlock()
	if rt != nil && rt.Heap != nil {
		rt.Heap.FreeAll()
	}
	return value.Nil, nil
}

func (m *Module) closeWindowLocked() {
	if !m.opened {
		return
	}
	m.shutdownAutomation()
	if m.inFrame {
		rl.EndDrawing()
		m.inFrame = false
	}
	if m.onAudioClose != nil {
		m.onAudioClose()
	}
	rl.CloseWindow()
	m.opened = false
}

func (m *Module) wShouldClose(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.SHOULDCLOSE expects 0 arguments")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.opened {
		// No window yet — not "closing"; avoids spinning a loop forever before OPEN.
		return value.FromBool(false), nil
	}
	return value.FromBool(rl.WindowShouldClose()), nil
}

func (m *Module) rClear(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	var c rl.Color
	switch len(args) {
	case 0:
		c = rl.Color{R: 0, G: 0, B: 0, A: 255}
	case 1:
		if rt == nil || rt.Heap == nil {
			return value.Nil, fmt.Errorf("RENDER.CLEAR: runtime heap not available for color handle")
		}
		if args[0].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("RENDER.CLEAR: single argument must be a color handle")
		}
		rgba, err := mbmatrix.HeapColorRGBA(rt.Heap, heap.Handle(args[0].IVal))
		if err != nil {
			return value.Nil, fmt.Errorf("RENDER.CLEAR: %w", err)
		}
		c = rl.Color{R: rgba.R, G: rgba.G, B: rgba.B, A: rgba.A}
	case 3:
		rn, ok1 := argInt(args[0])
		gn, ok2 := argInt(args[1])
		bn, ok3 := argInt(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("RENDER.CLEAR: r, g, b must be numeric")
		}
		c = rl.Color{R: clampU8(rn), G: clampU8(gn), B: clampU8(bn), A: 255}
	case 4:
		rn, ok1 := argInt(args[0])
		gn, ok2 := argInt(args[1])
		bn, ok3 := argInt(args[2])
		an, ok4 := argInt(args[3])
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return value.Nil, fmt.Errorf("RENDER.CLEAR: r, g, b, a must be numeric")
		}
		c = rl.Color{R: clampU8(rn), G: clampU8(gn), B: clampU8(bn), A: clampU8(an)}
	default:
		return value.Nil, fmt.Errorf("RENDER.CLEAR: expected 0, 1 (color handle), 3 (rgb), or 4 (rgba) arguments, got %d", len(args))
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.opened {
		return value.Nil, fmt.Errorf("RENDER.CLEAR: window is not open (call WINDOW.OPEN first)")
	}
	if !m.inFrame {
		rl.BeginDrawing()
		m.inFrame = true
	}
	if !postRenderTargetBegin(c) {
		rl.ClearBackground(c)
	}
	return value.Nil, nil
}

func (m *Module) rFrame(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("RENDER.FRAME expects 0 arguments")
	}
	m.mu.Lock()
	if !m.opened {
		m.mu.Unlock()
		return value.Nil, fmt.Errorf("RENDER.FRAME: window is not open")
	}
	if !m.inFrame {
		m.mu.Unlock()
		return value.Nil, fmt.Errorf("RENDER.FRAME: no active frame (call RENDER.CLEAR first)")
	}
	hook := m.frameEndHook
	m.mu.Unlock()

	postRenderTargetPresent()

	m.mu.Lock()
	n := len(m.frameDrawHooks)
	if cap(m.frameHookScratch) < n {
		m.frameHookScratch = make([]func(), n)
	} else {
		m.frameHookScratch = m.frameHookScratch[:n]
	}
	copy(m.frameHookScratch, m.frameDrawHooks)
	layers := m.frameHookScratch
	m.mu.Unlock()
	for _, fn := range layers {
		if fn != nil {
			fn()
		}
	}

	if hook != nil {
		hook()
	}

	if rt != nil {
		rt.FrameCount++
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.opened || !m.inFrame {
		return value.Nil, nil
	}
	rl.EndDrawing()
	m.inFrame = false

	if m.logFPS && m.diagOut != nil {
		m.fpsTick++
		if m.fpsTick%60 == 0 {
			fmt.Fprintf(m.diagOut, "[moonBASIC] raylib GetFPS: %d (WINDOW.SETFPS / default cap; vsync/GPU may vary)\n", rl.GetFPS())
		}
	}
	return value.Nil, nil
}

// Shutdown closes the window if it is still open.
func (m *Module) Shutdown() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closeWindowLocked()
}
