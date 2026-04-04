//go:build cgo

package mbcamera

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
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
