// Package mbtransition implements TRANSITION.* screen fades/wipes (CGO), drawn after the main scene.
package mbtransition

import (
	"image/color"
	"sync"
)

// Module registers TRANSITION.* builtins.
type Module struct{}

// NewModule creates the module.
func NewModule() *Module { return &Module{} }

// BindHeap is a no-op (transitions are global state).
func (m *Module) BindHeap(interface{}) {}

const (
	trIdle = iota
	trFadeOut
	trFadeIn
	trWipe
)

var (
	trMu     sync.Mutex
	trMode   int
	trElapsed, trDuration float32
	trColor  = color.RGBA{0, 0, 0, 255}
	trWipeDir string
	trDone   = true
)

func transitionResetIdle() {
	trMode = trIdle
	trElapsed = 0
	trDone = true
}

func clamp01(x float32) float32 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}
