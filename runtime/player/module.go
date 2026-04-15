// Package player registers PLAYER.* high-level kinematic / interaction helpers (Linux+Jolt KCC).
package player

import (
	mbcharcontroller "moonbasic/runtime/charcontroller"
	mbentity "moonbasic/runtime/mbentity"
	mbphysics3d "moonbasic/runtime/physics3d"
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
	// stepHeight mirrors CharacterVirtual ExtendedUpdate WalkStairsStepUp (Y) for PLAYER.SETSTEPHEIGHT.
	stepHeight map[int64]float64
	grab       map[int64]int64   // player entity# -> grabbed entity# (0 = none)
	fovKick    map[int64]float64 // degrees added to camera FOV (read via GETFOVKICK; apply in script or future hook)

	lastHero int64

	// kccNav: click-to-move / NAV targets for PLAYER.CREATE entities (see CHAR.NAVTO / CHAR.NAVUPDATE).
	kccNav map[int64]*kccNavState
	// kccLastGroundedAt: physics sim time when CharacterVirtual last reported grounded — for coyote grace in ISGROUNDED.
	kccLastGroundedAt map[int64]float64
	// swimManual: entity had PLAYER.SWIM enabled manually; skip ambient WATER volume swim override while true.
	swimManual map[int64]bool
}

// kccNavMode: NAV.GOTO / PLAYER.NAVTO (goto), NAV.CHASE, NAV.PATROL (KCC entities only).
const (
	kccNavGoto uint8 = iota
	kccNavChase
	kccNavPatrol
)

type kccNavState struct {
	mode   uint8
	active bool
	tx, tz float64
	speed  float64
	// arrival: stop radius for goto/patrol; chase uses chaseGap as target standoff.
	arrival, brake float64

	chaseTarget int64
	chaseGap    float64

	patrolAX, patrolAZ, patrolBX, patrolBZ float64
	patrolToB                              bool
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
	if m.kccNav == nil {
		m.kccNav = make(map[int64]*kccNavState)
	}
	if m.kccLastGroundedAt == nil {
		m.kccLastGroundedAt = make(map[int64]float64)
	}
	if m.swimManual == nil {
		m.swimManual = make(map[int64]bool)
	}
}

// BindWater wires the water module for PLAYER.ISSWIMMING (optional).
func (m *Module) BindWater(w *mwater.Module) { m.water = w }

// Bind wires character controller + entity modules (see compiler pipeline wirePlayerModules).
func (m *Module) Bind(char *mbcharcontroller.Module, ent *mbentity.Module) {
	m.char = char
	m.ent = ent
	mbentity.SetKinematicCharacterLookup(func(id int64) bool {
		if m.entToChar != nil {
			if _, ok := m.entToChar[id]; ok {
				return true
			}
		}
		return false
	})
	mbentity.SetCharacterGroundNormalResolver(func(id int64) (float64, float64, float64, bool) {
		if m.char != nil {
			if h, ok := m.entToChar[id]; ok {
				nx, ny, nz, ok2 := m.char.CharacterGroundNormal(h)
				if ok2 {
					return nx, ny, nz, true
				}
			}
		}
		return 0, 0, 0, false
	})

	mbphysics3d.SetPhysicsKCCFanIn(func(ph *mbphysics3d.Module) {
		if m.char == nil || ph == nil || m.entToChar == nil {
			return
		}
		for eid, ch := range m.entToChar {
			if eid < 1 {
				continue
			}
			ev := m.char.DrainCharacterContactsForFanIn(ch)
			if len(ev) == 0 {
				continue
			}
			ph.FanInCharacterContactEvents(eid, ev)
		}
	})
}

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	registerPlayerCommands(m, reg)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {
	mbentity.SetKinematicCharacterLookup(nil)
	mbentity.SetCharacterGroundNormalResolver(nil)
	mbphysics3d.SetPhysicsKCCFanIn(nil)
}

func (m *Module) Reset() {}

