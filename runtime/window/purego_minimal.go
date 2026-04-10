package window

import (
	"fmt"
	"os"
	goruntime "runtime"

	"moonbasic/internal/driver"
	"moonbasic/internal/raylibpurego"
	"moonbasic/runtime"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// Raylib 5.x config flags (match github.com/gen2brain/raylib-go/raylib constants; avoid importing raylib here for !cgo Unix builds).
const (
	raylibFlagWindowHighdpi = 0x00002000
	raylibFlagMsaa4xHint    = 0x00000020
)

// usePuregoDLL is true when the process should use internal/raylibpurego (sidecar DLL/SO) for core WINDOW/RENDER.
func (m *Module) usePuregoDLL() bool {
	return m.driverSel.Kind == driver.KindPuregoDLL
}

func (m *Module) ensureSidecar() (*raylibpurego.Game, error) {
	m.sidecarOnce.Do(func() {
		lib, err := raylibpurego.Load("")
		if err != nil {
			m.sidecarLoadErr = err
			return
		}
		m.sidecarLib = lib
		g := new(raylibpurego.Game)
		m.sidecarGame = g
		m.sidecarLoadErr = raylibpurego.RegisterGame(lib, g)
	})
	if m.sidecarLoadErr != nil {
		return nil, m.sidecarLoadErr
	}
	return m.sidecarGame, nil
}

func (m *Module) puregoWOpen(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := driver.CheckWindow(m.driverSel); err != nil {
		return value.Nil, err
	}
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

	g, err := m.ensureSidecar()
	if err != nil {
		return value.Nil, fmt.Errorf("WINDOW.OPEN: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.opened {
		m.closeWindowLockedPurego(g)
	}

	var flags uint32
	if windowOpenWantHighDPI() {
		flags |= uint32(raylibFlagWindowHighdpi)
	}
	if m.msaaSamples >= 2 {
		flags |= uint32(raylibFlagMsaa4xHint)
	}
	if g.SetConfigFlags != nil {
		g.SetConfigFlags(flags)
	}
	g.InitWindow(int32(w), int32(h), title)
	if g.SetTargetFPS != nil {
		g.SetTargetFPS(60)
	}
	if g.IsWindowReady != nil && !g.IsWindowReady() {
		g.CloseWindow()
		m.opened = false
		m.inFrame = false
		fmt.Fprintf(os.Stderr, "moonBASIC: could not open window %dx%d %q\n", w, h, title)
		os.Exit(1)
	}
	m.opened = true
	m.inFrame = false
	nWarmup := 2
	if s := os.Getenv("MOONBASIC_OPEN_WARMUP_FRAMES"); s != "" {
		fmt.Sscanf(s, "%d", &nWarmup)
	}
	for i := 0; i < nWarmup; i++ {
		g.BeginDrawing()
		g.ClearBackground(raylibpurego.ColorPtr(raylibpurego.Color{R: 0, G: 0, B: 0, A: 255}))
		g.EndDrawing()
	}

	return value.Nil, nil
}

func (m *Module) closeWindowLockedPurego(g *raylibpurego.Game) {
	if !m.opened {
		return
	}
	m.shutdownAutomation()
	if m.inFrame {
		g.EndDrawing()
		m.inFrame = false
	}
	if m.onAudioClose != nil {
		m.onAudioClose()
	}
	g.CloseWindow()
	m.opened = false
}

func (m *Module) puregoWClose(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.CLOSE expects 0 arguments")
	}
	g, err := m.ensureSidecar()
	if err != nil {
		return value.Nil, err
	}
	m.mu.Lock()
	m.closeWindowLockedPurego(g)
	m.mu.Unlock()
	if rt != nil && rt.Heap != nil {
		rt.Heap.FreeAll()
	}
	return value.Nil, nil
}

func (m *Module) puregoWShouldClose(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.SHOULDCLOSE expects 0 arguments")
	}
	g, err := m.ensureSidecar()
	if err != nil {
		return value.Nil, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.opened {
		return value.FromBool(false), nil
	}
	return value.FromBool(g.WindowShouldClose()), nil
}

func (m *Module) puregoWSetFPS(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WINDOW.SETFPS expects 1 argument (fps)")
	}
	fps, ok := argInt(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("WINDOW.SETFPS: fps must be numeric")
	}
	g, err := m.ensureSidecar()
	if err != nil {
		return value.Nil, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.opened {
		return value.Nil, fmt.Errorf("WINDOW.SETFPS: window is not open (call WINDOW.OPEN first)")
	}
	if g.SetTargetFPS != nil {
		g.SetTargetFPS(int32(fps))
	}
	return value.Nil, nil
}

func (m *Module) puregoRClear(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	g, err := m.ensureSidecar()
	if err != nil {
		return value.Nil, err
	}

	var c raylibpurego.Color
	switch len(args) {
	case 0:
		c = raylibpurego.Color{R: 0, G: 0, B: 0, A: 255}
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
		c = raylibpurego.Color{R: rgba.R, G: rgba.G, B: rgba.B, A: rgba.A}
	case 3:
		rn, ok1 := argInt(args[0])
		gn, ok2 := argInt(args[1])
		bn, ok3 := argInt(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("RENDER.CLEAR: r, g, b must be numeric")
		}
		c = raylibpurego.Color{R: clampU8(rn), G: clampU8(gn), B: clampU8(bn), A: 255}
	case 4:
		rn, ok1 := argInt(args[0])
		gn, ok2 := argInt(args[1])
		bn, ok3 := argInt(args[2])
		an, ok4 := argInt(args[3])
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return value.Nil, fmt.Errorf("RENDER.CLEAR: r, g, b, a must be numeric")
		}
		c = raylibpurego.Color{R: clampU8(rn), G: clampU8(gn), B: clampU8(bn), A: clampU8(an)}
	default:
		return value.Nil, fmt.Errorf("RENDER.CLEAR: expected 0, 1 (color handle), 3 (rgb), or 4 (rgba) arguments, got %d", len(args))
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.opened {
		return value.Nil, fmt.Errorf("RENDER.CLEAR: window is not open (call WINDOW.OPEN first)")
	}
	if !m.inFrame {
		g.BeginDrawing()
		m.inFrame = true
	}
	g.ClearBackground(raylibpurego.ColorPtr(c))
	return value.Nil, nil
}

func (m *Module) puregoRFrame(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("RENDER.FRAME expects 0 arguments")
	}
	g, err := m.ensureSidecar()
	if err != nil {
		return value.Nil, err
	}

	m.drainCleanupQueue()
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
	if rt != nil && g.DrawText != nil {
		if msg := rt.LastScriptErrorMessage(); msg != "" {
			const maxDraw = 220
			s := "Script error: " + msg
			if len(s) > maxDraw {
				s = s[:maxDraw] + "…"
			}
			red := raylibpurego.Color{R: 230, G: 41, B: 55, A: 255}
			g.DrawText(s, 8, 8, 18, raylibpurego.ColorPtr(red))
			if ln := rt.LastScriptErrorLine(); ln > 0 {
				maroon := raylibpurego.Color{R: 190, G: 33, B: 55, A: 255}
				g.DrawText(fmt.Sprintf("line %d", ln), 8, 30, 16, raylibpurego.ColorPtr(maroon))
			}
		}
	}
	g.EndDrawing()
	m.inFrame = false
	goruntime.Gosched()

	if m.logFPS && m.diagOut != nil && g.GetFPS != nil {
		m.fpsTick++
		if m.fpsTick%60 == 0 {
			fmt.Fprintf(m.diagOut, "[moonBASIC] raylib GetFPS: %d (WINDOW.SETFPS / default cap; vsync/GPU may vary)\n", g.GetFPS())
		}
	}
	return value.Nil, nil
}

func (m *Module) puregoShutdown() {
	g, err := m.ensureSidecar()
	if err != nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closeWindowLockedPurego(g)
}
