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

	"moonbasic/vm/heap"
)

// Module holds window/render state for one Registry. Use NewModule and RegisterModule.
type Module struct {
	mu      sync.Mutex
	inFrame bool
	opened  bool

	diagOut io.Writer
	logFPS  bool
	fpsTick int // incremented each RENDER.FRAME when logFPS

	onAudioOpen  func()
	onAudioClose func()

	// frameEndHook runs just before EndDrawing each RENDER.FRAME (e.g. DEBUG watch overlay).
	frameEndHook func()

	h *heap.Store

	autoMu           sync.Mutex
	automationRec    bool
	activeAutoHandle heap.Handle // heap handle of list passed to EVENT.SETACTIVELIST; 0 = none
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

// BindHeap implements runtime.HeapAware (automation event lists use handles).
func (m *Module) BindHeap(h *heap.Store) {
	m.h = h
}
