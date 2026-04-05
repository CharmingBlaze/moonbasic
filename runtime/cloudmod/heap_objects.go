package cloudmod

import "moonbasic/vm/heap"

type CloudObject struct {
	Coverage float32
	freed    bool
}

func (c *CloudObject) TypeName() string { return "Cloud" }
func (c *CloudObject) TypeTag() uint16  { return heap.TagCloud }
func (c *CloudObject) Free()            { c.freed = true }
