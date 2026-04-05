package input

import (
	"fmt"

	"moonbasic/vm/value"
)

// KeyCodeFromValue parses a numeric keyboard scancode (e.g. KEY_ESCAPE) from a VM value.
// Shared by INPUT.* and instant-game shortcuts so key handling stays in one place.
func KeyCodeFromValue(v value.Value) (int32, error) {
	if i, ok := v.ToInt(); ok {
		return int32(i), nil
	}
	if f, ok := v.ToFloat(); ok {
		return int32(f), nil
	}
	return 0, fmt.Errorf("expected numeric key code (use KEY_ESCAPE etc.)")
}
