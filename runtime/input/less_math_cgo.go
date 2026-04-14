//go:build cgo || (windows && !cgo)

package input

import (
	"fmt"
	"math"

	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// less math input implementation below

func (m *Module) inputMouseDelta(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("INPUT.MOUSEDELTA expects 0 arguments")
	}
	if m.h == nil {
		return value.Nil, fmt.Errorf("INPUT.MOUSEDELTA: heap not bound")
	}
	d := rl.GetMouseDelta()
	return allocInputTuple2(m.h, d.X, d.Y)
}

func (m *Module) inputMoveDir(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("INPUT.MOVEDIR expects (yaw, speed)")
	}
	if m.h == nil {
		return value.Nil, fmt.Errorf("INPUT.MOVEDIR: heap not bound")
	}
	yaw, _ := args[0].ToFloat()
	speed, _ := args[1].ToFloat()
	f := 0.0
	if rl.IsKeyDown(rl.KeyW) {
		f += 1
	}
	if rl.IsKeyDown(rl.KeyS) {
		f -= 1
	}
	s := 0.0
	if rl.IsKeyDown(rl.KeyD) {
		s += 1
	}
	if rl.IsKeyDown(rl.KeyA) {
		s -= 1
	}
	if f == 0 && s == 0 {
		return allocInputTuple2(m.h, 0, 0)
	}
	mag := math.Hypot(f, s)
	f /= mag
	s /= mag
	sy, cy := math.Sin(yaw), math.Cos(yaw)
	// Match terrain_chase: forward (sin, cos), right (cos, -sin) on XZ.
	stepX := (sy*f + cy*s) * speed
	stepZ := (cy*f - sy*s) * speed
	return allocInputTuple2(m.h, float32(stepX), float32(stepZ))
}

func allocInputTuple2(h *heap.Store, x, y float32) (value.Value, error) {
	arr, err := heap.NewArrayOfKind([]int64{2}, heap.ArrayKindFloat, 0)
	if err != nil {
		return value.Nil, err
	}
	arr.Floats[0] = float64(x)
	arr.Floats[1] = float64(y)
	id, err := h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}
