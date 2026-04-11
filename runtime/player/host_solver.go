package player

import (
	"math"
)

func (m *Module) Process(dt float64) {
	if dt <= 0 {
		return
	}

	// Update NAV intent first
	m.processNav(dt)

	// Update host KCC solvers (Windows/non-Jolt)
	if m.hostKCC != nil {
		for id, st := range m.hostKCC {
			m.updateHostKCC(id, st, dt)
		}
	}
}

func (m *Module) updateHostKCC(id int64, st *hostKCCState, dt float64) {
	if m.ent == nil {
		return
	}
	
	// 1. Gravity — while grounded, keep vy at 0 until airborne (avoids drift vs snap)
	if st.grounded {
		st.vy = 0
	} else {
		st.vy -= 28.0 * st.gravityScale * dt
	}

	// 2. Horizontal Movement
	m.ent.TranslateEntityByID(int(id), float32(st.vx*dt), 0, float32(st.vz*dt))

	// Sliding vs walls; ignore floor-like normals so we do not fight analytic snap (avoids rubber-banding on Y).
	m.ent.ResolveKinematicCollision(id, st.rad, true)

	// 3. Vertical Movement (no second sphere resolve — floor AABB pushes +Y each frame caused bouncing)
	m.ent.TranslateEntityByID(int(id), 0, float32(st.vy*dt), 0)

	// 4. Ground contact: pivot is capsule center (same as MODEL.CREATECAPSULE / DrawCapsule). Feet are Y − height/2.
	halfH := st.hei * 0.5
	if halfH < 1e-6 {
		halfH = st.rad
	}
	floorY, hit := m.ent.QueryKinematicFloor(id, st.rad, halfH)
	wp, ok := m.ent.GetWorldPosByID(int(id))
	if !ok {
		return
	}
	feetY := wp.Y - float32(halfH)
	const snapTol float32 = 0.12
	const leaveTol float32 = 0.22 // hysteresis: must leave this far above floor to become airborne

	onSupport := hit && feetY <= floorY+snapTol && st.vy <= 0
	stillSupported := hit && feetY <= floorY+leaveTol && st.vy <= 0 && st.grounded

	if onSupport || stillSupported {
		wp.Y = floorY + float32(halfH)
		m.ent.SetWorldPosByID(int(id), wp.X, wp.Y, wp.Z)
		st.grounded = true
		st.vy = 0
	} else {
		st.grounded = false
	}
}

func (m *Module) processNav(dt float64) {
	if m.kccNav == nil || m.ent == nil {
		return
	}
	for id, nav := range m.kccNav {
		if !nav.active {
			continue
		}
		
		wp, ok := m.ent.GetWorldPosByID(int(id))
		if !ok { continue }
		
		dx := float32(nav.tx) - wp.X
		dz := float32(nav.tz) - wp.Z
		dist := float32(math.Sqrt(float64(dx*dx + dz*dz)))
		
		if dist < float32(nav.arrival) {
			nav.active = false
			if st, ok := m.hostKCC[id]; ok {
				st.vx, st.vz = 0, 0
			}
			continue
		}

		// Soft stop damping (Implementation Detail: 0.2 units)
		spd := float32(nav.speed)
		if dist < 0.2 {
			spd *= (dist / 0.2)
		}

		nx := dx / dist
		nz := dz / dist
		
		// Look at target
		yaw := float32(math.Atan2(float64(nx), float64(nz))) * (180.0 / math.Pi)
		m.ent.RotateEntityAbsByID(int(id), 0, yaw, 0)
		
		if st, ok := m.hostKCC[id]; ok {
			st.vx = float64(nx) * float64(spd)
			st.vz = float64(nz) * float64(spd)
		}
	}
}
