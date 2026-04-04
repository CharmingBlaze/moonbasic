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
	const pad int32 = 6
	y := int32(8)
	c := rl.Color{R: 255, G: 255, B: 255, A: 235}
	shadow := rl.Color{R: 0, G: 0, B: 0, A: 200}
	maxW := int32(0)
	totalH := pad
	for _, e := range lines {
		s := fmt.Sprintf("%s: %s", e.label, e.text)
		w := rl.MeasureText(s, fontSize)
		if int32(w) > maxW {
			maxW = int32(w)
		}
		totalH += fontSize + 6
	}
	totalH += pad
	bg := rl.Color{R: 12, G: 14, B: 20, A: 200}
	rl.DrawRectangle(0, 0, maxW+pad*2+8, totalH, bg)
	for _, e := range lines {
		s := fmt.Sprintf("%s: %s", e.label, e.text)
		rl.DrawText(s, 9, y+1, fontSize, shadow)
		rl.DrawText(s, 8, y, fontSize, c)
		y += fontSize + 6
	}
}
