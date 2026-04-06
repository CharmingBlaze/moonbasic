//go:build cgo

package input

import rl "github.com/gen2brain/raylib-go/raylib"

func setMousePositionCompat(x, y int) {
	rl.SetMousePosition(x, y)
}

func setMouseOffsetCompat(x, y int) {
	rl.SetMouseOffset(x, y)
}
