package window

import (
	"strings"
	"testing"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

func TestStubRegistersCommands(t *testing.T) {
	reg := runtime.NewRegistryHeadless(heap.New())
	NewModule().Register(reg)
	for _, key := range []string{"WINDOW.OPEN", "WINDOW.SETFPS", "WINDOW.CLOSE", "WINDOW.SHOULDCLOSE", "RENDER.CLEAR", "RENDER.FRAME"} {
		_, err := reg.Call(key, nil)
		if err != nil && strings.Contains(err.Error(), "unknown command") {
			t.Fatalf("missing registration for %s: %v", key, err)
		}
	}
}
