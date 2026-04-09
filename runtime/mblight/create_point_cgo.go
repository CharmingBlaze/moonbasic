//go:build cgo || (windows && !cgo)

package mblight

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerPointLightBlitz(r runtime.Registrar) {
	r.Register("CreatePointLight", "light", m.createPointLightEntity)
}

// CreatePointLight(entity#, r, g, b) — point light that follows the entity each ENTITY.UPDATE (world position).
func (m *Module) createPointLightEntity(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CreatePointLight: heap not bound")
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CreatePointLight expects (entity#, r#, g#, b#)")
	}
	eid, ok := args[0].ToInt()
	if !ok || eid < 1 {
		return value.Nil, fmt.Errorf("CreatePointLight: entity# required")
	}
	rf, ok1 := argF32(args[1])
	gf, ok2 := argF32(args[2])
	bf, ok3 := argF32(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("CreatePointLight: r, g, b must be numeric")
	}
	o := newLightWithKind("point")
	o.parentEntID = eid
	if rf > 1 || gf > 1 || bf > 1 {
		o.r, o.g, o.b = rf/255, gf/255, bf/255
	} else {
		o.r, o.g, o.b = rf, gf, bf
	}
	o.colA = 1
	o.rangeDist = 25
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	o.self = id
	registerPointFollow(id)
	_ = rt
	return value.FromHandle(id), nil
}
