//go:build cgo

package mbmodel3d

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

// meshObj owns mesh GPU data uploaded for Raylib; typically freed before parent model is unloaded.
type meshObj struct {
	m       rl.Mesh
	release heap.ReleaseOnce
}

func (o *meshObj) TypeName() string { return "Mesh" }

func (o *meshObj) TypeTag() uint16 { return heap.TagMesh }

func (o *meshObj) Free() {
	o.release.Do(func() { rl.UnloadMesh(&o.m) })
}

// materialObj owns Raylib materials/shaders; if moved==true, ownership transferred to a model — do not unload here.
type materialObj struct {
	mat     rl.Material
	moved   bool // after MODEL.SETMATERIAL steal; skip UnloadMaterial on Free
	pbr     bool // MATERIAL.MAKEPBR: extra uniforms + optional shadow sampling
	release heap.ReleaseOnce
}

func (o *materialObj) TypeName() string { return "Material" }

func (o *materialObj) TypeTag() uint16 { return heap.TagMaterial }

func (o *materialObj) Free() {
	if o.moved {
		return
	}
	o.release.Do(func() { rl.UnloadMaterial(o.mat) })
}

// modelObj owns a loaded Raylib model (meshes/materials); unload model before freeing borrowed material handles.
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

	release heap.ReleaseOnce
}

func (o *modelObj) TypeName() string { return "Model" }

func (o *modelObj) TypeTag() uint16 { return heap.TagModel }

func (o *modelObj) Free() {
	o.release.Do(func() { rl.UnloadModel(o.model) })
}

// instancedModelObj draws many copies of one mesh/material via DrawMeshInstanced.
type instancedModelObj struct {
	model      rl.Model
	loadedPath string
	meshIdx    int32
	count      int
	px, py, pz []float32
	sx, sy, sz []float32
	transforms []rl.Matrix

	release heap.ReleaseOnce
}

func (o *instancedModelObj) TypeName() string { return "InstancedModel" }

func (o *instancedModelObj) TypeTag() uint16 { return heap.TagInstancedModel }

func (o *instancedModelObj) Free() {
	o.release.Do(func() { rl.UnloadModel(o.model) })
}

// lodModelObj holds three detail levels; MODEL.SETLODDISTANCES configures distance bands.
type lodModelObj struct {
	models     [3]rl.Model
	band0      float32 // [0, band0) → LOD0 (highest detail)
	band1      float32 // [band0, band1) → LOD1
	band2      float32 // [band1, band2) → LOD2; dist >= band2 → culled
	configured bool
	transform  rl.Matrix

	release heap.ReleaseOnce
}

func (o *lodModelObj) TypeName() string { return "LODModel" }

func (o *lodModelObj) TypeTag() uint16 { return heap.TagLODModel }

func (o *lodModelObj) Free() {
	o.release.Do(func() {
		for i := range o.models {
			if o.models[i].MeshCount > 0 {
				rl.UnloadModel(o.models[i])
				o.models[i] = rl.Model{}
			}
		}
	})
}

func (o *lodModelObj) worldPos() rl.Vector3 {
	return rl.Vector3{X: o.transform.M12, Y: o.transform.M13, Z: o.transform.M14}
}

// pickLOD returns model index 0..2, or -1 if culled / invalid.
func (o *lodModelObj) pickLOD(cam rl.Vector3) int {
	if o.models[0].MeshCount == 0 {
		return -1
	}
	d := rl.Vector3Distance(o.worldPos(), cam)
	if !o.configured {
		return 0
	}
	if d < o.band0 {
		return 0
	}
	if d < o.band1 {
		if o.models[1].MeshCount > 0 {
			return 1
		}
		return 0
	}
	if d < o.band2 {
		if o.models[2].MeshCount > 0 {
			return 2
		}
		if o.models[1].MeshCount > 0 {
			return 1
		}
		return 0
	}
	return -1
}

type shaderObj struct {
	sh      rl.Shader
	release heap.ReleaseOnce
}

func (o *shaderObj) TypeName() string { return "Shader" }

func (o *shaderObj) TypeTag() uint16 { return heap.TagShader }

func (o *shaderObj) Free() {
	o.release.Do(func() { rl.UnloadShader(o.sh) })
}
