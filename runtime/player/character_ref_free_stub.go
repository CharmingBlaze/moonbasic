//go:build (!linux && !windows) || !cgo

package player

func charRefHeapObjFree(m *Module, id int64) {
	_ = m
	_ = id
}
