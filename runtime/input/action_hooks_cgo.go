//go:build cgo

package input

import rl "github.com/gen2brain/raylib-go/raylib"

type actionQuery struct{}

func (actionQuery) keyPressed(k int32) bool {
	return rl.IsKeyPressed(rl.KeyboardKey(k))
}

func (actionQuery) keyDown(k int32) bool {
	return rl.IsKeyDown(rl.KeyboardKey(k))
}

func (actionQuery) keyReleased(k int32) bool {
	return rl.IsKeyReleased(rl.KeyboardKey(k))
}

func (actionQuery) gamepadBtnPressed(pad, btn int32) bool {
	if !rl.IsGamepadAvailable(pad) {
		return false
	}
	return rl.IsGamepadButtonPressed(pad, rl.GamepadButton(btn))
}

func (actionQuery) gamepadBtnDown(pad, btn int32) bool {
	if !rl.IsGamepadAvailable(pad) {
		return false
	}
	return rl.IsGamepadButtonDown(pad, rl.GamepadButton(btn))
}

func (actionQuery) gamepadBtnReleased(pad, btn int32) bool {
	if !rl.IsGamepadAvailable(pad) {
		return false
	}
	return rl.IsGamepadButtonReleased(pad, rl.GamepadButton(btn))
}

func (actionQuery) gamepadAxis(pad, axis int32) float32 {
	if !rl.IsGamepadAvailable(pad) {
		return 0
	}
	return rl.GetGamepadAxisMovement(pad, rl.GamepadAxis(axis))
}

func actionQueries() actionQuery { return actionQuery{} }
