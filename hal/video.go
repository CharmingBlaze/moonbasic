package hal

// VideoDevice abstracts rendering operations.
type VideoDevice interface {
	BeginDrawing()
	EndDrawing()
	ClearBackground(color RGBA)
	
	DrawRectangle(x, y, w, h int32, color RGBA)
	DrawRectangleLines(x, y, w, h int32, thick float32, color RGBA)
	DrawRectanglePro(rec Rect, origin V2, rotation float32, color RGBA)
	DrawCircle(x, y int32, radius float32, color RGBA)
	DrawCircleLines(x, y int32, radius float32, color RGBA)
	DrawTriangle(v1, v2, v3 V2, color RGBA)
	DrawPoly(center V2, sides int32, radius float32, rotation float32, color RGBA)
	DrawText(text string, x, y, size int32, color RGBA)
	
	// 3D Rendering
	BeginMode3D(camera Camera3D)
	EndMode3D()
	DrawCube(position V3, width, height, length float32, color RGBA)
	DrawGrid(slices int32, spacing float32)
}

// Camera3D mirrors Raylib's Camera3D for projection.
type Camera3D struct {
	Position V3
	Target   V3
	Up       V3
	Fovy     float32
	Projection int32
}
