//go:build (linux || windows) && cgo

package player

func charRefHeapObjFree(m *Module, id int64) {
	if m == nil {
		return
	}
	if m.char != nil {
		if h, ok := m.entToChar[id]; ok {
			_ = m.char.FreeCharacter(h)
			delete(m.entToChar, id)
		}
	}
}
