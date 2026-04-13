//go:build cgo || (windows && !cgo)

package player

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// registerPlayerKCCAliases wires short Player-centric names (Player.GetX-style) to the same handlers as PLAYER.GETPOSITIONX, etc.
func registerPlayerKCCAliases(m *Module, reg runtime.Registrar) {
	reg.Register("PLAYER.GETX", "player", m.playerGetPositionX)
	reg.Register("PLAYER.GETY", "player", m.playerGetPositionY)
	reg.Register("PLAYER.GETZ", "player", m.playerGetPositionZ)
	reg.Register("PLAYER.GETPITCH", "player", m.playerGetRotationPitch)
	reg.Register("PLAYER.GETYAW", "player", m.playerGetRotationYaw)
	reg.Register("PLAYER.GETROLL", "player", m.playerGetRotationRoll)
	reg.Register("PLAYER.GETGROUNDED", "player", m.playerIsGrounded)
	reg.Register("PLAYER.GETGRAVITY", "player", m.playerGetGravityScale)
	reg.Register("PLAYER.GETCAPSULERADIUS", "player", m.playerGetCapsuleRadius)
	reg.Register("PLAYER.GETCAPSULEHEIGHT", "player", m.playerGetCapsuleHeight)
	reg.Register("PLAYER.GETSHAPETYPE", "player", m.playerGetShapeType)

	reg.Register("CHAR.GETX", "player", m.playerGetPositionX)
	reg.Register("CHAR.GETY", "player", m.playerGetPositionY)
	reg.Register("CHAR.GETZ", "player", m.playerGetPositionZ)
	reg.Register("CHAR.GETPITCH", "player", m.playerGetRotationPitch)
	reg.Register("CHAR.GETYAW", "player", m.playerGetRotationYaw)
	reg.Register("CHAR.GETROLL", "player", m.playerGetRotationRoll)
	reg.Register("CHAR.GETGROUNDED", "player", m.playerIsGrounded)
	reg.Register("CHAR.GETGRAVITY", "player", m.playerGetGravityScale)
	reg.Register("CHAR.GETCAPSULERADIUS", "player", m.playerGetCapsuleRadius)
	reg.Register("CHAR.GETCAPSULEHEIGHT", "player", m.playerGetCapsuleHeight)
	reg.Register("CHAR.GETSHAPETYPE", "player", m.playerGetShapeType)
}

// playerGetShapeType reports the KCC collision shape; CharacterVirtual is always a capsule today.
func (m *Module) playerGetShapeType(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if _, ok := m.kccSubjectID(args); !ok {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("PLAYER.GETSHAPETYPE: %s", kccErrNoSubject)
		}
	}
	if m.h == nil {
		return value.Nil, fmt.Errorf("PLAYER.GETSHAPETYPE: heap not bound")
	}
	return value.FromStringIndex(m.h.Intern("capsule")), nil
}
