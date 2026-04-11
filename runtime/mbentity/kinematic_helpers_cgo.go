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
func (m *Module) ResolveKinematicCollision(id int64, radius float64, ignoreFloorLikeNormals bool) (hit bool) {
	st := m.store()
	e := st.ents[id]
	if e == nil {
		return false
	}

	r := float32(radius)
	if r <= 0 {
		return false
	}

	wp := m.worldPos(e)
	hasHit := false

	// Check against static entities
	for _, s := range st.staticEnts {
		if s.hidden {
			continue
		}
		smn, smx := m.aabbWorldMinMax(s)
		closest := rl.Vector3{
			X: float32(math.Max(float64(smn.X), math.Min(float64(wp.X), float64(smx.X)))),
			Y: float32(math.Max(float64(smn.Y), math.Min(float64(wp.Y), float64(smx.Y)))),
			Z: float32(math.Max(float64(smn.Z), math.Min(float64(wp.Z), float64(smx.Z)))),
		}
		d := rl.Vector3Distance(wp, closest)
		if d < r && d > 1e-6 {
			n := rl.Vector3Subtract(wp, closest)
			var nUse rl.Vector3
			if ignoreFloorLikeNormals {
				nUse = rl.Vector3{X: n.X, Y: 0, Z: n.Z}
				if rl.Vector3Length(nUse) < 1e-5 {
					continue
				}
				nUse = rl.Vector3Normalize(nUse)
			} else {
				nUse = rl.Vector3Normalize(n)
			}
			pen := r - d
			nwp := rl.Vector3Add(wp, rl.Vector3Scale(nUse, pen))
			m.setLocalFromWorld(e, nwp.X, nwp.Y, nwp.Z)
			wp = nwp // Update wp for next iteration
			hasHit = true
		} else if d <= 1e-6 {
			if ignoreFloorLikeNormals {
				continue
			}
			nwp := wp
			nwp.Y = smx.Y + r + 0.01
			m.setLocalFromWorld(e, nwp.X, nwp.Y, nwp.Z)
			wp = nwp
			hasHit = true
		}
	}

	return hasHit
}

// QueryKinematicFloor finds the highest static box top under the entity’s horizontal footprint.
// horizRadius expands the XZ footprint (capsule / cylinder radius).
// feetFromPivotDown is the distance from the entity pivot (world Y) down to the feet contact point —
// for MODEL.CREATECAPSULE and Jolt capsules this is height/2 (pivot at capsule center); do not use radius here.
func (m *Module) QueryKinematicFloor(id int64, horizRadius, feetFromPivotDown float64) (float32, bool) {
	st := m.store()
	e := st.ents[id]
	if e == nil {
		return 0, false
	}

	wp := m.worldPos(e)
	hR := float32(horizRadius)
	if hR < 0 {
		hR = 0
	}
	if feetFromPivotDown < 1e-6 {
		feetFromPivotDown = float64(hR)
	}

	px, py, pz := float64(wp.X), float64(wp.Y), float64(wp.Z)
	var best float64
	found := false

	for _, s := range st.ents {
		if !s.static {
			continue
		}
		sp := s.getPos()
		bx, by, bz := float64(sp.X), float64(sp.Y), float64(sp.Z)
		bw, bh, bd := float64(s.w), float64(s.h), float64(s.d)
		top := by + bh*0.5

		halfW := bw*0.5 + float64(hR)
		halfD := bd*0.5 + float64(hR)

		if math.Abs(px-bx) > halfW || math.Abs(pz-bz) > halfD {
			continue
		}

		feet := py - feetFromPivotDown
		if feet >= top-0.5 && feet <= top+0.5 {
			if !found || top > best {
				best = top
				found = true
			}
		}
	}

	return float32(best), found
}
