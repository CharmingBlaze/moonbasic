//go:build cgo || (windows && !cgo)

package mbcamera

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func RegisterCameraQoLAPI(m *Module, r runtime.Registrar) {
	r.Register("WORLD.TOSCREEN", "camera", runtime.AdaptLegacy(m.worldToScreen))
	r.Register("WORLD.TOWORLD", "camera", runtime.AdaptLegacy(m.worldToWorld))
	r.Register("WORLD.GETRAY", "camera", runtime.AdaptLegacy(m.worldGetRay))
	r.Register("CAMERA.SHAKE", "camera", runtime.AdaptLegacy(m.cameraShake))
}

func (m *Module) worldToScreen(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("WORLD.TOSCREEN expects (x, y, z, cameraHandle)")
	}
	x, _ := args[0].ToFloat()
	y, _ := args[1].ToFloat()
	z, _ := args[2].ToFloat()
	
	id := heap.Handle(args[3].IVal)
	cam, err := heap.Cast[*camObj](m.h, id)
	if err != nil { return value.Nil, fmt.Errorf("invalid camera") }

	pos2d := rl.GetWorldToScreen(rl.Vector3{X: float32(x), Y: float32(y), Z: float32(z)}, cam.cam)
	
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
	// Assigning direct local offsets into the cam structs...
	cam.cam.Position.Y += float32(intensity)

	return value.Nil, nil
}
