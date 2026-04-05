package mbdb

import (
	"database/sql"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

type dbObj struct {
	db    *sql.DB
	path  string
	freed bool
	stmts map[string]*sql.Stmt
	tx    *sql.Tx
}

func (d *dbObj) TypeName() string { return "DB" }
func (d *dbObj) TypeTag() uint16 { return heap.TagDB }

func (d *dbObj) Free() {
	if d.freed {
		return
	}
	d.freed = true
	if d.tx != nil {
		_ = d.tx.Rollback()
		d.tx = nil
	}
	for _, s := range d.stmts {
		_ = s.Close()
	}
	d.stmts = nil
	if d.db != nil {
		_ = d.db.Close()
		d.db = nil
	}
}

func (d *dbObj) assertLive() error {
	if d.freed {
		return runtime.Errorf("DB: use after free")
	}
	return nil
}

type stmtObj struct {
	db     *dbObj
	sqlKey string
	stmt   *sql.Stmt
	freed  bool
}

func (s *stmtObj) TypeName() string { return "STMT" }
func (s *stmtObj) TypeTag() uint16  { return heap.TagDBStmt }

func (s *stmtObj) Free() {
	if s.freed {
		return
	}
	s.freed = true
	// Cached stmt owned by dbObj; do not Close here.
	s.stmt = nil
	s.db = nil
}

type rowsObj struct {
	rows   *sql.Rows
	cols   []string
	values []interface{}
	ptrs   []interface{}
	freed  bool
}

func (r *rowsObj) TypeName() string { return "ROWS" }
func (r *rowsObj) TypeTag() uint16  { return heap.TagDBRows }

func (r *rowsObj) Free() {
	if r.freed {
		return
	}
	r.freed = true
	if r.rows != nil {
		_ = r.rows.Close()
		r.rows = nil
	}
	r.values = nil
	r.ptrs = nil
}

type txObj struct {
	db        *dbObj
	tx        *sql.Tx
	committed bool
	freed     bool
}

func (t *txObj) TypeName() string { return "TX" }
func (t *txObj) TypeTag() uint16  { return heap.TagDBTx }

func (t *txObj) Free() {
	if t.freed {
		return
	}
	t.freed = true
	if t.tx != nil && !t.committed {
		_ = t.tx.Rollback()
		if t.db != nil {
			t.db.tx = nil
		}
	}
	t.tx = nil
	t.db = nil
}
