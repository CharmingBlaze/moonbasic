//go:build cgo

package mbmodel3d

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

type meshObj struct {
	m rl.Mesh
}

func (o *meshObj) TypeName() string { return "Mesh" }

func (o *meshObj) TypeTag() uint16 { return heap.TagMesh }

func (o *meshObj) Free() {
	rl.UnloadMesh(&o.m)
}

type materialObj struct {
	mat   rl.Material
	moved bool // after MODEL.SETMATERIAL steal; skip UnloadMaterial on Free
}

func (o *materialObj) TypeName() string { return "Material" }

func (o *materialObj) TypeTag() uint16 { return heap.TagMaterial }

func (o *materialObj) Free() {
	if o.moved {
		return
	}
	rl.UnloadMaterial(o.mat)
}

type modelObj struct {
	model      rl.Model
	loadedPath string // set by MODEL.LOAD; used for MODEL.CLONE / MODEL.INSTANCE

	parent heap.Handle // scene graph (0 = none); host-driven, not used by Raylib

	// Stashed render hints (Raylib has no per-model pipeline; materials updated separately).
	wireframe bool
	cull      bool
	lighting  bool
	fog       bool
	blendMode int32 // -1 = unset; else rl.BlendMode value when applied by host
	depthBits int32 // bitmask: 1=no depth test, 2=no depth write (convention for future draw helpers)

	ambientR, ambientG, ambientB int32
}

func (o *modelObj) TypeName() string { return "Model" }

func (o *modelObj) TypeTag() uint16 { return heap.TagModel }

func (o *modelObj) Free() {
	rl.UnloadModel(o.model)
}

type shaderObj struct {
	sh rl.Shader
}

func (o *shaderObj) TypeName() string { return "Shader" }

func (o *shaderObj) TypeTag() uint16 { return heap.TagShader }

func (o *shaderObj) Free() {
	rl.UnloadShader(o.sh)
}
