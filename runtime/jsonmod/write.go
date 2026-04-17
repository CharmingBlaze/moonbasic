package mbjson

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func jSetString(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.SETSTRING: heap not bound")
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString || args[2].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.SETSTRING expects (handle, path$, val$)")
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
	val, err := rt.ArgString(args, 2)
	if err != nil {
		return value.Nil, err
	}
	if err := setPath(&j.root, segs, val); err != nil {
		return value.Nil, err
	}
	return args[0], nil
}

func jSetInt(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.SETINT: heap not bound")
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.SETINT expects (handle, path$, int)")
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
	var n int64
	switch args[2].Kind {
	case value.KindInt:
		n = args[2].IVal
	case value.KindFloat:
		n = int64(args[2].FVal)
	default:
		return value.Nil, runtime.Errorf("JSON.SETINT: value must be numeric")
	}
	if err := setPath(&j.root, segs, float64(n)); err != nil {
		return value.Nil, err
	}
	return args[0], nil
}

func jSetFloat(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.SETFLOAT: heap not bound")
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.SETFLOAT expects (handle, path$, float)")
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
	f, ok := args[2].ToFloat()
	if !ok {
		return value.Nil, runtime.Errorf("JSON.SETFLOAT: value must be numeric")
	}
	if err := setPath(&j.root, segs, f); err != nil {
		return value.Nil, err
	}
	return args[0], nil
}

func jSetBool(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.SETBOOL: heap not bound")
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.SETBOOL expects (handle, path$, bool)")
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
	var b bool
	switch args[2].Kind {
	case value.KindBool:
		b = args[2].IVal != 0
	case value.KindInt:
		b = args[2].IVal != 0
	default:
		return value.Nil, runtime.Errorf("JSON.SETBOOL: value must be bool")
	}
	if err := setPath(&j.root, segs, b); err != nil {
		return value.Nil, err
	}
	return args[0], nil
}

func jSetNull(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.SETNULL: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.SETNULL expects (handle, path$)")
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
	if err := setPath(&j.root, segs, nil); err != nil {
		return value.Nil, err
	}
	return args[0], nil
}

func jDelete(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.DELETE: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.DELETE expects (handle, path$)")
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
	if err := deletePath(j, segs); err != nil {
		return value.Nil, err
	}
	return args[0], nil
}

func jClear(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.CLEAR: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.CLEAR expects (handle, path$)")
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
	if err := clearContainer(j, segs); err != nil {
		return value.Nil, err
	}
	return args[0], nil
}

func jAppend(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.APPEND: heap not bound")
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.APPEND expects (handle, path$, value)")
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
	val, err := valueToJSON(rt, args[2])
	if err != nil {
		return value.Nil, err
	}
	if err := appendValue(j, segs, val); err != nil {
		return value.Nil, err
	}
	return args[0], nil
}

func valueToJSON(rt *runtime.Runtime, v value.Value) (interface{}, error) {
	switch v.Kind {
	case value.KindString:
		return rt.ArgString([]value.Value{v}, 0)
	case value.KindInt:
		return float64(v.IVal), nil
	case value.KindFloat:
		return v.FVal, nil
	case value.KindBool:
		return v.IVal != 0, nil
	default:
		return nil, runtime.Errorf("JSON.APPEND: unsupported value type")
	}
}
