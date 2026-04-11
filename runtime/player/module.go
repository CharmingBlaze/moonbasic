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
	// stepHeight mirrors CharacterVirtual ExtendedUpdate WalkStairsStepUp (Y) for PLAYER.SETSTEPHEIGHT.
	stepHeight map[int64]float64
	grab       map[int64]int64   // player entity# -> grabbed entity# (0 = none)
	fovKick    map[int64]float64 // degrees added to camera FOV (read via GETFOVKICK; apply in script or future hook)

	// hostKCC: non-Jolt KCC backing (Windows/macOS CGO). Linux+Jolt uses entToChar + mbcharcontroller.
	hostKCC map[int64]*hostKCCState

	// kccNav: click-to-move / NAV targets for PLAYER.CREATE entities (see CHAR.NAVTO / CHAR.NAVUPDATE).
	kccNav map[int64]*kccNavState
	// kccLastGroundedAt: wall time (seconds) when CharacterVirtual last reported grounded — for optional coyote grace in ISGROUNDED.
	kccLastGroundedAt map[int64]float64
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

// hostKCCState is a software kinematic solver used when Jolt CharacterVirtual is unavailable
// (Windows/macOS fullruntime). Same script surface as Linux+Jolt; collision is approximate (static boxes).
//
// rad: horizontal capsule radius (same as CHAR.MAKE / MODEL.CREATECAPSULE radius).
// hei: total capsule height; pivot is the capsule center, so feet are hei/2 below the pivot (matches primitive draw).
type hostKCCState struct {
	rad, hei        float64
	stepH, slopeDeg float64
	stickDown, pad  float64
	gravityScale    float64
	mass            float64
	vx, vy, vz      float64
	crouch          bool
	swimBuoy        float64
	swimDrag        float64
	swimOn          bool
	grounded        bool
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
	if m.hostKCC == nil {
		m.hostKCC = make(map[int64]*hostKCCState)
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
		if m.hostKCC != nil {
			_, ok := m.hostKCC[id]
			return ok
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
		if m.hostKCC != nil {
			if st, ok := m.hostKCC[id]; ok && st.grounded {
				return 0, 1, 0, true
			}
		}
		return 0, 0, 0, false
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
}

func (m *Module) Reset() {}

