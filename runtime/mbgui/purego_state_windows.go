//go:build !cgo && windows

package mbgui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Rounded pixel rect for stable map keys (avoids fmt per call and unbounded string keys when
// coordinates drift). Collides only if widgets share the same integer pixel rect.
type puregoRectKey struct {
	x, y, w, h int32
}

func makePuregoRectKey(b rl.Rectangle) puregoRectKey {
	return puregoRectKey{
		x: int32(b.X + 0.5),
		y: int32(b.Y + 0.5),
		w: int32(b.Width + 0.5),
		h: int32(b.Height + 0.5),
	}
}

const puregoWidgetMapMax = 2048

// Minimal immediate-mode GUI state for Windows + CGO_ENABLED=0 (Raylib purego).
// Not API-identical to raygui; see docs/reference/GUI.md for full CGO behavior.
var pg struct {
	enabled  bool
	disabled bool
	locked   bool

	alpha    float32
	guiState int32

	textSize    float32
	textSpacing float32
	lineExtra   float32
	wrapMode    int32
	alignH      int32
	alignV      int32
	tooltipOn   bool
	tooltipText string
	iconScale   int32
	styleInt    map[int64]int64
	styleColor  map[int64]rl.Color
	lastTheme   string

	// Per-rectangle widget state for simplified controls (toggle group without explicit index).
	widgetInt map[puregoRectKey]int32
}

// puregoNoteWidgetWrite caps widgetInt so pathological UIs cannot grow memory without bound.
func puregoNoteWidgetWrite() {
	if len(pg.widgetInt) <= puregoWidgetMapMax {
		return
	}
	pg.widgetInt = make(map[puregoRectKey]int32, puregoWidgetMapMax/2)
}

func init() {
	pg.enabled = true
	pg.alpha = 1
	pg.textSize = 12
	pg.textSpacing = 1
	pg.widgetInt = make(map[puregoRectKey]int32)
	puregoResetDefaultTheme()
}

func styleKey(control, property int32) int64 {
	return (int64(control) << 32) | (int64(property) & 0xffffffff)
}

func puregoCanInteract() bool {
	return pg.enabled && !pg.disabled && !pg.locked
}

func puregoMulAlpha(c rl.Color, a float32) rl.Color {
	if a < 0 {
		a = 0
	}
	if a > 1 {
		a = 1
	}
	return rl.Color{R: c.R, G: c.G, B: c.B, A: uint8(float32(c.A) * a)}
}

func puregoBaseTextColor() rl.Color {
	c := pg.styleColor[styleKey(0, 2)] // GPROP_TEXT_COLOR_NORMAL on DEFAULT
	if c.A == 0 && c.R == 0 && c.G == 0 && c.B == 0 {
		return rl.Color{R: 230, G: 230, B: 235, A: 255}
	}
	return c
}

func puregoPanelColor() rl.Color {
	c := pg.styleColor[styleKey(0, 1)]
	if c.A == 0 && c.R == 0 && c.G == 0 && c.B == 0 {
		return rl.Color{R: 48, G: 54, B: 64, A: 255}
	}
	return c
}

func puregoDrawLabelText(text string, b rl.Rectangle, col rl.Color) {
	font := rl.GetFontDefault()
	size := pg.textSize
	if size < 8 {
		size = 8
	}
	sp := pg.textSpacing
	col = puregoMulAlpha(col, pg.alpha)
	m := rl.MeasureTextEx(font, text, size, sp)
	x := b.X + 4
	y := b.Y + (b.Height-size)/2
	switch pg.alignH {
	case 1:
		x = b.X + (b.Width-m.X)/2
	case 2:
		x = b.X + b.Width - m.X - 4
	}
	rl.DrawTextEx(font, text, rl.Vector2{X: x, Y: y}, size, sp, col)
}

func puregoDrawButtonChrome(b rl.Rectangle, pressed bool) {
	bg := puregoPanelColor()
	border := rl.Color{R: 90, G: 100, B: 115, A: uint8(float32(255) * pg.alpha)}
	if pressed {
		bg = rl.Color{R: 70, G: 82, B: 100, A: uint8(float32(255) * pg.alpha)}
	}
	if pg.disabled {
		bg = rl.Color{R: 40, G: 44, B: 50, A: uint8(float32(255) * pg.alpha)}
		border = rl.Color{R: 70, G: 74, B: 80, A: uint8(float32(255) * pg.alpha)}
	}
	rl.DrawRectangleRec(b, bg)
	rl.DrawRectangleLinesEx(b, 1, border)
}

func puregoClickIn(b rl.Rectangle) bool {
	if !puregoCanInteract() {
		return false
	}
	mp := rl.GetMousePosition()
	if !rl.CheckCollisionPointRec(mp, b) {
		return false
	}
	return rl.IsMouseButtonPressed(rl.MouseLeftButton)
}

func puregoDragValue(b rl.Rectangle, val, minV, maxV float32) float32 {
	if !puregoCanInteract() {
		return val
	}
	if !rl.IsMouseButtonDown(rl.MouseLeftButton) {
		return val
	}
	mp := rl.GetMousePosition()
	if !rl.CheckCollisionPointRec(mp, b) {
		return val
	}
	t := (mp.X - b.X) / b.Width
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return minV + t*(maxV-minV)
}
