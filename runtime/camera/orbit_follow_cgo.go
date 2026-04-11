//go:build cgo || (windows && !cgo)

package mbcamera

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	mbtime "moonbasic/runtime/time"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// camOrbitDispatch implements CAMERA.ORBIT: (cam, entity, distance#) for automated orbit, or
// (cam, tx, ty, tz, yaw#, pitch#, distance#) for explicit spherical placement (Blitz / SETORBIT).
func (m *Module) camOrbitDispatch(args []value.Value) (value.Value, error) {
	switch len(args) {
	case 3:
		return m.camOrbitFollowEntity(args)
	case 7:
		return m.camSetOrbit(args)
	default:
		return value.Nil, fmt.Errorf("CAMERA.ORBIT expects (camera, entity, distance#) or (camera, tx, ty, tz, yaw#, pitch#, distance#)")
	}
}

// camOrbitFollowEntity: third-person orbit with internal yaw/pitch/dist, mouse (R-drag), Q/E, wheel.
// Matches GAMEHELPERS orbit defaults (pitch clamp, sensitivity) so BASIC stays math-free.
func (m *Module) camOrbitFollowEntity(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.ORBIT: heap not bound")
	}
	if m.entityWorldPos == nil {
		return value.Nil, fmt.Errorf("CAMERA.ORBIT: entity bridge not ready (ensure entity module binds camera)")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CAMERA.ORBIT (entity) expects (camera, entity, distance#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.ORBIT: invalid camera handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	if args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CAMERA.ORBIT: entity must be a handle")
	}
	eh := heap.Handle(args[1].IVal)
	baseDist, ok := argF(args[2])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.ORBIT: distance must be numeric")
	}
	tx, ty, tz, ok := m.entityWorldPos(m.h, eh)
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.ORBIT: could not resolve entity world position")
	}

	const defaultPitch = 0.2

	dt := mbtime.DeltaSeconds(nil)
	if dt <= 0 {
		dt = 1.0 / 60.0
	}

	if !o.orbitInited {
		o.orbitDist = baseDist
		o.orbitPitch = defaultPitch
		o.orbitYaw = 0
		o.orbitInited = true
	}
	if o.orbitDist <= 0 {
		o.orbitDist = baseDist
	}

	minPitch := o.orbitMinPitch
	maxPitch := o.orbitMaxPitch
	distMin := o.orbitMinDist
	distMax := o.orbitMaxDist
	targetYOff := o.orbitTargetYOff
	mouseSens := o.orbitMouseSens
	wheelSens := o.orbitWheelSens

	// Mouse yaw / pitch
	md := rl.GetMouseDelta()
	mouseActive := o.orbitUseMouse && (!o.orbitRightMouseForDrag || rl.IsMouseButtonDown(rl.MouseRightButton))
	if mouseActive {
		o.orbitYaw -= float32(md.X) * mouseSens
		o.orbitPitch += float32(md.Y) * mouseSens
	}

	// Keyboard orbit (Q/E or custom keys): exclusive, same polarity as before.
	var ax float64
	kL, kR := o.orbitKeyLeft, o.orbitKeyRight
	if kL != 0 || kR != 0 {
		downL := kL != 0 && rl.IsKeyDown(kL)
		downR := kR != 0 && rl.IsKeyDown(kR)
		if downL && !downR {
			ax = -1
		} else if downR && !downL {
			ax = 1
		}
	}
	o.orbitYaw += float32(ax * float64(o.orbitKeyRadPerSec) * dt)

	o.orbitPitch = clamp32(o.orbitPitch, minPitch, maxPitch)

	o.orbitDist -= float32(rl.GetMouseWheelMove()) * wheelSens
	o.orbitDist = clamp32(o.orbitDist, distMin, distMax)

	tyLook := ty + targetYOff

	// Spherical coordinates (Y-up): standard yaw/pitch orbit around target.
	cp := math.Cos(float64(o.orbitPitch))
	sp := math.Sin(float64(o.orbitPitch))
	sy := math.Sin(float64(o.orbitYaw))
	cy := math.Cos(float64(o.orbitYaw))
	d := float64(o.orbitDist)

	o.cam.Position.X = tx + float32(d*cp*sy)
	o.cam.Position.Y = tyLook + float32(d*sp)
	o.cam.Position.Z = tz + float32(d*cp*cy)
	o.cam.Target = rl.Vector3{X: tx, Y: tyLook, Z: tz}
	return value.Nil, nil
}

func clamp32(x, lo, hi float32) float32 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

// camUseMouseOrbit maps to cam.UseMouseOrbit / CAMERA.USEMOUSEORBIT: enable/disable mouse contribution to orbit-follow.
func (m *Module) camUseMouseOrbit(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.USEMOUSEORBIT: heap not bound")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CAMERA.USEMOUSEORBIT expects (camera handle, useMouse#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.USEMOUSEORBIT: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	o.orbitUseMouse = argBool(args[1])
	return value.Nil, nil
}

// camUseOrbitRightMouse maps to CAMERA.USEORBITRIGHTMOUSE: when true (default), mouse orbit only while RMB is held; when false, mouse moves orbit without RMB.
func (m *Module) camUseOrbitRightMouse(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.USEORBITRIGHTMOUSE: heap not bound")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CAMERA.USEORBITRIGHTMOUSE expects (camera handle, requireRightMouseDrag#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.USEORBITRIGHTMOUSE: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	o.orbitRightMouseForDrag = argBool(args[1])
	return value.Nil, nil
}

// camSetOrbitKeys maps to CAMERA.SETORBITKEYS: raylib key codes; use 0 for a side to disable keyboard orbit on that side.
func (m *Module) camSetOrbitKeys(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.SETORBITKEYS: heap not bound")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CAMERA.SETORBITKEYS expects (camera handle, leftKey#, rightKey#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.SETORBITKEYS: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	kL, ok1 := argKeyCode(args[1])
	kR, ok2 := argKeyCode(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("CAMERA.SETORBITKEYS: key codes must be numeric")
	}
	o.orbitKeyLeft = kL
	o.orbitKeyRight = kR
	return value.Nil, nil
}

// camSetOrbitLimits maps to CAMERA.SETORBITLIMITS: pitch (radians), distance clamps.
func (m *Module) camSetOrbitLimits(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.SETORBITLIMITS: heap not bound")
	}
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("CAMERA.SETORBITLIMITS expects (camera handle, minPitch#, maxPitch#, minDist#, maxDist#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.SETORBITLIMITS: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	minP, ok1 := argF(args[1])
	maxP, ok2 := argF(args[2])
	minD, ok3 := argF(args[3])
	maxD, ok4 := argF(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("CAMERA.SETORBITLIMITS: limits must be numeric")
	}
	o.orbitMinPitch = minP
	o.orbitMaxPitch = maxP
	o.orbitMinDist = minD
	o.orbitMaxDist = maxD
	return value.Nil, nil
}

// camSetOrbitSpeed maps to CAMERA.SETORBITSPEED: mouse drag sensitivity and scroll-wheel zoom scale.
func (m *Module) camSetOrbitSpeed(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.SETORBITSPEED: heap not bound")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CAMERA.SETORBITSPEED expects (camera handle, mouseSens#, wheelSens#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.SETORBITSPEED: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	ms, ok1 := argF(args[1])
	ws, ok2 := argF(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("CAMERA.SETORBITSPEED: sensitivity must be numeric")
	}
	o.orbitMouseSens = ms
	o.orbitWheelSens = ws
	return value.Nil, nil
}

// camSetOrbitKeySpeed maps to CAMERA.SETORBITKEYSPEED: keyboard yaw rate in radians per second.
func (m *Module) camSetOrbitKeySpeed(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.SETORBITKEYSPEED: heap not bound")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CAMERA.SETORBITKEYSPEED expects (camera handle, keyRadPerSec#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.SETORBITKEYSPEED: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	kr, ok1 := argF(args[1])
	if !ok1 {
		return value.Nil, fmt.Errorf("CAMERA.SETORBITKEYSPEED: keyRadPerSec must be numeric")
	}
	o.orbitKeyRadPerSec = kr
	return value.Nil, nil
}

// camGetYaw returns internal orbit yaw (radians) for aligning entities with the orbit camera.
func (m *Module) camGetYaw(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.YAW: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CAMERA.YAW expects (camera handle)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.YAW: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(o.orbitYaw)), nil
}
