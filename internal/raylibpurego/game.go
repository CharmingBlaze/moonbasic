package raylibpurego

import "unsafe"

// Game holds dynamically bound raylib symbols for a minimal CGO-free loop.
type Game struct {
	SetConfigFlags       func(flags uint32)
	InitWindow           func(width int32, height int32, title string)
	CloseWindow          func()
	IsWindowReady        func() bool
	WindowShouldClose    func() bool
	SetTargetFPS         func(fps int32)
	GetFPS               func() int32
	BeginDrawing           func()
	EndDrawing             func()
	ClearBackground        func(col uintptr)
	GetFrameTime           func() float32
	GetMouseX              func() int32
	GetMouseY              func() int32
	IsKeyDown              func(key int32) bool
	DrawTexture            func(tex uintptr, posX int32, posY int32, tint uintptr)
	DrawRectangle          func(posX int32, posY int32, width int32, height int32, col uintptr)
	DrawText               func(text string, posX int32, posY int32, fontSize int32, col uintptr)
	LoadTexture            func(outTex uintptr, fileName string)
	UnloadTexture          func(tex uintptr)
	GenImageColor          func(outImg uintptr, width int32, height int32, col uintptr)
	LoadTextureFromImage   func(outTex uintptr, image uintptr)
	UnloadImage            func(image uintptr)
}

// ColorPtr returns a uintptr to c for functions that take Color by value as uintptr in raylib-go purego.
func ColorPtr(c Color) uintptr {
	return uintptr(unsafe.Pointer(&c))
}

// TexturePtr returns uintptr to texture struct for Draw/Unload.
func TexturePtr(t *Texture2D) uintptr {
	return uintptr(unsafe.Pointer(t))
}

// ImagePtr returns uintptr to Image.
func ImagePtr(i *Image) uintptr {
	return uintptr(unsafe.Pointer(i))
}
