//go:build cgo || (windows && !cgo)

package mbmatrix

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerQuatExtras(reg runtime.Registrar) {
	reg.Register("QUAT.TOEULER", "quat", runtime.AdaptLegacy(m.quatToEuler))
	reg.Register("QUAT.FROMVEC3TOVEC3", "quat", runtime.AdaptLegacy(m.quatFromVec3ToVec3))
	reg.Register("QUAT.FROMMAT4", "quat", runtime.AdaptLegacy(m.quatFromMat4))
	reg.Register("QUAT.TRANSFORM", "quat", runtime.AdaptLegacy(m.quatTransform))
}

// QUAT.TOEULER returns a new Vec3 handle: X=roll, Y=pitch, Z=yaw (radians), Raylib QuaternionToEuler.
func (m *Module) quatToEuler(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("QUAT.TOEULER expects quaternion handle")
	}
	q, err := m.quatFromArgs(args, 0, "QUAT.TOEULER")
	if err != nil {
		return value.Nil, err
	}
	v := rl.QuaternionToEuler(toQ(q))
	return m.allocVec3(fromV3(v))
}

func (m *Module) quatFromVec3ToVec3(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("QUAT.FROMVEC3TOVEC3 expects (fromVec3, toVec3)")
	}
	from, err := m.vec3FromArgs(args, 0, "QUAT.FROMVEC3TOVEC3")
	if err != nil {
		return value.Nil, err
	}
	to, err := m.vec3FromArgs(args, 1, "QUAT.FROMVEC3TOVEC3")
	if err != nil {
		return value.Nil, err
	}
	return m.allocQuat(fromQ(rl.QuaternionFromVector3ToVector3(toV3(from), toV3(to))))
}

func (m *Module) quatFromMat4(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("QUAT.FROMMAT4 expects matrix handle")
	}
	mat, err := m.matrixFromArgs(args, 0, "QUAT.FROMMAT4")
	if err != nil {
		return value.Nil, err
	}
	return m.allocQuat(fromQ(rl.QuaternionFromMatrix(toM(mat))))
}

// QUAT.TRANSFORM applies a 4×4 matrix to a quaternion (rotation composition helper).
func (m *Module) quatTransform(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("QUAT.TRANSFORM expects (quat, mat4)")
	}
	q, err := m.quatFromArgs(args, 0, "QUAT.TRANSFORM")
	if err != nil {
		return value.Nil, err
	}
	mat, err := m.matrixFromArgs(args, 1, "QUAT.TRANSFORM")
	if err != nil {
		return value.Nil, err
	}
	return m.allocQuat(fromQ(rl.QuaternionTransform(toQ(q), toM(mat))))
}
