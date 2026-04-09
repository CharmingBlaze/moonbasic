package symtable

import (
	"testing"

	"moonbasic/compiler/ast"
	"moonbasic/compiler/types"
)

func TestBuilderImplicitGlobalFloat(t *testing.T) {
	b := NewBuilder()
	prog := &ast.Program{
		Stmts: []ast.Stmt{
			&ast.AssignNode{Name: "N", Expr: &ast.IntLitNode{Value: 42}},
		},
	}
	tab := b.Build(prog)
	s := tab.Resolve("N")
	if s == nil {
		t.Fatal("expected N")
	}
	if s.Type != types.Float {
		t.Fatalf("implicit unsuffixed global from int literal should infer float, got %v", s.Type)
	}
	if !s.Persistent {
		t.Fatal("implicit global N should be marked Persistent")
	}
}
