package include

import (
	"os"
	"path/filepath"
	"testing"

	"moonbasic/compiler/arena"
	"moonbasic/compiler/parser"
)

func TestExpand_duplicateIncludeSkipped(t *testing.T) {
	dir := t.TempDir()
	lib := filepath.Join(dir, "lib.mb")
	if err := os.WriteFile(lib, []byte("PRINT(1)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mainSrc := `INCLUDE "lib.mb"
INCLUDE "lib.mb"
`
	mainPath := filepath.Join(dir, "main.mb")
	ar := arena.NewArena()
	defer ar.Reset()
	prog, err := parser.ParseSourceWithArena(mainPath, mainSrc, ar)
	if err != nil {
		t.Fatal(err)
	}
	out, err := ExpandWithArena(mainPath, prog, ar)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Stmts) != 1 {
		t.Fatalf("expected 1 statement after merge, got %d (duplicate INCLUDE should not duplicate body)", len(out.Stmts))
	}
}
