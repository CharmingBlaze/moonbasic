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

func registerModelMaterial(m *Module, reg runtime.Registrar) {
	// SETMATERIAL transfers GPU material ownership from the heap material into the model slot
	// (source handle must not be used afterward; MATERIAL.FREE on it is a no-op).
	reg.Register("MODEL.SETMATERIAL", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("MODEL.SETMATERIAL expects model handle, materialIndex, materialHandle")
		}
		modObj, err := m.getModel(args, 0, "MODEL.SETMATERIAL")
		if err != nil {
			return value.Nil, err
		}
		matIdx, ok := argInt(args[1])
		if !ok || matIdx < 0 {
			return value.Nil, fmt.Errorf("MODEL.SETMATERIAL: invalid material index")
		}
		matObj, err := m.getMaterial(args, 2, "MODEL.SETMATERIAL")
		if err != nil {
			return value.Nil, err
		}
		mats := modObj.model.GetMaterials()
		if int(matIdx) >= len(mats) {
			return value.Nil, fmt.Errorf("MODEL.SETMATERIAL: material index out of range")
		}
		rl.UnloadMaterial(mats[matIdx])
		mats[matIdx] = matObj.mat
		matObj.mat = rl.Material{}
		matObj.moved = true
		return value.Nil, nil
	}))

	reg.Register("MODEL.SETMATERIALTEXTURE", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("MODEL.SETMATERIALTEXTURE expects model handle, materialIndex, slot, texture handle")
		}
		modObj, err := m.getModel(args, 0, "MODEL.SETMATERIALTEXTURE")
		if err != nil {
			return value.Nil, err
		}
		matIdx, ok := argInt(args[1])
		if !ok || matIdx < 0 {
			return value.Nil, fmt.Errorf("MODEL.SETMATERIALTEXTURE: invalid material index")
		}
		slot, ok2 := argInt(args[2])
		if !ok2 || slot < 0 || slot >= int32(rl.MaxMaterialMaps) {
			return value.Nil, fmt.Errorf("MODEL.SETMATERIALTEXTURE: invalid map slot")
		}
		if args[3].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("MODEL.SETMATERIALTEXTURE: texture handle required")
		}
		tex, err := texture.ForBinding(m.h, heap.Handle(args[3].IVal))
		if err != nil {
			return value.Nil, fmt.Errorf("MODEL.SETMATERIALTEXTURE: %w", err)
		}
		mats := modObj.model.GetMaterials()
		if int(matIdx) >= len(mats) {
			return value.Nil, fmt.Errorf("MODEL.SETMATERIALTEXTURE: material index out of range")
		}
		rl.SetMaterialTexture(&mats[matIdx], slot, tex)
		return value.Nil, nil
	}))

	reg.Register("MODEL.SETMATERIALSHADER", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("MODEL.SETMATERIALSHADER expects model handle, materialIndex, shader handle")
		}
		modObj, err := m.getModel(args, 0, "MODEL.SETMATERIALSHADER")
		if err != nil {
			return value.Nil, err
		}
		matIdx, ok := argInt(args[1])
		if !ok || matIdx < 0 {
			return value.Nil, fmt.Errorf("MODEL.SETMATERIALSHADER: invalid material index")
		}
		shObj, err := m.getShader(args, 2, "MODEL.SETMATERIALSHADER")
		if err != nil {
			return value.Nil, err
		}
		mats := modObj.model.GetMaterials()
		if int(matIdx) >= len(mats) {
			return value.Nil, fmt.Errorf("MODEL.SETMATERIALSHADER: material index out of range")
		}
		mats[matIdx].Shader = shObj.sh
		return value.Nil, nil
	}))

	reg.Register("MODEL.SETMODELMESHMATERIAL", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("MODEL.SETMODELMESHMATERIAL expects model handle, meshId, materialId")
		}
		modObj, err := m.getModel(args, 0, "MODEL.SETMODELMESHMATERIAL")
		if err != nil {
			return value.Nil, err
		}
		meshID, ok1 := argInt(args[1])
		matID, ok2 := argInt(args[2])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("MODEL.SETMODELMESHMATERIAL: meshId and materialId must be numeric")
		}
		rl.SetModelMeshMaterial(&modObj.model, meshID, matID)
		return value.Nil, nil
	}))
}
