//go:build !cgo && !windows

package window

import "moonbasic/vm/value"

// SeedWindowFlagGlobals installs FLAG_* values matching raylib (for semantic check / stub runs).
func SeedWindowFlagGlobals(globals map[string]value.Value) {
	if globals == nil {
		return
	}
	globals["FLAG_VSYNC_HINT"] = value.FromInt(0x00000040)
	globals["FLAG_FULLSCREEN_MODE"] = value.FromInt(0x00000002)
	globals["FLAG_WINDOW_RESIZABLE"] = value.FromInt(0x00000004)
	globals["FLAG_WINDOW_UNDECORATED"] = value.FromInt(0x00000008)
	globals["FLAG_WINDOW_HIDDEN"] = value.FromInt(0x00000080)
	globals["FLAG_WINDOW_MINIMIZED"] = value.FromInt(0x00000200)
	globals["FLAG_WINDOW_MAXIMIZED"] = value.FromInt(0x00000400)
	globals["FLAG_WINDOW_UNFOCUSED"] = value.FromInt(0x00000800)
	globals["FLAG_WINDOW_TOPMOST"] = value.FromInt(0x00001000)
	globals["FLAG_WINDOW_TRANSPARENT"] = value.FromInt(0x00000010)
	globals["FLAG_WINDOW_HIGHDPI"] = value.FromInt(0x00002000)
	globals["FLAG_WINDOW_MOUSE_PASSTHROUGH"] = value.FromInt(0x00004000)
	globals["FLAG_MSAA_4X_HINT"] = value.FromInt(0x00000020)
}
