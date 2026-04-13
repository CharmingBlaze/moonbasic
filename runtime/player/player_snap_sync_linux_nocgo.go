//go:build linux && !cgo

package player

func playerSnapSyncCharacter(m *Module, id int64, x, y, z float64) {
	if m.hostKCC != nil {
		if st, ok := m.hostKCC[id]; ok {
			st.vx, st.vy, st.vz = 0, 0, 0
			st.grounded = true
			return
		}
	}
	_ = id
	_ = x
	_ = y
	_ = z
}
