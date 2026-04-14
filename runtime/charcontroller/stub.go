//go:build (!linux && !windows) || !cgo

package mbcharcontroller

import (
	"fmt"
	"math"
	"moonbasic/runtime"
	mbphysics3d "moonbasic/runtime/physics3d"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type Vec3 struct {
	X, Y, Z float32
}

type charObj struct {
	pos      Vec3
	vel      Vec3
	grounded bool
	radius   float32
	height   float32
	stepH    float32
	snapD    float32
	gravityG float32
	maxSlope float32
	friction float32
	bounce   float32
	release  heap.ReleaseOnce
}

func (c *charObj) TypeName() string { return "CharController" }
func (c *charObj) TypeTag() uint16  { return heap.TagCharController }
func (c *charObj) Free()            {}

const stubCharDt = float32(1.0 / 60.0)

func registerCharControllerCommands(m *Module, reg runtime.Registrar) {
	// Host KCC owns CHARACTER.* on this build; CHARCONTROLLER.* uses the lightweight AABB capsule stub.
	reg.Register("CHARCONTROLLER.MAKE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcMake(m, args) }))
	reg.Register("CHARCONTROLLER.SETPOS", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcSetPos(m, args) }))
	reg.Register("CHARCONTROLLER.SETPOSITION", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcSetPos(m, args) }))
	reg.Register("CHARCONTROLLER.GETPOS", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcGetPos(m, args) }))
	reg.Register("CHARCONTROLLER.MOVE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcMove(m, args) }))
	reg.Register("CHARCONTROLLER.ISGROUNDED", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcIsGrounded(m, args) }))
	reg.Register("CHARCONTROLLER.X", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcAxis(m, args, 0) }))
	reg.Register("CHARCONTROLLER.Y", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcAxis(m, args, 1) }))
	reg.Register("CHARCONTROLLER.Z", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcAxis(m, args, 2) }))
	reg.Register("CHARCONTROLLER.FREE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcFree(m, args) }))
	reg.Register("CHARCONTROLLER.GETLINEARVEL", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcGetLinearVel(m, args) }))
	reg.Register("CHARCONTROLLER.GETGROUNDVELOCITY", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcGetGroundVel(m, args) }))
	reg.Register("CHARCONTROLLER.GETGROUNDNORMAL", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcGetGroundNormalArr(m, args) }))
	reg.Register("CHARCONTROLLER.GROUNDSTATE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcGroundState(m, args) }))
	reg.Register("CHARCONTROLLER.TELEPORT", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcTeleport(m, args) }))

	reg.Register("CHARACTERREF.SETPOSITION", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcSetPos(m, args) }))
	reg.Register("CHARACTERREF.MOVE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcMove(m, args) }))
	reg.Register("CHARACTERREF.SETVELOCITY", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcSetVel(m, args) }))
	reg.Register("CHARACTERREF.ADDVELOCITY", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcAddVel(m, args) }))
	reg.Register("CHARACTERREF.JUMP", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcJump(m, args) }))
	reg.Register("CHARACTERREF.UPDATE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcUpdate(m, args) }))
	reg.Register("CHARACTERREF.ISGROUNDED", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcIsGrounded(m, args) }))
	reg.Register("CHARACTERREF.ONSLOPE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcOnSlope(m, args) }))
	reg.Register("CHARACTERREF.ONWALL", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcOnWall(m, args) }))
	reg.Register("CHARACTERREF.GETSLOPEANGLE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcSlopeAngle(m, args) }))
	reg.Register("CHARACTERREF.GETSPEED", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcGetSpeed(m, args) }))
	reg.Register("CHARACTERREF.SETGRAVITY", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcSetGravity(m, args) }))
	reg.Register("CHARACTERREF.SETFRICTION", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcSetFriction(m, args) }))
	reg.Register("CHARACTERREF.SETMAXSLOPE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcSetMaxSlope(m, args) }))
	reg.Register("CHARACTERREF.SETSTEPHEIGHT", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcSetStepHeight(m, args) }))
	reg.Register("CHARACTERREF.SETSNAPDISTANCE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcSetSnapDist(m, args) }))
	reg.Register("CHARACTERREF.FREE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcFree(m, args) }))
	reg.Register("CHARACTERREF.GETPOSITION", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) { return stubCcGetPos(m, args) }))
}

func stubAllocFloat3(m *Module, x, y, z float64) (value.Value, error) {
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

func stubCcMake(m *Module, args []value.Value) (value.Value, error) {
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
	if radius <= 0 || height <= 0 {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MAKE: radius and height must be positive")
	}
	co := &charObj{
		pos:      Vec3{X: float32(x), Y: float32(y), Z: float32(z)},
		vel:      Vec3{},
		radius:   float32(radius),
		height:   float32(height),
		stepH:    0.3,
		snapD:    0.2,
		gravityG: 1,
		maxSlope: 45,
		friction: 0.5,
	}
	h, err := m.h.Alloc(co)
	if err != nil {
		return value.Nil, err
	}
	if co2, err2 := heap.Cast[*charObj](m.h, h); err2 == nil {
		co2.hostUpdate(stubCharDt)
	}
	return value.FromHandle(h), nil
}

func stubCcSetPos(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.SETPOS: heap not bound")
	}
	if len(args) < 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.SETPOS: need handle, x, y, z")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
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
	co.pos = Vec3{X: float32(x), Y: float32(y), Z: float32(z)}
	co.hostUpdate(stubCharDt)
	return value.Nil, nil
}

func stubCcMove(m *Module, args []value.Value) (value.Value, error) {
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
	co.pos.X += float32(dx)
	co.pos.Y += float32(dy)
	co.pos.Z += float32(dz)
	co.hostUpdate(stubCharDt)
	return value.Nil, nil
}

func stubCcAxis(m *Module, args []value.Value, axis int) (value.Value, error) {
	if m.h == nil {
		return value.FromFloat(0), nil
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.FromFloat(0), nil
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.FromFloat(0), nil
	}
	switch axis {
	case 0:
		return value.FromFloat(float64(co.pos.X)), nil
	case 1:
		return value.FromFloat(float64(co.pos.Y)), nil
	default:
		return value.FromFloat(float64(co.pos.Z)), nil
	}
}

func stubCcIsGrounded(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.ISGROUNDED: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.ISGROUNDED: need handle")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.ISGROUNDED: invalid handle")
	}
	return value.FromBool(co.grounded), nil
}

func stubCcGetPos(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETPOS: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETPOS: need handle")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETPOS: invalid handle")
	}
	return stubAllocFloat3(m, float64(co.pos.X), float64(co.pos.Y), float64(co.pos.Z))
}

func stubCcFree(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.FREE: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.FREE: need handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func stubCcGetLinearVel(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETLINEARVEL: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETLINEARVEL: need handle")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETLINEARVEL: invalid handle")
	}
	return stubAllocFloat3(m, float64(co.vel.X), float64(co.vel.Y), float64(co.vel.Z))
}

func stubCcGetGroundVel(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETGROUNDVELOCITY: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETGROUNDVELOCITY: need handle")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETGROUNDVELOCITY: invalid handle")
	}
	if co.grounded {
		return stubAllocFloat3(m, float64(co.vel.X), 0, float64(co.vel.Z))
	}
	return stubAllocFloat3(m, float64(co.vel.X), float64(co.vel.Y), float64(co.vel.Z))
}

func stubCcGetGroundNormalArr(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETGROUNDNORMAL: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETGROUNDNORMAL: need handle")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETGROUNDNORMAL: invalid handle")
	}
	if co.grounded {
		return stubAllocFloat3(m, 0, 1, 0)
	}
	return stubAllocFloat3(m, 0, 0, 0)
}

func stubCcGroundState(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GROUNDSTATE: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GROUNDSTATE: need handle")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GROUNDSTATE: invalid handle")
	}
	if co.grounded {
		return value.FromInt(0), nil
	}
	return value.FromInt(3), nil
}

func stubCcTeleport(m *Module, args []value.Value) (value.Value, error) {
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
	return value.Nil, nil
}

func stubCcSetVel(m *Module, args []value.Value) (value.Value, error) {
	co, err := hGet(m, args[0])
	if err != nil { return value.Nil, err }
	vx, _ := args[1].ToFloat()
	vy, _ := args[2].ToFloat()
	vz, _ := args[3].ToFloat()
	co.vel = Vec3{X: float32(vx), Y: float32(vy), Z: float32(vz)}
	return value.Nil, nil
}

func stubCcAddVel(m *Module, args []value.Value) (value.Value, error) {
	co, err := hGet(m, args[0])
	if err != nil { return value.Nil, err }
	vx, _ := args[1].ToFloat()
	vy, _ := args[2].ToFloat()
	vz, _ := args[3].ToFloat()
	co.vel.X += float32(vx)
	co.vel.Y += float32(vy)
	co.vel.Z += float32(vz)
	return value.Nil, nil
}

func stubCcJump(m *Module, args []value.Value) (value.Value, error) {
	co, err := hGet(m, args[0])
	if err != nil { return value.Nil, err }
	f, _ := args[1].ToFloat()
	co.vel.Y = float32(f)
	co.grounded = false
	return value.Nil, nil
}

func stubCcUpdate(m *Module, args []value.Value) (value.Value, error) {
	co, err := hGet(m, args[0])
	if err != nil { return value.Nil, err }
	dt, _ := args[1].ToFloat()
	co.hostUpdate(float32(dt))
	return value.Nil, nil
}

func stubCcOnSlope(m *Module, args []value.Value) (value.Value, error) {
	return value.FromBool(false), nil
}

func stubCcOnWall(m *Module, args []value.Value) (value.Value, error) {
	return value.FromBool(false), nil
}

func stubCcSlopeAngle(m *Module, args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func stubCcGetSpeed(m *Module, args []value.Value) (value.Value, error) {
	co, err := hGet(m, args[0])
	if err != nil { return value.FromFloat(0), nil }
	s := math.Sqrt(float64(co.vel.X*co.vel.X + co.vel.Y*co.vel.Y + co.vel.Z*co.vel.Z))
	return value.FromFloat(s), nil
}

func stubCcSetGravity(m *Module, args []value.Value) (value.Value, error) {
	co, err := hGet(m, args[0])
	if err != nil { return value.Nil, err }
	g, _ := args[1].ToFloat()
	co.gravityG = float32(g)
	return value.Nil, nil
}

func stubCcSetFriction(m *Module, args []value.Value) (value.Value, error) {
	co, err := hGet(m, args[0])
	if err != nil { return value.Nil, err }
	f, _ := args[1].ToFloat()
	co.friction = float32(f)
	return value.Nil, nil
}

func stubCcSetMaxSlope(m *Module, args []value.Value) (value.Value, error) {
	co, err := hGet(m, args[0])
	if err != nil { return value.Nil, err }
	deg, _ := args[1].ToFloat()
	co.maxSlope = float32(deg)
	return value.Nil, nil
}

func stubCcSetStepHeight(m *Module, args []value.Value) (value.Value, error) {
	co, err := hGet(m, args[0])
	if err != nil { return value.Nil, err }
	h, _ := args[1].ToFloat()
	co.stepH = float32(h)
	return value.Nil, nil
}

func stubCcSetSnapDist(m *Module, args []value.Value) (value.Value, error) {
	co, err := hGet(m, args[0])
	if err != nil { return value.Nil, err }
	d, _ := args[1].ToFloat()
	co.snapD = float32(d)
	return value.Nil, nil
}

// CharacterGroundNormal reports a world-up normal when grounded (stub KCC); otherwise zero.
func (m *Module) CharacterGroundNormal(h heap.Handle) (nx, ny, nz float64, ok bool) {
	if m.h == nil {
		return 0, 0, 0, false
	}
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil {
		return 0, 0, 0, false
	}
	if !co.grounded {
		return 0, 0, 0, false
	}
	return 0, 1, 0, true
}

// CharacterGroundStateInt matches Jolt EGroundState when applicable: 0 OnGround, 3 InAir.
func (m *Module) CharacterGroundStateInt(h heap.Handle) (int, bool) {
	if m.h == nil {
		return 0, false
	}
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil {
		return 0, false
	}
	if co.grounded {
		return 0, true
	}
	return 3, true
}

func (m *Module) SetCharacterRestitution(h heap.Handle, bounce float64) error {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil {
		return fmt.Errorf("invalid character handle")
	}
	co.bounce = float32(bounce)
	return nil
}

func (m *Module) SetCharacterFriction(h heap.Handle, friction float64) error {
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil {
		return fmt.Errorf("invalid character handle")
	}
	co.friction = float32(friction)
	return nil
}

// CharacterTeleport snaps position and clears velocity (stub).
func (m *Module) CharacterTeleport(h heap.Handle, x, y, z float64) error {
	if m.h == nil {
		return fmt.Errorf("CharacterTeleport: heap not bound")
	}
	co, err := heap.Cast[*charObj](m.h, h)
	if err != nil {
		return fmt.Errorf("invalid character handle")
	}
	co.pos = Vec3{X: float32(x), Y: float32(y), Z: float32(z)}
	co.vel = Vec3{}
	co.hostUpdate(stubCharDt)
	return nil
}

func hGet(m *Module, v value.Value) (*charObj, error) {
	return heap.Cast[*charObj](m.h, heap.Handle(v.IVal))
}

func (c *charObj) hostUpdate(dt float32) {
	if c.grounded {
		c.vel.Y = -0.1
	} else {
		c.vel.Y -= 32.0 * c.gravityG * dt
	}

	// 1. Horizontal Phase (Iterative Slide)
	moveX := c.vel.X * dt
	moveZ := c.vel.Z * dt
	
	c.pos.X += moveX
	c.pos.Z += moveZ

	statics := mbphysics3d.GetStaticBodyRegistry()
	
	// Iterative Slide (Simplified for AABB-stubs)
	for i := 0; i < 3; i++ {
		hit := false
		var nx, nz float32
		
		for _, b := range statics {
			if b.Shape == nil || b.Shape.Kind != 1 { continue }
			hx, hy, hz := b.Shape.F1, b.Shape.F2, b.Shape.F3
			
			// Y-range check for body
			if c.pos.Y+c.height*0.5 < b.Pos.Y-hy || c.pos.Y-c.height*0.5 > b.Pos.Y+hy {
				continue
			}

			// AABB overlap check
			if c.pos.X > b.Pos.X-hx-c.radius && c.pos.X < b.Pos.X+hx+c.radius &&
			   c.pos.Z > b.Pos.Z-hz-c.radius && c.pos.Z < b.Pos.Z+hz+c.radius {
				
				hit = true
				// Calculate push-out normal
				dx := c.pos.X - b.Pos.X
				dz := c.pos.Z - b.Pos.Z
				
				if math.Abs(float64(dx))/(float64(hx)+float64(c.radius)) > math.Abs(float64(dz))/(float64(hz)+float64(c.radius)) {
					if dx > 0 { nx = 1; c.pos.X = b.Pos.X + hx + c.radius + 0.001 } else { nx = -1; c.pos.X = b.Pos.X - hx - c.radius - 0.001 }
				} else {
					if dz > 0 { nz = 1; c.pos.Z = b.Pos.Z + hz + c.radius + 0.001 } else { nz = -1; c.pos.Z = b.Pos.Z - hz - c.radius - 0.001 }
				}
				break
			}
		}
		if !hit { break }
		
		// Project velocity onto plane
		dot := c.vel.X*nx + c.vel.Z*nz
		if dot < 0 {
			c.vel.X -= nx * dot
			c.vel.Z -= nz * dot
		}
	}

	// 2. Vertical Phase
	c.pos.Y += c.vel.Y * dt

	// 3. Ground Snapping
	c.grounded = false
	feetY := c.pos.Y - c.height*0.5
	
	for _, b := range statics {
		if b.Shape == nil || b.Shape.Kind != 1 { continue }
		hx, hy, hz := b.Shape.F1, b.Shape.F2, b.Shape.F3
		if c.pos.X > b.Pos.X-hx-c.radius && c.pos.X < b.Pos.X+hx+c.radius &&
		   c.pos.Z > b.Pos.Z-hz-c.radius && c.pos.Z < b.Pos.Z+hz+c.radius {
			
			topY := b.Pos.Y + hy
			if feetY <= topY + c.snapD && feetY >= topY - 0.5 {
				c.pos.Y = topY + c.height*0.5
				c.vel.Y = 0
				c.grounded = true
				break
			}
		}
	}
}

func shutdownCharController(m *Module) { _ = m }
