//go:build cgo || windows

package raylib

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"moonbasic/hal"
)

type Driver struct{}

func NewDriver() *Driver {
	return &Driver{}
}

// --- SystemDevice ---

func (d *Driver) InitWindow(width, height int, title string) {
	rl.InitWindow(int32(width), int32(height), title)
}

func (d *Driver) CloseWindow() {
	rl.CloseWindow()
}

func (d *Driver) WindowShouldClose() bool {
	return rl.WindowShouldClose()
}

func (d *Driver) SetTargetFPS(fps int) {
	rl.SetTargetFPS(int32(fps))
}

func (d *Driver) GetFPS() int {
	return int(rl.GetFPS())
}

func (d *Driver) GetFrameTime() float32 {
	return rl.GetFrameTime()
}

func (d *Driver) PollInputEvents() {
	rl.PollInputEvents()
}

func (d *Driver) SetWindowSize(width, height int) {
	rl.SetWindowSize(width, height)
}

func (d *Driver) GetScreenWidth() int {
	return int(rl.GetScreenWidth())
}

func (d *Driver) GetScreenHeight() int {
	return int(rl.GetScreenHeight())
}

// --- VideoDevice ---

func (d *Driver) BeginDrawing() {
	rl.BeginDrawing()
}

func (d *Driver) EndDrawing() {
	rl.EndDrawing()
}

func (d *Driver) ClearBackground(c hal.RGBA) {
	rl.ClearBackground(rl.Color{R: c.R, G: c.G, B: c.B, A: c.A})
}

func (d *Driver) DrawRectangle(x, y, w, h int32, c hal.RGBA) {
	rl.DrawRectangle(x, y, w, h, rl.Color{R: c.R, G: c.G, B: c.B, A: c.A})
}

func (d *Driver) DrawRectanglePro(rec hal.Rect, origin hal.V2, rotation float32, c hal.RGBA) {
	rl.DrawRectanglePro(
		rl.Rectangle{X: rec.X, Y: rec.Y, Width: rec.Width, Height: rec.Height},
		rl.Vector2{X: origin.X, Y: origin.Y},
		rotation,
		rl.Color{R: c.R, G: c.G, B: c.B, A: c.A},
	)
}

func (d *Driver) DrawRectangleLines(x, y, w, h int32, thick float32, c hal.RGBA) {
	rl.DrawRectangleLinesEx(
		rl.Rectangle{X: float32(x), Y: float32(y), Width: float32(w), Height: float32(h)},
		thick,
		rl.Color{R: c.R, G: c.G, B: c.B, A: c.A},
	)
}

func (d *Driver) DrawCircle(x, y int32, radius float32, c hal.RGBA) {
	rl.DrawCircle(x, y, radius, rl.Color{R: c.R, G: c.G, B: c.B, A: c.A})
}

func (d *Driver) DrawCircleLines(x, y int32, radius float32, c hal.RGBA) {
	rl.DrawCircleLines(x, y, radius, rl.Color{R: c.R, G: c.G, B: c.B, A: c.A})
}

func (d *Driver) DrawTriangle(v1, v2, v3 hal.V2, c hal.RGBA) {
	rl.DrawTriangle(
		rl.Vector2{X: v1.X, Y: v1.Y},
		rl.Vector2{X: v2.X, Y: v2.Y},
		rl.Vector2{X: v3.X, Y: v3.Y},
		rl.Color{R: c.R, G: c.G, B: c.B, A: c.A},
	)
}

func (d *Driver) DrawPoly(center hal.V2, sides int32, radius float32, rotation float32, c hal.RGBA) {
	rl.DrawPoly(
		rl.Vector2{X: center.X, Y: center.Y},
		sides,
		radius,
		rotation,
		rl.Color{R: c.R, G: c.G, B: c.B, A: c.A},
	)
}

func (d *Driver) DrawText(text string, x, y, size int32, c hal.RGBA) {
	rl.DrawText(text, x, y, size, rl.Color{R: c.R, G: c.G, B: c.B, A: c.A})
}

func (d *Driver) BeginMode3D(cam hal.Camera3D) {
	rl.BeginMode3D(rl.Camera3D{
		Position:   rl.Vector3{X: cam.Position.X, Y: cam.Position.Y, Z: cam.Position.Z},
		Target:     rl.Vector3{X: cam.Target.X, Y: cam.Target.Y, Z: cam.Target.Z},
		Up:         rl.Vector3{X: cam.Up.X, Y: cam.Up.Y, Z: cam.Up.Z},
		Fovy:       cam.Fovy,
		Projection: rl.CameraProjection(cam.Projection),
	})
}

func (d *Driver) EndMode3D() {
	rl.EndMode3D()
}

func (d *Driver) DrawCube(position hal.V3, width, height, length float32, c hal.RGBA) {
	rl.DrawCube(
		rl.Vector3{X: position.X, Y: position.Y, Z: position.Z},
		width, height, length,
		rl.Color{R: c.R, G: c.G, B: c.B, A: c.A},
	)
}

func (d *Driver) DrawGrid(slices int32, spacing float32) {
	rl.DrawGrid(slices, spacing)
}

// --- InputDevice ---

func (d *Driver) IsKeyDown(key int32) bool {
	return rl.IsKeyDown(key)
}

func (d *Driver) IsKeyPressed(key int32) bool {
	return rl.IsKeyPressed(key)
}

func (d *Driver) IsMouseButtonPressed(button int32) bool {
	return rl.IsMouseButtonPressed(rl.MouseButton(button))
}

func (d *Driver) GetMousePosition() hal.V2 {
	v := rl.GetMousePosition()
	return hal.V2{X: v.X, Y: v.Y}
}

func (d *Driver) GetMouseWheelMove() float32 {
	return rl.GetMouseWheelMove()
}

func (d *Driver) IsGamepadAvailable(id int32) bool {
	return rl.IsGamepadAvailable(id)
}

func (d *Driver) GetGamepadAxisMovement(id, axis int32) float32 {
	return rl.GetGamepadAxisMovement(id, axis)
}
