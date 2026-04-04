//go:build cgo

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
	frameW    int32
	frameH    int32
	numFrames int
	curFrame  int
	playing   bool
	fps       float32
	accum     float32

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

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	reg.Register("SPRITE.LOAD", "sprite", m.spLoad)
	reg.Register("SPRITE.DRAW", "sprite", m.spDraw)
	reg.Register("SPRITE.SETPOS", "sprite", m.spSetPos)
	reg.Register("SPRITE.SETPOSITION", "sprite", m.spSetPos)
	reg.Register("SPRITE.DEFANIM", "sprite", m.spDefAnim)
	reg.Register("SPRITE.PLAYANIM", "sprite", m.spPlayAnim)
	reg.Register("SPRITE.UPDATEANIM", "sprite", m.spUpdateAnim)
	reg.Register("SPRITE.HIT", "sprite", m.spHit)
	m.registerAtlas(reg)
	m.registerAnim(reg)
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
	if s.numFrames < 1 && s.anim == nil {
		return value.Nil, nil
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
	pos := rl.Vector2{X: float32(x) + s.x, Y: float32(y) + s.y}
	tint := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	rl.DrawTextureRec(s.tex, rec, pos, tint)
	return value.Nil, nil
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
