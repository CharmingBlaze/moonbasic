// Package window implements WINDOW.* and a minimal RENDER.* slice backed by Raylib (CGO)
// or explicit build errors when CGO is disabled.
//
// Frame contract: RENDER.CLEAR begins a frame (BeginDrawing on first CLEAR after FRAME or OPEN)
// and clears the backbuffer. RENDER.FRAME ends the frame (EndDrawing). Each loop iteration
// should use CLEAR then FRAME.
//
// Exit semantics: WINDOW.CLOSE (and Registry.Shutdown) call EndDrawing first if a frame is
// open, then CloseWindow exactly once. Re-opening after OPEN replaces the previous window
// safely. Scripts typically exit the main loop on INPUT.KEYDOWN(KEY_ESCAPE) and then call
// WINDOW.CLOSE; the process exit code is the CLI’s (0 on normal completion), not a special
// “ESC exit code.” Closing the OS window sets Raylib’s should-close flag; scripts that only
// poll KEY_ESCAPE keep running until that key is pressed unless they also use WINDOW.SHOULDCLOSE.
package window

import (
	"io"
	"sync"

	"moonbasic/internal/driver"
	"moonbasic/internal/raylibpurego"

	"moonbasic/vm/heap"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Module holds window/render state for one Registry. Use NewModule and RegisterModule.
type Module struct {
	mu      sync.Mutex
	inFrame bool
	opened  bool

	diagOut io.Writer
	logFPS  bool
	fpsTick int // incremented each RENDER.FRAME when logFPS
	msaaSamples int32

	onAudioOpen  func()
	onAudioClose func()

	// frameEndHook runs just before EndDrawing each RENDER.FRAME (e.g. DEBUG watch overlay).
	frameEndHook func()

	// frameDrawHooks run after postRenderTargetPresent, before frameEndHook (2D lights, transitions).
	frameDrawHooks []func()

	// frameHookScratch holds a copy of frameDrawHooks for RENDER.FRAME without allocating each frame.
	frameHookScratch []func()

	h *heap.Store

	autoMu           sync.Mutex
	automationRec    bool
	activeAutoHandle heap.Handle // heap handle of list passed to EVENT.SETACTIVELIST; 0 = none

	cleanupMu    sync.Mutex
	cleanupQueue []func()

	// Driver selection from [driver.GetDefaultDriver]; set via [Module.BindDriverSelection] before Register.
	driverSel driver.Selection

	sidecarOnce   sync.Once
	sidecarLib    *raylibpurego.LoadResult
	sidecarGame   *raylibpurego.Game
	sidecarLoadErr error

	// Visual Polish: World Flash
	flashColor   rl.Color
	flashDur     float32
	flashElapsed float32
}

// NewModule allocates state for the window/render builtins.
func NewModule() *Module {
	return &Module{}
}

// SetDiagnostics configures optional stderr-style diagnostics. When logFPS is true and
// diagOut is non-nil, the CGO build logs raylib GetFPS about once per second while the
// program runs (throttled inside RENDER.FRAME). Used when the host passes pipeline.Options.Debug
// (CLI --info).
func (m *Module) SetDiagnostics(diagOut io.Writer, logFPS bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.diagOut = diagOut
	m.logFPS = logFPS
	m.fpsTick = 0
}

// SetAudioHooks registers optional callbacks invoked after WINDOW.OPEN and before WINDOW.CLOSE
// (used by runtime/audio for InitAudioDevice / CloseAudioDevice).
func (m *Module) SetAudioHooks(onOpen, onClose func()) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onAudioOpen = onOpen
	m.onAudioClose = onClose
}

// SetFrameEndHook registers a callback invoked at the end of each frame, after the
// game has drawn and before the backbuffer is presented. Pass nil to disable.
func (m *Module) SetFrameEndHook(h func()) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.frameEndHook = h
}

// AppendFrameDrawHook adds a callback run each RENDER.FRAME after the main scene
// (post RT present) and before the frameEndHook. Used for 2D lighting and transitions.
func (m *Module) AppendFrameDrawHook(h func()) {
	if h == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.frameDrawHooks = append(m.frameDrawHooks, h)
}

// BindHeap implements runtime.HeapAware (automation event lists use handles).
func (m *Module) BindHeap(h *heap.Store) {
	m.h = h
}

// BindDriverSelection records [driver.GetDefaultDriver] output for WINDOW/RENDER dispatch (CGO vs sidecar purego).
func (m *Module) BindDriverSelection(sel driver.Selection) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.driverSel = sel
}
