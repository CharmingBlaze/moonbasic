//go:build cgo || (windows && !cgo)

package window

import (
	"fmt"
	"os"
	goruntime "runtime"

	"moonbasic/internal/driver"
	"moonbasic/runtime"
	"moonbasic/runtime/mbmatrix"
	mbphysics3d "moonbasic/runtime/physics3d"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Register wires Raylib-backed WINDOW.* and minimal RENDER.* handlers.

func (m *Module) Register(reg runtime.Registrar) {
	reg.Register("WINDOW.OPEN", "window", m.wOpen)
	reg.Register("WINDOW.CANOPEN", "window", m.wCanOpen)
	reg.Register("WINDOW.SETFPS", "window", m.wSetFPS)
	reg.Register("WINDOW.SETMSAA", "window", m.wSetMSAA)
	reg.Register("WINDOW.CLOSE", "window", m.wClose)
	reg.Register("WINDOW.SHOULDCLOSE", "window", m.wShouldClose)
	reg.Register("RENDER.CLEAR", "render", m.rClear)
	reg.Register("RENDER.FRAME", "render", m.rFrame)
	reg.Register("RENDER.BEGIN3D", "render", m.rBegin3D)
	reg.Register("RENDER.END3D", "render", m.rEnd3D)
	m.registerRenderAdvanced(reg)
	m.registerPostCommands(reg)
	m.registerEffectCommands(reg)
	m.registerComputeShaderCommands(reg)
	m.registerDecalCommands(reg)
	m.registerWindowStateCommands(reg)
	m.registerWindowMetricsCommands(reg)
	m.registerWindowPlacementCommands(reg)
	m.registerLoadingModeCommands(reg)
	m.registerAutomationCommands(reg)
	m.registerBlitzSysCommands(reg)
	m.registerBlitzDisplayQueries(reg)

	// Global shorthands (Easy Mode)
	reg.Register("SKYCOLOR", "draw", m.rClear)
	reg.Register("FPS", "window", runtime.AdaptLegacy(m.wGetFPS))
	reg.Register("AMBIENTLIGHT", "draw", runtime.AdaptLegacy(m.rSetAmbient))
}

func (m *Module) wOpen(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("WINDOW.OPEN expects 3 arguments (width, height, title)")
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

	if m.usePuregoDLL() {
		return m.puregoWOpen(rt, args...)
	}
	if err := driver.CheckWindow(m.driverSel); err != nil {
		return value.Nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.opened {
		m.closeWindowLocked()
	}

	var flags uint32
	// High-DPI is opt-in (MOONBASIC_ENABLE_HIGHDPI=1). Always forcing FLAG_WINDOW_HIGHDPI breaks many
	// Windows + Intel setups (white window, wrong render buffer vs client size).
	if windowOpenWantHighDPI() {
		flags |= uint32(rl.FlagWindowHighdpi)
	}
	// Raylib exposes a single MSAA hint bit; treat 2/4/8 sample requests as enabling it.
	if m.msaaSamples >= 2 {
		flags |= uint32(rl.FlagMsaa4xHint)
	}
	rl.SetConfigFlags(flags)
	rl.InitWindow(int32(w), int32(h), title)
	rl.SetTargetFPS(60)
	if !rl.IsWindowReady() {
		rl.CloseWindow()
		m.opened = false
		m.inFrame = false
		fmt.Fprintf(os.Stderr, "moonBASIC: could not open window %dx%d %q\n", w, h, title)
		os.Exit(1)
	}
	// Nudge framebuffer to match requested client pixels (fixes 1280x720 window vs scaled render size on some drivers).
	rl.SetWindowSize(int(w), int(h))

	m.opened = true
	m.inFrame = false
	m.fpsTick = 0
	mbphysics3d.LogJoltPhysicsBackendHint()
	if minimalOpenHandshake() {
		rl.PollInputEvents()
	} else {
		presentWindowOpenIrisGuard()
		portabilityDrainAfterOpen()
		openWarmupBlankFrames(openWarmupBlankFrameCount())
	}
	if m.onAudioOpen != nil {
		m.onAudioOpen()
	}
	return value.Nil, nil
}

// presentWindowOpenIrisGuard primes presentation on strict GL drivers (e.g. Intel Iris on Windows).
// Does not set m.inFrame. MOONBASIC_SKIP_OPEN_PRESENT_KICK=1 skips (poll only).
// MOONBASIC_SKIP_WINDOW_FOCUS=1 skips SetWindowFocused if focus causes issues.
func presentWindowOpenIrisGuard() {
	if envTruthy("MOONBASIC_SKIP_OPEN_PRESENT_KICK") {
		rl.PollInputEvents()
		return
	}
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.EndDrawing()
	rl.SetWindowState(rl.FlagWindowResizable)
	rl.ClearWindowState(rl.FlagWindowResizable)
	if !envTruthy("MOONBASIC_SKIP_WINDOW_FOCUS") {
		rl.SetWindowFocused()
	}
	rl.PollInputEvents()
}

func portabilityDrainAfterOpen() {
	n := 32
	if envTruthy("MOONBASIC_SAFE_WINDOW") {
		n = 128
	}
	for i := 0; i < n; i++ {
		rl.PollInputEvents()
	}
}

func openWarmupBlankFrames(n int) {
	for i := 0; i < n; i++ {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.EndDrawing()
		rl.PollInputEvents()
	}
}

func (m *Module) wCanOpen(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("WINDOW.CANOPEN expects 3 arguments (width, height, title)")
	}
	w, okw := argInt(args[0])
	h, okh := argInt(args[1])
	if !okw || !okh {
		return value.Nil, fmt.Errorf("WINDOW.CANOPEN: width and height must be numeric")
	}
	if args[2].Kind != value.KindString {
		return value.Nil, fmt.Errorf("WINDOW.CANOPEN: title must be a string")
	}
	title, err := rt.ArgString(args, 2)
	if err != nil {
		return value.Nil, err
	}
	ok := w > 0 && h > 0 && title != ""
	return value.FromBool(ok), nil
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
	if m.usePuregoDLL() {
		return m.puregoWSetFPS(rt, args...)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.opened {
		return value.Nil, fmt.Errorf("WINDOW.SETFPS: window is not open (call WINDOW.OPEN first)")
	}
	rl.SetTargetFPS(int32(fps))
	return value.Nil, nil
}

func (m *Module) wSetMSAA(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WINDOW.SETMSAA expects 1 argument (samples)")
	}
	s, ok := argInt(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("WINDOW.SETMSAA: samples must be numeric")
	}
	m.mu.Lock()
	m.msaaSamples = int32(s)
	m.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) wClose(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.CLOSE expects 0 arguments")
	}
	if m.usePuregoDLL() {
		return m.puregoWClose(rt, args...)
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
	if m.usePuregoDLL() {
		g, err := m.ensureSidecar()
		if err != nil {
			m.opened = false
			m.inFrame = false
			return
		}
		m.closeWindowLockedPurego(g)
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
	if m.usePuregoDLL() {
		return m.puregoWShouldClose(rt, args...)
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
	if m.usePuregoDLL() {
		return m.puregoRClear(rt, args...)
	}
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
	if m.usePuregoDLL() {
		return m.puregoRFrame(rt, args...)
	}
	m.drainCleanupQueue()
	if rt != nil {
		// Advance time lerps (SlowMotion, etc.)
		rt.Call("TIME.UPDATE", []value.Value{value.FromFloat(float64(rl.GetFrameTime()))})
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
	if rt != nil {
		if msg := rt.LastScriptErrorMessage(); msg != "" {
			const maxDraw = 220
			s := "Script error: " + msg
			if len(s) > maxDraw {
				s = s[:maxDraw] + "…"
			}
			rl.DrawText(s, 8, 8, 18, rl.Red)
			if ln := rt.LastScriptErrorLine(); ln > 0 {
				rl.DrawText(fmt.Sprintf("line %d", ln), 8, 30, 16, rl.Maroon)
			}
		}
	}

	// Visual Polish: World Flash
	if m.flashDur > 0 && m.flashElapsed < m.flashDur {
		dt := rl.GetFrameTime()
		m.flashElapsed += dt
		alpha := 1.0 - (m.flashElapsed / m.flashDur)
		if alpha > 0 {
			fCol := m.flashColor
			fCol.A = uint8(float32(fCol.A) * alpha)
			sw := int32(rl.GetScreenWidth())
			sh := int32(rl.GetScreenHeight())
			rl.DrawRectangle(0, 0, sw, sh, fCol)
		}
	}

	rl.EndDrawing()
	m.inFrame = false
	goruntime.Gosched()

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

// Reset clears per-frame state.
func (m *Module) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.inFrame = false
	m.flashDur = 0
}
