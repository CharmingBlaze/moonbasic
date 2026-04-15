//go:build (linux || windows) && cgo

package player

import (
	"fmt"
	"math"

	mbcamera "moonbasic/runtime/camera"
	mbentity "moonbasic/runtime/mbentity"
	mbmatrix "moonbasic/runtime/mbmatrix"
	mbphysics3d "moonbasic/runtime/physics3d"
	mwater "moonbasic/runtime/water"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

const defaultEyeY = 1.65

const playerCapsuleRadius = 0.4
const playerCapsuleHeight = 1.75

func (m *Module) physicsFixedDt() float64 {
	if m.h == nil {
		return 1.0 / 60.0
	}
	ph := mbphysics3d.GetModule(m.h)
	if ph == nil {
		return 1.0 / 60.0
	}
	return ph.FixedStepSeconds()
}

// playerCreateInternal allocates Jolt KCC for an entity; used by PLAYER.CREATE and CHARACTER.CREATE.

func (m *Module) playerCreateInternal(args []value.Value) (int64, error) {
	if m.h == nil || m.char == nil || m.ent == nil {
		return 0, fmt.Errorf("PLAYER.CREATE: not available (requires Linux+Jolt fullruntime)")
	}
	if len(args) != 1 && len(args) != 3 {
		return 0, fmt.Errorf("PLAYER.CREATE / CHAR.MAKE expects (entity) or (entity, radius#, height#)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return 0, fmt.Errorf("PLAYER.CREATE: invalid entity")
	}
	if _, dup := m.entToChar[id]; dup {
		return 0, fmt.Errorf("PLAYER.CREATE: entity already has a character controller")
	}
	rad, hei := float64(playerCapsuleRadius), float64(playerCapsuleHeight)
	if len(args) == 3 {
		r, ok1 := args[1].ToFloat()
		h, ok2 := args[2].ToFloat()
		if !ok1 || !ok2 || r <= 0 || h <= 0 {
			return 0, fmt.Errorf("CHAR.MAKE: radius and height must be positive numbers")
		}
		rad, hei = r, h
	}
	px, py, pz, ok := m.ent.PlayerBridgeWorldPos(id)
	if !ok {
		return 0, fmt.Errorf("PLAYER.CREATE: unknown entity")
	}
	h, err := m.char.AllocCharacter(rad, hei, px, py, pz, 0, 0.02, 100.0, 0.05)
	if err != nil {
		return 0, err
	}
	m.entToChar[id] = h
	_ = m.ent.PlayerBridgeClearScriptedMotion(id)
	m.lastHero = id
	return id, nil
}

func (m *Module) playerCreate(args []value.Value) (value.Value, error) {
	if _, err := m.playerCreateInternal(args); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) playerMove(args []value.Value) (value.Value, error) {
	if m.char == nil || m.ent == nil {
		return value.Nil, fmt.Errorf("PLAYER.MOVE: not available on this platform")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("PLAYER.MOVE expects (entity, velocityX, velocityZ) world units/sec")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.MOVE: invalid entity")
	}
	vx, ok1 := args[1].ToFloat()
	vz, ok2 := args[2].ToFloat()
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("PLAYER.MOVE: velocities must be numeric")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.MOVE: call PLAYER.CREATE first")
	}
	m.syncKCCAmbientWater(id, ch)
	dt := m.physicsFixedDt()
	if err := m.char.CharacterMoveXZVelocity(ch, vx, vz, dt); err != nil {
		return value.Nil, err
	}
	x, y, z, ok := m.char.CharacterPosition(ch)
	if ok {
		_ = m.ent.PlayerBridgeSetWorldPos(id, float32(x), float32(y), float32(z))
	}
	if geid, ok := m.grab[id]; ok && geid > 0 && m.ent != nil {
		cx, cy, cz, ok := m.ent.PlayerBridgeWorldPos(id)
		if ok {
			_, _, _, dx, _, dz, ok2 := m.ent.PlayerBridgeEyeRay(id, 0.15)
			if ok2 {
				flen := math.Hypot(dx, dz)
				if flen > 1e-6 {
					fx := float32(dx / flen * 0.55)
					fz := float32(dz / flen * 0.55)
					_ = m.ent.PlayerBridgeSetWorldPos(geid, float32(cx)+fx, float32(cy)+0.35, float32(cz)+fz)
				}
			}
		}
	}
	return value.Nil, nil
}

// playerCharMoveDir implements CHAR.MOVE(entity, dirX, dirZ, speed): world XZ velocity = dir * speed (typ. dir ∈ {-1,0,1}).
func (m *Module) playerCharMoveDir(args []value.Value) (value.Value, error) {
	if m.char == nil || m.ent == nil {
		return value.Nil, fmt.Errorf("CHAR.MOVE: not available on this platform")
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CHAR.MOVE expects (entity, dirX#, dirZ#, speed#)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("CHAR.MOVE: invalid entity")
	}
	dx, ok1 := args[1].ToFloat()
	dz, ok2 := args[2].ToFloat()
	spd, ok3 := args[3].ToFloat()
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("CHAR.MOVE: numeric arguments required")
	}
	if spd < 0 {
		return value.Nil, fmt.Errorf("CHAR.MOVE: speed must be non-negative")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("CHAR.MOVE: call PLAYER.CREATE / CHAR.MAKE first")
	}
	m.syncKCCAmbientWater(id, ch)
	vx := dx * spd
	vz := dz * spd
	dt := m.physicsFixedDt()
	if err := m.char.CharacterMoveXZVelocity(ch, vx, vz, dt); err != nil {
		return value.Nil, err
	}
	x, y, z, ok := m.char.CharacterPosition(ch)
	if ok {
		_ = m.ent.PlayerBridgeSetWorldPos(id, float32(x), float32(y), float32(z))
	}
	return value.Nil, nil
}

// playerMoveWithCamera implements CHAR.MOVEWITHCAMERA / PLAYER.MOVEWITHCAMERA — camera XZ walk basis × input axes × speed.
func (m *Module) playerMoveWithCamera(args []value.Value) (value.Value, error) {
	if m.h == nil || m.char == nil || m.ent == nil {
		return value.Nil, fmt.Errorf("PLAYER.MOVEWITHCAMERA: not available on this platform")
	}
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("PLAYER.MOVEWITHCAMERA expects (entity, camera, forwardAxis#, strafeAxis#, speed#)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.MOVEWITHCAMERA: invalid entity")
	}
	if args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("PLAYER.MOVEWITHCAMERA: camera handle required")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.MOVEWITHCAMERA: call PLAYER.CREATE first")
	}
	f, ok1 := args[2].ToFloat()
	s, ok2 := args[3].ToFloat()
	spd, ok3 := args[4].ToFloat()
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("PLAYER.MOVEWITHCAMERA: forward/strafe/speed must be numeric")
	}
	if spd < 0 {
		return value.Nil, fmt.Errorf("PLAYER.MOVEWITHCAMERA: speed must be non-negative")
	}
	camH := heap.Handle(args[1].IVal)
	fwd, right, err := mbcamera.CameraXZWalkBasis(m.h, camH)
	if err != nil {
		return value.Nil, err
	}
	m.syncKCCAmbientWater(id, ch)
	vx := (float64(fwd.X)*f + float64(right.X)*s) * spd
	vz := (float64(fwd.Z)*f + float64(right.Z)*s) * spd
	dt := m.physicsFixedDt()
	if err := m.char.CharacterMoveXZVelocity(ch, vx, vz, dt); err != nil {
		return value.Nil, err
	}
	x, y, z, ok := m.char.CharacterPosition(ch)
	if ok {
		_ = m.ent.PlayerBridgeSetWorldPos(id, float32(x), float32(y), float32(z))
	}
	return value.Nil, nil
}

func (m *Module) playerNavTo(args []value.Value) (value.Value, error) {
	if m.kccNav == nil {
		m.kccNav = make(map[int64]*kccNavState)
	}
	if len(args) < 4 || len(args) > 6 {
		return value.Nil, fmt.Errorf("PLAYER.NAVTO expects (entity, targetX#, targetZ#, speed# [, arrivalXZ# [, brakeDist#]])")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.NAVTO: invalid entity")
	}
	if _, ok := m.entToChar[id]; !ok {
		return value.Nil, fmt.Errorf("PLAYER.NAVTO: call PLAYER.CREATE first")
	}
	tx, _ := args[1].ToFloat()
	tz, _ := args[2].ToFloat()
	spd, _ := args[3].ToFloat()
	if spd < 0 {
		return value.Nil, fmt.Errorf("PLAYER.NAVTO: speed must be non-negative")
	}
	arr := 0.2
	if len(args) >= 5 {
		if a, _ := args[4].ToFloat(); a > 0 {
			arr = a
		}
	}
	brake := 0.75
	if len(args) == 6 {
		if b, _ := args[5].ToFloat(); b > 0 {
			brake = b
		}
	}
	m.kccNav[id] = &kccNavState{mode: kccNavGoto, active: true, tx: tx, tz: tz, speed: spd, arrival: arr, brake: brake}
	return value.Nil, nil
}

func (m *Module) playerNavChase(args []value.Value) (value.Value, error) {
	if m.kccNav == nil {
		m.kccNav = make(map[int64]*kccNavState)
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("NAV.CHASE expects (entity, targetEntity#, standoffGap#, speed#)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("NAV.CHASE: invalid entity")
	}
	if _, ok := m.entToChar[id]; !ok {
		return value.Nil, fmt.Errorf("NAV.CHASE: call PLAYER.CREATE / CHAR.MAKE first")
	}
	tid, ok := m.playerEntID(args[1])
	if !ok || tid < 1 {
		return value.Nil, fmt.Errorf("NAV.CHASE: invalid target entity")
	}
	gap, ok := args[2].ToFloat()
	if !ok || gap < 0 {
		return value.Nil, fmt.Errorf("NAV.CHASE: gap must be non-negative")
	}
	spd, ok := args[3].ToFloat()
	if !ok || spd < 0 {
		return value.Nil, fmt.Errorf("NAV.CHASE: speed must be non-negative")
	}
	m.kccNav[id] = &kccNavState{
		mode: kccNavChase, active: true, chaseTarget: tid, chaseGap: gap, speed: spd,
		arrival: 0.2, brake: 0.75,
	}
	return value.Nil, nil
}

func (m *Module) playerNavPatrol(args []value.Value) (value.Value, error) {
	if m.kccNav == nil {
		m.kccNav = make(map[int64]*kccNavState)
	}
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("NAV.PATROL expects (entity, ax#, az#, bx#, bz#, speed#)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("NAV.PATROL: invalid entity")
	}
	if _, ok := m.entToChar[id]; !ok {
		return value.Nil, fmt.Errorf("NAV.PATROL: call PLAYER.CREATE / CHAR.MAKE first")
	}
	ax, _ := args[1].ToFloat()
	az, _ := args[2].ToFloat()
	bx, _ := args[3].ToFloat()
	bz, _ := args[4].ToFloat()
	spd, ok := args[5].ToFloat()
	if !ok || spd < 0 {
		return value.Nil, fmt.Errorf("NAV.PATROL: speed must be non-negative")
	}
	m.kccNav[id] = &kccNavState{
		mode: kccNavPatrol, active: true,
		patrolAX: ax, patrolAZ: az, patrolBX: bx, patrolBZ: bz,
		speed: spd, arrival: 0.2, brake: 0.75,
		patrolToB: true,
	}
	return value.Nil, nil
}

func (m *Module) playerNavUpdate(args []value.Value) (value.Value, error) {
	if m.char == nil || m.ent == nil || m.kccNav == nil {
		return value.Nil, nil
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.NAVUPDATE expects (entity)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.NAVUPDATE: invalid entity")
	}
	st, ok := m.kccNav[id]
	if !ok || st == nil || !st.active {
		return value.Nil, nil
	}
	ch, ok := m.entToChar[id]
	if !ok {
		st.active = false
		return value.Nil, nil
	}
	m.syncKCCAmbientWater(id, ch)
	px, _, pz, ok := m.ent.PlayerBridgeWorldPos(id)
	if !ok {
		return value.Nil, nil
	}

	var tx, tz float64
	switch st.mode {
	case kccNavGoto:
		tx, tz = st.tx, st.tz
	case kccNavChase:
		var ok2 bool
		tx, _, tz, ok2 = m.ent.PlayerBridgeWorldPos(st.chaseTarget)
		if !ok2 {
			return value.Nil, nil
		}
		if math.Hypot(tx-px, tz-pz) <= st.chaseGap {
			dt := m.physicsFixedDt()
			_ = m.char.CharacterMoveXZVelocity(ch, 0, 0, dt)
			x, y, z, ok2 := m.char.CharacterPosition(ch)
			if ok2 {
				_ = m.ent.PlayerBridgeSetWorldPos(id, float32(x), float32(y), float32(z))
			}
			return value.Nil, nil
		}
	case kccNavPatrol:
		if st.patrolToB {
			tx, tz = st.patrolBX, st.patrolBZ
		} else {
			tx, tz = st.patrolAX, st.patrolAZ
		}
	default:
		tx, tz = st.tx, st.tz
	}

	dx := tx - px
	dz := tz - pz
	dist := math.Hypot(dx, dz)
	switch st.mode {
	case kccNavPatrol:
		if dist <= st.arrival {
			dt := m.physicsFixedDt()
			_ = m.char.CharacterMoveXZVelocity(ch, 0, 0, dt)
			x, y, z, ok2 := m.char.CharacterPosition(ch)
			if ok2 {
				_ = m.ent.PlayerBridgeSetWorldPos(id, float32(x), float32(y), float32(z))
			}
			st.patrolToB = !st.patrolToB
			return value.Nil, nil
		}
	case kccNavGoto:
		if dist <= st.arrival {
			dt := m.physicsFixedDt()
			_ = m.char.CharacterMoveXZVelocity(ch, 0, 0, dt)
			x, y, z, ok2 := m.char.CharacterPosition(ch)
			if ok2 {
				_ = m.ent.PlayerBridgeSetWorldPos(id, float32(x), float32(y), float32(z))
			}
			st.active = false
			return value.Nil, nil
		}
	}

	spd := st.speed
	if st.brake > 0 && dist < st.brake {
		t := dist / st.brake
		spd *= t * t
	}
	vx := (dx / dist) * spd
	vz := (dz / dist) * spd
	dt := m.physicsFixedDt()
	if err := m.char.CharacterMoveXZVelocity(ch, vx, vz, dt); err != nil {
		return value.Nil, err
	}
	x, y, z, ok := m.char.CharacterPosition(ch)
	if ok {
		_ = m.ent.PlayerBridgeSetWorldPos(id, float32(x), float32(y), float32(z))
	}
	return value.Nil, nil
}

func (m *Module) playerSetPadding(args []value.Value) (value.Value, error) {
	if m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.SETPADDING: not available on this platform")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETPADDING expects (entity, padding#)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.SETPADDING: invalid entity")
	}
	pad, ok := args[1].ToFloat()
	if !ok || pad <= 0 {
		return value.Nil, fmt.Errorf("PLAYER.SETPADDING: padding must be > 0")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SETPADDING: call PLAYER.CREATE first")
	}
	newH, err := m.char.SetCharacterPadding(ch, float32(pad))
	if err != nil {
		return value.Nil, err
	}
	m.entToChar[id] = newH
	return value.Nil, nil
}

func (m *Module) playerJump(args []value.Value) (value.Value, error) {
	if m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.JUMP: not available on this platform")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.JUMP expects (entity, impulseY)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.JUMP: invalid entity")
	}
	imp, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.JUMP: impulse must be numeric")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.JUMP: call PLAYER.CREATE first")
	}
	if err := m.char.CharacterJump(ch, imp); err != nil {
		return value.Nil, err
	}
	if m.ent != nil {
		x, y, z, ok := m.char.CharacterPosition(ch)
		if ok {
			_ = m.ent.PlayerBridgeSetWorldPos(id, float32(x), float32(y), float32(z))
		}
	}
	return value.Nil, nil
}

func (m *Module) playerIsGrounded(args []value.Value) (value.Value, error) {
	if m.char == nil {
		return value.FromBool(false), nil
	}
	if len(args) > 2 {
		return value.Nil, fmt.Errorf("PLAYER.ISGROUNDED expects (), (entity), (entity, coyoteTimeSec#)")
	}
	var id int64
	var ok bool
	switch len(args) {
	case 0:
		id, ok = m.kccSubjectID(args)
		if !ok {
			return value.Nil, fmt.Errorf("PLAYER.ISGROUNDED: %s", kccErrNoSubject)
		}
	case 1, 2:
		id, ok = m.kccSubjectID(args[:1])
		if !ok || id < 1 {
			return value.Nil, fmt.Errorf("PLAYER.ISGROUNDED: invalid entity")
		}
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.FromBool(false), nil
	}
	g, err := m.char.CharacterIsGrounded(ch)
	if err != nil {
		return value.Nil, err
	}
	var now float64
	if m.h != nil {
		if ph := mbphysics3d.GetModule(m.h); ph != nil {
			now = ph.SimTimeSeconds()
		}
	}
	if m.kccLastGroundedAt == nil {
		m.kccLastGroundedAt = make(map[int64]float64)
	}
	if g {
		m.kccLastGroundedAt[id] = now
		return value.FromBool(true), nil
	}
	if len(args) == 2 {
		grace, ok := args[1].ToFloat()
		if ok && grace > 0 {
			if t, ok := m.kccLastGroundedAt[id]; ok && now-t <= grace {
				return value.FromBool(true), nil
			}
		}
	}
	return value.FromBool(false), nil
}

func (m *Module) playerGetLookTarget(args []value.Value) (value.Value, error) {
	if m.ent == nil {
		return value.FromInt(0), nil
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.GETLOOKTARGET expects (entity, maxDist)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETLOOKTARGET: invalid entity")
	}
	maxd, ok := args[1].ToFloat()
	if !ok || maxd <= 0 {
		return value.Nil, fmt.Errorf("PLAYER.GETLOOKTARGET: maxDist must be > 0")
	}
	ox, oy, oz, dx, dy, dz, ok := m.ent.PlayerBridgeEyeRay(id, defaultEyeY)
	if !ok {
		return value.FromInt(0), nil
	}
	hit := mbphysics3d.PickCastEntityID(ox, oy, oz, dx, dy, dz, maxd)
	if hit == id {
		// First hit is often the character capsule if registered; fall back to mesh AABB pick.
		hit = 0
	}
	if hit == 0 {
		hit = m.ent.PlayerBridgePickForward(id, float32(maxd))
		if hit == id {
			hit = 0
		}
	}
	return value.FromInt(hit), nil
}

func (m *Module) playerGetNearby(args []value.Value) (value.Value, error) {
	if m.ent == nil {
		return value.Nil, nil
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("PLAYER.GETNEARBY expects (entity, radius, tag)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETNEARBY: invalid entity")
	}
	rad, ok := args[1].ToFloat()
	if !ok || rad < 0 {
		return value.Nil, fmt.Errorf("PLAYER.GETNEARBY: radius must be >= 0")
	}
	if args[2].Kind != value.KindString {
		return value.Nil, fmt.Errorf("PLAYER.GETNEARBY: tag must be string")
	}
	tag, ok := m.h.GetString(int32(args[2].IVal))
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.GETNEARBY: invalid tag string")
	}
	cx, cy, cz, ok := m.ent.PlayerBridgeWorldPos(id)
	if !ok {
		return value.Nil, nil
	}
	ids := m.ent.PlayerBridgeNearbyTagged(cx, cy, cz, rad, tag)
	return allocFloatArray(m, ids)
}

func (m *Module) playerOnTrigger(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.ONTRIGGER expects (entity, callbackFunc)")
	}
	return value.Nil, fmt.Errorf("PLAYER.ONTRIGGER: VM callback from physics not wired — use LEVEL.BINDSCRIPT + collision checks or PHYSICS3D callbacks")
}

func (m *Module) playerSetState(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETSTATE expects (entity, state)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.SETSTATE: invalid entity")
	}
	st, ok := args[1].ToInt()
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SETSTATE: state must be int (use STATE_IDLE, STATE_WALKING, …)")
	}
	m.state[id] = int32(st)
	return value.Nil, nil
}

func (m *Module) playerSyncAnim(args []value.Value) (value.Value, error) {
	if m.char == nil || m.ent == nil {
		return value.Nil, fmt.Errorf("PLAYER.SYNCANIM: not available on this platform")
	}
	if len(args) != 1 && len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SYNCANIM expects (entity [, scale])")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.SYNCANIM: invalid entity")
	}
	scale := 1.0
	if len(args) == 2 {
		if s, ok := args[1].ToFloat(); ok {
			scale = s
		}
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SYNCANIM: call PLAYER.CREATE first")
	}
	vx, _, vz, ok := m.char.CharacterLinearVelocity(ch)
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SYNCANIM: internal")
	}
	hs := mbentity.PlayerBridgeHorizontalSpeed(float32(vx), float32(vz))
	sp := float32(hs * float32(scale))
	_ = m.ent.PlayerBridgeSetAnimSpeed(id, sp)
	return value.Nil, nil
}

func (m *Module) playerSetStepHeight(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETSTEPHEIGHT expects (entity, height)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.SETSTEPHEIGHT: invalid entity")
	}
	h, ok := args[1].ToFloat()
	if !ok || h < 0 {
		return value.Nil, fmt.Errorf("PLAYER.SETSTEPHEIGHT: height must be >= 0")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SETSTEPHEIGHT: call PLAYER.CREATE first")
	}
	m.stepHeight[id] = h
	if err := m.char.SetCharacterWalkStairsStepUp(ch, float32(h)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) playerSetSlopeLimit(args []value.Value) (value.Value, error) {
	if m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.SETSLOPELIMIT: not available on this platform")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETSLOPELIMIT expects (entity, maxSlopeDegrees)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.SETSLOPELIMIT: invalid entity")
	}
	deg, ok := args[1].ToFloat()
	if !ok || deg <= 0 || deg >= 90 {
		return value.Nil, fmt.Errorf("PLAYER.SETSLOPELIMIT: angle must be in (0, 90) degrees")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SETSLOPELIMIT: call PLAYER.CREATE first")
	}
	rad, fh := float64(playerCapsuleRadius), float64(playerCapsuleHeight)
	if cr, chh, ok := m.char.CharacterCapsuleDims(ch); ok {
		rad, fh = cr, chh
	}
	newH, err := m.char.RecreateCharacterWithSlope(ch, rad, fh, deg)
	if err != nil {
		return value.Nil, err
	}
	m.entToChar[id] = newH
	return value.Nil, nil
}

func (m *Module) playerGetGroundState(args []value.Value) (value.Value, error) {
	if len(args) > 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETGROUNDSTATE expects () or (entity)")
	}
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETGROUNDSTATE: %s", kccErrNoSubject)
		}
		return value.FromInt(3), nil
	}
	if id < 1 {
		return value.FromInt(3), nil
	}
	ch, ok := m.entToChar[id]
	if !ok || m.char == nil {
		return value.FromInt(3), nil
	}
	gi, ok := m.char.CharacterGroundStateInt(ch)
	if !ok {
		return value.FromInt(3), nil
	}
	return value.FromInt(int64(gi)), nil
}

func (m *Module) playerIsOnSteepSlope(args []value.Value) (value.Value, error) {
	if len(args) > 1 {
		return value.Nil, fmt.Errorf("PLAYER.ISONSTEEPSLOPE expects () or (entity)")
	}
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.ISONSTEEPSLOPE: %s", kccErrNoSubject)
		}
		return value.FromBool(false), nil
	}
	if id < 1 {
		return value.FromBool(false), nil
	}
	ch, ok := m.entToChar[id]
	if !ok || m.char == nil {
		return value.FromBool(false), nil
	}
	gi, ok := m.char.CharacterGroundStateInt(ch)
	if !ok {
		return value.FromBool(false), nil
	}
	// 1 = Jolt GroundStateOnSteepGround
	return value.FromBool(gi == 1), nil
}

func (m *Module) playerGetVelocity(args []value.Value) (value.Value, error) {
	if m.char == nil || m.h == nil {
		return value.Nil, fmt.Errorf("PLAYER.GETVELOCITY: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETVELOCITY expects (entity)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETVELOCITY: invalid entity")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.GETVELOCITY: call PLAYER.CREATE first")
	}
	vx, vy, vz, ok := m.char.CharacterLinearVelocity(ch)
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.GETVELOCITY: internal")
	}
	return mbmatrix.AllocVec3Value(m.h, float32(vx), float32(vy), float32(vz))
}

func (m *Module) playerTeleport(args []value.Value) (value.Value, error) {
	if m.char == nil || m.ent == nil {
		return value.Nil, fmt.Errorf("PLAYER.TELEPORT: not available on this platform")
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("PLAYER.TELEPORT expects (entity, x, y, z)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.TELEPORT: invalid entity")
	}
	x, ok1 := args[1].ToFloat()
	y, ok2 := args[2].ToFloat()
	z, ok3 := args[3].ToFloat()
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("PLAYER.TELEPORT: position must be numeric")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.TELEPORT: call PLAYER.CREATE first")
	}
	if err := m.char.CharacterTeleport(ch, x, y, z); err != nil {
		return value.Nil, err
	}
	_ = m.ent.PlayerBridgeSetWorldPos(id, float32(x), float32(y), float32(z))
	return value.Nil, nil
}

func (m *Module) playerSetGravityScale(args []value.Value) (value.Value, error) {
	if m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.SETGRAVITYSCALE: not available on this platform")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETGRAVITYSCALE expects (entity, scale)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.SETGRAVITYSCALE: invalid entity")
	}
	sc, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SETGRAVITYSCALE: scale must be numeric")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SETGRAVITYSCALE: call PLAYER.CREATE first")
	}
	if err := m.char.SetCharacterGravityScale(ch, sc); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) playerGetCrouch(args []value.Value) (value.Value, error) {
	if m.char == nil {
		return value.FromBool(false), nil
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETCROUCH expects (entity)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETCROUCH: invalid entity")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.FromBool(false), nil
	}
	return value.FromBool(m.char.CharacterCrouch(ch)), nil
}

func (m *Module) playerSetCrouch(args []value.Value) (value.Value, error) {
	if m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.SETCROUCH: not available on this platform")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETCROUCH expects (entity, enabled)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.SETCROUCH: invalid entity")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SETCROUCH: call PLAYER.CREATE first")
	}
	var en bool
	switch args[1].Kind {
	case value.KindBool:
		en = args[1].IVal != 0
	case value.KindInt:
		en = args[1].IVal != 0
	default:
		if f, ok := args[1].ToFloat(); ok {
			en = f != 0
		}
	}
	newH, err := m.char.SetCharacterCrouch(ch, en)
	if err != nil {
		return value.Nil, err
	}
	if newH != ch {
		m.entToChar[id] = newH
	}
	return value.Nil, nil
}

func (m *Module) playerSetJumpBuffer(args []value.Value) (value.Value, error) {
	if m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.SETJUMPBUFFER: not available on this platform")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETJUMPBUFFER expects (entity, seconds#)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.SETJUMPBUFFER: invalid entity")
	}
	sec, ok := args[1].ToFloat()
	if !ok || sec < 0 {
		return value.Nil, fmt.Errorf("PLAYER.SETJUMPBUFFER: seconds must be >= 0")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SETJUMPBUFFER: call PLAYER.CREATE first")
	}
	if err := m.char.SetCharacterJumpBuffer(ch, sec); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) playerSetAirControl(args []value.Value) (value.Value, error) {
	if m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.SETAIRCONTROL: not available on this platform")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETAIRCONTROL expects (entity, scale#)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.SETAIRCONTROL: invalid entity")
	}
	s, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SETAIRCONTROL: scale must be numeric")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SETAIRCONTROL: call PLAYER.CREATE first")
	}
	if err := m.char.SetCharacterAirControl(ch, s); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) playerSetGroundControl(args []value.Value) (value.Value, error) {
	if m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.SETGROUNDCONTROL: not available on this platform")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETGROUNDCONTROL expects (entity, scale#)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.SETGROUNDCONTROL: invalid entity")
	}
	s, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SETGROUNDCONTROL: scale must be numeric")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SETGROUNDCONTROL: call PLAYER.CREATE first")
	}
	if err := m.char.SetCharacterGroundControl(ch, s); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) playerSwim(args []value.Value) (value.Value, error) {
	if m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.SWIM: not available on this platform")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("PLAYER.SWIM expects (entity, buoyancy, drag)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.SWIM: invalid entity")
	}
	buoy, ok1 := args[1].ToFloat()
	drag, ok2 := args[2].ToFloat()
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("PLAYER.SWIM: buoyancy and drag must be numeric")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SWIM: call PLAYER.CREATE first")
	}
	on := buoy > 1e-9 || drag > 1e-9
	if m.swimManual == nil {
		m.swimManual = make(map[int64]bool)
	}
	m.swimManual[id] = on
	if err := m.char.SetCharacterSwim(ch, buoy, drag, on); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) playerSetVelocity(args []value.Value) (value.Value, error) {
	if m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.SETVELOCITY: not available on this platform")
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("PLAYER.SETVELOCITY expects (entity, vx#, vy#, vz#)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.SETVELOCITY: invalid entity")
	}
	vx, _ := args[1].ToFloat()
	vy, _ := args[2].ToFloat()
	vz, _ := args[3].ToFloat()
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SETVELOCITY: call PLAYER.CREATE first")
	}
	if err := m.char.SetCharacterLinearVelocity(ch, vx, vy, vz); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) playerAddImpulse(args []value.Value) (value.Value, error) {
	if m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.ADDIMPULSE: not available on this platform")
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("PLAYER.ADDIMPULSE expects (entity, dvx#, dvy#, dvz#)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.ADDIMPULSE: invalid entity")
	}
	ix, _ := args[1].ToFloat()
	iy, _ := args[2].ToFloat()
	iz, _ := args[3].ToFloat()
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.ADDIMPULSE: call PLAYER.CREATE first")
	}
	cvx, cvy, cvz, ok := m.char.CharacterLinearVelocity(ch)
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.ADDIMPULSE: internal")
	}
	if err := m.char.SetCharacterLinearVelocity(ch, cvx+ix, cvy+iy, cvz+iz); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) playerGetSubmergedFraction(args []value.Value) (value.Value, error) {
	if m.h == nil || m.ent == nil || m.char == nil {
		return value.FromFloat(0), nil
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETSUBMERGEDFACTOR expects (entity)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETSUBMERGEDFACTOR: invalid entity")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.FromFloat(0), nil
	}
	x, y, z, ok := m.ent.PlayerBridgeWorldPos(id)
	if !ok {
		return value.FromFloat(0), nil
	}
	_, fh, ok2 := m.char.CharacterCapsuleDims(ch)
	if !ok2 {
		fh = playerCapsuleHeight
	}
	mnY := float32(y) - float32(fh)*0.5
	mxY := float32(y) + float32(fh)*0.5
	f := mwater.EntitySubmergedFraction(m.h, mnY, mxY, float32(x), float32(z))
	return value.FromFloat(float64(f)), nil
}

func (m *Module) playerIsSubmerged(args []value.Value) (value.Value, error) {
	v, err := m.playerGetSubmergedFraction(args)
	if err != nil {
		return value.Nil, err
	}
	f, ok := v.ToFloat()
	if !ok {
		return value.FromBool(false), nil
	}
	return value.FromBool(f > 0.45), nil
}

func (m *Module) syncKCCAmbientWater(id int64, ch heap.Handle) {
	if m.char == nil || m.ent == nil || m.h == nil {
		return
	}
	if m.swimManual != nil && m.swimManual[id] {
		return
	}
	x, y, z, ok := m.ent.PlayerBridgeWorldPos(id)
	if !ok {
		return
	}
	_, fh, ok2 := m.char.CharacterCapsuleDims(ch)
	if !ok2 {
		fh = playerCapsuleHeight
	}
	mnY := float32(y) - float32(fh)*0.5
	mxY := float32(y) + float32(fh)*0.5
	frac := mwater.EntitySubmergedFraction(m.h, mnY, mxY, float32(x), float32(z))
	if frac <= 0.02 {
		_ = m.char.SetCharacterSwim(ch, 0, 0, false)
		return
	}
	_ = m.char.SetCharacterSwim(ch, float64(frac)*0.75, float64(frac)*1.85, true)
}

func (m *Module) playerSetStepOffset(args []value.Value) (value.Value, error) {
	return m.playerSetStepHeight(args)
}

func (m *Module) playerSetStickFloor(args []value.Value) (value.Value, error) {
	if m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.SETSTICKFLOOR: not available on this platform")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETSTICKFLOOR expects (entity, downDistance#)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.SETSTICKFLOOR: invalid entity")
	}
	down, ok := args[1].ToFloat()
	if !ok || down < 0 {
		return value.Nil, fmt.Errorf("PLAYER.SETSTICKFLOOR: downDistance must be >= 0")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SETSTICKFLOOR: call PLAYER.CREATE first")
	}
	if err := m.char.SetCharacterStickToFloorDown(ch, float32(down)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) playerGetStandNormal(args []value.Value) (value.Value, error) {
	if m.char == nil || m.h == nil || m.ent == nil {
		return value.Nil, fmt.Errorf("PLAYER.GETSTANDNORMAL: not available")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETSTANDNORMAL expects (entity)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETSTANDNORMAL: invalid entity")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.GETSTANDNORMAL: call PLAYER.CREATE first")
	}
	nx, ny, nz, ok := m.char.CharacterGroundNormal(ch)
	if ok {
		return mbmatrix.AllocVec3Value(m.h, float32(nx), float32(ny), float32(nz))
	}
	cx, cy, cz, ok := m.ent.PlayerBridgeWorldPos(id)
	if !ok {
		return mbmatrix.AllocVec3Value(m.h, 0, 1, 0)
	}
	nx2, ny2, nz2, hit := mbphysics3d.RaycastDownNormal(cx, cy+0.35, cz, 4.0)
	if !hit {
		return mbmatrix.AllocVec3Value(m.h, 0, 1, 0)
	}
	return mbmatrix.AllocVec3Value(m.h, float32(nx2), float32(ny2), float32(nz2))
}

func (m *Module) playerPush(args []value.Value) (value.Value, error) {
	if m.ent == nil || m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.PUSH: not available")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("PLAYER.PUSH expects (playerEntity, targetEntity, force)")
	}
	pid, ok1 := m.playerEntID(args[0])
	tid, ok2 := m.playerEntID(args[1])
	force, ok3 := args[2].ToFloat()
	if !ok1 || !ok2 || !ok3 || pid < 1 || tid < 1 {
		return value.Nil, fmt.Errorf("PLAYER.PUSH: invalid arguments")
	}
	_, _, _, dx, _, dz, ok := m.ent.PlayerBridgeEyeRay(pid, defaultEyeY)
	if !ok {
		return value.Nil, nil
	}
	flen := math.Hypot(dx, dz)
	if flen < 1e-9 {
		return value.Nil, nil
	}
	pm := 70.0
	if ch, ok := m.entToChar[pid]; ok {
		pm = m.char.CharacterMass(ch)
	}
	scale := force * (pm / 70.0)
	fx := float32(dx / flen * scale)
	fz := float32(dz / flen * scale)
	_ = m.ent.PlayerBridgeApplyForce(tid, fx, 0, fz)
	return value.Nil, nil
}

func (m *Module) playerGrab(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.GRAB expects (playerEntity, targetEntity) — use target 0 to release")
	}
	pid, ok1 := m.playerEntID(args[0])
	tid, ok2 := m.playerEntID(args[1])
	if !ok1 || !ok2 || pid < 1 {
		return value.Nil, fmt.Errorf("PLAYER.GRAB: invalid player entity")
	}
	if _, ok := m.entToChar[pid]; !ok {
		return value.Nil, fmt.Errorf("PLAYER.GRAB: call PLAYER.CREATE first")
	}
	if tid < 1 {
		delete(m.grab, pid)
		return value.Nil, nil
	}
	m.grab[pid] = tid
	return value.Nil, nil
}

func (m *Module) playerSetMass(args []value.Value) (value.Value, error) {
	if m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.SETMASS: not available on this platform")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETMASS expects (entity, mass)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.SETMASS: invalid entity")
	}
	mass, ok := args[1].ToFloat()
	if !ok || mass <= 0 {
		return value.Nil, fmt.Errorf("PLAYER.SETMASS: mass must be > 0")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SETMASS: call PLAYER.CREATE first")
	}
	if err := m.char.SetCharacterMass(ch, mass); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) playerGetSurfaceType(args []value.Value) (value.Value, error) {
	if m.h == nil || m.ent == nil {
		return value.Nil, fmt.Errorf("PLAYER.GETSURFACETYPE: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETSURFACETYPE expects (entity)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETSURFACETYPE: invalid entity")
	}
	cx, cy, cz, ok := m.ent.PlayerBridgeWorldPos(id)
	if !ok {
		return value.FromStringIndex(m.h.Intern("Default")), nil
	}
	hit := mbphysics3d.PickCastEntityID(cx, cy+0.25, cz, 0, -1, 0, 3.0)
	if hit <= 0 {
		return value.FromStringIndex(m.h.Intern("Default")), nil
	}
	s := m.ent.SurfaceMaterialHint(hit)
	return value.FromStringIndex(m.h.Intern(s)), nil
}

func (m *Module) playerSetFovKick(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETFOVKICK expects (entity, degrees)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.SETFOVKICK: invalid entity")
	}
	deg, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SETFOVKICK: degrees must be numeric")
	}
	m.fovKick[id] = deg
	return value.Nil, nil
}

func (m *Module) playerGetFovKick(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETFOVKICK expects (entity)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETFOVKICK: invalid entity")
	}
	return value.FromFloat(m.fovKick[id]), nil
}

func (m *Module) playerIsMoving(args []value.Value) (value.Value, error) {
	if m.char == nil {
		return value.FromBool(false), nil
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.ISMOVING expects (entity)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.ISMOVING: invalid entity")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.FromBool(false), nil
	}
	vx, _, vz, ok := m.char.CharacterLinearVelocity(ch)
	if !ok {
		return value.FromBool(false), nil
	}
	hs := math.Hypot(vx, vz)
	return value.FromBool(hs > 0.05), nil
}

func allocFloatArray(m *Module, ids []int64) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("heap not bound")
	}
	if len(ids) == 0 {
		return value.Nil, nil
	}
	arr, err := heap.NewArray([]int64{int64(len(ids))})
	if err != nil {
		return value.Nil, err
	}
	for i, id := range ids {
		_ = arr.Set([]int64{int64(i)}, float64(id))
	}
	h, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}
