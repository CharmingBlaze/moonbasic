//go:build cgo || (windows && !cgo)

package mbentity

import (
	"math"
	"path"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// PlayerBridgeWorldPos returns the entity world position (feet / pivot).
func (m *Module) PlayerBridgeWorldPos(id int64) (x, y, z float64, ok bool) {
	e := m.store().ents[id]
	if e == nil {
		return 0, 0, 0, false
	}
	wp := m.worldPos(e)
	return float64(wp.X), float64(wp.Y), float64(wp.Z), true
}

// PlayerBridgeApplyForce applies ENTITY.ADDFORCE-style velocity integration (host entity store).
func (m *Module) PlayerBridgeApplyForce(id int64, fx, fy, fz float32) bool {
	e := m.store().ents[id]
	if e == nil {
		return false
	}
	invM := float32(1)
	if e.mass > 1e-6 {
		invM = 1 / e.mass
	}
	e.vel.X += fx * invM
	e.vel.Y += fy * invM
	e.vel.Z += fz * invM
	e.static = false
	return true
}

// SurfaceMaterialHint returns a surface / footstep label from glTF metadata or Blender tag ("Default" if unknown).
func (m *Module) SurfaceMaterialHint(entityID int64) string {
	st := m.store()
	e := st.ents[entityID]
	if e == nil {
		return "Default"
	}
	if st.entMeta != nil && st.entMeta[entityID] != nil {
		row := st.entMeta[entityID]
		for _, k := range []string{"material", "Material", "surface", "Surface", "footstep", "Footstep", "MaterialName"} {
			if v, ok := metaGetCI(row, k); ok && strings.TrimSpace(v) != "" {
				return strings.TrimSpace(v)
			}
		}
	}
	if strings.TrimSpace(e.blenderTag) != "" {
		return strings.TrimSpace(e.blenderTag)
	}
	return "Default"
}

// PlayerBridgeSetWorldPos moves the entity root to world (x,y,z).
func (m *Module) PlayerBridgeSetWorldPos(id int64, x, y, z float32) bool {
	e := m.store().ents[id]
	if e == nil {
		return false
	}
	m.setLocalFromWorld(e, x, y, z)
	return true
}

// PlayerBridgeNearbyTagged returns entity ids within radius of (cx,cy,cz) whose name or Blender tag matches tagPat (case-insensitive path.Match glob).
func (m *Module) PlayerBridgeNearbyTagged(cx, cy, cz, radius float64, tagPat string) []int64 {
	tagPat = strings.TrimSpace(tagPat)
	tp := strings.ToUpper(tagPat)
	r2 := radius * radius
	var out []int64
	st := m.store()
	for _, e := range st.ents {
		if e == nil {
			continue
		}
		wp := m.worldPos(e)
		dx := float64(wp.X) - cx
		dy := float64(wp.Y) - cy
		dz := float64(wp.Z) - cz
		if dx*dx+dy*dy+dz*dz > r2 {
			continue
		}
		name := strings.ToUpper(strings.TrimSpace(e.name))
		btag := strings.ToUpper(strings.TrimSpace(e.blenderTag))
		ok, _ := path.Match(tp, name)
		if !ok {
			ok, _ = path.Match(tp, btag)
		}
		if !ok {
			continue
		}
		out = append(out, e.id)
	}
	return out
}

// PlayerBridgeClosestTagged returns the nearest matching entity id within radius, or 0 if none.
func (m *Module) PlayerBridgeClosestTagged(cx, cy, cz, radius float64, tagPat string) int64 {
	tagPat = strings.TrimSpace(tagPat)
	tp := strings.ToUpper(tagPat)
	r2 := radius * radius
	var best int64
	bestD2 := r2 + 1
	st := m.store()
	for _, e := range st.ents {
		if e == nil {
			continue
		}
		wp := m.worldPos(e)
		dx := float64(wp.X) - cx
		dy := float64(wp.Y) - cy
		dz := float64(wp.Z) - cz
		d2 := dx*dx + dy*dy + dz*dz
		if d2 > r2 {
			continue
		}
		name := strings.ToUpper(strings.TrimSpace(e.name))
		btag := strings.ToUpper(strings.TrimSpace(e.blenderTag))
		ok, _ := path.Match(tp, name)
		if !ok {
			ok, _ = path.Match(tp, btag)
		}
		if !ok {
			continue
		}
		if best == 0 || d2 < bestD2 {
			bestD2 = d2
			best = e.id
		}
	}
	return best
}

// PlayerBridgeEyeRay returns eye position (feet + eyeY) and forward unit vector from world rotation.
func (m *Module) PlayerBridgeEyeRay(id int64, eyeY float32) (ox, oy, oz, dx, dy, dz float64, ok bool) {
	e := m.store().ents[id]
	if e == nil {
		return 0, 0, 0, 0, 0, 0, false
	}
	wp := m.worldPos(e)
	pitch, yaw, _ := m.worldEuler(e)
	fwd := forwardFromYawPitch(yaw, pitch)
	ox = float64(wp.X)
	oy = float64(wp.Y) + float64(eyeY)
	oz = float64(wp.Z)
	dx = float64(fwd.X)
	dy = float64(fwd.Y)
	dz = float64(fwd.Z)
	return ox, oy, oz, dx, dy, dz, true
}

// PlayerBridgePickForward performs ENTITY.PICK-style AABB raycast (no Jolt). Returns 0 if none.
func (m *Module) PlayerBridgePickForward(id int64, rng float32) int64 {
	e := m.store().ents[id]
	if e == nil {
		return 0
	}
	pitch, yaw, _ := m.worldEuler(e)
	fwd := forwardFromYawPitch(yaw, pitch)
	origin := m.worldPos(e)
	end := rl.Vector3Add(origin, rl.Vector3Scale(fwd, rng))
	bestID := int64(0)
	bestT := float32(1e20)
	for _, s := range m.store().ents {
		if !s.static || s.id == e.id {
			continue
		}
		smn, smx := m.aabbWorldMinMax(s)
		t := rayAABB(origin, end, smn, smx)
		if t >= 0 && t < bestT {
			bestT = t
			bestID = s.id
		}
	}
	return bestID
}

// PlayerBridgeSetAnimSpeed sets animation playback speed multiplier for an entity.
func (m *Module) PlayerBridgeSetAnimSpeed(id int64, speed float32) bool {
	e := m.store().ents[id]
	if e == nil {
		return false
	}
	e.animSpeed = speed
	return true
}

// PlayerBridgeHorizontalSpeed returns sqrt(vx^2+vz^2) for SyncAnim helper.
func PlayerBridgeHorizontalSpeed(vx, vz float32) float32 {
	return float32(math.Sqrt(float64(vx*vx + vz*vz)))
}
