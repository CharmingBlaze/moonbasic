//go:build cgo || (windows && !cgo)

package mbentity

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (m *Module) worldRotQuat(e *ent) rl.Quaternion {
	p, w, r := e.getRot()
	ql := rl.QuaternionFromEuler(p, w, r)
	if e.parentID() == 0 {
		return ql
	}
	parent := m.store().ents[e.parentID()]
	if parent == nil {
		return ql
	}
	qp := m.worldRotQuat(parent)
	return rl.QuaternionMultiply(qp, ql)
}

func (m *Module) worldPos(e *ent) rl.Vector3 {
	// When any ancestor is a bone socket, TRS must match worldMatrix (matrix path).
	if m.entityUsesBoneMatrixChain(e) {
		wmat := m.worldMatrix(e)
		return rl.Vector3{X: wmat.M12, Y: wmat.M13, Z: wmat.M14}
	}
	if e.parentID() == 0 {
		return e.getPos()
	}
	p := m.store().ents[e.parentID()]
	if p == nil {
		return e.getPos()
	}
	pw := m.worldPos(p)
	pq := m.worldRotQuat(p)
	off := rl.Vector3RotateByQuaternion(e.getPos(), pq)
	return rl.Vector3Add(pw, off)
}

func (m *Module) entityUsesBoneMatrixChain(e *ent) bool {
	for e != nil {
		if e.boneWorldValid() {
			return true
		}
		if e.parentID() == 0 {
			break
		}
		e = m.store().ents[e.parentID()]
	}
	return false
}

func (m *Module) worldEuler(e *ent) (pitch, yaw, roll float32) {
	q := m.worldRotQuat(e)
	v := rl.QuaternionToEuler(q)
	// Raylib: X=roll, Y=pitch, Z=yaw (see QUAT.TOEULER in mbmatrix)
	return v.Y, v.Z, v.X
}

// WorldEulerForEntityID returns world-space pitch, yaw, roll in degrees (same convention as ENTITY.ENTITYPITCH/YAW/ROLL).
func (m *Module) WorldEulerForEntityID(id int64) (pitch, yaw, roll float32, ok bool) {
	e := m.store().ents[id]
	if e == nil {
		return 0, 0, 0, false
	}
	p, y, r := m.worldEuler(e)
	return p, y, r, true
}

func (m *Module) setLocalFromWorld(e *ent, wx, wy, wz float32) {
	if e.parentID() == 0 {
		e.setPos(rl.Vector3{X: wx, Y: wy, Z: wz})
		return
	}
	p := m.store().ents[e.parentID()]
	if p == nil {
		e.setPos(rl.Vector3{X: wx, Y: wy, Z: wz})
		return
	}
	pw := m.worldPos(p)
	pq := m.worldRotQuat(p)
	rel := rl.Vector3Subtract(rl.Vector3{X: wx, Y: wy, Z: wz}, pw)
	inv := rl.QuaternionInvert(pq)
	e.setPos(rl.Vector3RotateByQuaternion(rel, inv))
}

// setLocalRotFromWorldQuat applies a world-space orientation to the entity's local pitch/yaw/roll,
// accounting for parent chain (same euler convention as worldEuler / QuaternionFromEuler).
func (m *Module) setLocalRotFromWorldQuat(e *ent, qWorld rl.Quaternion) {
	if e.parentID() == 0 {
		v := rl.QuaternionToEuler(qWorld)
		e.setRot(v.Y, v.Z, v.X)
		return
	}
	parent := m.store().ents[e.parentID()]
	if parent == nil {
		v := rl.QuaternionToEuler(qWorld)
		e.setRot(v.Y, v.Z, v.X)
		return
	}
	qp := m.worldRotQuat(parent)
	inv := rl.QuaternionInvert(qp)
	qLocal := rl.QuaternionMultiply(inv, qWorld)
	v := rl.QuaternionToEuler(qLocal)
	e.setRot(v.Y, v.Z, v.X)
}

func forwardFromYawPitch(yaw, pitch float32) rl.Vector3 {
	cp := float32(math.Cos(float64(pitch)))
	sp := float32(math.Sin(float64(pitch)))
	sy := float32(math.Sin(float64(yaw)))
	cy := float32(math.Cos(float64(yaw)))
	return rl.Vector3Normalize(rl.Vector3{X: sy * cp, Y: sp, Z: cy * cp})
}
