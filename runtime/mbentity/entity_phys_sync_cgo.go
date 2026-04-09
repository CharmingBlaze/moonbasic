//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"

	"moonbasic/runtime"
	mbphysics3d "moonbasic/runtime/physics3d"
	"moonbasic/vm/value"
)

// syncEntitiesFromPhysics copies world translation from the Jolt matrix buffer into linked entities
// (parent-aware local pose via setLocalFromWorld). Rotation in the buffer is not applied (jolt-go
// currently exposes position; buffer holds translation-only matrices from syncSharedBuffers).
func (m *Module) syncEntitiesFromPhysics() {
	buf := mbphysics3d.MatrixBufferForEntitySync()
	if len(buf) == 0 {
		return
	}
	st := m.store()
	for _, e := range st.ents {
		if e == nil || e.physBufIndex < 0 {
			continue
		}
		idx := e.physBufIndex * 16
		if idx+16 > len(buf) {
			continue
		}
		tx := buf[idx+12]
		ty := buf[idx+13]
		tz := buf[idx+14]
		m.setLocalFromWorld(e, tx, ty, tz)
	}
}

func (m *Module) entLinkPhysBuffer(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.LINKPHYSBUFFER expects 2 arguments (entity#, bufferIndex#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.LINKPHYSBUFFER: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.LINKPHYSBUFFER: unknown entity %d", id)
	}
	bi, ok := args[1].ToInt()
	if !ok || bi < 0 {
		return value.Nil, fmt.Errorf("ENTITY.LINKPHYSBUFFER: bufferIndex must be int >= 0")
	}
	e.physBufIndex = int(bi)
	mbphysics3d.RegisterEntityBufferLink(id, int(bi))
	return value.Nil, nil
}

func (m *Module) entClearPhysBuffer(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.CLEARPHYSBUFFER expects 1 argument (entity#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.CLEARPHYSBUFFER: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.CLEARPHYSBUFFER: unknown entity")
	}
	e.physBufIndex = -1
	mbphysics3d.UnregisterEntityCollision(id)
	return value.Nil, nil
}

func registerPhysicsEntitySync(m *Module, r runtime.Registrar) {
	r.Register("ENTITY.LINKPHYSBUFFER", "entity", runtime.AdaptLegacy(m.entLinkPhysBuffer))
	r.Register("ENTITY.CLEARPHYSBUFFER", "entity", runtime.AdaptLegacy(m.entClearPhysBuffer))
	mbphysics3d.SetAfterPhysicsMatrixSync(m.syncEntitiesFromPhysics)
	m.installPickLayerLookup()
}

func (m *Module) installPickLayerLookup() {
	mbphysics3d.SetPickLayerLookup(func(id int64) (uint8, bool) {
		if m.h == nil {
			return 0, false
		}
		st := entityStores[m.h]
		if st == nil {
			return 0, false
		}
		e := st.ents[id]
		if e == nil {
			return 0, false
		}
		return e.collisionLayer, true
	})
}
