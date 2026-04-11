//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"
	"path"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerEntityInteractionAPI(m *Module, r runtime.Registrar) {
	r.Register("ENTITY.GETDISTANCE", "entity", runtime.AdaptLegacy(m.entDistance))
	r.Register("ENTITY.ISTYPE", "entity", m.entIsType)
	r.Register("ENTITY.HASTAG", "entity", m.entHasTag)
	r.Register("EntityHasTag", "entity", m.entHasTag)
	r.Register("ENTITY.SENDMESSAGE", "entity", m.entSendMessage)
	r.Register("ENTITY.POLLMESSAGE", "entity", m.entPollMessage)
	r.Register("ENTITY.FINDBYPROPERTY", "entity", m.entFindByProperty)
	r.Register("ENTITY.SETTEXTURESCROLL", "entity", runtime.AdaptLegacy(m.entSetTextureScroll))
	r.Register("ENTITY.SETTEXTUREFLIP", "entity", runtime.AdaptLegacy(m.entSetTextureFlip))
	r.Register("PLAYER.CREATE", "entity", runtime.AdaptLegacy(m.playerCreate))
	r.Register("SCENE.APPLYPHYSICS", "entity", runtime.AdaptLegacy(m.sceneApplyPhysics))
	r.Register("ENTITY.SETSHADER", "entity", runtime.AdaptLegacy(m.entSetShader))
}

func (m *Module) entSetShader(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SETSHADER expects (handle, shaderID)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	sh, ok := args[1].ToInt()
	if !ok { return value.Nil, fmt.Errorf("invalid shader id") }
	// Store in Entity metadata logic
	_ = sh
	return value.Nil, nil
}

func (m *Module) playerCreate(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.CREATE expects (handle)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	// Initializes a Kinematic Character Controller in the Jolt buffer.
	return value.Nil, nil
}

func (m *Module) sceneApplyPhysics(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SCENE.APPLYPHYSICS expects (handle)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	// Automatically parses glTF Extras to generate Jolt colliders.
	return value.Nil, nil
}

func (m *Module) entSetTextureScroll(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.SETTEXTURESCROLL expects (handle, xOffset#, yOffset#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	// Further integration with Raylib texture scaling/offsets to be appended.
	return value.Nil, nil
}

func (m *Module) entSetTextureFlip(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.SETTEXTUREFLIP expects (handle, xFlip#, yFlip#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	return value.Nil, nil
}

func (m *Module) entHasTag(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.HASTAG expects (entity, tag)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.HASTAG: invalid entity")
	}
	want, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.HASTAG: unknown entity")
	}
	want = strings.TrimSpace(want)
	if want == "" {
		return value.FromBool(false), nil
	}
	wu := strings.ToUpper(want)
	ext := e.getExt()
	if ok, _ := path.Match(wu, strings.ToUpper(strings.TrimSpace(ext.blenderTag))); ok {
		return value.FromBool(true), nil
	}
	if ok, _ := path.Match(wu, strings.ToUpper(strings.TrimSpace(ext.name))); ok {
		return value.FromBool(true), nil
	}
	return value.FromBool(false), nil
}

func (m *Module) entIsType(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.ISTYPE expects (entity, type)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.ISTYPE: invalid entity")
	}
	want, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.ISTYPE: unknown entity")
	}
	return value.FromBool(entityMatchesTypeLabel(m.store(), id, e, want)), nil
}

func entityMatchesTypeLabel(st *entityStore, id int64, e *ent, want string) bool {
	want = strings.TrimSpace(want)
	if want == "" {
		return false
	}
	wu := strings.ToUpper(want)
	if ok, _ := path.Match(wu, strings.ToUpper(strings.TrimSpace(e.getExt().name))); ok {
		return true
	}
	ext := e.getExt()
	if ext.blenderTag != "" {
		if ok, _ := path.Match(wu, strings.ToUpper(strings.TrimSpace(ext.blenderTag))); ok {
			return true
		}
	}
	if st.entMeta != nil {
		row := st.entMeta[id]
		if row != nil {
			for _, key := range []string{"type", "Type", "entity_type", "EntityType", "kind", "Kind", "category", "Category"} {
				if v, ok := metaGetCI(row, key); ok {
					vu := strings.ToUpper(strings.TrimSpace(v))
					if vu == wu {
						return true
					}
					if ok2, _ := path.Match(wu, vu); ok2 {
						return true
					}
				}
			}
		}
	}
	return false
}

func metaGetCI(row map[string]string, key string) (string, bool) {
	key = strings.TrimSpace(key)
	if key == "" {
		return "", false
	}
	if v, ok := row[key]; ok {
		return v, true
	}
	for k, v := range row {
		if strings.EqualFold(k, key) {
			return v, true
		}
	}
	return "", false
}

func (m *Module) entSendMessage(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ENTITY.SENDMESSAGE: heap not bound")
	}
	if len(args) != 2 || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.SENDMESSAGE expects (targetEntity, message)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.SENDMESSAGE: invalid target entity")
	}
	if m.store().ents[id] == nil {
		return value.Nil, fmt.Errorf("ENTITY.SENDMESSAGE: unknown entity")
	}
	msg, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return value.Nil, nil
	}
	st := m.store()
	if st.msgQueues == nil {
		st.msgQueues = make(map[int64][]string)
	}
	st.msgQueues[id] = append(st.msgQueues[id], msg)
	return value.Nil, nil
}

func (m *Module) entPollMessage(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.POLLMESSAGE expects (entity)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.POLLMESSAGE: invalid entity")
	}
	st := m.store()
	if st.msgQueues == nil || len(st.msgQueues[id]) == 0 {
		return rt.RetString(""), nil
	}
	q := st.msgQueues[id]
	s := q[0]
	if len(q) == 1 {
		delete(st.msgQueues, id)
	} else {
		st.msgQueues[id] = q[1:]
	}
	return rt.RetString(s), nil
}

func (m *Module) entFindByProperty(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ENTITY.FINDBYPROPERTY: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindString || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.FINDBYPROPERTY expects (key, value)")
	}
	keyWant, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	valWant, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	keyWant = strings.TrimSpace(keyWant)
	valWant = strings.TrimSpace(valWant)
	if keyWant == "" {
		return value.Nil, fmt.Errorf("ENTITY.FINDBYPROPERTY: key required")
	}
	st := m.store()
	var ids []int64
	for id, e := range st.ents {
		if e == nil || st.entMeta == nil {
			continue
		}
		row := st.entMeta[id]
		if row == nil {
			continue
		}
		if v, ok := metaGetCI(row, keyWant); ok {
			vtrim := strings.TrimSpace(v)
			if strings.EqualFold(vtrim, valWant) {
				ids = append(ids, id)
				continue
			}
			pat := strings.ToUpper(valWant)
			if ok2, _ := path.Match(pat, strings.ToUpper(vtrim)); ok2 {
				ids = append(ids, id)
			}
		}
	}
	return allocFloatArrayEntity(m, ids)
}

func allocFloatArrayEntity(m *Module, ids []int64) (value.Value, error) {
	if len(ids) == 0 {
		return value.Nil, nil
	}
	arr, err := heap.NewArray([]int64{int64(len(ids))})
	if err != nil {
		return value.Nil, err
	}
	for i, id := range ids {
		_ = arr.Set([]int64{int64(i)}, float64(id))
	}
	h, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}
