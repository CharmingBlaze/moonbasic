package player

import "moonbasic/vm/value"

// playerEntID resolves a script entity value (numeric id or EntityRef handle) to an internal entity id.
func (m *Module) playerEntID(v value.Value) (int64, bool) {
	if m.ent == nil {
		id, ok := v.ToInt()
		if !ok || id < 1 {
			return 0, false
		}
		return id, true
	}
	return m.ent.ResolveEntityID(v)
}
