package mbjson

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func jToString(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.TOSTRING: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("JSON.TOSTRING expects handle")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := j.assertLive(); err != nil {
		return value.Nil, err
	}
	b, err := json.Marshal(j.root)
	if err != nil {
		return value.Nil, err
	}
	return rt.RetString(string(b)), nil
}

func jPretty(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.PRETTY: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("JSON.PRETTY expects handle")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := j.assertLive(); err != nil {
		return value.Nil, err
	}
	b, err := json.MarshalIndent(j.root, "", "  ")
	if err != nil {
		return value.Nil, err
	}
	return rt.RetString(string(b)), nil
}

func jMinify(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	// Same as TOSTRING — compact JSON text.
	return jToString(m, rt, args...)
}

func jToFile(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.TOFILE: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.TOFILE expects (handle, path$)")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := j.assertLive(); err != nil {
		return value.Nil, err
	}
	path, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	b, err := json.Marshal(j.root)
	if err != nil {
		return value.Nil, err
	}
	if err := os.WriteFile(strings.TrimSpace(path), b, 0o644); err != nil {
		return value.Nil, err
	}
	return args[0], nil
}

func jToFilePretty(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.TOFILEPRETTY: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("JSON.TOFILEPRETTY expects (handle, path$)")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := j.assertLive(); err != nil {
		return value.Nil, err
	}
	path, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	b, err := json.MarshalIndent(j.root, "", "  ")
	if err != nil {
		return value.Nil, err
	}
	if err := os.WriteFile(strings.TrimSpace(path), b, 0o644); err != nil {
		return value.Nil, err
	}
	return args[0], nil
}

// jToCSV exports a JSON array of objects to CSV text (RFC 4180 via encoding/csv).
// Optional path selects a sub-array; empty path uses root (must be array of objects).
func jToCSV(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("JSON.TOCSV: heap not bound")
	}
	if len(args) != 1 && len(args) != 2 {
		return value.Nil, runtime.Errorf("JSON.TOCSV expects (handle) or (handle, path$)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("JSON.TOCSV: first argument must be handle")
	}
	j, err := castJSON(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if err := j.assertLive(); err != nil {
		return value.Nil, err
	}
	var arr []interface{}
	if len(args) == 1 {
		a, ok := j.root.([]interface{})
		if !ok {
			return value.Nil, runtime.Errorf("JSON.TOCSV: root must be array of objects")
		}
		arr = a
	} else {
		if args[1].Kind != value.KindString {
			return value.Nil, runtime.Errorf("JSON.TOCSV: path must be string")
		}
		segs, err := pathFromArgs(rt, args, 1)
		if err != nil {
			return value.Nil, err
		}
		v, ok := getValue(j.root, segs)
		if !ok {
			return value.Nil, runtime.Errorf("JSON.TOCSV: path not found")
		}
		a, ok := v.([]interface{})
		if !ok {
			return value.Nil, runtime.Errorf("JSON.TOCSV: path must be array")
		}
		arr = a
	}
	if len(arr) == 0 {
		return rt.RetString(""), nil
	}
	var headers []string
	keySet := map[string]bool{}
	for _, row := range arr {
		mo, ok := row.(map[string]interface{})
		if !ok {
			return value.Nil, runtime.Errorf("JSON.TOCSV: each row must be object")
		}
		for k := range mo {
			if !keySet[k] {
				keySet[k] = true
				headers = append(headers, k)
			}
		}
	}
	for i := 0; i < len(headers); i++ {
		for j := i + 1; j < len(headers); j++ {
			if headers[j] < headers[i] {
				headers[i], headers[j] = headers[j], headers[i]
			}
		}
	}
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.Write(headers); err != nil {
		return value.Nil, err
	}
	for _, row := range arr {
		mo := row.(map[string]interface{})
		rec := make([]string, len(headers))
		for i, h := range headers {
			rec[i] = cellString(mo[h])
		}
		if err := w.Write(rec); err != nil {
			return value.Nil, err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return value.Nil, err
	}
	return rt.RetString(buf.String()), nil
}

func cellString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case bool:
		if t {
			return "true"
		}
		return "false"
	case float64:
		return fmt.Sprint(t)
	case json.Number:
		return t.String()
	default:
		b, _ := json.Marshal(t)
		return string(b)
	}
}
