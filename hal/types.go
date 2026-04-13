package hal

// RGBA represents a 32-bit color.
type RGBA struct {
	R, G, B, A uint8
}

// V2 represents a 2D vector.
type V2 struct {
	X, Y float32
}

// V3 represents a 3D vector.
type V3 struct {
	X, Y, Z float32
}

// V4 represents a 4D vector.
type V4 struct {
	X, Y, Z, W float32
}

// Rect represents a 2D rectangle.
type Rect struct {
	X, Y, Width, Height float32
}

// Matrix stores a 4×4 transform in the same **column-major** element order as Raylib’s rl.Matrix:
// column 0 is (M0,M1,M2,M3), column 1 is (M4,M5,M6,M7), etc. Field names match rl.M0..M15.
type Matrix struct {
	M0, M4, M8, M12  float32
	M1, M5, M9, M13  float32
	M2, M6, M10, M14 float32
	M3, M7, M11, M15 float32
}

// Common colors
var (
	Black = RGBA{0, 0, 0, 255}
	White = RGBA{255, 255, 255, 255}
	Red   = RGBA{255, 0, 0, 255}
	Green = RGBA{0, 255, 0, 255}
	Blue  = RGBA{0, 0, 255, 255}
)
