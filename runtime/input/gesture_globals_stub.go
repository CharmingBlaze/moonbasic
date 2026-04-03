//go:build !cgo

package input

import "moonbasic/vm/value"

// SeedGestureGlobals installs GESTURE_* values matching raylib (for semantic check / stub runs).
func SeedGestureGlobals(globals map[string]value.Value) {
	if globals == nil {
		return
	}
	globals["GESTURE_NONE"] = value.FromInt(0)
	globals["GESTURE_TAP"] = value.FromInt(1)
	globals["GESTURE_DOUBLETAP"] = value.FromInt(2)
	globals["GESTURE_HOLD"] = value.FromInt(4)
	globals["GESTURE_DRAG"] = value.FromInt(8)
	globals["GESTURE_SWIPE_RIGHT"] = value.FromInt(16)
	globals["GESTURE_SWIPE_LEFT"] = value.FromInt(32)
	globals["GESTURE_SWIPE_UP"] = value.FromInt(64)
	globals["GESTURE_SWIPE_DOWN"] = value.FromInt(128)
	globals["GESTURE_PINCH_IN"] = value.FromInt(256)
	globals["GESTURE_PINCH_OUT"] = value.FromInt(512)
}
