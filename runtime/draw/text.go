//go:build cgo || (windows && !cgo)

package mbdraw

import (
	"fmt"
	"image/color"
	"unicode/utf8"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	mbfont "moonbasic/runtime/font"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func stringFromRT(rt *runtime.Runtime, v value.Value) string {
	var pool []string
	var hg value.StringGetter
	if rt != nil {
		if rt.Prog != nil {
			pool = rt.Prog.StringTable
		}
		if rt.Heap != nil {
			hg = rt.Heap
		}
	}
	return value.StringAt(v, pool, hg)
}

func registerTextCmds(m *Module, r runtime.Registrar) {
	r.Register("DRAW.TEXT", "draw", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.text(rt, args)
	})
	r.Register("DRAW.TEXTEX", "draw", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.textEx(rt, args)
	})
	r.Register("DRAW.TEXTFONT", "draw", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.textEx(rt, args)
	})
	r.Register("DRAW.TEXTPRO", "draw", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.textPro(rt, args)
	})
	r.Register("DRAW.TEXTWIDTH", "draw", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.measureText(rt, args)
	})
	r.Register("DRAW.TEXTFONTWIDTH", "draw", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.measureTextFontWidth(rt, args)
	})
	r.Register("RENDER.DRAWFPS", "render", runtime.AdaptLegacy(m.drawFPS))
	r.Register("MEASURETEXT", "draw", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.measureText(rt, args)
	})
	r.Register("MEASURETEXTEX", "draw", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.measureTextEx(rt, args)
	})
	r.Register("GETTEXTCODEPOINTCOUNT", "draw", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.getCodepointCount(rt, args)
	})
}

func (m *Module) measureText(rt *runtime.Runtime, args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MEASURETEXT expects 2 arguments (text$, size)")
	}
	text := stringFromRT(rt, args[0])
	size, ok2 := argInt(args[1])
	if !ok2 {
		return value.Nil, fmt.Errorf("MEASURETEXT: arguments must be string and numeric")
	}
	width := rl.MeasureText(text, size)
	return value.FromInt(int64(width)), nil
}

func (m *Module) measureTextEx(rt *runtime.Runtime, args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("MEASURETEXTEX expects 4 arguments (font, text$, size, spacing)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("MEASURETEXTEX: font must be a handle")
	}
	font, err := mbfont.FontForHandle(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	text := stringFromRT(rt, args[1])
	size, ok2 := argFloat(args[2])
	spacing, ok3 := argFloat(args[3])
	if !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("MEASURETEXTEX: arguments must be correct types")
	}
	vec := rl.MeasureTextEx(font, text, size, spacing)
	arr, err := heap.NewArray([]int64{2})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, float64(vec.X))
	_ = arr.Set([]int64{1}, float64(vec.Y))
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) measureTextFontWidth(rt *runtime.Runtime, args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("DRAW.TEXTFONTWIDTH expects 4 arguments (font, text$, size, spacing)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAW.TEXTFONTWIDTH: font must be a handle")
	}
	font, err := mbfont.FontForHandle(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	text := stringFromRT(rt, args[1])
	size, ok2 := argFloat(args[2])
	spacing, ok3 := argFloat(args[3])
	if !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("DRAW.TEXTFONTWIDTH: arguments must be correct types")
	}
	vec := rl.MeasureTextEx(font, text, size, spacing)
	return value.FromFloat(float64(vec.X)), nil
}

func (m *Module) getCodepointCount(rt *runtime.Runtime, args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("GETTEXTCODEPOINTCOUNT expects 1 argument (text$)")
	}
	text := stringFromRT(rt, args[0])
	count := utf8.RuneCountInString(text)
	return value.FromInt(int64(count)), nil
}

func (m *Module) textEx(rt *runtime.Runtime, args []value.Value) (value.Value, error) {
	if len(args) != 10 {
		return value.Nil, fmt.Errorf("DRAW.TEXTEX expects 10 arguments (font, text$, x, y, size, spacing, r, g, b, a)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAW.TEXTEX: font must be a handle")
	}
	font, err := mbfont.FontForHandle(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	text := stringFromRT(rt, args[1])
	x, ok2 := argFloat(args[2])
	y, ok3 := argFloat(args[3])
	size, ok4 := argFloat(args[4])
	spacing, ok5 := argFloat(args[5])
	if !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("DRAW.TEXTEX: text, x, y, size, and spacing must be correct types")
	}
	r, ok6 := argInt(args[6])
	g, ok7 := argInt(args[7])
	b, ok8 := argInt(args[8])
	a, ok9 := argInt(args[9])
	if !ok6 || !ok7 || !ok8 || !ok9 {
		return value.Nil, fmt.Errorf("DRAW.TEXTEX: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawTextEx(font, text, rl.Vector2{X: x, Y: y}, size, spacing, col)
	return value.Nil, nil
}

func (m *Module) textPro(rt *runtime.Runtime, args []value.Value) (value.Value, error) {
	if len(args) != 13 {
		return value.Nil, fmt.Errorf("DRAW.TEXTPRO expects 13 arguments (font, text$, x, y, ox, oy, rot, size, spacing, r, g, b, a)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAW.TEXTPRO: font must be a handle")
	}
	font, err := mbfont.FontForHandle(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	text := stringFromRT(rt, args[1])
	x, ok2 := argFloat(args[2])
	y, ok3 := argFloat(args[3])
	ox, ok4 := argFloat(args[4])
	oy, ok5 := argFloat(args[5])
	rot, ok6 := argFloat(args[6])
	size, ok7 := argFloat(args[7])
	spacing, ok8 := argFloat(args[8])
	if !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 {
		return value.Nil, fmt.Errorf("DRAW.TEXTPRO: geometry arguments must be numeric")
	}
	r, ok9 := argInt(args[9])
	g, ok10 := argInt(args[10])
	b, ok11 := argInt(args[11])
	a, ok12 := argInt(args[12])
	if !ok9 || !ok10 || !ok11 || !ok12 {
		return value.Nil, fmt.Errorf("DRAW.TEXTPRO: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawTextPro(font, text, rl.Vector2{X: x, Y: y}, rl.Vector2{X: ox, Y: oy}, rot, size, spacing, col)
	return value.Nil, nil
}

func (m *Module) text(rt *runtime.Runtime, args []value.Value) (value.Value, error) {
	if len(args) != 8 {
		return value.Nil, fmt.Errorf("DRAW.TEXT expects 8 arguments (text$, x, y, size, r,g,b,a)")
	}
	text := stringFromRT(rt, args[0])
	x, ok2 := argInt(args[1])
	y, ok3 := argInt(args[2])
	size, ok4 := argInt(args[3])
	if !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAW.TEXT: text, x, y, and size must be correct types")
	}
	r, ok5 := argInt(args[4])
	g, ok6 := argInt(args[5])
	b, ok7 := argInt(args[6])
	a, ok8 := argInt(args[7])
	if !ok5 || !ok6 || !ok7 || !ok8 {
		return value.Nil, fmt.Errorf("DRAW.TEXT: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawText(text, x, y, size, col)
	return value.Nil, nil
}

func (m *Module) drawFPS(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("RENDER.DRAWFPS expects 2 arguments (x, y)")
	}
	x, ok1 := argInt(args[0])
	y, ok2 := argInt(args[1])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("RENDER.DRAWFPS: x and y must be numeric")
	}
	rl.DrawFPS(x, y)
	return value.Nil, nil
}
