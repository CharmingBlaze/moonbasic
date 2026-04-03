// Package mbaudio registers AUDIO.* (Raylib when CGO; hooks WINDOW open/close).
package mbaudio

import "moonbasic/vm/heap"

// Module manages audio device lifecycle and heap-backed streams/waves/sounds.
type Module struct {
	h *heap.Store
}

// NewModule creates an audio module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// OnWindowOpen is registered on the window module (after InitWindow).
func (m *Module) OnWindowOpen() { raylibAudioOpen() }

// OnWindowClose is registered on the window module (before CloseWindow).
func (m *Module) OnWindowClose() { raylibAudioClose() }
