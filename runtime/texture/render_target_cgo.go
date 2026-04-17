//go:build cgo || (windows && !cgo)

package texture

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerRenderTargetCmds(m *Module, r runtime.Registrar) {
	r.Register("RENDERTARGET.CREATE", "texture", runtime.AdaptLegacy(m.rtMake))
	r.Register("RENDERTARGET.MAKE", "texture", runtime.AdaptLegacy(m.rtMake))
	r.Register("RENDERTARGET.FREE", "texture", runtime.AdaptLegacy(m.rtFree))
	r.Register("RENDERTARGET.BEGIN", "texture", runtime.AdaptLegacy(m.rtBegin))
	r.Register("RENDERTARGET.END", "texture", runtime.AdaptLegacy(m.rtEnd))
	r.Register("RENDERTARGET.TEXTURE", "texture", runtime.AdaptLegacy(m.rtTexture))
}

func (m *Module) rtMake(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("RENDERTARGET.MAKE: heap not bound")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("RENDERTARGET.MAKE expects (width, height)")
	}
	w, ok0 := argI32(args[0])
	h, ok1 := argI32(args[1])
	if !ok0 || !ok1 || w <= 0 || h <= 0 {
		return value.Nil, fmt.Errorf("RENDERTARGET.MAKE: width and height must be positive integers")
	}
	rt := rl.LoadRenderTexture(w, h)
	o := &RenderTargetObject{RT: rt}
	o.setFinalizer()
	id, err := m.h.Alloc(o)
	if err != nil {
		rl.UnloadRenderTexture(rt)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) rtFree(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("RENDERTARGET.FREE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("RENDERTARGET.FREE expects render target handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) rtBegin(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("RENDERTARGET.BEGIN expects render target handle")
	}
	o, err := heap.Cast[*RenderTargetObject](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("RENDERTARGET.BEGIN: %w", err)
	}
	rl.BeginTextureMode(o.RT)
	return args[0], nil
}

func (m *Module) rtEnd(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("RENDERTARGET.END expects 0 arguments")
	}
	rl.EndTextureMode()
	return value.Nil, nil
}

// rtTexture returns a Texture handle referencing the color attachment; free it with TEXTURE.FREE (no GPU unload).
// The view is invalid after RENDERTARGET.FREE on the source render target.
func (m *Module) rtTexture(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("RENDERTARGET.TEXTURE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("RENDERTARGET.TEXTURE expects render target handle")
	}
	o, err := heap.Cast[*RenderTargetObject](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("RENDERTARGET.TEXTURE: %w", err)
	}
	tex := o.RT.Texture
	obj := &TextureObject{Tex: tex, Borrowed: true}
	obj.setFinalizer()
	id, err := m.h.Alloc(obj)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}
