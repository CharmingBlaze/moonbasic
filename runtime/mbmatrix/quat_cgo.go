//go:build cgo || (windows && !cgo)

package mbmatrix

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerQuat(reg runtime.Registrar) {
	reg.Register("QUAT.IDENTITY", "quat", runtime.AdaptLegacy(m.quatIdentity))
	reg.Register("QUAT.FROMEULER", "quat", runtime.AdaptLegacy(m.quatFromEuler))
	reg.Register("QUAT.FROMAXISANGLE", "quat", runtime.AdaptLegacy(m.quatFromAxisAngle))
	reg.Register("QUAT.MULTIPLY", "quat", runtime.AdaptLegacy(m.quatMultiply))
	reg.Register("QUAT.SLERP", "quat", runtime.AdaptLegacy(m.quatSlerp))
	reg.Register("QUAT.TOMAT4", "quat", runtime.AdaptLegacy(m.quatToMat4))
	reg.Register("QUAT.NORMALIZE", "quat", runtime.AdaptLegacy(m.quatNormalize))
	reg.Register("QUAT.INVERT", "quat", runtime.AdaptLegacy(m.quatInvert))
	reg.Register("QUAT.FREE", "quat", runtime.AdaptLegacy(m.quatFree))
}

func (m *Module) quatIdentity(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("QUAT.IDENTITY expects 0 arguments")
	}
	return m.allocQuat(fromQ(rl.QuaternionIdentity()))
}

func (m *Module) quatFromEuler(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("QUAT.FROMEULER expects (pitch#, yaw#, roll#) radians")
	}
	px, ok1 := argF(args[0])
	py, ok2 := argF(args[1])
	pz, ok3 := argF(args[2])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("QUAT.FROMEULER: angles must be numeric")
	}
	return m.allocQuat(fromQ(rl.QuaternionFromEuler(px, py, pz)))
}

func (m *Module) quatFromAxisAngle(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("QUAT.FROMAXISANGLE expects (ax#, ay#, az#, angle#)")
	}
	ax, ok1 := argF(args[0])
	ay, ok2 := argF(args[1])
	az, ok3 := argF(args[2])
	ang, ok4 := argF(args[3])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("QUAT.FROMAXISANGLE: arguments must be numeric")
	}
	axis := rl.Vector3{X: ax, Y: ay, Z: az}
	return m.allocQuat(fromQ(rl.QuaternionFromAxisAngle(axis, ang)))
}

func (m *Module) quatMultiply(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("QUAT.MULTIPLY expects two quaternion handles")
	}
	a, err := m.quatFromArgs(args, 0, "QUAT.MULTIPLY")
	if err != nil {
		return value.Nil, err
	}
	b, err := m.quatFromArgs(args, 1, "QUAT.MULTIPLY")
	if err != nil {
		return value.Nil, err
	}
	return m.allocQuat(fromQ(rl.QuaternionMultiply(toQ(a), toQ(b))))
}

func (m *Module) quatSlerp(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("QUAT.SLERP expects (q1, q2, t#)")
	}
	a, err := m.quatFromArgs(args, 0, "QUAT.SLERP")
	if err != nil {
		return value.Nil, err
	}
	b, err := m.quatFromArgs(args, 1, "QUAT.SLERP")
	if err != nil {
		return value.Nil, err
	}
	t, ok := argF(args[2])
	if !ok {
		return value.Nil, fmt.Errorf("QUAT.SLERP: t must be numeric")
	}
	return m.allocQuat(fromQ(rl.QuaternionSlerp(toQ(a), toQ(b), t)))
}

func (m *Module) quatToMat4(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("QUAT.TOMAT4 expects quaternion handle")
	}
	q, err := m.quatFromArgs(args, 0, "QUAT.TOMAT4")
	if err != nil {
		return value.Nil, err
	}
	id, err := m.allocMat(fromM(rl.QuaternionToMatrix(toQ(q))))
	if err != nil {
		return value.Nil, err
	}
	return id, nil
}

func (m *Module) quatNormalize(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("QUAT.NORMALIZE expects quaternion handle")
	}
	q, err := m.quatFromArgs(args, 0, "QUAT.NORMALIZE")
	if err != nil {
		return value.Nil, err
	}
	return m.allocQuat(fromQ(rl.QuaternionNormalize(toQ(q))))
}

func (m *Module) quatInvert(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("QUAT.INVERT expects quaternion handle")
	}
	q, err := m.quatFromArgs(args, 0, "QUAT.INVERT")
	if err != nil {
		return value.Nil, err
	}
	return m.allocQuat(fromQ(rl.QuaternionInvert(toQ(q))))
}

func (m *Module) quatFree(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("QUAT.FREE expects quaternion handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}
