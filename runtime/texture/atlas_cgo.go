//go:build cgo || (windows && !cgo)

package texture

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerTextureAtlasCmds(m *Module, r runtime.Registrar) {
	reg := func(key string, fn func(*Module, []value.Value) (value.Value, error)) {
		r.Register(key, "texture", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return fn(m, a) }))
	}
	reg("TEXTURE.SETGRID", (*Module).texSetGrid)
	reg("TEXTURE.SETFRAME", (*Module).texSetFrame)
	reg("TEXTURE.LOADANIM", (*Module).texLoadAnim)
	reg("TEXTURE.PLAY", (*Module).texPlay)
	reg("TEXTURE.STOPANIM", (*Module).texStopAnim)
	reg("TEXTURE.TICKALL", (*Module).texTickAll)
	reg("TEXTURE.SETUVSCROLL", (*Module).texSetUVScroll)
	reg("TEXTURE.SETDISTORTION", (*Module).texSetDistortion)
}

func (m *Module) texObjTex(args []value.Value, op string) (*TextureObject, error) {
	if m.h == nil {
		return nil, fmt.Errorf("%s: heap not bound", op)
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s expects texture handle", op)
	}
	o, err := heap.Cast[*TextureObject](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (m *Module) texSetGrid(args []value.Value) (value.Value, error) {
	o, err := m.texObjTex(args, "TEXTURE.SETGRID")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("TEXTURE.SETGRID expects (texture, columns#, rows#)")
	}
	c, ok1 := args[1].ToInt()
	r, ok2 := args[2].ToInt()
	if !ok1 || !ok2 || c < 1 || r < 1 {
		return value.Nil, fmt.Errorf("TEXTURE.SETGRID: columns and rows must be >= 1")
	}
	o.mu.Lock()
	o.AtlasCols = int32(c)
	o.AtlasRows = int32(r)
	o.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) texSetFrame(args []value.Value) (value.Value, error) {
	o, err := m.texObjTex(args, "TEXTURE.SETFRAME")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("TEXTURE.SETFRAME expects (texture, frameIndex#)")
	}
	fi, ok := args[1].ToInt()
	if !ok {
		return value.Nil, fmt.Errorf("TEXTURE.SETFRAME: frameIndex must be numeric")
	}
	o.mu.Lock()
	o.FrameIndex = int32(fi)
	o.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) texLoadAnim(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TEXTURE.LOADANIM: heap not bound")
	}
	if len(args) != 3 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("TEXTURE.LOADANIM expects (path$, columns#, rows#)")
	}
	path, ok := m.h.GetString(int32(args[0].IVal))
	if !ok || path == "" {
		return value.Nil, fmt.Errorf("TEXTURE.LOADANIM: path string required")
	}
	cols, ok1 := args[1].ToInt()
	rows, ok2 := args[2].ToInt()
	if !ok1 || !ok2 || cols < 1 || rows < 1 {
		return value.Nil, fmt.Errorf("TEXTURE.LOADANIM: columns and rows must be >= 1")
	}
	t := rl.LoadTexture(path)
	if t.ID <= 0 {
		return value.Nil, fmt.Errorf("TEXTURE.LOADANIM: failed to load %q", path)
	}
	obj := &TextureObject{
		Tex:         t,
		loaded:      true,
		SourcePath:  path,
		Flags:       1,
		UScl:        1,
		VScl:        1,
		AtlasCols:   int32(cols),
		AtlasRows:   int32(rows),
		FrameIndex:  0,
		AnimLoop:    true,
	}
	obj.setFinalizer()
	texApplyLoadFlags(&obj.Tex, 1)
	id, err := m.h.Alloc(obj)
	if err != nil {
		rl.UnloadTexture(t)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) texPlay(args []value.Value) (value.Value, error) {
	o, err := m.texObjTex(args, "TEXTURE.PLAY")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("TEXTURE.PLAY expects (texture, fps#, loop#)")
	}
	fps, ok1 := args[1].ToFloat()
	if !ok1 {
		if fi, ok := args[1].ToInt(); ok {
			fps = float64(fi)
			ok1 = true
		}
	}
	if !ok1 || fps <= 0 {
		return value.Nil, fmt.Errorf("TEXTURE.PLAY: fps must be > 0")
	}
	loop := true
	if args[2].Kind == value.KindBool {
		loop = args[2].IVal != 0
	} else if bi, ok := args[2].ToInt(); ok {
		loop = bi != 0
	}
	o.mu.Lock()
	o.AnimFPS = float32(fps)
	o.AnimLoop = loop
	o.AnimPlaying = true
	o.animTime = 0
	o.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) texStopAnim(args []value.Value) (value.Value, error) {
	o, err := m.texObjTex(args, "TEXTURE.STOPANIM")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("TEXTURE.STOPANIM expects (texture)")
	}
	o.mu.Lock()
	o.AnimPlaying = false
	o.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) texTickAll(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TEXTURE.TICKALL: heap not bound")
	}
	var dt float32
	if len(args) == 0 {
		dt = rl.GetFrameTime()
	} else if len(args) == 1 {
		f, ok := args[0].ToFloat()
		if !ok {
			return value.Nil, fmt.Errorf("TEXTURE.TICKALL: optional dt# must be numeric")
		}
		dt = float32(f)
	} else {
		return value.Nil, fmt.Errorf("TEXTURE.TICKALL expects 0 or 1 arguments (dt#)")
	}
	m.h.RangeObjects(func(_ heap.Handle, obj heap.HeapObject) bool {
		if obj.TypeTag() != heap.TagTexture {
			return true
		}
		to, ok := obj.(*TextureObject)
		if !ok || to == nil {
			return true
		}
		to.Tick(dt)
		return true
	})
	return value.Nil, nil
}

func (m *Module) texSetUVScroll(args []value.Value) (value.Value, error) {
	o, err := m.texObjTex(args, "TEXTURE.SETUVSCROLL")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("TEXTURE.SETUVSCROLL expects (texture, speedU#, speedV#) — pixels/sec in source space")
	}
	su, ok1 := args[1].ToFloat()
	sv, ok2 := args[2].ToFloat()
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("TEXTURE.SETUVSCROLL: speeds must be numeric")
	}
	o.mu.Lock()
	o.ScrollSpeedU = float32(su)
	o.ScrollSpeedV = float32(sv)
	o.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) texSetDistortion(args []value.Value) (value.Value, error) {
	o, err := m.texObjTex(args, "TEXTURE.SETDISTORTION")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("TEXTURE.SETDISTORTION expects (texture, amount#)")
	}
	a, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("TEXTURE.SETDISTORTION: amount must be numeric")
	}
	o.mu.Lock()
	o.Distortion = float32(a)
	o.mu.Unlock()
	return value.Nil, nil
}
