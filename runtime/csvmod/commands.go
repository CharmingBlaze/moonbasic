package mbcsv

import (
	"bytes"
	"encoding/csv"
	"os"
	"strings"

	mbjson "moonbasic/runtime/jsonmod"
	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func castCSV(m *Module, h heap.Handle) (*csvObj, error) {
	return heap.Cast[*csvObj](m.h, h)
}

func registerCSVCommands(m *Module, r runtime.Registrar) {
	r.Register("CSV.LOAD", "csv", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return csvLoad(m, rt, args...) })
	r.Register("CSV.SAVE", "csv", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return csvSave(m, rt, args...) })
	r.Register("CSV.FROMSTRING", "csv", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return csvFromString(m, rt, args...) })
	r.Register("CSV.TOSTRING", "csv", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return csvToString(m, rt, args...) })
	r.Register("CSV.FREE", "csv", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return csvFree(m, rt, args...) })
	r.Register("CSV.ROWCOUNT", "csv", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return csvRowCount(m, rt, args...) })
	r.Register("CSV.COLCOUNT", "csv", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return csvColCount(m, rt, args...) })
	r.Register("CSV.GET", "csv", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return csvGet(m, rt, args...) })
	r.Register("CSV.SET", "csv", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return csvSet(m, rt, args...) })
	r.Register("CSV.TOJSON", "csv", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return csvToJSON(m, rt, args...) })
}

func csvLoad(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CSV.LOAD: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("CSV.LOAD expects path$")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	data, err := os.ReadFile(strings.TrimSpace(path))
	if err != nil {
		return value.Nil, err
	}
	return parseCSVBytes(m, data)
}

func csvFromString(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CSV.FROMSTRING: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("CSV.FROMSTRING expects string")
	}
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	return parseCSVBytes(m, []byte(s))
}

func parseCSVBytes(m *Module, data []byte) (value.Value, error) {
	r := csv.NewReader(bytes.NewReader(data))
	r.FieldsPerRecord = -1
	rows, err := r.ReadAll()
	if err != nil {
		return value.Nil, err
	}
	id, err := m.h.Alloc(&csvObj{rows: rows})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func csvSave(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CSV.SAVE: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("CSV.SAVE expects (handle, path$)")
	}
	c, err := castCSV(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := c.assertLive(); err != nil {
		return value.Nil, err
	}
	path, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.WriteAll(c.rows); err != nil {
		return value.Nil, err
	}
	w.Flush()
	if err := os.WriteFile(strings.TrimSpace(path), buf.Bytes(), 0o644); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func csvToString(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CSV.TOSTRING: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("CSV.TOSTRING expects handle")
	}
	c, err := castCSV(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := c.assertLive(); err != nil {
		return value.Nil, err
	}
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.WriteAll(c.rows); err != nil {
		return value.Nil, err
	}
	w.Flush()
	return rt.RetString(buf.String()), nil
}

func csvFree(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CSV.FREE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("CSV.FREE expects handle")
	}
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func csvRowCount(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CSV.ROWCOUNT: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("CSV.ROWCOUNT expects handle")
	}
	c, err := castCSV(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := c.assertLive(); err != nil {
		return value.Nil, err
	}
	return value.FromInt(int64(len(c.rows))), nil
}

func csvColCount(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CSV.COLCOUNT: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("CSV.COLCOUNT expects handle")
	}
	c, err := castCSV(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := c.assertLive(); err != nil {
		return value.Nil, err
	}
	if len(c.rows) == 0 {
		return value.FromInt(0), nil
	}
	return value.FromInt(int64(len(c.rows[0]))), nil
}

// Row and column indices are 1-based (first data row is 1).
func csvGet(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CSV.GET: heap not bound")
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("CSV.GET expects (handle, row, col)")
	}
	c, err := castCSV(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := c.assertLive(); err != nil {
		return value.Nil, err
	}
	row, ok1 := args[1].ToInt()
	col, ok2 := args[2].ToInt()
	if !ok1 {
		if f, okf := args[1].ToFloat(); okf {
			row = int64(f)
			ok1 = true
		}
	}
	if !ok2 {
		if f, okf := args[2].ToFloat(); okf {
			col = int64(f)
			ok2 = true
		}
	}
	if !ok1 || !ok2 {
		return value.Nil, runtime.Errorf("CSV.GET: row/col must be numeric")
	}
	ri, ci := int(row)-1, int(col)-1
	if ri < 0 || ri >= len(c.rows) || ci < 0 || ci >= len(c.rows[ri]) {
		return rt.RetString(""), nil
	}
	return rt.RetString(c.rows[ri][ci]), nil
}

func csvSet(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CSV.SET: heap not bound")
	}
	if len(args) != 4 || args[0].Kind != value.KindHandle || args[3].Kind != value.KindString {
		return value.Nil, runtime.Errorf("CSV.SET expects (handle, row, col, val$)")
	}
	c, err := castCSV(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := c.assertLive(); err != nil {
		return value.Nil, err
	}
	row, ok1 := args[1].ToInt()
	col, ok2 := args[2].ToInt()
	if !ok1 {
		if f, okf := args[1].ToFloat(); okf {
			row = int64(f)
			ok1 = true
		}
	}
	if !ok2 {
		if f, okf := args[2].ToFloat(); okf {
			col = int64(f)
			ok2 = true
		}
	}
	if !ok1 || !ok2 {
		return value.Nil, runtime.Errorf("CSV.SET: row/col must be numeric")
	}
	val, err := rt.ArgString(args, 3)
	if err != nil {
		return value.Nil, err
	}
	ri, ci := int(row)-1, int(col)-1
	if ri < 0 || ri >= len(c.rows) {
		return value.Nil, runtime.Errorf("CSV.SET: row out of range")
	}
	for len(c.rows[ri]) <= ci {
		c.rows[ri] = append(c.rows[ri], "")
	}
	c.rows[ri][ci] = val
	return value.Nil, nil
}

// CSV.TOJSON builds a JSON array of objects using row 1 as field names; data rows start at row 2.
func csvToJSON(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CSV.TOJSON: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("CSV.TOJSON expects handle")
	}
	c, err := castCSV(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := c.assertLive(); err != nil {
		return value.Nil, err
	}
	if len(c.rows) < 2 {
		return mbjson.AllocJSONRoot(m.h, []interface{}{})
	}
	headers := c.rows[0]
	var out []interface{}
	for r := 1; r < len(c.rows); r++ {
		row := c.rows[r]
		mo := make(map[string]interface{})
		for i, name := range headers {
			if name == "" {
				continue
			}
			if i < len(row) {
				mo[name] = row[i]
			} else {
				mo[name] = ""
			}
		}
		out = append(out, mo)
	}
	return mbjson.AllocJSONRoot(m.h, out)
}
