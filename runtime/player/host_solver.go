package player

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (m *Module) Process(dt float64) {
	if dt <= 0 || dt > 0.5 { // ignore invalid or massive delta-times
		return
	}

	if m.hostKCC != nil {
		for id, st := range m.hostKCC {
			if st == nil { continue }
			m.updateHostKCC(id, st, dt)
		}
	}
}

func (m *Module) updateHostKCC(id int64, st *hostKCCState, dt float64) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[KCC_FATAL] Panic in updateHostKCC(id=%d): %v\n", id, r)
		}
	}()
	if m.ent == nil {
		return
	}

	// 1. Current Position Context
	var wp rl.Vector3
	if id < 0 {
		wp = rl.Vector3{X: float32(st.x), Y: float32(st.y), Z: float32(st.z)}
	} else {
		p, ok := m.ent.GetWorldPosByID(int(id))
		if !ok { return }
		wp = p
	}
	
	// 2. Gravity & Ground State
	if st.grounded {
		st.vy = -0.1
	} else {
		st.vy -= 32.0 * st.gravityScale * dt
	}

	// 3. Horizontal Displacement & Stair Climbing
	moveX := st.vx * dt
	moveZ := st.vz * dt
	
	if math.Abs(moveX) > 1e-5 || math.Abs(moveZ) > 1e-5 {
		oldPos := wp
		wp.Y += float32(st.stepH)
		wp.X += float32(moveX)
		wp.Z += float32(moveZ)
		
		nwp, hit, _ := m.ent.ResolveKinematicCollisionAt(wp, id, st.rad, false)
		if hit {
			wp = oldPos
			wp.X += float32(moveX)
			wp.Z += float32(moveZ)
		} else {
			wp = nwp
		}
	}

	// 4. Iterative Collision Resolution (Slide)
	for i := 0; i < 4; i++ {
		nwp, hit, normal := m.ent.ResolveKinematicCollisionAt(wp, id, st.rad, true)
		if !hit { break }
		wp = nwp
		dot := float32(st.vx)*normal.X + float32(st.vz)*normal.Z
		if dot < 0 {
			st.vx -= float64(normal.X * dot)
			st.vz -= float64(normal.Z * dot)
		}
	}

	// 5. Vertical Movement
	wp.Y += float32(st.vy * dt)

	// 6. Ground Support & Snapping
	halfH := st.hei * 0.5
	if halfH < st.rad { halfH = st.rad }
	
	floorY, hit := m.ent.QueryKinematicFloorAt(wp, id, st.rad, halfH)
	feetY := wp.Y - float32(halfH)
	snapDist := float32(st.stepH + 0.2)

	if hit && feetY <= floorY+snapDist && (st.vy <= 0 || st.grounded) {
		st.grounded = true
		st.vy = 0
		wp.Y = floorY + float32(halfH)
	} else {
		st.grounded = false
	}

	// 7. Commit Position
	if id < 0 {
		st.x, st.y, st.z = float64(wp.X), float64(wp.Y), float64(wp.Z)
	} else {
		m.ent.SetWorldPosByID(int(id), wp.X, wp.Y, wp.Z)
	}
}

func (m *Module) processNav(dt float64) {
	if m.kccNav == nil || m.ent == nil {
		return
	}
	for id, nav := range m.kccNav {
		if !nav.active { continue }
		
		var wx, wz float32
		if id < 0 {
			if st, ok := m.hostKCC[id]; ok {
				wx, wz = float32(st.x), float32(st.z)
			} else {
				continue
			}
		} else {
			p, ok := m.ent.GetWorldPosByID(int(id))
			if !ok {
				continue
			}
			wx, wz = p.X, p.Z
		}

		dx := float32(nav.tx) - wx
		dz := float32(nav.tz) - wz
		dist := float32(math.Sqrt(float64(dx*dx + dz*dz)))
		
		if dist < float32(nav.arrival) {
			nav.active = false
			if st, ok := m.hostKCC[id]; ok {
				st.vx, st.vz = 0, 0
			}
			continue
		}

		spd := float32(nav.speed)
		if dist < 0.2 { spd *= (dist / 0.2) }

		nx := dx / dist
		nz := dz / dist
		
		// Look at target (Entity bound only)
		if id >= 0 {
			yaw := float32(math.Atan2(float64(nx), float64(nz))) * (180.0 / math.Pi)
			m.ent.RotateEntityAbsByID(int(id), 0, yaw, 0)
		}
		
		if st, ok := m.hostKCC[id]; ok {
			st.vx = float64(nx) * float64(spd)
			st.vz = float64(nz) * float64(spd)
		}
	}
}
