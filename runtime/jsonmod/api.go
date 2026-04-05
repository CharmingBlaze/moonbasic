package mbjson

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// AllocJSONRoot allocates a JSON handle wrapping an arbitrary decoded value (for CSV/TABLE bridges).
func AllocJSONRoot(h *heap.Store, root interface{}) (value.Value, error) {
	if h == nil {
		return value.Nil, runtime.Errorf("JSON: heap not bound")
	}
	id, err := h.Alloc(&jsonObj{root: root})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

// Root returns the decoded value inside a JSON handle (read-only for callers; do not mutate).
func Root(h *heap.Store, jsonHandle heap.Handle) (interface{}, error) {
	j, err := heap.Cast[*jsonObj](h, jsonHandle)
	if err != nil {
		return nil, err
	}
	if err := j.assertLive(); err != nil {
		return nil, err
	}
	return j.root, nil
}
