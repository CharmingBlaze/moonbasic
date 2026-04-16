//go:build !cgo && !windows

package raylib

import "moonbasic/hal"

type Driver struct{}

func NewDriver() *Driver { return &Driver{} }

func (d *Driver) InitWindow(width, height int, title string) {}
func (d *Driver) CloseWindow()                                {}
func (d *Driver) WindowShouldClose() bool                     { return false }
func (d *Driver) SetTargetFPS(fps int)                        {}
func (d *Driver) GetFPS() int                                 { return 0 }
func (d *Driver) GetFrameTime() float32                       { return 0 }
func (d *Driver) PollInputEvents()                            {}
func (d *Driver) SetWindowSize(width, height int)             {}
func (d *Driver) GetScreenWidth() int                         { return 0 }
func (d *Driver) GetScreenHeight() int                        { return 0 }

func (d *Driver) BeginDrawing() {}
func (d *Driver) EndDrawing()   {}

func (d *Driver) ClearBackground(c hal.RGBA) {}
func (d *Driver) DrawRectangle(x, y, w, h int32, c hal.RGBA) {}
func (d *Driver) DrawRectanglePro(rec hal.Rect, origin hal.V2, rotation float32, c hal.RGBA) {
}
func (d *Driver) DrawRectangleLines(x, y, w, h int32, thick float32, c hal.RGBA) {}
func (d *Driver) DrawCircle(x, y int32, radius float32, c hal.RGBA)               {}
func (d *Driver) DrawCircleLines(x, y int32, radius float32, c hal.RGBA)          {}
func (d *Driver) DrawTriangle(v1, v2, v3 hal.V2, c hal.RGBA)                      {}
func (d *Driver) DrawPoly(center hal.V2, sides int32, radius float32, rotation float32, c hal.RGBA) {
}
func (d *Driver) DrawText(text string, x, y, size int32, c hal.RGBA) {}

func (d *Driver) BeginMode3D(cam hal.Camera3D) {}
func (d *Driver) EndMode3D()                   {}
func (d *Driver) DrawCube(position hal.V3, width, height, length float32, c hal.RGBA) {
}
func (d *Driver) DrawGrid(slices int32, spacing float32) {}

func (d *Driver) IsKeyDown(key int32) bool                  { return false }
func (d *Driver) IsKeyPressed(key int32) bool               { return false }
func (d *Driver) IsMouseButtonPressed(button int32) bool    { return false }
func (d *Driver) GetMousePosition() hal.V2                  { return hal.V2{} }
func (d *Driver) GetMouseWheelMove() float32                { return 0 }
func (d *Driver) IsGamepadAvailable(id int32) bool          { return false }
func (d *Driver) GetGamepadAxisMovement(id, axis int32) float32 { return 0 }
