//go:build (linux || windows) && cgo

package mbcharcontroller

import (
	"fmt"
	"math"
	"sync"

	"github.com/bbitechnologies/jolt-go/jolt"

	"moonbasic/runtime"
	mbnav "moonbasic/runtime/mbnav"
	mbphysics3d "moonbasic/runtime/physics3d"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// Handles still alive at shutdown must be freed before PHYSICS3D.STOP destroys Jolt.
var charTrackMu sync.Mutex
var charTracked = map[heap.Handle]struct{}{}

func trackChar(h heap.Handle) {
	charTrackMu.Lock()
	charTracked[h] = struct{}{}
	charTrackMu.Unlock()
}

func untrackChar(h heap.Handle) {
	charTrackMu.Lock()
	delete(charTracked, h)
	charTrackMu.Unlock()
}

// charVerticalSleepVy: when Jolt reports OnGround and |vy| is below this, zero vy (stationary bounce).
const charVerticalSleepVy = float32(0.25)

// groundedMaxDownVy: when OnGround, do not integrate past this downward speed (reduces gravity-induced floor separation).
const groundedMaxDownVy = float32(0.1)

type charObj struct {
	cv      *jolt.CharacterVirtual
	release heap.ReleaseOnce

	// extUpdate is passed to CharacterVirtual.ExtendedUpdate (StickToFloor / WalkStairs tuning).
	extUpdate jolt.ExtendedUpdateSettings

	gravityScale float32 // 1 = default world gravity response on Y
	swimMode     bool
	swimBuoyancy float32 // reduces downward acceleration when swimming (0..1 typical)
	swimDrag     float32 // horizontal velocity damping per second when swimming
	crouch       bool    // when true, capsule uses standCapFullH * crouchHeightMul (feet-anchored rebuild)
	charMass     float32 // stored for gameplay; Jolt CharacterVirtual mass is fixed at create

	// Capsule + padding (for rebuild after PLAYER.SETSLOPE / CHAR.SETPADDING).
	capRadius       float32
	capFullH        float32 // total capsule height passed to CreateCapsule
	standCapFullH   float32 // standing height (PLAYER.CREATE); used to restore after crouch
	crouchHeightMul float32 // crouch height = standCapFullH * mul (default 0.55)
	maxSlopeDeg     float32
	charPad         float32 // CharacterVirtualSettings.CharacterPadding
	charFriction    float32 // gameplay scalar (CHARACTERREF.SETFRICTION); not Jolt body friction
	charBounce      float32 // gameplay scalar (CHARACTERREF.SETBOUNCE)
	maxStrength     float32 // New: Max force character can push with
	backoffDist     float32 // New: PredictiveContactDistance

	// Jump buffer (seconds): if jump pressed in air, apply impulse on landing before deadline (sim time).
	jumpBufferSec   float32
	jumpPendingImp  float32
	jumpDeadlineSim float64 // sim-time deadline; 0 = none

	// Coyote time (seconds, physics clock): jump allowed briefly after IsSupported becomes false.
	coyoteSec        float32
	lastSupportedSim float64 // mbphysics3d sim time when IsSupported was last true; -1 = unset

	// Horizontal control multipliers (PLAYER.SETAIRCONTROL / SETGROUNDCONTROL).
	airControl    float32
	groundControl float32

	touchCeiling bool // last move: head/capsule contact with downward-facing surface
}

func (m *Module) kccDt() float32 {
	if m == nil || m.h == nil {
		return float32(1.0 / 60.0)
	}
	ph := mbphysics3d.GetModule(m.h)
	if ph == nil {
		return float32(1.0 / 60.0)
	}
	return float32(ph.FixedStepSeconds())
}

func (m *Module) simTimeSeconds() float64 {
	if m == nil || m.h == nil {
		return 0
	}
	ph := mbphysics3d.GetModule(m.h)
	if ph == nil {
		return 0
	}
	return ph.SimTimeSeconds()
}

func (m *Module) updateKCCSupportedClock(co *charObj) {
	if co == nil || co.cv == nil {
		return
	}
	st := m.simTimeSeconds()
	if co.cv.IsSupported() {
		co.lastSupportedSim = st
	}
}

func (c *charObj) TypeName() string { return "CharController" }

func (c *charObj) TypeTag() uint16 { return heap.TagCharController }

func (c *charObj) Free() {
	c.release.Do(func() {
		if c.cv != nil {
			c.cv.Destroy()
			c.cv = nil
		}
	})
}

func registerCharControllerCommands(m *Module, reg runtime.Registrar) {
	reg.Register("CHARACTER.MAKE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) == 3 && args[0].Kind != value.KindHandle {
			return ccMake(m, args)
		}
		return ccMakeLegacy(m, args)
	}))
	reg.Register("CHARACTER.CREATE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) == 1 || len(args) == 3 {
			return ccMakeLegacy(m, args)
		}
		return ccMake(m, args)
	}))
	// CHARCONTROLLER.* — Jolt CharacterVirtual helpers (same behavior as cc* below).
	reg.Register("CHARCONTROLLER.CREATE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccMake(m, args) }))
	reg.Register("CHARCONTROLLER.MAKE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccMake(m, args) }))
	reg.Register("CHARCONTROLLER.SETPOS", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccSetPos(m, args) }))
	reg.Register("CHARCONTROLLER.SETPOSITION", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccSetPos(m, args) }))
	reg.Register("CHARCONTROLLER.GETPOS", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccGetPos(m, args) }))
	reg.Register("CHARCONTROLLER.MOVE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccMove(m, args) }))
	reg.Register("CHARCONTROLLER.ISGROUNDED", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccIsGrounded(m, args) }))
	reg.Register("CHARCONTROLLER.X", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccAxis(m, args, 0) }))
	reg.Register("CHARCONTROLLER.Y", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccAxis(m, args, 1) }))
	reg.Register("CHARCONTROLLER.Z", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccAxis(m, args, 2) }))
	reg.Register("CHARCONTROLLER.FREE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccFree(m, args) }))
	reg.Register("CHARCONTROLLER.GETLINEARVEL", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccGetLinearVel(m, args) }))
	reg.Register("CHARCONTROLLER.GETGROUNDVELOCITY", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccGetGroundVel(m, args) }))
	reg.Register("CHARCONTROLLER.GETGROUNDNORMAL", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccGetGroundNormalArr(m, args) }))
	reg.Register("CHARCONTROLLER.GROUNDSTATE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccGroundState(m, args) }))
	reg.Register("CHARCONTROLLER.TELEPORT", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccTeleportScript(m, args) }))

	// CONTROLLER.* — short aliases (same handlers as CHARCONTROLLER.*).
	reg.Register("CONTROLLER.CREATE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccMake(m, args) }))
	reg.Register("CONTROLLER.MAKE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccMake(m, args) }))
	reg.Register("CONTROLLER.MOVE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccMove(m, args) }))
	reg.Register("CONTROLLER.GROUNDED", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccIsGrounded(m, args) }))
	reg.Register("CONTROLLER.JUMP", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccJump(m, args) }))
	reg.Register("CONTROLLER.FREE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccFree(m, args) }))

	reg.Register("CHARACTERREF.SETPOS", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccSetPos(m, args) }))
	reg.Register("CHARACTERREF.SETPOSITION", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccSetPos(m, args) }))
	reg.Register("CHARACTERREF.MOVE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccMove(m, args) }))
	reg.Register("CHARACTERREF.SETVELOCITY", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccSetVel(m, args) }))
	reg.Register("CHARACTERREF.ADDVELOCITY", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccAddVel(m, args) }))
	reg.Register("CHARACTERREF.JUMP", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccJump(m, args) }))
	reg.Register("CHARACTERREF.UPDATE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccUpdate(m, args) }))
	reg.Register("CHARACTERREF.ISGROUNDED", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccIsGrounded(m, args) }))
	reg.Register("CHARACTERREF.ONSLOPE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccOnSlope(m, args) }))
	reg.Register("CHARACTERREF.ONWALL", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccOnWall(m, args) }))
	reg.Register("CHARACTERREF.GETSLOPEANGLE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccSlopeAngle(m, args) }))
	reg.Register("CHARACTERREF.GETSPEED", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccGetSpeed(m, args) }))
	reg.Register("CHARACTERREF.SETGRAVITY", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccSetGravity(m, args) }))
	reg.Register("CHARACTERREF.SETFRICTION", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccSetFriction(m, args) }))
	reg.Register("CHARACTERREF.SETMAXSLOPE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccSetMaxSlope(m, args) }))
	reg.Register("CHARACTERREF.SETSTEPHEIGHT", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccSetStepHeight(m, args) }))
	reg.Register("CHARACTERREF.SETSNAPDISTANCE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccSetSnapDist(m, args) }))
	reg.Register("CHARACTERREF.SETSETTING", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccSetSetting(m, args) }))
	reg.Register("CHARACTERREF.SETCONTACTLISTENER", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccSetContactListener(m, args) }))
	reg.Register("CHARACTERREF.DRAINCONTACTS", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccDrainContacts(m, args) }))
	reg.Register("CHARACTERREF.FREE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccFree(m, args) }))
	reg.Register("CHARACTERREF.GETPOSITION", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccGetPos(m, args) }))
	reg.Register("CHARACTERREF.GETROT", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccGetRot(m, args) }))
	reg.Register("CHARACTERREF.SETJUMPBUFFER", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccSetJumpBuffer(m, args) }))
	reg.Register("CHARACTERREF.SETAIRCONTROL", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccSetAirControl(m, args) }))
	reg.Register("CHARACTERREF.SETGROUNDCONTROL", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccSetGroundControl(m, args) }))
	reg.Register("CHARACTERREF.GETISSLIDING", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccGetIsSliding(m, args) }))
	reg.Register("CHARACTERREF.GETCEILING", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccGetCeiling(m, args) }))
	reg.Register("CHARACTERREF.GETGROUNDVELOCITY", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccGetGroundVel(m, args) }))

	reg.Register("CHARACTERREF.GETGRAVITY", "charcontroller", runtime.AdaptLegacy(m.ccGetGravity))
	reg.Register("CHARACTERREF.GETMAXSLOPE", "charcontroller", runtime.AdaptLegacy(m.ccGetMaxSlope))
	reg.Register("CHARACTERREF.GETSTEPHEIGHT", "charcontroller", runtime.AdaptLegacy(m.ccGetStepHeight))
	reg.Register("CHARACTERREF.GETSNAPDISTANCE", "charcontroller", runtime.AdaptLegacy(m.ccGetSnapDist))
	reg.Register("CHARACTERREF.GETFRICTION", "charcontroller", runtime.AdaptLegacy(m.ccGetFriction))
	reg.Register("CHARACTERREF.GETJUMPBUFFER", "charcontroller", runtime.AdaptLegacy(m.ccGetJumpBuffer))
	reg.Register("CHARACTERREF.GETAIRCONTROL", "charcontroller", runtime.AdaptLegacy(m.ccGetAirControl))
	reg.Register("CHARACTERREF.GETGROUNDCONTROL", "charcontroller", runtime.AdaptLegacy(m.ccGetGroundControl))
}

func (m *Module) ccGetGravity(args []value.Value) (value.Value, error) { return ccGetGravity(m, args) }
func (m *Module) ccGetMaxSlope(args []value.Value) (value.Value, error) { return ccGetMaxSlope(m, args) }
func (m *Module) ccGetStepHeight(args []value.Value) (value.Value, error) {
	return ccGetStepHeight(m, args)
}
func (m *Module) ccGetSnapDist(args []value.Value) (value.Value, error)    { return ccGetSnapDist(m, args) }
func (m *Module) ccGetFriction(args []value.Value) (value.Value, error)    { return ccGetFriction(m, args) }
func (m *Module) ccGetJumpBuffer(args []value.Value) (value.Value, error)  { return ccGetJumpBuffer(m, args) }
func (m *Module) ccGetAirControl(args []value.Value) (value.Value, error)  { return ccGetAirControl(m, args) }
func (m *Module) ccGetGroundControl(args []value.Value) (value.Value, error) {
	return ccGetGroundControl(m, args)
}

func ccSetJumpBuffer(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil || len(args) < 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETJUMPBUFFER: need handle, seconds#")
	}
	sec, ok := args[1].ToFloat()
	if !ok || sec < 0 {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETJUMPBUFFER: seconds must be >= 0")
	}
	if err := m.SetCharacterJumpBuffer(heap.Handle(args[0].IVal), sec); err != nil {
		return value.Nil, err
	}
	return args[0], nil
}

func ccSetAirControl(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil || len(args) < 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETAIRCONTROL: need handle, scale#")
	}
	s, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETAIRCONTROL: scale must be numeric")
	}
	if err := m.SetCharacterAirControl(heap.Handle(args[0].IVal), s); err != nil {
		return value.Nil, err
	}
	return args[0], nil
}

func ccSetGroundControl(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil || len(args) < 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETGROUNDCONTROL: need handle, scale#")
	}
	s, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETGROUNDCONTROL: scale must be numeric")
	}
	if err := m.SetCharacterGroundControl(heap.Handle(args[0].IVal), s); err != nil {
		return value.Nil, err
	}
	return args[0], nil
}

func ccGetIsSliding(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil || len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARACTERREF.GETISSLIDING: need handle")
	}
	return value.FromBool(m.CharacterIsSliding(heap.Handle(args[0].IVal))), nil
}

func ccGetCeiling(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil || len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARACTERREF.GETCEILING: need handle")
	}
	return value.FromBool(m.CharacterTouchingCeiling(heap.Handle(args[0].IVal))), nil
}

func (m *Module) extendedUpdateChar(co *charObj, dt float32) {
	if co == nil || co.cv == nil {
		return
	}
	g := mbphysics3d.GravityVec()
	co.cv.ExtendedUpdate(dt, g, &co.extUpdate)
	m.characterVerticalSleep(co)
}

func (m *Module) characterVerticalSleep(co *charObj) {
	if co == nil || co.cv == nil {
		return
	}
	if co.cv.GetGroundState() != jolt.GroundStateOnGround {
		return
	}
	v := co.cv.GetLinearVelocity()
	if math.Abs(float64(v.Y)) >= float64(charVerticalSleepVy) {
		return
	}
	v.Y = 0
	co.cv.SetLinearVelocity(v)
}

// SetCharacterRestitution stores a gameplay bounciness factor (CHARACTERREF.SETBOUNCE / player bridge).
func (m *Module) SetCharacterRestitution(h heap.Handle, bounce float64) error {
	if m.h == nil {
		return fmt.Errorf("heap not bound")
	}
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return fmt.Errorf("invalid character handle")
	}
	if bounce < 0 {
		bounce = 0
	}
	if bounce > 2 {
		bounce = 2
	}
	co.charBounce = float32(bounce)
	return nil
}

// CharacterRestitution returns gameplay bounciness set via CHARACTERREF.SETBOUNCE/SETBOUNCINESS.
func (m *Module) CharacterRestitution(h heap.Handle) (float64, bool) {
	if m.h == nil {
		return 0, false
	}
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return 0, false
	}
	return float64(co.charBounce), true
}

// SetCharacterFriction stores a gameplay friction factor (CHARACTERREF.SETFRICTION / player bridge).
func (m *Module) SetCharacterFriction(h heap.Handle, friction float64) error {
	if m.h == nil {
		return fmt.Errorf("heap not bound")
	}
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return fmt.Errorf("invalid character handle")
	}
	if friction < 0 {
		friction = 0
	}
	if friction > 2 {
		friction = 2
	}
	co.charFriction = float32(friction)
	return nil
}

func shutdownCharController(m *Module) {
	if m.h == nil {
		return
	}
	charTrackMu.Lock()
	hs := make([]heap.Handle, 0, len(charTracked))
	for h := range charTracked {
		hs = append(hs, h)
	}
	charTracked = make(map[heap.Handle]struct{})
	charTrackMu.Unlock()
	for _, h := range hs {
		m.h.Free(h)
	}
}

func ccMake(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MAKE: heap not bound")
	}
	if len(args) < 5 {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MAKE: need radius, height, x, y, z")
	}
	radius, ok := args[0].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MAKE: radius must be numeric")
	}
	height, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MAKE: height must be numeric")
	}
	x, ok := args[2].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MAKE: x must be numeric")
	}
	y, ok := args[3].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MAKE: y must be numeric")
	}
	z, ok := args[4].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MAKE: z must be numeric")
	}
	h, err := createCharacterVirtualFromParams(m, radius, height, x, y, z, 50, 0.02, 100.0, 0.05)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

// createCharacterVirtualFromParams allocates a CharacterVirtual. maxSlopeDeg ≤ 0 uses 50°.
// padding ≤ 0 uses Jolt default CharacterPadding (0.02).
// strength ≤ 0 uses 100. padding2 ≤ 0 uses 0.05.
func createCharacterVirtualFromParams(m *Module, radius, height, x, y, z, maxSlopeDeg, padding, strength, padding2 float64) (heap.Handle, error) {
	ps := mbphysics3d.ActiveJoltPhysics()
	if ps == nil {
		return 0, fmt.Errorf("CHARCONTROLLER: PHYSICS3D not started")
	}
	if maxSlopeDeg <= 0 {
		maxSlopeDeg = 45
	}
	fr, fh := float32(radius), float32(height)
	hh := fh/2 - fr
	if hh < 0.05 {
		hh = 0.05
	}
	capsule := jolt.CreateCapsule(hh, fr)
	settings := jolt.NewCharacterVirtualSettings(capsule)
	if maxSlopeDeg <= 0 {
		maxSlopeDeg = 50
	}
	settings.MaxSlopeAngle = jolt.DegreesToRadians(float32(maxSlopeDeg))
	settings.MaxStrength = 500.0 // Professional Architect Default
	if strength > 0 {
		settings.MaxStrength = float32(strength)
	}
	// Bug Fix: Explicitly set recovery speed and contact distance to ensure zero-bounce
	settings.PenetrationRecoverySpeed = 1.0
	settings.PredictiveContactDistance = 0.1 // Professional Architect Default
	if padding2 > 0 {
		settings.PredictiveContactDistance = float32(padding2)
	}
	if padding > 0 {
		settings.CharacterPadding = float32(padding)
	}
	cv := ps.CreateCharacterVirtual(settings, jolt.Vec3{X: float32(x), Y: float32(y), Z: float32(z)})
	cv.SetContactListenerEnabled(true)
	ext := jolt.DefaultExtendedUpdateSettings()
	h, err := m.h.Alloc(&charObj{
		cv: cv, extUpdate: ext, gravityScale: 1, charMass: 70, charFriction: 0.5, charBounce: 0,
		capRadius: fr, capFullH: fh, standCapFullH: fh, crouchHeightMul: 0.55, maxSlopeDeg: float32(maxSlopeDeg), charPad: settings.CharacterPadding,
		maxStrength: settings.MaxStrength, backoffDist: settings.PredictiveContactDistance,
		airControl: 1, groundControl: 1,
		jumpBufferSec: 0.1, coyoteSec: 0.1, lastSupportedSim: -1,
	})
	if err != nil {
		if cv != nil {
			cv.Destroy()
		}
		return 0, err
	}
	trackChar(h)
	return h, nil
}

func ccSetPos(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.SETPOS: heap not bound")
	}
	if len(args) < 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.SETPOS: need handle, x, y, z")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil || co.cv == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.SETPOS: invalid handle")
	}
	x, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.SETPOS: x must be numeric")
	}
	y, ok := args[2].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.SETPOS: y must be numeric")
	}
	z, ok := args[3].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.SETPOS: z must be numeric")
	}
	co.cv.SetPosition(jolt.Vec3{X: float32(x), Y: float32(y), Z: float32(z)})
	m.extendedUpdateChar(co, m.kccDt())
	return args[0], nil
}

func ccMove(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MOVE: heap not bound")
	}
	if len(args) < 4 {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MOVE: need handle, dx, dy, dz")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MOVE: first arg must be handle")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if co.cv == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MOVE: invalid handle")
	}
	dx, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MOVE: dx must be numeric")
	}
	dy, ok := args[2].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MOVE: dy must be numeric")
	}
	dz, ok := args[3].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MOVE: dz must be numeric")
	}
	p := co.cv.GetPosition()
	co.cv.SetPosition(jolt.Vec3{
		X: p.X + float32(dx),
		Y: p.Y + float32(dy),
		Z: p.Z + float32(dz),
	})
	m.extendedUpdateChar(co, m.kccDt())
	return args[0], nil
}

func ccAxis(m *Module, args []value.Value, axis int) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER axis: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.FromFloat(0), nil
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil || co.cv == nil {
		return value.FromFloat(0), nil
	}
	p := co.cv.GetPosition()
	switch axis {
	case 0:
		return value.FromFloat(float64(p.X)), nil
	case 1:
		return value.FromFloat(float64(p.Y)), nil
	default:
		return value.FromFloat(float64(p.Z)), nil
	}
}

func ccIsGrounded(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.ISGROUNDED: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.ISGROUNDED: need handle")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if co.cv == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.ISGROUNDED: invalid handle")
	}
	if co.cv.IsSupported() {
		return value.FromBool(true), nil
	}
	return value.FromBool(false), nil
}

func ccGetPos(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETPOS: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETPOS: need handle")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if co.cv == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETPOS: invalid handle")
	}
	p := co.cv.GetPosition()
	arr, err := heap.NewArray([]int64{3})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, float64(p.X))
	_ = arr.Set([]int64{1}, float64(p.Y))
	_ = arr.Set([]int64{2}, float64(p.Z))
	ah, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(ah), nil
}

func ccGetRot(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARACTERREF.GETROT: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARACTERREF.GETROT: need handle")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if co.cv == nil {
		return value.Nil, fmt.Errorf("CHARACTERREF.GETROT: invalid handle")
	}
	v := co.cv.GetLinearVelocity()
	p, y, r := mbnav.EulerFromWorldDirection(float64(v.X), float64(v.Y), float64(v.Z))
	arr, err := heap.NewArray([]int64{3})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, p)
	_ = arr.Set([]int64{1}, y)
	_ = arr.Set([]int64{2}, r)
	ah, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(ah), nil
}

func ccFree(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.FREE: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.FREE: need handle")
	}
	hid := heap.Handle(args[0].IVal)
	untrackChar(hid)
	if err := m.h.Free(hid); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func ccAllocFloat3(m *Module, x, y, z float64) (value.Value, error) {
	arr, err := heap.NewArray([]int64{3})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, x)
	_ = arr.Set([]int64{1}, y)
	_ = arr.Set([]int64{2}, z)
	ah, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(ah), nil
}

func ccGetLinearVel(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETLINEARVEL: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETLINEARVEL: need handle")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil || co.cv == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETLINEARVEL: invalid handle")
	}
	v := co.cv.GetLinearVelocity()
	return ccAllocFloat3(m, float64(v.X), float64(v.Y), float64(v.Z))
}

func ccGetGroundVel(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETGROUNDVELOCITY: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETGROUNDVELOCITY: need handle")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil || co.cv == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETGROUNDVELOCITY: invalid handle")
	}
	v := co.cv.GetGroundVelocity()
	return ccAllocFloat3(m, float64(v.X), float64(v.Y), float64(v.Z))
}

func ccGetGroundNormalArr(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETGROUNDNORMAL: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETGROUNDNORMAL: need handle")
	}
	h := heap.Handle(args[0].IVal)
	nx, ny, nz, ok := m.CharacterGroundNormal(h)
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETGROUNDNORMAL: invalid handle")
	}
	return ccAllocFloat3(m, nx, ny, nz)
}

func ccGroundState(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GROUNDSTATE: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GROUNDSTATE: need handle")
	}
	gi, ok := m.CharacterGroundStateInt(heap.Handle(args[0].IVal))
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GROUNDSTATE: invalid handle")
	}
	return value.FromInt(int64(gi)), nil
}

func ccTeleportScript(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.TELEPORT: heap not bound")
	}
	if len(args) < 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.TELEPORT: need handle, x, y, z")
	}
	x, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.TELEPORT: x must be numeric")
	}
	y, ok := args[2].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.TELEPORT: y must be numeric")
	}
	z, ok := args[3].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.TELEPORT: z must be numeric")
	}
	if err := m.CharacterTeleport(heap.Handle(args[0].IVal), x, y, z); err != nil {
		return value.Nil, err
	}
	return args[0], nil
}

// AllocCharacter creates a Jolt CharacterVirtual (used by runtime/player).
// maxSlopeDeg ≤ 0 defaults to 45° (matches prior CHARCONTROLLER.MAKE behavior).
// padding ≤ 0 uses default skin width; >0 sets CharacterPadding.
func (m *Module) AllocCharacter(radius, height, x, y, z float64, maxSlopeDeg, padding, strength, padding2 float64) (heap.Handle, error) {
	return createCharacterVirtualFromParams(m, radius, height, x, y, z, maxSlopeDeg, padding, strength, padding2)
}

// RecreateCharacterWithSlope destroys a character and creates a new one at the same pose with a new max walk slope.
func (m *Module) RecreateCharacterWithSlope(oldH heap.Handle, radius, height, maxSlopeDeg float64) (heap.Handle, error) {
	co, err := heap.Cast[*charObj](m.h, oldH)
	if err != nil || co.cv == nil {
		return 0, fmt.Errorf("RecreateCharacterWithSlope: invalid handle")
	}
	p := co.cv.GetPosition()
	return m.RecreateCharacterWithSlopeAt(oldH, radius, height, maxSlopeDeg, float64(p.X), float64(p.Y), float64(p.Z))
}

// RecreateCharacterWithSlopeAt recreates the character at an explicit position (feet-anchored crouch uses a new center Y).
func (m *Module) RecreateCharacterWithSlopeAt(oldH heap.Handle, radius, height, maxSlopeDeg, atX, atY, atZ float64) (heap.Handle, error) {
	if m.h == nil {
		return 0, fmt.Errorf("RecreateCharacterWithSlopeAt: heap not bound")
	}
	co, err := heap.Cast[*charObj](m.h, oldH)
	if err != nil || co.cv == nil {
		return 0, fmt.Errorf("RecreateCharacterWithSlopeAt: invalid handle")
	}
	v := co.cv.GetLinearVelocity()
	extCopy := co.extUpdate
	gs := co.gravityScale
	sm := co.swimMode
	sb, sd := co.swimBuoyancy, co.swimDrag
	crouchOn := co.crouch
	cm := co.charMass
	cf := co.charFriction
	cb := co.charBounce
	ms := co.maxStrength
	bd := co.backoffDist
	standH := co.standCapFullH
	crouchMul := co.crouchHeightMul
	jb := co.jumpBufferSec
	jpi := co.jumpPendingImp
	jds := co.jumpDeadlineSim
	coy := co.coyoteSec
	lss := co.lastSupportedSim
	ac, gc := co.airControl, co.groundControl
	untrackChar(oldH)
	co.cv.Destroy()
	co.cv = nil
	if err := m.h.Free(oldH); err != nil {
		return 0, err
	}
	h, err := createCharacterVirtualFromParams(m, radius, height, atX, atY, atZ, maxSlopeDeg, float64(co.charPad), float64(ms), float64(bd))
	if err != nil {
		return 0, err
	}
	co2, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co2.cv == nil {
		return 0, fmt.Errorf("RecreateCharacterWithSlopeAt: internal")
	}
	co2.cv.SetLinearVelocity(v)
	co2.extUpdate = extCopy
	co2.gravityScale = gs
	co2.swimMode = sm
	co2.swimBuoyancy, co2.swimDrag = sb, sd
	co2.crouch = crouchOn
	co2.charMass = cm
	co2.charFriction = cf
	co2.charBounce = cb
	co2.maxStrength = ms
	co2.backoffDist = bd
	co2.standCapFullH = standH
	co2.crouchHeightMul = crouchMul
	co2.jumpBufferSec = jb
	co2.jumpPendingImp = jpi
	co2.jumpDeadlineSim = jds
	co2.coyoteSec = coy
	co2.lastSupportedSim = lss
	co2.airControl, co2.groundControl = ac, gc
	return h, nil
}

// CharacterGroundStateInt returns Jolt EGroundState as int: 0 OnGround, 1 OnSteepGround, 2 NotSupported, 3 InAir.
func (m *Module) CharacterGroundStateInt(h heap.Handle) (int, bool) {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return 0, false
	}
	return int(co.cv.GetGroundState()), true
}

// CharacterCapsuleDims returns stored capsule radius and total height (for rebuilds that must preserve PLAYER.CREATE sizing).
func (m *Module) CharacterCapsuleDims(h heap.Handle) (radius, fullHeight float64, ok bool) {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return 0, 0, false
	}
	return float64(co.capRadius), float64(co.capFullH), true
}

// SetCharacterLinearVelocity sets world linear velocity on the CharacterVirtual (used by CHARACTERREF.*).
func (m *Module) SetCharacterLinearVelocity(h heap.Handle, vx, vy, vz float64) error {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return fmt.Errorf("invalid character handle")
	}
	co.cv.SetLinearVelocity(jolt.Vec3{X: float32(vx), Y: float32(vy), Z: float32(vz)})
	return nil
}

// CharacterIntegrateStep applies gravity on Y for dt and runs ExtendedUpdate (CHARACTERREF.UPDATE after SETLINEARVELOCITY).
func (m *Module) CharacterIntegrateStep(h heap.Handle, dt float64) error {
	if m.h == nil {
		return fmt.Errorf("CharacterIntegrateStep: heap not bound")
	}
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return fmt.Errorf("invalid character handle")
	}
	dtUse := float64(m.kccDt())
	if dtUse <= 0 {
		return nil
	}
	_ = dt // script/frame delta ignored; KCC uses physics fixed step (same as PHYSICS3D.STEP)

	vel := co.cv.GetLinearVelocity()
	gst0 := co.cv.GetGroundState()
	if gst0 == jolt.GroundStateOnGround || gst0 == jolt.GroundStateOnSteepGround {
		gv := co.cv.GetGroundVelocity()
		vel.X += gv.X
		vel.Z += gv.Z
		vel.Y += gv.Y
	}
	g := mbphysics3d.GravityVec()
	gs := co.gravityScale
	if gs <= 0 {
		gs = 1
	}

	// Grounding & Slopes (Vertical Clamping)
	// If stationary and SlopeAngle < MaxSlopeAngle, lock horizontal and zero Y jitter.
	groundState := co.cv.GetGroundState()
	onGround := (groundState == jolt.GroundStateOnGround)

	if onGround {
		// Professional Grade: Vertical Clamping & Anti-Jitter
		horizSpeedSq := vel.X*vel.X + vel.Z*vel.Z
		if horizSpeedSq < 0.01 {
			vel.X = 0
			vel.Z = 0
		}
	} else if groundState == jolt.GroundStateOnSteepGround {
		// Professional Grade: Slope Clipping & Sliding Math
		// Specification: V_clipped = V - (V . N)N
		n := co.cv.GetGroundNormal()
		dot := vel.X*n.X + vel.Y*n.Y + vel.Z*n.Z
		if dot < 0 {
			vel.X -= dot * n.X
			vel.Y -= dot * n.Y
			vel.Z -= dot * n.Z
		}

		// Sliding Gravity: Add component of gravity along the steep slope plane
		// G_plane = G - (G . N)N
		gdot := g.X*n.X + g.Y*n.Y + g.Z*n.Z
		gx := (g.X - gdot*n.X) * gs
		gy := (g.Y - gdot*n.Y) * gs
		gz := (g.Z - gdot*n.Z) * gs

		vel.X += gx * float32(dtUse)
		vel.Y += gy * float32(dtUse)
		vel.Z += gz * float32(dtUse)
	}

	// Physics logic (Gravity + Swim)
	gy := g.Y * float32(dtUse) * gs
	if co.swimMode {
		// Blueprint: Decrease gravity and add linear damping
		b := co.swimBuoyancy
		if b > 1 {
			b = 1
		}
		gy *= (1.0 - b)
		d := float64(co.swimDrag)
		if d > 0 && dtUse > 0 {
			damp := math.Max(0, 1-d*dtUse)
			vel.X *= float32(damp)
			vel.Z *= float32(damp)
		}
	}
	vel.Y += gy
	if co.cv.GetGroundState() == jolt.GroundStateOnGround && vel.Y < -groundedMaxDownVy {
		vel.Y = -groundedMaxDownVy
	}
	co.cv.SetLinearVelocity(vel)
	m.extendedUpdateChar(co, float32(dtUse))
	m.tryConsumeBufferedJump(co)
	m.updateTouchCeiling(co)
	m.updateKCCSupportedClock(co)
	return nil
}

// CharacterSlopeAngleDegrees returns ground tilt from vertical (0 = flat), degrees, when on walkable ground.
func (m *Module) CharacterSlopeAngleDegrees(h heap.Handle) float64 {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return 0
	}
	if co.cv.GetGroundState() != jolt.GroundStateOnGround {
		return 0
	}
	n := co.cv.GetGroundNormal()
	return math.Acos(float64(n.Y)) * (180.0 / math.Pi)
}

// CharacterGameplayFriction returns stored gameplay friction (CHARACTERREF.SETFRICTION).
func (m *Module) CharacterGameplayFriction(h heap.Handle) float64 {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return 0
	}
	return float64(co.charFriction)
}

// CharacterMaxSlopeDegrees returns configured max walk slope (degrees).
func (m *Module) CharacterMaxSlopeDegrees(h heap.Handle) (float64, bool) {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return 0, false
	}
	return float64(co.maxSlopeDeg), true
}

// CharacterStepHeightY returns WalkStairsStepUp Y (stairs / curbs).
func (m *Module) CharacterStepHeightY(h heap.Handle) (float64, bool) {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return 0, false
	}
	return float64(co.extUpdate.WalkStairsStepUp.Y), true
}

// CharacterSnapDownDistance returns stick-to-floor max step down (positive world units).
func (m *Module) CharacterSnapDownDistance(h heap.Handle) (float64, bool) {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return 0, false
	}
	return float64(-co.extUpdate.StickToFloorStepDown.Y), true
}

// CharacterGravityScaleVal returns PLAYER.SETGRAVITYSCALE factor.
func (m *Module) CharacterGravityScaleVal(h heap.Handle) (float64, bool) {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return 0, false
	}
	return float64(co.gravityScale), true
}

// CharacterGroundNormal returns the ground surface normal under a CharacterVirtual (Jolt), or false if unavailable.
func (m *Module) CharacterGroundNormal(h heap.Handle) (nx, ny, nz float64, ok bool) {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return 0, 0, 0, false
	}
	n := co.cv.GetGroundNormal()
	return float64(n.X), float64(n.Y), float64(n.Z), true
}

// SetCharacterGravityScale scales accumulated gravity on Y (1 = default). Values <= 0 are treated as 1.
func (m *Module) SetCharacterGravityScale(h heap.Handle, scale float64) error {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return fmt.Errorf("invalid character handle")
	}
	co.gravityScale = float32(scale)
	if co.gravityScale <= 0 {
		co.gravityScale = 1
	}
	return nil
}

// SetCharacterSwim enables swim mode: buoyancy reduces downward gravity; drag damps horizontal velocity (per second).
func (m *Module) SetCharacterSwim(h heap.Handle, buoyancy, drag float64, on bool) error {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return fmt.Errorf("invalid character handle")
	}
	co.swimMode = on
	co.swimBuoyancy = float32(buoyancy)
	co.swimDrag = float32(drag)
	return nil
}

// SetCharacterPadding rebuilds the CharacterVirtual with a new CharacterPadding (collision skin).
func (m *Module) SetCharacterPadding(h heap.Handle, pad float32) (heap.Handle, error) {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return 0, fmt.Errorf("invalid character handle")
	}
	co.charPad = pad
	// Recreating with existing slope
	return m.RecreateCharacterWithSlope(h, float64(co.capRadius), float64(co.capFullH), float64(co.maxSlopeDeg))
}


// CharacterPadding returns the current CharacterPadding (skin width).
func (m *Module) CharacterPadding(h heap.Handle) (float32, bool) {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return 0, false
	}
	return co.charPad, true
}

func ccSetSetting(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETSETTING: heap not bound")
	}
	if len(args) < 3 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETSETTING: need handle, settingName$, value")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil || co.cv == nil {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETSETTING: invalid handle")
	}
	name := runtime.ArgString(args[1])
	val, ok := args[2].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETSETTING: value must be numeric")
	}

	switch name {
	case "MAXSTRENGTH":
		co.maxStrength = float32(val)
		// Refactor: We would need a Recreate or a Jolt API to update these live.
		// For now, we store it. If the blueprint implies live updates, we should rebuild.
		_, _ = m.RecreateCharacterWithSlope(heap.Handle(args[0].IVal), float64(co.capRadius), float64(co.capFullH), float64(co.maxSlopeDeg))
	case "BACKOFFDIST", "PREDICTIVEDISTANCE":
		co.backoffDist = float32(val)
		_, _ = m.RecreateCharacterWithSlope(heap.Handle(args[0].IVal), float64(co.capRadius), float64(co.capFullH), float64(co.maxSlopeDeg))
	case "PADDING":
		_, _ = m.SetCharacterPadding(heap.Handle(args[0].IVal), float32(val))
	case "MAXSLOPE":
		_, _ = m.RecreateCharacterWithSlope(heap.Handle(args[0].IVal), float64(co.capRadius), float64(co.capFullH), val)
	case "GRAVITY":
		co.gravityScale = float32(val)
	case "STEPHEIGHT":
		co.extUpdate.WalkStairsStepUp.Y = float32(val)
	case "SNAPDIST":
		co.extUpdate.StickToFloorStepDown.Y = -float32(val)
	case "SWIMMODE":
		co.swimMode = val != 0
	case "BUOYANCY":
		co.swimBuoyancy = float32(val)
	case "DRAG":
		co.swimDrag = float32(val)
	default:
		return value.Nil, fmt.Errorf("CHARACTERREF.SETSETTING: unknown setting '%s'", name)
	}

	return value.Nil, nil
}

func ccSetContactListener(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETCONTACTLISTENER: heap not bound")
	}
	if len(args) < 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETCONTACTLISTENER: need handle, enabled?")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil || co.cv == nil {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETCONTACTLISTENER: invalid handle")
	}
	co.cv.SetContactListenerEnabled(args[1].IVal != 0)
	return args[0], nil
}

func ccDrainContacts(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARACTERREF.DRAINCONTACTS: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARACTERREF.DRAINCONTACTS: need handle")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil || co.cv == nil {
		return value.Nil, fmt.Errorf("CHARACTERREF.DRAINCONTACTS: invalid handle")
	}

	events := co.cv.DrainContactQueue(256)
	if len(events) == 0 {
		return value.Nil, nil
	}

	// Return a 2D array: [[bodyB, px, py, pz, nx, ny, nz, dist], ...]
	arr, err := heap.NewArray([]int64{int64(len(events)), 8})
	if err != nil {
		return value.Nil, err
	}

	for i, e := range events {
		idx := int64(i)
		_ = arr.Set([]int64{idx, 0}, float64(e.BodyB))
		_ = arr.Set([]int64{idx, 1}, float64(e.Position.X))
		_ = arr.Set([]int64{idx, 2}, float64(e.Position.Y))
		_ = arr.Set([]int64{idx, 3}, float64(e.Position.Z))
		_ = arr.Set([]int64{idx, 4}, float64(e.Normal.X))
		_ = arr.Set([]int64{idx, 5}, float64(e.Normal.Y))
		_ = arr.Set([]int64{idx, 6}, float64(e.Normal.Z))
		_ = arr.Set([]int64{idx, 7}, float64(e.Distance))
	}

	ah, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(ah), nil
}

// SetCharacterWalkStairsStepUp sets ExtendedUpdateSettings.mWalkStairsStepUp (typically Y-only curb height).
func (m *Module) SetCharacterWalkStairsStepUp(h heap.Handle, stepUpY float32) error {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return fmt.Errorf("invalid character handle")
	}
	if stepUpY < 0 {
		stepUpY = 0
	}
	co.extUpdate.WalkStairsStepUp = jolt.Vec3{X: 0, Y: stepUpY, Z: 0}
	return nil
}

// SetCharacterStickToFloorDown sets how far down StickToFloor searches (positive distance → world -Y step vector).
func (m *Module) SetCharacterStickToFloorDown(h heap.Handle, down float32) error {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return fmt.Errorf("invalid character handle")
	}
	if down < 0 {
		down = 0
	}
	co.extUpdate.StickToFloorStepDown = jolt.Vec3{X: 0, Y: -down, Z: 0}
	return nil
}

// SetCharacterMass stores gameplay mass (used by PLAYER.Push scaling; Jolt capsule mass is fixed at create).
func (m *Module) SetCharacterMass(h heap.Handle, mass float64) error {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return fmt.Errorf("invalid character handle")
	}
	if mass < 0.001 {
		mass = 0.001
	}
	co.charMass = float32(mass)
	return nil
}

// characterCeilingBlocksStand returns true if a standing-height capsule would overlap geometry above (cannot un-crouch).
func (m *Module) characterCeilingBlocksStand(co *charObj) bool {
	if co == nil || co.cv == nil {
		return false
	}
	ps := co.cv.PhysicsSystem()
	if ps == nil {
		return false
	}
	p := co.cv.GetPosition()
	fr := co.capRadius
	standH := co.standCapFullH
	hh := standH*0.5 - fr
	if hh < 0.05 {
		hh = 0.05
	}
	capSh := jolt.CreateCapsule(hh, fr)
	defer capSh.Destroy()
	feetY := p.Y - co.capFullH*0.5
	cy := feetY + standH*0.5
	pos := jolt.Vec3{X: p.X, Y: cy, Z: p.Z}
	hits := ps.CollideShapeGetHits(capSh, pos, 16, 1e-3)
	curTop := p.Y + co.capFullH*0.5
	for _, hit := range hits {
		if hit.BodyID == nil {
			continue
		}
		if hit.ContactPoint.Y < feetY+0.05 {
			continue
		}
		if hit.ContactPoint.Y > curTop+0.02 {
			return true
		}
	}
	return false
}

// SetCharacterCrouch toggles crouch: rebuilds capsule at standCapFullH * crouchHeightMul (feet-anchored).
// Returns the handle to use (may differ from h after rebuild).
func (m *Module) SetCharacterCrouch(h heap.Handle, on bool) (heap.Handle, error) {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return 0, fmt.Errorf("invalid character handle")
	}
	if co.crouch == on {
		return h, nil
	}
	if !on && co.crouch && m.characterCeilingBlocksStand(co) {
		return h, nil
	}
	targetH := float64(co.standCapFullH)
	if on {
		targetH = float64(co.standCapFullH) * float64(co.crouchHeightMul)
	}
	minH := float64(co.capRadius)*2 + 0.05
	if targetH < minH {
		targetH = minH
	}
	p := co.cv.GetPosition()
	feetY := float64(p.Y) - float64(co.capFullH)*0.5
	newY := feetY + targetH*0.5
	newH, err := m.RecreateCharacterWithSlopeAt(h, float64(co.capRadius), targetH, float64(co.maxSlopeDeg), float64(p.X), newY, float64(p.Z))
	if err != nil {
		return 0, err
	}
	co2, err2 := heap.Cast[*charObj](m.h, newH)
	if err2 != nil || co2 == nil {
		return newH, nil
	}
	co2.crouch = on
	return newH, nil
}

// CharacterCrouch reports the stored crouch flag.
func (m *Module) CharacterCrouch(h heap.Handle) bool {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return false
	}
	return co.crouch
}

// CharacterMass returns stored gameplay mass (kg).
func (m *Module) CharacterMass(h heap.Handle) float64 {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return 70
	}
	return float64(co.charMass)
}

// SetCharacterJumpBuffer sets coyote-style jump buffer duration (seconds) for air presses.
func (m *Module) SetCharacterJumpBuffer(h heap.Handle, sec float64) error {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return fmt.Errorf("invalid character handle")
	}
	if sec < 0 {
		sec = 0
	}
	co.jumpBufferSec = float32(sec)
	return nil
}

// SetCharacterAirControl scales horizontal input while airborne (1 = default).
func (m *Module) SetCharacterAirControl(h heap.Handle, scale float64) error {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return fmt.Errorf("invalid character handle")
	}
	if scale < 0 {
		scale = 0
	}
	co.airControl = float32(scale)
	return nil
}

// SetCharacterGroundControl scales horizontal input while on ground (1 = default).
func (m *Module) SetCharacterGroundControl(h heap.Handle, scale float64) error {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return fmt.Errorf("invalid character handle")
	}
	if scale < 0 {
		scale = 0
	}
	co.groundControl = float32(scale)
	return nil
}

// CharacterTouchingCeiling reports whether the last move step saw a strong downward contact normal (head bump).
func (m *Module) CharacterTouchingCeiling(h heap.Handle) bool {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return false
	}
	return co.touchCeiling
}

// CharacterIsSliding reports Jolt steep-ground / sliding state.
func (m *Module) CharacterIsSliding(h heap.Handle) bool {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return false
	}
	return co.cv.GetGroundState() == jolt.GroundStateOnSteepGround
}

// CharacterGroundVelocityVec returns Jolt GetGroundVelocity (platform motion on ground plane).
func (m *Module) CharacterGroundVelocityVec(h heap.Handle) (vx, vy, vz float64, ok bool) {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return 0, 0, 0, false
	}
	v := co.cv.GetGroundVelocity()
	return float64(v.X), float64(v.Y), float64(v.Z), true
}

// CharacterTeleport snaps the capsule to a world position and clears velocity (no interpolation smoothing).
func (m *Module) CharacterTeleport(h heap.Handle, x, y, z float64) error {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return fmt.Errorf("invalid character handle")
	}
	co.cv.SetPosition(jolt.Vec3{X: float32(x), Y: float32(y), Z: float32(z)})
	co.cv.SetLinearVelocity(jolt.Vec3{})
	m.extendedUpdateChar(co, m.kccDt())
	return nil
}

func (m *Module) tryConsumeBufferedJump(co *charObj) {
	if co == nil || co.cv == nil || co.jumpPendingImp <= 0 {
		return
	}
	st := m.simTimeSeconds()
	if co.jumpBufferSec > 0 && co.jumpDeadlineSim > 0 && st > co.jumpDeadlineSim {
		co.jumpPendingImp = 0
		co.jumpDeadlineSim = 0
		return
	}
	if !co.cv.IsSupported() {
		return
	}
	v := co.cv.GetLinearVelocity()
	v.Y += co.jumpPendingImp
	co.cv.SetLinearVelocity(v)
	co.jumpPendingImp = 0
	co.jumpDeadlineSim = 0
}

func (m *Module) updateTouchCeiling(co *charObj) {
	if co == nil {
		return
	}
	co.touchCeiling = false
	if co.cv == nil {
		return
	}
	for _, c := range co.cv.GetActiveContacts(48) {
		if c.ContactNormal.Y < -0.35 {
			co.touchCeiling = true
			return
		}
	}
}

// CharacterMoveXZVelocity sets horizontal world velocity (units/s) and integrates for dt seconds (gravity on Y).
func (m *Module) CharacterMoveXZVelocity(h heap.Handle, vx, vz, dt float64) error {
	if m.h == nil {
		return fmt.Errorf("CharacterMoveXZVelocity: heap not bound")
	}
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return fmt.Errorf("invalid character handle")
	}
	dtUse := float64(m.kccDt())
	_ = dt
	ac := co.airControl
	gc := co.groundControl
	if ac <= 0 {
		ac = 1
	}
	if gc <= 0 {
		gc = 1
	}
	gst := co.cv.GetGroundState()
	scale := ac
	if gst == jolt.GroundStateOnGround {
		scale = gc
	}
	inX := float32(vx) * scale
	inZ := float32(vz) * scale
	vel := co.cv.GetLinearVelocity()
	vel.X = inX
	vel.Z = inZ
	if gst == jolt.GroundStateOnGround || gst == jolt.GroundStateOnSteepGround {
		gv := co.cv.GetGroundVelocity()
		vel.X += gv.X
		vel.Z += gv.Z
		vel.Y += gv.Y
	}
	g := mbphysics3d.GravityVec()
	gs := co.gravityScale
	if gs <= 0 {
		gs = 1
	}
	gy := g.Y * float32(dtUse) * gs
	if co.swimMode {
		b := co.swimBuoyancy
		if b < 0 {
			b = 0
		}
		if b > 1 {
			b = 1
		}
		gy *= (1.0 - b)
		d := float64(co.swimDrag)
		if d > 0 && dtUse > 0 {
			damp := math.Max(0, 1-d*dtUse)
			vel.X *= float32(damp)
			vel.Z *= float32(damp)
		}
	}
	vel.Y += gy
	if co.cv.GetGroundState() == jolt.GroundStateOnGround && vel.Y < -groundedMaxDownVy {
		vel.Y = -groundedMaxDownVy
	}
	co.cv.SetLinearVelocity(vel)
	m.extendedUpdateChar(co, float32(dtUse))
	m.tryConsumeBufferedJump(co)
	m.updateTouchCeiling(co)
	m.updateKCCSupportedClock(co)
	return nil
}

// CharacterJump adds upward impulse to linear velocity Y (used by runtime/player).
func (m *Module) CharacterJump(h heap.Handle, impulseY float64) error {
	if m.h == nil {
		return fmt.Errorf("CharacterJump: heap not bound")
	}
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return fmt.Errorf("invalid character handle")
	}
	st := m.simTimeSeconds()
	coy := float64(co.coyoteSec)
	if coy <= 0 {
		coy = 0.1
	}
	supported := co.cv.IsSupported()
	if !supported && co.lastSupportedSim >= 0 && st-co.lastSupportedSim <= coy {
		supported = true
	}
	if supported {
		v := co.cv.GetLinearVelocity()
		v.Y += float32(impulseY)
		co.cv.SetLinearVelocity(v)
		co.jumpPendingImp = 0
		co.jumpDeadlineSim = 0
		m.extendedUpdateChar(co, m.kccDt())
		return nil
	}
	if co.jumpBufferSec > 0 {
		co.jumpPendingImp = float32(impulseY)
		co.jumpDeadlineSim = st + float64(co.jumpBufferSec)
		return nil
	}
	return nil
}

// CharacterIsGrounded reports Jolt ground support for the character.
func (m *Module) CharacterIsGrounded(h heap.Handle) (bool, error) {
	v, err := ccIsGrounded(m, []value.Value{value.FromHandle(int32(h))})
	if err != nil {
		return false, err
	}
	if v.Kind == value.KindBool {
		return v.IVal != 0, nil
	}
	return false, nil
}

// CharacterPosition returns the character capsule position.
func (m *Module) CharacterPosition(h heap.Handle) (x, y, z float64, ok bool) {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return 0, 0, 0, false
	}
	p := co.cv.GetPosition()
	return float64(p.X), float64(p.Y), float64(p.Z), true
}

// CharacterLinearVelocity returns current linear velocity (for PLAYER.SYNCANIM).
func (m *Module) CharacterLinearVelocity(h heap.Handle) (vx, vy, vz float64, ok bool) {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return 0, 0, 0, false
	}
	v := co.cv.GetLinearVelocity()
	return float64(v.X), float64(v.Y), float64(v.Z), true
}

// FreeCharacter destroys a character (same as CHARCONTROLLER.FREE).
func (m *Module) FreeCharacter(h heap.Handle) error {
	_, err := ccFree(m, []value.Value{value.FromHandle(int32(h))})
	return err
}

// DrainCharacterContactsForFanIn drains pending CharacterVirtual contacts (same data as CHARACTERREF.DRAINCONTACTS) for physics ONCOLLISION fan-in.
func (m *Module) DrainCharacterContactsForFanIn(h heap.Handle) []jolt.CharacterContactEvent {
	if m.h == nil {
		return nil
	}
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return nil
	}
	return co.cv.DrainContactQueue(256)
}

// LinkCharacterToEntity is reserved for wiring a CharacterVirtual handle to a scene entity id (visual sync hooks).
// Currently a no-op; gameplay uses existing PLAYER/CHARACTERREF paths when the capsule is created with an entity.
func (m *Module) LinkCharacterToEntity(h heap.Handle, entityID int64) error {
	_, _ = h, entityID
	return nil
}

func ccMakeLegacy(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 1 {
		return value.Nil, fmt.Errorf("CHARACTER.CREATE(entity) needs at least 1 arg")
	}
	r := 0.4
	h_val := 1.0
	if len(args) >= 3 {
		r, _ = args[1].ToFloat()
		h_val, _ = args[2].ToFloat()
	}
	// Legacy call: default strength 100, backoff 0.05
	h, err := createCharacterVirtualFromParams(m, r, h_val, 0, 0, 0, 50, 0.02, 100.0, 0.05)
	if err != nil {
		return value.Nil, err
	}
	// Bind to entity if arg0 is present
	eid, _ := args[0].ToInt()
	if eid > 0 {
		_ = m.LinkCharacterToEntity(h, eid)
	}
	return value.FromHandle(h), nil
}

func ccSetVel(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 4 {
		return value.Nil, fmt.Errorf("CHARACTER.SETVELOCITY needs handle, vx, vy, vz")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	vx, _ := args[1].ToFloat()
	vy, _ := args[2].ToFloat()
	vz, _ := args[3].ToFloat()
	co.cv.SetLinearVelocity(jolt.Vec3{X: float32(vx), Y: float32(vy), Z: float32(vz)})
	return args[0], nil
}

func ccAddVel(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 4 {
		return value.Nil, fmt.Errorf("CHARACTER.ADDVELOCITY needs handle, vx, vy, vz")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	vx, _ := args[1].ToFloat()
	vy, _ := args[2].ToFloat()
	vz, _ := args[3].ToFloat()
	v := co.cv.GetLinearVelocity()
	v.X += float32(vx)
	v.Y += float32(vy)
	v.Z += float32(vz)
	co.cv.SetLinearVelocity(v)
	return args[0], nil
}

func ccJump(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 2 {
		return value.Nil, fmt.Errorf("CHARACTER.JUMP needs handle, impulse")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	j, _ := args[1].ToFloat()
	v := co.cv.GetLinearVelocity()
	v.Y = float32(j)
	co.cv.SetLinearVelocity(v)
	return args[0], nil
}

func ccUpdate(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 1 {
		return value.Nil, fmt.Errorf("CHARACTER.UPDATE needs handle")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	dt := m.kccDt()
	m.extendedUpdateChar(co, dt)
	return args[0], nil
}

func ccOnSlope(m *Module, args []value.Value) (value.Value, error) {
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.FromBool(false), nil
	}
	if co.cv.GetGroundState() != jolt.GroundStateOnGround {
		return value.FromBool(false), nil
	}
	n := co.cv.GetGroundNormal()
	angle := math.Acos(float64(n.Y)) * (180.0 / math.Pi)
	return value.FromBool(angle > 1.0), nil
}

func ccOnWall(m *Module, args []value.Value) (value.Value, error) {
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.FromBool(false), nil
	}
	st := co.cv.GetGroundState()
	return value.FromBool(st == jolt.GroundStateOnSteepGround), nil
}

func ccSlopeAngle(m *Module, args []value.Value) (value.Value, error) {
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.FromFloat(0), nil
	}
	if co.cv.GetGroundState() != jolt.GroundStateOnGround {
		return value.FromFloat(0), nil
	}
	n := co.cv.GetGroundNormal()
	angle := math.Acos(float64(n.Y)) * (180.0 / math.Pi)
	return value.FromFloat(angle), nil
}

func ccGetSpeed(m *Module, args []value.Value) (value.Value, error) {
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.FromFloat(0), nil
	}
	v := co.cv.GetLinearVelocity()
	speed := math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z))
	return value.FromFloat(speed), nil
}

func ccSetGravity(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 2 {
		return value.Nil, nil
	}
	co, _ := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	g, _ := args[1].ToFloat()
	co.gravityScale = float32(g)
	return args[0], nil
}

func ccSetFriction(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 2 {
		return value.Nil, nil
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	f, _ := args[1].ToFloat()
	if f < 0 {
		f = 0
	}
	if f > 2 {
		f = 2
	}
	co.charFriction = float32(f)
	return args[0], nil
}

func ccSetMaxSlope(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 2 {
		return value.Nil, nil
	}
	co, _ := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	angle, _ := args[1].ToFloat()
	co.maxSlopeDeg = float32(angle)
	return args[0], nil
}

func ccSetStepHeight(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 2 {
		return value.Nil, nil
	}
	co, _ := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	h, _ := args[1].ToFloat()
	co.extUpdate.WalkStairsStepUp = jolt.Vec3{Y: float32(h)}
	return args[0], nil
}

func ccSetSnapDist(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 2 {
		return value.Nil, nil
	}
	co, _ := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	d, _ := args[1].ToFloat()
	co.extUpdate.StickToFloorStepDown = jolt.Vec3{Y: -float32(d)}
	return args[0], nil
}

func ccGetGravity(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 1 {
		return value.FromFloat(0), nil
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.FromFloat(1), nil
	}
	return value.FromFloat(float64(co.gravityScale)), nil
}

func ccGetMaxSlope(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 1 {
		return value.FromFloat(0), nil
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.FromFloat(45), nil
	}
	return value.FromFloat(float64(co.maxSlopeDeg)), nil
}

func ccGetStepHeight(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 1 {
		return value.FromFloat(0), nil
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(float64(co.extUpdate.WalkStairsStepUp.Y)), nil
}

func ccGetSnapDist(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 1 {
		return value.FromFloat(0), nil
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(float64(-co.extUpdate.StickToFloorStepDown.Y)), nil
}

func ccGetFriction(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 1 {
		return value.FromFloat(0.5), nil
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.FromFloat(0.5), nil
	}
	return value.FromFloat(float64(co.charFriction)), nil
}

func ccGetJumpBuffer(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 1 {
		return value.FromFloat(0.1), nil
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.FromFloat(0.1), nil
	}
	return value.FromFloat(float64(co.jumpBufferSec)), nil
}

func ccGetAirControl(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 1 {
		return value.FromFloat(1), nil
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.FromFloat(1), nil
	}
	return value.FromFloat(float64(co.airControl)), nil
}

func ccGetGroundControl(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 1 {
		return value.FromFloat(1), nil
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.FromFloat(1), nil
	}
	return value.FromFloat(float64(co.groundControl)), nil
}
