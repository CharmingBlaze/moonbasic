package mbjson

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func castJSON(m *Module, h heap.Handle) (*jsonObj, error) {
	return heap.Cast[*jsonObj](m.h, h)
}

func jParse(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.PARSE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.PARSE expects file path string")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	path = strings.TrimSpace(path)
	data, err := os.ReadFile(path)
	if err != nil {
		return value.Nil, fmt.Errorf("JSON.PARSE: %w", err)
	}
	return parseBytes(m, data)
}

func jParseString(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.PARSESTRING: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.PARSESTRING expects string")
	}
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	return parseBytes(m, []byte(s))
}

func parseBytes(m *Module, data []byte) (value.Value, error) {
	var root interface{}
	if err := json.Unmarshal(data, &root); err != nil {
		return value.Nil, fmt.Errorf("JSON: %w", err)
	}
	id, err := m.h.Alloc(&jsonObj{root: root})
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
	id, err := m.h.Alloc(&jsonObj{root: map[string]interface{}{}})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func jMakeArray(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.MAKEARRAY: heap not bound")
	}
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("JSON.MAKEARRAY expects 0 arguments")
	}
	id, err := m.h.Alloc(&jsonObj{root: []interface{}{}})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
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
