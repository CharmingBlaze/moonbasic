//go:build cgo || (windows && !cgo)

package mbcamera

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type cam2dObj struct {
	cam rl.Camera2D
}

func (o *cam2dObj) TypeName() string { return "Camera2D" }

func (o *cam2dObj) TypeTag() uint16 { return heap.TagCamera2D }

func (o *cam2dObj) Free() {}

func (m *Module) getCam2D(args []value.Value, ix int, op string) (*cam2dObj, error) {
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: expected Camera2D handle", op)
	}
	return heap.Cast[*cam2dObj](m.h, heap.Handle(args[ix].IVal))
}

func (m *Module) registerCamera2D(reg runtime.Registrar) {
	reg.Register("CAMERA2D.MAKE", "camera", runtime.AdaptLegacy(m.cam2dMake))
	reg.Register("CAMERA2D.SETTARGET", "camera", runtime.AdaptLegacy(m.cam2dSetTarget))
	reg.Register("CAMERA2D.SETOFFSET", "camera", runtime.AdaptLegacy(m.cam2dSetOffset))
	reg.Register("CAMERA2D.SETZOOM", "camera", runtime.AdaptLegacy(m.cam2dSetZoom))
	reg.Register("CAMERA2D.SETROTATION", "camera", runtime.AdaptLegacy(m.cam2dSetRotation))
	reg.Register("CAMERA2D.BEGIN", "camera", runtime.AdaptLegacy(m.cam2dBegin))
	reg.Register("CAMERA2D.END", "camera", runtime.AdaptLegacy(m.cam2dEnd))
	reg.Register("CAMERA2D.GETMATRIX", "camera", runtime.AdaptLegacy(m.cam2dGetMatrix))
	reg.Register("CAMERA2D.WORLDTOSCREEN", "camera", runtime.AdaptLegacy(m.cam2dWorldToScreen))
	reg.Register("CAMERA2D.SCREENTOWORLD", "camera", runtime.AdaptLegacy(m.cam2dScreenToWorld))
	reg.Register("CAMERA2D.FREE", "camera", runtime.AdaptLegacy(m.cam2dFree))
}

func (m *Module) cam2dMake(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA2D.MAKE: heap not bound")
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CAMERA2D.MAKE expects 0 arguments")
	}
	w := float32(rl.GetScreenWidth())
	h := float32(rl.GetScreenHeight())
	if w < 1 {
		w = 800
	}
	if h < 1 {
		h = 450
	}
	o := &cam2dObj{
		cam: rl.Camera2D{
			Offset:   rl.Vector2{X: w * 0.5, Y: h * 0.5},
			Target:   rl.Vector2{X: 0, Y: 0},
			Rotation: 0,
			Zoom:     1,
		},
	}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) cam2dSetTarget(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CAMERA2D.SETTARGET expects (camera, x, y)")
	}
	o, err := m.getCam2D(args, 0, "CAMERA2D.SETTARGET")
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argF(args[1])
	y, ok2 := argF(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("CAMERA2D.SETTARGET: x, y must be numeric")
	}
	o.cam.Target = rl.Vector2{X: x, Y: y}
	return value.Nil, nil
}

func (m *Module) cam2dSetOffset(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CAMERA2D.SETOFFSET expects (camera, x, y)")
	}
	o, err := m.getCam2D(args, 0, "CAMERA2D.SETOFFSET")
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argF(args[1])
	y, ok2 := argF(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("CAMERA2D.SETOFFSET: x, y must be numeric")
	}
	o.cam.Offset = rl.Vector2{X: x, Y: y}
	return value.Nil, nil
}

func (m *Module) cam2dSetZoom(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CAMERA2D.SETZOOM expects (camera, zoom)")
	}
	o, err := m.getCam2D(args, 0, "CAMERA2D.SETZOOM")
	if err != nil {
		return value.Nil, err
	}
	z, ok := argF(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA2D.SETZOOM: zoom must be numeric")
	}
	if z <= 0 {
		z = 0.01
	}
	o.cam.Zoom = z
	return value.Nil, nil
}

func (m *Module) cam2dSetRotation(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CAMERA2D.SETROTATION expects (camera, degrees)")
	}
	o, err := m.getCam2D(args, 0, "CAMERA2D.SETROTATION")
	if err != nil {
		return value.Nil, err
	}
	rad, ok := argF(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA2D.SETROTATION: degrees must be numeric")
	}
	o.cam.Rotation = rad
	return value.Nil, nil
}

func (m *Module) cam2dBegin(args []value.Value) (value.Value, error) {
	switch len(args) {
	case 0:
		cam := rl.Camera2D{
			Offset:   rl.Vector2{X: 0, Y: 0},
			Target:   rl.Vector2{X: 0, Y: 0},
			Rotation: 0,
			Zoom:     1,
		}
		rl.BeginMode2D(cam)
		return value.Nil, nil
	case 1:
		h, ok := argHandle(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("CAMERA2D.BEGIN: invalid handle")
		}
		o, err := heap.Cast[*cam2dObj](m.h, h)
		if err != nil {
			return value.Nil, err
		}
		rl.BeginMode2D(o.cam)
		return value.Nil, nil
	default:
		return value.Nil, fmt.Errorf("CAMERA2D.BEGIN expects 0 arguments (identity) or 1 (camera handle)")
	}
}

func (m *Module) cam2dEnd(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CAMERA2D.END expects 0 arguments")
	}
	rl.EndMode2D()
	return value.Nil, nil
}

func (m *Module) cam2dFree(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA2D.FREE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CAMERA2D.FREE expects camera handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) cam2dGetMatrix(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA2D.GETMATRIX: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CAMERA2D.GETMATRIX expects camera handle")
	}
	o, err := m.getCam2D(args, 0, "CAMERA2D.GETMATRIX")
	if err != nil {
		return value.Nil, err
	}
	mat := rl.GetCameraMatrix2D(o.cam)
	id, err := mbmatrix.AllocMatrix(m.h, mat)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) cam2dWorldToScreen(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA2D.WORLDTOSCREEN: heap not bound")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CAMERA2D.WORLDTOSCREEN expects (camera, worldX, worldY)")
	}
	o, err := m.getCam2D(args, 0, "CAMERA2D.WORLDTOSCREEN")
	if err != nil {
		return value.Nil, err
	}
	wx, ok1 := argF(args[1])
	wy, ok2 := argF(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("CAMERA2D.WORLDTOSCREEN: coordinates must be numeric")
	}
	v := rl.GetWorldToScreen2D(rl.Vector2{X: wx, Y: wy}, o.cam)
	return m.allocVec2Arr(v.X, v.Y)
}

func (m *Module) cam2dScreenToWorld(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA2D.SCREENTOWORLD: heap not bound")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CAMERA2D.SCREENTOWORLD expects (camera, screenX, screenY)")
	}
	o, err := m.getCam2D(args, 0, "CAMERA2D.SCREENTOWORLD")
	if err != nil {
		return value.Nil, err
	}
	sx, ok1 := argF(args[1])
	sy, ok2 := argF(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("CAMERA2D.SCREENTOWORLD: coordinates must be numeric")
	}
	v := rl.GetScreenToWorld2D(rl.Vector2{X: sx, Y: sy}, o.cam)
	return m.allocVec2Arr(v.X, v.Y)
}

func (m *Module) allocVec2Arr(x, y float32) (value.Value, error) {
	a, err := heap.NewArray([]int64{2})
	if err != nil {
		return value.Nil, err
	}
	_ = a.Set([]int64{0}, float64(x))
	_ = a.Set([]int64{1}, float64(y))
	id, err := m.h.Alloc(a)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}
