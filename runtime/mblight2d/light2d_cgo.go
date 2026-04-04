//go:build cgo

package mblight2d

import (
	"fmt"
	"image/color"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/window"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

var (
	lightMu sync.Mutex
	ambient = struct{ R, G, B, A uint8 }{30, 30, 50, 255}
	lights  []*light2dObj
)

func registerLight(l *light2dObj) {
	lightMu.Lock()
	defer lightMu.Unlock()
	lights = append(lights, l)
}

func unregisterLight(l *light2dObj) {
	lightMu.Lock()
	defer lightMu.Unlock()
	for i, x := range lights {
		if x == l {
			lights = append(lights[:i], lights[i+1:]...)
			return
		}
	}
}

func setAmbient(r, g, b, a uint8) {
	lightMu.Lock()
	defer lightMu.Unlock()
	ambient.R, ambient.G, ambient.B, ambient.A = r, g, b, a
}

// RegisterFrameHook registers the 2D lighting overlay (CGO only).
func RegisterFrameHook(w *window.Module) {
	if w == nil {
		return
	}
	w.AppendFrameDrawHook(drawLightOverlay)
}

type light2dObj struct {
	x, y            float32
	r, g, b, a      uint8
	radius          float32
	intensity       float32
	registered      bool
}

func (o *light2dObj) TypeName() string { return "Light2D" }

func (o *light2dObj) TypeTag() uint16 { return heap.TagLight2D }

func (o *light2dObj) Free() {
	if o.registered {
		unregisterLight(o)
		o.registered = false
	}
}

func (m *Module) requireHeap() error {
	if m.h == nil {
		return runtime.Errorf("LIGHT2D.*: heap not bound")
	}
	return nil
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

func (m *Module) Register(reg runtime.Registrar) {
	reg.Register("LIGHT2D.MAKE", "light2d", runtime.AdaptLegacy(m.ldMake))
	reg.Register("LIGHT2D.FREE", "light2d", runtime.AdaptLegacy(m.ldFree))
	reg.Register("LIGHT2D.SETPOS", "light2d", runtime.AdaptLegacy(m.ldSetPos))
	reg.Register("LIGHT2D.SETCOLOR", "light2d", runtime.AdaptLegacy(m.ldSetColor))
	reg.Register("LIGHT2D.SETRADIUS", "light2d", runtime.AdaptLegacy(m.ldSetRadius))
	reg.Register("LIGHT2D.SETINTENSITY", "light2d", runtime.AdaptLegacy(m.ldSetIntensity))
	reg.Register("RENDER.SET2DAMBIENT", "render", runtime.AdaptLegacy(m.set2DAmbient))
}

func (m *Module) set2DAmbient(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("RENDER.SET2DAmbIENT expects (r, g, b, a)")
	}
	r0, _ := argInt32(args[0])
	g0, _ := argInt32(args[1])
	b0, _ := argInt32(args[2])
	a0, _ := argInt32(args[3])
	setAmbient(clampU8(r0), clampU8(g0), clampU8(b0), clampU8(a0))
	return value.Nil, nil
}

func (m *Module) Shutdown() {}

func (m *Module) getLight(args []value.Value, ix int, op string) (*light2dObj, error) {
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: expected Light2D handle", op)
	}
	return heap.Cast[*light2dObj](m.h, heap.Handle(args[ix].IVal))
}

func (m *Module) ldMake(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("LIGHT2D.MAKE expects 0 arguments")
	}
	o := &light2dObj{r: 255, g: 255, b: 255, a: 255, radius: 100, intensity: 1}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	registerLight(o)
	o.registered = true
	return value.FromHandle(id), nil
}

func (m *Module) ldFree(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("LIGHT2D.FREE expects handle")
	}
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func (m *Module) ldSetPos(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("LIGHT2D.SETPOS expects (light, x, y)")
	}
	o, err := m.getLight(args, 0, "LIGHT2D.SETPOS")
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argF(args[1])
	y, ok2 := argF(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("LIGHT2D.SETPOS: x, y must be numeric")
	}
	o.x, o.y = x, y
	return value.Nil, nil
}

func (m *Module) ldSetColor(args []value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("LIGHT2D.SETCOLOR expects (light, r, g, b, a)")
	}
	o, err := m.getLight(args, 0, "LIGHT2D.SETCOLOR")
	if err != nil {
		return value.Nil, err
	}
	r0, _ := argInt32(args[1])
	g0, _ := argInt32(args[2])
	b0, _ := argInt32(args[3])
	a0, _ := argInt32(args[4])
	o.r, o.g, o.b, o.a = clampU8(r0), clampU8(g0), clampU8(b0), clampU8(a0)
	return value.Nil, nil
}

func (m *Module) ldSetRadius(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("LIGHT2D.SETRADIUS expects (light, radius)")
	}
	o, err := m.getLight(args, 0, "LIGHT2D.SETRADIUS")
	if err != nil {
		return value.Nil, err
	}
	rad, ok := argF(args[1])
	if !ok || rad < 0 {
		return value.Nil, fmt.Errorf("LIGHT2D.SETRADIUS: radius must be numeric")
	}
	o.radius = rad
	return value.Nil, nil
}

func (m *Module) ldSetIntensity(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("LIGHT2D.SETINTENSITY expects (light, intensity)")
	}
	o, err := m.getLight(args, 0, "LIGHT2D.SETINTENSITY")
	if err != nil {
		return value.Nil, err
	}
	in, ok := argF(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("LIGHT2D.SETINTENSITY: intensity must be numeric")
	}
	if in < 0 {
		in = 0
	}
	o.intensity = in
	return value.Nil, nil
}

func drawLightOverlay() {
	w := rl.GetRenderWidth()
	h := rl.GetRenderHeight()
	if w < 1 || h < 1 {
		return
	}
	lightMu.Lock()
	ar, ag, ab := ambient.R, ambient.G, ambient.B
	ls := append([]*light2dObj(nil), lights...)
	lightMu.Unlock()

	avg := (int(ar) + int(ag) + int(ab)) / 3
	if avg < 255 {
		overlayA := uint8(255 - avg)
		if overlayA > 0 {
			rl.DrawRectangle(0, 0, int32(w), int32(h), color.RGBA{R: 0, G: 0, B: 0, A: overlayA})
		}
	}

	if len(ls) == 0 {
		return
	}
	rl.SetBlendMode(rl.BlendAdditive)
	for _, l := range ls {
		if l.radius < 1 {
			continue
		}
		alpha := float32(l.a) * l.intensity * 0.004
		if alpha > 255 {
			alpha = 255
		}
		inner := color.RGBA{R: l.r, G: l.g, B: l.b, A: uint8(alpha)}
		outer := color.RGBA{R: l.r, G: l.g, B: l.b, A: 0}
		rl.DrawCircleGradient(int32(l.x), int32(l.y), l.radius, inner, outer)
	}
	rl.SetBlendMode(rl.BlendAlpha)
}
