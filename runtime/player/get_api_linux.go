//go:build linux && cgo

package player

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerPlayerCharGetAPI(m *Module, reg runtime.Registrar) {
	reg.Register("PLAYER.GETPOSITIONX", "player", m.playerGetPositionX)
	reg.Register("PLAYER.GETPOSITIONY", "player", m.playerGetPositionY)
	reg.Register("PLAYER.GETPOSITIONZ", "player", m.playerGetPositionZ)
	reg.Register("PLAYER.GETROTATIONPITCH", "player", m.playerGetRotationPitch)
	reg.Register("PLAYER.GETROTATIONYAW", "player", m.playerGetRotationYaw)
	reg.Register("PLAYER.GETROTATIONROLL", "player", m.playerGetRotationRoll)
	reg.Register("PLAYER.GETVELOCITYX", "player", m.playerGetVelocityX)
	reg.Register("PLAYER.GETVELOCITYY", "player", m.playerGetVelocityY)
	reg.Register("PLAYER.GETVELOCITYZ", "player", m.playerGetVelocityZ)
	reg.Register("PLAYER.GETSPEED", "player", m.playerGetSpeed)
	reg.Register("PLAYER.GETONSLOPE", "player", m.playerGetOnSlope)
	reg.Register("PLAYER.GETONWALL", "player", m.playerGetOnWall)
	reg.Register("PLAYER.GETSLOPEANGLE", "player", m.playerGetSlopeAngle)
	reg.Register("PLAYER.GETISJUMPING", "player", m.playerGetIsJumping)
	reg.Register("PLAYER.GETISFALLING", "player", m.playerGetIsFalling)
	reg.Register("PLAYER.GETMAXSLOPE", "player", m.playerGetMaxSlope)
	reg.Register("PLAYER.GETSTEPHEIGHT", "player", m.playerGetStepHeight)
	reg.Register("PLAYER.GETGRAVITYSCALE", "player", m.playerGetGravityScale)
	reg.Register("PLAYER.GETFRICTION", "player", m.playerGetFriction)
	reg.Register("PLAYER.GETSNAPDISTANCE", "player", m.playerGetSnapDistance)
	reg.Register("PLAYER.GETHEIGHT", "player", m.playerGetCapsuleHeight)
	reg.Register("PLAYER.GETRADIUS", "player", m.playerGetCapsuleRadius)
	reg.Register("PLAYER.GETLAYER", "player", m.playerGetLayerStub)
	reg.Register("PLAYER.GETMASK", "player", m.playerGetMaskStub)
	reg.Register("PLAYER.GETCOLLISIONENABLED", "player", m.playerGetCollisionEnabled)

	reg.Register("CHAR.GETPOSITIONX", "player", m.playerGetPositionX)
	reg.Register("CHAR.GETPOSITIONY", "player", m.playerGetPositionY)
	reg.Register("CHAR.GETPOSITIONZ", "player", m.playerGetPositionZ)
	reg.Register("CHAR.GETROTATIONPITCH", "player", m.playerGetRotationPitch)
	reg.Register("CHAR.GETROTATIONYAW", "player", m.playerGetRotationYaw)
	reg.Register("CHAR.GETROTATIONROLL", "player", m.playerGetRotationRoll)
	reg.Register("CHAR.GETVELOCITYX", "player", m.playerGetVelocityX)
	reg.Register("CHAR.GETVELOCITYY", "player", m.playerGetVelocityY)
	reg.Register("CHAR.GETVELOCITYZ", "player", m.playerGetVelocityZ)
	reg.Register("CHAR.GETSPEED", "player", m.playerGetSpeed)
	reg.Register("CHAR.GETONSLOPE", "player", m.playerGetOnSlope)
	reg.Register("CHAR.GETONWALL", "player", m.playerGetOnWall)
	reg.Register("CHAR.GETSLOPEANGLE", "player", m.playerGetSlopeAngle)
	reg.Register("CHAR.GETISJUMPING", "player", m.playerGetIsJumping)
	reg.Register("CHAR.GETISFALLING", "player", m.playerGetIsFalling)
	reg.Register("CHAR.GETMAXSLOPE", "player", m.playerGetMaxSlope)
	reg.Register("CHAR.GETSTEPHEIGHT", "player", m.playerGetStepHeight)
	reg.Register("CHAR.GETGRAVITYSCALE", "player", m.playerGetGravityScale)
	reg.Register("CHAR.GETFRICTION", "player", m.playerGetFriction)
	reg.Register("CHAR.GETSNAPDISTANCE", "player", m.playerGetSnapDistance)
	reg.Register("CHAR.GETHEIGHT", "player", m.playerGetCapsuleHeight)
	reg.Register("CHAR.GETRADIUS", "player", m.playerGetCapsuleRadius)
	reg.Register("CHAR.GETLAYER", "player", m.playerGetLayerStub)
	reg.Register("CHAR.GETMASK", "player", m.playerGetMaskStub)
	reg.Register("CHAR.GETCOLLISIONENABLED", "player", m.playerGetCollisionEnabled)

	registerPlayerKCCAliases(m, reg)
}

func (m *Module) playerGetPosAxis(id int64, axis int) (value.Value, error) {
	ch, ok := m.entToChar[id]
	if !ok || m.char == nil {
		return value.FromFloat(0), nil
	}
	x, y, z, ok := m.char.CharacterPosition(ch)
	if !ok {
		return value.FromFloat(0), nil
	}
	switch axis {
	case 0:
		return value.FromFloat(x), nil
	case 1:
		return value.FromFloat(y), nil
	default:
		return value.FromFloat(z), nil
	}
}

func (m *Module) playerGetPositionX(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETPOSITIONX: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	return m.playerGetPosAxis(id, 0)
}

func (m *Module) playerGetPositionY(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETPOSITIONY: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	return m.playerGetPosAxis(id, 1)
}

func (m *Module) playerGetPositionZ(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETPOSITIONZ: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	return m.playerGetPosAxis(id, 2)
}

func (m *Module) playerGetRotationPitch(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.ent == nil {
		return value.FromFloat(0), nil
	}
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETROTATIONPITCH: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	if id < 1 {
		return value.FromFloat(0), nil
	}
	p, _, _, ok := m.ent.WorldEulerForEntityID(id)
	if !ok {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(float64(p)), nil
}

func (m *Module) playerGetRotationYaw(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.ent == nil {
		return value.FromFloat(0), nil
	}
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETROTATIONYAW: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	if id < 1 {
		return value.FromFloat(0), nil
	}
	_, y, _, ok := m.ent.WorldEulerForEntityID(id)
	if !ok {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(float64(y)), nil
}

func (m *Module) playerGetRotationRoll(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.ent == nil {
		return value.FromFloat(0), nil
	}
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETROTATIONROLL: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	if id < 1 {
		return value.FromFloat(0), nil
	}
	_, _, r, ok := m.ent.WorldEulerForEntityID(id)
	if !ok {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(float64(r)), nil
}

func (m *Module) playerGetVelAxis(id int64, axis int) (value.Value, error) {
	ch, ok := m.entToChar[id]
	if !ok || m.char == nil {
		return value.FromFloat(0), nil
	}
	vx, vy, vz, ok := m.char.CharacterLinearVelocity(ch)
	if !ok {
		return value.FromFloat(0), nil
	}
	switch axis {
	case 0:
		return value.FromFloat(vx), nil
	case 1:
		return value.FromFloat(vy), nil
	default:
		return value.FromFloat(vz), nil
	}
}

func (m *Module) playerGetVelocityX(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETVELOCITYX: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	return m.playerGetVelAxis(id, 0)
}

func (m *Module) playerGetVelocityY(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETVELOCITYY: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	return m.playerGetVelAxis(id, 1)
}

func (m *Module) playerGetVelocityZ(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETVELOCITYZ: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	return m.playerGetVelAxis(id, 2)
}

func (m *Module) playerGetSpeed(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETSPEED: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	ch, ok := m.entToChar[id]
	if !ok || m.char == nil {
		return value.FromFloat(0), nil
	}
	vx, vy, vz, ok := m.char.CharacterLinearVelocity(ch)
	if !ok {
		return value.FromFloat(0), nil
	}
	sp := math.Sqrt(vx*vx + vy*vy + vz*vz)
	return value.FromFloat(sp), nil
}

func (m *Module) playerGetOnSlope(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.playerIsOnSteepSlope(rt, args...)
}

func (m *Module) playerGetOnWall(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETONWALL: %s", kccErrNoSubject)
		}
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
	// Jolt NotSupported: touching but not walkable floor (e.g. vertical surface)
	return value.FromBool(gi == 2), nil
}

func (m *Module) playerGetSlopeAngle(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETSLOPEANGLE: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	ch, ok := m.entToChar[id]
	if !ok || m.char == nil {
		return value.FromFloat(0), nil
	}
	ang := m.char.CharacterSlopeAngleDegrees(ch)
	return value.FromFloat(ang), nil
}

func (m *Module) playerGetIsJumping(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETISJUMPING: %s", kccErrNoSubject)
		}
		return value.FromBool(false), nil
	}
	ch, ok := m.entToChar[id]
	if !ok || m.char == nil {
		return value.FromBool(false), nil
	}
	g, err := m.char.CharacterIsGrounded(ch)
	if err != nil {
		return value.FromBool(false), nil
	}
	_, vy, _, ok := m.char.CharacterLinearVelocity(ch)
	if !ok {
		return value.FromBool(false), nil
	}
	return value.FromBool(!g && vy > 0.4), nil
}

func (m *Module) playerGetIsFalling(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETISFALLING: %s", kccErrNoSubject)
		}
		return value.FromBool(false), nil
	}
	ch, ok := m.entToChar[id]
	if !ok || m.char == nil {
		return value.FromBool(false), nil
	}
	g, err := m.char.CharacterIsGrounded(ch)
	if err != nil {
		return value.FromBool(false), nil
	}
	_, vy, _, ok := m.char.CharacterLinearVelocity(ch)
	if !ok {
		return value.FromBool(false), nil
	}
	return value.FromBool(!g && vy < -0.4), nil
}

func (m *Module) playerGetMaxSlope(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETMAXSLOPE: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	ch, ok := m.entToChar[id]
	if !ok || m.char == nil {
		return value.FromFloat(0), nil
	}
	if ms, ok := m.char.CharacterMaxSlopeDegrees(ch); ok {
		return value.FromFloat(ms), nil
	}
	return value.FromFloat(45), nil
}

func (m *Module) playerGetStepHeight(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETSTEPHEIGHT: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	ch, ok := m.entToChar[id]
	if !ok || m.char == nil {
		return value.FromFloat(0), nil
	}
	if h, ok := m.char.CharacterStepHeightY(ch); ok {
		return value.FromFloat(h), nil
	}
	return value.FromFloat(0), nil
}

func (m *Module) playerGetGravityScale(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETGRAVITYSCALE: %s", kccErrNoSubject)
		}
		return value.FromFloat(1), nil
	}
	ch, ok := m.entToChar[id]
	if !ok || m.char == nil {
		return value.FromFloat(1), nil
	}
	if g, ok := m.char.CharacterGravityScaleVal(ch); ok {
		return value.FromFloat(g), nil
	}
	return value.FromFloat(1), nil
}

func (m *Module) playerGetFriction(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETFRICTION: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	ch, ok := m.entToChar[id]
	if !ok || m.char == nil {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(m.char.CharacterGameplayFriction(ch)), nil
}

func (m *Module) playerGetSnapDistance(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETSNAPDISTANCE: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	ch, ok := m.entToChar[id]
	if !ok || m.char == nil {
		return value.FromFloat(0), nil
	}
	if d, ok := m.char.CharacterSnapDownDistance(ch); ok {
		return value.FromFloat(d), nil
	}
	return value.FromFloat(0), nil
}

func (m *Module) playerGetCapsuleHeight(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETHEIGHT: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	ch, ok := m.entToChar[id]
	if !ok || m.char == nil {
		return value.FromFloat(0), nil
	}
	_, fh, ok := m.char.CharacterCapsuleDims(ch)
	if !ok {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(fh), nil
}

func (m *Module) playerGetCapsuleRadius(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETRADIUS: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	ch, ok := m.entToChar[id]
	if !ok || m.char == nil {
		return value.FromFloat(0), nil
	}
	r, _, ok := m.char.CharacterCapsuleDims(ch)
	if !ok {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(r), nil
}

func (m *Module) playerGetLayerStub(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.FromInt(0), nil
}

func (m *Module) playerGetMaskStub(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.FromInt(0), nil
}

func (m *Module) playerGetCollisionEnabled(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.FromBool(true), nil
}
