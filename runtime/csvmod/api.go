package mbcsv

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// AllocCSV allocates a CSV handle from raw rows (for TABLE.TOCSV bridge).
func AllocCSV(h *heap.Store, rows [][]string) (value.Value, error) {
	if h == nil {
		return value.Nil, runtime.Errorf("CSV: heap not bound")
	}
	id, err := h.Alloc(&csvObj{rows: rows})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

// Rows returns a copy of CSV data (read-only use).
func Rows(h *heap.Store, id heap.Handle) ([][]string, error) {
	c, err := heap.Cast[*csvObj](h, id)
	if err != nil {
		return nil, err
	}
	if err := c.assertLive(); err != nil {
		return nil, err
	}
	out := make([][]string, len(c.rows))
	for i := range c.rows {
		out[i] = append([]string(nil), c.rows[i]...)
	}
	return out, nil
}
