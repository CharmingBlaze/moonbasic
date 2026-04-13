//go:build !linux

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
	reg.Register("CHARACTERREF.X", "player", m.playerGetPositionX)
	reg.Register("CHARACTERREF.Y", "player", m.playerGetPositionY)
	reg.Register("CHARACTERREF.Z", "player", m.playerGetPositionZ)
	reg.Register("CHARACTERREF.GETPOSITIONX", "player", m.playerGetPositionX)
	reg.Register("CHARACTERREF.GETPOSITIONY", "player", m.playerGetPositionY)
	reg.Register("CHARACTERREF.GETPOSITIONZ", "player", m.playerGetPositionZ)
	reg.Register("CHARACTERREF.VX", "player", m.playerGetVelocityX)
	reg.Register("CHARACTERREF.VY", "player", m.playerGetVelocityY)
	reg.Register("CHARACTERREF.VZ", "player", m.playerGetVelocityZ)

	registerPlayerKCCAliases(m, reg)
}

func (m *Module) hostGetPosAxis(id int64, axis int) (value.Value, error) {
	st, ok := m.hostKCC[id]
	if !ok {
		return value.FromFloat(0), nil
	}
	
	var val float64
	if id < 0 {
		// Standalone
		switch axis {
		case 0: val = st.x
		case 1: val = st.y
		default: val = st.z
		}
	} else {
		// Entity-bound
		if m.ent == nil { return value.FromFloat(0), nil }
		wp, ok := m.ent.GetWorldPosByID(int(id))
		if !ok { return value.FromFloat(0), nil }
		switch axis {
		case 0: val = float64(wp.X)
		case 1: val = float64(wp.Y)
		default: val = float64(wp.Z)
		}
	}
	return value.FromFloat(val), nil
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
	return m.hostGetPosAxis(id, 0)
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
	return m.hostGetPosAxis(id, 1)
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
	return m.hostGetPosAxis(id, 2)
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

func (m *Module) hostGetVelAxis(id int64, axis int) (value.Value, error) {
	st, ok := m.hostKCC[id]
	if !ok {
		return value.FromFloat(0), nil
	}
	switch axis {
	case 0:
		return value.FromFloat(st.vx), nil
	case 1:
		return value.FromFloat(st.vy), nil
	default:
		return value.FromFloat(st.vz), nil
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
	return m.hostGetVelAxis(id, 0)
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
	return m.hostGetVelAxis(id, 1)
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
	return m.hostGetVelAxis(id, 2)
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
	st, ok := m.hostKCC[id]
	if !ok {
		return value.FromFloat(0), nil
	}
	sp := math.Sqrt(st.vx*st.vx + st.vy*st.vy + st.vz*st.vz)
	return value.FromFloat(sp), nil
}

func (m *Module) playerGetOnSlope(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.playerIsOnSteepSlope(rt, args...)
}

func (m *Module) playerGetOnWall(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if _, ok := m.kccSubjectID(args); !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETONWALL: %s", kccErrNoSubject)
		}
	}
	return value.FromBool(false), nil
}

func (m *Module) playerGetSlopeAngle(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if _, ok := m.kccSubjectID(args); !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETSLOPEANGLE: %s", kccErrNoSubject)
		}
	}
	return value.FromFloat(0), nil
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
	st, ok := m.hostKCC[id]
	if !ok {
		return value.FromBool(false), nil
	}
	return value.FromBool(!st.grounded && st.vy > 0.4), nil
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
	st, ok := m.hostKCC[id]
	if !ok {
		return value.FromBool(false), nil
	}
	return value.FromBool(!st.grounded && st.vy < -0.4), nil
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
	st, ok := m.hostKCC[id]
	if !ok {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(st.slopeDeg), nil
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
	st, ok := m.hostKCC[id]
	if !ok {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(st.stepH), nil
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
	st, ok := m.hostKCC[id]
	if !ok {
		return value.FromFloat(1), nil
	}
	return value.FromFloat(st.gravityScale), nil
}

func (m *Module) playerGetFriction(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	id, ok := m.kccSubjectID(args)
	if !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETFRICTION: %s", kccErrNoSubject)
		}
		return value.FromFloat(0.5), nil
	}
	if _, ok := m.hostKCC[id]; !ok {
		return value.FromFloat(0.5), nil
	}
	return value.FromFloat(0.5), nil
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
	st, ok := m.hostKCC[id]
	if !ok {
		return value.FromFloat(0), nil
	}
	if st.stickDown > 0 {
		return value.FromFloat(st.stickDown), nil
	}
	return value.FromFloat(st.stepH + 0.2), nil
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
	st, ok := m.hostKCC[id]
	if !ok {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(st.hei), nil
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
	st, ok := m.hostKCC[id]
	if !ok {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(st.rad), nil
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
