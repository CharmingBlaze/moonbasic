package pipeline

import (
	"encoding/json"
	"testing"
)

func TestExportSymbolTableJSON(t *testing.T) {
	src := `x = 1.0
FUNCTION FOO()
  y = 2.0
ENDFUNCTION
`
	raw, err := ExportSymbolTableJSON("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	var v map[string]any
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Fatal(err)
	}
	if v["globals"] == nil {
		t.Fatal("expected globals in export")
	}
	if v["path"] != "t.mb" {
		t.Fatalf("expected path t.mb, got %v", v["path"])
	}
}

func TestDocumentSymbols(t *testing.T) {
	src := `CONST N = 10
TYPE T
  X
ENDTYPE
`
	syms, err := DocumentSymbols("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	if len(syms) < 2 {
		t.Fatalf("expected multiple symbols, got %d", len(syms))
	}
}
