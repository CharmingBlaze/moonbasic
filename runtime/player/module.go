// Package player registers PLAYER.* high-level kinematic / interaction helpers (Linux+Jolt KCC).
package player

import (
	mbcharcontroller "moonbasic/runtime/charcontroller"
	mbentity "moonbasic/runtime/mbentity"
	"moonbasic/runtime"
	mwater "moonbasic/runtime/water"
	"moonbasic/vm/heap"
)

// Module implements PLAYER.* builtins.
type Module struct {
	h    *heap.Store
	char *mbcharcontroller.Module
	ent  *mbentity.Module
	water *mwater.Module

	entToChar map[int64]heap.Handle
	state     map[int64]int32
	// stepHeight stores desired max stair step (reserved; Jolt wrapper does not expose runtime step height yet).
	stepHeight map[int64]float64
	grab       map[int64]int64   // player entity# -> grabbed entity# (0 = none)
	fovKick    map[int64]float64 // degrees added to camera FOV (read via GETFOVKICK; apply in script or future hook)
}

// NewModule constructs the player module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) {
	m.h = h
	if m.entToChar == nil {
		m.entToChar = make(map[int64]heap.Handle)
	}
	if m.state == nil {
		m.state = make(map[int64]int32)
	}
	if m.stepHeight == nil {
		m.stepHeight = make(map[int64]float64)
	}
	if m.grab == nil {
		m.grab = make(map[int64]int64)
	}
	if m.fovKick == nil {
		m.fovKick = make(map[int64]float64)
	}
}

// BindWater wires the water module for PLAYER.ISSWIMMING (optional).
func (m *Module) BindWater(w *mwater.Module) { m.water = w }

// Bind wires character controller + entity modules (see compiler pipeline wirePlayerModules).
func (m *Module) Bind(char *mbcharcontroller.Module, ent *mbentity.Module) {
	m.char = char
	m.ent = ent
	mbentity.SetCharacterGroundNormalResolver(func(id int64) (float64, float64, float64, bool) {
		if m.char == nil {
			return 0, 0, 0, false
		}
		h, ok := m.entToChar[id]
		if !ok {
			return 0, 0, 0, false
		}
		return m.char.CharacterGroundNormal(h)
	})
}

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	registerPlayerCommands(m, reg)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {
	mbentity.SetCharacterGroundNormalResolver(nil)
}

func (m *Module) Reset() {}

