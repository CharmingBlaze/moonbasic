//go:build cgo || (windows && !cgo)

package mbdebug

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
)

// DrawFrameOverlay draws DEBUG.WATCH lines on top of the current frame.
// Called from window.RENDER.FRAME before EndDrawing; requires an active Raylib frame.
func (m *Module) DrawFrameOverlay() {
	reg := runtime.ActiveRegistry()
	dt := rl.GetFrameTime()
	m.mu.Lock()
	user := m.overlayUser
	showGraph := m.showFPSGraph
	monitor := m.monitorOn
	inspectID := m.inspectID
	logs := make([]logLine, len(m.logLines))
	copy(logs, m.logLines)
	m.fpsHistory[m.fpsIdx] = dt
	m.fpsIdx = (m.fpsIdx + 1) % len(m.fpsHistory)

	lines := make([]watchEntry, len(m.watches))
	copy(lines, m.watches)
	m.mu.Unlock()

	if showGraph {
		m.drawFPSGraph()
	}
	
	if monitor {
		m.drawSystemMonitor()
	}
	
	if len(logs) > 0 {
		m.drawConsoleLogs(logs)
	}

	if inspectID > 0 {
		m.drawEntityInspector(inspectID)
	}

	if len(lines) == 0 {
		return
	}
	if reg != nil && !reg.DebugMode && !user {
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

func (m *Module) drawFPSGraph() {
	sw := int32(rl.GetScreenWidth())
	const width int32 = 120
	const height int32 = 60
	const pad int32 = 10
	x := sw - width - pad
	y := pad

	rl.DrawRectangle(x, y, width, height, rl.Color{R: 0, G: 0, B: 0, A: 160})
	rl.DrawRectangleLines(x, y, width, height, rl.Color{R: 200, G: 200, B: 200, A: 100})

	m.mu.Lock()
	defer m.mu.Unlock()

	// Draw 60fps / 16.6ms target line
	targetY := y + height - int32(0.5*float32(height)) // 16.6ms is half of 33.3ms scale
	if targetY >= y && targetY < y+height {
		rl.DrawLine(x, targetY, x+width, targetY, rl.Color{R: 0, G: 255, B: 0, A: 100})
	}

	for i := int32(0); i < width; i++ {
		idx := (m.fpsIdx + int(i)) % len(m.fpsHistory)
		val := m.fpsHistory[idx]
		// Map 0-33ms to 0-height
		h := int32(val * float32(height) / 0.0333)
		if h > height {
			h = height
		}
		color := rl.Color{R: 0, G: 255, B: 0, A: 200}
		if val > 0.017 { // > 60fps
			color = rl.Color{R: 255, G: 255, B: 0, A: 200}
		}
		if val > 0.033 { // > 30fps
			color = rl.Color{R: 255, G: 0, B: 0, A: 200}
		}
		rl.DrawLine(x+i, y+height, x+i, y+height-h, color)
	}

	fps := int32(0)
	if rl.GetFrameTime() > 0 {
		fps = int32(1.0 / rl.GetFrameTime())
	}
	rl.DrawText(fmt.Sprintf("%d FPS", fps), x+5, y+5, 12, rl.White)
}

func (m *Module) drawSystemMonitor() {
	sw := int32(rl.GetScreenWidth())
	const width int32 = 180
	const height int32 = 80
	const pad int32 = 10
	x := sw - width - pad
	y := int32(rl.GetScreenHeight()) - height - pad
	
	rl.DrawRectangle(x, y, width, height, rl.Color{R: 20, G: 20, B: 30, A: 200})
	rl.DrawRectangleLines(x, y, width, height, rl.Color{R: 100, G: 100, B: 255, A: 200})
	
	fps := int32(0)
	if rl.GetFrameTime() > 0 { fps = int32(1.0 / rl.GetFrameTime()) }

	rl.DrawText("SYSTEM MONITOR", x+10, y+10, 14, rl.SkyBlue)
	rl.DrawText(fmt.Sprintf("FPS: %d", fps), x+10, y+30, 16, rl.White)
	
	// Stub RAM for now (would need runtime.ReadMemStats)
	rl.DrawText("RAM: 124MB / 1GB", x+10, y+50, 12, rl.Gray)
}

func (m *Module) drawConsoleLogs(logs []logLine) {
	sh := int32(rl.GetScreenHeight())
	y := sh - 30
	for i := len(logs)-1; i >= 0; i-- {
		l := logs[i]
		rl.DrawText(l.msg, 10, y, 16, l.color)
		y -= 20
	}
}

func (m *Module) drawEntityInspector(id int64) {
	reg := runtime.ActiveRegistry()
	if reg == nil || reg.ResolveEntityWorldPos == nil { return }
	
	pos, ok := reg.ResolveEntityWorldPos(id)
	if !ok { return }
	
	const width int32 = 200
	const height int32 = 100
	x := int32(10)
	y := int32(100)
	
	rl.DrawRectangle(x, y, width, height, rl.Color{R: 40, G: 0, B: 0, A: 180})
	rl.DrawRectangleLines(x, y, width, height, rl.Red)
	
	rl.DrawText(fmt.Sprintf("ENTITY INSPECTOR #%d", id), x+10, y+10, 14, rl.Yellow)
	rl.DrawText(fmt.Sprintf("WorldPos: %.2f, %.2f, %.2f", pos.X, pos.Y, pos.Z), x+10, y+35, 14, rl.White)
	
	// If we had ResolveEntityName, we'd use it here.
	rl.DrawText("State: ACTIVE", x+10, y+60, 14, rl.Green)
}
