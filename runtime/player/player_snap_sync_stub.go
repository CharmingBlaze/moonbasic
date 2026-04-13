//go:build !linux

package player

func playerSnapSyncCharacter(m *Module, id int64, x, y, z float64) {
	if m.hostKCC != nil {
		if st, ok := m.hostKCC[id]; ok {
			st.vx, st.vy, st.vz = 0, 0, 0
			st.grounded = true
			if m.ent != nil {
				_ = m.ent.PlayerBridgeSetWorldPos(id, float32(x), float32(y), float32(z))
			}
			return
		}
	}
	_ = m
	_ = id
	_ = x
	_ = y
	_ = z
}
