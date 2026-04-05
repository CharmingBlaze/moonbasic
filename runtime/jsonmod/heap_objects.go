package mbjson

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

// jsonObj holds a decoded JSON value (object, array, or scalar) from encoding/json.
type jsonObj struct {
	root  interface{}
	freed bool
}

func (j *jsonObj) TypeName() string { return "JSON" }

func (j *jsonObj) TypeTag() uint16 { return heap.TagJSON }

func (j *jsonObj) Free() {
	if j.freed {
		return
	}
	j.freed = true
	j.root = nil
}

func (j *jsonObj) assertLive() error {
	if j.freed {
		return runtime.Errorf("JSON: use after free")
	}
	return nil
}
