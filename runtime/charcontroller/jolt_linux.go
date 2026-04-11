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
	reg.Register("CHARCONTROLLER.MAKE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return ccMake(m, args)
	}))
	reg.Register("CHARCONTROLLER.SETPOSITION", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return ccSetPos(m, args)
	}))
	reg.Register("CHARCONTROLLER.SETPOS", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return ccSetPos(m, args)
	}))
	reg.Register("CHARCONTROLLER.GETPOS", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return ccGetPos(m, args)
	}))
	reg.Register("CHARCONTROLLER.MOVE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return ccMove(m, args)
	}))
	reg.Register("CHARCONTROLLER.ISGROUNDED", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return ccIsGrounded(m, args)
	}))
	reg.Register("CHARCONTROLLER.X", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return ccAxis(m, args, 0)
	}))
	reg.Register("CHARCONTROLLER.Y", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return ccAxis(m, args, 1)
	}))
	reg.Register("CHARCONTROLLER.Z", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return ccAxis(m, args, 2)
	}))
	reg.Register("CHARCONTROLLER.FREE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return ccFree(m, args)
	}))

	reg.Register("CONTROLLER.CREATE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccMake(m, args) }))
	reg.Register("CONTROLLER.MOVE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccMove(m, args) }))
	reg.Register("CONTROLLER.GROUNDED", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccIsGrounded(m, args) }))
	reg.Register("CONTROLLER.JUMP", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("CONTROLLER.JUMP: not implemented; apply upward velocity via CHARCONTROLLER.MOVE or physics")
	}))
	reg.Register("CONTROLLER.FREE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return ccFree(m, args) }))
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
	if padding > 0 {
		settings.CharacterPadding = float32(padding)
	}
	cv := ps.CreateCharacterVirtual(settings, jolt.Vec3{X: float32(x), Y: float32(y), Z: float32(z)})
	ext := jolt.DefaultExtendedUpdateSettings()
	pad := settings.CharacterPadding
	h, err := m.h.Alloc(&charObj{
		cv: cv, extUpdate: ext, gravityScale: 1, charMass: 70,
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
	return h, nil
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
