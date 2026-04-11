//go:build cgo || (windows && !cgo)

package mbentity

// kinematicCharacterEntity reports entity ids that use Jolt CharacterVirtual (PLAYER.CREATE) so
// ENTITY.UPDATE does not apply scripted gravity/velocity to the same mesh.
var kinematicCharacterEntity func(int64) bool

// SetKinematicCharacterLookup registers a predicate from runtime/player (optional).
func SetKinematicCharacterLookup(fn func(int64) bool) {
	kinematicCharacterEntity = fn
}

func entityUsesKinematicCharacter(id int64) bool {
	if kinematicCharacterEntity == nil {
		return false
	}
	return kinematicCharacterEntity(id)
}
