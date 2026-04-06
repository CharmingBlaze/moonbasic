//go:build cgo || (windows && !cgo)

package mbentity

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (m *Module) worldRotQuat(e *ent) rl.Quaternion {
	ql := rl.QuaternionFromEuler(e.pitch, e.yaw, e.roll)
	if e.parentID == 0 {
		return ql
	}
	p := m.store().ents[e.parentID]
	if p == nil {
		return ql
	}
	qp := m.worldRotQuat(p)
	return rl.QuaternionMultiply(qp, ql)
}

func (m *Module) worldPos(e *ent) rl.Vector3 {
	if e.parentID == 0 {
		return e.pos
	}
	p := m.store().ents[e.parentID]
	if p == nil {
		return e.pos
	}
	pw := m.worldPos(p)
	pq := m.worldRotQuat(p)
	off := rl.Vector3RotateByQuaternion(e.pos, pq)
	return rl.Vector3Add(pw, off)
}

func (m *Module) worldEuler(e *ent) (pitch, yaw, roll float32) {
	q := m.worldRotQuat(e)
	v := rl.QuaternionToEuler(q)
	// Raylib: X=roll, Y=pitch, Z=yaw (see QUAT.TOEULER in mbmatrix)
	return v.Y, v.Z, v.X
}

func (m *Module) setLocalFromWorld(e *ent, wx, wy, wz float32) {
	if e.parentID == 0 {
		e.pos = rl.Vector3{X: wx, Y: wy, Z: wz}
		return
	}
	p := m.store().ents[e.parentID]
	if p == nil {
		e.pos = rl.Vector3{X: wx, Y: wy, Z: wz}
		return
	}
	pw := m.worldPos(p)
	pq := m.worldRotQuat(p)
	rel := rl.Vector3Subtract(rl.Vector3{X: wx, Y: wy, Z: wz}, pw)
	inv := rl.QuaternionInvert(pq)
	e.pos = rl.Vector3RotateByQuaternion(rel, inv)
}

func forwardFromYawPitch(yaw, pitch float32) rl.Vector3 {
	cp := float32(math.Cos(float64(pitch)))
	sp := float32(math.Sin(float64(pitch)))
	sy := float32(math.Sin(float64(yaw)))
	cy := float32(math.Cos(float64(yaw)))
	return rl.Vector3Normalize(rl.Vector3{X: sy * cp, Y: sp, Z: cy * cp})
}
