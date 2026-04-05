package mbtable

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

type tableObj struct {
	cols []string
	rows [][]interface{}
	freed bool
}

func (t *tableObj) TypeName() string { return "TABLE" }

func (t *tableObj) TypeTag() uint16 { return heap.TagTable }

func (t *tableObj) Free() {
	if t.freed {
		return
	}
	t.freed = true
	t.cols = nil
	t.rows = nil
}

func (t *tableObj) assertLive() error {
	if t.freed {
		return runtime.Errorf("TABLE: use after free")
	}
	return nil
}

func colIndex(cols []string, name string) int {
	for i, c := range cols {
		if c == name {
			return i
		}
	}
	return -1
}
