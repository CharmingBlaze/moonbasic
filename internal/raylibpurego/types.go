package raylibpurego

// Color matches raylib Color (RGBA bytes).
type Color struct {
	R, G, B, A uint8
}

// Texture2D matches raylib Texture2D layout for LoadTexture / DrawTexture.
type Texture2D struct {
	ID      uint32
	Width   int32
	Height  int32
	Mipmaps int32
	Format  int32 // PixelFormat
}

// Image matches raylib Image for GenImageColor.
type Image struct {
	Format int32
	Mipmaps int32
	Width   int32
	Height  int32
	Data    uintptr
}
