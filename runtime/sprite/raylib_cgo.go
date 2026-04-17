//go:build cgo || (windows && !cgo)

package mbsprite

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

const spriteDefaultAlpha = float32(1)

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

	scaleX, scaleY float32
	rotRad         float32 // radians, CCW (DrawTexturePro rotation in degrees)
	tr, tg, tb     uint8
	alpha          float32 // 0–1

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
	reg.Register("SPRITE.GETPOS", "sprite", m.spGetPos)
	reg.Register("SPRITE.DEFANIM", "sprite", m.spDefAnim)
	reg.Register("SPRITE.PLAYANIM", "sprite", m.spPlayAnim)
	reg.Register("SPRITE.UPDATEANIM", "sprite", m.spUpdateAnim)
	reg.Register("SPRITE.SETFRAME", "sprite", m.spSetFrame)
	reg.Register("SPRITE.PLAY", "sprite", m.spPlayRange)
	reg.Register("SPRITE.SETORIGIN", "sprite", m.spSetOrigin)
	reg.Register("SPRITE.GETSCALE", "sprite", m.spGetScale)
	reg.Register("SPRITE.SETSCALE", "sprite", m.spSetScale)
	reg.Register("SPRITE.GETROT", "sprite", m.spGetRot)
	reg.Register("SPRITE.SETROT", "sprite", m.spSetRot)
	reg.Register("SPRITE.GETCOLOR", "sprite", m.spGetColor)
	reg.Register("SPRITE.SETCOLOR", "sprite", m.spSetColor)
	reg.Register("SPRITE.GETALPHA", "sprite", m.spGetAlpha)
	reg.Register("SPRITE.SETALPHA", "sprite", m.spSetAlpha)
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
		scaleX:       1,
		scaleY:       1,
		tr:           255,
		tg:           255,
		tb:           255,
		alpha:        spriteDefaultAlpha,
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
	return args[0], nil
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
	src := rl.Rectangle{
		X:      srcX,
		Y:      float32(s.srcY),
		Width:  float32(s.frameW),
		Height: float32(s.frameH),
	}
	sx := s.scaleX
	sy := s.scaleY
	if sx == 0 {
		sx = 1
	}
	if sy == 0 {
		sy = 1
	}
	dw := float32(s.frameW) * sx
	dh := float32(s.frameH) * sy
	dest := rl.Rectangle{
		X:      float32(screenX) + s.x,
		Y:      float32(screenY) + s.y,
		Width:  dw,
		Height: dh,
	}
	origin := rl.Vector2{X: s.originX * sx, Y: s.originY * sy}
	rotDeg := float32(s.rotRad * 180.0 / math.Pi)
	a := s.alpha
	if a <= 0 {
		a = spriteDefaultAlpha
	}
	if a > 1 {
		a = 1
	}
	ai := int32(a * 255)
	if ai < 0 {
		ai = 0
	}
	if ai > 255 {
		ai = 255
	}
	tint := color.RGBA{R: s.tr, G: s.tg, B: s.tb, A: uint8(ai)}
	rl.DrawTexturePro(s.tex, src, dest, origin, rotDeg, tint)
	return nil
}

func (m *Module) spGetPos(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.h == nil {
		return value.Nil, runtime.Errorf("SPRITE.GETPOS: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SPRITE.GETPOS expects (handle)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.GETPOS: invalid handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	return mbmatrix.AllocVec3Value(m.h, s.x, s.y, 0)
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
	return args[0], nil
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
	return args[0], nil
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
	return args[0], nil
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
	return args[0], nil
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
	return args[0], nil
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
	return args[0], nil
}

func clampSpriteU8(n int64) uint8 {
	if n < 0 {
		return 0
	}
	if n > 255 {
		return 255
	}
	return uint8(n)
}

func spriteEffScale(s *spriteObj) (float32, float32) {
	sx, sy := s.scaleX, s.scaleY
	if sx == 0 {
		sx = 1
	}
	if sy == 0 {
		sy = 1
	}
	return sx, sy
}

func (m *Module) spGetScale(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.h == nil {
		return value.Nil, runtime.Errorf("SPRITE.GETSCALE: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SPRITE.GETSCALE expects (handle)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.GETSCALE: invalid handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	sx, sy := spriteEffScale(s)
	return mbmatrix.AllocVec3Value(m.h, sx, sy, 1)
}

func (m *Module) spSetScale(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("SPRITE.SETSCALE expects (handle, sx#, sy#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.SETSCALE: invalid handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	x1, ok1 := argF(args[1])
	y1, ok2 := argF(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("SPRITE.SETSCALE: sx, sy must be numeric")
	}
	if x1 <= 0 || y1 <= 0 {
		return value.Nil, fmt.Errorf("SPRITE.SETSCALE: sx, sy must be > 0")
	}
	s.scaleX = x1
	s.scaleY = y1
	return args[0], nil
}

func (m *Module) spGetRot(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.h == nil {
		return value.Nil, runtime.Errorf("SPRITE.GETROT: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SPRITE.GETROT expects (handle)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.GETROT: invalid handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	return mbmatrix.AllocVec3Value(m.h, 0, 0, s.rotRad)
}

func (m *Module) spSetRot(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SPRITE.SETROT expects (handle, radians#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.SETROT: invalid handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	rad, ok := argF(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.SETROT: radians must be numeric")
	}
	s.rotRad = rad
	return args[0], nil
}

func (m *Module) spGetColor(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.h == nil {
		return value.Nil, runtime.Errorf("SPRITE.GETCOLOR: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SPRITE.GETCOLOR expects (handle)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.GETCOLOR: invalid handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	a := s.alpha
	if a <= 0 {
		a = spriteDefaultAlpha
	}
	if a > 1 {
		a = 1
	}
	instance := heap.NewInstance("Color")
	instance.SetField("r", value.FromInt(int64(s.tr)))
	instance.SetField("g", value.FromInt(int64(s.tg)))
	instance.SetField("b", value.FromInt(int64(s.tb)))
	instance.SetField("a", value.FromInt(int64(a * 255)))
	id, err := m.h.Alloc(instance)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) spSetColor(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 4 && len(args) != 5 {
		return value.Nil, fmt.Errorf("SPRITE.SETCOLOR expects (handle, r, g, b [, a])")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.SETCOLOR: invalid handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	ri, ok1 := args[1].ToInt()
	gi, ok2 := args[2].ToInt()
	bi, ok3 := args[3].ToInt()
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("SPRITE.SETCOLOR: RGB must be integer")
	}
	s.tr = clampSpriteU8(ri)
	s.tg = clampSpriteU8(gi)
	s.tb = clampSpriteU8(bi)
	if len(args) == 5 {
		if af, ok := args[4].ToFloat(); ok {
			s.alpha = float32(af)
		} else if ai, ok := args[4].ToInt(); ok {
			s.alpha = float32(ai) / 255
		} else {
			return value.Nil, fmt.Errorf("SPRITE.SETCOLOR: alpha must be numeric")
		}
		if s.alpha < 0 {
			s.alpha = 0
		}
		if s.alpha > 1 {
			s.alpha = 1
		}
	}
	return args[0], nil
}

func (m *Module) spGetAlpha(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.h == nil {
		return value.Nil, runtime.Errorf("SPRITE.GETALPHA: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SPRITE.GETALPHA expects (handle)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.GETALPHA: invalid handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	a := s.alpha
	if a <= 0 {
		a = spriteDefaultAlpha
	}
	if a > 1 {
		a = 1
	}
	return value.FromFloat(float64(a)), nil
}

func (m *Module) spSetAlpha(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SPRITE.SETALPHA expects (handle, alpha#)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.SETALPHA: invalid handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	if af, ok := args[1].ToFloat(); ok {
		a := float32(af)
		if a > 1 && a <= 255 {
			a /= 255
		}
		if a < 0 {
			a = 0
		}
		if a > 1 {
			a = 1
		}
		s.alpha = a
		return args[0], nil
	}
	if ai, ok := args[1].ToInt(); ok {
		s.alpha = float32(ai) / 255
		if s.alpha < 0 {
			s.alpha = 0
		}
		if s.alpha > 1 {
			s.alpha = 1
		}
		return args[0], nil
	}
	return value.Nil, fmt.Errorf("SPRITE.SETALPHA: alpha must be numeric")
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
		return args[0], nil
	}
	if s.rangePlaying {
		if s.rangeSpeed <= 0 {
			return args[0], nil
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
		return args[0], nil
	}
	if !s.playing || s.numFrames < 1 || s.fps <= 0 {
		return args[0], nil
	}
	s.accum += dt * s.fps
	for s.accum >= 1 {
		s.accum--
		s.curFrame = (s.curFrame + 1) % s.numFrames
	}
	return args[0], nil
}
