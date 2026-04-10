//go:build cgo || (windows && !cgo)

package mbmatrix

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerColor(reg runtime.Registrar) {
	reg.Register("COLOR.RGB", "color", runtime.AdaptLegacy(m.colorRGB))
	reg.Register("COLOR.RGBA", "color", runtime.AdaptLegacy(m.colorRGBA))
	reg.Register("COLOR.HEX", "color", m.colorHex)
	reg.Register("COLOR.HSV", "color", runtime.AdaptLegacy(m.colorHSV))
	reg.Register("COLOR.FROMHSV", "color", runtime.AdaptLegacy(m.colorHSV))
	reg.Register("COLOR.CLAMP", "color", runtime.AdaptLegacy(m.colorClamp))
	reg.Register("COLOR.FREE", "color", runtime.AdaptLegacy(m.colorFree))
	reg.Register("COLOR.R", "color", runtime.AdaptLegacy(m.colorR))
	reg.Register("COLOR.G", "color", runtime.AdaptLegacy(m.colorG))
	reg.Register("COLOR.B", "color", runtime.AdaptLegacy(m.colorB))
	reg.Register("COLOR.A", "color", runtime.AdaptLegacy(m.colorA))
	reg.Register("COLOR.LERP", "color", runtime.AdaptLegacy(m.colorLerp))
	reg.Register("COLOR.FADE", "color", runtime.AdaptLegacy(m.colorFade))
	reg.Register("COLOR.TOHSVX", "color", runtime.AdaptLegacy(m.colorToHSVX))
	reg.Register("COLOR.TOHSVY", "color", runtime.AdaptLegacy(m.colorToHSVY))
	reg.Register("COLOR.TOHSVZ", "color", runtime.AdaptLegacy(m.colorToHSVZ))
	reg.Register("COLOR.TOHSV", "color", runtime.AdaptLegacy(m.colorToHSVTuple))
	reg.Register("COLOR.TOHEX", "color", m.colorToHex)
	reg.Register("COLOR.INVERT", "color", runtime.AdaptLegacy(m.colorInvert))
	reg.Register("COLOR.CONTRAST", "color", runtime.AdaptLegacy(m.colorContrast))
	reg.Register("COLOR.BRIGHTNESS", "color", runtime.AdaptLegacy(m.colorBrightness))
}

func (m *Module) colorFromArgs(args []value.Value, idx int, op string) (color.RGBA, error) {
	if err := m.requireHeap(); err != nil {
		return color.RGBA{}, err
	}
	if idx >= len(args) || args[idx].Kind != value.KindHandle {
		return color.RGBA{}, fmt.Errorf("%s: argument %d must be color handle", op, idx+1)
	}
	o, err := heap.Cast[*colorObj](m.h, heap.Handle(args[idx].IVal))
	if err != nil {
		return color.RGBA{}, fmt.Errorf("%s: %w", op, err)
	}
	return o.c, nil
}

func (m *Module) allocColor(c color.RGBA) (value.Value, error) {
	id, err := m.h.Alloc(&colorObj{c: c})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func argU8(v value.Value) (uint8, bool) {
	if i, ok := v.ToInt(); ok {
		return clampU8(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return clampU8(int64(f)), true
	}
	return 0, false
}

func clampU8(x int64) uint8 {
	if x < 0 {
		return 0
	}
	if x > 255 {
		return 255
	}
	return uint8(x)
}

func (m *Module) colorRGB(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("COLOR.RGB expects 3 arguments (r, g, b)")
	}
	r, ok1 := argU8(args[0])
	g, ok2 := argU8(args[1])
	b, ok3 := argU8(args[2])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("COLOR.RGB: components must be numeric")
	}
	return m.allocColor(rl.NewColor(r, g, b, 255))
}

func (m *Module) colorRGBA(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("COLOR.RGBA expects 4 arguments (r, g, b, a)")
	}
	r, ok1 := argU8(args[0])
	g, ok2 := argU8(args[1])
	b, ok3 := argU8(args[2])
	a, ok4 := argU8(args[3])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("COLOR.RGBA: components must be numeric")
	}
	return m.allocColor(rl.NewColor(r, g, b, a))
}

func parseHexColorString(s string) (color.RGBA, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(strings.ToLower(s), "#")
	if s == "" {
		return color.RGBA{}, fmt.Errorf("empty hex string")
	}
	var n uint64
	var err error
	switch len(s) {
	case 6:
		n, err = strconv.ParseUint(s, 16, 24)
		if err != nil {
			return color.RGBA{}, err
		}
		return color.RGBA{R: uint8(n >> 16), G: uint8(n >> 8), B: uint8(n), A: 255}, nil
	case 8:
		n, err = strconv.ParseUint(s, 16, 32)
		if err != nil {
			return color.RGBA{}, err
		}
		return color.RGBA{R: uint8(n >> 24), G: uint8(n >> 16), B: uint8(n >> 8), A: uint8(n)}, nil
	default:
		return color.RGBA{}, fmt.Errorf("hex color must be #RRGGBB or #RRGGBBAA (%d digits)", len(s))
	}
}

func (m *Module) colorHex(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("COLOR.HEX expects string (#RRGGBB or #RRGGBBAA)")
	}
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	c, err := parseHexColorString(s)
	if err != nil {
		return value.Nil, fmt.Errorf("COLOR.HEX: %w", err)
	}
	return m.allocColor(c)
}

func (m *Module) colorHSV(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) == 2 {
		// COLOR.HSV(index, total) — evenly spaced hues on the wheel (1-based index).
		idx, ok1 := args[0].ToFloat()
		total, ok2 := args[1].ToFloat()
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("COLOR.HSV: index and total must be numeric")
		}
		if total < 1 {
			return value.Nil, fmt.Errorf("COLOR.HSV: total must be >= 1")
		}
		h := (idx - 1) / total * 360
		if h < 0 {
			h = 0
		}
		c := rl.ColorFromHSV(float32(h), 0.85, 1.0)
		return m.allocColor(c)
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("COLOR.HSV expects 2 arguments (index, total) or 3 arguments (h, s, v)")
	}
	h, ok1 := argF(args[0])
	s, ok2 := argF(args[1])
	v, ok3 := argF(args[2])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("COLOR.HSV: components must be numeric")
	}
	// Accept normalized hue 0..1 as a convenience form for procedural palettes.
	if h >= 0 && h <= 1 {
		h *= 360
	}
	c := rl.ColorFromHSV(float32(h), float32(s), float32(v))
	return m.allocColor(c)
}

func (m *Module) colorClamp(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("COLOR.CLAMP expects 3 arguments (r, g, b)")
	}
	r, ok1 := argU8(args[0])
	g, ok2 := argU8(args[1])
	b, ok3 := argU8(args[2])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("COLOR.CLAMP: components must be numeric")
	}
	arr, err := heap.NewArrayOfKind([]int64{3}, heap.ArrayKindFloat, 0)
	if err != nil {
		return value.Nil, err
	}
	arr.Floats[0] = float64(r)
	arr.Floats[1] = float64(g)
	arr.Floats[2] = float64(b)
	h, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func (m *Module) colorFree(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("COLOR.FREE expects color handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) colorR(args []value.Value) (value.Value, error) {
	c, err := m.colorFromArgs(args, 0, "COLOR.R")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("COLOR.R expects color handle")
	}
	return value.FromInt(int64(c.R)), nil
}

func (m *Module) colorG(args []value.Value) (value.Value, error) {
	c, err := m.colorFromArgs(args, 0, "COLOR.G")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("COLOR.G expects color handle")
	}
	return value.FromInt(int64(c.G)), nil
}

func (m *Module) colorB(args []value.Value) (value.Value, error) {
	c, err := m.colorFromArgs(args, 0, "COLOR.B")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("COLOR.B expects color handle")
	}
	return value.FromInt(int64(c.B)), nil
}

func (m *Module) colorA(args []value.Value) (value.Value, error) {
	c, err := m.colorFromArgs(args, 0, "COLOR.A")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("COLOR.A expects color handle")
	}
	return value.FromInt(int64(c.A)), nil
}

func (m *Module) colorLerp(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("COLOR.LERP expects (a, b, t)")
	}
	a, err := m.colorFromArgs(args, 0, "COLOR.LERP")
	if err != nil {
		return value.Nil, err
	}
	b, err := m.colorFromArgs(args, 1, "COLOR.LERP")
	if err != nil {
		return value.Nil, err
	}
	t, ok := argF(args[2])
	if !ok {
		return value.Nil, fmt.Errorf("COLOR.LERP: t must be numeric")
	}
	return m.allocColor(rl.ColorLerp(a, b, t))
}

func (m *Module) colorFade(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("COLOR.FADE expects (color, alpha) with alpha 0..1")
	}
	c, err := m.colorFromArgs(args, 0, "COLOR.FADE")
	if err != nil {
		return value.Nil, err
	}
	a, ok := argF(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("COLOR.FADE: alpha must be numeric")
	}
	return m.allocColor(rl.Fade(c, a))
}

func (m *Module) colorToHSVX(args []value.Value) (value.Value, error) {
	return m.colorToHSVComponent(args, 0, "COLOR.TOHSVX")
}

func (m *Module) colorToHSVY(args []value.Value) (value.Value, error) {
	return m.colorToHSVComponent(args, 1, "COLOR.TOHSVY")
}

func (m *Module) colorToHSVZ(args []value.Value) (value.Value, error) {
	return m.colorToHSVComponent(args, 2, "COLOR.TOHSVZ")
}

// colorToHSVTuple returns (h, s, v) as Raylib ColorToHSV — same components as TOHSVX/Y/Z in one tuple.
func (m *Module) colorToHSVTuple(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("COLOR.TOHSV expects color handle")
	}
	c, err := m.colorFromArgs(args, 0, "COLOR.TOHSV")
	if err != nil {
		return value.Nil, err
	}
	v := rl.ColorToHSV(c)
	arr, err := heap.NewArrayOfKind([]int64{3}, heap.ArrayKindFloat, 0)
	if err != nil {
		return value.Nil, err
	}
	arr.Floats[0] = float64(v.X)
	arr.Floats[1] = float64(v.Y)
	arr.Floats[2] = float64(v.Z)
	h, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func (m *Module) colorToHSVComponent(args []value.Value, axis int, op string) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("%s expects color handle", op)
	}
	c, err := m.colorFromArgs(args, 0, op)
	if err != nil {
		return value.Nil, err
	}
	v := rl.ColorToHSV(c)
	var out float32
	switch axis {
	case 0:
		out = v.X
	case 1:
		out = v.Y
	default:
		out = v.Z
	}
	return value.FromFloat(float64(out)), nil
}

func (m *Module) colorToHex(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("COLOR.TOHEX expects color handle")
	}
	c, err := m.colorFromArgs(args, 0, "COLOR.TOHEX")
	if err != nil {
		return value.Nil, err
	}
	s := fmt.Sprintf("#%02X%02X%02X%02X", c.R, c.G, c.B, c.A)
	return rt.RetString(s), nil
}

func (m *Module) colorInvert(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("COLOR.INVERT expects color handle")
	}
	c, err := m.colorFromArgs(args, 0, "COLOR.INVERT")
	if err != nil {
		return value.Nil, err
	}
	inv := color.RGBA{R: 255 - c.R, G: 255 - c.G, B: 255 - c.B, A: c.A}
	return m.allocColor(inv)
}

func (m *Module) colorContrast(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("COLOR.CONTRAST expects (color, contrast) with contrast roughly -1..1")
	}
	c, err := m.colorFromArgs(args, 0, "COLOR.CONTRAST")
	if err != nil {
		return value.Nil, err
	}
	v, ok := argF(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("COLOR.CONTRAST: contrast must be numeric")
	}
	return m.allocColor(rl.ColorContrast(c, v))
}

func (m *Module) colorBrightness(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("COLOR.BRIGHTNESS expects (color, brightness) with factor roughly -1..1")
	}
	c, err := m.colorFromArgs(args, 0, "COLOR.BRIGHTNESS")
	if err != nil {
		return value.Nil, err
	}
	v, ok := argF(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("COLOR.BRIGHTNESS: factor must be numeric")
	}
	return m.allocColor(rl.ColorBrightness(c, v))
}
