//go:build unix && !cgo

package window

import (
	"testing"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func TestStubOpenReturnsHint(t *testing.T) {
	reg := runtime.NewRegistryHeadless(heap.New())
	NewModule().Register(reg)
	_, err := reg.Call("WINDOW.OPEN", []value.Value{
		value.FromInt(800),
		value.FromInt(600),
		value.FromStringIndex(0),
	})
	if err == nil {
		t.Fatal("expected error without CGO")
	}
}
