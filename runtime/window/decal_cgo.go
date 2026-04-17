//go:build cgo || (windows && !cgo)

package window

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	mbdraw "moonbasic/runtime/draw"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/runtime/mbmodel3d"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type decalObj struct {
	texH heap.Handle
	pos  rl.Vector3
	sx   float32
	sy   float32
	life float32 // seconds; 0 = no fade-out
	age  float32
}

func (d *decalObj) TypeName() string { return "Decal" }

func (d *decalObj) TypeTag() uint16 { return heap.TagDecal }

func (d *decalObj) Free() {}

func (m *Module) registerDecalCommands(r runtime.Registrar) {
	r.Register("DECAL.CREATE", "decal", m.decalMake)
	r.Register("DECAL.MAKE", "decal", m.decalMake)
	r.Register("DECAL.FREE", "decal", m.decalFree)
	r.Register("DECAL.SETPOS", "decal", m.decalSetPos)
	r.Register("DECAL.GETPOS", "decal", m.decalGetPos)
	r.Register("DECAL.SETSIZE", "decal", m.decalSetSize)
	r.Register("DECAL.GETSIZE", "decal", m.decalGetSize)
	r.Register("DECAL.SETLIFETIME", "decal", m.decalSetLifetime)
	r.Register("DECAL.GETLIFETIME", "decal", m.decalGetLifetime)
	r.Register("DECAL.GETROT", "decal", m.decalGetRot)
	r.Register("DECAL.DRAW", "decal", m.decalDraw)
}

func (m *Module) requireHeapDecal(rt *runtime.Runtime) (*heap.Store, error) {
	if rt != nil && rt.Heap != nil {
		return rt.Heap, nil
	}
	if m.h != nil {
		return m.h, nil
	}
	return nil, fmt.Errorf("DECAL.*: heap not bound")
}

func (m *Module) getDecal(store *heap.Store, args []value.Value, at int, op string) (*decalObj, error) {
	if len(args) <= at || args[at].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: expected decal handle", op)
	}
	return heap.Cast[*decalObj](store, heap.Handle(args[at].IVal))
}

func (m *Module) decalMake(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	store, err := m.requireHeapDecal(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DECAL.MAKE expects texture handle (from TEXTURE.LOAD)")
	}
	if _, err := mbdraw.TextureForBinding(store, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, fmt.Errorf("DECAL.MAKE: %w", err)
	}
	d := &decalObj{texH: heap.Handle(args[0].IVal), sx: 1, sy: 1}
	id, err := store.Alloc(d)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) decalFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	store, err := m.requireHeapDecal(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DECAL.FREE expects decal handle")
	}
	_ = store.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func (m *Module) decalSetPos(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	store, err := m.requireHeapDecal(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("DECAL.SETPOS expects (decal, x, y, z)")
	}
	d, err := m.getDecal(store, args, 0, "DECAL.SETPOS")
	if err != nil {
		return value.Nil, err
	}
	var xf, yf, zf float32
	for i, a := range args[1:4] {
		if f, ok := a.ToFloat(); ok {
			switch i {
			case 0:
				xf = float32(f)
			case 1:
				yf = float32(f)
			case 2:
				zf = float32(f)
			}
		} else if iv, ok := a.ToInt(); ok {
			switch i {
			case 0:
				xf = float32(iv)
			case 1:
				yf = float32(iv)
			case 2:
				zf = float32(iv)
			}
		} else {
			return value.Nil, fmt.Errorf("DECAL.SETPOS: x, y, z must be numeric")
		}
	}
	d.pos = rl.Vector3{X: xf, Y: yf, Z: zf}
	return args[0], nil
}

func (m *Module) decalGetPos(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	store, err := m.requireHeapDecal(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DECAL.GETPOS expects decal handle")
	}
	d, err := m.getDecal(store, args, 0, "DECAL.GETPOS")
	if err != nil {
		return value.Nil, err
	}
	return mbmatrix.AllocVec3Value(store, d.pos.X, d.pos.Y, d.pos.Z)
}

func (m *Module) decalSetSize(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	store, err := m.requireHeapDecal(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("DECAL.SETSIZE expects (decal, width, height)")
	}
	d, err := m.getDecal(store, args, 0, "DECAL.SETSIZE")
	if err != nil {
		return value.Nil, err
	}
	var wf, hf float32
	if f, ok := args[1].ToFloat(); ok {
		wf = float32(f)
	} else if i, ok := args[1].ToInt(); ok {
		wf = float32(i)
	} else {
		return value.Nil, fmt.Errorf("DECAL.SETSIZE: width must be numeric")
	}
	if f, ok := args[2].ToFloat(); ok {
		hf = float32(f)
	} else if i, ok := args[2].ToInt(); ok {
		hf = float32(i)
	} else {
		return value.Nil, fmt.Errorf("DECAL.SETSIZE: height must be numeric")
	}
	if wf <= 0 || hf <= 0 {
		return value.Nil, fmt.Errorf("DECAL.SETSIZE: size must be positive")
	}
	d.sx, d.sy = wf, hf
	return args[0], nil
}

func (m *Module) decalGetSize(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	store, err := m.requireHeapDecal(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DECAL.GETSIZE expects decal handle")
	}
	d, err := m.getDecal(store, args, 0, "DECAL.GETSIZE")
	if err != nil {
		return value.Nil, err
	}
	return mbmatrix.AllocVec2Value(store, d.sx, d.sy)
}

func (m *Module) decalSetLifetime(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	store, err := m.requireHeapDecal(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DECAL.SETLIFETIME expects (decal, seconds)")
	}
	d, err := m.getDecal(store, args, 0, "DECAL.SETLIFETIME")
	if err != nil {
		return value.Nil, err
	}
	var sec float32
	if f, ok := args[1].ToFloat(); ok {
		sec = float32(f)
	} else if i, ok := args[1].ToInt(); ok {
		sec = float32(i)
	} else {
		return value.Nil, fmt.Errorf("DECAL.SETLIFETIME: seconds must be numeric")
	}
	if sec < 0 {
		sec = 0
	}
	d.life = sec
	d.age = 0
	return args[0], nil
}

func (m *Module) decalGetLifetime(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	store, err := m.requireHeapDecal(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DECAL.GETLIFETIME expects decal handle")
	}
	d, err := m.getDecal(store, args, 0, "DECAL.GETLIFETIME")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(d.life)), nil
}

// decalGetRot returns [0,0,0]; the decal object has no stored rotation (2D screen-space quad).
func (m *Module) decalGetRot(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	store, err := m.requireHeapDecal(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DECAL.GETROT expects decal handle")
	}
	if _, err := m.getDecal(store, args, 0, "DECAL.GETROT"); err != nil {
		return value.Nil, err
	}
	return mbmatrix.AllocVec3Value(store, 0, 0, 0)
}

func (m *Module) decalDraw(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	store, err := m.requireHeapDecal(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DECAL.DRAW expects decal handle")
	}
	d, err := m.getDecal(store, args, 0, "DECAL.DRAW")
	if err != nil {
		return value.Nil, err
	}
	cam, ok := mbmodel3d.ActiveCamera3D()
	if !ok {
		return value.Nil, fmt.Errorf("DECAL.DRAW: no active 3D camera (use CAMERA.BEGIN first)")
	}
	tex, err := mbdraw.TextureForBinding(store, d.texH)
	if err != nil {
		return value.Nil, fmt.Errorf("DECAL.DRAW: %w", err)
	}
	dt := rl.GetFrameTime()
	d.age += dt
	var alpha uint8 = 255
	if d.life > 0 {
		if d.age >= d.life {
			return value.Nil, nil
		}
		t := 1.0 - float64(d.age/d.life)
		if t < 0 {
			t = 0
		}
		if t > 1 {
			t = 1
		}
		alpha = uint8(255.0 * t)
	}
	tint := color.RGBA{R: 255, G: 255, B: 255, A: alpha}
	src := rl.NewRectangle(0, 0, float32(tex.Width), float32(tex.Height))
	rl.DrawBillboardRec(cam, tex, src, d.pos, rl.NewVector2(d.sx, d.sy), tint)
	return args[0], nil
}
