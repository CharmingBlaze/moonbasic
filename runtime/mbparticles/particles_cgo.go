//go:build cgo

package mbparticles

import (
	"fmt"
	"image/color"
	"math/rand/v2"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	mbdraw "moonbasic/runtime/draw"
	"moonbasic/runtime/mbmodel3d"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

const maxParticles = 16000

var (
	partDefMu     sync.Mutex
	partDefTex    rl.Texture2D
	partDefTexOK  bool
)

func particleFallbackTex() rl.Texture2D {
	partDefMu.Lock()
	defer partDefMu.Unlock()
	if !partDefTexOK {
		img := rl.GenImageColor(1, 1, rl.White)
		partDefTex = rl.LoadTextureFromImage(img)
		rl.UnloadImage(img)
		partDefTexOK = true
	}
	return partDefTex
}

func unloadParticleFallbackTex() {
	partDefMu.Lock()
	defer partDefMu.Unlock()
	if partDefTexOK {
		rl.UnloadTexture(partDefTex)
		partDefTex = rl.Texture2D{}
		partDefTexOK = false
	}
}

type particle struct {
	x, y, z    float32
	vx, vy, vz float32
	life, age  float32
}

type particleObj struct {
	texH heap.Handle

	emitRate float32
	emitAcc  float32

	lifeMin, lifeMax float32

	vx0, vy0, vz0 float32
	vspread       float32

	sr, sg, sb, sa uint8
	er, eg, eb, ea uint8

	sizeStart, sizeEnd float32
	gravity            float32

	px, py, pz float32

	playing bool
	parts   []particle
}

func (o *particleObj) TypeName() string { return "Particle" }

func (o *particleObj) TypeTag() uint16 { return heap.TagParticle }

func (o *particleObj) Free() {
	o.parts = nil
}

func (m *Module) requireHeap() error {
	if m.h == nil {
		return runtime.Errorf("PARTICLE.* builtins: heap not bound")
	}
	return nil
}

func argFloat(v value.Value) (float32, bool) {
	if f, ok := v.ToFloat(); ok {
		return float32(f), true
	}
	if i, ok := v.ToInt(); ok {
		return float32(i), true
	}
	return 0, false
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

func clampU8(n int32) uint8 {
	if n < 0 {
		return 0
	}
	if n > 255 {
		return 255
	}
	return uint8(n)
}

func (m *Module) getParticle(args []value.Value, ix int, op string) (*particleObj, error) {
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: expected particle handle", op)
	}
	return heap.Cast[*particleObj](m.h, heap.Handle(args[ix].IVal))
}

func (o *particleObj) spawnOne() {
	if len(o.parts) >= maxParticles {
		return
	}
	life := o.lifeMin + (o.lifeMax-o.lifeMin)*rand.Float32()
	if life <= 0 {
		life = 0.01
	}
	sp := o.vspread
	vx := o.vx0 + (rand.Float32()*2-1)*sp
	vy := o.vy0 + (rand.Float32()*2-1)*sp
	vz := o.vz0 + (rand.Float32()*2-1)*sp
	o.parts = append(o.parts, particle{
		x: o.px, y: o.py, z: o.pz,
		vx: vx, vy: vy, vz: vz,
		life: life, age: 0,
	})
}

func (o *particleObj) update(dt float32) {
	// Integrate & cull dead
	dst := o.parts[:0]
	for i := range o.parts {
		p := &o.parts[i]
		p.age += dt
		if p.age >= p.life {
			continue
		}
		p.vy += o.gravity * dt
		p.x += p.vx * dt
		p.y += p.vy * dt
		p.z += p.vz * dt
		dst = append(dst, *p)
	}
	o.parts = dst

	if o.playing && o.emitRate > 0 && dt > 0 {
		o.emitAcc += o.emitRate * dt
		for o.emitAcc >= 1 && len(o.parts) < maxParticles {
			o.spawnOne()
			o.emitAcc -= 1
		}
	}
}

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	reg.Register("PARTICLE.MAKE", "particle", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("PARTICLE.MAKE expects no arguments")
		}
		o := &particleObj{
			lifeMin: 1, lifeMax: 1,
			sr: 255, sg: 255, sb: 255, sa: 255,
			er: 255, eg: 255, eb: 255, ea: 0,
			sizeStart: 0.2, sizeEnd: 0,
		}
		id, err := m.h.Alloc(o)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	})

	reg.Register("PARTICLE.FREE", "particle", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 || args[0].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("PARTICLE.FREE expects particle handle")
		}
		m.h.Free(heap.Handle(args[0].IVal))
		return value.Nil, nil
	}))

	reg.Register("PARTICLE.SETTEXTURE", "particle", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("PARTICLE.SETTEXTURE expects (particle, textureHandle)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.SETTEXTURE")
		if err != nil {
			return value.Nil, err
		}
		if args[1].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("PARTICLE.SETTEXTURE: texture must be a handle")
		}
		o.texH = heap.Handle(args[1].IVal)
		return value.Nil, nil
	}))

	reg.Register("PARTICLE.SETEMITRATE", "particle", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("PARTICLE.SETEMITRATE expects (particle, rate)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.SETEMITRATE")
		if err != nil {
			return value.Nil, err
		}
		rate, ok := argFloat(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("PARTICLE.SETEMITRATE: rate must be numeric")
		}
		if rate < 0 {
			rate = 0
		}
		o.emitRate = rate
		return value.Nil, nil
	}))

	reg.Register("PARTICLE.SETLIFETIME", "particle", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("PARTICLE.SETLIFETIME expects (particle, min, max)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.SETLIFETIME")
		if err != nil {
			return value.Nil, err
		}
		a, ok1 := argFloat(args[1])
		b, ok2 := argFloat(args[2])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("PARTICLE.SETLIFETIME: min, max must be numeric")
		}
		if a <= 0 {
			a = 0.01
		}
		if b <= 0 {
			b = 0.01
		}
		if a > b {
			a, b = b, a
		}
		o.lifeMin, o.lifeMax = a, b
		return value.Nil, nil
	}))

	reg.Register("PARTICLE.SETVELOCITY", "particle", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("PARTICLE.SETVELOCITY expects (particle, vx, vy, vz, spread)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.SETVELOCITY")
		if err != nil {
			return value.Nil, err
		}
		vx, ok1 := argFloat(args[1])
		vy, ok2 := argFloat(args[2])
		vz, ok3 := argFloat(args[3])
		sp, ok4 := argFloat(args[4])
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return value.Nil, fmt.Errorf("PARTICLE.SETVELOCITY: arguments must be numeric")
		}
		if sp < 0 {
			sp = 0
		}
		o.vx0, o.vy0, o.vz0 = vx, vy, vz
		o.vspread = sp
		return value.Nil, nil
	}))

	reg.Register("PARTICLE.SETCOLOR", "particle", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("PARTICLE.SETCOLOR expects (particle, r, g, b, a)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.SETCOLOR")
		if err != nil {
			return value.Nil, err
		}
		r0, ok1 := argInt32(args[1])
		g0, ok2 := argInt32(args[2])
		b0, ok3 := argInt32(args[3])
		a0, ok4 := argInt32(args[4])
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return value.Nil, fmt.Errorf("PARTICLE.SETCOLOR: r,g,b,a must be numeric")
		}
		o.sr, o.sg, o.sb, o.sa = clampU8(r0), clampU8(g0), clampU8(b0), clampU8(a0)
		return value.Nil, nil
	}))

	reg.Register("PARTICLE.SETCOLOREND", "particle", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("PARTICLE.SETCOLOREND expects (particle, r, g, b, a)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.SETCOLOREND")
		if err != nil {
			return value.Nil, err
		}
		r0, ok1 := argInt32(args[1])
		g0, ok2 := argInt32(args[2])
		b0, ok3 := argInt32(args[3])
		a0, ok4 := argInt32(args[4])
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return value.Nil, fmt.Errorf("PARTICLE.SETCOLOREND: r,g,b,a must be numeric")
		}
		o.er, o.eg, o.eb, o.ea = clampU8(r0), clampU8(g0), clampU8(b0), clampU8(a0)
		return value.Nil, nil
	}))

	reg.Register("PARTICLE.SETSIZE", "particle", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("PARTICLE.SETSIZE expects (particle, startSize, endSize)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.SETSIZE")
		if err != nil {
			return value.Nil, err
		}
		s0, ok1 := argFloat(args[1])
		s1, ok2 := argFloat(args[2])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("PARTICLE.SETSIZE: sizes must be numeric")
		}
		if s0 < 0 {
			s0 = 0
		}
		if s1 < 0 {
			s1 = 0
		}
		o.sizeStart, o.sizeEnd = s0, s1
		return value.Nil, nil
	}))

	reg.Register("PARTICLE.SETGRAVITY", "particle", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("PARTICLE.SETGRAVITY expects (particle, g)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.SETGRAVITY")
		if err != nil {
			return value.Nil, err
		}
		g, ok := argFloat(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("PARTICLE.SETGRAVITY: g must be numeric")
		}
		o.gravity = g
		return value.Nil, nil
	}))

	reg.Register("PARTICLE.SETPOS", "particle", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("PARTICLE.SETPOS expects (particle, x, y, z)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.SETPOS")
		if err != nil {
			return value.Nil, err
		}
		x, ok1 := argFloat(args[1])
		y, ok2 := argFloat(args[2])
		z, ok3 := argFloat(args[3])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("PARTICLE.SETPOS: x, y, z must be numeric")
		}
		o.px, o.py, o.pz = x, y, z
		return value.Nil, nil
	}))

	reg.Register("PARTICLE.PLAY", "particle", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("PARTICLE.PLAY expects particle handle")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.PLAY")
		if err != nil {
			return value.Nil, err
		}
		o.playing = true
		return value.Nil, nil
	}))

	reg.Register("PARTICLE.UPDATE", "particle", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("PARTICLE.UPDATE expects (particle, dt)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.UPDATE")
		if err != nil {
			return value.Nil, err
		}
		dt, ok := argFloat(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("PARTICLE.UPDATE: dt must be numeric")
		}
		if dt < 0 {
			dt = 0
		}
		if dt > 0.25 {
			dt = 0.25
		}
		o.update(dt)
		return value.Nil, nil
	}))

	reg.Register("PARTICLE.DRAW", "particle", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("PARTICLE.DRAW expects particle handle")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.DRAW")
		if err != nil {
			return value.Nil, err
		}
		cam, ok := mbmodel3d.ActiveCamera3D()
		if !ok {
			return value.Nil, fmt.Errorf("PARTICLE.DRAW: must be called inside CAMERA.Begin ... CAMERA.End")
		}
		var tex rl.Texture2D
		if o.texH != 0 {
			tex, err = mbdraw.TextureForBinding(m.h, o.texH)
			if err != nil {
				return value.Nil, fmt.Errorf("PARTICLE.DRAW: %w", err)
			}
		} else {
			tex = particleFallbackTex()
		}
		for i := range o.parts {
			p := &o.parts[i]
			if p.life <= 0 {
				continue
			}
			t := p.age / p.life
			if t > 1 {
				t = 1
			}
			inv := 1 - t
			r := uint8(float32(o.sr)*inv + float32(o.er)*t)
			g := uint8(float32(o.sg)*inv + float32(o.eg)*t)
			b := uint8(float32(o.sb)*inv + float32(o.eb)*t)
			a := uint8(float32(o.sa)*inv + float32(o.ea)*t)
			sz := o.sizeStart + (o.sizeEnd-o.sizeStart)*t
			if sz < 0 {
				sz = 0
			}
			if a == 0 {
				continue
			}
			rl.DrawBillboard(cam, tex, rl.Vector3{X: p.x, Y: p.y, Z: p.z}, sz, color.RGBA{R: r, G: g, B: b, A: a})
		}
		return value.Nil, nil
	}))
}

// Shutdown releases module-level GPU resources.
func (m *Module) Shutdown() { unloadParticleFallbackTex() }
