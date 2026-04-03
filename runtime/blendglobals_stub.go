//go:build !cgo

package runtime

import "moonbasic/vm/value"

// SeedBlendModeGlobals is a no-op without Raylib (CGO).
func SeedBlendModeGlobals(globals map[string]value.Value) {}
