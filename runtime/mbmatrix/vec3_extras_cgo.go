//go:build cgo || (windows && !cgo)

package mbmatrix

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerVec3Extras(reg runtime.Registrar) {
	reg.Register("VEC3.TRANSFORMMAT4", "vec3", runtime.AdaptLegacy(m.vec3TransformMat4))
	reg.Register("VEC3.ANGLE", "vec3", runtime.AdaptLegacy(m.vec3Angle))
	reg.Register("VEC3.PROJECT", "vec3", runtime.AdaptLegacy(m.vec3Project))
	reg.Register("VEC3.ORTHONORMALIZE", "vec3", runtime.AdaptLegacy(m.vec3OrthoNormalize))
	reg.Register("VEC3.ROTATEBYQUAT", "vec3", runtime.AdaptLegacy(m.vec3RotateByQuat))
}

func (m *Module) vec3TransformMat4(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC3.TRANSFORMMAT4 expects (vec3, mat4)")
	}
	v, err := m.vec3FromArgs(args, 0, "VEC3.TRANSFORMMAT4")
	if err != nil {
		return value.Nil, err
	}
	mat, err := m.matrixFromArgs(args, 1, "VEC3.TRANSFORMMAT4")
	if err != nil {
		return value.Nil, err
	}
	return m.allocVec3(fromV3(rl.Vector3Transform(toV3(v), toM(mat))))
}

func (m *Module) vec3Angle(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC3.ANGLE expects two vec3 handles")
	}
	a, err := m.vec3FromArgs(args, 0, "VEC3.ANGLE")
	if err != nil {
		return value.Nil, err
	}
	b, err := m.vec3FromArgs(args, 1, "VEC3.ANGLE")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(rl.Vector3Angle(toV3(a), toV3(b)))), nil
}

func (m *Module) vec3Project(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC3.PROJECT expects (vec3, ontoVec3)")
	}
	v1, err := m.vec3FromArgs(args, 0, "VEC3.PROJECT")
	if err != nil {
		return value.Nil, err
	}
	v2, err := m.vec3FromArgs(args, 1, "VEC3.PROJECT")
	if err != nil {
		return value.Nil, err
	}
	return m.allocVec3(fromV3(rl.Vector3Project(toV3(v1), toV3(v2))))
}

func (m *Module) vec3OrthoNormalize(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("VEC3.ORTHONORMALIZE expects two vec3 handles (modified in place)")
	}
	o1, err := heap.Cast[*vec3Obj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("VEC3.ORTHONORMALIZE: %w", err)
	}
	o2, err := heap.Cast[*vec3Obj](m.h, heap.Handle(args[1].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("VEC3.ORTHONORMALIZE: %w", err)
	}
	rv1, rv2 := toV3(o1.v), toV3(o2.v)
	rl.Vector3OrthoNormalize(&rv1, &rv2)
	o1.v, o2.v = fromV3(rv1), fromV3(rv2)
	return value.Nil, nil
}

func (m *Module) vec3RotateByQuat(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC3.ROTATEBYQUAT expects (vec3, quat)")
	}
	v, err := m.vec3FromArgs(args, 0, "VEC3.ROTATEBYQUAT")
	if err != nil {
		return value.Nil, err
	}
	q, err := m.quatFromArgs(args, 1, "VEC3.ROTATEBYQUAT")
	if err != nil {
		return value.Nil, err
	}
	return m.allocVec3(fromV3(rl.Vector3RotateByQuaternion(toV3(v), toQ(q))))
}
