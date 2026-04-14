//go:build cgo || (windows && !cgo)

package player

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// registerPlayerKCCAliases wires short Player-centric names (Player.GetX-style) to the same handlers as PLAYER.GETPOSITIONX, etc.
func registerPlayerKCCAliases(m *Module, reg runtime.Registrar) {
	reg.Register("PLAYER.GETX", "player", runtime.AdaptLegacy(m.playerGetPositionX))
	reg.Register("PLAYER.GETY", "player", runtime.AdaptLegacy(m.playerGetPositionY))
	reg.Register("PLAYER.GETZ", "player", runtime.AdaptLegacy(m.playerGetPositionZ))
	reg.Register("PLAYER.GETPITCH", "player", runtime.AdaptLegacy(m.playerGetRotationPitch))
	reg.Register("PLAYER.GETYAW", "player", runtime.AdaptLegacy(m.playerGetRotationYaw))
	reg.Register("PLAYER.GETROLL", "player", runtime.AdaptLegacy(m.playerGetRotationRoll))
	reg.Register("PLAYER.GETGROUNDED", "player", runtime.AdaptLegacy(m.playerIsGrounded))
	reg.Register("PLAYER.GETGRAVITY", "player", runtime.AdaptLegacy(m.playerGetGravityScale))
	reg.Register("PLAYER.GETCAPSULERADIUS", "player", runtime.AdaptLegacy(m.playerGetCapsuleRadius))
	reg.Register("PLAYER.GETCAPSULEHEIGHT", "player", runtime.AdaptLegacy(m.playerGetCapsuleHeight))
	reg.Register("PLAYER.GETSHAPETYPE", "player", runtime.AdaptLegacy(m.playerGetShapeType))

	reg.Register("CHAR.GETX", "player", runtime.AdaptLegacy(m.playerGetPositionX))
	reg.Register("CHAR.GETY", "player", runtime.AdaptLegacy(m.playerGetPositionY))
	reg.Register("CHAR.GETZ", "player", runtime.AdaptLegacy(m.playerGetPositionZ))
	reg.Register("CHAR.GETPITCH", "player", runtime.AdaptLegacy(m.playerGetRotationPitch))
	reg.Register("CHAR.GETYAW", "player", runtime.AdaptLegacy(m.playerGetRotationYaw))
	reg.Register("CHAR.GETROLL", "player", runtime.AdaptLegacy(m.playerGetRotationRoll))
	reg.Register("CHAR.GETGROUNDED", "player", runtime.AdaptLegacy(m.playerIsGrounded))
	reg.Register("CHAR.GETGRAVITY", "player", runtime.AdaptLegacy(m.playerGetGravityScale))
	reg.Register("CHAR.GETCAPSULERADIUS", "player", runtime.AdaptLegacy(m.playerGetCapsuleRadius))
	reg.Register("CHAR.GETCAPSULEHEIGHT", "player", runtime.AdaptLegacy(m.playerGetCapsuleHeight))
	reg.Register("CHAR.GETSHAPETYPE", "player", runtime.AdaptLegacy(m.playerGetShapeType))
}

// playerGetShapeType reports the KCC collision shape; CharacterVirtual is always a capsule today.
func (m *Module) playerGetShapeType(args []value.Value) (value.Value, error) {
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
