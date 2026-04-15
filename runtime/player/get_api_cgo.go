//go:build (linux || windows) && cgo

package player

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerPlayerCharGetAPI(m *Module, reg runtime.Registrar) {
	reg.Register("PLAYER.GETPOSITIONX", "player", runtime.AdaptLegacy(m.playerGetPositionX))
	reg.Register("PLAYER.GETPOSITIONY", "player", runtime.AdaptLegacy(m.playerGetPositionY))
	reg.Register("PLAYER.GETPOSITIONZ", "player", runtime.AdaptLegacy(m.playerGetPositionZ))
	reg.Register("PLAYER.GETROTATIONPITCH", "player", runtime.AdaptLegacy(m.playerGetRotationPitch))
	reg.Register("PLAYER.GETROTATIONYAW", "player", runtime.AdaptLegacy(m.playerGetRotationYaw))
	reg.Register("PLAYER.GETROTATIONROLL", "player", runtime.AdaptLegacy(m.playerGetRotationRoll))
	reg.Register("PLAYER.GETVELOCITYX", "player", runtime.AdaptLegacy(m.playerGetVelocityX))
	reg.Register("PLAYER.GETVELOCITYY", "player", runtime.AdaptLegacy(m.playerGetVelocityY))
	reg.Register("PLAYER.GETVELOCITYZ", "player", runtime.AdaptLegacy(m.playerGetVelocityZ))
	reg.Register("PLAYER.GETSPEED", "player", runtime.AdaptLegacy(m.playerGetSpeed))
	reg.Register("PLAYER.GETONSLOPE", "player", runtime.AdaptLegacy(m.playerGetOnSlope))
	reg.Register("PLAYER.GETONWALL", "player", runtime.AdaptLegacy(m.playerGetOnWall))
	reg.Register("PLAYER.GETSLOPEANGLE", "player", runtime.AdaptLegacy(m.playerGetSlopeAngle))
	reg.Register("PLAYER.GETISJUMPING", "player", runtime.AdaptLegacy(m.playerGetIsJumping))
	reg.Register("PLAYER.GETISFALLING", "player", runtime.AdaptLegacy(m.playerGetIsFalling))
	reg.Register("PLAYER.GETMAXSLOPE", "player", runtime.AdaptLegacy(m.playerGetMaxSlope))
	reg.Register("PLAYER.GETSTEPHEIGHT", "player", runtime.AdaptLegacy(m.playerGetStepHeight))
	reg.Register("PLAYER.GETGRAVITYSCALE", "player", runtime.AdaptLegacy(m.playerGetGravityScale))
	reg.Register("PLAYER.GETFRICTION", "player", runtime.AdaptLegacy(m.playerGetFriction))
	reg.Register("PLAYER.GETSNAPDISTANCE", "player", runtime.AdaptLegacy(m.playerGetSnapDistance))
	reg.Register("PLAYER.GETISSLIDING", "player", runtime.AdaptLegacy(m.playerGetIsSliding))
	reg.Register("PLAYER.GETCEILING", "player", runtime.AdaptLegacy(m.playerGetCeiling))
	reg.Register("PLAYER.GETGROUNDVELOCITYX", "player", runtime.AdaptLegacy(m.playerGetGroundVelocityX))
	reg.Register("PLAYER.GETGROUNDVELOCITYY", "player", runtime.AdaptLegacy(m.playerGetGroundVelocityY))
	reg.Register("PLAYER.GETGROUNDVELOCITYZ", "player", runtime.AdaptLegacy(m.playerGetGroundVelocityZ))
	reg.Register("PLAYER.GETHEIGHT", "player", runtime.AdaptLegacy(m.playerGetCapsuleHeight))
	reg.Register("PLAYER.GETRADIUS", "player", runtime.AdaptLegacy(m.playerGetCapsuleRadius))
	reg.Register("PLAYER.GETLAYER", "player", runtime.AdaptLegacy(m.playerGetKCCObjectLayer))
	reg.Register("PLAYER.GETMASK", "player", runtime.AdaptLegacy(m.playerGetKCCCollisionMask))
	reg.Register("PLAYER.GETCOLLISIONENABLED", "player", runtime.AdaptLegacy(m.playerGetCollisionEnabled))

	reg.Register("CHAR.GETPOSITIONX", "player", runtime.AdaptLegacy(m.playerGetPositionX))
	reg.Register("CHAR.GETPOSITIONY", "player", runtime.AdaptLegacy(m.playerGetPositionY))
	reg.Register("CHAR.GETPOSITIONZ", "player", runtime.AdaptLegacy(m.playerGetPositionZ))
	reg.Register("CHAR.GETROTATIONPITCH", "player", runtime.AdaptLegacy(m.playerGetRotationPitch))
	reg.Register("CHAR.GETROTATIONYAW", "player", runtime.AdaptLegacy(m.playerGetRotationYaw))
	reg.Register("CHAR.GETROTATIONROLL", "player", runtime.AdaptLegacy(m.playerGetRotationRoll))
	reg.Register("CHAR.GETVELOCITYX", "player", runtime.AdaptLegacy(m.playerGetVelocityX))
	reg.Register("CHAR.GETVELOCITYY", "player", runtime.AdaptLegacy(m.playerGetVelocityY))
	reg.Register("CHAR.GETVELOCITYZ", "player", runtime.AdaptLegacy(m.playerGetVelocityZ))
	reg.Register("CHAR.GETSPEED", "player", runtime.AdaptLegacy(m.playerGetSpeed))
	reg.Register("CHAR.GETONSLOPE", "player", runtime.AdaptLegacy(m.playerGetOnSlope))
	reg.Register("CHAR.GETONWALL", "player", runtime.AdaptLegacy(m.playerGetOnWall))
	reg.Register("CHAR.GETSLOPEANGLE", "player", runtime.AdaptLegacy(m.playerGetSlopeAngle))
	reg.Register("CHAR.GETISJUMPING", "player", runtime.AdaptLegacy(m.playerGetIsJumping))
	reg.Register("CHAR.GETISFALLING", "player", runtime.AdaptLegacy(m.playerGetIsFalling))
	reg.Register("CHAR.GETMAXSLOPE", "player", runtime.AdaptLegacy(m.playerGetMaxSlope))
	reg.Register("CHAR.GETSTEPHEIGHT", "player", runtime.AdaptLegacy(m.playerGetStepHeight))
	reg.Register("CHAR.GETGRAVITYSCALE", "player", runtime.AdaptLegacy(m.playerGetGravityScale))
	reg.Register("CHAR.GETFRICTION", "player", runtime.AdaptLegacy(m.playerGetFriction))
	reg.Register("CHAR.GETSNAPDISTANCE", "player", runtime.AdaptLegacy(m.playerGetSnapDistance))
	reg.Register("CHAR.GETISSLIDING", "player", runtime.AdaptLegacy(m.playerGetIsSliding))
	reg.Register("CHAR.GETCEILING", "player", runtime.AdaptLegacy(m.playerGetCeiling))
	reg.Register("CHAR.GETGROUNDVELOCITYX", "player", runtime.AdaptLegacy(m.playerGetGroundVelocityX))
	reg.Register("CHAR.GETGROUNDVELOCITYY", "player", runtime.AdaptLegacy(m.playerGetGroundVelocityY))
	reg.Register("CHAR.GETGROUNDVELOCITYZ", "player", runtime.AdaptLegacy(m.playerGetGroundVelocityZ))
	reg.Register("CHAR.GETHEIGHT", "player", runtime.AdaptLegacy(m.playerGetCapsuleHeight))
	reg.Register("CHAR.GETRADIUS", "player", runtime.AdaptLegacy(m.playerGetCapsuleRadius))
	reg.Register("CHAR.GETLAYER", "player", runtime.AdaptLegacy(m.playerGetKCCObjectLayer))
	reg.Register("CHAR.GETMASK", "player", runtime.AdaptLegacy(m.playerGetKCCCollisionMask))
	reg.Register("CHAR.GETCOLLISIONENABLED", "player", runtime.AdaptLegacy(m.playerGetCollisionEnabled))

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

func (m *Module) playerGetPositionX(args []value.Value) (value.Value, error) {
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETPOSITIONX: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	return m.playerGetPosAxis(id, 0)
}

func (m *Module) playerGetPositionY(args []value.Value) (value.Value, error) {
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETPOSITIONY: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	return m.playerGetPosAxis(id, 1)
}

func (m *Module) playerGetPositionZ(args []value.Value) (value.Value, error) {
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETPOSITIONZ: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	return m.playerGetPosAxis(id, 2)
}

func (m *Module) playerGetRotationPitch(args []value.Value) (value.Value, error) {
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

func (m *Module) playerGetRotationYaw(args []value.Value) (value.Value, error) {
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

func (m *Module) playerGetRotationRoll(args []value.Value) (value.Value, error) {
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

func (m *Module) playerGetVelocityX(args []value.Value) (value.Value, error) {
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETVELOCITYX: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	return m.playerGetVelAxis(id, 0)
}

func (m *Module) playerGetVelocityY(args []value.Value) (value.Value, error) {
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETVELOCITYY: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	return m.playerGetVelAxis(id, 1)
}

func (m *Module) playerGetVelocityZ(args []value.Value) (value.Value, error) {
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETVELOCITYZ: %s", kccErrNoSubject)
		}
		return value.FromFloat(0), nil
	}
	return m.playerGetVelAxis(id, 2)
}

func (m *Module) playerGetSpeed(args []value.Value) (value.Value, error) {
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

func (m *Module) playerGetOnSlope(args []value.Value) (value.Value, error) {
	return m.playerIsOnSteepSlope(args)
}

func (m *Module) playerGetOnWall(args []value.Value) (value.Value, error) {
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

func (m *Module) playerGetSlopeAngle(args []value.Value) (value.Value, error) {
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

func (m *Module) playerGetIsJumping(args []value.Value) (value.Value, error) {
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

func (m *Module) playerGetIsFalling(args []value.Value) (value.Value, error) {
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

func (m *Module) playerGetMaxSlope(args []value.Value) (value.Value, error) {
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

func (m *Module) playerGetStepHeight(args []value.Value) (value.Value, error) {
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

func (m *Module) playerGetGravityScale(args []value.Value) (value.Value, error) {
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

func (m *Module) playerGetFriction(args []value.Value) (value.Value, error) {
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

func (m *Module) playerGetSnapDistance(args []value.Value) (value.Value, error) {
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

func (m *Module) playerGetCapsuleHeight(args []value.Value) (value.Value, error) {
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

func (m *Module) playerGetCapsuleRadius(args []value.Value) (value.Value, error) {
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

// KCC queries use Jolt object layer id 2 (CHARACTER) in physics_layers.h; mask includes all object layers (0..4 with ONE_WAY).
func (m *Module) playerGetKCCObjectLayer(args []value.Value) (value.Value, error) {
	_ = args
	return value.FromInt(2), nil
}

func (m *Module) playerGetKCCCollisionMask(args []value.Value) (value.Value, error) {
	_ = args
	return value.FromInt(31), nil
}

func (m *Module) playerGetIsSliding(args []value.Value) (value.Value, error) {
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETISSLIDING: %s", kccErrNoSubject)
		}
		return value.FromBool(false), nil
	}
	ch, ok := m.entToChar[id]
	if !ok || m.char == nil {
		return value.FromBool(false), nil
	}
	return value.FromBool(m.char.CharacterIsSliding(ch)), nil
}

func (m *Module) playerGetCeiling(args []value.Value) (value.Value, error) {
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETCEILING: %s", kccErrNoSubject)
		}
		return value.FromBool(false), nil
	}
	ch, ok := m.entToChar[id]
	if !ok || m.char == nil {
		return value.FromBool(false), nil
	}
	return value.FromBool(m.char.CharacterTouchingCeiling(ch)), nil
}

func (m *Module) playerGetGroundVelocityX(args []value.Value) (value.Value, error) {
	v, ok := m.playerGetGroundVelocityAxis(args, 0)
	if !ok {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(v), nil
}

func (m *Module) playerGetGroundVelocityY(args []value.Value) (value.Value, error) {
	v, ok := m.playerGetGroundVelocityAxis(args, 1)
	if !ok {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(v), nil
}

func (m *Module) playerGetGroundVelocityZ(args []value.Value) (value.Value, error) {
	v, ok := m.playerGetGroundVelocityAxis(args, 2)
	if !ok {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(v), nil
}

func (m *Module) playerGetGroundVelocityAxis(args []value.Value, axis int) (float64, bool) {
	id, ok := m.kccSubjectID(args)
	if !ok {
		return 0, false
	}
	ch, ok := m.entToChar[id]
	if !ok || m.char == nil {
		return 0, false
	}
	vx, vy, vz, ok := m.char.CharacterGroundVelocityVec(ch)
	if !ok {
		return 0, false
	}
	switch axis {
	case 0:
		return vx, true
	case 1:
		return vy, true
	default:
		return vz, true
	}
}

func (m *Module) playerGetCollisionEnabled(args []value.Value) (value.Value, error) {
	_ = args
	return value.FromBool(true), nil
}
