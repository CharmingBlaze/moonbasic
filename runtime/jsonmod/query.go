package mbjson

import (
	"fmt"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// JSON.QUERY supports a small pattern language:
//   "path.to.array[*].field" — for each element of the array at path.to.array, collect .field (stringified).
// Bracket indices use normal path rules: "items[0].x" returns one value as a single-element list.
func jQuery(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.QUERY: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.QUERY expects (handle, pattern$)")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := j.assertLive(); err != nil {
		return value.Nil, err
	}
	pat, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	pat = strings.TrimSpace(pat)
	if pat == "" {
		id, err := m.h.Alloc(&heap.StringList{Items: nil})
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}
	const wild = "[*]"
	if i := strings.Index(pat, wild); i >= 0 {
		prefix := pat[:i]
		suffix := pat[i+len(wild):]
		if suffix != "" && !strings.HasPrefix(suffix, ".") && !strings.HasPrefix(suffix, "[") {
			return value.Nil, fmt.Errorf("JSON.QUERY: after [*] expect '.' or '['")
		}
		segs, err := parseJSONPath(prefix)
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
		arr, ok := v.([]interface{})
		if !ok {
			return value.Nil, runtime.Errorf("JSON.QUERY: [*] requires array at prefix path")
		}
		suffix = strings.TrimPrefix(suffix, ".")
		var rest []seg
		if suffix != "" {
			rest, err = parseJSONPath(suffix)
			if err != nil {
				return value.Nil, err
			}
		}
		var out []string
		for _, el := range arr {
			s := cellString(el)
			if len(rest) > 0 {
				v2, ok2 := getValue(el, rest)
				if !ok2 {
					out = append(out, "")
					continue
				}
				s = cellString(v2)
			}
			out = append(out, s)
		}
		id, err := m.h.Alloc(&heap.StringList{Items: out})
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}
	segs, err := parseJSONPath(pat)
	if err != nil {
		return value.Nil, err
	}
	v, ok := getValue(j.root, segs)
	if !ok {
		id, err := m.h.Alloc(&heap.StringList{Items: []string{""}})
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}
	id, err := m.h.Alloc(&heap.StringList{Items: []string{cellString(v)}})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}
