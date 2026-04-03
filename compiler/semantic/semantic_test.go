package semantic

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"moonbasic/compiler/ast"
	"moonbasic/compiler/parser"
)

func TestFoldAssignInt(t *testing.T) {
	prog, err := parser.ParseSource("t.mbc", "x = 5 + 10\n")
	if err != nil {
		t.Fatal(err)
	}
	FoldConstants(prog)
	as := prog.Stmts[0].(*ast.AssignNode)
	lit, ok := as.Expr.(*ast.IntLitNode)
	if !ok || lit.Value != 15 {
		t.Fatalf("expected IntLit 15, got %#v", as.Expr)
	}
}

func TestTypeCheckSetFPSString(t *testing.T) {
	prog, err := parser.ParseSource("t.mbc", `RENDER.SETFPS("Fast")`+"\n")
	if err != nil {
		t.Fatal(err)
	}
	a := DefaultAnalyzer("t.mbc", parser.SplitLines(`RENDER.SETFPS("Fast")`))
	if err := a.Run(prog); err == nil {
		t.Fatal("expected type error")
	}
}

func TestUnknownEngineCommandRejected(t *testing.T) {
	src := "FOO.BAR()\n"
	prog, err := parser.ParseSource("t.mbc", src)
	if err != nil {
		t.Fatal(err)
	}
	a := DefaultAnalyzer("t.mbc", parser.SplitLines(src))
	if err := a.Run(prog); err == nil {
		t.Fatal("expected error for unknown FOO.BAR")
	}
}

func TestTypeCheckSetFPSInt(t *testing.T) {
	prog, err := parser.ParseSource("t.mbc", "RENDER.SETFPS(60)\n")
	if err != nil {
		t.Fatal(err)
	}
	a := DefaultAnalyzer("t.mbc", parser.SplitLines("RENDER.SETFPS(60)"))
	if err := a.Run(prog); err != nil {
		t.Fatal(err)
	}
}

func TestReferenceSemantic(t *testing.T) {
	src := readReference(t)
	lines := parser.SplitLines(src)
	prog, err := parser.ParseSource("reference.mbc", src)
	if err != nil {
		t.Fatal(err)
	}
	a := DefaultAnalyzer("reference.mbc", lines)
	if err := a.Run(prog); err != nil {
		t.Fatal(err)
	}
}

func readReference(t *testing.T) string {
	t.Helper()
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	p := filepath.Join(dir, "..", "..", "testdata", "reference.mbc")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read reference: %v", err)
	}
	return string(b)
}
