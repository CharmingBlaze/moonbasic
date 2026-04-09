//go:build cgo || (windows && !cgo)

package texture

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"moonbasic/runtime"
	"moonbasic/runtime/mbimage"
	"moonbasic/runtime/mbjobs"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerTextureLoadCmds(m *Module, r runtime.Registrar) {
	r.Register("TEXTURE.LOAD", "texture", m.texLoad)
	r.Register("LOADTEXTURE", "texture", m.texLoad) // Blitz-style flat alias
	r.Register("LoadTexture", "texture", m.texLoad) // Modern Blitz alias (same registry id)
	r.Register("TEXTURE.LOADASYNC", "texture", m.texLoadAsync)
	r.Register("TEXTURE.ISLOADED", "texture", runtime.AdaptLegacy(m.texIsLoaded))
	r.Register("TEXTURE.FROMIMAGE", "texture", runtime.AdaptLegacy(m.texFromImage))
	r.Register("TEXTURE.FREE", "texture", runtime.AdaptLegacy(m.texFree))
	r.Register("FREETEXTURE", "texture", runtime.AdaptLegacy(m.texFree)) // Blitz-style flat alias
}

func (m *Module) texLoad(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TEXTURE.LOAD: heap not bound")
	}
	if len(args) < 1 || len(args) > 2 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("TEXTURE.LOAD expects (path$ [, flags#])")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	flags := int32(1)
	if len(args) == 2 {
		if fi, ok := args[1].ToInt(); ok {
			flags = int32(fi)
		}
	}
	t := rl.LoadTexture(path)
	obj := &TextureObject{Tex: t, loaded: true, SourcePath: path, Flags: flags, UScl: 1, VScl: 1}
	obj.setFinalizer()
	texApplyLoadFlags(&t, flags)
	obj.Tex = t
	id, err := m.h.Alloc(obj)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func texApplyLoadFlags(t *rl.Texture2D, flags int32) {
	// flags: documented presets — 1 = default trilinear + repeat
	if flags <= 0 {
		return
	}
	rl.SetTextureFilter(*t, rl.FilterTrilinear)
	rl.SetTextureWrap(*t, rl.WrapRepeat)
}

func (m *Module) texLoadAsync(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TEXTURE.LOADASYNC: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("TEXTURE.LOADASYNC expects 1 string path")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	obj := &TextureObject{isLoading: true}
	obj.setFinalizer()
	id, err := m.h.Alloc(obj)
	if err != nil {
		return value.Nil, err
	}

	mbjobs.EnqueueJob(func() {
		// Hand off to Main Thread for OpenGL calls (rl.LoadTexture requires context)
		enqueueOnMainThread(func() {
			t := rl.LoadTexture(path)
			obj.mu.Lock()
			obj.Tex = t
			obj.loaded = true
			obj.isLoading = false
			obj.mu.Unlock()
		})
	})

	return value.FromHandle(id), nil
}

func (m *Module) texIsLoaded(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TEXTURE.ISLOADED: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TEXTURE.ISLOADED expects 1 texture handle")
	}
	obj, ok := m.h.Get(heap.Handle(args[0].IVal))
	if !ok {
		return value.Nil, fmt.Errorf("TEXTURE.ISLOADED: invalid handle")
	}
	to, ok := obj.(*TextureObject)
	if !ok {
		return value.Nil, fmt.Errorf("TEXTURE.ISLOADED: handle is not a texture")
	}
	to.mu.RLock()
	defer to.mu.RUnlock()
	return value.FromBool(to.loaded), nil
}

func (m *Module) texFromImage(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TEXTURE.FROMIMAGE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TEXTURE.FROMIMAGE expects 1 image handle")
	}
	img, err := mbimage.RayImageForTexture(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("TEXTURE.FROMIMAGE: %w", err)
	}
	t := rl.LoadTextureFromImage(img)
	obj := &TextureObject{Tex: t}
	obj.setFinalizer()
	id, err := m.h.Alloc(obj)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) texFree(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TEXTURE.FREE expects handle")
	}
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}
