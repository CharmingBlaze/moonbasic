//go:build cgo || (windows && !cgo)
//
// Platform: builds on Windows and Linux fullruntime (not linux-only). Jolt-backed raycasts and
// matrix sync run only where physics3d links Jolt (Linux+CGO). On Windows, stubs return no ray hits;
// ENTITY.ADDPHYSICS is a no-op there—see AGENTS.md “Physics sync & Jolt”.
//
package mbentity

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	mbphysics3d "moonbasic/runtime/physics3d"
	"moonbasic/runtime/water"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ground probe for Jolt-linked bodies: ray length scales with collider half-height plus a small
// skin past the feet (industry-style buffer). Fixed-length rays mis-size tall/short characters and
// worsen grounded flicker next to solver slop. CharacterVirtual (PLAYER.*) uses Jolt ground state
// instead; this path is for dynamic rigid bodies + ENTITY.MOVEWITHCAMERA linear velocity.
const (
	joltGroundRayStartLift = 0.2  // origin above pivot (same order as prior +0.25; avoids odd self-hits)
	joltGroundPastFeetSkin = 0.12 // ~10–15% past feet: ray length uses startLift + halfExtent + this
	joltGroundRayMin       = 0.45
	joltGroundRayMax       = 5.0
	// Visual Y snap: grounded + near analytic stand height + not jumping → pin display ty to hitY+half
	// (masks solver micro-bounce; does not move the Jolt body).
	joltGroundVisualSnapBand = 0.14
	joltGroundSnapVyMax      = 0.35 // above this upward vy → skip visual snap (jump)
	// Jolt contact/solver slop (~2cm): allow snap when ty is within band+slop of analytic stand height.
	joltSolverContactSlop = 0.02
	// Grounded vertical sleep: only micro ±vy (solver slop); do not clamp real jumps (vy ≫ this).
	joltGroundVySleepMax = 0.25
)

// joltColliderHalfExtentDown estimates pivot→lowest point along +Y for upright primitives.
// Physics shapes use the same dimensions (see ENTITY.PHYSICS / BDAddCapsule).
func joltColliderHalfExtentDown(e *ent) float64 {
	if e == nil {
		return 0.5
	}
	// Smart bottom-pivot offset: detects distance from center to feet for Primitives
	// and distance from pivot to mesh-bottom for Models, multiplied by scale.
	off := float64(e.physBottomOffset)
	if off < 1e-4 {
		// Fallback for objects with no offset (pivot already at feet or uninitialized)
		// but check h/2 if it's a centered primitive kind.
		switch e.kind {
		case entKindBox, entKindSphere, entKindCylinder, entKindCapsule:
			off = float64(e.h) * 0.5
			if e.kind == entKindSphere { off = float64(e.radius) }
			if e.kind == entKindCapsule { off = float64(e.cylH)*0.5 + float64(e.radius) }
		}
	}
	res := off * float64(e.scale.Y)
	if res < 0.01 {
		return 0.01 // Avoid zero-length rays
	}
	return res
}

// syncEntitiesFromPhysics copies world pose from the Jolt matrix buffer + body rotation into linked
// entities (parent-aware local TRS). Shared buffer columns 0–10 hold the 3×3 rotation written by
// syncSharedBuffers; translation uses 12–14 (with optional visual Y snap for physics-driven bodies).
func (m *Module) syncEntitiesFromPhysics() {
	buf := mbphysics3d.MatrixBufferForEntitySync()
	if len(buf) == 0 {
		return
	}
	st := m.store()
	for _, e := range st.ents {
		if e == nil || e.physBufIndex < 0 {
			continue
		}
		idx := e.physBufIndex * 16
		if idx+16 > len(buf) {
			continue
		}
		tx := buf[idx+12]
		ty := buf[idx+13]
		tz := buf[idx+14]

		if e.physicsDriven {
			half := joltColliderHalfExtentDown(e)
			maxDown := joltGroundRayStartLift + half + joltGroundPastFeetSkin
			if maxDown < joltGroundRayMin {
				maxDown = joltGroundRayMin
			}
			if maxDown > joltGroundRayMax {
				maxDown = joltGroundRayMax
			}
			// Probe from Jolt world translation (buffer), not stale entity pose — matches final ty/tz we apply.
			rayOx := float64(tx)
			rayOy := float64(ty) + joltGroundRayStartLift
			rayOz := float64(tz)
			nx, ny, nz, hitY, hit := mbphysics3d.RaycastDownGroundProbe(rayOx, rayOy, rayOz, maxDown)
			// Count only floor-like contacts (normal mostly +Y). A bare ray hit can be a wall
			// or edge and would flutter ENTITY.GROUNDED / jump coyote every frame ("bunny hop").
			groundNormal := hit && ny >= 0.28 && math.Abs(nx) < 0.95 && math.Abs(nz) < 0.95
			e.onGround = groundNormal
			if groundNormal {
				vx, vy, vz := mbphysics3d.GetLinearVelocityToIndex(e.physBufIndex)
				// Standing contact: zero tiny vertical velocity from solver ping-pong (XZ unchanged).
				if math.Abs(float64(vy)) < joltGroundVySleepMax {
					mbphysics3d.SetVelocityToIndex(e.physBufIndex, vx, 0, vz)
				}
				idealY := hitY + half
				verticalDiff := math.Abs(float64(ty) - idealY)
				if verticalDiff < joltGroundVisualSnapBand+joltSolverContactSlop {
					if vy < joltGroundSnapVyMax {
						ty = float32(idealY)
					}
				}
			}
		}
		m.setLocalFromWorld(e, tx, ty, tz)
		if qx, qy, qz, qw, ok := mbphysics3d.GetBodyQuaternionForBufferIndex(e.physBufIndex); ok {
			m.setLocalRotFromWorldQuat(e, rl.Quaternion{X: qx, Y: qy, Z: qz, W: qw})
		}
	}
}

func (m *Module) processAutoBuoyancy(dt float32) {
	if !m.autoBuoyancy {
		return
	}
	st := m.store()
	grav := mbphysics3d.GravityVec()
	for _, e := range st.ents {
		if e == nil || e.physBufIndex < 0 || !e.physicsDriven {
			continue
		}
		half := joltColliderHalfExtentDown(e)
		wp := m.worldPos(e)

		// Fraction 0..1 based on vertical overlap with water volumes.
		frac := water.EntitySubmergedFraction(m.h, wp.Y-float32(half), wp.Y+float32(half), wp.X, wp.Z)
		if frac > 0.01 {
			// Apply upward buoyancy force.
			// Base buoyancy = counter-gravity + a bit extra to float at equilibrium.
			// Force = mass * acceleration.
			buoyY := -grav.Y * frac * 1.5
			mbphysics3d.ApplyImpulseToIndex(e.physBufIndex, 0, buoyY*float32(dt), 0)
			// Drag in water (linear damping equivalent)
			mbphysics3d.ApplyImpulseToIndex(e.physBufIndex, 0, -0.5*float32(dt), 0) // Tiny resistance
		}
	}
}

func (m *Module) entLinkPhysBuffer(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.LINKPHYSBUFFER expects 2 arguments (entity, bufferIndex)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.LINKPHYSBUFFER: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.LINKPHYSBUFFER: unknown entity %d", id)
	}
	bi, ok := args[1].ToInt()
	if !ok || bi < 0 {
		return value.Nil, fmt.Errorf("ENTITY.LINKPHYSBUFFER: bufferIndex must be int >= 0")
	}
	e.physBufIndex = int(bi)
	e.physicsDriven = true
	// Scripted gravity + Jolt world gravity causes depenetration “bunny hop” on stacked integration.
	e.gravity = 0
	e.vel = rl.Vector3{}
	mbphysics3d.RegisterEntityBufferLink(id, int(bi))
	return value.Nil, nil
}

func (m *Module) entClearPhysBuffer(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.CLEARPHYSBUFFER expects 1 argument (entity)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.CLEARPHYSBUFFER: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.CLEARPHYSBUFFER: unknown entity")
	}
	e.physBufIndex = -1
	e.physicsDriven = false
	mbphysics3d.UnregisterEntityCollision(id)
	return value.Nil, nil
}

func registerPhysicsEntitySync(m *Module, r runtime.Registrar) {
	r.Register("ENTITY.LINKPHYSBUFFER", "entity", runtime.AdaptLegacy(m.entLinkPhysBuffer))
	r.Register("ENTITY.CLEARPHYSBUFFER", "entity", runtime.AdaptLegacy(m.entClearPhysBuffer))
	mbphysics3d.SetAfterPhysicsMatrixSync(m.syncEntitiesFromPhysics)
	m.installPickLayerLookup()
}

func (m *Module) installPickLayerLookup() {
	mbphysics3d.SetPickLayerLookup(func(id int64) (uint8, bool) {
		if m.h == nil {
			return 0, false
		}
		st := entityStores[m.h]
		if st == nil {
			return 0, false
		}
		e := st.ents[id]
		if e == nil {
			return 0, false
		}
		return e.collisionLayer, true
	})
}
