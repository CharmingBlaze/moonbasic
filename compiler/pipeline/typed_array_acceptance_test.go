package pipeline

import (
	"os"
	"path/filepath"
	"testing"
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
