//go:build cgo

package mbmodel3d

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/convert"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) requireHeap() error {
	if m.h == nil {
		return runtime.Errorf("3D builtins: heap not bound")
	}
	return nil
}

func argInt(v value.Value) (int32, bool) {
	if i, ok := v.ToInt(); ok {
		return int32(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return int32(f), true
	}
	return 0, false
}

func argFloat(v value.Value) (float32, bool) {
	if f, ok := v.ToFloat(); ok {
		return float32(f), true
	}
	if i, ok := v.ToInt(); ok {
		return float32(i), true
	}
	return 0, false
}

func argBool(v value.Value) (bool, bool) {
	if v.Kind == value.KindBool {
		return v.IVal != 0, true
	}
	if i, ok := v.ToInt(); ok {
		return i != 0, true
	}
	return false, false
}

func rgbaFromArgs(r, g, b, a value.Value) (color.RGBA, error) {
	ri, ok1 := argInt(r)
	gi, ok2 := argInt(g)
	bi, ok3 := argInt(b)
	ai, ok4 := argInt(a)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return color.RGBA{}, fmt.Errorf("expected numeric RGBA components")
	}
	c := convert.NewColor4(ri, gi, bi, ai)
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}, nil
}

func (m *Module) getMesh(args []value.Value, ix int, op string) (*meshObj, error) {
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: argument %d must be mesh handle", op, ix+1)
	}
	return heap.Cast[*meshObj](m.h, heap.Handle(args[ix].IVal))
}

func (m *Module) getMaterial(args []value.Value, ix int, op string) (*materialObj, error) {
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: argument %d must be material handle", op, ix+1)
	}
	return heap.Cast[*materialObj](m.h, heap.Handle(args[ix].IVal))
}

func (m *Module) getModel(args []value.Value, ix int, op string) (*modelObj, error) {
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: argument %d must be model handle", op, ix+1)
	}
	return heap.Cast[*modelObj](m.h, heap.Handle(args[ix].IVal))
}

func (m *Module) getShader(args []value.Value, ix int, op string) (*shaderObj, error) {
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: argument %d must be shader handle", op, ix+1)
	}
	return heap.Cast[*shaderObj](m.h, heap.Handle(args[ix].IVal))
}

func (m *Module) allocMesh(mesh rl.Mesh, op string) (value.Value, error) {
	id, err := m.h.Alloc(&meshObj{m: mesh})
	if err != nil {
		return value.Nil, fmt.Errorf("%s: %w", op, err)
	}
	o, err := heap.Cast[*meshObj](m.h, id)
	if err != nil {
		return value.Nil, fmt.Errorf("%s: %w", op, err)
	}
	rl.UploadMesh(&o.m, false)
	return value.FromHandle(id), nil
}

// ModelFromArgs retrieves a Raylib model from a heap handle.
func ModelFromArgs(s *heap.Store, h heap.Handle) (rl.Model, error) {
	o, err := heap.Cast[*modelObj](s, h)
	if err != nil {
		return rl.Model{}, err
	}
	return o.model, nil
}
