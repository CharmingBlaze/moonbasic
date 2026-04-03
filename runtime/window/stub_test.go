//go:build !cgo

package window

import (
	"strings"
	"testing"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func TestStubRegistersCommands(t *testing.T) {
	reg := runtime.NewRegistry(heap.New())
	NewModule().Register(reg)
	for _, key := range []string{"WINDOW.OPEN", "WINDOW.SETFPS", "WINDOW.CLOSE", "WINDOW.SHOULDCLOSE", "RENDER.CLEAR", "RENDER.FRAME"} {
		_, err := reg.Call(key, nil)
		if err != nil && strings.Contains(err.Error(), "unknown command") {
			t.Fatalf("missing registration for %s: %v", key, err)
		}
	}
}

func TestStubOpenReturnsHint(t *testing.T) {
	reg := runtime.NewRegistry(heap.New())
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
