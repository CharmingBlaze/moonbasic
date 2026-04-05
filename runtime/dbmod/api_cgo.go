//go:build cgo

package mbdb

import (
	"database/sql"

	"moonbasic/vm/heap"
)

// Connection returns the underlying *sql.DB for an open DB handle (bridges only).
func Connection(h *heap.Store, id heap.Handle) (*sql.DB, error) {
	d, err := heap.Cast[*dbObj](h, id)
	if err != nil {
		return nil, err
	}
	if err := d.assertLive(); err != nil {
		return nil, err
	}
	return d.db, nil
}
