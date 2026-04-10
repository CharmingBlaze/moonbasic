package mblight

import (
	"fmt"
	"math"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func normalizeDir3(x, y, z float32) (float32, float32, float32) {
	l := float32(math.Sqrt(float64(x*x + y*y + z*z)))
	if l < 1e-8 {
		return 0, -1, 0
	}
	return x / l, y / l, z / l
}

func argF32(v value.Value) (float32, bool) {
	if f, ok := v.ToFloat(); ok {
		return float32(f), true
	}
	if i, ok := v.ToInt(); ok {
		return float32(i), true
	}
	return 0, false
}

func argBoolLight(v value.Value) (bool, bool) {
	if v.Kind == value.KindBool {
		return v.IVal != 0, true
	}
	if i, ok := v.ToInt(); ok {
		return i != 0, true
	}
	if f, ok := v.ToFloat(); ok {
		return f != 0, true
	}
	return false, false
}

func newLightWithKind(kind string) *lightObj {
	return &lightObj{
		kind:         kind,
		r:            1,
		g:            1,
		b:            1,
		colA:         1,
		intensity:    1,
		dirX:         -0.4082483,
		dirY:         -0.8164966,
		dirZ:         -0.4082483,
		targetX:      0,
		targetY:      2,
		targetZ:      0,
		shadowBiasK:  1,
		innerConeDeg: 25,
		outerConeDeg: 35,
		rangeDist:    10,
		enabled:      true,
	}
}

func (m *Module) lightCreatePoint(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LIGHT.CREATEPOINT: heap not bound")
	}
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("LIGHT.CREATEPOINT expects (x#, y#, z#, r#, g#, b#, energy#)")
	}
	x, ok1 := argF32(args[0])
	y, ok2 := argF32(args[1])
	z, ok3 := argF32(args[2])
	rf, ok4 := argF32(args[3])
	gf, ok5 := argF32(args[4])
	bf, ok6 := argF32(args[5])
	en, ok7 := argF32(args[6])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 {
		return value.Nil, fmt.Errorf("LIGHT.CREATEPOINT: arguments must be numeric")
	}
	o := newLightWithKind("point")
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	o.self = id
	o.posX, o.posY, o.posZ = x, y, z
	if rf > 1 || gf > 1 || bf > 1 {
		o.r = rf / 255
		o.g = gf / 255
		o.b = bf / 255
	} else {
		o.r, o.g, o.b = rf, gf, bf
	}
	o.colA = 1
	o.intensity = en
	if o.intensity < 0 {
		o.intensity = 0
	}
	return value.FromHandle(id), nil
}

func (m *Module) lightCreateDirectional(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LIGHT.CREATEDIRECTIONAL: heap not bound")
	}
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("LIGHT.CREATEDIRECTIONAL expects (dx#, dy#, dz#, r#, g#, b#, energy#)")
	}
	dx, ok1 := argF32(args[0])
	dy, ok2 := argF32(args[1])
	dz, ok3 := argF32(args[2])
	rf, ok4 := argF32(args[3])
	gf, ok5 := argF32(args[4])
	bf, ok6 := argF32(args[5])
	en, ok7 := argF32(args[6])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 {
		return value.Nil, fmt.Errorf("LIGHT.CREATEDIRECTIONAL: arguments must be numeric")
	}
	o := newLightWithKind("directional")
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	o.self = id
	o.dirX, o.dirY, o.dirZ = normalizeDir3(dx, dy, dz)
	if rf > 1 || gf > 1 || bf > 1 {
		o.r = rf / 255
		o.g = gf / 255
		o.b = bf / 255
	} else {
		o.r, o.g, o.b = rf, gf, bf
	}
	o.colA = 1
	o.intensity = en
	if o.intensity < 0 {
		o.intensity = 0
	}
	return value.FromHandle(id), nil
}

func (m *Module) lightCreateSpot(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LIGHT.CREATESPOT: heap not bound")
	}
	if len(args) != 11 {
		return value.Nil, fmt.Errorf("LIGHT.CREATESPOT expects (x#, y#, z#, tx#, ty#, tz#, r#, g#, b#, outerConeDeg#, energy#)")
	}
	x, ok1 := argF32(args[0])
	y, ok2 := argF32(args[1])
	z, ok3 := argF32(args[2])
	tx, ok4 := argF32(args[3])
	ty, ok5 := argF32(args[4])
	tz, ok6 := argF32(args[5])
	rf, ok7 := argF32(args[6])
	gf, ok8 := argF32(args[7])
	bf, ok9 := argF32(args[8])
	cone, ok10 := argF32(args[9])
	en, ok11 := argF32(args[10])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 || !ok9 || !ok10 || !ok11 {
		return value.Nil, fmt.Errorf("LIGHT.CREATESPOT: arguments must be numeric")
	}
	o := newLightWithKind("spot")
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	o.self = id
	o.posX, o.posY, o.posZ = x, y, z
	o.targetX, o.targetY, o.targetZ = tx, ty, tz
	ddx := tx - x
	ddy := ty - y
	ddz := tz - z
	o.dirX, o.dirY, o.dirZ = normalizeDir3(ddx, ddy, ddz)
	if cone < 0.5 {
		cone = 0.5
	}
	if cone > 89 {
		cone = 89
	}
	o.outerConeDeg = cone
	o.innerConeDeg = cone * 0.85
	if o.innerConeDeg < 0.5 {
		o.innerConeDeg = 0.5
	}
	if rf > 1 || gf > 1 || bf > 1 {
		o.r = rf / 255
		o.g = gf / 255
		o.b = bf / 255
	} else {
		o.r, o.g, o.b = rf, gf, bf
	}
	o.colA = 1
	o.intensity = en
	if o.intensity < 0 {
		o.intensity = 0
	}
	return value.FromHandle(id), nil
}

func (m *Module) lightMake(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LIGHT.MAKE: heap not bound")
	}
	kind := "directional"
	if len(args) > 1 {
		return value.Nil, fmt.Errorf("LIGHT.MAKE expects 0 arguments (default directional, white, intensity 1) or 1 kind string")
	}
	if len(args) == 1 {
		if args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("LIGHT.MAKE: optional argument must be a string kind (e.g. \"directional\")")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		kind = strings.ToLower(strings.TrimSpace(s))
	}
	o := newLightWithKind(kind)
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	o.self = id
	return value.FromHandle(id), nil
}

func (m *Module) lightFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LIGHT.FREE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("LIGHT.FREE expects (lightHandle)")
	}
	h := heap.Handle(args[0].IVal)
	return value.Nil, m.h.Free(h)
}

func (m *Module) lightSetDir(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LIGHT.SETDIR: heap not bound")
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("LIGHT.SETDIR expects (light, x, y, z)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("LIGHT.SETDIR: light must be a handle")
	}
	o, err := heap.Cast[*lightObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argF32(args[1])
	y, ok2 := argF32(args[2])
	z, ok3 := argF32(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("LIGHT.SETDIR: x, y, z must be numeric")
	}
	o.dirX, o.dirY, o.dirZ = normalizeDir3(x, y, z)
	return value.Nil, nil
}

func (m *Module) lightSetShadow(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LIGHT.SETSHADOW: heap not bound")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("LIGHT.SETSHADOW expects (light, enabled?)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("LIGHT.SETSHADOW: light must be a handle")
	}
	h := heap.Handle(args[0].IVal)
	o, err := heap.Cast[*lightObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	en, ok := argBoolLight(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("LIGHT.SETSHADOW: enabled must be bool or numeric")
	}
	shadowMu.Lock()
	defer shadowMu.Unlock()
	if en {
		o.shadow = true
		shadowCasterHandle = h
	} else {
		o.shadow = false
		if shadowCasterHandle == h {
			shadowCasterHandle = 0
		}
	}
	return value.Nil, nil
}

func normalizeColorChannel(v float32) float32 {
	if v > 1 {
		v /= 255
	}
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func (m *Module) lightSetColor(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LIGHT.SETCOLOR: heap not bound")
	}
	if len(args) != 4 && len(args) != 5 {
		return value.Nil, fmt.Errorf("LIGHT.SETCOLOR expects (light, r, g, b) or (light, r, g, b, a) with 0–255 or 0.0–1.0")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("LIGHT.SETCOLOR: light must be a handle")
	}
	o, err := heap.Cast[*lightObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	rf, g1, b1 := args[1], args[2], args[3]
	rx, ok1 := argF32(rf)
	gx, ok2 := argF32(g1)
	bx, ok3 := argF32(b1)
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("LIGHT.SETCOLOR: r, g, b must be numeric")
	}
	// Heuristic: values > 1 treated as 0–255.
	if rx > 1 || gx > 1 || bx > 1 {
		o.r = rx / 255
		o.g = gx / 255
		o.b = bx / 255
	} else {
		o.r, o.g, o.b = rx, gx, bx
	}
	o.colA = 1
	if len(args) == 5 {
		ax, ok := argF32(args[4])
		if !ok {
			return value.Nil, fmt.Errorf("LIGHT.SETCOLOR: a must be numeric")
		}
		o.colA = normalizeColorChannel(ax)
	}
	return value.Nil, nil
}

func (m *Module) lightSetIntensity(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LIGHT.SETINTENSITY: heap not bound")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("LIGHT.SETINTENSITY expects (light, amount#)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("LIGHT.SETINTENSITY: light must be a handle")
	}
	o, err := heap.Cast[*lightObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	v, ok := argF32(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("LIGHT.SETINTENSITY: amount must be numeric")
	}
	if v < 0 {
		v = 0
	}
	o.intensity = v
	return value.Nil, nil
}

func (m *Module) lightSetPosition(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LIGHT.SETPOSITION: heap not bound")
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("LIGHT.SETPOSITION expects (light, x#, y#, z#)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("LIGHT.SETPOSITION: light must be a handle")
	}
	o, err := heap.Cast[*lightObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argF32(args[1])
	y, ok2 := argF32(args[2])
	z, ok3 := argF32(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("LIGHT.SETPOSITION: x, y, z must be numeric")
	}
	o.posX, o.posY, o.posZ = x, y, z
	return value.Nil, nil
}

func (m *Module) lightSetTarget(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LIGHT.SETTARGET: heap not bound")
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("LIGHT.SETTARGET expects (light, x#, y#, z#) — shadow frustum look-at point")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("LIGHT.SETTARGET: light must be a handle")
	}
	o, err := heap.Cast[*lightObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argF32(args[1])
	y, ok2 := argF32(args[2])
	z, ok3 := argF32(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("LIGHT.SETTARGET: x, y, z must be numeric")
	}
	o.targetX, o.targetY, o.targetZ = x, y, z
	return value.Nil, nil
}

func (m *Module) lightSetShadowBias(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LIGHT.SETSHADOWBIAS: heap not bound")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("LIGHT.SETSHADOWBIAS expects (light, bias#) — multiplier for depth bias (default 1)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("LIGHT.SETSHADOWBIAS: light must be a handle")
	}
	o, err := heap.Cast[*lightObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	v, ok := argF32(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("LIGHT.SETSHADOWBIAS: bias must be numeric")
	}
	if v < 0.1 {
		v = 0.1
	}
	if v > 8 {
		v = 8
	}
	o.shadowBiasK = v
	return value.Nil, nil
}

func (m *Module) lightSetInnerCone(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.lightCone(rt, "LIGHT.SETINNERCONE", args, true)
}

func (m *Module) lightSetOuterCone(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.lightCone(rt, "LIGHT.SETOUTERCONE", args, false)
}

func (m *Module) lightCone(rt *runtime.Runtime, name string, args []value.Value, inner bool) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("%s: heap not bound", name)
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("%s expects (light, degrees#)", name)
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("%s: light must be a handle", name)
	}
	o, err := heap.Cast[*lightObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	v, ok := argF32(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("%s: degrees must be numeric", name)
	}
	if inner {
		o.innerConeDeg = v
	} else {
		o.outerConeDeg = v
	}
	return value.Nil, nil
}

func (m *Module) lightSetRange(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LIGHT.SETRANGE: heap not bound")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("LIGHT.SETRANGE expects (light, range#)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("LIGHT.SETRANGE: light must be a handle")
	}
	o, err := heap.Cast[*lightObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	v, ok := argF32(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("LIGHT.SETRANGE: range must be numeric")
	}
	if v < 0 {
		v = 0
	}
	o.rangeDist = v
	return value.Nil, nil
}

func (m *Module) lightEnable(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LIGHT.ENABLE: heap not bound")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("LIGHT.ENABLE expects (light, enabled?)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("LIGHT.ENABLE: light must be a handle")
	}
	o, err := heap.Cast[*lightObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	en, ok := argBoolLight(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("LIGHT.ENABLE: enabled must be bool or numeric")
	}
	o.enabled = en
	if !en {
		shadowMu.Lock()
		if shadowCasterHandle == o.self {
			shadowCasterHandle = 0
		}
		shadowMu.Unlock()
	}
	return value.Nil, nil
}

func (m *Module) lightIsEnabled(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LIGHT.ISENABLED: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("LIGHT.ISENABLED expects (light)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("LIGHT.ISENABLED: light must be a handle")
	}
	o, err := heap.Cast[*lightObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if o.enabled {
		return value.FromInt(1), nil
	}
	return value.FromInt(0), nil
}
