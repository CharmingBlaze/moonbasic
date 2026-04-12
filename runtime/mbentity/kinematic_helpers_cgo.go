//go:build cgo || (windows && !cgo)

package mbentity

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ResolveKinematicCollision (Internal Helper) performs sphere-vs-static sliding and pushes the entity out of geometry.
// Used by the Host KCC fallback on Windows.
//
// If ignoreFloorLikeNormals is true, separation uses XZ only (Y correction dropped). Otherwise any upward
// component from sphere-vs-AABB (top faces, edges, corners) fights host KCC’s analytic floor snap and reads as
// rubber-ball bouncing.
func (m *Module) ResolveKinematicCollision(id int64, radius float64, ignoreFloorLikeNormals bool) (hit bool, outNormal rl.Vector3) {
	st := m.store()
	e := st.ents[id]
	if e == nil {
		return false, rl.Vector3{}
	}

	wp := m.worldPos(e)
	newPos, hit, normal := m.ResolveKinematicCollisionAt(wp, id, radius, ignoreFloorLikeNormals)
	if hit {
		m.setLocalFromWorld(e, newPos.X, newPos.Y, newPos.Z)
	}
	return hit, normal
}

// ResolveKinematicCollisionAt performs the collision resolution at a specific world position.
func (m *Module) ResolveKinematicCollisionAt(wp rl.Vector3, id int64, radius float64, ignoreFloorLikeNormals bool) (rl.Vector3, bool, rl.Vector3) {
	st := m.store()
	r := float32(radius)
	if r <= 0 {
		return wp, false, rl.Vector3{}
	}

	hasHit := false
	var bestNormal rl.Vector3

	// Helper function for box vs sphere
	resolveBox := func(bx, by, bz, hw, hh, hd float32) {
		smn := rl.Vector3{X: bx - hw, Y: by - hh, Z: bz - hd}
		smx := rl.Vector3{X: bx + hw, Y: by + hh, Z: bz + hd}
		
		closest := rl.Vector3{
			X: float32(math.Max(float64(smn.X), math.Min(float64(wp.X), float64(smx.X)))),
			Y: float32(math.Max(float64(smn.Y), math.Min(float64(wp.Y), float64(smx.Y)))),
			Z: float32(math.Max(float64(smn.Z), math.Min(float64(wp.Z), float64(smx.Z)))),
		}
		d := rl.Vector3Distance(wp, closest)
		if d < r && d > 1e-6 {
			n := rl.Vector3Subtract(wp, closest)
			var nUse rl.Vector3
			var pen float32
			
			rawN := rl.Vector3Normalize(n)
			if rl.Vector3Length(rawN) < 1e-4 {
				return
			}
			if ignoreFloorLikeNormals {
				nUse = rl.Vector3{X: n.X, Y: 0, Z: n.Z}
				if rl.Vector3Length(nUse) < 1e-5 {
					return
				}
				nUse = rl.Vector3Normalize(nUse)
				dist2D := float32(math.Sqrt(float64(n.X*n.X + n.Z*n.Z)))
				pen = r - dist2D
			} else {
				nUse = rawN
				pen = r - d
			}
			
			if pen < 0 { pen = 0 }
			nwp := rl.Vector3Add(wp, rl.Vector3Scale(nUse, pen))
			wp = nwp
			hasHit = true
			if !hasHit || math.Abs(float64(rawN.Y)) < math.Abs(float64(bestNormal.Y)) {
				bestNormal = rawN
			}
		} else if d <= 1e-6 && !ignoreFloorLikeNormals {
			// Deep penetration fallback
			wp.Y = smx.Y + r + 0.01
			hasHit = true
			bestNormal = rl.Vector3{X: 0, Y: 1, Z: 0}
		}
	}

	// 1. Check against solid entities
	for _, s := range st.ents {
		if s == nil || s.id == id || s.hidden { continue }
		if !s.static && s.getExt().collType == 0 { continue }
		sp := s.getPos()
		resolveBox(sp.X, sp.Y, sp.Z, s.w*0.5, s.h*0.5, s.d*0.5)
	}

	// 2. Heap static bodies (physics3d stub Pos/Shape) — see kinematic_physicsstatic_*.go
	foreachPhysicsStaticBox(resolveBox)

	return wp, hasHit, bestNormal
}

func (m *Module) QueryKinematicFloor(id int64, horizRadius, feetFromPivotDown float64) (float32, bool) {
	st := m.store()
	e := st.ents[id]
	if e == nil {
		return 0, false
	}
	wp := m.worldPos(e)
	return m.QueryKinematicFloorAt(wp, id, horizRadius, feetFromPivotDown)
}

func (m *Module) QueryKinematicFloorAt(wp rl.Vector3, id int64, horizRadius, feetFromPivotDown float64) (float32, bool) {
	st := m.store()
	hR := float32(horizRadius)
	if hR < 0 { hR = 0 }
	if feetFromPivotDown < 1e-6 { feetFromPivotDown = float64(hR) }

	px, py, pz := float64(wp.X), float64(wp.Y), float64(wp.Z)
	var best float64
	found := false

	checkFloor := func(bx, by, bz, bw, bh, bd float64) {
		top := by + bh*0.5
		halfW := bw*0.5 + float64(hR)
		halfD := bd*0.5 + float64(hR)

		if math.Abs(px-bx) > halfW || math.Abs(pz-bz) > halfD { return }
		
		feet := py - feetFromPivotDown
		if feet >= top-2.0 && feet <= top+0.8 {
			if !found || top > best {
				best = top
				found = true
			}
		}
	}

	// 1. Check entities
	for _, s := range st.ents {
		if s == nil || s.id == id { continue }
		if !s.static && s.getExt().collType == 0 { continue }
		sp := s.getPos()
		checkFloor(float64(sp.X), float64(sp.Y), float64(sp.Z), float64(s.w), float64(s.h), float64(s.d))
	}

	// 2. Heap static bodies — see kinematic_physicsstatic_*.go
	foreachPhysicsStaticFloor(checkFloor)

	return float32(best), found
}
