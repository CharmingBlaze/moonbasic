//go:build darwin || freebsd || linux || windows

package raylibpurego

import (
	"fmt"

	"github.com/ebitengine/purego"
)

// RegisterGame binds a core set of symbols for window, input, drawing, and texture loading.
func RegisterGame(lib *LoadResult, g *Game) error {
	if lib == nil || g == nil {
		return fmt.Errorf("raylibpurego: nil lib or game")
	}
	h := lib.Handle
	purego.RegisterLibFunc(&g.SetConfigFlags, h, "SetConfigFlags")
	purego.RegisterLibFunc(&g.InitWindow, h, "InitWindow")
	purego.RegisterLibFunc(&g.CloseWindow, h, "CloseWindow")
	purego.RegisterLibFunc(&g.IsWindowReady, h, "IsWindowReady")
	purego.RegisterLibFunc(&g.WindowShouldClose, h, "WindowShouldClose")
	purego.RegisterLibFunc(&g.SetTargetFPS, h, "SetTargetFPS")
	purego.RegisterLibFunc(&g.GetFPS, h, "GetFPS")
	purego.RegisterLibFunc(&g.BeginDrawing, h, "BeginDrawing")
	purego.RegisterLibFunc(&g.EndDrawing, h, "EndDrawing")
	purego.RegisterLibFunc(&g.ClearBackground, h, "ClearBackground")
	purego.RegisterLibFunc(&g.GetFrameTime, h, "GetFrameTime")
	purego.RegisterLibFunc(&g.GetMouseX, h, "GetMouseX")
	purego.RegisterLibFunc(&g.GetMouseY, h, "GetMouseY")
	purego.RegisterLibFunc(&g.IsKeyDown, h, "IsKeyDown")
	purego.RegisterLibFunc(&g.DrawTexture, h, "DrawTexture")
	purego.RegisterLibFunc(&g.DrawRectangle, h, "DrawRectangle")
	purego.RegisterLibFunc(&g.DrawText, h, "DrawText")
	purego.RegisterLibFunc(&g.LoadTexture, h, "LoadTexture")
	purego.RegisterLibFunc(&g.UnloadTexture, h, "UnloadTexture")
	purego.RegisterLibFunc(&g.GenImageColor, h, "GenImageColor")
	purego.RegisterLibFunc(&g.LoadTextureFromImage, h, "LoadTextureFromImage")
	purego.RegisterLibFunc(&g.UnloadImage, h, "UnloadImage")
	return nil
}
