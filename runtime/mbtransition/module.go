package mbtransition

import (
	"image/color"
	"sync"

	"moonbasic/vm/heap"
)

type Module struct {
	h *heap.Store
}

func NewModule() *Module { return &Module{} }

func (m *Module) BindHeap(h *heap.Store) { m.h = h }

func (m *Module) Reset() {
	trMu.Lock()
	trMode = trIdle
	trElapsed = 0
	trDone = true
	trMu.Unlock()
}

// Global transition state (shared with CGO hooks)
var (
	trMu       sync.Mutex
	trMode     int
	trElapsed  float32
	trDuration float32
	trDone     bool = true
	trColor    color.RGBA = color.RGBA{0, 0, 0, 255}
	trWipeDir  string
)

const (
	trIdle    = 0
	trFadeOut = 1
	trFadeIn  = 2
	trWipe    = 3
)

func clamp01(v float32) float32 {
	if v < 0 { return 0 }
	if v > 1 { return 1 }
	return v
}
