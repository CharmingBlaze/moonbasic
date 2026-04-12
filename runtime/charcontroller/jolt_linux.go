//go:build linux && cgo

package mbcharcontroller

import (
	"fmt"
	"math"
	"sync"

	"github.com/bbitechnologies/jolt-go/jolt"

	"moonbasic/runtime"
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

const charDt = float32(1.0 / 60.0)

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
	crouch       bool    // gameplay flag (capsule resize not in wrapper yet)
	charMass     float32 // stored for gameplay; Jolt CharacterVirtual mass is fixed at create

	// Capsule + padding (for rebuild after PLAYER.SETSLOPE / CHAR.SETPADDING).
	capRadius   float32
	capFullH    float32 // total capsule height passed to CreateCapsule
	maxSlopeDeg float32
	charPad     float32 // CharacterVirtualSettings.CharacterPadding
	charFriction float32 // gameplay scalar (CHARACTERREF.SETFRICTION); not Jolt body friction
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
	reg.Register("CHARACTER.CREATE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) == 3 && args[0].Kind != value.KindHandle {
			return ccMake(m, args)
		}
		// Overload: CHARACTER.CREATE(entity, r, h)
		return ccMakeLegacy(m, args)
	}))
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
	reg.Register("CHARACTERREF.FREE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccFree(m, args) }))
	reg.Register("CHARACTERREF.GETPOSITION", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccGetPos(m, args) }))
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
	h, err := createCharacterVirtualFromParams(m, radius, height, x, y, z, 45, -1)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

// createCharacterVirtualFromParams allocates a CharacterVirtual. maxSlopeDeg ≤ 0 uses 45°.
// padding ≤ 0 uses Jolt default CharacterPadding (0.02); else sets CharacterVirtualSettings.CharacterPadding.
func createCharacterVirtualFromParams(m *Module, radius, height, x, y, z, maxSlopeDeg, padding float64) (heap.Handle, error) {
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
	settings.MaxSlopeAngle = jolt.DegreesToRadians(float32(maxSlopeDeg))
	// Bug Fix: Explicitly set recovery speed and contact distance to ensure zero-bounce
	settings.PenetrationRecoverySpeed = 1.0
	settings.PredictiveContactDistance = 0.1
	if padding > 0 {
		settings.CharacterPadding = float32(padding)
	}
	cv := ps.CreateCharacterVirtual(settings, jolt.Vec3{X: float32(x), Y: float32(y), Z: float32(z)})
	ext := jolt.DefaultExtendedUpdateSettings()
	pad := settings.CharacterPadding
	h, err := m.h.Alloc(&charObj{
		cv: cv, extUpdate: ext, gravityScale: 1, charMass: 70, charFriction: 0.5,
		capRadius: fr, capFullH: fh, maxSlopeDeg: float32(maxSlopeDeg), charPad: pad,
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
	m.extendedUpdateChar(co, charDt)
	return value.Nil, nil
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
	m.extendedUpdateChar(co, charDt)
	return value.Nil, nil
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

// AllocCharacter creates a Jolt CharacterVirtual (used by runtime/player).
// maxSlopeDeg ≤ 0 defaults to 45° (matches prior CHARCONTROLLER.MAKE behavior).
// padding ≤ 0 uses default skin width; >0 sets CharacterPadding.
func (m *Module) AllocCharacter(radius, height, x, y, z float64, maxSlopeDeg, padding float64) (heap.Handle, error) {
	return createCharacterVirtualFromParams(m, radius, height, x, y, z, maxSlopeDeg, padding)
}

// RecreateCharacterWithSlope destroys a character and creates a new one at the same pose with a new max walk slope.
func (m *Module) RecreateCharacterWithSlope(oldH heap.Handle, radius, height, maxSlopeDeg float64) (heap.Handle, error) {
	if m.h == nil {
		return 0, fmt.Errorf("RecreateCharacterWithSlope: heap not bound")
	}
	co, err := heap.Cast[*charObj](m.h, oldH)
	if err != nil || co.cv == nil {
		return 0, fmt.Errorf("RecreateCharacterWithSlope: invalid handle")
	}
	p := co.cv.GetPosition()
	v := co.cv.GetLinearVelocity()
	extCopy := co.extUpdate
	gs := co.gravityScale
	sm := co.swimMode
	sb, sd := co.swimBuoyancy, co.swimDrag
	cr := co.crouch
	cm := co.charMass
	cf := co.charFriction
	untrackChar(oldH)
	co.cv.Destroy()
	co.cv = nil
	if err := m.h.Free(oldH); err != nil {
		return 0, err
	}
	h, err := createCharacterVirtualFromParams(m, radius, height, float64(p.X), float64(p.Y), float64(p.Z), maxSlopeDeg, float64(co.charPad))
	if err != nil {
		return 0, err
	}
	co2, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co2.cv == nil {
		return 0, fmt.Errorf("RecreateCharacterWithSlope: internal")
	}
	co2.cv.SetLinearVelocity(v)
	co2.extUpdate = extCopy
	co2.gravityScale = gs
	co2.swimMode = sm
	co2.swimBuoyancy, co2.swimDrag = sb, sd
	co2.crouch = cr
	co2.charMass = cm
	co2.charFriction = cf
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
	if dt <= 0 {
		dt = 1.0 / 60.0
	}
	vel := co.cv.GetLinearVelocity()
	g := mbphysics3d.GravityVec()
	gs := co.gravityScale
	if gs <= 0 {
		gs = 1
	}
	gy := g.Y * float32(dt) * gs
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
		if d > 0 && dt > 0 {
			damp := math.Max(0, 1-d*dt)
			vel.X *= float32(damp)
			vel.Z *= float32(damp)
		}
	}
	vel.Y += gy
	if co.cv.GetGroundState() == jolt.GroundStateOnGround && vel.Y < -groundedMaxDownVy {
		vel.Y = -groundedMaxDownVy
	}
	co.cv.SetLinearVelocity(vel)
	m.extendedUpdateChar(co, float32(dt))
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
	if pad < 1e-4 {
		pad = 0.02
	}
	co.charPad = pad
	return m.RecreateCharacterWithSlope(h, float64(co.capRadius), float64(co.capFullH), float64(co.maxSlopeDeg))
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

// SetCharacterCrouch stores a crouch flag (capsule height change is not in the Jolt wrapper yet).
func (m *Module) SetCharacterCrouch(h heap.Handle, on bool) error {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return fmt.Errorf("invalid character handle")
	}
	co.crouch = on
	return nil
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

// CharacterTeleport snaps the capsule to a world position and clears velocity (no interpolation smoothing).
func (m *Module) CharacterTeleport(h heap.Handle, x, y, z float64) error {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil || co.cv == nil {
		return fmt.Errorf("invalid character handle")
	}
	co.cv.SetPosition(jolt.Vec3{X: float32(x), Y: float32(y), Z: float32(z)})
	co.cv.SetLinearVelocity(jolt.Vec3{})
	m.extendedUpdateChar(co, charDt)
	return nil
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
	vel := co.cv.GetLinearVelocity()
	vel.X = float32(vx)
	vel.Z = float32(vz)
	g := mbphysics3d.GravityVec()
	gs := co.gravityScale
	if gs <= 0 {
		gs = 1
	}
	gy := g.Y * float32(dt) * gs
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
		if d > 0 && dt > 0 {
			damp := math.Max(0, 1-d*dt)
			vel.X *= float32(damp)
			vel.Z *= float32(damp)
		}
	}
	vel.Y += gy
	if co.cv.GetGroundState() == jolt.GroundStateOnGround && vel.Y < -groundedMaxDownVy {
		vel.Y = -groundedMaxDownVy
	}
	co.cv.SetLinearVelocity(vel)
	m.extendedUpdateChar(co, float32(dt))
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
	v := co.cv.GetLinearVelocity()
	v.Y += float32(impulseY)
	co.cv.SetLinearVelocity(v)
	m.extendedUpdateChar(co, charDt)
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

func ccMakeLegacy(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 3 { return value.Nil, fmt.Errorf("CHARACTER.CREATE(entity, r, h) needs 3 args") }
	r, _ := args[1].ToFloat()
	h_val, _ := args[2].ToFloat()
	h, err := createCharacterVirtualFromParams(m, r, h_val, 0, 0, 0, 45, -1)
	if err != nil { return value.Nil, err }
	return value.FromHandle(h), nil
}

func ccSetVel(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 4 { return value.Nil, fmt.Errorf("CHARACTER.SETVELOCITY needs handle, vx, vy, vz") }
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil { return value.Nil, err }
	vx, _ := args[1].ToFloat()
	vy, _ := args[2].ToFloat()
	vz, _ := args[3].ToFloat()
	co.cv.SetLinearVelocity(jolt.Vec3{X: float32(vx), Y: float32(vy), Z: float32(vz)})
	return value.Nil, nil
}

func ccAddVel(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 4 { return value.Nil, fmt.Errorf("CHARACTER.ADDVELOCITY needs handle, vx, vy, vz") }
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil { return value.Nil, err }
	vx, _ := args[1].ToFloat()
	vy, _ := args[2].ToFloat()
	vz, _ := args[3].ToFloat()
	v := co.cv.GetLinearVelocity()
	v.X += float32(vx); v.Y += float32(vy); v.Z += float32(vz)
	co.cv.SetLinearVelocity(v)
	return value.Nil, nil
}

func ccJump(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 2 { return value.Nil, fmt.Errorf("CHARACTER.JUMP needs handle, impulse") }
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil { return value.Nil, err }
	j, _ := args[1].ToFloat()
	v := co.cv.GetLinearVelocity()
	v.Y = float32(j) 
	co.cv.SetLinearVelocity(v)
	return value.Nil, nil
}

func ccUpdate(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 1 { return value.Nil, fmt.Errorf("CHARACTER.UPDATE needs handle") }
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil { return value.Nil, err }
	dt := charDt 
	m.extendedUpdateChar(co, dt)
	return value.Nil, nil
}

func ccOnSlope(m *Module, args []value.Value) (value.Value, error) {
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil { return value.FromBool(false), nil }
	if co.cv.GetGroundState() != jolt.GroundStateOnGround { return value.FromBool(false), nil }
	n := co.cv.GetGroundNormal()
	angle := math.Acos(float64(n.Y)) * (180.0 / math.Pi)
	return value.FromBool(angle > 1.0), nil
}

func ccOnWall(m *Module, args []value.Value) (value.Value, error) {
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil { return value.FromBool(false), nil }
	st := co.cv.GetGroundState()
	return value.FromBool(st == jolt.GroundStateOnSteepGround), nil
}

func ccSlopeAngle(m *Module, args []value.Value) (value.Value, error) {
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil { return value.FromFloat(0), nil }
	if co.cv.GetGroundState() != jolt.GroundStateOnGround { return value.FromFloat(0), nil }
	n := co.cv.GetGroundNormal()
	angle := math.Acos(float64(n.Y)) * (180.0 / math.Pi)
	return value.FromFloat(angle), nil
}

func ccGetSpeed(m *Module, args []value.Value) (value.Value, error) {
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil { return value.FromFloat(0), nil }
	v := co.cv.GetLinearVelocity()
	speed := math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z))
	return value.FromFloat(speed), nil
}

func ccSetGravity(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 2 { return value.Nil, nil }
	co, _ := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	g, _ := args[1].ToFloat()
	co.gravityScale = float32(g)
	return value.Nil, nil
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
	return value.Nil, nil
}

func ccSetMaxSlope(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 2 { return value.Nil, nil }
	co, _ := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	angle, _ := args[1].ToFloat()
	co.maxSlopeDeg = float32(angle)
	return value.Nil, nil
}

func ccSetStepHeight(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 2 { return value.Nil, nil }
	co, _ := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	h, _ := args[1].ToFloat()
	co.extUpdate.WalkStairsStepUp = jolt.Vec3{Y: float32(h)}
	return value.Nil, nil
}

func ccSetSnapDist(m *Module, args []value.Value) (value.Value, error) {
	if len(args) < 2 { return value.Nil, nil }
	co, _ := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	d, _ := args[1].ToFloat()
	co.extUpdate.StickToFloorStepDown = jolt.Vec3{Y: -float32(d)}
	return value.Nil, nil
}
