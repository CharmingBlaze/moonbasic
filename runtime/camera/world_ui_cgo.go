//go:build cgo || (windows && !cgo)

package mbcamera

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func RegisterCameraQoLAPI(m *Module, r runtime.Registrar) {
	r.Register("WORLD.TOSCREEN", "camera", runtime.AdaptLegacy(m.worldToScreen))
	r.Register("WORLD.TOWORLD", "camera", runtime.AdaptLegacy(m.worldToWorld))
	r.Register("WORLD.GETRAY", "camera", runtime.AdaptLegacy(m.worldGetRay))
	r.Register("WORLD.MOUSEFLOOR", "camera", runtime.AdaptLegacy(m.worldMouseFloor))
	r.Register("WORLD.MOUSEPICK", "camera", runtime.AdaptLegacy(m.worldMousePick))
	r.Register("WORLD.SHAKE", "camera", runtime.AdaptLegacy(m.cameraShake))
	r.Register("WORLD.HITSTOP", "camera", m.worldHitStop)
	r.Register("CAMERA.SHAKE", "camera", runtime.AdaptLegacy(m.cameraShake))
}

func (m *Module) worldToScreen(args []value.Value) (value.Value, error) {
	if len(args) != 2 && len(args) != 4 {
		return value.Nil, fmt.Errorf("WORLD.TOSCREEN expects (x, y, z, camera) or (entity, camera)")
	}
	
	var pos rl.Vector3
	var camID heap.Handle
	
	if len(args) == 2 {
		// Shorthand: TOSCREEN(entity, camera)
		id, _ := args[0].ToInt()
		camID = heap.Handle(args[1].IVal)
		if m.entityWorldPos != nil {
			if px, py, pz, ok := m.entityWorldPos(m.h, heap.Handle(args[0].IVal)); ok {
				pos = rl.Vector3{X: px, Y: py, Z: pz}
			} else {
				// Try numeric lookup if handle fails
				if reg := runtime.ActiveRegistry(); reg != nil && reg.ResolveEntityWorldPos != nil {
					if wp, ok2 := reg.ResolveEntityWorldPos(id); ok2 {
						pos = wp
					}
				}
			}
		}
	} else {
		x, _ := args[0].ToFloat()
		y, _ := args[1].ToFloat()
		z, _ := args[2].ToFloat()
		pos = rl.Vector3{X: float32(x), Y: float32(y), Z: float32(z)}
		camID = heap.Handle(args[3].IVal)
	}
	
	cam, err := heap.Cast[*camObj](m.h, camID)
	if err != nil { return value.Nil, fmt.Errorf("invalid camera") }

	pos2d := rl.GetWorldToScreen(pos, cam.cam)
	
	arr, _ := heap.NewArrayOfKind([]int64{2}, heap.ArrayKindFloat, 0)
	arr.Floats[0], arr.Floats[1] = float64(pos2d.X), float64(pos2d.Y)
	h, _ := m.h.Alloc(arr)

	return value.FromHandle(h), nil
}

func (m *Module) worldToWorld(args []value.Value) (value.Value, error) {
	return value.Nil, nil // Stubbed mapping to rl.GetScreenToWorldRay
}

func (m *Module) worldGetRay(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("WORLD.GETRAY expects (mouseX, mouseY, cameraHandle)")
	}
	mx, _ := args[0].ToFloat()
	my, _ := args[1].ToFloat()
	id := heap.Handle(args[2].IVal)
	
	cam, err := heap.Cast[*camObj](m.h, id)
	if err != nil { return value.Nil, fmt.Errorf("invalid camera") }

	ray := rl.GetScreenToWorldRay(rl.Vector2{X: float32(mx), Y: float32(my)}, cam.cam)
	
	arr, _ := heap.NewArrayOfKind([]int64{6}, heap.ArrayKindFloat, 0)
	arr.Floats[0], arr.Floats[1], arr.Floats[2] = float64(ray.Position.X), float64(ray.Position.Y), float64(ray.Position.Z)
	arr.Floats[3], arr.Floats[4], arr.Floats[5] = float64(ray.Direction.X), float64(ray.Direction.Y), float64(ray.Direction.Z)
	h, _ := m.h.Alloc(arr)
	
	return value.FromHandle(h), nil
}

func (m *Module) cameraShake(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CAMERA.SHAKE expects (cameraHandle, intensity, duration)")
	}
	id := heap.Handle(args[0].IVal)
	cam, err := heap.Cast[*camObj](m.h, id)
	if err != nil { return value.Nil, fmt.Errorf("invalid camera") }
	
	intensity, _ := args[1].ToFloat()
	dur, _ := args[2].ToFloat()
	
	cam.shakeMag = float32(intensity)
	cam.shakeTime = float32(dur)

	return value.Nil, nil
}

func (m *Module) worldMouseFloor(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("WORLD.MOUSEFLOOR expects (camera, floorY#)")
	}
	cid := heap.Handle(args[0].IVal)
	fy, _ := args[1].ToFloat()
	
	cam, err := heap.Cast[*camObj](m.h, cid)
	if err != nil { return value.Nil, fmt.Errorf("invalid camera") }
	
	ray := rl.GetScreenToWorldRay(rl.GetMousePosition(), cam.cam)
	
	// Ray-Plane intersection: P = O + t*D
	// (P.y - floorY) = 0  =>  O.y + t*D.y - floorY = 0  => t = (floorY - O.y) / D.y
	if math.Abs(float64(ray.Direction.Y)) < 1e-6 {
		return value.Nil, nil // Parallel or pointing away
	}
	
	t := (float32(fy) - ray.Position.Y) / ray.Direction.Y
	if t < 0 { return value.Nil, nil }
	
	hit := rl.Vector3Add(ray.Position, rl.Vector3Scale(ray.Direction, t))
	
	arr, _ := heap.NewArrayOfKind([]int64{3}, heap.ArrayKindFloat, 0)
	arr.Floats[0], arr.Floats[1], arr.Floats[2] = float64(hit.X), float64(hit.Y), float64(hit.Z)
	h, _ := m.h.Alloc(arr)
	
	return value.FromHandle(h), nil
}

func (m *Module) worldMousePick(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WORLD.MOUSEPICK expects (camera)")
	}
	cid := heap.Handle(args[0].IVal)
	cam, err := heap.Cast[*camObj](m.h, cid)
	if err != nil { return value.Nil, fmt.Errorf("invalid camera") }
	
	ray := rl.GetScreenToWorldRay(rl.GetMousePosition(), cam.cam)
	
	// Delegate to physics engine raycast
	reg := runtime.ActiveRegistry()
	if reg == nil { return value.FromInt(0), nil }
	
	// Call ENTITY.RAYCAST internally or use physics helper
	id, err := reg.Call("ENTITY.RAYCAST", []value.Value{
		value.FromFloat(float64(ray.Position.X)), value.FromFloat(float64(ray.Position.Y)), value.FromFloat(float64(ray.Position.Z)),
		value.FromFloat(float64(ray.Direction.X)), value.FromFloat(float64(ray.Direction.Y)), value.FromFloat(float64(ray.Direction.Z)),
		value.FromFloat(1000.0),
	})
	return id, err
}

func (m *Module) worldHitStop(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WORLD.HITSTOP expects (durationSeconds#)")
	}
	dur, _ := args[0].ToFloat()
	if dur <= 0 { return value.Nil, nil }
	
	rt.HitStopEndAt = rl.GetTime() + dur
	return value.Nil, nil
}
