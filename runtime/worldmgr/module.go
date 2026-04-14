// Package worldmgr implements WORLD.* streaming helpers (uses active terrain from terrain module).
package worldmgr

import (
	"moonbasic/vm/heap"

	scat "moonbasic/runtime/scatter"
	terr "moonbasic/runtime/terrain"
)

// Module coordinates WORLD.* builtins.
type Module struct {
	h    *heap.Store
	terr *terr.Module
	scat *scat.Module

	// vegScatter is a lazily allocated SCATTER instance for WORLD.SETVEGETATION.
	vegScatter heap.Handle

	FogMode    int
	FogColor   [4]uint8
	FogDensity float32
}

// NewModule requires the terrain module instance registered in the same registry.
func NewModule(t *terr.Module) *Module {
	return &Module{terr: t}
}

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// BindScatter wires the scatter module for WORLD.SETVEGETATION (see compiler pipeline wireWorldModules).
func (m *Module) BindScatter(s *scat.Module) { m.scat = s }


// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {
	m.vegScatter = 0
}

// Reset clears per-session world state.
func (m *Module) Reset() {
	m.vegScatter = 0
	m.FogMode = 0
	m.FogColor = [4]uint8{0, 0, 0, 0}
	m.FogDensity = 0
}
