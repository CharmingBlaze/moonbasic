//go:build cgo || (windows && !cgo)

package mbsprite

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type spriteGroupObj struct {
	children []heap.Handle
}

func (o *spriteGroupObj) TypeName() string { return "SpriteGroup" }

func (o *spriteGroupObj) TypeTag() uint16 { return heap.TagSpriteGroup }

func (o *spriteGroupObj) Free() {}

type spriteLayerObj struct {
	z        float32
	children []heap.Handle
}

func (o *spriteLayerObj) TypeName() string { return "SpriteLayer" }

func (o *spriteLayerObj) TypeTag() uint16 { return heap.TagSpriteLayer }

func (o *spriteLayerObj) Free() {}

type batchEntry struct {
	h heap.Handle
	x int32
	y int32
}

type spriteBatchObj struct {
	entries []batchEntry
}

func (o *spriteBatchObj) TypeName() string { return "SpriteBatch" }

func (o *spriteBatchObj) TypeTag() uint16 { return heap.TagSpriteBatch }

func (o *spriteBatchObj) Free() {}

type spriteUIObj struct {
	spr heap.Handle
	ax  float32
	ay  float32
}

func (o *spriteUIObj) TypeName() string { return "SpriteUI" }

func (o *spriteUIObj) TypeTag() uint16 { return heap.TagSpriteUI }

func (o *spriteUIObj) Free() {}

type particle struct {
	x, y, vx, vy, life float32
}

type particle2DObj struct {
	max   int
	color rl.Color
	parts []particle
}

func (o *particle2DObj) TypeName() string { return "Particle2D" }

func (o *particle2DObj) TypeTag() uint16 { return heap.TagParticle2D }

func (o *particle2DObj) Free() {}

func (m *Module) registerSpriteExtras(reg runtime.Registrar) {
	reg.Register("SPRITEGROUP.CREATE", "sprite", m.sgMake)
	reg.Register("SPRITEGROUP.MAKE", "sprite", m.sgMake)
	reg.Register("SPRITEGROUP.ADD", "sprite", m.sgAdd)
	reg.Register("SPRITEGROUP.REMOVE", "sprite", m.sgRemove)
	reg.Register("SPRITEGROUP.CLEAR", "sprite", m.sgClear)
	reg.Register("SPRITEGROUP.DRAW", "sprite", m.sgDraw)
	reg.Register("SPRITEGROUP.FREE", "sprite", m.sgFree)

	reg.Register("SPRITELAYER.MAKE", "sprite", m.slMakeDeprecated)
	reg.Register("SPRITELAYER.CREATE", "sprite", m.slCreate)
	reg.Register("SPRITELAYER.ADD", "sprite", m.slAdd)
	reg.Register("SPRITELAYER.CLEAR", "sprite", m.slClear)
	reg.Register("SPRITELAYER.SETZ", "sprite", m.slSetZ)
	reg.Register("SPRITELAYER.DRAW", "sprite", m.slDraw)
	reg.Register("SPRITELAYER.FREE", "sprite", m.slFree)

	reg.Register("SPRITEBATCH.MAKE", "sprite", m.sbMakeDeprecated)
	reg.Register("SPRITEBATCH.CREATE", "sprite", m.sbCreate)
	reg.Register("SPRITEBATCH.ADD", "sprite", m.sbAdd)
	reg.Register("SPRITEBATCH.CLEAR", "sprite", m.sbClear)
	reg.Register("SPRITEBATCH.DRAW", "sprite", m.sbDraw)
	reg.Register("SPRITEBATCH.FREE", "sprite", m.sbFree)

	reg.Register("SPRITEUI.CREATE", "sprite", m.suiMake)
	reg.Register("SPRITEUI.MAKE", "sprite", m.suiMake)
	reg.Register("SPRITEUI.DRAW", "sprite", m.suiDraw)
	reg.Register("SPRITEUI.FREE", "sprite", m.suiFree)

	reg.Register("PARTICLE2D.CREATE", "sprite", m.p2Make)
	reg.Register("PARTICLE2D.MAKE", "sprite", m.p2Make)
	reg.Register("PARTICLE2D.EMIT", "sprite", m.p2Emit)
	reg.Register("PARTICLE2D.UPDATE", "sprite", m.p2Update)
	reg.Register("PARTICLE2D.DRAW", "sprite", m.p2Draw)
	reg.Register("PARTICLE2D.FREE", "sprite", m.p2Free)
}

func (m *Module) sgMake(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("SPRITEGROUP.MAKE expects 0 arguments")
	}
	if m.h == nil {
		return value.Nil, runtime.Errorf("SPRITEGROUP.MAKE: heap not bound")
	}
	o := &spriteGroupObj{}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) sgAdd(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SPRITEGROUP.ADD expects (group, sprite)")
	}
	gh, ok1 := argHandle(args[0])
	sh, ok2 := argHandle(args[1])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("SPRITEGROUP.ADD: handles required")
	}
	g, err := heap.Cast[*spriteGroupObj](m.h, gh)
	if err != nil {
		return value.Nil, err
	}
	if _, err := heap.Cast[*spriteObj](m.h, sh); err != nil {
		return value.Nil, fmt.Errorf("SPRITEGROUP.ADD: second argument must be a sprite")
	}
	g.children = append(g.children, sh)
	return value.Nil, nil
}

func (m *Module) sgRemove(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SPRITEGROUP.REMOVE expects (group, sprite)")
	}
	gh, ok1 := argHandle(args[0])
	sh, ok2 := argHandle(args[1])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("SPRITEGROUP.REMOVE: handles required")
	}
	g, err := heap.Cast[*spriteGroupObj](m.h, gh)
	if err != nil {
		return value.Nil, err
	}
	for i, h := range g.children {
		if h == sh {
			g.children = append(g.children[:i], g.children[i+1:]...)
			break
		}
	}
	return value.Nil, nil
}

func (m *Module) sgClear(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SPRITEGROUP.CLEAR expects (group)")
	}
	gh, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITEGROUP.CLEAR: invalid handle")
	}
	g, err := heap.Cast[*spriteGroupObj](m.h, gh)
	if err != nil {
		return value.Nil, err
	}
	g.children = g.children[:0]
	return value.Nil, nil
}

func (m *Module) sgDraw(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("SPRITEGROUP.DRAW expects (group, x, y)")
	}
	gh, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITEGROUP.DRAW: invalid group handle")
	}
	x, ok1 := argInt32(args[1])
	y, ok2 := argInt32(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("SPRITEGROUP.DRAW: x,y must be numeric")
	}
	g, err := heap.Cast[*spriteGroupObj](m.h, gh)
	if err != nil {
		return value.Nil, err
	}
	for _, ch := range g.children {
		s, err := heap.Cast[*spriteObj](m.h, ch)
		if err != nil {
			return value.Nil, err
		}
		if err := m.drawSpriteAtScreen(s, x, y); err != nil {
			return value.Nil, err
		}
	}
	return value.Nil, nil
}

func (m *Module) sgFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("SPRITEGROUP.FREE expects (group)")
	}
	return value.Nil, m.h.Free(heap.Handle(args[0].IVal))
}

func (m *Module) slCreate(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.slMakeWithOp(rt, args, "SPRITELAYER.CREATE")
}

func (m *Module) slMakeDeprecated(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.slMakeWithOp(rt, args, "SPRITELAYER.MAKE")
}

func (m *Module) slMakeWithOp(rt *runtime.Runtime, args []value.Value, op string) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("%s expects (z)", op)
	}
	z, ok := argF(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("%s: z must be numeric", op)
	}
	if m.h == nil {
		return value.Nil, runtime.Errorf("%s: heap not bound", op)
	}
	o := &spriteLayerObj{z: z}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) slAdd(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SPRITELAYER.ADD expects (layer, sprite)")
	}
	lh, ok1 := argHandle(args[0])
	sh, ok2 := argHandle(args[1])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("SPRITELAYER.ADD: handles required")
	}
	L, err := heap.Cast[*spriteLayerObj](m.h, lh)
	if err != nil {
		return value.Nil, err
	}
	if _, err := heap.Cast[*spriteObj](m.h, sh); err != nil {
		return value.Nil, fmt.Errorf("SPRITELAYER.ADD: sprite handle expected")
	}
	L.children = append(L.children, sh)
	return value.Nil, nil
}

func (m *Module) slClear(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SPRITELAYER.CLEAR expects (layer)")
	}
	lh, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITELAYER.CLEAR: invalid handle")
	}
	L, err := heap.Cast[*spriteLayerObj](m.h, lh)
	if err != nil {
		return value.Nil, err
	}
	L.children = L.children[:0]
	return value.Nil, nil
}

func (m *Module) slSetZ(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SPRITELAYER.SETZ expects (layer, z)")
	}
	lh, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITELAYER.SETZ: invalid layer")
	}
	z, ok2 := argF(args[1])
	if !ok2 {
		return value.Nil, fmt.Errorf("SPRITELAYER.SETZ: z must be numeric")
	}
	L, err := heap.Cast[*spriteLayerObj](m.h, lh)
	if err != nil {
		return value.Nil, err
	}
	L.z = z
	return value.Nil, nil
}

func (m *Module) slDraw(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("SPRITELAYER.DRAW expects (layer, x, y)")
	}
	lh, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITELAYER.DRAW: invalid layer")
	}
	x, ok1 := argInt32(args[1])
	y, ok2 := argInt32(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("SPRITELAYER.DRAW: x,y must be numeric")
	}
	L, err := heap.Cast[*spriteLayerObj](m.h, lh)
	if err != nil {
		return value.Nil, err
	}
	_ = L.z
	for _, ch := range L.children {
		s, err := heap.Cast[*spriteObj](m.h, ch)
		if err != nil {
			return value.Nil, err
		}
		if err := m.drawSpriteAtScreen(s, x, y); err != nil {
			return value.Nil, err
		}
	}
	return value.Nil, nil
}

func (m *Module) slFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("SPRITELAYER.FREE expects (layer)")
	}
	return value.Nil, m.h.Free(heap.Handle(args[0].IVal))
}

func (m *Module) sbCreate(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.sbMakeWithOp(rt, args, "SPRITEBATCH.CREATE")
}

func (m *Module) sbMakeDeprecated(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.sbMakeWithOp(rt, args, "SPRITEBATCH.MAKE")
}

func (m *Module) sbMakeWithOp(rt *runtime.Runtime, args []value.Value, op string) (value.Value, error) {
	_ = rt
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("%s expects 0 arguments", op)
	}
	if m.h == nil {
		return value.Nil, runtime.Errorf("%s: heap not bound", op)
	}
	o := &spriteBatchObj{}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) sbAdd(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("SPRITEBATCH.ADD expects (batch, sprite, x, y)")
	}
	bh, ok1 := argHandle(args[0])
	sh, ok2 := argHandle(args[1])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("SPRITEBATCH.ADD: batch and sprite handles required")
	}
	x, ok3 := argInt32(args[2])
	y, ok4 := argInt32(args[3])
	if !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("SPRITEBATCH.ADD: x,y must be numeric")
	}
	b, err := heap.Cast[*spriteBatchObj](m.h, bh)
	if err != nil {
		return value.Nil, err
	}
	if _, err := heap.Cast[*spriteObj](m.h, sh); err != nil {
		return value.Nil, fmt.Errorf("SPRITEBATCH.ADD: sprite handle expected")
	}
	b.entries = append(b.entries, batchEntry{h: sh, x: x, y: y})
	return value.Nil, nil
}

func (m *Module) sbClear(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SPRITEBATCH.CLEAR expects (batch)")
	}
	bh, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITEBATCH.CLEAR: invalid batch")
	}
	b, err := heap.Cast[*spriteBatchObj](m.h, bh)
	if err != nil {
		return value.Nil, err
	}
	b.entries = b.entries[:0]
	return value.Nil, nil
}

func (m *Module) sbDraw(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SPRITEBATCH.DRAW expects (batch)")
	}
	bh, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITEBATCH.DRAW: invalid batch")
	}
	b, err := heap.Cast[*spriteBatchObj](m.h, bh)
	if err != nil {
		return value.Nil, err
	}
	for _, e := range b.entries {
		s, err := heap.Cast[*spriteObj](m.h, e.h)
		if err != nil {
			return value.Nil, err
		}
		if err := m.drawSpriteAtScreen(s, e.x, e.y); err != nil {
			return value.Nil, err
		}
	}
	return value.Nil, nil
}

func (m *Module) sbFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("SPRITEBATCH.FREE expects (batch)")
	}
	return value.Nil, m.h.Free(heap.Handle(args[0].IVal))
}

func (m *Module) suiMake(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("SPRITEUI.MAKE expects (sprite, anchorX, anchorY)")
	}
	sh, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITEUI.MAKE: sprite handle required")
	}
	ax, ok1 := argF(args[1])
	ay, ok2 := argF(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("SPRITEUI.MAKE: anchors must be numeric")
	}
	if _, err := heap.Cast[*spriteObj](m.h, sh); err != nil {
		return value.Nil, err
	}
	if m.h == nil {
		return value.Nil, runtime.Errorf("SPRITEUI.MAKE: heap not bound")
	}
	o := &spriteUIObj{spr: sh, ax: ax, ay: ay}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) suiDraw(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("SPRITEUI.DRAW expects (ui, screenW, screenH)")
	}
	uh, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITEUI.DRAW: invalid handle")
	}
	sw, ok1 := argInt32(args[1])
	sh, ok2 := argInt32(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("SPRITEUI.DRAW: screen size must be numeric")
	}
	ui, err := heap.Cast[*spriteUIObj](m.h, uh)
	if err != nil {
		return value.Nil, err
	}
	s, err := heap.Cast[*spriteObj](m.h, ui.spr)
	if err != nil {
		return value.Nil, err
	}
	fw := float32(s.frameW)
	fh := float32(s.frameH)
	sx := int32(float32(sw)*ui.ax - fw*ui.ax)
	sy := int32(float32(sh)*ui.ay - fh*ui.ay)
	if err := m.drawSpriteAtScreen(s, sx, sy); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) suiFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("SPRITEUI.FREE expects (ui)")
	}
	return value.Nil, m.h.Free(heap.Handle(args[0].IVal))
}

func (m *Module) p2Make(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("PARTICLE2D.MAKE expects (max, r, g, b, a)")
	}
	maxN, ok0 := args[0].ToInt()
	if !ok0 || maxN < 1 {
		return value.Nil, fmt.Errorf("PARTICLE2D.MAKE: max must be a positive integer")
	}
	r, ok1 := argInt32(args[1])
	g, ok2 := argInt32(args[2])
	b, ok3 := argInt32(args[3])
	a, ok4 := argInt32(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("PARTICLE2D.MAKE: rgba must be numeric")
	}
	if m.h == nil {
		return value.Nil, runtime.Errorf("PARTICLE2D.MAKE: heap not bound")
	}
	o := &particle2DObj{
		max:   int(maxN),
		color: rl.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)},
		parts: nil,
	}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) p2Emit(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("PARTICLE2D.EMIT expects (p, x, y, vx, vy, life)")
	}
	ph, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("PARTICLE2D.EMIT: invalid handle")
	}
	x, ok1 := argF(args[1])
	y, ok2 := argF(args[2])
	vx, ok3 := argF(args[3])
	vy, ok4 := argF(args[4])
	life, ok5 := argF(args[5])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("PARTICLE2D.EMIT: numeric arguments required")
	}
	po, err := heap.Cast[*particle2DObj](m.h, ph)
	if err != nil {
		return value.Nil, err
	}
	if len(po.parts) >= po.max {
		return value.Nil, nil
	}
	po.parts = append(po.parts, particle{x: x, y: y, vx: vx, vy: vy, life: life})
	return value.Nil, nil
}

func (m *Module) p2Update(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PARTICLE2D.UPDATE expects (p, dt)")
	}
	ph, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("PARTICLE2D.UPDATE: invalid handle")
	}
	dt, ok2 := argF(args[1])
	if !ok2 {
		return value.Nil, fmt.Errorf("PARTICLE2D.UPDATE: dt must be numeric")
	}
	po, err := heap.Cast[*particle2DObj](m.h, ph)
	if err != nil {
		return value.Nil, err
	}
	dst := po.parts[:0]
	for _, q := range po.parts {
		q.life -= dt
		q.x += q.vx * dt
		q.y += q.vy * dt
		if q.life > 0 {
			dst = append(dst, q)
		}
	}
	po.parts = dst
	return value.Nil, nil
}

func (m *Module) p2Draw(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PARTICLE2D.DRAW expects (p)")
	}
	ph, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("PARTICLE2D.DRAW: invalid handle")
	}
	po, err := heap.Cast[*particle2DObj](m.h, ph)
	if err != nil {
		return value.Nil, err
	}
	for _, q := range po.parts {
		rl.DrawCircleV(rl.Vector2{X: q.x, Y: q.y}, 3, po.color)
	}
	return value.Nil, nil
}

func (m *Module) p2Free(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("PARTICLE2D.FREE expects (p)")
	}
	return value.Nil, m.h.Free(heap.Handle(args[0].IVal))
}
