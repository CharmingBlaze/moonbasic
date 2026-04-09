package mblight

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// CreateLight builds a light from a numeric type: 1=directional, 2=point, 3=spot.
// Optional parentEntity# is stored for future scene attachment (not yet applied each frame).
func (m *Module) blitzCreateLight(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CreateLight: heap not bound")
	}
	if len(args) < 1 || len(args) > 2 {
		return value.Nil, fmt.Errorf("CreateLight expects (type# [, parentEntity#]) — 1=directional, 2=point, 3=spot")
	}
	ti, ok := args[0].ToInt()
	if !ok {
		return value.Nil, fmt.Errorf("CreateLight: type must be numeric")
	}
	var kind string
	switch ti {
	case 1:
		kind = "directional"
	case 2:
		kind = "point"
	case 3:
		kind = "spot"
	default:
		return value.Nil, fmt.Errorf("CreateLight: type must be 1–3")
	}
	o := newLightWithKind(kind)
	if len(args) == 2 {
		if pe, ok2 := args[1].ToInt(); ok2 && pe >= 1 {
			o.parentEntID = pe
		}
	}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	o.self = id
	return value.FromHandle(id), nil
}

func (m *Module) blitzLightRange(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.lightSetRange(rt, args...)
}

func (m *Module) blitzLightColor(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.lightSetColor(rt, args...)
}
