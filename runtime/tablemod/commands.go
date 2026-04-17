package mbtable

import (
	"fmt"
	"strings"

	"moonbasic/runtime"
	mbcsv "moonbasic/runtime/csvmod"
	mbjson "moonbasic/runtime/jsonmod"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func castTable(m *Module, h heap.Handle) (*tableObj, error) {
	return heap.Cast[*tableObj](m.h, h)
}

func registerTableCommands(m *Module, r runtime.Registrar) {
	r.Register("TABLE.CREATE", "table", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return tabCreate(m, rt, args...) })
	r.Register("TABLE.MAKE", "table", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return tabCreate(m, rt, args...) })
	r.Register("TABLE.FREE", "table", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return tabFree(m, rt, args...) })
	r.Register("TABLE.ADDROW", "table", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return tabAddRow(m, rt, args...) })
	r.Register("TABLE.ROWCOUNT", "table", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return tabRowCount(m, rt, args...)
	})
	r.Register("TABLE.COLCOUNT", "table", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return tabColCount(m, rt, args...)
	})
	r.Register("TABLE.GET", "table", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return tabGet(m, rt, args...) })
	r.Register("TABLE.SET", "table", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return tabSet(m, rt, args...) })
	r.Register("TABLE.TOJSON", "table", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return tabToJSON(m, rt, args...) })
	r.Register("TABLE.FROMJSON", "table", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return tabFromJSON(m, rt, args...)
	})
	r.Register("TABLE.TOCSV", "table", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return tabToCSV(m, rt, args...) })
	r.Register("TABLE.FROMCSV", "table", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return tabFromCSV(m, rt, args...) })
}

func splitCols(s string) []string {
	parts := strings.Split(s, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func tabCreate(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TABLE.CREATE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("TABLE.CREATE expects cols$ (comma-separated)")
	}
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	cols := splitCols(s)
	if len(cols) == 0 {
		return value.Nil, runtime.Errorf("TABLE.CREATE: need at least one column")
	}
	id, err := m.h.Alloc(&tableObj{cols: cols, rows: nil})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func tabFree(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TABLE.FREE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("TABLE.FREE expects handle")
	}
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func tabAddRow(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TABLE.ADDROW: heap not bound")
	}
	if len(args) < 2 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("TABLE.ADDROW expects (handle, ...)")
	}
	t, err := castTable(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := t.assertLive(); err != nil {
		return value.Nil, err
	}
	if len(args)-1 != len(t.cols) {
		return value.Nil, fmt.Errorf("TABLE.ADDROW: expected %d values, got %d", len(t.cols), len(args)-1)
	}
	row := make([]interface{}, len(t.cols))
	for i := 1; i < len(args); i++ {
		row[i-1], err = cellFromValue(rt, args[i])
		if err != nil {
			return value.Nil, err
		}
	}
	t.rows = append(t.rows, row)
	return args[0], nil
}

func cellFromValue(rt *runtime.Runtime, v value.Value) (interface{}, error) {
	switch v.Kind {
	case value.KindString:
		return rt.ArgString([]value.Value{v}, 0)
	case value.KindInt:
		return float64(v.IVal), nil
	case value.KindFloat:
		return v.FVal, nil
	case value.KindBool:
		return v.IVal != 0, nil
	default:
		return nil, runtime.Errorf("TABLE: unsupported cell type")
	}
}

func tabRowCount(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("TABLE.ROWCOUNT expects handle")
	}
	t, err := castTable(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := t.assertLive(); err != nil {
		return value.Nil, err
	}
	return value.FromInt(int64(len(t.rows))), nil
}

func tabColCount(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("TABLE.COLCOUNT expects handle")
	}
	t, err := castTable(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := t.assertLive(); err != nil {
		return value.Nil, err
	}
	return value.FromInt(int64(len(t.cols))), nil
}

func tabGet(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 || args[0].Kind != value.KindHandle || args[2].Kind != value.KindString {
		return value.Nil, runtime.Errorf("TABLE.GET expects (handle, row, col$)")
	}
	t, err := castTable(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := t.assertLive(); err != nil {
		return value.Nil, err
	}
	ri, ok := args[1].ToInt()
	if !ok {
		if f, okf := args[1].ToFloat(); okf {
			ri = int64(f)
			ok = true
		}
	}
	if !ok {
		return value.Nil, runtime.Errorf("TABLE.GET: row must be numeric")
	}
	row := int(ri) - 1
	colName, err := rt.ArgString(args, 2)
	if err != nil {
		return value.Nil, err
	}
	ci := colIndex(t.cols, colName)
	if row < 0 || row >= len(t.rows) || ci < 0 {
		return rt.RetString(""), nil
	}
	return cellToValue(rt, t.rows[row][ci])
}

func cellToValue(rt *runtime.Runtime, v interface{}) (value.Value, error) {
	switch x := v.(type) {
	case string:
		return rt.RetString(x), nil
	case float64:
		return value.FromFloat(x), nil
	case bool:
		return value.FromBool(x), nil
	case nil:
		return rt.RetString(""), nil
	default:
		return rt.RetString(fmt.Sprint(x)), nil
	}
}

func tabSet(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 || args[0].Kind != value.KindHandle || args[2].Kind != value.KindString {
		return value.Nil, runtime.Errorf("TABLE.SET expects (handle, row, col$, value)")
	}
	t, err := castTable(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := t.assertLive(); err != nil {
		return value.Nil, err
	}
	ri, ok := args[1].ToInt()
	if !ok {
		if f, okf := args[1].ToFloat(); okf {
			ri = int64(f)
			ok = true
		}
	}
	if !ok {
		return value.Nil, runtime.Errorf("TABLE.SET: row must be numeric")
	}
	row := int(ri) - 1
	colName, err := rt.ArgString(args, 2)
	if err != nil {
		return value.Nil, err
	}
	ci := colIndex(t.cols, colName)
	if row < 0 || row >= len(t.rows) || ci < 0 {
		return value.Nil, runtime.Errorf("TABLE.SET: bad row/column")
	}
	val, err := cellFromValue(rt, args[3])
	if err != nil {
		return value.Nil, err
	}
	t.rows[row][ci] = val
	return args[0], nil
}

func tabToJSON(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("TABLE.TOJSON expects handle")
	}
	t, err := castTable(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := t.assertLive(); err != nil {
		return value.Nil, err
	}
	var arr []interface{}
	for _, row := range t.rows {
		mo := make(map[string]interface{})
		for i, c := range t.cols {
			if i < len(row) {
				mo[c] = row[i]
			}
		}
		arr = append(arr, mo)
	}
	return mbjson.AllocJSONRoot(m.h, arr)
}

func tabFromJSON(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("TABLE.FROMJSON expects json handle")
	}
	root, err := mbjson.Root(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	arr, ok := root.([]interface{})
	if !ok {
		return value.Nil, runtime.Errorf("TABLE.FROMJSON: JSON root must be array")
	}
	if len(arr) == 0 {
		id, err := m.h.Alloc(&tableObj{cols: nil, rows: nil})
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}
	first, ok := arr[0].(map[string]interface{})
	if !ok {
		return value.Nil, runtime.Errorf("TABLE.FROMJSON: array elements must be objects")
	}
	var cols []string
	for k := range first {
		cols = append(cols, k)
	}
	for i := 0; i < len(cols); i++ {
		for j := i + 1; j < len(cols); j++ {
			if cols[j] < cols[i] {
				cols[i], cols[j] = cols[j], cols[i]
			}
		}
	}
	var rows [][]interface{}
	for _, el := range arr {
		mo, ok := el.(map[string]interface{})
		if !ok {
			return value.Nil, runtime.Errorf("TABLE.FROMJSON: array elements must be objects")
		}
		row := make([]interface{}, len(cols))
		for i, c := range cols {
			row[i] = mo[c]
		}
		rows = append(rows, row)
	}
	id, err := m.h.Alloc(&tableObj{cols: cols, rows: rows})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func tabToCSV(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("TABLE.TOCSV expects handle")
	}
	t, err := castTable(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := t.assertLive(); err != nil {
		return value.Nil, err
	}
	var rows [][]string
	rows = append(rows, append([]string(nil), t.cols...))
	for _, row := range t.rows {
		rec := make([]string, len(t.cols))
		for i, c := range row {
			if i < len(rec) {
				rec[i] = fmt.Sprint(c)
			}
		}
		rows = append(rows, rec)
	}
	return mbcsv.AllocCSV(m.h, rows)
}

func tabFromCSV(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("TABLE.FROMCSV expects csv handle")
	}
	raw, err := mbcsv.Rows(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if len(raw) < 2 {
		id, e := m.h.Alloc(&tableObj{cols: nil, rows: nil})
		if e != nil {
			return value.Nil, e
		}
		return value.FromHandle(id), nil
	}
	cols := append([]string(nil), raw[0]...)
	var rows [][]interface{}
	for r := 1; r < len(raw); r++ {
		line := raw[r]
		row := make([]interface{}, len(cols))
		for i := range cols {
			if i < len(line) {
				row[i] = line[i]
			} else {
				row[i] = ""
			}
		}
		rows = append(rows, row)
	}
	id, err := m.h.Alloc(&tableObj{cols: cols, rows: rows})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}
