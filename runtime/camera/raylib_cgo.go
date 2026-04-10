//go:build cgo || (windows && !cgo)

package mbcamera

import (
	"fmt"
	"math"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/runtime/mbmodel3d"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type camObj struct {
	cam rl.Camera3D

	shakeMag  float32
	shakeTime float32

	useClip       bool
	clipNear      float64
	clipFar       float64
	fpsMode       bool
	fpsSensitivity float32
}

func (c *camObj) TypeName() string { return "Camera3D" }

func (c *camObj) TypeTag() uint16 { return heap.TagCamera }

func (c *camObj) Free() {}

func argF(v value.Value) (float32, bool) {
	if f, ok := v.ToFloat(); ok {
		return float32(f), true
	}
	if i, ok := v.ToInt(); ok {
		return float32(i), true
	}
	return 0, false
}

func argHandle(v value.Value) (heap.Handle, bool) {
	if v.Kind != value.KindHandle {
		return 0, false
	}
	return heap.Handle(v.IVal), true
}

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	r.Register("CAMERA.MAKE", "camera", runtime.AdaptLegacy(m.camMake))
	r.Register("CAMERA.CREATE", "camera", runtime.AdaptLegacy(m.camMake))
	r.Register("CAM", "camera", runtime.AdaptLegacy(m.camMake))
	r.Register("CAMERA.SETPOS", "camera", runtime.AdaptLegacy(m.camSetPos))
	r.Register("CAMERA.SETPOSITION", "camera", runtime.AdaptLegacy(m.camSetPos))
	r.Register("CAMERA.SETTARGET", "camera", runtime.AdaptLegacy(m.camSetTarget))
	r.Register("CAMERA.LOOKAT", "camera", runtime.AdaptLegacy(m.camSetTarget))
	r.Register("CAMERA.SETPROJECTION", "camera", runtime.AdaptLegacy(m.camSetProjection))
	r.Register("CAMERA.SETMODE", "camera", m.camSetMode)
	r.Register("CAMERA.SETFOV", "camera", runtime.AdaptLegacy(m.camSetFov))
	r.Register("CAMERA.BEGIN", "camera", runtime.AdaptLegacy(m.camBegin))
	r.Register("CAMERA.END", "camera", runtime.AdaptLegacy(m.camEnd))
	r.Register("CAMERA.MOVE", "camera", runtime.AdaptLegacy(m.camMove))
	r.Register("CAMERA.GETRAY", "camera", runtime.AdaptLegacy(m.camGetRay))
	r.Register("CAMERA.UNPROJECT", "camera", runtime.AdaptLegacy(m.camGetRay))
	r.Register("CAMERA.PICK", "camera", runtime.AdaptLegacy(m.camGetRay))
	r.Register("CAMERA.SHAKE", "camera", runtime.AdaptLegacy(m.camShake))
	r.Register("CameraFOV", "camera", runtime.AdaptLegacy(m.camSetFov))
	r.Register("CameraShake", "camera", runtime.AdaptLegacy(m.camShake))
	r.Register("CameraLookAt", "camera", runtime.AdaptLegacy(m.camSetTarget))
	r.Register("CAMERA.GETVIEWRAY", "camera", runtime.AdaptLegacy(m.camGetViewRay))
	r.Register("CAMERA.GETMATRIX", "camera", runtime.AdaptLegacy(m.camGetMatrix))
	r.Register("MATRIX.FREE", "camera", runtime.AdaptLegacy(m.matrixFree))
	m.registerCameraExtras(r)
	m.registerBlitzCamera(r)
	registerBlitzCameraExtras(m, r)
	m.registerScreenHelpers(r)
	m.registerCamera2D(r)
	m.registerCull(r)
	registerCameraMore(m, r)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

func (m *Module) camMake(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.MAKE: heap not bound")
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CAMERA.MAKE expects 0 arguments")
	}
	o := &camObj{
		cam: rl.Camera3D{
			Position:   rl.Vector3{X: 0, Y: 2, Z: 8},
			Target:     rl.Vector3{X: 0, Y: 0, Z: 0},
			Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
			Fovy:       45,
			Projection: rl.CameraPerspective,
		},
	}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) camSetPos(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CAMERA.SETPOS expects 4 arguments (handle, x, y, z)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.SETPOS: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argF(args[1])
	y, ok2 := argF(args[2])
	z, ok3 := argF(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("CAMERA.SETPOS: non-numeric position")
	}
	o.cam.Position = rl.Vector3{X: x, Y: y, Z: z}
	return value.Nil, nil
}

func (m *Module) camSetTarget(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CAMERA.SETTARGET expects 4 arguments (handle, x, y, z)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.SETTARGET: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argF(args[1])
	y, ok2 := argF(args[2])
	z, ok3 := argF(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("CAMERA.SETTARGET: non-numeric target")
	}
	o.cam.Target = rl.Vector3{X: x, Y: y, Z: z}
	return value.Nil, nil
}

func (m *Module) camSetProjection(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CAMERA.SETPROJECTION expects (handle, mode#): 0 perspective, 1 orthographic")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.SETPROJECTION: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	var mode int64
	if i, ok := args[1].ToInt(); ok {
		mode = i
	} else if f, ok := args[1].ToFloat(); ok {
		mode = int64(f)
	} else {
		return value.Nil, fmt.Errorf("CAMERA.SETPROJECTION: mode must be numeric (0 or 1)")
	}
	switch mode {
	case 0:
		o.cam.Projection = rl.CameraPerspective
	case 1:
		o.cam.Projection = rl.CameraOrthographic
	default:
		return value.Nil, fmt.Errorf("CAMERA.SETPROJECTION: mode must be 0 (perspective) or 1 (orthographic)")
	}
	return value.Nil, nil
}

// camSetMode: alias-friendly projection picker — numeric 0/1 or strings "perspective"/"orthographic".
func (m *Module) camSetMode(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CAMERA.SETMODE expects (handle, mode): 0 perspective, 1 orthographic, or those names as a string")
	}
	if args[1].Kind == value.KindString {
		s, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Nil, err
		}
		switch strings.ToLower(strings.TrimSpace(s)) {
		case "perspective", "camera_perspective", "persp", "proj_perspective":
			return m.camSetProjection([]value.Value{args[0], value.FromInt(0)})
		case "orthographic", "ortho", "camera_orthographic", "proj_orthographic":
			return m.camSetProjection([]value.Value{args[0], value.FromInt(1)})
		default:
			return value.Nil, fmt.Errorf("CAMERA.SETMODE: unknown mode %q (use perspective or orthographic)", s)
		}
	}
	return m.camSetProjection(args)
}

func (m *Module) camSetFov(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CAMERA.SETFOV expects 2 arguments (handle, fovy)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.SETFOV: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	fov, ok := argF(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.SETFOV: fovy must be numeric")
	}
	o.cam.Fovy = fov
	return value.Nil, nil
}

// CameraWorldPosition returns the world-space position of a CAMERA.MAKE handle (for ENTITY follow helpers).
func CameraWorldPosition(h *heap.Store, camH heap.Handle) (rl.Vector3, bool) {
	o, err := heap.Cast[*camObj](h, camH)
	if err != nil {
		return rl.Vector3{}, false
	}
	return o.cam.Position, true
}

func (m *Module) camShake(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CAMERA.SHAKE expects (camera, amount#, duration#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.SHAKE: invalid camera handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	amt, ok1 := argF(args[1])
	dur, ok2 := argF(args[2])
	if !ok1 || !ok2 || dur < 0 {
		return value.Nil, fmt.Errorf("CAMERA.SHAKE: amount and duration must be numeric")
	}
	o.shakeMag = amt
	o.shakeTime = dur
	return value.Nil, nil
}

func (m *Module) camBegin(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CAMERA.BEGIN expects 1 argument (handle)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.BEGIN: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	cam := o.cam
	if o.shakeTime > 0 && o.shakeMag > 0 {
		dt := rl.GetFrameTime()
		o.shakeTime -= dt
		if o.shakeTime < 0 {
			o.shakeTime = 0
		}
		t := float64(rl.GetTime())
		mag := float64(o.shakeMag) * 0.02
		ox := float32(math.Sin(t*50.0) * mag)
		oy := float32(math.Cos(t*43.0) * mag)
		oz := float32(math.Sin(t*37.0) * mag)
		cam.Position.X += ox
		cam.Position.Y += oy
		cam.Position.Z += oz
		cam.Target.X += ox * 0.5
		cam.Target.Y += oy * 0.5
		cam.Target.Z += oz * 0.5
	}
	if o.useClip {
		rl.SetClipPlanes(o.clipNear, o.clipFar)
	}
	mbmodel3d.MarkCamera3DBegin(cam.Position.X, cam.Position.Y, cam.Position.Z)
	mbmodel3d.StoreActiveCamera3D(cam)
	m.lastActive3D = h
	rl.BeginMode3D(cam)
	rw := float32(rl.GetRenderWidth())
	rh := float32(rl.GetRenderHeight())
	aspect := float32(16.0 / 9.0)
	if rh > 0 {
		aspect = rw / rh
	}
	setActiveFrustum(cam, aspect)
	return value.Nil, nil
}

func (m *Module) camEnd(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CAMERA.END expects 0 arguments")
	}
	mbmodel3d.FlushDeferred3D(m.h)
	clearActiveFrustum()
	mbmodel3d.MarkCamera3DEnd()
	rl.EndMode3D()
	return value.Nil, nil
}

func (m *Module) camMove(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CAMERA.MOVE expects 4 arguments (handle, dx, dy, dz)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.MOVE: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	dx, ok1 := argF(args[1])
	dy, ok2 := argF(args[2])
	dz, ok3 := argF(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("CAMERA.MOVE: non-numeric delta")
	}
	o.cam.Position.X += dx
	o.cam.Position.Y += dy
	o.cam.Position.Z += dz
	o.cam.Target.X += dx
	o.cam.Target.Y += dy
	o.cam.Target.Z += dz
	return value.Nil, nil
}

func (m *Module) camGetRay(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.GETRAY: heap not bound")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CAMERA.GETRAY expects 3 arguments (handle, screenX, screenY)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.GETRAY: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	sx, ok1 := argF(args[1])
	sy, ok2 := argF(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("CAMERA.GETRAY: screen position must be numeric")
	}
	ray := rl.GetScreenToWorldRayEx(rl.Vector2{X: sx, Y: sy}, o.cam, int32(rl.GetRenderWidth()), int32(rl.GetRenderHeight()))
	return m.allocRayHandle(ray)
}

func (m *Module) camGetViewRay(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.GETVIEWRAY: heap not bound")
	}
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("CAMERA.GETVIEWRAY expects 5 arguments (screenX, screenY, handle, width, height)")
	}
	sx, ok1 := argF(args[0])
	sy, ok2 := argF(args[1])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("CAMERA.GETVIEWRAY: screen position must be numeric")
	}
	h, ok := argHandle(args[2])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.GETVIEWRAY: invalid camera handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	var wi int64
	if i, ok := args[3].ToInt(); ok {
		wi = i
	} else if wf, ok := args[3].ToFloat(); ok {
		wi = int64(wf)
	} else {
		return value.Nil, fmt.Errorf("CAMERA.GETVIEWRAY: width must be numeric")
	}
	var hi int64
	if i, ok := args[4].ToInt(); ok {
		hi = i
	} else if hf, ok := args[4].ToFloat(); ok {
		hi = int64(hf)
	} else {
		return value.Nil, fmt.Errorf("CAMERA.GETVIEWRAY: height must be numeric")
	}
	if wi <= 0 || hi <= 0 {
		return value.Nil, fmt.Errorf("CAMERA.GETVIEWRAY: width and height must be positive")
	}
	ray := rl.GetScreenToWorldRayEx(rl.Vector2{X: sx, Y: sy}, o.cam, int32(wi), int32(hi))
	return m.allocRayHandle(ray)
}

func (m *Module) allocRayHandle(ray rl.Ray) (value.Value, error) {
	arr, err := heap.NewArray([]int64{6})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, float64(ray.Position.X))
	_ = arr.Set([]int64{1}, float64(ray.Position.Y))
	_ = arr.Set([]int64{2}, float64(ray.Position.Z))
	_ = arr.Set([]int64{3}, float64(ray.Direction.X))
	_ = arr.Set([]int64{4}, float64(ray.Direction.Y))
	_ = arr.Set([]int64{5}, float64(ray.Direction.Z))
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) camGetMatrix(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.GETMATRIX: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CAMERA.GETMATRIX expects camera handle")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.GETMATRIX: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	mat := rl.GetCameraMatrix(o.cam)
	id, err := mbmatrix.AllocMatrix(m.h, mat)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) matrixFree(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("MATRIX.FREE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("MATRIX.FREE expects matrix handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}
