//go:build cgo || (windows && !cgo)

package mbcamera

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerScreenHelpers(reg runtime.Registrar) {
	reg.Register("CAMERA.WORLDTOSCREEN", "camera", runtime.AdaptLegacy(m.camWorldToScreen))
	reg.Register("CAMERA.ISONSCREEN", "camera", runtime.AdaptLegacy(m.camIsOnScreen))
	reg.Register("CAMERA.MOUSERAY", "camera", runtime.AdaptLegacy(m.camMouseRay))
}

func (m *Module) camWorldToScreen(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.WORLDTOSCREEN: heap not bound")
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CAMERA.WORLDTOSCREEN expects (camera, wx, wy, wz)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.WORLDTOSCREEN: invalid camera handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	wx, ok1 := argF(args[1])
	wy, ok2 := argF(args[2])
	wz, ok3 := argF(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("CAMERA.WORLDTOSCREEN: world position must be numeric")
	}
	pos := rl.Vector3{X: wx, Y: wy, Z: wz}
	v := rl.GetWorldToScreen(pos, o.cam)
	arr, err := heap.NewArray([]int64{2})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, float64(v.X))
	_ = arr.Set([]int64{1}, float64(v.Y))
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) camIsOnScreen(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.ISONSCREEN: heap not bound")
	}
	if len(args) != 4 && len(args) != 5 {
		return value.Nil, fmt.Errorf("CAMERA.ISONSCREEN expects (camera, wx, wy, wz) or (camera, wx, wy, wz, margin#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.ISONSCREEN: invalid camera handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	wx, ok1 := argF(args[1])
	wy, ok2 := argF(args[2])
	wz, ok3 := argF(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("CAMERA.ISONSCREEN: world position must be numeric")
	}
	var margin float32
	if len(args) == 5 {
		mg, ok := argF(args[4])
		if !ok {
			return value.Nil, fmt.Errorf("CAMERA.ISONSCREEN: margin must be numeric")
		}
		margin = mg
	}
	pos := rl.Vector3{X: wx, Y: wy, Z: wz}
	v := rl.GetWorldToScreen(pos, o.cam)
	rw := float32(rl.GetRenderWidth())
	rh := float32(rl.GetRenderHeight())
	inside := v.X >= -margin && v.X <= rw+margin && v.Y >= -margin && v.Y <= rh+margin
	return value.FromBool(inside), nil
}

func (m *Module) camMouseRay(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.MOUSERAY: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CAMERA.MOUSERAY expects (camera)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.MOUSERAY: invalid camera handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	mp := rl.GetMousePosition()
	ray := rl.GetScreenToWorldRayEx(mp, o.cam, int32(rl.GetRenderWidth()), int32(rl.GetRenderHeight()))
	return m.allocRayHandle(ray)
}
