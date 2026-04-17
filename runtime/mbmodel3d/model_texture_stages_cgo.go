//go:build cgo || (windows && !cgo)

package mbmodel3d

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/texture"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerModelTextureStages(m *Module, reg runtime.Registrar) {
	reg.Register("MODEL.SETTEXTURESTAGE", "model", runtime.AdaptLegacy(m.modelSetTextureStage))
	reg.Register("MODEL.SETSTAGEBLEND", "model", runtime.AdaptLegacy(m.modelSetStageBlend))
	reg.Register("MODEL.SETSTAGESCROLL", "model", runtime.AdaptLegacy(m.modelSetStageScroll))
	reg.Register("MODEL.SETSTAGESCALE", "model", runtime.AdaptLegacy(m.modelSetStageScale))
	reg.Register("MODEL.SETSTAGEROTATE", "model", runtime.AdaptLegacy(m.modelSetStageRotate))
	reg.Register("MODEL.SCROLLTEXTURE", "model", runtime.AdaptLegacy(m.modelScrollTexture))
	reg.Register("MODEL.SCALETEXTURE", "model", runtime.AdaptLegacy(m.modelScaleTexture))
	reg.Register("MODEL.ROTATETEXTURE", "model", runtime.AdaptLegacy(m.modelRotateTexture))
}

// modelSetTextureStage binds a texture to material 0 at the given map slot (same as MODEL.SETMATERIALTEXTURE with matIndex 0).
func (m *Module) modelSetTextureStage(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("MODEL.SETTEXTURESTAGE expects (model, stage, tex)")
	}
	modObj, err := m.getModel(args, 0, "MODEL.SETTEXTURESTAGE")
	if err != nil {
		return value.Nil, err
	}
	slot, ok := argInt(args[1])
	if !ok || slot < 0 || slot >= int32(rl.MaxMaterialMaps) {
		return value.Nil, fmt.Errorf("MODEL.SETTEXTURESTAGE: invalid stage/map slot")
	}
	if args[2].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("MODEL.SETTEXTURESTAGE: texture handle required")
	}
	tex, err := texture.ForBinding(m.h, heap.Handle(args[2].IVal))
	if err != nil {
		return value.Nil, err
	}
	mats := modObj.model.GetMaterials()
	if len(mats) < 1 {
		return value.Nil, fmt.Errorf("MODEL.SETTEXTURESTAGE: model has no materials")
	}
	rl.SetMaterialTexture(&mats[0], slot, tex)
	return args[0], nil
}

// modelSetStageBlend writes mode into MaterialMap.Value (custom shaders often read this per map).
func (m *Module) modelSetStageBlend(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("MODEL.SETSTAGEBLEND expects (model, stage, mode)")
	}
	modObj, err := m.getModel(args, 0, "MODEL.SETSTAGEBLEND")
	if err != nil {
		return value.Nil, err
	}
	stage, ok := argInt(args[1])
	if !ok || stage < 0 || stage >= int32(rl.MaxMaterialMaps) {
		return value.Nil, fmt.Errorf("MODEL.SETSTAGEBLEND: invalid stage")
	}
	mode, ok := argFloat(args[2])
	if !ok {
		return value.Nil, fmt.Errorf("MODEL.SETSTAGEBLEND: mode must be numeric")
	}
	mats := modObj.model.GetMaterials()
	if len(mats) < 1 {
		return value.Nil, fmt.Errorf("MODEL.SETSTAGEBLEND: model has no materials")
	}
	mats[0].GetMap(stage).Value = mode
	return args[0], nil
}

// modelSetStageScroll adds (u,v) to map.Value and material.Params[stage%4] (u into Value, v into Params slot).
func (m *Module) modelSetStageScroll(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("MODEL.SETSTAGESCROLL expects (model, stage, u, v)")
	}
	modObj, err := m.getModel(args, 0, "MODEL.SETSTAGESCROLL")
	if err != nil {
		return value.Nil, err
	}
	stage, ok := argInt(args[1])
	if !ok || stage < 0 || stage >= int32(rl.MaxMaterialMaps) {
		return value.Nil, fmt.Errorf("MODEL.SETSTAGESCROLL: invalid stage")
	}
	u, ok1 := argFloat(args[2])
	v, ok2 := argFloat(args[3])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("MODEL.SETSTAGESCROLL: u and v must be numeric")
	}
	mats := modObj.model.GetMaterials()
	if len(mats) < 1 {
		return value.Nil, fmt.Errorf("MODEL.SETSTAGESCROLL: model has no materials")
	}
	mat := &mats[0]
	mat.GetMap(stage).Value += u
	mat.Params[stage%4] += v
	return args[0], nil
}

// modelSetStageScale multiplies map.Value by u and Params[stage%4] scale factor by v.
func (m *Module) modelSetStageScale(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("MODEL.SETSTAGESCALE expects (model, stage, u, v)")
	}
	modObj, err := m.getModel(args, 0, "MODEL.SETSTAGESCALE")
	if err != nil {
		return value.Nil, err
	}
	stage, ok := argInt(args[1])
	if !ok || stage < 0 || stage >= int32(rl.MaxMaterialMaps) {
		return value.Nil, fmt.Errorf("MODEL.SETSTAGESCALE: invalid stage")
	}
	u, ok1 := argFloat(args[2])
	v, ok2 := argFloat(args[3])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("MODEL.SETSTAGESCALE: u and v must be numeric")
	}
	if u == 0 || v == 0 {
		return value.Nil, fmt.Errorf("MODEL.SETSTAGESCALE: u and v must be non-zero")
	}
	mats := modObj.model.GetMaterials()
	if len(mats) < 1 {
		return value.Nil, fmt.Errorf("MODEL.SETSTAGESCALE: model has no materials")
	}
	mat := &mats[0]
	mat.GetMap(stage).Value *= u
	mat.Params[stage%4] *= v
	return args[0], nil
}

func (m *Module) modelSetStageRotate(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("MODEL.SETSTAGEROTATE expects (model, stage, angle)")
	}
	modObj, err := m.getModel(args, 0, "MODEL.SETSTAGEROTATE")
	if err != nil {
		return value.Nil, err
	}
	stage, ok := argInt(args[1])
	if !ok || stage < 0 || stage >= int32(rl.MaxMaterialMaps) {
		return value.Nil, fmt.Errorf("MODEL.SETSTAGEROTATE: invalid stage")
	}
	ang, ok := argFloat(args[2])
	if !ok {
		return value.Nil, fmt.Errorf("MODEL.SETSTAGEROTATE: angle must be numeric")
	}
	mats := modObj.model.GetMaterials()
	if len(mats) < 1 {
		return value.Nil, fmt.Errorf("MODEL.SETSTAGEROTATE: model has no materials")
	}
	mats[0].GetMap(stage).Value += ang
	return args[0], nil
}

// modelScrollTexture adjusts material 0 Params[0], Params[1] (common UV scroll convention for custom shaders).
func (m *Module) modelScrollTexture(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("MODEL.SCROLLTEXTURE expects (model, u, v)")
	}
	modObj, err := m.getModel(args, 0, "MODEL.SCROLLTEXTURE")
	if err != nil {
		return value.Nil, err
	}
	u, ok1 := argFloat(args[1])
	v, ok2 := argFloat(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("MODEL.SCROLLTEXTURE: u and v must be numeric")
	}
	mats := modObj.model.GetMaterials()
	if len(mats) < 1 {
		return value.Nil, fmt.Errorf("MODEL.SCROLLTEXTURE: model has no materials")
	}
	mat := &mats[0]
	mat.Params[0] += u
	mat.Params[1] += v
	return args[0], nil
}

func (m *Module) modelScaleTexture(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("MODEL.SCALETEXTURE expects (model, u, v)")
	}
	modObj, err := m.getModel(args, 0, "MODEL.SCALETEXTURE")
	if err != nil {
		return value.Nil, err
	}
	u, ok1 := argFloat(args[1])
	v, ok2 := argFloat(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("MODEL.SCALETEXTURE: u and v must be numeric")
	}
	if u == 0 || v == 0 {
		return value.Nil, fmt.Errorf("MODEL.SCALETEXTURE: scale factors must be non-zero")
	}
	mats := modObj.model.GetMaterials()
	if len(mats) < 1 {
		return value.Nil, fmt.Errorf("MODEL.SCALETEXTURE: model has no materials")
	}
	mat := &mats[0]
	mat.Params[2] *= u
	mat.Params[3] *= v
	return args[0], nil
}

func (m *Module) modelRotateTexture(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MODEL.ROTATETEXTURE expects (model, angle)")
	}
	modObj, err := m.getModel(args, 0, "MODEL.ROTATETEXTURE")
	if err != nil {
		return value.Nil, err
	}
	ang, ok := argFloat(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("MODEL.ROTATETEXTURE: angle must be numeric")
	}
	mats := modObj.model.GetMaterials()
	if len(mats) < 1 {
		return value.Nil, fmt.Errorf("MODEL.ROTATETEXTURE: model has no materials")
	}
	mat := &mats[0]
	mat.Params[3] += ang
	return args[0], nil
}
