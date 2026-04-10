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

func registerMaterialCmds(m *Module, reg runtime.Registrar) {
	matMakeDefault := func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("MATERIAL.MAKEDEFAULT expects no arguments")
		}
		mat := rl.LoadMaterialDefault()
		obj := &materialObj{mat: mat}
		obj.setFinalizer()
		id, err := m.h.Alloc(obj)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}
	reg.Register("MATERIAL.MAKEDEFAULT", "material", runtime.AdaptLegacy(matMakeDefault))
	reg.Register("MATERIAL.CREATE", "material", runtime.AdaptLegacy(matMakeDefault))

	reg.Register("MATERIAL.MAKEPBR", "material", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("MATERIAL.MAKEPBR expects no arguments")
		}
		mat := makePBRMaterial()
		obj := &materialObj{mat: mat, pbr: true}
		obj.setFinalizer()
		id, err := m.h.Alloc(obj)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}))

	reg.Register("MATERIAL.FREE", "material", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 || args[0].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("MATERIAL.FREE expects material handle")
		}
		if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
			return value.Nil, err
		}
		return value.Nil, nil
	}))

	reg.Register("MATERIAL.SETSHADER", "material", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("MATERIAL.SETSHADER expects material handle, shader handle")
		}
		mo, err := m.getMaterial(args, 0, "MATERIAL.SETSHADER")
		if err != nil {
			return value.Nil, err
		}
		so, err := m.getShader(args, 1, "MATERIAL.SETSHADER")
		if err != nil {
			return value.Nil, err
		}
		mo.mat.Shader = so.sh
		return value.Nil, nil
	}))

	reg.Register("MATERIAL.SETTEXTURE", "material", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("MATERIAL.SETTEXTURE expects material handle, slot, texture handle")
		}
		mo, err := m.getMaterial(args, 0, "MATERIAL.SETTEXTURE")
		if err != nil {
			return value.Nil, err
		}
		slot, ok := argInt(args[1])
		if !ok || slot < 0 || slot >= int32(rl.MaxMaterialMaps) {
			return value.Nil, fmt.Errorf("MATERIAL.SETTEXTURE: invalid map slot")
		}
		if args[2].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("MATERIAL.SETTEXTURE: texture handle required")
		}
		tex, err := texture.ForBinding(m.h, heap.Handle(args[2].IVal))
		if err != nil {
			return value.Nil, fmt.Errorf("MATERIAL.SETTEXTURE: %w", err)
		}
		rl.SetMaterialTexture(&mo.mat, slot, tex)
		return value.Nil, nil
	}))

	reg.Register("MATERIAL.SETCOLOR", "material", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("MATERIAL.SETCOLOR expects material handle, slot, r, g, b, a")
		}
		mo, err := m.getMaterial(args, 0, "MATERIAL.SETCOLOR")
		if err != nil {
			return value.Nil, err
		}
		slot, ok := argInt(args[1])
		if !ok || slot < 0 || slot >= int32(rl.MaxMaterialMaps) {
			return value.Nil, fmt.Errorf("MATERIAL.SETCOLOR: invalid map slot")
		}
		col, err := rgbaFromArgs(args[2], args[3], args[4], args[5])
		if err != nil {
			return value.Nil, fmt.Errorf("MATERIAL.SETCOLOR: %w", err)
		}
		mp := mo.mat.GetMap(slot)
		mp.Color = col
		return value.Nil, nil
	}))

	reg.Register("MATERIAL.SETFLOAT", "material", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("MATERIAL.SETFLOAT expects material handle, slot, value")
		}
		mo, err := m.getMaterial(args, 0, "MATERIAL.SETFLOAT")
		if err != nil {
			return value.Nil, err
		}
		slot, ok := argInt(args[1])
		if !ok || slot < 0 || slot >= int32(rl.MaxMaterialMaps) {
			return value.Nil, fmt.Errorf("MATERIAL.SETFLOAT: invalid map slot")
		}
		v, ok2 := argFloat(args[2])
		if !ok2 {
			return value.Nil, fmt.Errorf("MATERIAL.SETFLOAT: value must be numeric")
		}
		mp := mo.mat.GetMap(slot)
		mp.Value = v
		return value.Nil, nil
	}))

	reg.Register("MATERIAL.SETEFFECT", "material", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) < 2 {
			return value.Nil, fmt.Errorf("MATERIAL.SETEFFECT expects material handle, effectName$")
		}
		mo, err := m.getMaterial(args, 0, "MATERIAL.SETEFFECT")
		if err != nil {
			return value.Nil, err
		}
		eff, _ := m.h.GetString(int32(args[1].IVal))
		sh := globalShaderLib.GetShader(eff)
		if rl.IsShaderValid(sh) {
			mo.mat.Shader = sh
		} else {
			return value.Nil, fmt.Errorf("MATERIAL.SETEFFECT: unknown or invalid effect '%s'", eff)
		}
		return value.Nil, nil
	}))

	reg.Register("MATERIAL.SETEFFECTPARAM", "material", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("MATERIAL.SETEFFECTPARAM expects (matHandle, paramName$, value#)")
		}
		mo, err := m.getMaterial(args, 0, "MATERIAL.SETEFFECTPARAM")
		if err != nil {
			return value.Nil, err
		}
		pname, _ := m.h.GetString(int32(args[1].IVal))
		val, ok := args[2].ToFloat()
		if !ok {
			return value.Nil, fmt.Errorf("MATERIAL.SETEFFECTPARAM: value must be numeric")
		}
		if mo.params == nil {
			mo.params = make(map[string]float32)
		}
		mo.params[pname] = float32(val)
		return value.Nil, nil
	}))
}
