package mbjson

import (
	"encoding/json"
	"fmt"
	"strconv"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func pathFromArgs(rt *runtime.Runtime, args []value.Value, idx int) ([]seg, error) {
	s, err := rt.ArgString(args, idx)
	if err != nil {
		return nil, err
	}
	return parseJSONPath(s)
}

func jHas(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.HAS: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("JSON.HAS expects (handle, path$)")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := j.assertLive(); err != nil {
		return value.Nil, err
	}
	segs, err := pathFromArgs(rt, args, 1)
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(hasValue(j.root, segs)), nil
}

func jTypeStr(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.TYPE$: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("JSON.TYPE$ expects (handle, path$)")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := j.assertLive(); err != nil {
		return value.Nil, err
	}
	segs, err := pathFromArgs(rt, args, 1)
	if err != nil {
		return value.Nil, err
	}
	v, ok := getValue(j.root, segs)
	if !ok {
		return rt.RetString("missing"), nil
	}
	return rt.RetString(jsonTypeOf(v)), nil
}

func jsonTypeOf(v interface{}) string {
	if v == nil {
		return "null"
	}
	switch v.(type) {
	case map[string]interface{}:
		return "object"
	case []interface{}:
		return "array"
	case string:
		return "string"
	case bool:
		return "bool"
	case float64, json.Number:
		return "number"
	default:
		return fmt.Sprintf("%T", v)
	}
}

func jLen(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.LEN: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("JSON.LEN expects (handle, path$)")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := j.assertLive(); err != nil {
		return value.Nil, err
	}
	segs, err := pathFromArgs(rt, args, 1)
	if err != nil {
		return value.Nil, err
	}
	v, ok := getValue(j.root, segs)
	if !ok {
		return value.FromInt(0), nil
	}
	switch t := v.(type) {
	case map[string]interface{}:
		return value.FromInt(int64(len(t))), nil
	case []interface{}:
		return value.FromInt(int64(len(t))), nil
	case string:
		return value.FromInt(int64(len(t))), nil
	default:
		return value.FromInt(0), nil
	}
}

func jKeys(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.KEYS: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("JSON.KEYS expects (handle, path$)")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := j.assertLive(); err != nil {
		return value.Nil, err
	}
	segs, err := pathFromArgs(rt, args, 1)
	if err != nil {
		return value.Nil, err
	}
	v, ok := getValue(j.root, segs)
	if !ok {
		id, err := m.h.Alloc(&heap.StringList{Items: nil})
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}
	mo, ok := v.(map[string]interface{})
	if !ok {
		return value.Nil, runtime.Errorf("JSON.KEYS: value is not object")
	}
	keys := make([]string, 0, len(mo))
	for k := range mo {
		keys = append(keys, k)
	}
	// deterministic order for tests/docs
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[j] < keys[i] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	id, err := m.h.Alloc(&heap.StringList{Items: keys})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func jGetString(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.GETSTRING: heap not bound")
	}
	if len(args) != 2 && len(args) != 3 {
		return value.Nil, runtime.Errorf("JSON.GETSTRING expects (handle, path$) or (handle, path$, default$)")
	}
	if args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.GETSTRING: bad arguments")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := j.assertLive(); err != nil {
		return value.Nil, err
	}
	segs, err := pathFromArgs(rt, args, 1)
	if err != nil {
		return value.Nil, err
	}
	v, ok := getValue(j.root, segs)
	if !ok {
		if len(args) == 3 {
			s, err := rt.ArgString(args, 2)
			if err != nil {
				return value.Nil, err
			}
			return rt.RetString(s), nil
		}
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
	if len(args) != 2 && len(args) != 3 {
		return value.Nil, runtime.Errorf("JSON.GETINT expects (handle, path$) or (handle, path$, default)")
	}
	if args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.GETINT: bad arguments")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := j.assertLive(); err != nil {
		return value.Nil, err
	}
	segs, err := pathFromArgs(rt, args, 1)
	if err != nil {
		return value.Nil, err
	}
	v, ok := getValue(j.root, segs)
	if !ok {
		if len(args) == 3 {
			if n, ok2 := args[2].ToInt(); ok2 {
				return value.FromInt(n), nil
			}
			if f, okf := args[2].ToFloat(); okf {
				return value.FromInt(int64(f)), nil
			}
		}
		return value.FromInt(0), nil
	}
	n := coerceInt(v)
	return value.FromInt(n), nil
}

func coerceInt(v interface{}) int64 {
	switch t := v.(type) {
	case int:
		return int64(t)
	case int64:
		return t
	case float64:
		return int64(t)
	case json.Number:
		i, err := t.Int64()
		if err != nil {
			f, _ := t.Float64()
			return int64(f)
		}
		return i
	case string:
		n, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return 0
		}
		return n
	case bool:
		if t {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func jGetFloat(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.GETFLOAT: heap not bound")
	}
	if len(args) != 2 && len(args) != 3 {
		return value.Nil, runtime.Errorf("JSON.GETFLOAT expects (handle, path$) or (handle, path$, default)")
	}
	if args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.GETFLOAT: bad arguments")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := j.assertLive(); err != nil {
		return value.Nil, err
	}
	segs, err := pathFromArgs(rt, args, 1)
	if err != nil {
		return value.Nil, err
	}
	v, ok := getValue(j.root, segs)
	if !ok {
		if len(args) == 3 {
			if f, ok2 := args[2].ToFloat(); ok2 {
				return value.FromFloat(f), nil
			}
			if n, okn := args[2].ToInt(); okn {
				return value.FromFloat(float64(n)), nil
			}
		}
		return value.FromFloat(0), nil
	}
	return value.FromFloat(coerceFloat(v)), nil
}

func coerceFloat(v interface{}) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case int:
		return float64(t)
	case int64:
		return float64(t)
	case json.Number:
		f, err := t.Float64()
		if err != nil {
			return 0
		}
		return f
	case string:
		f, err := strconv.ParseFloat(t, 64)
		if err != nil {
			return 0
		}
		return f
	default:
		return 0
	}
}

func jGetBool(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.GETBOOL: heap not bound")
	}
	if len(args) != 2 && len(args) != 3 {
		return value.Nil, runtime.Errorf("JSON.GETBOOL expects (handle, path$) or (handle, path$, default)")
	}
	if args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.GETBOOL: bad arguments")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := j.assertLive(); err != nil {
		return value.Nil, err
	}
	segs, err := pathFromArgs(rt, args, 1)
	if err != nil {
		return value.Nil, err
	}
	v, ok := getValue(j.root, segs)
	if !ok {
		if len(args) == 3 && args[2].Kind == value.KindBool {
			return value.FromBool(args[2].IVal != 0), nil
		}
		return value.FromBool(false), nil
	}
	b, ok := v.(bool)
	if !ok {
		return value.FromBool(false), nil
	}
	return value.FromBool(b), nil
}

func jGetArray(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return jGetSub(m, rt, args, "array", "JSON.GETARRAY")
}

func jGetObject(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return jGetSub(m, rt, args, "object", "JSON.GETOBJECT")
}

func jGetSub(m *Module, rt *runtime.Runtime, args []value.Value, want, op string) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("%s: heap not bound", op)
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("%s expects (handle, path$)", op)
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := j.assertLive(); err != nil {
		return value.Nil, err
	}
	segs, err := pathFromArgs(rt, args, 1)
	if err != nil {
		return value.Nil, err
	}
	v, ok := getValue(j.root, segs)
	if !ok {
		return value.Nil, runtime.Errorf("%s: path not found", op)
	}
	if want == "array" {
		if _, ok := v.([]interface{}); !ok {
			return value.Nil, runtime.Errorf("%s: not an array", op)
		}
	} else {
		if _, ok := v.(map[string]interface{}); !ok {
			return value.Nil, runtime.Errorf("%s: not an object", op)
		}
	}
	id, err := m.h.Alloc(&jsonObj{root: v})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}
