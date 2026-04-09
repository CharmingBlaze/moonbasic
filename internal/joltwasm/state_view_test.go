package joltwasm

import (
	"testing"
)

func TestFloatsView(t *testing.T) {
	b := []byte{
		0, 0, 0x80, 0x3f, // 1.0f
		0, 0, 0, 0x40, // 2.0f
	}
	f := FloatsView(b)
	if len(f) != 2 || f[0] != 1 || f[1] != 2 {
		t.Fatalf("got %v", f)
	}
	if FloatsView([]byte{1, 2, 3}) != nil {
		t.Fatal("expected nil for bad length")
	}
}
