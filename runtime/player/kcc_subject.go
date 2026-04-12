package player

import "moonbasic/vm/value"

// kccErrNoSubject is returned when a PLAYER.* getter is called with no arguments
// but no implicit KCC subject was established (PLAYER.CREATE / Character.Create).
const kccErrNoSubject = "no KCC subject: pass (entity) or call PLAYER.CREATE / Character.Create first"

// kccSubjectID resolves the Kinematic Character subject for PLAYER.* commands that accept
// an optional (entity) argument.
//
//   - One argument: EntityRef, positive entity id, or (host only) negative virtual id for standalone Character.Create.
//   - Zero arguments: lastHero after PLAYER.CREATE / CHARACTER.CREATE, if that id has KCC state.
func (m *Module) kccSubjectID(args []value.Value) (int64, bool) {
	if len(args) >= 1 {
		v := args[0]
		if id, ok := v.ToInt(); ok && id < 0 {
			if m.hostKCC != nil {
				if _, ok := m.hostKCC[id]; ok {
					return id, true
				}
			}
			return 0, false
		}
		return m.playerEntID(v)
	}
	if m.lastHero == 0 {
		return 0, false
	}
	id := m.lastHero
	if m.entToChar != nil {
		if _, ok := m.entToChar[id]; ok {
			return id, true
		}
	}
	if m.hostKCC != nil {
		if _, ok := m.hostKCC[id]; ok {
			return id, true
		}
	}
	return 0, false
}
