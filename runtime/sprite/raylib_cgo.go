//go:build cgo || (windows && !cgo)

package mbsprite

import (
	"fmt"
	"image/color"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type spriteObj struct {
	tex       rl.Texture2D
	fromAtlas bool
	srcX      int32
	srcY      int32
	// atlasRegionW/H: full atlas rect (used to infer per-frame width for ANIM strips)
	atlasRegionW int32
	atlasRegionH int32

	x, y      float32
	originX   float32
	originY   float32
	frameW    int32
	frameH    int32
	numFrames int
	curFrame  int
	playing   bool
	fps       float32
	accum     float32

	rangeStart   int
	rangeEnd     int
	rangeSpeed   float32
	rangeLoop    bool
	rangePlaying bool
	rangeAccum   float32

	anim *animMachine

	release heap.ReleaseOnce
}

func (s *spriteObj) TypeName() string { return "Sprite" }

func (s *spriteObj) TypeTag() uint16 { return heap.TagSprite }

func (s *spriteObj) Free() {
	s.release.Do(func() {
		if !s.fromAtlas {
			rl.UnloadTexture(s.tex)
		}
	})
}

func argHandle(v value.Value) (heap.Handle, bool) {
	if v.Kind != value.KindHandle {
		return 0, false
	}
	return heap.Handle(v.IVal), true
}

func argInt32(v value.Value) (int32, bool) {
	if i, ok := v.ToInt(); ok {
		return int32(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return int32(f), true
	}
	return 0, false
}

func argF(v value.Value) (float32, bool) {
	if f, ok := v.ToFloat(); ok {
		return float32(f), true
	}
	if i, ok := v.ToInt(); ok {
		return float32(i), true
	}
	return 0, false
}

func argSpriteInt(v value.Value) (int, bool) {
	if i, ok := v.ToInt(); ok {
		return int(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return int(f), true
	}
	return 0, false
}

func argBool(v value.Value) (bool, bool) {
	if v.Kind == value.KindBool {
		return v.IVal != 0, true
	}
	if i, ok := v.ToInt(); ok {
		return i != 0, true
	}
	if f, ok := v.ToFloat(); ok {
		return f != 0, true
	}
	return false, false
}

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	reg.Register("SPRITE.LOAD", "sprite", m.spLoad)
	reg.Register("SPRITE.DRAW", "sprite", m.spDraw)
	reg.Register("SPRITE.SETPOS", "sprite", m.spSetPos)
	reg.Register("SPRITE.SETPOSITION", "sprite", m.spSetPos)
	reg.Register("SPRITE.DEFANIM", "sprite", m.spDefAnim)
	reg.Register("SPRITE.PLAYANIM", "sprite", m.spPlayAnim)
	reg.Register("SPRITE.UPDATEANIM", "sprite", m.spUpdateAnim)
	reg.Register("SPRITE.SETFRAME", "sprite", m.spSetFrame)
	reg.Register("SPRITE.PLAY", "sprite", m.spPlayRange)
	reg.Register("SPRITE.SETORIGIN", "sprite", m.spSetOrigin)
	reg.Register("SPRITE.HIT", "sprite", m.spHit)
	reg.Register("SPRITECOLLIDE", "sprite", m.spHit)
	reg.Register("SPRITE.POINTHIT", "sprite", m.spPointHit)
	reg.Register("SPRITE.FREE", "sprite", m.spFree)
	m.registerAtlas(reg)
	m.registerAnim(reg)
	m.registerSpriteExtras(reg)
}

func (m *Module) spFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.h == nil {
		return value.Nil, runtime.Errorf("SPRITE.FREE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("SPRITE.FREE expects (sprite)")
	}
	return value.Nil, m.h.Free(heap.Handle(args[0].IVal))
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

func (m *Module) spLoad(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("SPRITE.LOAD: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("SPRITE.LOAD expects 1 string path")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	t := rl.LoadTexture(path)
	s := &spriteObj{
		tex:          t,
		fromAtlas:    false,
		atlasRegionW: t.Width,
		atlasRegionH: t.Height,
		frameW:       t.Width,
		frameH:       t.Height,
		numFrames:    1,
		fps:          8,
	}
	id, err := m.h.Alloc(s)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) spDraw(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("SPRITE.DRAW expects 3 arguments (handle, x, y)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.DRAW: invalid handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argInt32(args[1])
	y, ok2 := argInt32(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("SPRITE.DRAW: x,y must be numeric")
	}
	if err := m.drawSpriteAtScreen(s, x, y); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

// drawSpriteAtScreen draws the sprite’s current frame at integer screen (x,y), plus SetPos offsets.
func (m *Module) drawSpriteAtScreen(s *spriteObj, screenX, screenY int32) error {
	if s.numFrames < 1 && s.anim == nil {
		return nil
	}
	if s.anim != nil {
		s.syncAnimFrame()
	}
	srcX := float32(s.srcX) + float32(s.curFrame)*float32(s.frameW)
	rec := rl.Rectangle{
		X:      srcX,
		Y:      float32(s.srcY),
		Width:  float32(s.frameW),
		Height: float32(s.frameH),
	}
	pos := rl.Vector2{X: float32(screenX) + s.x - s.originX, Y: float32(screenY) + s.y - s.originY}
	tint := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	rl.DrawTextureRec(s.tex, rec, pos, tint)
	return nil
}

func (m *Module) spSetPos(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("SPRITE.SETPOS expects 3 arguments (handle, x, y)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.SETPOS: invalid handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argF(args[1])
	y, ok2 := argF(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("SPRITE.SETPOS: non-numeric position")
	}
	s.x = x
	s.y = y
	return value.Nil, nil
}

func (m *Module) spDefAnim(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("SPRITE.DEFANIM expects 2 arguments (handle, frameCountString)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.DEFANIM: invalid handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	countStr, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	n, err := strconv.Atoi(countStr)
	if err != nil || n < 1 {
		return value.Nil, fmt.Errorf("SPRITE.DEFANIM: frame count must be a positive integer string")
	}
	s.numFrames = n
	avail := s.tex.Width - s.srcX
	if avail < int32(n) {
		return value.Nil, fmt.Errorf("SPRITE.DEFANIM: not enough width for frames")
	}
	s.frameW = avail / int32(n)
	s.frameH = s.tex.Height - s.srcY
	s.curFrame = 0
	s.accum = 0
	return value.Nil, nil
}

func (m *Module) spPlayAnim(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("SPRITE.PLAYANIM expects 2 arguments (handle, name)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.PLAYANIM: invalid handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	_, err = rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	s.playing = true
	s.curFrame = 0
	s.accum = 0
	return value.Nil, nil
}

func (m *Module) spSetFrame(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SPRITE.SETFRAME expects 2 arguments (handle, frameIndex)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.SETFRAME: invalid handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	idx, ok := argSpriteInt(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.SETFRAME: frame index must be numeric")
	}
	if s.numFrames > 0 {
		if idx < 0 {
			idx = 0
		}
		if idx >= s.numFrames {
			idx = s.numFrames - 1
		}
	}
	s.curFrame = idx
	s.rangePlaying = false
	return value.Nil, nil
}

func (m *Module) spPlayRange(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("SPRITE.PLAY expects 5 arguments (handle, start#, end#, speed#, loop)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.PLAY: invalid handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	start, ok1 := argSpriteInt(args[1])
	end, ok2 := argSpriteInt(args[2])
	speed, ok3 := argF(args[3])
	loop, ok4 := argBool(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("SPRITE.PLAY: start, end, speed, and loop must be valid")
	}
	if start > end {
		start, end = end, start
	}
	if s.numFrames > 0 {
		if start < 0 {
			start = 0
		}
		if end >= s.numFrames {
			end = s.numFrames - 1
		}
		if start > end {
			start = end
		}
	}
	s.rangeStart = start
	s.rangeEnd = end
	s.rangeSpeed = speed
	s.rangeLoop = loop
	s.rangeAccum = 0
	s.rangePlaying = true
	s.playing = false
	s.curFrame = start
	return value.Nil, nil
}

func (m *Module) spSetOrigin(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("SPRITE.SETORIGIN expects 3 arguments (handle, originX#, originY#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.SETORIGIN: invalid handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	ox, ok1 := argF(args[1])
	oy, ok2 := argF(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("SPRITE.SETORIGIN: origin must be numeric")
	}
	s.originX = ox
	s.originY = oy
	return value.Nil, nil
}

func (m *Module) spUpdateAnim(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SPRITE.UPDATEANIM expects 2 arguments (handle, deltaSeconds)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.UPDATEANIM: invalid handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	dt, ok := argF(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.UPDATEANIM: delta must be numeric")
	}
	if s.anim != nil {
		return value.Nil, nil
	}
	if s.rangePlaying {
		if s.rangeSpeed <= 0 {
			return value.Nil, nil
		}
		s.rangeAccum += dt * s.rangeSpeed
		for s.rangeAccum >= 1 {
			s.rangeAccum--
			if s.curFrame < s.rangeEnd {
				s.curFrame++
			} else {
				if s.rangeLoop {
					s.curFrame = s.rangeStart
				} else {
					s.rangePlaying = false
				}
			}
		}
		return value.Nil, nil
	}
	if !s.playing || s.numFrames < 1 || s.fps <= 0 {
		return value.Nil, nil
	}
	s.accum += dt * s.fps
	for s.accum >= 1 {
		s.accum--
		s.curFrame = (s.curFrame + 1) % s.numFrames
	}
	return value.Nil, nil
}

func (m *Module) spHit(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SPRITE.HIT expects 2 arguments (handleA, handleB)")
	}
	ha, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.HIT: invalid handle A")
	}
	hb, ok := argHandle(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.HIT: invalid handle B")
	}
	a, err := heap.Cast[*spriteObj](m.h, ha)
	if err != nil {
		return value.Nil, err
	}
	b, err := heap.Cast[*spriteObj](m.h, hb)
	if err != nil {
		return value.Nil, err
	}
	ax := float64(a.x)
	ay := float64(a.y)
	aw := float64(a.frameW)
	ah := float64(a.frameH)
	bx := float64(b.x)
	by := float64(b.y)
	bw := float64(b.frameW)
	bh := float64(b.frameH)
	hit := ax < bx+bw && ax+aw > bx && ay < by+bh && ay+ah > by
	return value.FromBool(hit), nil
}

func (m *Module) spPointHit(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("SPRITE.POINTHIT expects 3 arguments (handle, x, y)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.POINTHIT: invalid sprite handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	px, ok1 := argF(args[1])
	py, ok2 := argF(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("SPRITE.POINTHIT: x and y must be numeric")
	}
	sx := float64(s.x)
	sy := float64(s.y)
	sw := float64(s.frameW)
	sh := float64(s.frameH)
	pxf := float64(px)
	pyf := float64(py)
	inside := pxf >= sx && pxf < sx+sw && pyf >= sy && pyf < sy+sh
	return value.FromBool(inside), nil
}
