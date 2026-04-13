//go:build cgo || linux

package pipeline

import (
	"strings"
	"testing"

	"moonbasic/runtime"
	"moonbasic/vm"
	"moonbasic/vm/heap"
)

// Exercising the VM + registry pulls in moonbasic/runtime. On Windows with CGO_ENABLED=0,
// raylib-go purego loads raylib.dll in init — exclude this file there unless CGO is on.
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
