//go:build cgo

package runtime

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"moonbasic/vm/value"
)

// SeedBlendModeGlobals installs BLEND_* constants matching raylib BlendMode.
func SeedBlendModeGlobals(globals map[string]value.Value) {
	if globals == nil {
		return
	}
	globals["BLEND_ALPHA"] = value.FromInt(int64(rl.BlendAlpha))
	globals["BLEND_ADDITIVE"] = value.FromInt(int64(rl.BlendAdditive))
	globals["BLEND_MULTIPLIED"] = value.FromInt(int64(rl.BlendMultiplied))
	globals["BLEND_ADD_COLORS"] = value.FromInt(int64(rl.BlendAddColors))
	globals["BLEND_SUBTRACT_COLORS"] = value.FromInt(int64(rl.BlendSubtractColors))
	globals["BLEND_ALPHA_PREMULP"] = value.FromInt(int64(rl.BlendAlphaPremultiply))
	globals["BLEND_CUSTOM"] = value.FromInt(int64(rl.BlendCustom))
}
