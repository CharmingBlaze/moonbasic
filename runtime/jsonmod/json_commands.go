package mbjson

import (
	"encoding/json"
	"fmt"
	"strconv"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type jsonObj struct {
	m map[string]any
}

func (j *jsonObj) TypeName() string { return "JSON" }

func (j *jsonObj) TypeTag() uint16 { return heap.TagJSON }

func (j *jsonObj) Free() { j.m = nil }

func registerJSONCommands(m *Module, r runtime.Registrar) {
	r.Register("JSON.PARSE", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jParse(m, rt, args...) })
	r.Register("JSON.MAKE", "json", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return jMake(m, a) }))
	r.Register("JSON.GETSTRING", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jGetString(m, rt, args...) })
	r.Register("JSON.GETINT", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jGetInt(m, rt, args...) })
	r.Register("JSON.GETFLOAT", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jGetFloat(m, rt, args...) })
	r.Register("JSON.GETBOOL", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jGetBool(m, rt, args...) })
	r.Register("JSON.SETSTRING", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jSetString(m, rt, args...) })
	r.Register("JSON.SETINT", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jSetInt(m, rt, args...) })
	r.Register("JSON.SETFLOAT", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jSetFloat(m, rt, args...) })
	r.Register("JSON.SETBOOL", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jSetBool(m, rt, args...) })
	r.Register("JSON.TOSTRING", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jToString(m, rt, args...) })
	r.Register("JSON.FREE", "json", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return jFree(m, a) }))
}

func castJSON(m *Module, h heap.Handle) (*jsonObj, error) {
	return heap.Cast[*jsonObj](m.h, h)
}

func jParse(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.PARSE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.PARSE expects string")
	}
	var raw map[string]any
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := json.Unmarshal([]byte(s), &raw); err != nil {
		return value.Nil, fmt.Errorf("JSON.PARSE: %w", err)
	}
	if raw == nil {
		raw = make(map[string]any)
	}
	for k, v := range raw {
		if _, isMap := v.(map[string]any); isMap {
			return value.Nil, fmt.Errorf("JSON.PARSE: nested object at %q not supported", k)
		}
		if _, isArr := v.([]any); isArr {
			return value.Nil, fmt.Errorf("JSON.PARSE: nested array at %q not supported", k)
		}
	}
	id, err := m.h.Alloc(&jsonObj{m: raw})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func jMake(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.MAKE: heap not bound")
	}
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("JSON.MAKE expects 0 arguments")
	}
	id, err := m.h.Alloc(&jsonObj{m: make(map[string]any)})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func jGetString(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.GETSTRING: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.GETSTRING expects (handle, key$)")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if j.m == nil {
		return rt.RetString(""), nil
	}
	key, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	v, ok := j.m[key]
	if !ok {
		return rt.RetString(""), nil
	}
	switch t := v.(type) {
	case string:
		return rt.RetString(t), nil
	default:
		return rt.RetString(fmt.Sprint(t)), nil
	}
}

func jGetInt(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.GETINT: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.GETINT expects (handle, key$)")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if j.m == nil {
		return value.FromInt(0), nil
	}
	key, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	v, ok := j.m[key]
	if !ok {
		return value.FromInt(0), nil
	}
	switch t := v.(type) {
	case int:
		return value.FromInt(int64(t)), nil
	case int64:
		return value.FromInt(t), nil
	case float64:
		return value.FromInt(int64(t)), nil
	case json.Number:
		i, err := t.Int64()
		if err != nil {
			f, _ := t.Float64()
			return value.FromInt(int64(f)), nil
		}
		return value.FromInt(i), nil
	case string:
		n, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return value.FromInt(0), nil
		}
		return value.FromInt(n), nil
	default:
		return value.FromInt(0), nil
	}
}

func jGetFloat(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.GETFLOAT: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.GETFLOAT expects (handle, key$)")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if j.m == nil {
		return value.FromFloat(0), nil
	}
	key, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	v, ok := j.m[key]
	if !ok {
		return value.FromFloat(0), nil
	}
	switch t := v.(type) {
	case float64:
		return value.FromFloat(t), nil
	case int:
		return value.FromFloat(float64(t)), nil
	case int64:
		return value.FromFloat(float64(t)), nil
	case json.Number:
		f, err := t.Float64()
		if err != nil {
			return value.FromFloat(0), nil
		}
		return value.FromFloat(f), nil
	default:
		return value.FromFloat(0), nil
	}
}

func jGetBool(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.GETBOOL: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.GETBOOL expects (handle, key$)")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if j.m == nil {
		return value.FromBool(false), nil
	}
	key, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	v, ok := j.m[key]
	if !ok {
		return value.FromBool(false), nil
	}
	switch t := v.(type) {
	case bool:
		return value.FromBool(t), nil
	default:
		return value.FromBool(false), nil
	}
}

func jSetString(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.SETSTRING: heap not bound")
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString || args[2].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.SETSTRING expects (handle, key$, val$)")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	key, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	val, err := rt.ArgString(args, 2)
	if err != nil {
		return value.Nil, err
	}
	j.m[key] = val
	return value.Nil, nil
}

func jSetInt(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.SETINT: heap not bound")
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.SETINT expects (handle, key$, int)")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	var n int64
	switch args[2].Kind {
	case value.KindInt:
		n = args[2].IVal
	case value.KindFloat:
		n = int64(args[2].FVal)
	default:
		return value.Nil, runtime.Errorf("JSON.SETINT: value must be numeric")
	}
	key, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	j.m[key] = float64(n) // encode as JSON number
	return value.Nil, nil
}

func jSetFloat(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.SETFLOAT: heap not bound")
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.SETFLOAT expects (handle, key$, float)")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	f, ok := args[2].ToFloat()
	if !ok {
		return value.Nil, runtime.Errorf("JSON.SETFLOAT: value must be numeric")
	}
	key, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	j.m[key] = f
	return value.Nil, nil
}

func jSetBool(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.SETBOOL: heap not bound")
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.SETBOOL expects (handle, key$, bool)")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	var b bool
	switch args[2].Kind {
	case value.KindBool:
		b = args[2].IVal != 0
	case value.KindInt:
		b = args[2].IVal != 0
	default:
		return value.Nil, runtime.Errorf("JSON.SETBOOL: value must be bool")
	}
	key, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	j.m[key] = b
	return value.Nil, nil
}

func jToString(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.TOSTRING: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("JSON.TOSTRING expects handle")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	b, err := json.Marshal(j.m)
	if err != nil {
		return value.Nil, err
	}
	return rt.RetString(string(b)), nil
}

func jFree(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.FREE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("JSON.FREE expects handle")
	}
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}
