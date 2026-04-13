//go:build !linux

package player

func charRefHeapObjFree(m *Module, id int64) {
	if m == nil || m.hostKCC == nil {
		return
	}
	delete(m.hostKCC, id)
}
