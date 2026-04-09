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
	globals["KEY_1"] = value.FromInt(49)
	globals["KEY_2"] = value.FromInt(50)
	globals["KEY_3"] = value.FromInt(51)
	globals["KEY_4"] = value.FromInt(52)
	globals["KEY_5"] = value.FromInt(53)
	globals["KEY_6"] = value.FromInt(54)
	globals["KEY_W"] = value.FromInt(87)
	globals["KEY_A"] = value.FromInt(65)
	globals["KEY_S"] = value.FromInt(83)
	globals["KEY_D"] = value.FromInt(68)
	globals["KEY_Q"] = value.FromInt(81)
	globals["KEY_E"] = value.FromInt(69)
	globals["KEY_I"] = value.FromInt(73)
	globals["KEY_K"] = value.FromInt(75)
	globals["KEY_F1"] = value.FromInt(290)
	globals["KEY_F2"] = value.FromInt(291)
	globals["KEY_F3"] = value.FromInt(292)
	globals["KEY_F4"] = value.FromInt(293)
	globals["KEY_F5"] = value.FromInt(294)
	globals["KEY_F6"] = value.FromInt(295)
	globals["KEY_F7"] = value.FromInt(296)
	globals["KEY_F8"] = value.FromInt(297)
	globals["KEY_F9"] = value.FromInt(298)
	globals["KEY_F10"] = value.FromInt(299)
	globals["KEY_F11"] = value.FromInt(300)
	globals["KEY_F12"] = value.FromInt(301)
	// Raylib KeyboardKey — arrows (INPUT.ActionAxis digital defaults use Left/A and Right/D).

	globals["KEY_LEFT"] = value.FromInt(263)
	globals["KEY_RIGHT"] = value.FromInt(262)
	globals["KEY_UP"] = value.FromInt(265)
	globals["KEY_DOWN"] = value.FromInt(264)
	// Gamepad — numeric values match raylib GamepadAxis / GamepadButton (CGO builds).
	globals["GAMEPAD_AXIS_LEFT_X"] = value.FromInt(0)
	globals["GAMEPAD_AXIS_LEFT_Y"] = value.FromInt(1)
	globals["GAMEPAD_BUTTON_RIGHT_FACE_DOWN"] = value.FromInt(2)
	globals["GAMEPAD_BUTTON_RIGHT_FACE_RIGHT"] = value.FromInt(3)
	globals["GAMEPAD_BUTTON_RIGHT_FACE_LEFT"] = value.FromInt(4)
	globals["GAMEPAD_BUTTON_RIGHT_FACE_UP"] = value.FromInt(5)
}
