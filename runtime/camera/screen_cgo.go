//go:build cgo || (windows && !cgo)

package mbcamera

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerScreenHelpers(reg runtime.Registrar) {
	reg.Register("CAMERA.WORLDTOSCREEN", "camera", runtime.AdaptLegacy(m.camWorldToScreen))
	reg.Register("CAMERA.PROJECT", "camera", runtime.AdaptLegacy(m.camWorldToScreen))
	reg.Register("CAMERA.ISONSCREEN", "camera", runtime.AdaptLegacy(m.camIsOnScreen))
	reg.Register("CAMERA.MOUSERAY", "camera", runtime.AdaptLegacy(m.camMouseRay))
	reg.Register("WORLD.TOSCREEN", "world", runtime.AdaptLegacy(m.worldToScreenActive))
	reg.Register("WORLD.TOWORLD", "world", runtime.AdaptLegacy(m.worldFromScreenActive))
	reg.Register("WORLD.MOUSE2D", "world", runtime.AdaptLegacy(m.worldMouse2D))
	reg.Register("WORLD.MOUSEFLOOR3D", "world", runtime.AdaptLegacy(m.worldMouseFloor3D))
	reg.Register("WORLD.MOUSETOFLOOR", "world", runtime.AdaptLegacy(m.worldMouseFloor3D))
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

// worldToScreenActive uses the last CAMERA.BEGIN 3D camera (same as CAMERA.GETACTIVE).
func (m *Module) worldToScreenActive(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("WORLD.TOSCREEN: heap not bound")
	}
	if m.lastActive3D == 0 {
		return value.Nil, fmt.Errorf("WORLD.TOSCREEN: no active 3D camera (call CAMERA.BEGIN first)")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("WORLD.TOSCREEN expects wx#, wy#, wz#")
	}
	return m.camWorldToScreen([]value.Value{value.FromHandle(int32(m.lastActive3D)), args[0], args[1], args[2]})
}

// worldFromScreenActive unprojects screen pixels through the active 3D camera; depth is distance along the view ray.
func (m *Module) worldFromScreenActive(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("WORLD.TOWORLD: heap not bound")
	}
	if m.lastActive3D == 0 {
		return value.Nil, fmt.Errorf("WORLD.TOWORLD: no active 3D camera (call CAMERA.BEGIN first)")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("WORLD.TOWORLD expects screenX#, screenY#, depth#")
	}
	o, err := heap.Cast[*camObj](m.h, m.lastActive3D)
	if err != nil {
		return value.Nil, err
	}
	sx, ok1 := argF(args[0])
	sy, ok2 := argF(args[1])
	d, ok3 := argF(args[2])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("WORLD.TOWORLD: arguments must be numeric")
	}
	ray := rl.GetScreenToWorldRayEx(rl.Vector2{X: sx, Y: sy}, o.cam, int32(rl.GetRenderWidth()), int32(rl.GetRenderHeight()))
	dir := rl.Vector3Normalize(ray.Direction)
	p := rl.Vector3Add(ray.Position, rl.Vector3Scale(dir, d))
	arr, err := heap.NewArray([]int64{3})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, float64(p.X))
	_ = arr.Set([]int64{1}, float64(p.Y))
	_ = arr.Set([]int64{2}, float64(p.Z))
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

// worldMouse2D transforms the current mouse position through a Camera2D (same idea as CAMERA2D.SCREENTOWORLD at the mouse).
func (m *Module) worldMouse2D(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("WORLD.MOUSE2D: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WORLD.MOUSE2D expects (camera2D)")
	}
	o, err := m.getCam2D(args, 0, "WORLD.MOUSE2D")
	if err != nil {
		return value.Nil, err
	}
	mp := rl.GetMousePosition()
	w := rl.GetScreenToWorld2D(mp, o.cam)
	arr, err := heap.NewArray([]int64{2})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, float64(w.X))
	_ = arr.Set([]int64{1}, float64(w.Y))
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

// worldMouseFloor3D intersects the mouse ray from a Camera3D with the horizontal plane y=floorY.
func (m *Module) worldMouseFloor3D(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("WORLD.MOUSEFLOOR3D: heap not bound")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("WORLD.MOUSEFLOOR3D expects (camera3D, floorY#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("WORLD.MOUSEFLOOR3D: invalid camera handle")
	}
	co, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	floorY, ok1 := argF(args[1])
	if !ok1 {
		return value.Nil, fmt.Errorf("WORLD.MOUSEFLOOR3D: floorY must be numeric")
	}
	mp := rl.GetMousePosition()
	ray := rl.GetScreenToWorldRayEx(mp, co.cam, int32(rl.GetRenderWidth()), int32(rl.GetRenderHeight()))
	dy := float64(ray.Direction.Y)
	if math.Abs(dy) < 1e-8 {
		return value.Nil, nil
	}
	t := (float64(floorY) - float64(ray.Position.Y)) / dy
	if t < 0 {
		return value.Nil, nil
	}
	x := float64(ray.Position.X) + t*float64(ray.Direction.X)
	z := float64(ray.Position.Z) + t*float64(ray.Direction.Z)
	arr, err := heap.NewArray([]int64{2})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, x)
	_ = arr.Set([]int64{1}, z)
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}
