//go:build cgo

package mbmodel3d

import (
	"fmt"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerModelInstDraw(m *Module, reg runtime.Registrar) {
	reg.Register("MODEL.MAKEINSTANCED", "model", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 || args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("MODEL.MAKEINSTANCED expects path$, instanceCount")
		}
		path, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		n64, ok := args[1].ToInt()
		if !ok {
			if f, okf := args[1].ToFloat(); okf {
				n64 = int64(f)
			} else {
				return value.Nil, fmt.Errorf("MODEL.MAKEINSTANCED: instanceCount must be numeric")
			}
		}
		if n64 < 1 || n64 > 200000 {
			return value.Nil, fmt.Errorf("MODEL.MAKEINSTANCED: instanceCount must be in range 1..200000")
		}
		n := int(n64)
		mod := rl.LoadModel(path)
		px := make([]float32, n)
		py := make([]float32, n)
		pz := make([]float32, n)
		sx := make([]float32, n)
		sy := make([]float32, n)
		sz := make([]float32, n)
		tf := make([]rl.Matrix, n)
		for i := range sx {
			sx[i], sy[i], sz[i] = 1, 1, 1
			tf[i] = rl.MatrixIdentity()
		}
		io := &instancedModelObj{
			model:      mod,
			loadedPath: path,
			meshIdx:    0,
			count:      n,
			px:         px,
			py:         py,
			pz:         pz,
			sx:         sx,
			sy:         sy,
			sz:         sz,
			transforms: tf,
		}
		id, err := m.h.Alloc(io)
		if err != nil {
			rl.UnloadModel(mod)
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	})

	reg.Register("MODEL.SETINSTANCEPOS", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("MODEL.SETINSTANCEPOS expects (instancedModel, index, x, y, z)")
		}
		o, err := m.getInstancedModel(args, 0, "MODEL.SETINSTANCEPOS")
		if err != nil {
			return value.Nil, err
		}
		idx, ok := argInt(args[1])
		if !ok || int(idx) < 0 || int(idx) >= o.count {
			return value.Nil, fmt.Errorf("MODEL.SETINSTANCEPOS: invalid index")
		}
		x, ok1 := argFloat(args[2])
		y, ok2 := argFloat(args[3])
		z, ok3 := argFloat(args[4])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("MODEL.SETINSTANCEPOS: x, y, z must be numeric")
		}
		i := int(idx)
		o.px[i], o.py[i], o.pz[i] = x, y, z
		return value.Nil, nil
	}))

	reg.Register("MODEL.SETINSTANCESCALE", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("MODEL.SETINSTANCESCALE expects (instancedModel, index, sx, sy, sz)")
		}
		o, err := m.getInstancedModel(args, 0, "MODEL.SETINSTANCESCALE")
		if err != nil {
			return value.Nil, err
		}
		idx, ok := argInt(args[1])
		if !ok || int(idx) < 0 || int(idx) >= o.count {
			return value.Nil, fmt.Errorf("MODEL.SETINSTANCESCALE: invalid index")
		}
		x, ok1 := argFloat(args[2])
		y, ok2 := argFloat(args[3])
		z, ok3 := argFloat(args[4])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("MODEL.SETINSTANCESCALE: scales must be numeric")
		}
		i := int(idx)
		o.sx[i], o.sy[i], o.sz[i] = x, y, z
		return value.Nil, nil
	}))

	reg.Register("MODEL.UPDATEINSTANCES", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MODEL.UPDATEINSTANCES expects instanced model handle")
		}
		o, err := m.getInstancedModel(args, 0, "MODEL.UPDATEINSTANCES")
		if err != nil {
			return value.Nil, err
		}
		for i := 0; i < o.count; i++ {
			tr := rl.MatrixTranslate(o.px[i], o.py[i], o.pz[i])
			sc := rl.MatrixScale(o.sx[i], o.sy[i], o.sz[i])
			o.transforms[i] = rl.MatrixMultiply(tr, sc)
		}
		return value.Nil, nil
	}))

	reg.Register("MODEL.DRAW", "model", runtime.AdaptLegacy(m.modelDraw))
}

func (m *Module) modelDraw(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("MODEL.DRAW expects model or instanced model handle")
	}
	h := heap.Handle(args[0].IVal)

	if o, err := heap.Cast[*modelObj](m.h, h); err == nil {
		if shadowDeferActive() && InCamera3D() {
			draw3dMu.Lock()
			deferredModels = append(deferredModels, h)
			draw3dMu.Unlock()
			return value.Nil, nil
		}
		rl.DrawModel(o.model, rl.Vector3{}, 1, rl.White)
		return value.Nil, nil
	}

	if lo, err := heap.Cast[*lodModelObj](m.h, h); err == nil {
		cam, in3D := ViewerPositionForRendering()
		if !in3D {
			cam = lo.worldPos()
		}
		li := lo.pickLOD(cam)
		if li < 0 {
			return value.Nil, nil
		}
		if shadowDeferActive() && InCamera3D() {
			draw3dMu.Lock()
			deferredModels = append(deferredModels, h)
			draw3dMu.Unlock()
			return value.Nil, nil
		}
		mod := &lo.models[li]
		saved := mod.Transform
		mod.Transform = lo.transform
		rl.DrawModel(*mod, rl.Vector3{}, 1, rl.White)
		mod.Transform = saved
		return value.Nil, nil
	}

	if io, err := heap.Cast[*instancedModelObj](m.h, h); err == nil {
		if io.meshIdx < 0 || io.meshIdx >= io.model.MeshCount {
			return value.Nil, fmt.Errorf("MODEL.DRAW: invalid mesh index on instanced model")
		}
		mi := io.meshIdx
		meshes := io.model.GetMeshes()
		mats := io.model.GetMaterials()
		mm := unsafe.Slice(io.model.MeshMaterial, io.model.MeshCount)
		mid := mm[mi]
		mesh := meshes[mi]
		mat := mats[mid]
		n := io.count
		if n <= 0 || len(io.transforms) < n {
			return value.Nil, fmt.Errorf("MODEL.DRAW: instanced model not ready")
		}
		if shadowDeferActive() && InCamera3D() {
			draw3dMu.Lock()
			pendingInstDraw = append(pendingInstDraw, instancedDrawRec{instH: h})
			draw3dMu.Unlock()
			return value.Nil, nil
		}
		rl.DrawMeshInstanced(mesh, mat, io.transforms[:n], n)
		return value.Nil, nil
	}

	return value.Nil, fmt.Errorf("MODEL.DRAW: handle is not a Model, LODModel, or InstancedModel")
}
