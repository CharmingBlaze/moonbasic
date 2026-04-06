//go:build cgo || (windows && !cgo)

package input

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/value"
)

// SeedGestureGlobals installs GESTURE_* Raylib gesture ids / enable flags.
func SeedGestureGlobals(globals map[string]value.Value) {
	if globals == nil {
		return
	}
	globals["GESTURE_NONE"] = value.FromInt(int64(rl.GestureNone))
	globals["GESTURE_TAP"] = value.FromInt(int64(rl.GestureTap))
	globals["GESTURE_DOUBLETAP"] = value.FromInt(int64(rl.GestureDoubletap))
	globals["GESTURE_HOLD"] = value.FromInt(int64(rl.GestureHold))
	globals["GESTURE_DRAG"] = value.FromInt(int64(rl.GestureDrag))
	globals["GESTURE_SWIPE_RIGHT"] = value.FromInt(int64(rl.GestureSwipeRight))
	globals["GESTURE_SWIPE_LEFT"] = value.FromInt(int64(rl.GestureSwipeLeft))
	globals["GESTURE_SWIPE_UP"] = value.FromInt(int64(rl.GestureSwipeUp))
	globals["GESTURE_SWIPE_DOWN"] = value.FromInt(int64(rl.GestureSwipeDown))
	globals["GESTURE_PINCH_IN"] = value.FromInt(int64(rl.GesturePinchIn))
	globals["GESTURE_PINCH_OUT"] = value.FromInt(int64(rl.GesturePinchOut))
}
