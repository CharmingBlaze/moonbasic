//go:build cgo || (windows && !cgo)

package mbparticles

import (
	"fmt"
	"image/color"
	"math/rand/v2"
	"strings"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	mbcamera "moonbasic/runtime/camera"
	mbdraw "moonbasic/runtime/draw"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/runtime/mbmodel3d"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

const maxParticles = 16000

var (
	partDefMu    sync.Mutex
	partDefTex   rl.Texture2D
	partDefTexOK bool
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
	x, y, z      float32
	vx, vy, vz   float32
	life, age    float32
	size0, size1 float32
}

type particleObj struct {
	texH heap.Handle

	emitRate float32
	emitAcc  float32

	lifeMin, lifeMax float32

	vx0, vy0, vz0 float32
	vspread       float32

	speedMin, speedMax float32

	sr, sg, sb, sa uint8
	er, eg, eb, ea uint8

	sizeStartMin, sizeStartMax float32
	sizeEndMin, sizeEndMax     float32

	gx, gy, gz float32

	px, py, pz float32

	playing   bool
	billboard bool

	parts []particle
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

func truthy(v value.Value) bool {
	if v.Kind == value.KindInt {
		return v.IVal != 0
	}
	if v.Kind == value.KindFloat {
		return v.FVal != 0
	}
	return false
}

func (m *Module) getParticle(args []value.Value, ix int, op string) (*particleObj, error) {
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: expected particle handle", op)
	}
	return heap.Cast[*particleObj](m.h, heap.Handle(args[ix].IVal))
}

// regDual registers the same handler under PARTICLE.name and PARTICLE3D.name.
func regDual(reg runtime.Registrar, particleKey string, fn runtime.BuiltinFn) {
	reg.Register(particleKey, "particle", fn)
	k3 := strings.Replace(particleKey, "PARTICLE.", "PARTICLE3D.", 1)
	if k3 == particleKey {
		panic("regDual: expected PARTICLE. prefix: " + particleKey)
	}
	reg.Register(k3, "particle", fn)
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
	spd := o.speedMin + (o.speedMax-o.speedMin)*rand.Float32()
	if spd < 0 {
		spd = 0
	}
	vx *= spd
	vy *= spd
	vz *= spd

	s0 := o.sizeStartMin + (o.sizeStartMax-o.sizeStartMin)*rand.Float32()
	s1 := o.sizeEndMin + (o.sizeEndMax-o.sizeEndMin)*rand.Float32()
	if s0 < 0 {
		s0 = 0
	}
	if s1 < 0 {
		s1 = 0
	}

	o.parts = append(o.parts, particle{
		x: o.px, y: o.py, z: o.pz,
		vx: vx, vy: vy, vz: vz,
		life: life, age: 0,
		size0: s0, size1: s1,
	})
}

// Update steps simulation; exported so ENTITY.UPDATE can discover TagParticle via
// interface { Update(float32) } (unexported update() was never called).
func (o *particleObj) Update(dt float32) {
	o.update(dt)
}

func (o *particleObj) update(dt float32) {
	dst := o.parts[:0]
	for i := range o.parts {
		p := &o.parts[i]
		p.age += dt
		if p.age >= p.life {
			continue
		}
		p.vx += o.gx * dt
		p.vy += o.gy * dt
		p.vz += o.gz * dt
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
	// Define all handlers first to enable multiple registrations
	makeFn := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("PARTICLE emitter constructor expects no arguments")
		}
		o := &particleObj{
			lifeMin: 1, lifeMax: 1,
			sr: 255, sg: 255, sb: 255, sa: 255,
			er: 255, eg: 255, eb: 255, ea: 0,
			sizeStartMin: 0.2, sizeStartMax: 0.2,
			sizeEndMin: 0, sizeEndMax: 0,
			speedMin: 1, speedMax: 1,
			billboard: true,
		}
		id, err := m.h.Alloc(o)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}

	freeFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 || args[0].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("PARTICLE.FREE expects particle handle")
		}
		m.h.Free(heap.Handle(args[0].IVal))
		return value.Nil, nil
	})

	setTexFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
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
		return args[0], nil
	})

	setRateFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
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
		return args[0], nil
	})

	setLifeFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
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
		return args[0], nil
	})

	setVelFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
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
		return args[0], nil
	})

	setDirFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("PARTICLE.SETDIRECTION expects (particle, vx, vy, vz)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.SETDIRECTION")
		if err != nil {
			return value.Nil, err
		}
		vx, ok1 := argFloat(args[1])
		vy, ok2 := argFloat(args[2])
		vz, ok3 := argFloat(args[3])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("PARTICLE.SETDIRECTION: values must be numeric")
		}
		o.vx0, o.vy0, o.vz0 = vx, vy, vz
		return args[0], nil
	})

	setSpreadFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("PARTICLE.SETSPREAD expects (particle, angle)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.SETSPREAD")
		if err != nil {
			return value.Nil, err
		}
		sp, ok := argFloat(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("PARTICLE.SETSPREAD: angle must be numeric")
		}
		if sp < 0 {
			sp = 0
		}
		o.vspread = sp
		return args[0], nil
	})

	setSpeedFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("PARTICLE.SETSPEED expects (particle, min#, max#)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.SETSPEED")
		if err != nil {
			return value.Nil, err
		}
		a, ok1 := argFloat(args[1])
		b, ok2 := argFloat(args[2])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("PARTICLE.SETSPEED: min, max must be numeric")
		}
		if a < 0 {
			a = 0
		}
		if b < 0 {
			b = 0
		}
		if a > b {
			a, b = b, a
		}
		o.speedMin, o.speedMax = a, b
		return args[0], nil
	})

	setColorFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
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
		return args[0], nil
	})

	setColorEndFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
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
		return args[0], nil
	})

	setSizeFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
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
		o.sizeStartMin, o.sizeStartMax = s0, s0
		o.sizeEndMin, o.sizeEndMax = s1, s1
		return args[0], nil
	})

	setGravFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		o, err := m.getParticle(args, 0, "PARTICLE.SETGRAVITY")
		if err != nil {
			return value.Nil, err
		}
		if len(args) == 2 {
			g, ok := argFloat(args[1])
			if !ok {
				return value.Nil, fmt.Errorf("PARTICLE.SETGRAVITY: g must be numeric")
			}
			o.gx, o.gy, o.gz = 0, g, 0
			return args[0], nil
		}
		if len(args) == 4 {
			gx, ok1 := argFloat(args[1])
			gy, ok2 := argFloat(args[2])
			gz, ok3 := argFloat(args[3])
			if !ok1 || !ok2 || !ok3 {
				return value.Nil, fmt.Errorf("PARTICLE.SETGRAVITY: gx, gy, gz must be numeric")
			}
			o.gx, o.gy, o.gz = gx, gy, gz
			return args[0], nil
		}
		return value.Nil, fmt.Errorf("PARTICLE.SETGRAVITY expects (particle, g) or (particle, gx, gy, gz)")
	})

	setPosFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
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
		return args[0], nil
	})

	getPosFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("PARTICLE.GETPOS expects (particle)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.GETPOS")
		if err != nil {
			return value.Nil, err
		}
		return mbmatrix.AllocVec3Value(m.h, o.px, o.py, o.pz)
	})

	getColorFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("PARTICLE.GETCOLOR expects (particle)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.GETCOLOR")
		if err != nil {
			return value.Nil, err
		}
		instance := heap.NewInstance("Color")
		instance.SetField("r", value.FromInt(int64(o.sr)))
		instance.SetField("g", value.FromInt(int64(o.sg)))
		instance.SetField("b", value.FromInt(int64(o.sb)))
		instance.SetField("a", value.FromInt(int64(o.sa)))
		id, err := m.h.Alloc(instance)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	})

	getAlphaFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("PARTICLE.GETALPHA expects (particle)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.GETALPHA")
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(float64(o.sa) / 255.0), nil
	})

	setAlphaFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("PARTICLE.SETALPHA expects (particle, a)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.SETALPHA")
		if err != nil {
			return value.Nil, err
		}
		a0, ok := argInt32(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("PARTICLE.SETALPHA: a must be numeric")
		}
		o.sa = clampU8(a0)
		return args[0], nil
	})

	getVelFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("PARTICLE.GETVELOCITY expects (particle)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.GETVELOCITY")
		if err != nil {
			return value.Nil, err
		}
		return mbmatrix.AllocVec3Value(m.h, o.vx0, o.vy0, o.vz0)
	})

	getSizeFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("PARTICLE.GETSIZE expects (particle)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.GETSIZE")
		if err != nil {
			return value.Nil, err
		}
		// Matches PARTICLE.SETSIZE(start,end): lower bounds of start/end ranges (see SETSTARTSIZE/SETENDSIZE for ranges).
		return mbmatrix.AllocVec2Value(m.h, o.sizeStartMin, o.sizeEndMin)
	})

	setBurstFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("PARTICLE.SETBURST expects (particle, count)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.SETBURST")
		if err != nil {
			return value.Nil, err
		}
		n, ok := args[1].ToInt()
		if !ok {
			if f, okf := args[1].ToFloat(); okf {
				n = int64(f)
			} else {
				return value.Nil, fmt.Errorf("PARTICLE.SETBURST: count must be numeric")
			}
		}
		if n < 0 {
			n = 0
		}
		if n > 100000 {
			n = 100000
		}
		for i := int64(0); i < n && len(o.parts) < maxParticles; i++ {
			o.spawnOne()
		}
		return args[0], nil
	})

	setBillFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("PARTICLE.SETBILLBOARD expects (particle, enable)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.SETBILLBOARD")
		if err != nil {
			return value.Nil, err
		}
		o.billboard = truthy(args[1])
		return args[0], nil
	})

	playFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
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
		return args[0], nil
	})

	stopFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("PARTICLE.STOP expects particle handle")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.STOP")
		if err != nil {
			return value.Nil, err
		}
		o.playing = false
		return args[0], nil
	})

	updateFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
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
		return args[0], nil
	})

	drawFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 && len(args) != 2 {
			return value.Nil, fmt.Errorf("PARTICLE.DRAW expects particle or (particle, camera)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.DRAW")
		if err != nil {
			return value.Nil, err
		}
		var cam rl.Camera3D
		if len(args) == 1 {
			var ok bool
			cam, ok = mbmodel3d.ActiveCamera3D()
			if !ok {
				return value.Nil, fmt.Errorf("PARTICLE.DRAW: use CAMERA.Begin/End or pass camera handle")
			}
		} else {
			if args[1].Kind != value.KindHandle {
				return value.Nil, fmt.Errorf("PARTICLE.DRAW: camera must be a handle")
			}
			cam, err = mbcamera.RayCamera3D(m.h, heap.Handle(args[1].IVal))
			if err != nil {
				return value.Nil, fmt.Errorf("PARTICLE.DRAW: %w", err)
			}
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
			sz := p.size0 + (p.size1-p.size0)*t
			if sz < 0 {
				sz = 0
			}
			if a == 0 {
				continue
			}
			pos := rl.Vector3{X: p.x, Y: p.y, Z: p.z}
			if o.billboard {
				rl.DrawBillboard(cam, tex, pos, sz, color.RGBA{R: r, G: g, B: b, A: a})
			} else {
				rl.DrawCube(pos, sz*0.5, sz*0.5, sz*0.5, color.RGBA{R: r, G: g, B: b, A: a})
			}
		}
		return args[0], nil
	})

	isAliveFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("PARTICLE.ISALIVE expects particle handle")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.ISALIVE")
		if err != nil {
			return value.Nil, err
		}
		alive := o.playing || len(o.parts) > 0
		if alive {
			return value.FromInt(1), nil
		}
		return value.FromInt(0), nil
	})

	countFn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("PARTICLE.COUNT expects particle handle")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.COUNT")
		if err != nil {
			return value.Nil, err
		}
		return value.FromInt(int64(len(o.parts))), nil
	})

	setSize0Fn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("PARTICLE.SETSTARTSIZE expects (particle, min#, max#)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.SETSTARTSIZE")
		if err != nil {
			return value.Nil, err
		}
		a, ok1 := argFloat(args[1])
		b, ok2 := argFloat(args[2])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("PARTICLE.SETSTARTSIZE: min, max must be numeric")
		}
		if a > b {
			a, b = b, a
		}
		if a < 0 {
			a = 0
		}
		if b < 0 {
			b = 0
		}
		o.sizeStartMin, o.sizeStartMax = a, b
		return args[0], nil
	})

	setSize1Fn := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("PARTICLE.SETENDSIZE expects (particle, min#, max#)")
		}
		o, err := m.getParticle(args, 0, "PARTICLE.SETENDSIZE")
		if err != nil {
			return value.Nil, err
		}
		a, ok1 := argFloat(args[1])
		b, ok2 := argFloat(args[2])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("PARTICLE.SETENDSIZE: min, max must be numeric")
		}
		if a > b {
			a, b = b, a
		}
		if a < 0 {
			a = 0
		}
		if b < 0 {
			b = 0
		}
		o.sizeEndMin, o.sizeEndMax = a, b
		return args[0], nil
	})

	// Register all variations
	regDual(reg, "PARTICLE.MAKE", makeFn)
	regDual(reg, "PARTICLE.CREATE", makeFn)
	reg.Register("CREATEEMITTER", "particle", makeFn)

	regDual(reg, "PARTICLE.FREE", freeFn)
	regDual(reg, "PARTICLE.SETTEXTURE", setTexFn)

	regDual(reg, "PARTICLE.SETEMITRATE", setRateFn)
	reg.Register("PARTICLEEMITRATE", "particle", setRateFn)
	reg.Register("PARTICLE.SETRATE", "particle", setRateFn)
	reg.Register("PARTICLE3D.SETRATE", "particle", setRateFn)

	regDual(reg, "PARTICLE.SETLIFETIME", setLifeFn)
	reg.Register("PARTICLELIFE", "particle", setLifeFn)

	regDual(reg, "PARTICLE.SETVELOCITY", setVelFn)
	reg.Register("PARTICLEVELOCITY", "particle", setVelFn)

	regDual(reg, "PARTICLE.SETDIRECTION", setDirFn)
	regDual(reg, "PARTICLE.SETSPREAD", setSpreadFn)

	regDual(reg, "PARTICLE.SETSPEED", setSpeedFn)
	reg.Register("PARTICLESPEED", "particle", setSpeedFn)

	regDual(reg, "PARTICLE.SETSTARTSIZE", setSize0Fn)
	regDual(reg, "PARTICLE.SETENDSIZE", setSize1Fn)

	regDual(reg, "PARTICLE.SETCOLOR", setColorFn)
	reg.Register("PARTICLECOLOR", "particle", setColorFn)
	reg.Register("PARTICLE.SETSTARTCOLOR", "particle", setColorFn)
	reg.Register("PARTICLE3D.SETSTARTCOLOR", "particle", setColorFn)

	regDual(reg, "PARTICLE.SETCOLOREND", setColorEndFn)
	reg.Register("PARTICLE.SETENDCOLOR", "particle", setColorEndFn)
	reg.Register("PARTICLE3D.SETENDCOLOR", "particle", setColorEndFn)

	regDual(reg, "PARTICLE.SETSIZE", setSizeFn)
	regDual(reg, "PARTICLE.SETGRAVITY", setGravFn)

	regDual(reg, "PARTICLE.SETPOS", setPosFn)
	regDual(reg, "PARTICLE.GETPOS", getPosFn)
	regDual(reg, "PARTICLE.GETCOLOR", getColorFn)
	regDual(reg, "PARTICLE.GETALPHA", getAlphaFn)
	regDual(reg, "PARTICLE.SETALPHA", setAlphaFn)
	regDual(reg, "PARTICLE.GETVELOCITY", getVelFn)
	regDual(reg, "PARTICLE.GETSIZE", getSizeFn)
	reg.Register("EMITTERPOS", "particle", setPosFn)

	regDual(reg, "PARTICLE.SETBURST", setBurstFn)
	reg.Register("EMITPARTICLE", "particle", setBurstFn)

	regDual(reg, "PARTICLE.SETBILLBOARD", setBillFn)
	regDual(reg, "PARTICLE.PLAY", playFn)
	regDual(reg, "PARTICLE.STOP", stopFn)

	regDual(reg, "PARTICLE.UPDATE", updateFn)
	reg.Register("UPDATEEMITTER", "particle", updateFn)

	regDual(reg, "PARTICLE.DRAW", drawFn)
	reg.Register("DRAWEMITTER", "particle", drawFn)
	reg.Register("PARTICLES.DRAWEMITTER", "particle", drawFn)

	regDual(reg, "PARTICLE.ISALIVE", isAliveFn)
	reg.Register("EMITTERALIVE", "particle", isAliveFn)

	regDual(reg, "PARTICLE.COUNT", countFn)
	reg.Register("EMITTERCOUNT", "particle", countFn)
}

// Shutdown releases module-level GPU resources.
func (m *Module) Shutdown() { unloadParticleFallbackTex() }
