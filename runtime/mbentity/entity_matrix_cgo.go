//go:build cgo || (windows && !cgo)

package mbentity

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// localMatrix builds TRS from the entity's local position, euler angles, and scale.
// Order matches worldRotQuat: local quaternion from pitch/yaw/roll, then scale.
func (m *Module) localMatrix(e *ent) rl.Matrix {
	p, w, r := e.getRot()
	pos := e.getPos()
	q := rl.QuaternionFromEuler(p, w, r)
	rotM := rl.QuaternionToMatrix(q)
	sc := rl.MatrixScale(e.scale.X, e.scale.Y, e.scale.Z)
	tr := rl.MatrixTranslate(pos.X, pos.Y, pos.Z)
	return rl.MatrixMultiply(tr, rl.MatrixMultiply(rotM, sc))
}

// worldMatrix returns parentWorld * local for a hierarchy chain (column-major / raylib order).
// Bone sockets (FindBone): boneWorld is updated each frame from skeletal pose × host world; local is a user offset.
func (m *Module) worldMatrix(e *ent) rl.Matrix {
	L := m.localMatrix(e)
	ext := e.ext
	if ext != nil && ext.boneWorldValid {
		return rl.MatrixMultiply(ext.boneWorld, L)
	}
	if e.parentID() == 0 {
		return L
	}
	p := m.store().ents[e.parentID()]
	if p == nil {
		return L
	}
	return rl.MatrixMultiply(m.worldMatrix(p), L)
}
