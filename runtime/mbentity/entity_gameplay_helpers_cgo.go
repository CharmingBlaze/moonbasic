//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerEntityGameplayHelpersAPI(m *Module, r runtime.Registrar) {
	r.Register("ENTITY.NAVTO", "entity", runtime.AdaptLegacy(m.entNavTo))
	r.Register("ENT.NAVTO", "entity", runtime.AdaptLegacy(m.entNavTo))
	r.Register("ENTITY.DIST", "entity", runtime.AdaptLegacy(m.entDistance))
	r.Register("ENT.DIST", "entity", runtime.AdaptLegacy(m.entDistance))
	r.Register("CHAR.DIST", "entity", runtime.AdaptLegacy(m.entDistance))
	r.Register("ENTITY.SETHEALTH", "entity", runtime.AdaptLegacy(m.entSetHealth))
	r.Register("ENT.SET_HP", "entity", runtime.AdaptLegacy(m.entSetHealth))
	r.Register("ENT.SETHP", "entity", runtime.AdaptLegacy(m.entSetHealth))
	r.Register("ENTITY.DAMAGE", "entity", runtime.AdaptLegacy(m.entDamage))
	r.Register("ENT.DAMAGE", "entity", runtime.AdaptLegacy(m.entDamage))
	r.Register("ENTITY.ISALIVE", "entity", runtime.AdaptLegacy(m.entIsAlive))
	r.Register("ENT.ISALIVE", "entity", runtime.AdaptLegacy(m.entIsAlive))
	r.Register("ENT.SET_TEAM", "entity", runtime.AdaptLegacy(m.entSetTeam))
	r.Register("ENT.SETTEAM", "entity", runtime.AdaptLegacy(m.entSetTeam))
	r.Register("ENTITY.ONDEATHDROP", "entity", runtime.AdaptLegacy(m.entOnDeathDrop))
	r.Register("ENT.ONDEATH", "entity", runtime.AdaptLegacy(m.entOnDeath))
	r.Register("ENTITY.MAGNETTO", "entity", runtime.AdaptLegacy(m.entMagnetTo))
	r.Register("ENTITY.SETTAG", "entity", runtime.AdaptLegacy(m.entSetTag))
	r.Register("ENTITY.ADDWOBBLE", "entity", runtime.AdaptLegacy(m.entAddWobble))
	r.Register("ENT.WOBBLE", "entity", runtime.AdaptLegacy(m.entAddWobble))
	r.Register("ENTITY.ADDTRAIL", "entity", runtime.AdaptLegacy(m.entAddTrail))
	r.Register("ENTITY.WASGROUNDED", "entity", runtime.AdaptLegacy(m.entWasGrounded))
	r.Register("ENTITY.ISWALLSLIDING", "entity", runtime.AdaptLegacy(m.entIsWallSliding))
	r.Register("ENTITY.CUTJUMP", "entity", runtime.AdaptLegacy(m.entCutJump))
	r.Register("ENTITY.SETGRAVITYSCALE", "entity", runtime.AdaptLegacy(m.entSetGravityScaleEnt))
	r.Register("SPAWNER.MAKE", "entity", runtime.AdaptLegacy(m.spawnerMake))
	r.Register("ENT.SHOOT", "entity", runtime.AdaptLegacy(m.entShoot))
}

func (m *Module) entNavTo(args []value.Value) (value.Value, error) {
	if len(args) != 4 && len(args) != 5 && len(args) != 6 {
		return value.Nil, fmt.Errorf("ENTITY.NAVTO expects (entity, targetX, targetZ, speed# [, arrivalXZ# [, brakeDist#]])")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.NAVTO: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.NAVTO: unknown entity")
	}
	if e.static || e.physicsDriven {
		return value.Nil, fmt.Errorf("ENTITY.NAVTO: entity must be non-static and not physics-driven (Jolt body)")
	}
	tx, _ := args[1].ToFloat()
	tz, _ := args[2].ToFloat()
	spd, _ := args[3].ToFloat()
	if spd < 0 {
		return value.Nil, fmt.Errorf("ENTITY.NAVTO: speed must be non-negative")
	}
	ext := e.getExt()
	ext.patrolActive = false
	ext.navActive = true
	ext.navTX = float32(tx)
	ext.navTZ = float32(tz)
	ext.navSpeed = float32(spd)
	ext.navArrival = 0
	ext.navBrake = 0.75
	if len(args) >= 5 {
		ar, _ := args[4].ToFloat()
		if ar > 0 {
			ext.navArrival = float32(ar)
		}
	}
	if len(args) == 6 {
		br, _ := args[5].ToFloat()
		if br > 0 {
			ext.navBrake = float32(br)
		}
	}
	return m.chainEntityRef(args[0])
}

func (m *Module) entSetTeam(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENT.SET_TEAM expects (entity, teamId#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENT.SET_TEAM: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENT.SET_TEAM: unknown entity")
	}
	tid, ok := args[1].ToInt()
	if !ok {
		return value.Nil, fmt.Errorf("ENT.SET_TEAM: teamId must be numeric")
	}
	e.getExt().teamID = int32(tid)
	return value.Nil, nil
}

func (m *Module) entSetHealth(args []value.Value) (value.Value, error) {
	if len(args) != 2 && len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.SETHEALTH expects (entity, maxHealth#) or (entity, currentHealth#, maxHealth#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.SETHEALTH: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.SETHEALTH: unknown entity")
	}
	ext := e.getExt()
	if len(args) == 2 {
		mx, _ := args[1].ToFloat()
		if mx <= 0 {
			return value.Nil, fmt.Errorf("ENTITY.SETHEALTH: maxHealth must be positive")
		}
		ext.hpMax = float32(mx)
		ext.hpCur = ext.hpMax
	} else {
		cur, _ := args[1].ToFloat()
		mx, _ := args[2].ToFloat()
		if mx <= 0 {
			return value.Nil, fmt.Errorf("ENTITY.SETHEALTH: maxHealth must be positive")
		}
		ext.hpMax = float32(mx)
		ext.hpCur = float32(cur)
		if ext.hpCur < 0 {
			ext.hpCur = 0
		}
		if ext.hpCur > ext.hpMax {
			ext.hpCur = ext.hpMax
		}
	}
	return m.chainEntityRef(args[0])
}

func (m *Module) entDamage(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.DAMAGE expects (entity, amount#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.DAMAGE: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.DAMAGE: unknown entity")
	}
	ext := e.getExt()
	if ext.hpMax <= 0 {
		return value.Nil, fmt.Errorf("ENTITY.DAMAGE: call ENTITY.SETHEALTH first")
	}
	amt, _ := args[1].ToFloat()
	if amt < 0 {
		return value.Nil, fmt.Errorf("ENTITY.DAMAGE: amount must be non-negative")
	}
	ext.hpCur -= float32(amt)
	if ext.hpCur < 0 {
		ext.hpCur = 0
	}
	if amt > 0 {
		if ext.damageBlinkRemain <= 0 {
			ext.damageBlinkR0, ext.damageBlinkG0, ext.damageBlinkB0 = e.r, e.g, e.b
		}
		ext.damageBlinkRemain = 0.1
		e.r, e.g, e.b = 255, 48, 48
	}
	if ext.hpCur <= 0 && ext.deathDropPrefab >= 1 && ext.deathDropChance > 0 {
		if rand.Float64()*100 < float64(ext.deathDropChance) {
			wp := m.worldPos(e)
			v, err := m.entCopy([]value.Value{value.FromInt(ext.deathDropPrefab)})
			if err == nil && v.Kind == value.KindInt {
				nid, _ := v.ToInt()
				if ne := m.store().ents[nid]; ne != nil {
					m.setLocalFromWorld(ne, wp.X, wp.Y, wp.Z)
				}
			}
		}
	}
	return value.Nil, nil
}

func (m *Module) entIsAlive(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.ISALIVE expects (entity)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.ISALIVE: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.FromBool(false), nil
	}
	ext := e.getExt()
	if ext.hpMax <= 0 {
		return value.FromBool(true), nil
	}
	return value.FromBool(ext.hpCur > 0), nil
}

func (m *Module) resolvePrefabNameOrEntity(v value.Value) (int64, error) {
	if v.Kind == value.KindString {
		if m.h == nil {
			return 0, fmt.Errorf("heap not bound")
		}
		tag, ok := m.h.GetString(int32(v.IVal))
		if !ok {
			return 0, fmt.Errorf("invalid prefab name string")
		}
		id, ok2 := m.store().byName[strings.ToUpper(tag)]
		if !ok2 || id < 1 {
			return 0, fmt.Errorf("unknown prefab name %q (use ENTITY.SETNAME on a template entity)", tag)
		}
		return id, nil
	}
	id, ok := m.entID(v)
	if !ok || id < 1 {
		return 0, fmt.Errorf("invalid prefab entity")
	}
	return id, nil
}

func (m *Module) entOnDeath(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENT.ONDEATH expects (entity, prefabEntity# or prefabName$); use ENTITY.ONDEATHDROP for custom drop chance")
	}
	pid, err := m.resolvePrefabNameOrEntity(args[1])
	if err != nil {
		return value.Nil, err
	}
	return m.entOnDeathDrop([]value.Value{args[0], value.FromInt(pid), value.FromFloat(100)})
}

func (m *Module) entShoot(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENT.SHOOT expects (shooterEntity, prefabEntity# or prefabName$, speed#)")
	}
	sid, ok := m.entID(args[0])
	if !ok || sid < 1 {
		return value.Nil, fmt.Errorf("ENT.SHOOT: invalid shooter entity")
	}
	pid, err := m.resolvePrefabNameOrEntity(args[1])
	if err != nil {
		return value.Nil, err
	}
	spd, ok3 := args[2].ToFloat()
	if !ok3 || spd < 0 {
		return value.Nil, fmt.Errorf("ENT.SHOOT: speed must be non-negative")
	}
	sh := m.store().ents[sid]
	pref := m.store().ents[pid]
	if sh == nil || pref == nil {
		return value.Nil, fmt.Errorf("ENT.SHOOT: unknown entity")
	}
	v, err := m.entCopy([]value.Value{value.FromInt(pid)})
	if err != nil {
		return value.Nil, err
	}
	newID, ok4 := v.ToInt()
	if !ok4 || newID < 1 {
		return value.Nil, fmt.Errorf("ENT.SHOOT: copy failed")
	}
	bullet := m.store().ents[newID]
	if bullet == nil {
		return value.Nil, fmt.Errorf("ENT.SHOOT: internal error")
	}
	wp := m.worldPos(sh)
	p, w, _ := sh.getRot()
	fwd := forwardFromYawPitch(w, p)
	off := rl.Vector3Scale(fwd, 0.6)
	bullet.setPos(rl.Vector3Add(wp, off))
	f32 := float32(spd)
	bullet.vel = rl.Vector3{X: fwd.X * f32, Y: fwd.Y * f32, Z: fwd.Z * f32}
	bullet.static = false
	return value.FromInt(newID), nil
}

func (m *Module) entOnDeathDrop(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.ONDEATHDROP expects (entity, prefabEntity, chancePercent#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.ONDEATHDROP: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.ONDEATHDROP: unknown entity")
	}
	pid, ok2 := m.entID(args[1])
	if !ok2 || pid < 1 {
		return value.Nil, fmt.Errorf("ENTITY.ONDEATHDROP: invalid prefab entity")
	}
	if m.store().ents[pid] == nil {
		return value.Nil, fmt.Errorf("ENTITY.ONDEATHDROP: unknown prefab")
	}
	ch, _ := args[2].ToFloat()
	if ch < 0 || ch > 100 {
		return value.Nil, fmt.Errorf("ENTITY.ONDEATHDROP: chance must be 0..100")
	}
	ext := e.getExt()
	ext.deathDropPrefab = pid
	ext.deathDropChance = float32(ch)
	return m.chainEntityRef(args[0])
}

func (m *Module) entMagnetTo(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("ENTITY.MAGNETTO expects (entity, targetEntity, radius#, speed#)")
	}
	id, ok := m.entID(args[0])
	tid, ok2 := m.entID(args[1])
	if !ok || !ok2 || id < 1 || tid < 1 {
		return value.Nil, fmt.Errorf("ENTITY.MAGNETTO: invalid entity")
	}
	if m.store().ents[id] == nil || m.store().ents[tid] == nil {
		return value.Nil, fmt.Errorf("ENTITY.MAGNETTO: unknown entity")
	}
	rad, _ := args[2].ToFloat()
	spd, _ := args[3].ToFloat()
	if rad < 0 || spd < 0 {
		return value.Nil, fmt.Errorf("ENTITY.MAGNETTO: radius and speed must be non-negative")
	}
	ext := m.store().ents[id].getExt()
	ext.magnetActive = true
	ext.magnetTarget = tid
	ext.magnetRadius = float32(rad)
	ext.magnetSpeed = float32(spd)
	return m.chainEntityRef(args[0])
}

func (m *Module) entSetTag(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SETTAG expects (entity, tag$)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.SETTAG: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.SETTAG: unknown entity")
	}
	if args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.SETTAG: tag must be a string")
	}
	tag, ok := m.h.GetString(int32(args[1].IVal))
	if !ok {
		return value.Nil, fmt.Errorf("ENTITY.SETTAG: invalid tag string")
	}
	e.getExt().blenderTag = tag
	return m.chainEntityRef(args[0])
}

func (m *Module) entAddWobble(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.ADDWOBBLE expects (entity, height#, speed#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.ADDWOBBLE: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.ADDWOBBLE: unknown entity")
	}
	if e.physicsDriven {
		return value.Nil, fmt.Errorf("ENTITY.ADDWOBBLE: not for physics-driven entities")
	}
	h, _ := args[1].ToFloat()
	sp, _ := args[2].ToFloat()
	if h < 0 || sp < 0 {
		return value.Nil, fmt.Errorf("ENTITY.ADDWOBBLE: height and speed must be non-negative")
	}
	ext := e.getExt()
	ext.wobbleAmp = float32(h)
	ext.wobbleSpeed = float32(sp)
	ext.wobblePhase = 0
	ext.wobbleLastOff = 0
	return value.Nil, nil
}

func (m *Module) entAddTrail(args []value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("ENTITY.ADDTRAIL expects (entity, length#, r#, g#, b#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.ADDTRAIL: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.ADDTRAIL: unknown entity")
	}
	n, _ := args[1].ToInt()
	if n < 2 || n > 256 {
		return value.Nil, fmt.Errorf("ENTITY.ADDTRAIL: length must be 2..256")
	}
	rv, _ := args[2].ToFloat()
	gv, _ := args[3].ToFloat()
	bv, _ := args[4].ToFloat()
	ext := e.getExt()
	ext.trailCap = int(n)
	ext.trailSeg = make([]rl.Vector3, n)
	ext.trailHead = 0
	ext.trailCount = 0
	ext.trailR = uint8(rv)
	ext.trailG = uint8(gv)
	ext.trailB = uint8(bv)
	return value.Nil, nil
}

func (m *Module) entWasGrounded(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.WASGROUNDED expects (entity, graceSeconds#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.WASGROUNDED: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.FromBool(false), nil
	}
	grace, _ := args[1].ToFloat()
	if grace < 0 {
		return value.Nil, fmt.Errorf("ENTITY.WASGROUNDED: grace must be non-negative")
	}
	if e.onGround {
		return value.FromBool(true), nil
	}
	// Approximate "coyote time" using the same frame budget as ENTITY.UPDATE (groundCoyoteMax).
	if grace > 0 && e.groundCoyoteLeft > 0 {
		return value.FromBool(true), nil
	}
	return value.FromBool(false), nil
}

func (m *Module) entIsWallSliding(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.ISWALLSLIDING expects (entity)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.ISWALLSLIDING: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.FromBool(false), nil
	}
	ext := e.getExt()
	// Scripted sphere collisions: wall contact has mostly horizontal normals while falling.
	if e.onGround || e.vel.Y >= -0.01 {
		return value.FromBool(false), nil
	}
	if ext.hasHit && math.Abs(float64(ext.hitNY)) < 0.45 {
		return value.FromBool(true), nil
	}
	return value.FromBool(false), nil
}

func (m *Module) entCutJump(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.CUTJUMP expects (entity)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.CUTJUMP: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.CUTJUMP: unknown entity")
	}
	if e.vel.Y > 0.1 {
		e.vel.Y *= 0.5
	}
	return value.Nil, nil
}

func (m *Module) spawnerMake(args []value.Value) (value.Value, error) {
	if len(args) != 2 && len(args) != 4 {
		return value.Nil, fmt.Errorf("SPAWNER.MAKE expects (prefabEntity, intervalSec#) or (prefabEntity, intervalSec#, x#, z#)")
	}
	pid, ok := m.entID(args[0])
	if !ok || pid < 1 {
		return value.Nil, fmt.Errorf("SPAWNER.MAKE: invalid prefab entity")
	}
	if m.store().ents[pid] == nil {
		return value.Nil, fmt.Errorf("SPAWNER.MAKE: unknown prefab entity")
	}
	iv, _ := args[1].ToFloat()
	var x, z float64
	if len(args) == 4 {
		x, _ = args[2].ToFloat()
		z, _ = args[3].ToFloat()
	}
	if iv <= 0 {
		return value.Nil, fmt.Errorf("SPAWNER.MAKE: interval must be positive")
	}
	st := m.store()
	st.spawners = append(st.spawners, spawnerRec{
		prefabID: pid,
		interval: iv,
		remain:   iv,
		x:        float32(x),
		z:        float32(z),
	})
	return value.Nil, nil
}

func (m *Module) entSetGravityScaleEnt(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SETGRAVITYSCALE expects (entity, scale#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.SETGRAVITYSCALE: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.SETGRAVITYSCALE: unknown entity")
	}
	s, _ := args[1].ToFloat()
	if s < 0 {
		return value.Nil, fmt.Errorf("ENTITY.SETGRAVITYSCALE: scale must be non-negative")
	}
	e.gravScale = float32(s)
	return value.Nil, nil
}

func (m *Module) processDamageBlink(dt float32) {
	st := m.store()
	for _, e := range st.ents {
		if e == nil || e.ext == nil {
			continue
		}
		ext := e.ext
		if ext.damageBlinkRemain <= 0 {
			continue
		}
		ext.damageBlinkRemain -= dt
		if ext.damageBlinkRemain <= 0 {
			e.r, e.g, e.b = ext.damageBlinkR0, ext.damageBlinkG0, ext.damageBlinkB0
		}
	}
}

// --- per-frame processing from ENTITY.UPDATE ---

func (m *Module) processSpawners(dt float32) {
	st := m.store()
	dtf := float64(dt)
	for i := range st.spawners {
		s := &st.spawners[i]
		s.remain -= dtf
		if s.remain > 0 {
			continue
		}
		s.remain += s.interval
		v, err := m.entCopy([]value.Value{value.FromInt(s.prefabID)})
		if err != nil || v.Kind != value.KindInt {
			continue
		}
		nid, _ := v.ToInt()
		if ne := st.ents[nid]; ne != nil {
			wp := m.worldPos(ne)
			m.setLocalFromWorld(ne, s.x, wp.Y, s.z)
		}
	}
}

func (m *Module) processGameplayMotion(dt float32) {
	st := m.store()
	for _, e := range st.ents {
		if e == nil || e.ext == nil {
			continue
		}
		ext := e.ext
		if ext.navActive && !e.static && !e.physicsDriven {
			tx, tz := ext.navTX, ext.navTZ
			wp := m.worldPos(e)
			dx := tx - wp.X
			dz := tz - wp.Z
			dist := float32(math.Sqrt(float64(dx*dx + dz*dz)))
			arr := ext.navArrival
			if arr <= 0 {
				arr = 0.05
			}
			if dist < arr {
				ext.navActive = false
			} else {
				spd := ext.navSpeed
				br := ext.navBrake
				if br <= 0 {
					br = 0.75
				}
				if dist < br {
					t := dist / br
					spd *= t * t
				}
				step := spd * dt
				var nx, nz float32
				if step >= dist {
					nx, nz = tx, tz
					ext.navActive = false
				} else {
					nx = wp.X + dx/dist*step
					nz = wp.Z + dz/dist*step
				}
				m.setLocalFromWorld(e, nx, wp.Y, nz)
				nw := m.worldPos(e)
				yaw := math.Atan2(float64(tx-nw.X), float64(tz-nw.Z))
				p, _, r := e.getRot()
				e.setRot(p, float32(yaw), r)
			}
		}

		if ext.patrolActive && !e.static && !e.physicsDriven {
			var tx, tz float32
			if ext.patrolToB {
				tx, tz = ext.patrolBX, ext.patrolBZ
			} else {
				tx, tz = ext.patrolAX, ext.patrolAZ
			}
			wp := m.worldPos(e)
			dx := tx - wp.X
			dz := tz - wp.Z
			dist := float32(math.Sqrt(float64(dx*dx + dz*dz)))
			arr := float32(0.05)
			if dist < arr {
				ext.patrolToB = !ext.patrolToB
			} else {
				step := ext.patrolSpeed * dt
				var nx, nz float32
				if step >= dist {
					nx, nz = tx, tz
					ext.patrolToB = !ext.patrolToB
				} else {
					nx = wp.X + dx/dist*step
					nz = wp.Z + dz/dist*step
				}
				m.setLocalFromWorld(e, nx, wp.Y, nz)
				nw := m.worldPos(e)
				yaw := math.Atan2(float64(tx-nw.X), float64(tz-nw.Z))
				p, _, r := e.getRot()
				e.setRot(p, float32(yaw), r)
			}
		}

		if ext.magnetActive && !e.static && !e.physicsDriven {
			tgt := st.ents[ext.magnetTarget]
			if tgt == nil {
				ext.magnetActive = false
				continue
			}
			wp := m.worldPos(e)
			tw := m.worldPos(tgt)
			dx := tw.X - wp.X
			dy := tw.Y - wp.Y
			dz := tw.Z - wp.Z
			dist := float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
			if dist <= 0.01 {
				ext.magnetActive = false
				continue
			}
			if dist > ext.magnetRadius {
				continue
			}
			step := ext.magnetSpeed * dt
			if step >= dist {
				m.setLocalFromWorld(e, tw.X, tw.Y, tw.Z)
				ext.magnetActive = false
			} else {
				t := step / dist
				m.setLocalFromWorld(e, wp.X+dx*t, wp.Y+dy*t, wp.Z+dz*t)
			}
		}

	}
}

func (m *Module) processWobble(dt float32) {
	st := m.store()
	for _, e := range st.ents {
		if e == nil || e.ext == nil || e.physicsDriven {
			continue
		}
		ext := e.ext
		if ext.wobbleAmp <= 0 {
			continue
		}
		ext.wobblePhase += dt * ext.wobbleSpeed
		off := float32(math.Sin(float64(ext.wobblePhase))) * ext.wobbleAmp
		wp := m.worldPos(e)
		corr := off - ext.wobbleLastOff
		ext.wobbleLastOff = off
		m.setLocalFromWorld(e, wp.X, wp.Y+corr, wp.Z)
	}
}

func (m *Module) recordEntityTrails() {
	st := m.store()
	for _, e := range st.ents {
		if e == nil || e.ext == nil {
			continue
		}
		ext := e.ext
		if ext.trailCap < 2 {
			continue
		}
		wp := m.worldPos(e)
		ext.trailSeg[ext.trailHead] = wp
		ext.trailHead = (ext.trailHead + 1) % ext.trailCap
		if ext.trailCount < ext.trailCap {
			ext.trailCount++
		}
	}
}

func (m *Module) drawEntityTrail(e *ent) {
	ext := e.ext
	if ext == nil || ext.trailCap < 2 || ext.trailCount < 2 {
		return
	}
	col := rl.Color{R: ext.trailR, G: ext.trailG, B: ext.trailB, A: 255}
	n := ext.trailCount
	start := (ext.trailHead - n + ext.trailCap) % ext.trailCap
	var prev rl.Vector3
	first := true
	for i := 0; i < n; i++ {
		idx := (start + i) % ext.trailCap
		p := ext.trailSeg[idx]
		if !first {
			rl.DrawLine3D(prev, p, col)
		}
		first = false
		prev = p
	}
}
