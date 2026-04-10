package pipeline

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"moonbasic/runtime"
	"moonbasic/vm"
	"moonbasic/vm/heap"
)

func TestTypedArrayAcceptanceCompile(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "typed_array_acceptance.mb")
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = CompileSource("typed_array_acceptance.mb", string(src))
	if err != nil {
		t.Fatal(err)
	}
}

func TestArrayBoundsErrorMessage(t *testing.T) {
	src := `
TYPE Cell
    value AS INTEGER
    visited AS INTEGER
ENDTYPE
map AS Cell(10, 10, 5)
map(11, 1, 1).value = 0
`
	prog, err := CompileSource("bounds.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	prog.SourcePath = "bounds.mb"
	h := heap.New()
	reg := runtime.NewRegistry(h)
	reg.InitCore()
	v := vm.New(reg, h)
	err = v.Execute(prog)
	reg.Shutdown()
	if err == nil {
		t.Fatal("expected runtime error for out-of-bounds index")
	}
	s := err.Error()
	if !strings.Contains(s, "out of bounds") {
		t.Fatalf("expected 'out of bounds' in error, got:\n%s", s)
	}
}

func TestConvenienceCompile(t *testing.T) {
	src := `
CONST N = 8
vals AS INTEGER(N)
FOR i = 1 TO vals.length
    vals(i) = CLAMP(i * 2, 1, N)
NEXT i

x = 3.0
y = 4.0
dist = VEC2.LENGTH(x, y)
x, y = VEC2.NORMALIZE(x, y)
a, b, c = Entity.GetPos(1)
`
	if _, err := CompileSource("convenience_compile.mb", src); err != nil {
		t.Fatal(err)
	}
}
