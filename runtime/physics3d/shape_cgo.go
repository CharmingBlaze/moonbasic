//go:build (linux || windows) && cgo
package mbphysics3d

import (
	"fmt"
	"github.com/bbitechnologies/jolt-go/jolt"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type ShapeObj struct {
	Shape *jolt.Shape
	Kind  int // 1: Box, 2: Sphere, 3: Capsule, 4: Cylinder
	Dim1  float32
	Dim2  float32
	Dim3  float32
}

func (s *ShapeObj) TypeName() string { return "PhysicsShape" }
func (s *ShapeObj) TypeTag() uint16  { return heap.TagShape }
func (s *ShapeObj) Free() {
	if s.Shape != nil {
		s.Shape.Destroy()
		s.Shape = nil
	}
}

func (m *Module) shCreateBox(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("SHAPE.CREATEBOX expects (hx, hy, hz)")
	}
	hx, _ := args[0].ToFloat()
	hy, _ := args[1].ToFloat()
	hz, _ := args[2].ToFloat()
	sh := jolt.CreateBox(jolt.Vec3{X: float32(hx), Y: float32(hy), Z: float32(hz)})
	obj := &ShapeObj{Shape: sh, Kind: 1, Dim1: float32(hx), Dim2: float32(hy), Dim3: float32(hz)}
	id, err := m.h.Alloc(obj)
	if err != nil {
		sh.Destroy()
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) shCreateSphere(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SHAPE.CREATESPHERE expects (radius)")
	}
	r, _ := args[0].ToFloat()
	sh := jolt.CreateSphere(float32(r))
	obj := &ShapeObj{Shape: sh, Kind: 2, Dim1: float32(r)}
	id, err := m.h.Alloc(obj)
	if err != nil {
		sh.Destroy()
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) shCreateCapsule(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SHAPE.CREATECAPSULE expects (radius, height)")
	}
	r, _ := args[0].ToFloat()
	h_val, _ := args[1].ToFloat()
	hh := float32(h_val)/2 - float32(r)
	if hh < 0.05 { hh = 0.05 }
	sh := jolt.CreateCapsule(hh, float32(r))
	obj := &ShapeObj{Shape: sh, Kind: 3, Dim1: float32(r), Dim2: float32(h_val)}
	id, err := m.h.Alloc(obj)
	if err != nil {
		sh.Destroy()
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) shCreateCylinder(args []value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("SHAPE.CREATECYLINDER: modern cylinder shapes not exposed in this jolt-go wrapper; use CreateBox or CreateCapsule")
}

func (m *Module) shFree(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("SHAPE.FREE expects shape handle")
	}
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func (m *Module) shGetType(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("SHAPE.GETTYPE expects (handle)")
	}
	sh, err := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(int64(sh.Kind)), nil
}

func (m *Module) shGetDim1(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("SHAPE.GETWIDTH expects (handle)")
	}
	sh, err := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(sh.Dim1)), nil
}

func (m *Module) shGetDim2(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("SHAPE.GETHEIGHT expects (handle)")
	}
	sh, err := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(sh.Dim2)), nil
}

func (m *Module) shGetDim3(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("SHAPE.GETDEPTH expects (handle)")
	}
	sh, err := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(sh.Dim3)), nil
}
