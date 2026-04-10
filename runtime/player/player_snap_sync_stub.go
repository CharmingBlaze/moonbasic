//go:build !(linux && cgo)

package player

func playerSnapSyncCharacter(m *Module, id int64, x, y, z float64) {
	_ = m
	_ = id
	_ = x
	_ = y
	_ = z
}
