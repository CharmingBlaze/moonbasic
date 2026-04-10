//go:build linux && cgo

package player

import (
	"fmt"
	"math"

	mbentity "moonbasic/runtime/mbentity"
	mbmatrix "moonbasic/runtime/mbmatrix"
	mbphysics3d "moonbasic/runtime/physics3d"
	mbtime "moonbasic/runtime/time"
	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

const defaultEyeY = 1.65

const playerCapsuleRadius = 0.4
const playerCapsuleHeight = 1.75

func registerPlayerCommands(m *Module, reg runtime.Registrar) {
	reg.Register("PLAYER.CREATE", "player", m.playerCreate)
	reg.Register("PLAYER.MOVE", "player", m.playerMove)
	reg.Register("PLAYER.JUMP", "player", m.playerJump)
	reg.Register("PLAYER.ISGROUNDED", "player", m.playerIsGrounded)
	reg.Register("PLAYER.GETLOOKTARGET", "player", m.playerGetLookTarget)
	reg.Register("PLAYER.GETNEARBY", "player", m.playerGetNearby)
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

func (m *Module) playerCreate(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.h == nil || m.char == nil || m.ent == nil {
		return value.Nil, fmt.Errorf("PLAYER.CREATE: not available (requires Linux+Jolt fullruntime)")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.CREATE expects (entity#)")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.CREATE: invalid entity")
	}
	if _, dup := m.entToChar[id]; dup {
		return value.Nil, fmt.Errorf("PLAYER.CREATE: entity already has a character controller")
	}
	px, py, pz, ok := m.ent.PlayerBridgeWorldPos(id)
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.CREATE: unknown entity")
	}
	h, err := m.char.AllocCharacter(playerCapsuleRadius, playerCapsuleHeight, px, py, pz, 0)
	if err != nil {
		return value.Nil, err
	}
	m.entToChar[id] = h
	return value.Nil, nil
}

func (m *Module) playerMove(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.char == nil || m.ent == nil {
		return value.Nil, fmt.Errorf("PLAYER.MOVE: not available on this platform")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("PLAYER.MOVE expects (entity#, velocityX#, velocityZ#) world units/sec")
	}
	id, ok := args[0].ToInt()
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
	dt := mbtime.DeltaSeconds(rt)
	if dt <= 0 {
		dt = 1.0 / 60.0
	}
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

func (m *Module) playerJump(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.JUMP: not available on this platform")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.JUMP expects (entity#, impulseY#)")
	}
	id, ok := args[0].ToInt()
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

func (m *Module) playerIsGrounded(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.char == nil {
		return value.FromBool(false), nil
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.ISGROUNDED expects (entity#)")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.ISGROUNDED: invalid entity")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.FromBool(false), nil
	}
	g, err := m.char.CharacterIsGrounded(ch)
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(g), nil
}

func (m *Module) playerGetLookTarget(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.ent == nil {
		return value.FromInt(0), nil
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.GETLOOKTARGET expects (entity#, maxDist#)")
	}
	id, ok := args[0].ToInt()
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

func (m *Module) playerGetNearby(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.ent == nil {
		return value.Nil, nil
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("PLAYER.GETNEARBY expects (entity#, radius#, tag$)")
	}
	id, ok := args[0].ToInt()
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

func (m *Module) playerOnTrigger(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.ONTRIGGER expects (entity#, callbackFunc$)")
	}
	return value.Nil, fmt.Errorf("PLAYER.ONTRIGGER: VM callback from physics not wired — use LEVEL.BINDSCRIPT + collision checks or PHYSICS3D callbacks")
}

func (m *Module) playerSetState(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETSTATE expects (entity#, state#)")
	}
	id, ok := args[0].ToInt()
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

func (m *Module) playerSyncAnim(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.char == nil || m.ent == nil {
		return value.Nil, fmt.Errorf("PLAYER.SYNCANIM: not available on this platform")
	}
	if len(args) != 1 && len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SYNCANIM expects (entity# [, scale#])")
	}
	id, ok := args[0].ToInt()
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

func (m *Module) playerSetStepHeight(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETSTEPHEIGHT expects (entity#, height#)")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.SETSTEPHEIGHT: invalid entity")
	}
	h, ok := args[1].ToFloat()
	if !ok || h < 0 {
		return value.Nil, fmt.Errorf("PLAYER.SETSTEPHEIGHT: height must be >= 0")
	}
	if _, ok := m.entToChar[id]; !ok {
		return value.Nil, fmt.Errorf("PLAYER.SETSTEPHEIGHT: call PLAYER.CREATE first")
	}
	m.stepHeight[id] = h
	return value.Nil, nil
}

func (m *Module) playerSetSlopeLimit(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.SETSLOPELIMIT: not available on this platform")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETSLOPELIMIT expects (entity#, maxSlopeDegrees#)")
	}
	id, ok := args[0].ToInt()
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
	newH, err := m.char.RecreateCharacterWithSlope(ch, playerCapsuleRadius, playerCapsuleHeight, deg)
	if err != nil {
		return value.Nil, err
	}
	m.entToChar[id] = newH
	return value.Nil, nil
}

func (m *Module) playerGetVelocity(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.char == nil || m.h == nil {
		return value.Nil, fmt.Errorf("PLAYER.GETVELOCITY: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETVELOCITY expects (entity#)")
	}
	id, ok := args[0].ToInt()
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

func (m *Module) playerTeleport(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.char == nil || m.ent == nil {
		return value.Nil, fmt.Errorf("PLAYER.TELEPORT: not available on this platform")
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("PLAYER.TELEPORT expects (entity#, x#, y#, z#)")
	}
	id, ok := args[0].ToInt()
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

func (m *Module) playerSetGravityScale(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.SETGRAVITYSCALE: not available on this platform")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETGRAVITYSCALE expects (entity#, scale#)")
	}
	id, ok := args[0].ToInt()
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

func (m *Module) playerGetCrouch(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.char == nil {
		return value.FromBool(false), nil
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETCROUCH expects (entity#)")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETCROUCH: invalid entity")
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return value.FromBool(false), nil
	}
	return value.FromBool(m.char.CharacterCrouch(ch)), nil
}

func (m *Module) playerSetCrouch(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.SETCROUCH: not available on this platform")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETCROUCH expects (entity#, enabled)")
	}
	id, ok := args[0].ToInt()
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
	if err := m.char.SetCharacterCrouch(ch, en); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) playerSwim(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.SWIM: not available on this platform")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("PLAYER.SWIM expects (entity#, buoyancy#, drag#)")
	}
	id, ok := args[0].ToInt()
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
	if err := m.char.SetCharacterSwim(ch, buoy, drag, on); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) playerSetStepOffset(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.playerSetStepHeight(rt, args...)
}

func (m *Module) playerGetStandNormal(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.char == nil || m.h == nil || m.ent == nil {
		return value.Nil, fmt.Errorf("PLAYER.GETSTANDNORMAL: not available")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETSTANDNORMAL expects (entity#)")
	}
	id, ok := args[0].ToInt()
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

func (m *Module) playerPush(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.ent == nil || m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.PUSH: not available")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("PLAYER.PUSH expects (playerEntity#, targetEntity#, force#)")
	}
	pid, ok1 := args[0].ToInt()
	tid, ok2 := args[1].ToInt()
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

func (m *Module) playerGrab(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.GRAB expects (playerEntity#, targetEntity#) — use target 0 to release")
	}
	pid, ok1 := args[0].ToInt()
	tid, ok2 := args[1].ToInt()
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

func (m *Module) playerSetMass(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.char == nil {
		return value.Nil, fmt.Errorf("PLAYER.SETMASS: not available on this platform")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETMASS expects (entity#, mass#)")
	}
	id, ok := args[0].ToInt()
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

func (m *Module) playerGetSurfaceType(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.h == nil || m.ent == nil {
		return value.Nil, fmt.Errorf("PLAYER.GETSURFACETYPE: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETSURFACETYPE expects (entity#)")
	}
	id, ok := args[0].ToInt()
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

func (m *Module) playerSetFovKick(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.SETFOVKICK expects (entity#, degrees#)")
	}
	id, ok := args[0].ToInt()
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

func (m *Module) playerGetFovKick(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETFOVKICK expects (entity#)")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETFOVKICK: invalid entity")
	}
	return value.FromFloat(m.fovKick[id]), nil
}

func (m *Module) playerIsMoving(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.char == nil {
		return value.FromBool(false), nil
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.ISMOVING expects (entity#)")
	}
	id, ok := args[0].ToInt()
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
