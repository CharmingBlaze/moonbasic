//go:build cgo

package mbdebug

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// DrawFrameOverlay draws DEBUG.WATCH lines on top of the current frame.
// Called from window.RENDER.FRAME before EndDrawing; requires an active Raylib frame.
func (m *Module) DrawFrameOverlay() {
	m.mu.Lock()
	lines := make([]watchEntry, len(m.watches))
	copy(lines, m.watches)
	m.mu.Unlock()
	if len(lines) == 0 {
		return
	}
	const fontSize int32 = 18
	y := int32(8)
	c := rl.Color{R: 255, G: 255, B: 255, A: 220}
	shadow := rl.Color{R: 0, G: 0, B: 0, A: 200}
	for _, e := range lines {
		s := fmt.Sprintf("%s: %s", e.label, e.text)
		rl.DrawText(s, 9, y+1, fontSize, shadow)
		rl.DrawText(s, 8, y, fontSize, c)
		y += fontSize + 6
	}
}
