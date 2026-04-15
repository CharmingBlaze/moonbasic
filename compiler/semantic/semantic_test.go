package semantic

import (
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"moonbasic/compiler/ast"
	"moonbasic/compiler/parser"
)

func TestFoldTurnEntityGroupedThirdArg(t *testing.T) {
	src := "TURNENTITY(pivot, 0, (1.0 + 0.5), 0)\n"
	prog, err := parser.ParseSource("t.mbc", src)
	if err != nil {
		t.Fatal(err)
	}
	FoldConstants(prog)
	cs := prog.Stmts[0].(*ast.CallStmtNode)
	if cs.Name != "TURNENTITY" {
		t.Fatalf("command name: got %q", cs.Name)
	}
	if len(cs.Args) != 4 {
		t.Fatalf("args: got %d", len(cs.Args))
	}
	fl, ok := cs.Args[2].(*ast.FloatLitNode)
	if !ok {
		t.Fatalf("third arg after fold: want *FloatLitNode, got %T %#v", cs.Args[2], cs.Args[2])
	}
	if math.Abs(fl.Value-1.5) > 1e-12 {
		t.Fatalf("third arg value: want 1.5, got %v", fl.Value)
	}
}

func TestParseBareDrawEntitiesEqualsDrawEntitiesParens(t *testing.T) {
	src := "DrawEntities\nDrawEntities()\n"
	prog, err := parser.ParseSource("t.mbc", src)
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Stmts) != 2 {
		t.Fatalf("stmts: %d", len(prog.Stmts))
	}
	a, ok := prog.Stmts[0].(*ast.CallStmtNode)
	if !ok || a.Name != "DRAWENTITIES" || len(a.Args) != 0 {
		t.Fatalf("stmt0: %#v", prog.Stmts[0])
	}
	b, ok := prog.Stmts[1].(*ast.CallStmtNode)
	if !ok || b.Name != "DRAWENTITIES" || len(b.Args) != 0 {
		t.Fatalf("stmt1: %#v", prog.Stmts[1])
	}
}

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

func TestEntitySpatialLiteralOOB(t *testing.T) {
	src := "ENTITY.X(33554432) = 1.0\n"
	prog, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	a := DefaultAnalyzer("t.mb", parser.SplitLines(src))
	if err := a.Run(prog); err == nil {
		t.Fatal("expected type error for ENTITY spatial literal above MaxEntitySpatialIndex")
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

func TestDeprecatedMakeProducesWarning(t *testing.T) {
	src := "cam = CAMERA.MAKE()\n"
	prog, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	a := DefaultAnalyzer("t.mb", parser.SplitLines(src))
	if err := a.Run(prog); err != nil {
		t.Fatal(err)
	}
	ws := a.DeprecationWarnings()
	if len(ws) == 0 {
		t.Fatal("expected deprecation warning for CAMERA.MAKE")
	}
	if !strings.Contains(ws[0], "use CAMERA.CREATE") {
		t.Fatalf("unexpected warning: %q", ws[0])
	}
}

func TestStrictDeprecatedErrorsOnMake(t *testing.T) {
	src := "CAMERA.MAKE()\n"
	prog, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	a := DefaultAnalyzer("t.mb", parser.SplitLines(src))
	a.StrictDeprecated = true
	if err := a.Run(prog); err == nil {
		t.Fatal("expected strict deprecation error for CAMERA.MAKE")
	}
}

func TestDeprecatedSetPositionProducesWarning(t *testing.T) {
	src := "c = CAMERA.CREATE()\nCAMERA.SETPOSITION(c, 0, 1, 2)\n"
	prog, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	a := DefaultAnalyzer("t.mb", parser.SplitLines(src))
	if err := a.Run(prog); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, w := range a.DeprecationWarnings() {
		if strings.Contains(w, "CAMERA.SETPOSITION") && strings.Contains(w, "CAMERA.SETPOS") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected SETPOSITION deprecation warning, got: %#v", a.DeprecationWarnings())
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
