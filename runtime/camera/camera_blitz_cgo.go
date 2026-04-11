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

func (m *Module) registerBlitzCamera(reg runtime.Registrar) {
	reg.Register("CAMERA.TURN", "camera", runtime.AdaptLegacy(m.camTurn))
	reg.Register("CAMERA.ROTATE", "camera", runtime.AdaptLegacy(m.camRotateAbs))
	reg.Register("CAMERA.ORBIT", "camera", runtime.AdaptLegacy(m.camOrbitDispatch))
	reg.Register("CAMERA.ZOOM", "camera", runtime.AdaptLegacy(m.camZoom))
	reg.Register("CAMERA.FOLLOW", "camera", m.camFollow)
}

// ThirdPersonFollowStep lerps the camera toward a third-person position behind (tx,ty,tz) on the XZ plane
// at yaw, with fixed camera height. Target lerps toward the subject. Exported for mbentity.
func ThirdPersonFollowStep(h *heap.Store, camHandle heap.Handle, tx, ty, tz, yaw, dist, height, smooth float32, dt float64) error {
	o, err := heap.Cast[*camObj](h, camHandle)
	if err != nil {
		return err
	}
	if dist < 0.05 {
		dist = 0.05
	}
	sy, cy := math.Sin(float64(yaw)), math.Cos(float64(yaw))
	wantPX := float32(float64(tx) - sy*float64(dist))
	wantPZ := float32(float64(tz) - cy*float64(dist))
	wantPY := height
	wantTgt := rl.Vector3{X: tx, Y: ty, Z: tz}

	k := smooth * 8.0 * float32(dt)
	if k > 1 {
		k = 1
	}
	if k < 0 {
		k = 0
	}
	wantPos := rl.Vector3{X: wantPX, Y: wantPY, Z: wantPZ}
	o.cam.Position = rl.Vector3Lerp(o.cam.Position, wantPos, k)
	o.cam.Target = rl.Vector3Lerp(o.cam.Target, wantTgt, k)
	return nil
}

// camTurn: incremental yaw (around world +Y), pitch (around camera right), roll (around view axis). Angles in radians.
func (m *Module) camTurn(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CAMERA.TURN expects 4 arguments (handle, dpitch#, dyaw#, droll#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.TURN: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	dp, ok1 := argF(args[1])
	dy, ok2 := argF(args[2])
	dr, ok3 := argF(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("CAMERA.TURN: angles must be numeric")
	}

	pos := o.cam.Position
	tgt := o.cam.Target
	fwd := rl.Vector3Subtract(tgt, pos)
	dist := rl.Vector3Length(fwd)
	if dist < 1e-4 {
		dist = 1
	}
	fwd = rl.Vector3Normalize(fwd)
	worldUp := rl.Vector3{X: 0, Y: 1, Z: 0}

	fwd = rl.Vector3RotateByAxisAngle(fwd, worldUp, dy)

	right := rl.Vector3Normalize(rl.Vector3CrossProduct(worldUp, fwd))
	if rl.Vector3Length(right) < 1e-6 {
		right = rl.Vector3{X: 1, Y: 0, Z: 0}
	}
	fwd = rl.Vector3RotateByAxisAngle(fwd, right, dp)
	fwd = rl.Vector3Normalize(fwd)

	up := rl.Vector3Normalize(rl.Vector3CrossProduct(right, fwd))
	up = rl.Vector3RotateByAxisAngle(up, fwd, dr)
	up = rl.Vector3Normalize(up)

	o.cam.Position = pos
	o.cam.Target = rl.Vector3Add(pos, rl.Vector3Scale(fwd, dist))
	o.cam.Up = up
	return value.Nil, nil
}

// camRotateAbs: absolute orientation from pitch, yaw, roll (radians). Y-up, yaw about Y, pitch about X, roll about view.
// Camera position unchanged; target = position + forward * current distance.
func (m *Module) camRotateAbs(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CAMERA.ROTATE expects 4 arguments (handle, pitch#, yaw#, roll#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.ROTATE: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	pitch, ok1 := argF(args[1])
	yaw, ok2 := argF(args[2])
	roll, ok3 := argF(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("CAMERA.ROTATE: angles must be numeric")
	}

	pos := o.cam.Position
	tgt := o.cam.Target
	fwd := rl.Vector3Subtract(tgt, pos)
	dist := rl.Vector3Length(fwd)
	if dist < 1e-4 {
		dist = 1
	}

	cp := float64(math.Cos(float64(pitch)))
	sp := float64(math.Sin(float64(pitch)))
	sy := float64(math.Sin(float64(yaw)))
	cy := float64(math.Cos(float64(yaw)))

	fx := float32(sy * cp)
	fy := float32(sp)
	fz := float32(cy * cp)
	fwd = rl.Vector3Normalize(rl.Vector3{X: fx, Y: fy, Z: fz})

	worldUp := rl.Vector3{X: 0, Y: 1, Z: 0}
	right := rl.Vector3Normalize(rl.Vector3CrossProduct(worldUp, fwd))
	if rl.Vector3Length(right) < 1e-6 {
		right = rl.Vector3{X: 1, Y: 0, Z: 0}
	}
	up := rl.Vector3Normalize(rl.Vector3CrossProduct(right, fwd))
	up = rl.Vector3RotateByAxisAngle(up, fwd, roll)
	up = rl.Vector3Normalize(up)

	o.cam.Position = pos
	o.cam.Target = rl.Vector3Add(pos, rl.Vector3Scale(fwd, dist))
	o.cam.Up = up
	return value.Nil, nil
}

func (m *Module) camZoom(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CAMERA.ZOOM expects 2 arguments (handle, amount#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.ZOOM: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	amt, ok1 := argF(args[1])
	if !ok1 {
		return value.Nil, fmt.Errorf("CAMERA.ZOOM: amount must be numeric")
	}
	o.cam.Fovy += amt
	if o.cam.Fovy < 10 {
		o.cam.Fovy = 10
	}
	if o.cam.Fovy > 120 {
		o.cam.Fovy = 120
	}
	return value.Nil, nil
}

// camFollow: third-person follow using world target (tx,ty,tz), horizontal yaw, distance, absolute camera Y, smoothness 0..1, dt from runtime.
func (m *Module) camFollow(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 8 {
		return value.Nil, fmt.Errorf("CAMERA.FOLLOW expects 8 arguments (handle, tx#, ty#, tz#, yaw#, dist#, height#, smooth#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.FOLLOW: invalid handle")
	}
	tx, ok1 := argF(args[1])
	ty, ok2 := argF(args[2])
	tz, ok3 := argF(args[3])
	yaw, ok4 := argF(args[4])
	dist, ok5 := argF(args[5])
	height, ok6 := argF(args[6])
	smooth, ok7 := argF(args[7])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 {
		return value.Nil, fmt.Errorf("CAMERA.FOLLOW: numeric arguments required")
	}
	dt := mbtime.DeltaSeconds(rt)
	if err := ThirdPersonFollowStep(m.h, h, tx, ty, tz, yaw, dist, height, smooth, dt); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}
