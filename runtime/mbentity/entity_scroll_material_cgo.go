//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/texture"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerEntityMaterialScrollAPI(m *Module, r runtime.Registrar) {
	r.Register("ENTITY.SCROLLMATERIAL", "entity", runtime.AdaptLegacy(m.entScrollMaterial))
	r.Register("ENTITY.SETDETAILTEXTURE", "entity", runtime.AdaptLegacy(m.entSetDetailTexture))
}

// entScrollMaterial adds (du,dv) to material 0 scroll params (same convention as MODEL.SCROLLTEXTURE).
func (m *Module) entScrollMaterial(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ENTITY.SCROLLMATERIAL: heap not bound")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.SCROLLMATERIAL expects (entity#, du#, dv#)")
	}
	eid, ok := m.entID(args[0])
	if !ok || eid < 1 {
		return value.Nil, fmt.Errorf("ENTITY.SCROLLMATERIAL: invalid entity")
	}
	e := m.store().ents[eid]
	if e == nil || !e.hasRLModel {
		return value.Nil, fmt.Errorf("ENTITY.SCROLLMATERIAL: entity has no 3D model")
	}
	du, ok1 := argF32(args[1])
	dv, ok2 := argF32(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("ENTITY.SCROLLMATERIAL: du, dv must be numeric")
	}
	mats := e.rlModel.GetMaterials()
	if len(mats) < 1 {
		return value.Nil, fmt.Errorf("ENTITY.SCROLLMATERIAL: no materials")
	}
	mat := &mats[0]
	mat.Params[0] += du
	mat.Params[1] += dv
	return value.Nil, nil
}

// entSetDetailTexture binds a secondary map (noise/detail) to MATERIAL_MAP_NORMAL for material 0.
func (m *Module) entSetDetailTexture(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ENTITY.SETDETAILTEXTURE: heap not bound")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SETDETAILTEXTURE expects (entity#, textureHandle)")
	}
	eid, ok := m.entID(args[0])
	if !ok || eid < 1 {
		return value.Nil, fmt.Errorf("ENTITY.SETDETAILTEXTURE: invalid entity")
	}
	e := m.store().ents[eid]
	if e == nil || !e.hasRLModel {
		return value.Nil, fmt.Errorf("ENTITY.SETDETAILTEXTURE: entity has no 3D model")
	}
	if args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("ENTITY.SETDETAILTEXTURE: texture handle required")
	}
	tex, err := texture.ForBinding(m.h, heap.Handle(args[1].IVal))
	if err != nil {
		return value.Nil, err
	}
	mats := e.rlModel.GetMaterials()
	if len(mats) < 1 {
		return value.Nil, fmt.Errorf("ENTITY.SETDETAILTEXTURE: no materials")
	}
	rl.SetMaterialTexture(&mats[0], rl.MapNormal, tex)
	return value.Nil, nil
}
