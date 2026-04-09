// Command puregohello opens a window using only internal/raylibpurego (no CGO).
// Build: CGO_ENABLED=0 go build -o puregohello ./cmd/puregohello
// Runtime: place raylib.dll (Windows), libraylib.so (Linux), or libraylib.dylib (macOS) next to the binary or on PATH.
package main

import (
	"fmt"
	"log"
	"os"

	"moonbasic/internal/raylibpurego"
)

// Raylib keyboard constants (subset of raylib.Key*; avoid importing raylib-go which loads the DLL at init).
const (
	keyRight = 262
	keyLeft  = 263
	keyDown  = 264
	keyUp    = 265
)

func main() {
	lib, err := raylibpurego.Load("")
	if err != nil {
		log.Fatal(err)
	}
	var g raylibpurego.Game
	if err := raylibpurego.RegisterGame(lib, &g); err != nil {
		log.Fatal(err)
	}
	var img raylibpurego.Image
	white := raylibpurego.Color{R: 240, G: 240, B: 255, A: 255}
	g.GenImageColor(raylibpurego.ImagePtr(&img), 64, 64, raylibpurego.ColorPtr(white))

	var tex raylibpurego.Texture2D
	g.LoadTextureFromImage(raylibpurego.TexturePtr(&tex), raylibpurego.ImagePtr(&img))
	g.UnloadImage(raylibpurego.ImagePtr(&img))
	defer g.UnloadTexture(raylibpurego.TexturePtr(&tex))

	g.InitWindow(800, 450, "moonbasic purego hello (CGO_ENABLED=0)")
	defer g.CloseWindow()

	var x, y float32 = 200, 175
	fmt.Fprintf(os.Stderr, "puregohello: arrow keys move the square; close window to exit.\n")

	for !g.WindowShouldClose() {
		dt := g.GetFrameTime()
		speed := float32(280)
		if g.IsKeyDown(keyRight) {
			x += speed * dt
		}
		if g.IsKeyDown(keyLeft) {
			x -= speed * dt
		}
		if g.IsKeyDown(keyDown) {
			y += speed * dt
		}
		if g.IsKeyDown(keyUp) {
			y -= speed * dt
		}

		g.BeginDrawing()
		gray := raylibpurego.Color{R: 30, G: 30, B: 40, A: 255}
		g.ClearBackground(raylibpurego.ColorPtr(gray))
		tint := raylibpurego.Color{R: 255, G: 255, B: 255, A: 255}
		g.DrawTexture(raylibpurego.TexturePtr(&tex), int32(x), int32(y), raylibpurego.ColorPtr(tint))
		hint := raylibpurego.Color{R: 200, G: 200, B: 200, A: 255}
		g.DrawText("purego + sidecar raylib — no C compiler", 12, 12, 20, raylibpurego.ColorPtr(hint))
		g.EndDrawing()
	}
}
