package blitzengine

import (
	"sync"
)

// pen2D holds Blitz-style 2D pen state for SETCOLOR / SETALPHA / SETORIGIN / SETVIEWPORT.
type pen2D struct {
	mu sync.Mutex

	r, g, b   int
	a         float64 // 0..1, stored as float; drawn as uint8 a*255
	ox, oy    float64
	vx, vy, vw, vh int32
	hasViewport bool
}

func (p *pen2D) setColor(r, g, b int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.r, p.g, p.b = r, g, b
}

func (p *pen2D) setAlpha(a float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.a = a
}

func (p *pen2D) setOrigin(x, y float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ox, p.oy = x, y
}

func (p *pen2D) setViewport(x, y, w, h int32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.vx, p.vy, p.vw, p.vh = x, y, w, h
	p.hasViewport = true
}

func (p *pen2D) clearViewport() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.hasViewport = false
}

func (p *pen2D) rgbaA() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	a := int(p.a * 255)
	if a < 0 {
		return 0
	}
	if a > 255 {
		return 255
	}
	return a
}

func (p *pen2D) rgb() (int, int, int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.r, p.g, p.b
}

func (p *pen2D) offset() (float64, float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.ox, p.oy
}

func (p *pen2D) viewport() (int32, int32, int32, int32, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.vx, p.vy, p.vw, p.vh, p.hasViewport
}
