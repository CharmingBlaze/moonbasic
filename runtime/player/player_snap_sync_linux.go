//go:build linux && cgo

package player

func playerSnapSyncCharacter(m *Module, id int64, x, y, z float64) {
	if m.char == nil {
		return
	}
	ch, ok := m.entToChar[id]
	if !ok {
		return
	}
	_ = m.char.CharacterTeleport(ch, x, y, z)
}
