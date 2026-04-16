package parser

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestParseReference(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	refPath := filepath.Join(filepath.Dir(file), "..", "..", "testdata", "reference.mbc")
	src, err := os.ReadFile(refPath)
	if err != nil {
		t.Fatal(err)
	}
	prog, err := ParseSource(refPath, string(src))
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Stmts) < 10 {
		t.Fatalf("expected many top-level stmts, got %d", len(prog.Stmts))
	}
	if len(prog.Functions) != 1 || prog.Functions[0].Name != "onland" {
		t.Fatalf("functions: %+v", prog.Functions)
	}
}

func TestParseSimpleAssign(t *testing.T) {
	prog, err := ParseSource("x.mbc", "x = 1\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Stmts) != 1 {
		t.Fatal(len(prog.Stmts))
	}
}

func TestSingleLineIf(t *testing.T) {
	_, err := ParseSource("x.mbc", "IF TRUE THEN x = 1\n")
	if err != nil {
		t.Fatal(err)
	}
}

func TestColonSeparatedStatements(t *testing.T) {
	prog, err := ParseSource("x.mbc", "x = 1 : y = 2 : z = 3\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Stmts) != 3 {
		t.Fatalf("want 3 stmts, got %d", len(prog.Stmts))
	}
}

func TestColonSeparatedStatementsInWhileBody(t *testing.T) {
	src := "WHILE TRUE\n" +
		"    vx = 0 : vz = 0\n" +
		"WEND\n"
	prog, err := ParseSource("x.mbc", src)
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Stmts) != 1 {
		t.Fatalf("want 1 stmt, got %d", len(prog.Stmts))
	}
}
