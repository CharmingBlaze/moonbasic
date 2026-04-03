package codegen

import (
	"strings"
	"testing"

	"moonbasic/compiler/parser"
	"moonbasic/vm/opcode"
)

func TestCompilePrintLiteral(t *testing.T) {
	src := "PRINT(42)\n"
	prog, err := parser.ParseSource("t.mbc", src)
	if err != nil {
		t.Fatal(err)
	}
	g := New("t.mbc", parser.SplitLines(src))
	out, err := g.Compile(prog)
	if err != nil {
		t.Fatal(err)
	}
	if out == nil || out.Main == nil {
		t.Fatal("nil program or main chunk")
	}
	code := out.Main.Instructions
	if len(code) < 2 {
		t.Fatalf("expected instructions, got %d", len(code))
	}
	if code[0].Op != opcode.OpPushInt || code[1].Op != opcode.OpCallBuiltin {
		t.Fatalf("unexpected ops: %v", code)
	}
}

func TestCompileLogicalOrAndXor(t *testing.T) {
	src := "x? = TRUE OR FALSE\ny? = TRUE AND FALSE\nz? = TRUE XOR FALSE\n"
	lines := parser.SplitLines(src)
	tree, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	g := New("t.mb", lines)
	out, err := g.Compile(tree)
	if err != nil {
		t.Fatal(err)
	}
	d := out.Main.Disassemble()
	for _, op := range []string{"OR", "AND", "XOR"} {
		if !strings.Contains(d, op) {
			t.Fatalf("disassembly missing %s:\n%s", op, d)
		}
	}
}

func TestCompileAssignAdd(t *testing.T) {
	src := "x = 2 + 3\nPRINT(x)\n"
	lines := parser.SplitLines(src)
	tree, err := parser.ParseSource("t.mbc", src)
	if err != nil {
		t.Fatal(err)
	}
	g := New("t.mbc", lines)
	out, err := g.Compile(tree)
	if err != nil {
		t.Fatal(err)
	}
	d := out.Main.Disassemble()
	if !strings.Contains(d, "HALT") || !strings.Contains(d, "ADD") {
		t.Log(d)
		t.Fatal("expected folded or runtime add in disassembly")
	}
}
