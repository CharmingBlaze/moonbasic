//go:build cgo || modernc_sqlite

package mbdb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

var spaceCollapse = regexp.MustCompile(`\s+`)

func normalizeSQL(s string) string {
	s = strings.TrimSpace(s)
	return spaceCollapse.ReplaceAllString(s, " ")
}

func registerDBCommands(m *Module, r runtime.Registrar) {
	r.Register("DB.OPEN", "db", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return dbOpen(m, rt, args...) })
	r.Register("DB.CLOSE", "db", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return dbClose(m, rt, args...) })
	r.Register("DB.ISOPEN", "db", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return dbIsOpen(m, rt, args...) })
	r.Register("DB.EXEC", "db", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return dbExec(m, rt, args...) })
	r.Register("DB.QUERY", "db", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return dbQuery(m, rt, args...) })
	r.Register("DB.QUERYJSON", "db", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return dbQueryJSON(m, rt, args...) })
	r.Register("ROWS.NEXT", "db", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return rowsNext(m, rt, args...) })
	r.Register("ROWS.CLOSE", "db", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return rowsClose(m, rt, args...) })
	r.Register("ROWS.GETSTRING", "db", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return rowsGetString(m, rt, args...) })
	r.Register("ROWS.GETINT", "db", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return rowsGetInt(m, rt, args...) })
	r.Register("ROWS.GETFLOAT", "db", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return rowsGetFloat(m, rt, args...) })
	r.Register("DB.PREPARE", "db", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return dbPrepare(m, rt, args...) })
	r.Register("DB.STMTCLOSE", "db", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return dbStmtClose(m, rt, args...) })
	r.Register("DB.STMTEXEC", "db", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return dbStmtExec(m, rt, args...) })
	r.Register("DB.BEGIN", "db", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return dbBegin(m, rt, args...) })
	r.Register("DB.COMMIT", "db", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return dbCommit(m, rt, args...) })
	r.Register("DB.ROLLBACK", "db", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return dbRollback(m, rt, args...) })
	r.Register("DB.LASTINSERTID", "db", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return dbLastInsertID(m, rt, args...) })
	r.Register("DB.CHANGES", "db", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return dbChanges(m, rt, args...) })
}

func castDB(m *Module, h heap.Handle) (*dbObj, error) {
	return heap.Cast[*dbObj](m.h, h)
}

func castRows(m *Module, h heap.Handle) (*rowsObj, error) {
	return heap.Cast[*rowsObj](m.h, h)
}

func castStmt(m *Module, h heap.Handle) (*stmtObj, error) {
	return heap.Cast[*stmtObj](m.h, h)
}

func castTx(m *Module, h heap.Handle) (*txObj, error) {
	return heap.Cast[*txObj](m.h, h)
}

func dbOpen(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("DB.OPEN: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("DB.OPEN expects path$")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	path = strings.TrimSpace(path)
	d, err := sql.Open(sqliteDriverName, path)
	if err != nil {
		return value.Nil, err
	}
	if err := d.Ping(); err != nil {
		_ = d.Close()
		return value.Nil, err
	}
	id, err := m.h.Alloc(&dbObj{db: d, path: path, stmts: make(map[string]*sql.Stmt)})
	if err != nil {
		_ = d.Close()
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func dbClose(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("DB.CLOSE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("DB.CLOSE expects handle")
	}
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func dbIsOpen(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("DB.ISOPEN: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("DB.ISOPEN expects handle")
	}
	d, err := castDB(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.FromBool(false), nil
	}
	if err := d.assertLive(); err != nil {
		return value.FromBool(false), nil
	}
	return value.FromBool(d.db != nil), nil
}

func argsToAny(args []value.Value, start int) []any {
	out := make([]any, 0, len(args)-start)
	for i := start; i < len(args); i++ {
		out = append(out, valueToAny(args[i]))
	}
	return out
}

func valueToAny(v value.Value) any {
	switch v.Kind {
	case value.KindString:
		// handled by caller with rt — not here
		return nil
	case value.KindInt:
		return v.IVal
	case value.KindFloat:
		return v.FVal
	case value.KindBool:
		if v.IVal != 0 {
			return 1
		}
		return 0
	default:
		return nil
	}
}

func dbExec(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("DB.EXEC: heap not bound")
	}
	if len(args) < 2 {
		return value.Nil, runtime.Errorf("DB.EXEC expects (db, sql$, ...params)")
	}
	if args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("DB.EXEC: bad arguments")
	}
	d, err := castDB(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := d.assertLive(); err != nil {
		return value.Nil, err
	}
	sqlStr, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	params := make([]any, 0, len(args)-2)
	for i := 2; i < len(args); i++ {
		switch args[i].Kind {
		case value.KindString:
			s, e := rt.ArgString(args, i)
			if e != nil {
				return value.Nil, e
			}
			params = append(params, s)
		case value.KindInt:
			params = append(params, args[i].IVal)
		case value.KindFloat:
			params = append(params, args[i].FVal)
		case value.KindBool:
			if args[i].IVal != 0 {
				params = append(params, 1)
			} else {
				params = append(params, 0)
			}
		default:
			return value.Nil, runtime.Errorf("DB.EXEC: unsupported param type")
		}
	}
	var res sql.Result
	if d.tx != nil {
		res, err = d.tx.Exec(sqlStr, params...)
	} else {
		res, err = d.db.Exec(sqlStr, params...)
	}
	if err != nil {
		return value.Nil, err
	}
	_, _ = res.RowsAffected()
	return value.Nil, nil
}

func dbQuery(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("DB.QUERY: heap not bound")
	}
	if len(args) < 2 {
		return value.Nil, runtime.Errorf("DB.QUERY expects (db, sql$, ...params)")
	}
	d, err := castDB(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := d.assertLive(); err != nil {
		return value.Nil, err
	}
	sqlStr, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	params := make([]any, 0, len(args)-2)
	for i := 2; i < len(args); i++ {
		switch args[i].Kind {
		case value.KindString:
			s, e := rt.ArgString(args, i)
			if e != nil {
				return value.Nil, e
			}
			params = append(params, s)
		case value.KindInt:
			params = append(params, args[i].IVal)
		case value.KindFloat:
			params = append(params, args[i].FVal)
		case value.KindBool:
			if args[i].IVal != 0 {
				params = append(params, 1)
			} else {
				params = append(params, 0)
			}
		default:
			return value.Nil, runtime.Errorf("DB.QUERY: unsupported param type")
		}
	}
	var rows *sql.Rows
	if d.tx != nil {
		rows, err = d.tx.Query(sqlStr, params...)
	} else {
		rows, err = d.db.Query(sqlStr, params...)
	}
	if err != nil {
		return value.Nil, err
	}
	id, err := m.h.Alloc(&rowsObj{rows: rows})
	if err != nil {
		_ = rows.Close()
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func dbQueryJSON(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("DB.QUERYJSON: heap not bound")
	}
	if len(args) < 2 {
		return value.Nil, runtime.Errorf("DB.QUERYJSON expects (db, sql$, ...params)")
	}
	d, err := castDB(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := d.assertLive(); err != nil {
		return value.Nil, err
	}
	sqlStr, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	params := make([]any, 0, len(args)-2)
	for i := 2; i < len(args); i++ {
		switch args[i].Kind {
		case value.KindString:
			s, e := rt.ArgString(args, i)
			if e != nil {
				return value.Nil, e
			}
			params = append(params, s)
		case value.KindInt:
			params = append(params, args[i].IVal)
		case value.KindFloat:
			params = append(params, args[i].FVal)
		case value.KindBool:
			if args[i].IVal != 0 {
				params = append(params, 1)
			} else {
				params = append(params, 0)
			}
		default:
			return value.Nil, runtime.Errorf("DB.QUERYJSON: unsupported param type")
		}
	}
	var rows *sql.Rows
	if d.tx != nil {
		rows, err = d.tx.Query(sqlStr, params...)
	} else {
		rows, err = d.db.Query(sqlStr, params...)
	}
	if err != nil {
		return value.Nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return value.Nil, err
	}
	var out []interface{}
	n := len(cols)
	vals := make([]interface{}, n)
	ptrs := make([]interface{}, n)
	for i := range vals {
		ptrs[i] = &vals[i]
	}
	for rows.Next() {
		for i := range vals {
			vals[i] = nil
		}
		if err := rows.Scan(ptrs...); err != nil {
			return value.Nil, err
		}
		row := make(map[string]interface{})
		for i, c := range cols {
			row[c] = vals[i]
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		return value.Nil, err
	}
	b, err := json.Marshal(out)
	if err != nil {
		return value.Nil, err
	}
	return rt.RetString(string(b)), nil
}

func rowsNext(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ROWS.NEXT: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("ROWS.NEXT expects handle")
	}
	r, err := castRows(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if r.freed || r.rows == nil {
		return value.FromBool(false), nil
	}
	if !r.rows.Next() {
		if err := r.rows.Err(); err != nil {
			return value.Nil, err
		}
		return value.FromBool(false), nil
	}
	if r.values == nil {
		cols, err := r.rows.Columns()
		if err != nil {
			return value.Nil, err
		}
		r.cols = cols
		r.values = make([]interface{}, len(cols))
		r.ptrs = make([]interface{}, len(cols))
		for i := range r.values {
			r.ptrs[i] = &r.values[i]
		}
	}
	for i := range r.values {
		r.values[i] = nil
	}
	if err := r.rows.Scan(r.ptrs...); err != nil {
		return value.Nil, err
	}
	return value.FromBool(true), nil
}

func rowsClose(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ROWS.CLOSE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("ROWS.CLOSE expects handle")
	}
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func cellString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case []byte:
		return string(t)
	case string:
		return t
	default:
		return fmt.Sprint(t)
	}
}

func rowsGetString(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("ROWS.GETSTRING expects (rows, col)")
	}
	r, err := castRows(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if r.values == nil {
		return rt.RetString(""), nil
	}
	col, ok := args[1].ToInt()
	if !ok {
		if f, okf := args[1].ToFloat(); okf {
			col = int64(f)
			ok = true
		}
	}
	if !ok {
		return value.Nil, runtime.Errorf("ROWS.GETSTRING: col must be numeric")
	}
	ci := int(col) - 1
	if ci < 0 || ci >= len(r.values) {
		return rt.RetString(""), nil
	}
	return rt.RetString(cellString(r.values[ci])), nil
}

func rowsGetInt(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("ROWS.GETINT expects (rows, col)")
	}
	r, err := castRows(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if r.values == nil {
		return value.FromInt(0), nil
	}
	col, ok := args[1].ToInt()
	if !ok {
		if f, okf := args[1].ToFloat(); okf {
			col = int64(f)
			ok = true
		}
	}
	if !ok {
		return value.Nil, runtime.Errorf("ROWS.GETINT: col must be numeric")
	}
	ci := int(col) - 1
	if ci < 0 || ci >= len(r.values) {
		return value.FromInt(0), nil
	}
	s := cellString(r.values[ci])
	var n int64
	fmt.Sscan(s, &n)
	return value.FromInt(n), nil
}

func rowsGetFloat(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("ROWS.GETFLOAT expects (rows, col)")
	}
	r, err := castRows(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if r.values == nil {
		return value.FromFloat(0), nil
	}
	col, ok := args[1].ToInt()
	if !ok {
		if f, okf := args[1].ToFloat(); okf {
			col = int64(f)
			ok = true
		}
	}
	if !ok {
		return value.Nil, runtime.Errorf("ROWS.GETFLOAT: col must be numeric")
	}
	ci := int(col) - 1
	if ci < 0 || ci >= len(r.values) {
		return value.FromFloat(0), nil
	}
	var f float64
	fmt.Sscan(cellString(r.values[ci]), &f)
	return value.FromFloat(f), nil
}

func dbPrepare(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("DB.PREPARE: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("DB.PREPARE expects (db, sql$)")
	}
	d, err := castDB(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := d.assertLive(); err != nil {
		return value.Nil, err
	}
	sqlStr, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	key := normalizeSQL(sqlStr)
	if st, ok := d.stmts[key]; ok {
		id, err := m.h.Alloc(&stmtObj{db: d, sqlKey: key, stmt: st})
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}
	var st *sql.Stmt
	if d.tx != nil {
		st, err = d.tx.Prepare(sqlStr)
	} else {
		st, err = d.db.Prepare(sqlStr)
	}
	if err != nil {
		return value.Nil, err
	}
	d.stmts[key] = st
	id, err := m.h.Alloc(&stmtObj{db: d, sqlKey: key, stmt: st})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func dbStmtClose(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("DB.STMTCLOSE expects handle")
	}
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func dbStmtExec(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("DB.STMTEXEC expects (stmt, ...params)")
	}
	s, err := castStmt(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if s.freed || s.stmt == nil {
		return value.Nil, runtime.Errorf("DB.STMTEXEC: bad statement")
	}
	d := s.db
	if err := d.assertLive(); err != nil {
		return value.Nil, err
	}
	params := make([]any, 0, len(args)-1)
	for i := 1; i < len(args); i++ {
		switch args[i].Kind {
		case value.KindString:
			str, e := rt.ArgString(args, i)
			if e != nil {
				return value.Nil, e
			}
			params = append(params, str)
		case value.KindInt:
			params = append(params, args[i].IVal)
		case value.KindFloat:
			params = append(params, args[i].FVal)
		case value.KindBool:
			if args[i].IVal != 0 {
				params = append(params, 1)
			} else {
				params = append(params, 0)
			}
		default:
			return value.Nil, runtime.Errorf("DB.STMTEXEC: unsupported param type")
		}
	}
	var res sql.Result
	if d.tx != nil {
		ts := d.tx.Stmt(s.stmt)
		res, err = ts.Exec(params...)
		_ = ts.Close()
	} else {
		res, err = s.stmt.Exec(params...)
	}
	if err != nil {
		return value.Nil, err
	}
	_, _ = res.RowsAffected()
	return value.Nil, nil
}

func dbBegin(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("DB.BEGIN expects db handle")
	}
	d, err := castDB(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := d.assertLive(); err != nil {
		return value.Nil, err
	}
	if d.tx != nil {
		return value.Nil, runtime.Errorf("DB.BEGIN: transaction already active")
	}
	tx, err := d.db.Begin()
	if err != nil {
		return value.Nil, err
	}
	d.tx = tx
	id, err := m.h.Alloc(&txObj{db: d, tx: tx})
	if err != nil {
		_ = tx.Rollback()
		d.tx = nil
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func dbCommit(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("DB.COMMIT expects tx handle")
	}
	t, err := castTx(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if t.freed || t.tx == nil {
		return value.Nil, runtime.Errorf("DB.COMMIT: invalid tx")
	}
	d := t.db
	if err := d.assertLive(); err != nil {
		return value.Nil, err
	}
	if d.tx != t.tx {
		return value.Nil, runtime.Errorf("DB.COMMIT: stale transaction")
	}
	if err := t.tx.Commit(); err != nil {
		return value.Nil, err
	}
	t.committed = true
	d.tx = nil
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func dbRollback(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("DB.ROLLBACK expects tx handle")
	}
	t, err := castTx(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if t.freed || t.tx == nil {
		return value.Nil, runtime.Errorf("DB.ROLLBACK: invalid tx")
	}
	d := t.db
	if err := d.assertLive(); err != nil {
		return value.Nil, err
	}
	if d.tx != t.tx {
		return value.Nil, runtime.Errorf("DB.ROLLBACK: stale transaction")
	}
	if err := t.tx.Rollback(); err != nil {
		return value.Nil, err
	}
	t.committed = true
	d.tx = nil
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func dbLastInsertID(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("DB.LASTINSERTID expects db handle")
	}
	d, err := castDB(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := d.assertLive(); err != nil {
		return value.Nil, err
	}
	var id int64
	err = d.db.QueryRow("SELECT last_insert_rowid()").Scan(&id)
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(id), nil
}

func dbChanges(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("DB.CHANGES expects db handle")
	}
	d, err := castDB(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := d.assertLive(); err != nil {
		return value.Nil, err
	}
	var n int64
	err = d.db.QueryRow("SELECT changes()").Scan(&n)
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(n), nil
}
