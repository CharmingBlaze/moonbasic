package mbjson

import (
	"fmt"
	"strconv"
	"strings"
)

// seg is one path step: object key or array index.
type seg struct {
	key   string
	idx   int
	isIdx bool
}

// parseJSONPath parses dotted paths with bracket indices: "a.b[2].c", "[0].x".
func parseJSONPath(path string) ([]seg, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, nil
	}
	var out []seg
	i := 0
	for i < len(path) {
		if path[i] == '.' {
			i++
			continue
		}
		if path[i] == '[' {
			j := i + 1
			for j < len(path) && path[j] != ']' {
				j++
			}
			if j >= len(path) || path[j] != ']' {
				return nil, fmt.Errorf("JSON path: unclosed '['")
			}
			numStr := strings.TrimSpace(path[i+1 : j])
			idx, err := strconv.Atoi(numStr)
			if err != nil {
				return nil, fmt.Errorf("JSON path: invalid array index %q", numStr)
			}
			out = append(out, seg{isIdx: true, idx: idx})
			i = j + 1
			continue
		}
		start := i
		for i < len(path) && path[i] != '.' && path[i] != '[' {
			i++
		}
		key := path[start:i]
		if key == "" {
			return nil, fmt.Errorf("JSON path: empty key segment")
		}
		out = append(out, seg{key: key})
	}
	return out, nil
}

func getValue(root interface{}, segs []seg) (interface{}, bool) {
	cur := root
	for _, s := range segs {
		if s.isIdx {
			arr, ok := cur.([]interface{})
			if !ok || s.idx < 0 || s.idx >= len(arr) {
				return nil, false
			}
			cur = arr[s.idx]
			continue
		}
		m, ok := cur.(map[string]interface{})
		if !ok {
			return nil, false
		}
		v, ok := m[s.key]
		if !ok {
			return nil, false
		}
		cur = v
	}
	return cur, true
}

func hasValue(root interface{}, segs []seg) bool {
	_, ok := getValue(root, segs)
	return ok
}

func setPath(root *interface{}, segs []seg, val interface{}) error {
	if len(segs) == 0 {
		*root = val
		return nil
	}
	sg := segs[0]
	rest := segs[1:]
	if sg.isIdx {
		arr, ok := (*root).([]interface{})
		if !ok || arr == nil {
			arr = []interface{}{}
		}
		for len(arr) <= sg.idx {
			arr = append(arr, nil)
		}
		*root = arr
		return setPath(&arr[sg.idx], rest, val)
	}
	m, ok := (*root).(map[string]interface{})
	if !ok || m == nil {
		m = make(map[string]interface{})
		*root = m
	}
	child := m[sg.key]
	if err := setPath(&child, rest, val); err != nil {
		return err
	}
	m[sg.key] = child
	return nil
}

func deletePath(j *jsonObj, segs []seg) error {
	if len(segs) == 0 {
		return fmt.Errorf("JSON.DELETE: path required")
	}
	last := segs[len(segs)-1]
	prefix := segs[:len(segs)-1]
	parent, ok := getValue(j.root, prefix)
	if !ok {
		return nil
	}
	if last.isIdx {
		arr, ok := parent.([]interface{})
		if !ok || last.idx < 0 || last.idx >= len(arr) {
			return nil
		}
		copy(arr[last.idx:], arr[last.idx+1:])
		arr[len(arr)-1] = nil
		arr = arr[:len(arr)-1]
		if len(prefix) == 0 {
			j.root = arr
		}
		return nil
	}
	m, ok := parent.(map[string]interface{})
	if !ok {
		return fmt.Errorf("JSON.DELETE: parent is not object")
	}
	delete(m, last.key)
	return nil
}

func clearContainer(j *jsonObj, segs []seg) error {
	if len(segs) == 0 {
		switch v := j.root.(type) {
		case map[string]interface{}:
			for k := range v {
				delete(v, k)
			}
			return nil
		case []interface{}:
			j.root = v[:0]
			return nil
		default:
			return fmt.Errorf("JSON.CLEAR: root must be object or array")
		}
	}
	v, ok := getValue(j.root, segs)
	if !ok {
		return nil
	}
	switch t := v.(type) {
	case map[string]interface{}:
		for k := range t {
			delete(t, k)
		}
	case []interface{}:
		return setPath(&j.root, segs, []interface{}{})
	default:
		return fmt.Errorf("JSON.CLEAR: value at path is not object or array")
	}
	return nil
}

func appendValue(j *jsonObj, segs []seg, val interface{}) error {
	v, ok := getValue(j.root, segs)
	if !ok {
		// create array at path
		return setPath(&j.root, segs, []interface{}{val})
	}
	arr, ok := v.([]interface{})
	if !ok {
		return fmt.Errorf("JSON.APPEND: path is not an array")
	}
	arr = append(arr, val)
	if len(segs) == 0 {
		j.root = arr
		return nil
	}
	return setPath(&j.root, segs, arr)
}
