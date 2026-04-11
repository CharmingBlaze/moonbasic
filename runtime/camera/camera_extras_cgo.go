//go:build cgo || (windows && !cgo)

package mbcamera

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// camTurnLeft returns radians to add to orbit yaw (negative amount).
func (m *Module) camTurnLeft(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CAMERA.TURNLEFT expects (cam, amount#)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CAMERA.TURNLEFT: first argument must be a camera handle")
	}
	if _, err := heap.Cast[*camObj](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	amt, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.TURNLEFT: amount must be numeric")
	}
	return value.FromFloat(-math.Abs(amt)), nil
}

// camTurnRight returns radians to add to orbit yaw (positive amount).
func (m *Module) camTurnRight(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CAMERA.TURNRIGHT expects (cam, amount#)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CAMERA.TURNRIGHT: first argument must be a camera handle")
	}
	if _, err := heap.Cast[*camObj](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	amt, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.TURNRIGHT: amount must be numeric")
	}
	return value.FromFloat(math.Abs(amt)), nil
}

// camOrbitCamera combines mouse X delta × sensitivity with Q/E yaw (degrees/sec), same as
// FLOAT(Input.MouseDeltaX()) * sens + Input.Orbit(KEY_Q, KEY_E, keyDegPerSec#, dt#).
// Validates the camera handle; returns radians to add to your camYaw# variable.
func (m *Module) camOrbitCamera(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CAMERA.ORBITCAMERA expects (cam, mouseSens#, keyDegPerSec#, dt#)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CAMERA.ORBITCAMERA: first argument must be a camera handle")
	}
	if _, err := heap.Cast[*camObj](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	mouseSens, ok1 := args[1].ToFloat()
	keyDegPerSec, ok2 := args[2].ToFloat()
	dt, ok3 := args[3].ToFloat()
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("CAMERA.ORBITCAMERA: mouseSens, keyDegPerSec, dt must be numeric")
	}
	d := rl.GetMouseDelta()
	mousePart := float64(d.X) * mouseSens
	var ax float64
	if rl.IsKeyDown(rl.KeyQ) && !rl.IsKeyDown(rl.KeyE) {
		ax = -1
	} else if rl.IsKeyDown(rl.KeyE) && !rl.IsKeyDown(rl.KeyQ) {
		ax = 1
	}
	keyPart := ax * (keyDegPerSec * math.Pi / 180.0) * dt
	return value.FromFloat(mousePart + keyPart), nil
}

func (m *Module) registerCameraExtras(reg runtime.Registrar) {
	reg.Register("CAMERA.GETPOS", "camera", runtime.AdaptLegacy(m.camGetPos))
	reg.Register("CAMERA.GETTARGET", "camera", runtime.AdaptLegacy(m.camGetTarget))
	reg.Register("CAMERA.SETUP", "camera", runtime.AdaptLegacy(m.camSetUp))
	reg.Register("CAMERA.SETORBIT", "camera", runtime.AdaptLegacy(m.camSetOrbit))
	reg.Register("CAMERA.ORBITAROUND", "camera", runtime.AdaptLegacy(m.camOrbitAround))
	reg.Register("CAMERA.ORBITAROUNDEG", "camera", runtime.AdaptLegacy(m.camOrbitAroundDeg))
	reg.Register("CAMERA.TURNLEFT", "camera", m.camTurnLeft)
	reg.Register("CAMERA.TURNRIGHT", "camera", m.camTurnRight)
	reg.Register("CAMERA.ORBITCAMERA", "camera", m.camOrbitCamera)
	reg.Register("CAMERA.SMOOTHEXP", "camera", m.camSmoothExp)
	reg.Register("CAMERA.FREE", "camera", runtime.AdaptLegacy(m.camFree))
	reg.Register("CAMERA.YAW", "camera", runtime.AdaptLegacy(m.camGetYaw))
	reg.Register("CAMERA.GETYAW", "camera", runtime.AdaptLegacy(m.camGetYaw))
	reg.Register("CAMERA.USEMOUSEORBIT", "camera", runtime.AdaptLegacy(m.camUseMouseOrbit))
	reg.Register("CAMERA.USEORBITRIGHTMOUSE", "camera", runtime.AdaptLegacy(m.camUseOrbitRightMouse))
	reg.Register("CAMERA.SETORBITKEYS", "camera", runtime.AdaptLegacy(m.camSetOrbitKeys))
	reg.Register("CAMERA.SETORBITLIMITS", "camera", runtime.AdaptLegacy(m.camSetOrbitLimits))
	reg.Register("CAMERA.SETORBITSPEED", "camera", runtime.AdaptLegacy(m.camSetOrbitSpeed))
	reg.Register("CAMERA.SETORBITKEYSPEED", "camera", runtime.AdaptLegacy(m.camSetOrbitKeySpeed))
}

func (m *Module) camGetPos(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.GETPOS: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CAMERA.GETPOS expects camera handle")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.GETPOS: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	p := o.cam.Position
	return mbmatrix.AllocVec3Value(m.h, p.X, p.Y, p.Z)
}

func (m *Module) camGetTarget(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.GETTARGET: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CAMERA.GETTARGET expects camera handle")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.GETTARGET: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	t := o.cam.Target
	return mbmatrix.AllocVec3Value(m.h, t.X, t.Y, t.Z)
}

// camSetOrbit places the camera on a spherical shell around (tx,ty,tz): yaw about +Y, pitch
// elevation, distance. Yaw 0 and pitch 0 puts the camera on +Z relative to the target; yaw
// follows the same convention as typical third-person demos (sin/cos on XZ).
func (m *Module) camSetOrbit(args []value.Value) (value.Value, error) {
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("CAMERA.SETORBIT expects 7 arguments (handle, tx, ty, tz, yaw#, pitch#, distance#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.SETORBIT: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	tx, ok1 := argF(args[1])
	ty, ok2 := argF(args[2])
	tz, ok3 := argF(args[3])
	yaw, ok4 := argF(args[4])
	pitch, ok5 := argF(args[5])
	dist, ok6 := argF(args[6])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("CAMERA.SETORBIT: tx..distance must be numeric")
	}
	if dist < 0.15 {
		dist = 0.15
	}
	maxPitch := float32(1.45)
	if pitch > maxPitch {
		pitch = maxPitch
	}
	if pitch < -maxPitch {
		pitch = -maxPitch
	}
	sy, cy := math.Sin(float64(yaw)), math.Cos(float64(yaw))
	sp, cp := math.Sin(float64(pitch)), math.Cos(float64(pitch))
	hdist := float64(dist) * cp
	px := float32(float64(tx) + sy*hdist)
	py := float32(float64(ty) + float64(dist)*sp)
	pz := float32(float64(tz) + cy*hdist)
	o.cam.Position = rl.Vector3{X: px, Y: py, Z: pz}
	o.cam.Target = rl.Vector3{X: tx, Y: ty, Z: tz}
	return value.Nil, nil
}

// ApplySetOrbit applies the same math as CAMERA.SETORBIT (exported for ENTITY CAMERA.ORBITENTITY).
func ApplySetOrbit(h *heap.Store, camHandle heap.Handle, tx, ty, tz, yaw, pitch, dist float32) error {
	o, err := heap.Cast[*camObj](h, camHandle)
	if err != nil {
		return err
	}
	if dist < 0.15 {
		dist = 0.15
	}
	maxPitch := float32(1.45)
	if pitch > maxPitch {
		pitch = maxPitch
	}
	if pitch < -maxPitch {
		pitch = -maxPitch
	}
	sy, cy := math.Sin(float64(yaw)), math.Cos(float64(yaw))
	sp, cp := math.Sin(float64(pitch)), math.Cos(float64(pitch))
	hdist := float64(dist) * cp
	px := float32(float64(tx) + sy*hdist)
	py := float32(float64(ty) + float64(dist)*sp)
	pz := float32(float64(tz) + cy*hdist)
	o.cam.Position = rl.Vector3{X: px, Y: py, Z: pz}
	o.cam.Target = rl.Vector3{X: tx, Y: ty, Z: tz}
	return nil
}

// ApplySetTarget sets the camera look-at target in world space (for ENTITY.CAMERA.SETTARGETENTITY).
func ApplySetTarget(h *heap.Store, camHandle heap.Handle, tx, ty, tz float32) error {
	o, err := heap.Cast[*camObj](h, camHandle)
	if err != nil {
		return err
	}
	o.cam.Target = rl.Vector3{X: tx, Y: ty, Z: tz}
	return nil
}

// camOrbitAround: third-person on XZ ring — camera Y is absolute height#; yaw radians on XZ at distance dist.
func (m *Module) camOrbitAround(args []value.Value) (value.Value, error) {
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("CAMERA.ORBITAROUND expects 7 arguments (handle, tx, ty, tz, yaw#, distance#, cameraY#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.ORBITAROUND: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	tx, ok1 := argF(args[1])
	ty, ok2 := argF(args[2])
	tz, ok3 := argF(args[3])
	yaw, ok4 := argF(args[4])
	dist, ok5 := argF(args[5])
	cy, ok6 := argF(args[6])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("CAMERA.ORBITAROUND: numeric arguments required")
	}
	if dist < 0.05 {
		dist = 0.05
	}
	sy, cyaw := math.Sin(float64(yaw)), math.Cos(float64(yaw))
	px := float32(float64(tx) + sy*float64(dist))
	pz := float32(float64(tz) + cyaw*float64(dist))
	o.cam.Position = rl.Vector3{X: px, Y: cy, Z: pz}
	o.cam.Target = rl.Vector3{X: tx, Y: ty, Z: tz}
	return value.Nil, nil
}

func (m *Module) camOrbitAroundDeg(args []value.Value) (value.Value, error) {
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("CAMERA.ORBITAROUNDEG expects 7 arguments (handle, tx, ty, tz, yawDeg#, distance#, cameraY#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.ORBITAROUNDEG: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	tx, ok1 := argF(args[1])
	ty, ok2 := argF(args[2])
	tz, ok3 := argF(args[3])
	yawDeg, ok4 := argF(args[4])
	dist, ok5 := argF(args[5])
	camY, ok6 := argF(args[6])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("CAMERA.ORBITAROUNDEG: numeric arguments required")
	}
	yaw := float64(yawDeg) * math.Pi / 180.0
	if dist < 0.05 {
		dist = 0.05
	}
	sy, cyaw := math.Sin(yaw), math.Cos(yaw)
	px := float32(float64(tx) + sy*float64(dist))
	pz := float32(float64(tz) + cyaw*float64(dist))
	o.cam.Position = rl.Vector3{X: px, Y: camY, Z: pz}
	o.cam.Target = rl.Vector3{X: tx, Y: ty, Z: tz}
	return value.Nil, nil
}

func (m *Module) camSetUp(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CAMERA.SETUP expects (handle, ux, uy, uz)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.SETUP: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	ux, ok1 := argF(args[1])
	uy, ok2 := argF(args[2])
	uz, ok3 := argF(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("CAMERA.SETUP: up vector must be numeric")
	}
	o.cam.Up = rl.Vector3{X: ux, Y: uy, Z: uz}
	return value.Nil, nil
}

func (m *Module) camFree(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.FREE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CAMERA.FREE expects camera handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}
