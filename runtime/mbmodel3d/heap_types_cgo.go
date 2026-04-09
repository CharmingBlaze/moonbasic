//go:build cgo || (windows && !cgo)

package mbmodel3d

import (
	"runtime"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

// meshObj owns mesh GPU data uploaded for Raylib; typically freed before parent model is unloaded.
// If consumedByModel is set, MODEL.MAKE(LoadModelFromMesh) shares the same Mesh GPU data with a Model —
// UnloadMesh is deferred to MODEL.FREE; MESH.FREE on the source handle skips GPU unload.
// If backingModel.MeshCount > 0, this mesh came from MESH.LOAD (first submesh of a loaded model);
// Free() calls UnloadModel(backingModel) and must not UnloadMesh separately.
// pin* slices are non-nil for MESH.MAKECUSTOM — keep vertex RAM alive for Raylib pointers.
type meshObj struct {
	m               rl.Mesh
	consumedByModel bool
	backingModel    rl.Model
	pinVerts        []float32
	pinNorms        []float32
	pinUVs          []float32
	pinIdx          []uint16
	release         heap.ReleaseOnce
}

func (o *meshObj) setFinalizer() {
	runtime.SetFinalizer(o, func(m *meshObj) {
		enqueueOnMainThread(func() { m.Free() })
	})
}

func (o *meshObj) TypeName() string { return "Mesh" }

func (o *meshObj) TypeTag() uint16 { return heap.TagMesh }

func (o *meshObj) Free() {
	o.release.Do(func() {
		if o.backingModel.MeshCount > 0 {
			rl.UnloadModel(o.backingModel)
			o.backingModel = rl.Model{}
			return
		}
		if o.consumedByModel {
			return
		}
		rl.UnloadMesh(&o.m)
	})
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

func (o *materialObj) setFinalizer() {
	runtime.SetFinalizer(o, func(m *materialObj) {
		enqueueOnMainThread(func() { m.Free() })
	})
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

	hidden bool // when true, MODEL.DRAW skips

	// Optional: loaded via MODEL.LOADANIMATIONS(model, path$); freed before model.
	anims       []rl.ModelAnimation
	animIdx     int
	animFrame   float32
	animPlaying bool
	animLoop    bool
	animSpeed   float32 // 1 = default

	// Asynchronous state
	mu        sync.RWMutex
	isLoading bool
	loaded    bool
	loadError string

	release heap.ReleaseOnce
}

func (o *modelObj) TypeName() string { return "Model" }

func (o *modelObj) TypeTag() uint16 { return heap.TagModel }

func (o *modelObj) Free() {
	o.release.Do(func() {
		if len(o.anims) > 0 {
			rl.UnloadModelAnimations(o.anims)
			o.anims = nil
		}
		rl.UnloadModel(o.model)
	})
}

func (o *modelObj) setFinalizer() {
	runtime.SetFinalizer(o, func(m *modelObj) {
		enqueueOnMainThread(func() { m.Free() })
	})
}

// instancedModelObj draws many copies of one mesh/material via DrawMeshInstanced.
type instancedModelObj struct {
	model      rl.Model
	loadedPath string
	meshIdx    int32
	count      int
	px, py, pz []float32
	sx, sy, sz []float32
	rx, ry, rz []float32 // radians (MatrixRotateXYZ); used when manual[i] is false
	cr, cg, cb, ca []float32
	manual         []bool // true: transforms[i] owned by INSTANCE.SETMATRIX until cleared by SetPos/Rot/Scale
	transforms     []rl.Matrix

	cullDistance float32 // if > 0, skip draw when camera is farther than this from instance centroid

	release heap.ReleaseOnce
}

func (o *instancedModelObj) TypeName() string { return "InstancedModel" }

func (o *instancedModelObj) TypeTag() uint16 { return heap.TagInstancedModel }

func (o *instancedModelObj) Free() {
	o.release.Do(func() { rl.UnloadModel(o.model) })
}

func (o *instancedModelObj) setFinalizer() {
	runtime.SetFinalizer(o, func(m *instancedModelObj) {
		enqueueOnMainThread(func() { m.Free() })
	})
}

func (o *instancedModelObj) anchorPos() rl.Vector3 {
	if o == nil || o.count <= 0 {
		return rl.Vector3{}
	}
	var sx, sy, sz float32
	for i := 0; i < o.count; i++ {
		sx += o.px[i]
		sy += o.py[i]
		sz += o.pz[i]
	}
	n := float32(o.count)
	return rl.Vector3{X: sx / n, Y: sy / n, Z: sz / n}
}

func (o *instancedModelObj) uniformInstanceColors() bool {
	if o == nil || o.count <= 0 {
		return true
	}
	r0, g0, b0, a0 := o.cr[0], o.cg[0], o.cb[0], o.ca[0]
	for i := 1; i < o.count; i++ {
		if o.cr[i] != r0 || o.cg[i] != g0 || o.cb[i] != b0 || o.ca[i] != a0 {
			return false
		}
	}
	return true
}

func (o *instancedModelObj) shouldCull() bool {
	if o == nil || o.cullDistance <= 0 {
		return false
	}
	cam, _ := ViewerPositionForRendering()
	return rl.Vector3Distance(o.anchorPos(), cam) > o.cullDistance
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

func (o *lodModelObj) setFinalizer() {
	runtime.SetFinalizer(o, func(m *lodModelObj) {
		enqueueOnMainThread(func() { m.Free() })
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

func (o *shaderObj) setFinalizer() {
	runtime.SetFinalizer(o, func(m *shaderObj) {
		enqueueOnMainThread(func() { m.Free() })
	})
}
