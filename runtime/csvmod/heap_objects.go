package mbcsv

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

// csvObj holds tabular CSV data (rows of string fields), including a header row when present.
type csvObj struct {
	rows  [][]string
	freed bool
}

func (c *csvObj) TypeName() string { return "CSV" }

func (c *csvObj) TypeTag() uint16 { return heap.TagCSV }

func (c *csvObj) Free() {
	if c.freed {
		return
	}
	c.freed = true
	c.rows = nil
}

func (c *csvObj) assertLive() error {
	if c.freed {
		return runtime.Errorf("CSV: use after free")
	}
	return nil
}
