//go:build !linux || !cgo

package player

import (
	"fmt"

	"moonbasic/runtime"
	mbcamera "moonbasic/runtime/camera"
	mbphysics3d "moonbasic/runtime/physics3d"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerPlayerCommands(m *Module, reg runtime.Registrar) {
	reg.Register("PLAYER.CREATE", "player", m.playerCreate)
	reg.Register("PLAYER.MOVE", "player", m.playerMove)
	reg.Register("PLAYER.JUMP", "player", m.playerJump)
	reg.Register("PLAYER.ISGROUNDED", "player", m.playerIsGrounded)
	reg.Register("PLAYER.GETLOOKTARGET", "player", m.playerGetLookTarget)
	reg.Register("PLAYER.GETNEARBY", "player", m.playerGetNearby)
	reg.Register("ENT.GET_NEAREST", "player", m.playerGetNearby)
	reg.Register("ENT.GETNEAREST", "player", m.playerGetNearby)
	reg.Register("PLAYER.ONTRIGGER", "player", m.playerOnTrigger)
	reg.Register("PLAYER.SETSTATE", "player", m.playerSetState)
	reg.Register("PLAYER.SYNCANIM", "player", m.playerSyncAnim)
	reg.Register("PLAYER.SETSTEPHEIGHT", "player", m.playerSetStepHeight)
	reg.Register("PLAYER.SETSLOPELIMIT", "player", m.playerSetSlopeLimit)
	reg.Register("PLAYER.GETVELOCITY", "player", m.playerGetVelocity)
	reg.Register("PLAYER.TELEPORT", "player", m.playerTeleport)
	reg.Register("PLAYER.SETGRAVITYSCALE", "player", m.playerSetGravityScale)
	reg.Register("PLAYER.GETCROUCH", "player", m.playerGetCrouch)
	reg.Register("PLAYER.SETCROUCH", "player", m.playerSetCrouch)
	reg.Register("PLAYER.SWIM", "player", m.playerSwim)
	reg.Register("PLAYER.SETSTEPOFFSET", "player", m.playerSetStepOffset)
	reg.Register("PLAYER.SETSTICKFLOOR", "player", m.playerSetStickFloor)
	reg.Register("PLAYER.NAVTO", "player", m.playerNavTo)
	reg.Register("PLAYER.NAVUPDATE", "player", m.playerNavUpdate)
	reg.Register("PLAYER.SETPADDING", "player", m.playerSetPadding)
	reg.Register("PLAYER.MOVEWITHCAMERA", "player", m.playerMoveWithCamera)
	reg.Register("NAV.GOTO", "player", m.playerNavTo)
	reg.Register("NAV.UPDATE", "player", m.playerNavUpdate)
	reg.Register("NAV.CHASE", "player", m.playerNavChase)
	reg.Register("NAV.PATROL", "player", m.playerNavPatrol)
	reg.Register("CHAR.MAKE", "player", m.playerCreate)
	reg.Register("CHAR.SETSTEP", "player", m.playerSetStepOffset)
	reg.Register("CHAR.SETSLOPE", "player", m.playerSetSlopeLimit)
	reg.Register("CHAR.SETPADDING", "player", m.playerSetPadding)
	reg.Register("CHAR.MOVE", "player", m.playerCharMoveDir)
	reg.Register("CHAR.MOVEWITHCAMERA", "player", m.playerMoveWithCamera)
	reg.Register("CHAR.MOVEWITHCAM", "player", m.playerMoveWithCamera)
	reg.Register("CHAR.NAVTO", "player", m.playerNavTo)
	reg.Register("CHAR.NAVUPDATE", "player", m.playerNavUpdate)
	reg.Register("CHAR.STICK", "player", m.playerSetStickFloor)
	reg.Register("CHAR.ISGROUNDED", "player", m.playerIsGrounded)
	reg.Register("CHAR.JUMP", "player", m.playerJump)
	reg.Register("PLAYER.GETSTANDNORMAL", "player", m.playerGetStandNormal)
	reg.Register("PLAYER.PUSH", "player", m.playerPush)
	reg.Register("PLAYER.GRAB", "player", m.playerGrab)
	reg.Register("PLAYER.SETMASS", "player", m.playerSetMass)
	reg.Register("PLAYER.GETSURFACETYPE", "player", m.playerGetSurfaceType)
	reg.Register("PLAYER.SETFOVKICK", "player", m.playerSetFovKick)
	reg.Register("PLAYER.GETFOVKICK", "player", m.playerGetFovKick)
	reg.Register("PLAYER.ISMOVING", "player", m.playerIsMoving)
	registerPlayerTerrainCommands(m, reg)
}

// hostEntID maps VM entity values (int or EntityRef handle from MODEL.*) to the internal entity store id.
func (m *Module) hostEntID(v value.Value) (int64, bool) {
	if m.ent == nil {
		return 0, false
	}
	return m.ent.ResolveEntityID(v)
}

func (m *Module) playerCreate(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CHAR.MAKE expects (entity, radius#, height#)")
	}
	id, ok := m.hostEntID(args[0])
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
		stepH:        0.3, // default
		slopeDeg:     45.0,
		gravityScale: 1.0,
		pad:          0.02,
	}

	// Disable standard physics/gravity for the entity so the KCC can take over
	if m.ent != nil {
		m.ent.DisablePhysicsByID(int(id))
	}

	return value.FromInt(id), nil
}

func (m *Module) playerMove(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.MOVE: Jolt kinematic character requires Linux+CGO")
}

func (m *Module) playerCharMoveDir(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CHAR.MOVE expects (entity, dx#, dz#, speed#)")
	}
	id, ok := m.hostEntID(args[0])
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
	id, ok := m.hostEntID(args[0])
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
	id, ok := m.hostEntID(args[0])
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
	id, ok := m.hostEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.NAVUPDATE: invalid entity")
	}
	// NAV intent is applied in Module.Process (UPDATEPHYSICS) via processNav + updateHostKCC (host_solver.go).
	return value.Nil, nil
}

func (m *Module) playerSetPadding(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.SETPADDING: Jolt kinematic character requires Linux+CGO")
}

func (m *Module) playerJump(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CHAR.JUMP expects (entity, force#)")
	}
	id, ok := m.hostEntID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("CHAR.JUMP: invalid entity")
	}
	st, ok := m.hostKCC[id]
	if !ok {
		return value.Nil, fmt.Errorf("CHAR.JUMP: call CHAR.MAKE first")
	}
	force, _ := args[1].ToFloat()
	if st.grounded {
		st.vy = force
		st.grounded = false
	}
	return value.Nil, nil
}

func (m *Module) playerIsGrounded(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 && len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.ISGROUNDED expects (entity) or (entity, coyoteTimeSec#)")
	}
	id, ok := m.hostEntID(args[0])
	if !ok || id < 1 {
		return value.FromBool(false), nil
	}
	st, ok := m.hostKCC[id]
	if !ok {
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
	id, ok := m.hostEntID(args[0])
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
	id, ok := m.hostEntID(args[0])
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
	id, ok := m.hostEntID(args[0])
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
	id, ok := m.hostEntID(args[0])
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
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.TELEPORT: requires PLAYER.CREATE (Linux+Jolt)")
}

func (m *Module) playerSetGravityScale(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.SETGRAVITYSCALE: requires PLAYER.CREATE (Linux+Jolt)")
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
	id, ok := m.hostEntID(args[0])
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
	id, ok := m.hostEntID(args[0])
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
	_, ok := m.hostEntID(args[0])
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
