//go:build cgo || (windows && !cgo)

package window

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/value"
)

// SeedWindowFlagGlobals installs FLAG_* raylib config flag bitmasks for WINDOW.SETFLAG / CHECKFLAG.
func SeedWindowFlagGlobals(globals map[string]value.Value) {
	if globals == nil {
		return
	}
	globals["FLAG_VSYNC_HINT"] = value.FromInt(int64(rl.FlagVsyncHint))
	globals["FLAG_FULLSCREEN_MODE"] = value.FromInt(int64(rl.FlagFullscreenMode))
	globals["FLAG_WINDOW_RESIZABLE"] = value.FromInt(int64(rl.FlagWindowResizable))
	globals["FLAG_WINDOW_UNDECORATED"] = value.FromInt(int64(rl.FlagWindowUndecorated))
	globals["FLAG_WINDOW_HIDDEN"] = value.FromInt(int64(rl.FlagWindowHidden))
	globals["FLAG_WINDOW_MINIMIZED"] = value.FromInt(int64(rl.FlagWindowMinimized))
	globals["FLAG_WINDOW_MAXIMIZED"] = value.FromInt(int64(rl.FlagWindowMaximized))
	globals["FLAG_WINDOW_UNFOCUSED"] = value.FromInt(int64(rl.FlagWindowUnfocused))
	globals["FLAG_WINDOW_TOPMOST"] = value.FromInt(int64(rl.FlagWindowTopmost))
	globals["FLAG_WINDOW_TRANSPARENT"] = value.FromInt(int64(rl.FlagWindowTransparent))
	globals["FLAG_WINDOW_HIGHDPI"] = value.FromInt(int64(rl.FlagWindowHighdpi))
	globals["FLAG_WINDOW_MOUSE_PASSTHROUGH"] = value.FromInt(int64(rl.FlagWindowMousePassthrough))
	globals["FLAG_MSAA_4X_HINT"] = value.FromInt(int64(rl.FlagMsaa4xHint))
}
