//go:build !linux || !cgo

package player

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	mbcamera "moonbasic/runtime/camera"
	mbphysics3d "moonbasic/runtime/physics3d"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerPlayerCommands(m *Module, reg runtime.Registrar) {
	reg.Register("CHARACTER.CREATE", "player", m.playerCharacterCreate)
	reg.Register("CHARACTERREF.ADDVELOCITY", "player", m.charRefAddVel)
	reg.Register("CHARACTERREF.SETLINEARVELOCITY", "player", m.charRefSetVel)
	reg.Register("CHARACTERREF.SETVELOCITY", "player", m.charRefSetVel)
	reg.Register("CHARACTERREF.SETSNAPDISTANCE", "player", m.charRefSetSnapDist)
	reg.Register("CHARACTERREF.SETSTICKDOWN", "player", m.charRefSetSnapDist) // Alias
	reg.Register("CHARACTERREF.UPDATE", "player", m.charRefUpdate)
	reg.Register("CHARACTERREF.UPDATEMOVE", "player", m.charRefUpdate)
	reg.Register("CHARACTERREF.JUMP", "player", m.charRefJump)
	reg.Register("CHARACTERREF.MOVEWITHCAMERA", "player", m.charRefMoveWithCam)
	reg.Register("CHARACTERREF.SETMAXSLOPE", "player", m.charRefSetMaxSlope)
	reg.Register("CHARACTERREF.SETSTEPHEIGHT", "player", m.charRefSetStepHeight)
	reg.Register("CHARACTERREF.ISGROUNDED", "player", m.charRefIsGrounded)
	reg.Register("CHARACTERREF.SETPOSITION", "player", m.charRefSetPos)
	reg.Register("CHARACTERREF.GETPOSITION", "player", m.charRefGetPos)
	reg.Register("CHARACTERREF.FREE", "player", m.charRefFree)
	reg.Register("CHARACTERREF.GETGROUNDSTATE", "player", m.charRefGetGroundState)
	reg.Register("PLAYER.GETGROUNDSTATE", "player", m.playerGetGroundState)
	reg.Register("PLAYER.ISONSTEEPSLOPE", "player", m.playerIsOnSteepSlope)
	reg.Register("CHAR.GETGROUNDSTATE", "player", m.playerGetGroundState)
	reg.Register("CHAR.ISONSTEEPSLOPE", "player", m.playerIsOnSteepSlope)

	reg.Register("CHARACTERREF.SETGRAVITY", "player", m.charRefSetGravityScale)
	reg.Register("CHARACTERREF.SETGRAVITYSCALE", "player", m.charRefSetGravityScale)
	reg.Register("CHARACTERREF.SETFRICTION", "player", m.charRefSetFriction)
	reg.Register("CHARACTERREF.SETPADDING", "player", m.charRefSetPadding)
	reg.Register("CHARACTERREF.SETBOUNCE", "player", m.charRefSetBounce)
	reg.Register("CHARACTERREF.GETSPEED", "player", m.charRefGetSpeed)
	reg.Register("CHARACTERREF.ISMOVING", "player", m.charRefIsMoving)

	registerPlayerCharGetAPI(m, reg)
	registerPlayerTerrainCommands(m, reg)
}

func (m *Module) playerCharacterCreate(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	// Polymorphism: Character.Create(x, y, z) vs Character.Create(entity, radius, height)
	if len(args) == 3 && args[0].Kind != value.KindHandle {
		_, isEnt := m.playerEntID(args[0])
		if !isEnt {
			// Standalone (pos X, Y, Z)
			return m.playerStandaloneCreate(args)
		}
	}

	v, err := m.playerCreate(rt, args...)
	if err != nil {
		return v, err
	}
	id := v.IVal
	m.lastHero = id
	obj := &charRefHeapObj{id: id, m: m}
	h, err := m.h.Alloc(obj)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func (m *Module) playerStandaloneCreate(args []value.Value) (value.Value, error) {
	x, _ := args[0].ToFloat(); y, _ := args[1].ToFloat(); z, _ := args[2].ToFloat()
	id := m.nextStandaloneID
	m.nextStandaloneID--

	m.hostKCC[id] = &hostKCCState{
		x:            x,
		y:            y,
		z:            z,
		rad:          0.5,
		hei:          2.0,
		stepH:        0.3,
		slopeDeg:     45.0,
		gravityScale: 1.0,
		pad:          0.02,
		vx: 0, vy: 0, vz: 0,
		grounded: false,
	}
	m.lastHero = id
	
	obj := &charRefHeapObj{id: id, m: m}
	h, err := m.h.Alloc(obj)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func (m *Module) playerCreate(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CHARACTER.CREATE/CHAR.MAKE expects (entity, radius#, height#) or (x, y, z)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("CHAR.MAKE: invalid entity")
	}
	rad, _ := args[1].ToFloat()
	hei, _ := args[2].ToFloat()
	if rad <= 0 || hei <= 0 {
		return value.Nil, fmt.Errorf("CHAR.MAKE: radius and height must be positive")
	}

	m.hostKCC[id] = &hostKCCState{
		rad:          rad,
		hei:          hei,
		stepH:        0.3,
		slopeDeg:     45.0,
		gravityScale: 1.0,
		pad:          0.02,
	}
	m.lastHero = id

	// Disable standard physics/gravity for the entity so the KCC can take over
	if m.ent != nil {
		m.ent.DisablePhysicsByID(int(id))
	}

	return value.FromInt(id), nil
}

func (m *Module) playerMove(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("PLAYER.MOVE expects (entity, velocityX, velocityZ) world units/sec")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.MOVE: invalid entity")
	}
	st, ok := m.hostKCC[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.MOVE: call PLAYER.CREATE / CHAR.MAKE first")
	}
	vx, _ := args[1].ToFloat()
	vz, _ := args[2].ToFloat()
	st.vx = vx
	st.vz = vz
	return value.Nil, nil
}

func (m *Module) playerCharMoveDir(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CHAR.MOVE expects (entity, dx#, dz#, speed#)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("CHAR.MOVE: invalid entity")
	}
	st, ok := m.hostKCC[id]
	if !ok {
		return value.Nil, fmt.Errorf("CHAR.MOVE: entity has no host KCC (call CHAR.MAKE first)")
	}

	dx, _ := args[1].ToFloat()
	dz, _ := args[2].ToFloat()
	spd, _ := args[3].ToFloat()

	st.vx = dx * spd
	st.vz = dz * spd

	return value.Nil, nil
}

func (m *Module) playerMoveWithCamera(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("CHAR.MOVEWITHCAM expects (entity, camera, fwd#, side#, speed#)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("CHAR.MOVEWITHCAMERA: invalid entity")
	}
	st, ok := m.hostKCC[id]
	if !ok {
		return value.Nil, fmt.Errorf("CHAR.MOVEWITHCAMERA: call CHAR.MAKE first")
	}

	camH := heap.Handle(args[1].IVal)
	fwd, _ := args[2].ToFloat()
	side, _ := args[3].ToFloat()
	spd, _ := args[4].ToFloat()

	fwdVec, sideVec, err := mbcamera.CameraXZWalkBasis(m.h, camH)
	if err != nil {
		return value.Nil, err
	}

	st.vx = (float64(fwdVec.X)*fwd + float64(sideVec.X)*side) * spd
	st.vz = (float64(fwdVec.Z)*fwd + float64(sideVec.Z)*side) * spd

	return value.Nil, nil
}

func (m *Module) playerNavChase(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("NAV.CHASE: Jolt kinematic character requires Linux+CGO")
}

func (m *Module) playerNavPatrol(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("NAV.PATROL: Jolt kinematic character requires Linux+CGO")
}

func (m *Module) playerNavTo(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) < 4 {
		return value.Nil, fmt.Errorf("NAV.GOTO expects (entity, tx#, tz#, speed# [, arrival#])")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("NAV.GOTO: invalid entity")
	}
	tx, _ := args[1].ToFloat()
	tz, _ := args[2].ToFloat()
	spd, _ := args[3].ToFloat()
	arr := 0.2
	if len(args) >= 5 {
		arrVal, _ := args[4].ToFloat()
		arr = arrVal
	}

	m.kccNav[id] = &kccNavState{
		active:  true,
		tx:      tx,
		tz:      tz,
		speed:   spd,
		arrival: arr,
	}

	return value.Nil, nil
}

func (m *Module) playerNavUpdate(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.NAVUPDATE expects (entity)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.NAVUPDATE: invalid entity")
	}
	// NAV intent is applied in Module.Process (UPDATEPHYSICS) via processNav + updateHostKCC (host_solver.go).
	return value.Nil, nil
}

func (m *Module) playerUpdate(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CHAR.UPDATE expects (dt#)")
	}
	dt, _ := args[0].ToFloat()
	m.Process(dt)
	return value.Nil, nil
}


func (m *Module) charRefAddVel(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("Character.AddVelocity expects (vx, vy, vz)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	st, ok := m.hostKCC[obj.id]
	if !ok {
		return value.Nil, fmt.Errorf("Character: host KCC state missing")
	}
	vx, _ := args[1].ToFloat()
	vy, _ := args[2].ToFloat()
	vz, _ := args[3].ToFloat()
	st.vx += vx
	st.vy += vy
	st.vz += vz
	return value.Nil, nil
}

func (m *Module) charRefSetVel(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("Character.SetLinearVelocity expects (vx, vy, vz)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	st, ok := m.hostKCC[obj.id]
	if !ok {
		return value.Nil, fmt.Errorf("Character: host KCC state missing")
	}
	vx, _ := args[1].ToFloat()
	vy, _ := args[2].ToFloat()
	vz, _ := args[3].ToFloat()
	st.vx, st.vy, st.vz = vx, vy, vz
	return value.Nil, nil
}

func (m *Module) charRefSetPos(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("Character.SetPos expects (x, y, z)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	z, _ := args[3].ToFloat()
	
	if obj.id < 0 {
		if st, ok := m.hostKCC[obj.id]; ok {
			st.x, st.y, st.z = x, y, z
			st.vx, st.vy, st.vz = 0, 0, 0
		}
	} else {
		if m.ent != nil {
			m.ent.PlayerBridgeSetWorldPos(obj.id, float32(x), float32(y), float32(z))
		}
		if st, ok := m.hostKCC[obj.id]; ok {
			st.vx, st.vy, st.vz = 0, 0, 0
		}
	}
	return value.Nil, nil
}

func (m *Module) charRefFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("Character.Free expects no args (handle receiver)")
	}
	h := heap.Handle(args[0].IVal)
	if err := m.h.Free(h); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) charRefJump(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("Character.Jump expects (force#)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	f, _ := args[1].ToFloat()
	if st, ok := m.hostKCC[obj.id]; ok {
		st.vy = f
		st.grounded = false
	}
	return value.Nil, nil
}

func (m *Module) charRefUpdate(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("Character.Update expects (dt)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	dt, _ := args[1].ToFloat()
	st, ok := m.hostKCC[obj.id]
	if ok {
		m.updateHostKCC(obj.id, st, dt)
	}
	return value.Nil, nil
}

func (m *Module) charRefMoveWithCam(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("Character.MoveWithCamera expects (cam_handle, fwd#, side#, speed#)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	st, ok := m.hostKCC[obj.id]
	if !ok {
		return value.Nil, fmt.Errorf("Character: host KCC state missing")
	}
	
	camH := heap.Handle(args[1].IVal)
	fwd, _ := args[2].ToFloat()
	side, _ := args[3].ToFloat()
	spd, _ := args[4].ToFloat()
	
	fwdVec, sideVec, err := mbcamera.CameraXZWalkBasis(m.h, camH)
	if err != nil {
		return value.Nil, err
	}
	
	st.vx = (float64(fwdVec.X)*fwd + float64(sideVec.X)*side) * spd
	st.vz = (float64(fwdVec.Z)*fwd + float64(sideVec.Z)*side) * spd
	
	return value.Nil, nil
}

func (m *Module) charRefSetMaxSlope(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	deg, _ := args[1].ToFloat()
	if st, ok := m.hostKCC[obj.id]; ok {
		st.slopeDeg = deg
	}
	return value.Nil, nil
}

func (m *Module) charRefSetStepHeight(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	h, _ := args[1].ToFloat()
	if st, ok := m.hostKCC[obj.id]; ok {
		st.stepH = h
	}
	return value.Nil, nil
}

func (m *Module) charRefIsGrounded(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if st, ok := m.hostKCC[obj.id]; ok {
		return value.FromBool(st.grounded), nil
	}
	return value.FromBool(false), nil
}

func (m *Module) charRefGetPos(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	arr, _ := heap.NewArray([]int64{3})
	x, _ := m.hostGetPosAxis(obj.id, 0)
	y, _ := m.hostGetPosAxis(obj.id, 1)
	z, _ := m.hostGetPosAxis(obj.id, 2)
	_ = arr.Set([]int64{0}, x.FVal)
	_ = arr.Set([]int64{1}, y.FVal)
	_ = arr.Set([]int64{2}, z.FVal)
	h, _ := m.h.Alloc(arr)
	return value.FromHandle(h), nil
}

func (m *Module) charRefSetSnapDist(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 { return value.Nil, fmt.Errorf("Character.SetSnapDistance expects (dist#)") }
	obj, _ := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	v, _ := args[1].ToFloat()
	if st, ok := m.hostKCC[obj.id]; ok { st.stickDown = v }
	return value.Nil, nil
}

func (m *Module) playerGetGroundState(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
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
	st, ok := m.hostKCC[id]
	if !ok {
		return value.FromInt(3), nil
	}
	if st.grounded {
		return value.FromInt(0), nil
	}
	return value.FromInt(3), nil
}

func (m *Module) playerIsOnSteepSlope(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) > 1 {
		return value.Nil, fmt.Errorf("PLAYER.ISONSTEEPSLOPE expects () or (entity)")
	}
	if _, ok := m.kccSubjectID(args); !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.ISONSTEEPSLOPE: %s", kccErrNoSubject)
		}
	}
	return value.FromBool(false), nil
}

func (m *Module) charRefGetGroundState(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CHARACTERREF.GETGROUNDSTATE expects (handle)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	st, ok := m.hostKCC[obj.id]
	if !ok {
		return value.FromInt(3), nil
	}
	if st.grounded {
		return value.FromInt(0), nil
	}
	return value.FromInt(3), nil
}

func (m *Module) playerSetPadding(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.SETPADDING: Jolt kinematic character requires Linux+CGO")
}

func (m *Module) playerJump(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.JUMP expects (entity, force#)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.JUMP: invalid entity")
	}
	st, ok := m.hostKCC[id]
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.JUMP: call CHAR.MAKE first")
	}
	force, _ := args[1].ToFloat()
	// Bug Fix: Remove strict grounding check for script-triggered jumps or double-jumps if needed.
	// But mostly: zero out grounded immediately so solver doesn't snap us back.
	st.vy = force
	st.grounded = false
	return value.Nil, nil
}

func (m *Module) playerIsGrounded(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
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
		if !ok {
			return value.FromBool(false), nil
		}
	}
	st, ok2 := m.hostKCC[id]
	if !ok2 {
		return value.FromBool(false), nil
	}
	return value.FromBool(st.grounded), nil
}

const defaultEyeY = 1.65

func (m *Module) playerGetLookTarget(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.ent == nil || m.h == nil {
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

func (m *Module) playerGetNearby(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.ent == nil || m.h == nil {
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
	return allocFloatArrayStub(m, ids)
}

func (m *Module) playerOnTrigger(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.ONTRIGGER: VM callback from physics not wired — use LEVEL.BINDSCRIPT + collision checks")
}

func (m *Module) playerSetState(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETSTATE expects (entity, state)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.SETSTATE: invalid entity")
	}
	st, ok := args[1].ToInt()
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SETSTATE: state must be numeric")
	}
	m.state[id] = int32(st)
	return value.Nil, nil
}

func (m *Module) playerSyncAnim(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.SYNCANIM: requires PLAYER.CREATE (Linux+Jolt)")
}

func (m *Module) playerSetStepHeight(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.SETSTEPHEIGHT: requires PLAYER.CREATE (Linux+Jolt)")
}

func (m *Module) playerSetSlopeLimit(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CHAR.SETSLOPE expects (entity, deg#)")
	}
	id, ok := m.playerEntID(args[0])
	if ok {
		d, _ := args[1].ToFloat()
		if st, ok2 := m.hostKCC[id]; ok2 {
			st.slopeDeg = d
		}
	}
	return value.Nil, nil
}

func (m *Module) playerGetVelocity(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.GETVELOCITY: requires PLAYER.CREATE (Linux+Jolt)")
}

func (m *Module) playerTeleport(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("PLAYER.TELEPORT expects (entity, x, y, z)")
	}
	id, ok := m.playerEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.TELEPORT: invalid entity")
	}
	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	z, _ := args[3].ToFloat()
	
	if m.ent != nil {
		m.ent.PlayerBridgeSetWorldPos(id, float32(x), float32(y), float32(z))
	}
	if st, ok := m.hostKCC[id]; ok {
		st.vx, st.vy, st.vz = 0, 0, 0
	}
	return value.Nil, nil
}

func (m *Module) playerSetGravityScale(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETGRAVITYSCALE expects (entity, scale#)")
	}
	id, ok := m.playerEntID(args[0])
	if ok {
		s, _ := args[1].ToFloat()
		if st, ok2 := m.hostKCC[id]; ok2 {
			st.gravityScale = s
		}
	}
	return value.Nil, nil
}

func (m *Module) playerGetCrouch(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETCROUCH expects (entity)")
	}
	return value.FromBool(false), nil
}

func (m *Module) playerSetCrouch(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.SETCROUCH: requires PLAYER.CREATE (Linux+Jolt)")
}

func (m *Module) playerSwim(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.SWIM: requires PLAYER.CREATE (Linux+Jolt)")
}

func (m *Module) playerSetStepOffset(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CHAR.SETSTEP expects (entity, height#)")
	}
	id, ok := m.playerEntID(args[0])
	if ok {
		h, _ := args[1].ToFloat()
		if st, ok2 := m.hostKCC[id]; ok2 {
			st.stepH = h
		}
	}
	return value.Nil, nil
}

func (m *Module) playerSetStickFloor(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CHAR.STICK expects (entity, dist#)")
	}
	id, ok := m.playerEntID(args[0])
	if ok {
		d, _ := args[1].ToFloat()
		if st, ok2 := m.hostKCC[id]; ok2 {
			st.stickDown = d
		}
	}
	return value.Nil, nil
}

func (m *Module) playerGetStandNormal(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.GETSTANDNORMAL: requires PLAYER.CREATE (Linux+Jolt)")
}

func (m *Module) playerPush(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.PUSH: requires PLAYER.CREATE (Linux+Jolt)")
}

func (m *Module) playerGrab(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.GRAB: requires PLAYER.CREATE (Linux+Jolt)")
}

func (m *Module) playerSetMass(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.SETMASS: requires PLAYER.CREATE (Linux+Jolt)")
}

func (m *Module) playerGetSurfaceType(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETSURFACETYPE expects (entity)")
	}
	_, ok := m.playerEntID(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.GETSURFACETYPE: invalid entity")
	}
	return value.FromStringIndex(m.h.Intern("Default")), nil
}

func (m *Module) playerSetFovKick(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, nil
}

func (m *Module) playerGetFovKick(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	return value.FromFloat(0), nil
}

func (m *Module) playerIsMoving(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	return value.FromBool(false), nil
}

func allocFloatArrayStub(m *Module, ids []int64) (value.Value, error) {
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
func (m *Module) charRefSetGravityScale(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 { return value.Nil, fmt.Errorf("Character.SetGravity expects (handle, scale#); got %d args", len(args)) }
	obj, _ := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	v, _ := args[1].ToFloat()
	if st, ok := m.hostKCC[obj.id]; ok { st.gravityScale = v }
	return value.Nil, nil
}
func (m *Module) charRefSetFriction(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, nil // Host KCC friction not yet analytic
}
func (m *Module) charRefSetPadding(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	var id int64
	var v float64
	if len(args) == 2 {
		obj, _ := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
		id = obj.id
		v, _ = args[1].ToFloat()
	} else if len(args) == 1 {
		id = m.lastHero
		v, _ = args[0].ToFloat()
	} else {
		return value.Nil, fmt.Errorf("Character.SetPadding expects (handle, pad#) or (pad#)")
	}
	if st, ok := m.hostKCC[id]; ok { st.pad = v }
	return value.Nil, nil
}
func (m *Module) charRefSetBounce(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, nil
}
func (m *Module) charRefGetSpeed(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 { return value.Nil, fmt.Errorf("Character.GetSpeed expects receiver only") }
	obj, _ := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if st, ok := m.hostKCC[obj.id]; ok {
		sp := math.Sqrt(st.vx*st.vx + st.vy*st.vy + st.vz*st.vz)
		return value.FromFloat(sp), nil
	}
	return value.FromFloat(0), nil
}
func (m *Module) charRefIsMoving(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 { return value.Nil, fmt.Errorf("Character.IsMoving expects receiver only") }
	obj, _ := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if st, ok := m.hostKCC[obj.id]; ok {
		sp := math.Sqrt(st.vx*st.vx + st.vz*st.vz)
		return value.FromBool(sp > 0.01), nil
	}
	return value.FromBool(false), nil
}
