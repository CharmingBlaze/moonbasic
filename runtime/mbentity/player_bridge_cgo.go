//go:build cgo || (windows && !cgo)

package mbentity

import (
	"math"
	"path"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// PlayerBridgeAnalyticFloorTopForFeet returns the highest static box top Y under (px,pz) when feet at footY
// can stand on that surface (same static AABB scan as ENTITY.FLOOR / queryFloorY, but pivot at feet).
func (m *Module) PlayerBridgeAnalyticFloorTopForFeet(px, footY, pz, horizRadius float64) (topY float64, ok bool) {
	var best float64
	found := false
	for _, s := range m.store().ents {
		if s == nil || !s.static {
			continue
		}
		sp := s.getPos()
		bx, by, bz := float64(sp.X), float64(sp.Y), float64(sp.Z)
		bw, bh, bd := float64(s.w), float64(s.h), float64(s.d)
		top := by + bh*0.5
		halfW := bw*0.5 + horizRadius
		halfD := bd*0.5 + horizRadius
		if math.Abs(px-bx) > halfW || math.Abs(pz-bz) > halfD {
			continue
		}
		if footY <= top+3.0 && footY >= top-0.25 {
			if !found || top > best {
				best = top
				found = true
			}
		}
	}
	if !found {
		return 0, false
	}
	return best, true
}

// PlayerBridgeClearScriptedMotion zeros scripted gravity and velocity for PLAYER.CREATE / KCC entities.
func (m *Module) PlayerBridgeClearScriptedMotion(id int64) bool {
	e := m.store().ents[id]
	if e == nil {
		return false
	}
	e.gravity = 0
	e.vel = rl.Vector3{}
	return true
}

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
	ext := e.getExt()
	if strings.TrimSpace(ext.blenderTag) != "" {
		return strings.TrimSpace(ext.blenderTag)
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
		ext := e.getExt()
		name := strings.ToUpper(strings.TrimSpace(ext.name))
		btag := strings.ToUpper(strings.TrimSpace(ext.blenderTag))
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
		ext := e.getExt()
		name := strings.ToUpper(strings.TrimSpace(ext.name))
		btag := strings.ToUpper(strings.TrimSpace(ext.blenderTag))
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
	e.getExt().animSpeed = speed
	return true
}

// PlayerBridgeHorizontalSpeed returns sqrt(vx^2+vz^2) for SyncAnim helper.
func PlayerBridgeHorizontalSpeed(vx, vz float32) float32 {
	return float32(math.Sqrt(float64(vx*vx + vz*vz)))
}
