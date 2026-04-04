// Package mbtilemap implements TILEMAP.* for Tiled TMX (orthogonal, CSV) when CGO is enabled.
package mbtilemap

import "moonbasic/vm/heap"

// Module registers tilemap builtins.
type Module struct {
	h *heap.Store
}

// NewModule creates the module.
func NewModule() *Module { return &Module{} }

// BindHeap binds the VM heap before Register.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }
