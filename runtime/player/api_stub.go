//go:build !linux || !cgo

package player

import (
	"fmt"

	mbphysics3d "moonbasic/runtime/physics3d"
	"moonbasic/runtime"
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
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.CREATE: Jolt kinematic character requires Linux+CGO (same as CHARCONTROLLER.*)")
}

func (m *Module) playerMove(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.MOVE: Jolt kinematic character requires Linux+CGO")
}

func (m *Module) playerJump(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.JUMP: Jolt kinematic character requires Linux+CGO")
}

func (m *Module) playerIsGrounded(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.ISGROUNDED expects (entity#)")
	}
	return value.FromBool(false), nil
}

const defaultEyeY = 1.65

func (m *Module) playerGetLookTarget(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.ent == nil || m.h == nil {
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
	return allocFloatArrayStub(m, ids)
}

func (m *Module) playerOnTrigger(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLAYER.ONTRIGGER expects (entity#, callbackFunc$)")
	}
	return value.Nil, fmt.Errorf("PLAYER.ONTRIGGER: VM callback from physics not wired — use LEVEL.BINDSCRIPT + collision checks")
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
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.SETSLOPELIMIT: requires PLAYER.CREATE (Linux+Jolt)")
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
		return value.Nil, fmt.Errorf("PLAYER.GETCROUCH expects (entity#)")
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
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("PLAYER.SETSTEPOFFSET: requires PLAYER.CREATE (Linux+Jolt)")
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
	_ = rt
	if m.h == nil {
		return value.Nil, fmt.Errorf("heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.GETSURFACETYPE expects (entity#)")
	}
	_, ok := args[0].ToInt()
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
