//go:build cgo || (windows && !cgo)

package texture

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"moonbasic/runtime"
	"moonbasic/runtime/mbimage"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerTextureLoadCmds(m *Module, r runtime.Registrar) {
	r.Register("TEXTURE.LOAD", "texture", m.texLoad)
	r.Register("LOADTEXTURE", "texture", m.texLoad) // Blitz-style flat alias
	r.Register("TEXTURE.FROMIMAGE", "texture", runtime.AdaptLegacy(m.texFromImage))
	r.Register("TEXTURE.FREE", "texture", runtime.AdaptLegacy(m.texFree))
	r.Register("FREETEXTURE", "texture", runtime.AdaptLegacy(m.texFree)) // Blitz-style flat alias
}

func (m *Module) texLoad(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TEXTURE.LOAD: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("TEXTURE.LOAD expects 1 string path")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	t := rl.LoadTexture(path)
	id, err := m.h.Alloc(&TextureObject{Tex: t})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
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
	id, err := m.h.Alloc(&TextureObject{Tex: t})
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
