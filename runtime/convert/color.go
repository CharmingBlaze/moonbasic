// Package convert holds small helpers shared by draw/camera modules.
package convert

// Color4 is an 8-bit RGBA tuple (0–255 per channel).
type Color4 struct {
	R, G, B, A uint8
}

// NewColor4 builds Color4 from int32 components with clamping.
func NewColor4(r, g, b, a int32) Color4 {
	return Color4{
		R: uint8(clampU8(r)),
		G: uint8(clampU8(g)),
		B: uint8(clampU8(b)),
		A: uint8(clampU8(a)),
	}
}

func clampU8(v int32) int32 {
	switch {
	case v < 0:
		return 0
	case v > 255:
		return 255
	default:
		return v
	}
}
