//go:build cgo || (windows && !cgo)

package mbfont

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type fontObj struct {
	f       rl.Font
	release heap.ReleaseOnce
}

func (o *fontObj) TypeName() string { return "Font" }

func (o *fontObj) TypeTag() uint16 { return heap.TagFont }

func (o *fontObj) Free() {
	o.release.Do(func() { rl.UnloadFont(o.f) })
}

func argHandle(v value.Value) (heap.Handle, bool) {
	if v.Kind != value.KindHandle {
		return 0, false
	}
	return heap.Handle(v.IVal), true
}

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	r.Register("FONT.LOAD", "font", m.fontLoad)
	r.Register("FONT.LOADBDF", "font", m.fontLoadBDF)
	r.Register("FONT.FREE", "font", runtime.AdaptLegacy(m.fontFree))
	r.Register("FONT.DRAWDEFAULT", "font", runtime.AdaptLegacy(m.fontDrawDefault))
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

func (m *Module) fontLoad(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("FONT.LOAD: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("FONT.LOAD expects 1 string path")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	f := rl.LoadFont(path)
	o := &fontObj{f: f}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) fontLoadBDF(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("FONT.LOADBDF: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("FONT.LOADBDF expects (path$, size)")
	}
	var sz int64
	if i, ok := args[1].ToInt(); ok {
		sz = i
	} else if f, ok := args[1].ToFloat(); ok {
		sz = int64(f)
	} else {
		return value.Nil, fmt.Errorf("FONT.LOADBDF: size must be numeric")
	}
	if sz <= 0 || sz > 1024 {
		return value.Nil, fmt.Errorf("FONT.LOADBDF: size must be in 1..1024")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	f := rl.LoadFontEx(path, int32(sz), nil)
	o := &fontObj{f: f}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) fontFree(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("FONT.FREE expects 1 argument (handle)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("FONT.FREE: invalid handle")
	}
	if h == 0 {
		return value.Nil, nil
	}
	m.h.Free(h)
	return value.Nil, nil
}

// fontDrawDefault returns handle 0 — use with Raylib GetFontDefault() when drawing.
func (m *Module) fontDrawDefault(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("FONT.DRAWDEFAULT expects 0 arguments")
	}
	return value.FromHandle(0), nil
}

// FontForHandle returns the raylib Font for a given heap handle.
func FontForHandle(store *heap.Store, h heap.Handle) (rl.Font, error) {
	if h == 0 {
		return rl.GetFontDefault(), nil
	}
	o, err := heap.Cast[*fontObj](store, h)
	if err != nil {
		return rl.Font{}, err
	}
	return o.f, nil
}
