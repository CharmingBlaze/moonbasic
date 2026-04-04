//go:build cgo

package mbmatrix

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// AllocVec3Value allocates a heap-backed vec3 (VEC3.*) for other runtime packages (e.g. steering).
func AllocVec3Value(h *heap.Store, x, y, z float32) (value.Value, error) {
	if h == nil {
		return value.Nil, fmt.Errorf("AllocVec3Value: heap is nil")
	}
	id, err := h.Alloc(&vec3Obj{v: rl.Vector3{X: x, Y: y, Z: z}})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
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

func (m *Module) matrixFromArgs(args []value.Value, ix int, op string) (rl.Matrix, error) {
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return rl.Matrix{}, fmt.Errorf("%s: argument %d must be matrix handle", op, ix+1)
	}
	o, err := heap.Cast[*matObj](m.h, heap.Handle(args[ix].IVal))
	if err != nil {
		return rl.Matrix{}, err
	}
	return o.m, nil
}

func (m *Module) vec3FromArgs(args []value.Value, idx int, op string) (rl.Vector3, error) {
	if idx >= len(args) || args[idx].Kind != value.KindHandle {
		return rl.Vector3{}, fmt.Errorf("%s: argument %d must be vec3 handle", op, idx+1)
	}
	o, err := heap.Cast[*vec3Obj](m.h, heap.Handle(args[idx].IVal))
	if err != nil {
		return rl.Vector3{}, fmt.Errorf("%s: %w", op, err)
	}
	return o.v, nil
}

func (m *Module) vec2FromArgs(args []value.Value, idx int, op string) (rl.Vector2, error) {
	if idx >= len(args) || args[idx].Kind != value.KindHandle {
		return rl.Vector2{}, fmt.Errorf("%s: argument %d must be vec2 handle", op, idx+1)
	}
	o, err := heap.Cast[*vec2Obj](m.h, heap.Handle(args[idx].IVal))
	if err != nil {
		return rl.Vector2{}, fmt.Errorf("%s: %w", op, err)
	}
	return o.v, nil
}

func matElement(mat rl.Matrix, row, col int32) float32 {
	// Column-major (OpenGL / Raylib): columns are (M0,M1,M2,M3), (M4,…), …
	switch col {
	case 0:
		switch row {
		case 0:
			return mat.M0
		case 1:
			return mat.M1
		case 2:
			return mat.M2
		case 3:
			return mat.M3
		}
	case 1:
		switch row {
		case 0:
			return mat.M4
		case 1:
			return mat.M5
		case 2:
			return mat.M6
		case 3:
			return mat.M7
		}
	case 2:
		switch row {
		case 0:
			return mat.M8
		case 1:
			return mat.M9
		case 2:
			return mat.M10
		case 3:
			return mat.M11
		}
	case 3:
		switch row {
		case 0:
			return mat.M12
		case 1:
			return mat.M13
		case 2:
			return mat.M14
		case 3:
			return mat.M15
		}
	}
	return 0
}
