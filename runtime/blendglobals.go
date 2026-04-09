package runtime

import "moonbasic/vm/value"

// Raylib BlendMode enum (raylib 5.x) — inlined so package runtime does not import
// github.com/gen2brain/raylib-go/raylib (that init() loads raylib.dll and breaks headless tests/CI).
const (
	blendAlpha            int64 = 0
	blendAdditive         int64 = 1
	blendMultiplied       int64 = 2
	blendAddColors        int64 = 3
	blendSubtractColors   int64 = 4
	blendAlphaPremultiply int64 = 5
	blendCustom           int64 = 6
)

// SeedBlendModeGlobals installs BLEND_* constants matching raylib BlendMode.
func SeedBlendModeGlobals(globals map[string]value.Value) {
	if globals == nil {
		return
	}
	globals["BLEND_ALPHA"] = value.FromInt(blendAlpha)
	globals["BLEND_ADDITIVE"] = value.FromInt(blendAdditive)
	globals["BLEND_MULTIPLIED"] = value.FromInt(blendMultiplied)
	globals["BLEND_ADD_COLORS"] = value.FromInt(blendAddColors)
	globals["BLEND_SUBTRACT_COLORS"] = value.FromInt(blendSubtractColors)
	globals["BLEND_ALPHA_PREMULP"] = value.FromInt(blendAlphaPremultiply)
	globals["BLEND_CUSTOM"] = value.FromInt(blendCustom)
}
