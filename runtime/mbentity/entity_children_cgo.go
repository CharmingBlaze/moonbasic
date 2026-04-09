//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func childLinkRemove(st *entityStore, parentID, childID int64) {
	if st == nil || parentID < 1 {
		return
	}
	lst := st.children[parentID]
	for i, id := range lst {
		if id == childID {
			st.children[parentID] = append(lst[:i], lst[i+1:]...)
			if len(st.children[parentID]) == 0 {
				delete(st.children, parentID)
			}
			return
		}
	}
}

func childLinkAdd(st *entityStore, parentID, childID int64) {
	if st == nil || parentID < 1 {
		return
	}
	st.children[parentID] = append(st.children[parentID], childID)
}

func (m *Module) entVisible(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.VISIBLE expects (entity#, visible)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	var vis bool
	switch args[1].Kind {
	case value.KindBool:
		vis = args[1].IVal != 0
	case value.KindInt:
		vis = args[1].IVal != 0
	default:
		return value.Nil, fmt.Errorf("visible must be bool or 0/1")
	}
	e.hidden = !vis
	return value.Nil, nil
}

func (m *Module) entCountChildren(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.COUNTCHILDREN expects (entity#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	if m.store().ents[id] == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	n := len(m.store().children[id])
	return value.FromInt(int64(n)), nil
}

func (m *Module) entGetChild(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.GETCHILD expects (entity#, index#)")
	}
	pid, ok := m.entID(args[0])
	if !ok || pid < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	if m.store().ents[pid] == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	idx, ok2 := args[1].ToInt()
	if !ok2 || idx < 0 {
		return value.Nil, fmt.Errorf("index must be non-negative int")
	}
	lst := m.store().children[pid]
	if int(idx) >= len(lst) {
		return value.Nil, fmt.Errorf("ENTITY.GETCHILD: index out of range")
	}
	return value.FromInt(lst[idx]), nil
}

func (m *Module) entFindChild(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.FINDCHILD expects (entity#, name$)")
	}
	rootID, ok := m.entID(args[0])
	if !ok || rootID < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	if m.store().ents[rootID] == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	if args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("name must be string")
	}
	want, ok2 := rt.Heap.GetString(int32(args[1].IVal))
	if !ok2 {
		return value.Nil, fmt.Errorf("invalid string")
	}
	want = strings.TrimSpace(want)
	st := m.store()
	q := append([]int64(nil), st.children[rootID]...)
	for len(q) > 0 {
		cid := q[0]
		q = q[1:]
		ce := st.ents[cid]
		if ce == nil {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(ce.name), want) {
			return value.FromInt(cid), nil
		}
		q = append(q, st.children[cid]...)
	}
	return value.FromInt(0), nil
}
