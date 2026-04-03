package runtime

import "moonbasic/vm/value"

// SeedInputKeyGlobals installs uppercase KEY_* identifiers into the VM global map
// so scripts can write Input.KeyDown(KEY_ESCAPE). Values match raylib KeyboardKey
// (github.com/gen2brain/raylib-go/raylib) for CGO builds.
func SeedInputKeyGlobals(globals map[string]value.Value) {
	if globals == nil {
		return
	}
	// Subset — extend as INPUT.KEYDOWN coverage grows (see raylib KeyboardKey enum).
	globals["KEY_ESCAPE"] = value.FromInt(256)
	globals["KEY_SPACE"] = value.FromInt(32)
	globals["KEY_W"] = value.FromInt(87)
	globals["KEY_A"] = value.FromInt(65)
	globals["KEY_S"] = value.FromInt(83)
	globals["KEY_D"] = value.FromInt(68)
	globals["KEY_I"] = value.FromInt(73)
	globals["KEY_K"] = value.FromInt(75)
}
