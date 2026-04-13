package hal

// InputDevice abstracts keyboard, mouse, and gamepad polling.
type InputDevice interface {
	IsKeyDown(key int32) bool
	IsKeyPressed(key int32) bool
	IsMouseButtonPressed(button int32) bool
	GetMousePosition() V2
	GetMouseWheelMove() float32
	IsGamepadAvailable(id int32) bool
	GetGamepadAxisMovement(id, axis int32) float32
}
