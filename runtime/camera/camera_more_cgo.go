//go:build cgo || (windows && !cgo)

package mbcamera

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	mbsprite "moonbasic/runtime/sprite"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerCameraMore(m *Module, reg runtime.Registrar) {
	reg.Register("CAMERA.SETRANGE", "camera", runtime.AdaptLegacy(m.camSetRange))
	reg.Register("CAMERA.SETACTIVE", "camera", runtime.AdaptLegacy(m.camSetActive))
	reg.Register("CAMERA.GETACTIVE", "camera", runtime.AdaptLegacy(m.camGetActive))
	reg.Register("CAMERA.WORLDTOSCREEN2D", "camera", runtime.AdaptLegacy(m.camWorldToScreen))
	reg.Register("CAMERA.SETFPSMODE", "camera", runtime.AdaptLegacy(m.camSetFPSMode))
	reg.Register("CAMERA.CLEARFPSMODE", "camera", runtime.AdaptLegacy(m.camClearFPSMode))
	reg.Register("CAMERA.UPDATEFPS", "camera", runtime.AdaptLegacy(m.camUpdateFPS))

	reg.Register("CAMERA2D.FOLLOW", "camera", m.cam2dFollow)
	reg.Register("CAMERA2D.ZOOMTOMOUSE", "camera", runtime.AdaptLegacy(m.cam2dZoomToMouse))
	reg.Register("CAMERA2D.ZOOMIN", "camera", runtime.AdaptLegacy(m.cam2dZoomIn))
	reg.Register("CAMERA2D.ZOOMOUT", "camera", runtime.AdaptLegacy(m.cam2dZoomOut))
	reg.Register("CAMERA2D.ROTATION", "camera", runtime.AdaptLegacy(m.cam2dRotation))
	reg.Register("CAMERA2D.TARGETX", "camera", runtime.AdaptLegacy(m.cam2dTargetX))
	reg.Register("CAMERA2D.TARGETY", "camera", runtime.AdaptLegacy(m.cam2dTargetY))
}

func (m *Module) camSetRange(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CAMERA.SETRANGE expects (camera, near#, far#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.SETRANGE: invalid camera handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	near, ok1 := argF(args[1])
	far, ok2 := argF(args[2])
	if !ok1 || !ok2 || near <= 0 || far <= near {
		return value.Nil, fmt.Errorf("CAMERA.SETRANGE: need 0 < near < far")
	}
	o.useClip = true
	o.clipNear = float64(near)
	o.clipFar = float64(far)
	return value.Nil, nil
}

func (m *Module) camSetActive(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CAMERA.SETACTIVE expects camera handle")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.SETACTIVE: invalid handle")
	}
	if _, err := heap.Cast[*camObj](m.h, h); err != nil {
		return value.Nil, err
	}
	m.lastActive3D = h
	return value.Nil, nil
}

func (m *Module) camGetActive(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CAMERA.GETACTIVE expects 0 arguments")
	}
	if m.lastActive3D == 0 {
		return value.FromHandle(0), nil
	}
	return value.FromHandle(int32(m.lastActive3D)), nil
}

func (m *Module) camSetFPSMode(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CAMERA.SETFPSMODE expects (camera, sensitivity#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.SETFPSMODE: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	s, ok := argF(args[1])
	if !ok || s <= 0 {
		return value.Nil, fmt.Errorf("CAMERA.SETFPSMODE: sensitivity must be positive")
	}
	o.fpsMode = true
	o.fpsSensitivity = s
	rl.DisableCursor()
	return value.Nil, nil
}

func (m *Module) camClearFPSMode(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CAMERA.CLEARFPSMODE expects camera handle")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.CLEARFPSMODE: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	o.fpsMode = false
	rl.EnableCursor()
	return value.Nil, nil
}

func (m *Module) camUpdateFPS(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CAMERA.UPDATEFPS expects camera handle")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.UPDATEFPS: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	if !o.fpsMode {
		return value.Nil, nil
	}
	c := o.cam
	rl.UpdateCamera(&c, rl.CameraFirstPerson)
	o.cam = c
	return value.Nil, nil
}

func (m *Module) cam2dFollow(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA2D.FOLLOW: heap not bound")
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CAMERA2D.FOLLOW expects (camera2d, sprite, speed#, dt#)")
	}
	o, err := m.getCam2D(args, 0, "CAMERA2D.FOLLOW")
	if err != nil {
		return value.Nil, err
	}
	if args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CAMERA2D.FOLLOW: sprite must be handle")
	}
	sp := heap.Handle(args[1].IVal)
	speed, ok1 := argF(args[2])
	dt, ok2 := argF(args[3])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("CAMERA2D.FOLLOW: speed and dt must be numeric")
	}
	sx, sy, err := mbsprite.WorldXY(m.h, sp)
	if err != nil {
		return value.Nil, err
	}
	t := o.cam.Target
	t.X += (sx - t.X) * speed * dt
	t.Y += (sy - t.Y) * speed * dt
	o.cam.Target = t
	return value.Nil, nil
}

func (m *Module) cam2dZoomToMouse(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CAMERA2D.ZOOMTOMOUSE expects (camera, delta#)")
	}
	o, err := m.getCam2D(args, 0, "CAMERA2D.ZOOMTOMOUSE")
	if err != nil {
		return value.Nil, err
	}
	dz, ok := argF(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA2D.ZOOMTOMOUSE: delta must be numeric")
	}
	mp := rl.GetMousePosition()
	before := rl.GetScreenToWorld2D(mp, o.cam)
	z := o.cam.Zoom + dz
	if z <= 0.01 {
		z = 0.01
	}
	o.cam.Zoom = z
	after := rl.GetScreenToWorld2D(mp, o.cam)
	o.cam.Target.X += before.X - after.X
	o.cam.Target.Y += before.Y - after.Y
	return value.Nil, nil
}

func (m *Module) cam2dZoomIn(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CAMERA2D.ZOOMIN expects (camera, amount#)")
	}
	o, err := m.getCam2D(args, 0, "CAMERA2D.ZOOMIN")
	if err != nil {
		return value.Nil, err
	}
	a, ok := argF(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA2D.ZOOMIN: amount must be numeric")
	}
	o.cam.Zoom += a
	if o.cam.Zoom <= 0.01 {
		o.cam.Zoom = 0.01
	}
	return value.Nil, nil
}

func (m *Module) cam2dZoomOut(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CAMERA2D.ZOOMOUT expects (camera, amount#)")
	}
	o, err := m.getCam2D(args, 0, "CAMERA2D.ZOOMOUT")
	if err != nil {
		return value.Nil, err
	}
	a, ok := argF(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA2D.ZOOMOUT: amount must be numeric")
	}
	o.cam.Zoom -= a
	if o.cam.Zoom <= 0.01 {
		o.cam.Zoom = 0.01
	}
	return value.Nil, nil
}

func (m *Module) cam2dRotation(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CAMERA2D.ROTATION expects camera handle")
	}
	o, err := m.getCam2D(args, 0, "CAMERA2D.ROTATION")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(o.cam.Rotation)), nil
}

func (m *Module) cam2dTargetX(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CAMERA2D.TARGETX expects camera handle")
	}
	o, err := m.getCam2D(args, 0, "CAMERA2D.TARGETX")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(o.cam.Target.X)), nil
}

func (m *Module) cam2dTargetY(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CAMERA2D.TARGETY expects camera handle")
	}
	o, err := m.getCam2D(args, 0, "CAMERA2D.TARGETY")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(o.cam.Target.Y)), nil
}
