package player

import "moonbasic/vm/heap"

// charRefHeapObj is the heap object for CHARACTER.CREATE (handle-typed Character API).
// TagCharController matches CHARACTER. / CHARACTERREF.* handle dispatch.
type charRefHeapObj struct {
	id int64
	m  *Module
	release heap.ReleaseOnce
}

func (c *charRefHeapObj) TypeName() string { return "Character" }

func (c *charRefHeapObj) TypeTag() uint16 { return heap.TagCharController }

func (c *charRefHeapObj) Free() {
	c.release.Do(func() {
		charRefHeapObjFree(c.m, c.id)
	})
}
