//go:build cgo || (windows && !cgo)

package mbmatrix

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/hal"
	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// AllocVec3Value allocates a heap-backed vec3 (VEC3.*) for other runtime packages (e.g. steering).
func AllocVec3Value(h *heap.Store, x, y, z float32) (value.Value, error) {
	if h == nil {
		return value.Nil, fmt.Errorf("AllocVec3Value: heap is nil")
	}
	id, err := h.Alloc(&vec3Obj{v: hal.V3{X: x, Y: y, Z: z}})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

// Conversion helpers (hal <-> rl)
func toV2(v hal.V2) rl.Vector2 { return rl.Vector2{X: v.X, Y: v.Y} }
func toV3(v hal.V3) rl.Vector3 { return rl.Vector3{X: v.X, Y: v.Y, Z: v.Z} }
func toQ(v hal.V4) rl.Quaternion { return rl.Quaternion{X: v.X, Y: v.Y, Z: v.Z, W: v.W} }
func toM(m hal.Matrix) rl.Matrix {
	return rl.Matrix{
		M0: m.M0, M4: m.M4, M8: m.M8, M12: m.M12,
		M1: m.M1, M5: m.M5, M9: m.M9, M13: m.M13,
		M2: m.M2, M6: m.M6, M10: m.M10, M14: m.M14,
		M3: m.M3, M7: m.M7, M11: m.M11, M15: m.M15,
	}
}

func fromV2(v rl.Vector2) hal.V2 { return hal.V2{X: v.X, Y: v.Y} }
func fromV3(v rl.Vector3) hal.V3 { return hal.V3{X: v.X, Y: v.Y, Z: v.Z} }
func fromQ(v rl.Quaternion) hal.V4 { return hal.V4{X: v.X, Y: v.Y, Z: v.Z, W: v.W} }
func fromM(m rl.Matrix) hal.Matrix {
	return hal.Matrix{
		M0: m.M0, M4: m.M4, M8: m.M8, M12: m.M12,
		M1: m.M1, M5: m.M5, M9: m.M9, M13: m.M13,
		M2: m.M2, M6: m.M6, M10: m.M10, M14: m.M14,
		M3: m.M3, M7: m.M7, M11: m.M11, M15: m.M15,
	}
}

func (m *Module) requireHeap() error {
	if m.h == nil {
		return runtime.Errorf("matrix/vector builtins: heap not bound")
	}
	return nil
}

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

func (m *Module) matrixFromArgs(args []value.Value, ix int, op string) (hal.Matrix, error) {
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return hal.Matrix{}, fmt.Errorf("%s: argument %d must be matrix handle", op, ix+1)
	}
	o, err := heap.Cast[*matObj](m.h, heap.Handle(args[ix].IVal))
	if err != nil {
		return hal.Matrix{}, err
	}
	return o.m, nil
}

func (m *Module) vec3FromArgs(args []value.Value, idx int, op string) (hal.V3, error) {
	if idx >= len(args) || args[idx].Kind != value.KindHandle {
		return hal.V3{}, fmt.Errorf("%s: argument %d must be vec3 handle", op, idx+1)
	}
	o, err := heap.Cast[*vec3Obj](m.h, heap.Handle(args[idx].IVal))
	if err != nil {
		return hal.V3{}, fmt.Errorf("%s: %w", op, err)
	}
	return o.v, nil
}

func (m *Module) quatFromArgs(args []value.Value, idx int, op string) (hal.V4, error) {
	if idx >= len(args) || args[idx].Kind != value.KindHandle {
		return hal.V4{}, fmt.Errorf("%s: argument %d must be quaternion handle", op, idx+1)
	}
	o, err := heap.Cast[*quatObj](m.h, heap.Handle(args[idx].IVal))
	if err != nil {
		return hal.V4{}, fmt.Errorf("%s: %w", op, err)
	}
	return o.q, nil
}

func (m *Module) vec2FromArgs(args []value.Value, idx int, op string) (hal.V2, error) {
	if idx >= len(args) || args[idx].Kind != value.KindHandle {
		return hal.V2{}, fmt.Errorf("%s: argument %d must be vec2 handle", op, idx+1)
	}
	o, err := heap.Cast[*vec2Obj](m.h, heap.Handle(args[idx].IVal))
	if err != nil {
		return hal.V2{}, fmt.Errorf("%s: %w", op, err)
	}
	return o.v, nil
}

func (m *Module) allocVec3(v hal.V3) (value.Value, error) {
	id, err := m.h.Alloc(&vec3Obj{v: v})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) allocVec2(v hal.V2) (value.Value, error) {
	id, err := m.h.Alloc(&vec2Obj{v: v})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) allocMat(mat hal.Matrix) (value.Value, error) {
	id, err := m.h.Alloc(&matObj{m: mat})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) allocQuat(q hal.V4) (value.Value, error) {
	id, err := m.h.Alloc(&quatObj{q: q})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func matElement(mat hal.Matrix, row, col int32) float32 {
	switch col {
	case 0:
		switch row {
		case 0: return mat.M0
		case 1: return mat.M1
		case 2: return mat.M2
		case 3: return mat.M3
		}
	case 1:
		switch row {
		case 0: return mat.M4
		case 1: return mat.M5
		case 2: return mat.M6
		case 3: return mat.M7
		}
	case 2:
		switch row {
		case 0: return mat.M8
		case 1: return mat.M9
		case 2: return mat.M10
		case 3: return mat.M11
		}
	case 3:
		switch row {
		case 0: return mat.M12
		case 1: return mat.M13
		case 2: return mat.M14
		case 3: return mat.M15
		}
	}
	return 0
}
