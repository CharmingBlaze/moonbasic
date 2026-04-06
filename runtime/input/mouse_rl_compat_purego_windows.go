//go:build !cgo && windows

package input

import rl "github.com/gen2brain/raylib-go/raylib"

func setMousePositionCompat(x, y int) {
	rl.SetMousePosition(int32(x), int32(y))
}

func setMouseOffsetCompat(x, y int) {
	rl.SetMouseOffset(int32(x), int32(y))
}
