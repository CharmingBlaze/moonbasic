//go:build cgo || (windows && !cgo)

package mbsprite

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type atlasObj struct {
	tex     rl.Texture2D
	rects   map[string]rectI32
	release heap.ReleaseOnce
}

type rectI32 struct {
	x, y, w, h int32
}

func (o *atlasObj) TypeName() string { return "Atlas" }

func (o *atlasObj) TypeTag() uint16 { return heap.TagAtlas }

func (o *atlasObj) Free() {
	o.release.Do(func() { rl.UnloadTexture(o.tex) })
}

func parseAtlasJSON(data []byte) (map[string]rectI32, error) {
	var root map[string]json.RawMessage
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, err
	}
	rawFrames, ok := root["frames"]
	if !ok {
		return nil, fmt.Errorf("atlas JSON: missing \"frames\" object")
	}
	var framesObj map[string]json.RawMessage
	if err := json.Unmarshal(rawFrames, &framesObj); err != nil {
		return nil, fmt.Errorf("atlas JSON: frames: %w", err)
	}
	out := make(map[string]rectI32)
	for name, fr := range framesObj {
		var m map[string]interface{}
		if err := json.Unmarshal(fr, &m); err != nil {
			continue
		}
		x, y, w, h := extractFrameRect(m)
		if w > 0 && h > 0 {
			out[strings.TrimSuffix(name, filepath.Ext(name))] = rectI32{x: x, y: y, w: w, h: h}
			if _, has := out[name]; !has {
				out[name] = rectI32{x: x, y: y, w: w, h: h}
			}
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("atlas JSON: no frame rectangles parsed")
	}
	return out, nil
}

func extractFrameRect(m map[string]interface{}) (x, y, w, h int32) {
	if inner, ok := m["frame"].(map[string]interface{}); ok {
		m = inner
	}
	x = toI32(m["x"])
	y = toI32(m["y"])
	w = toI32(m["w"])
	h = toI32(m["h"])
	if w == 0 {
		w = toI32(m["width"])
	}
	if h == 0 {
		h = toI32(m["height"])
	}
	return
}

func toI32(v interface{}) int32 {
	switch t := v.(type) {
	case float64:
		return int32(t)
	case int:
		return int32(t)
	case int64:
		return int32(t)
	default:
		return 0
	}
}

func (m *Module) registerAtlas(reg runtime.Registrar) {
	reg.Register("ATLAS.LOAD", "sprite", m.atlasLoad)
	reg.Register("ATLAS.FREE", "sprite", runtime.AdaptLegacy(m.atlasFree))
	reg.Register("ATLAS.GETSPRITE", "sprite", m.atlasGetSprite)
}

func (m *Module) atlasLoad(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ATLAS.LOAD: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindString || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ATLAS.LOAD expects (imagePath, jsonPath)")
	}
	imgPath, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	jsonPath, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return value.Nil, fmt.Errorf("ATLAS.LOAD: read json: %w", err)
	}
	rects, err := parseAtlasJSON(data)
	if err != nil {
		return value.Nil, err
	}
	tex := rl.LoadTexture(imgPath)
	o := &atlasObj{tex: tex, rects: rects}
	id, err := m.h.Alloc(o)
	if err != nil {
		rl.UnloadTexture(tex)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) atlasFree(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("ATLAS.FREE expects atlas handle")
	}
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func (m *Module) atlasGetSprite(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ATLAS.GETSPRITE: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ATLAS.GETSPRITE expects (atlas, name)")
	}
	a, err := heap.Cast[*atlasObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	name, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	r, ok := a.rects[name]
	if !ok {
		return value.Nil, fmt.Errorf("ATLAS.GETSPRITE: unknown frame %q", name)
	}
	s := &spriteObj{
		tex:            a.tex,
		fromAtlas:      true,
		x:              0,
		y:              0,
		frameW:         r.w,
		frameH:         r.h,
		srcX:           r.x,
		srcY:           r.y,
		atlasRegionW:   r.w,
		atlasRegionH:   r.h,
		numFrames:      1,
		fps:            8,
		curFrame:       0,
		scaleX:         1,
		scaleY:         1,
		tr:             255,
		tg:             255,
		tb:             255,
		alpha:          spriteDefaultAlpha,
	}
	id, err := m.h.Alloc(s)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}
