package moon

import (
	"testing"

	"moonbasic/compiler/codegen"
	"moonbasic/compiler/parser"
	"moonbasic/vm/opcode"
)

// IR v3 bytecode and MOON v3 containers are OS-agnostic (same layout on Windows and Linux).
// This test lives in package moon (no compiler/pipeline import) so it runs without linking
// the full runtime registry / Raylib — including local `go test ./vm/moon` on Windows without raylib.dll.

func TestIRV3MOONEncodeDecodeRoundTrip(t *testing.T) {
	src := "a = 2: b = 3\nx = a + b\n"
	prog := mustCompile(t, src)
	assertMoonRoundTrip(t, prog)
}

func TestIRV3MOONRoundTripCRLFSource(t *testing.T) {
	src := "a = 2: b = 3\r\nx = a + b\r\n"
	prog := mustCompile(t, src)
	assertMoonRoundTrip(t, prog)
}

func mustCompile(t *testing.T, src string) *opcode.Program {
	t.Helper()
	lines := parser.SplitLines(src)
	tree, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	g := codegen.New("t.mb", lines)
	prog, err := g.Compile(tree)
	if err != nil {
		t.Fatal(err)
	}
	if prog == nil || prog.Main == nil {
		t.Fatal("nil program")
	}
	return prog
}

func assertMoonRoundTrip(t *testing.T, want *opcode.Program) {
	t.Helper()
	data, err := Encode(want)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Main.Instructions) != len(want.Main.Instructions) {
		t.Fatalf("main instructions: got %d want %d", len(got.Main.Instructions), len(want.Main.Instructions))
	}
	for i := range got.Main.Instructions {
		gi, wi := got.Main.Instructions[i], want.Main.Instructions[i]
		if gi != wi {
			t.Fatalf("instr[%d]: got %+v want %+v", i, gi, wi)
		}
	}
}
