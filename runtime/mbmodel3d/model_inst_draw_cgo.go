//go:build cgo || (windows && !cgo)

package mbmodel3d

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func newInstancedFromLoadedModel(mod rl.Model, path string, n int) *instancedModelObj {
	px := make([]float32, n)
	py := make([]float32, n)
	pz := make([]float32, n)
	sx := make([]float32, n)
	sy := make([]float32, n)
	sz := make([]float32, n)
	rx := make([]float32, n)
	ry := make([]float32, n)
	rz := make([]float32, n)
	cr := make([]float32, n)
	cg := make([]float32, n)
	cb := make([]float32, n)
	ca := make([]float32, n)
	manual := make([]bool, n)
	tf := make([]rl.Matrix, n)
	for i := range sx {
		sx[i], sy[i], sz[i] = 1, 1, 1
		cr[i], cg[i], cb[i], ca[i] = 255, 255, 255, 255
		tf[i] = rl.MatrixIdentity()
	}
	return &instancedModelObj{
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
		rx:         rx,
		ry:         ry,
		rz:         rz,
		cr:         cr,
		cg:         cg,
		cb:         cb,
		ca:         ca,
		manual:     manual,
		transforms: tf,
	}
}

func registerModelInstDraw(m *Module, reg runtime.Registrar) {
	makeInstanced := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
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
		io := newInstancedFromLoadedModel(mod, path, n)
		id, err := m.h.Alloc(io)
		if err != nil {
			rl.UnloadModel(mod)
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}

	makeInstancedFromModel := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 || args[0].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("INSTANCE.MAKE expects (model, instanceCount)")
		}
		mo, err := m.getModel(args, 0, "INSTANCE.MAKE")
		if err != nil {
			return value.Nil, err
		}
		if mo.loadedPath == "" {
			return value.Nil, fmt.Errorf("INSTANCE.MAKE: model must come from MODEL.LOAD (file path); use INSTANCE.MAKEINSTANCED(path$, count) for assets")
		}
		n64, ok := args[1].ToInt()
		if !ok {
			if f, okf := args[1].ToFloat(); okf {
				n64 = int64(f)
			} else {
				return value.Nil, fmt.Errorf("INSTANCE.MAKE: instanceCount must be numeric")
			}
		}
		if n64 < 1 || n64 > 200000 {
			return value.Nil, fmt.Errorf("INSTANCE.MAKE: instanceCount must be in range 1..200000")
		}
		n := int(n64)
		mod := rl.LoadModel(mo.loadedPath)
		io := newInstancedFromLoadedModel(mod, mo.loadedPath, n)
		id, err := m.h.Alloc(io)
		if err != nil {
			rl.UnloadModel(mod)
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}

	reg.Register("MODEL.MAKEINSTANCED", "model", makeInstanced)
	reg.Register("INSTANCE.MAKEINSTANCED", "model", makeInstanced)
	reg.Register("INSTANCE.MAKE", "model", makeInstancedFromModel)

	setInstancePos := func(args []value.Value) (value.Value, error) {
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
		o.manual[i] = false
		return value.Nil, nil
	}
	reg.Register("MODEL.SETINSTANCEPOS", "model", runtime.AdaptLegacy(setInstancePos))
	reg.Register("INSTANCE.SETINSTANCEPOS", "model", runtime.AdaptLegacy(setInstancePos))
	reg.Register("INSTANCE.SETPOS", "model", runtime.AdaptLegacy(setInstancePos))

	setInstanceRot := func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("INSTANCE.SETROT expects (inst, index, rx, ry, rz) radians")
		}
		o, err := m.getInstancedModel(args, 0, "INSTANCE.SETROT")
		if err != nil {
			return value.Nil, err
		}
		idx, ok := argInt(args[1])
		if !ok || int(idx) < 0 || int(idx) >= o.count {
			return value.Nil, fmt.Errorf("INSTANCE.SETROT: invalid index")
		}
		rx, ok1 := argFloat(args[2])
		ry, ok2 := argFloat(args[3])
		rz, ok3 := argFloat(args[4])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("INSTANCE.SETROT: rotation must be numeric")
		}
		i := int(idx)
		o.rx[i], o.ry[i], o.rz[i] = rx, ry, rz
		o.manual[i] = false
		return value.Nil, nil
	}
	reg.Register("INSTANCE.SETROT", "model", runtime.AdaptLegacy(setInstanceRot))

	setInstanceScale := func(args []value.Value) (value.Value, error) {
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
		o.manual[i] = false
		return value.Nil, nil
	}
	reg.Register("MODEL.SETINSTANCESCALE", "model", runtime.AdaptLegacy(setInstanceScale))
	reg.Register("INSTANCE.SETINSTANCESCALE", "model", runtime.AdaptLegacy(setInstanceScale))
	reg.Register("INSTANCE.SETSCALE", "model", runtime.AdaptLegacy(setInstanceScale))

	setInstanceMatrix := func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("INSTANCE.SETMATRIX expects (inst, index, matrix)")
		}
		o, err := m.getInstancedModel(args, 0, "INSTANCE.SETMATRIX")
		if err != nil {
			return value.Nil, err
		}
		idx, ok := argInt(args[1])
		if !ok || int(idx) < 0 || int(idx) >= o.count {
			return value.Nil, fmt.Errorf("INSTANCE.SETMATRIX: invalid index")
		}
		mat, err := mbmatrix.MatrixRaylib(m.h, heap.Handle(args[2].IVal))
		if err != nil {
			return value.Nil, fmt.Errorf("INSTANCE.SETMATRIX: %w", err)
		}
		i := int(idx)
		o.transforms[i] = mat
		o.manual[i] = true
		return value.Nil, nil
	}
	reg.Register("INSTANCE.SETMATRIX", "model", runtime.AdaptLegacy(setInstanceMatrix))

	setInstanceColor := func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 6 {
			return value.Nil, fmt.Errorf("INSTANCE.SETCOLOR expects (inst, index, r, g, b, a)")
		}
		o, err := m.getInstancedModel(args, 0, "INSTANCE.SETCOLOR")
		if err != nil {
			return value.Nil, err
		}
		idx, ok := argInt(args[1])
		if !ok || int(idx) < 0 || int(idx) >= o.count {
			return value.Nil, fmt.Errorf("INSTANCE.SETCOLOR: invalid index")
		}
		r, ok1 := argFloat(args[2])
		g, ok2 := argFloat(args[3])
		b, ok3 := argFloat(args[4])
		a, ok4 := argFloat(args[5])
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return value.Nil, fmt.Errorf("INSTANCE.SETCOLOR: rgba must be numeric")
		}
		i := int(idx)
		o.cr[i], o.cg[i], o.cb[i], o.ca[i] = r, g, b, a
		return value.Nil, nil
	}
	reg.Register("INSTANCE.SETCOLOR", "model", runtime.AdaptLegacy(setInstanceColor))

	updateInstances := func(args []value.Value) (value.Value, error) {
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
			if o.manual[i] {
				continue
			}
			tr := rl.MatrixTranslate(o.px[i], o.py[i], o.pz[i])
			rot := rl.MatrixRotateXYZ(rl.Vector3{X: o.rx[i], Y: o.ry[i], Z: o.rz[i]})
			sc := rl.MatrixScale(o.sx[i], o.sy[i], o.sz[i])
			o.transforms[i] = rl.MatrixMultiply(rl.MatrixMultiply(tr, rot), sc)
		}
		return value.Nil, nil
	}
	reg.Register("MODEL.UPDATEINSTANCES", "model", runtime.AdaptLegacy(updateInstances))
	reg.Register("INSTANCE.UPDATEINSTANCES", "model", runtime.AdaptLegacy(updateInstances))
	reg.Register("INSTANCE.UPDATEBUFFER", "model", runtime.AdaptLegacy(updateInstances))

	setCullDistance := func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("INSTANCE.SETCULLDISTANCE expects (inst, distance)")
		}
		o, err := m.getInstancedModel(args, 0, "INSTANCE.SETCULLDISTANCE")
		if err != nil {
			return value.Nil, err
		}
		d, ok := argFloat(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("INSTANCE.SETCULLDISTANCE: distance must be numeric")
		}
		if d < 0 {
			d = 0
		}
		o.cullDistance = d
		return value.Nil, nil
	}
	reg.Register("INSTANCE.SETCULLDISTANCE", "model", runtime.AdaptLegacy(setCullDistance))

	instCount := func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("INSTANCE.COUNT expects instanced model handle")
		}
		o, err := m.getInstancedModel(args, 0, "INSTANCE.COUNT")
		if err != nil {
			return value.Nil, err
		}
		return value.FromInt(int64(o.count)), nil
	}
	reg.Register("INSTANCE.COUNT", "model", runtime.AdaptLegacy(instCount))

	instDrawLOD := func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("INSTANCE.DRAWLOD expects (inst, lodMesh, distance)")
		}
		io, err := m.getInstancedModel(args, 0, "INSTANCE.DRAWLOD")
		if err != nil {
			return value.Nil, err
		}
		lod, err := m.getMesh(args, 1, "INSTANCE.DRAWLOD")
		if err != nil {
			return value.Nil, err
		}
		dist, ok := argFloat(args[2])
		if !ok {
			return value.Nil, fmt.Errorf("INSTANCE.DRAWLOD: distance must be numeric")
		}
		h := heap.Handle(args[0].IVal)
		if shadowDeferActive() && InCamera3D() {
			draw3dMu.Lock()
			pendingInstDraw = append(pendingInstDraw, instancedDrawRec{instH: h, lodMeshH: heap.Handle(args[1].IVal), lodDist: dist})
			draw3dMu.Unlock()
			return value.Nil, nil
		}
		drawInstancedRaster(io, lod, dist, false)
		return value.Nil, nil
	}
	reg.Register("INSTANCE.DRAWLOD", "model", runtime.AdaptLegacy(instDrawLOD))

	reg.Register("INSTANCE.DRAW", "model", runtime.AdaptLegacy(m.modelDrawInstOnly))
	reg.Register("MODEL.DRAW", "model", runtime.AdaptLegacy(m.modelDraw))
}

func (m *Module) modelDrawInstOnly(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("INSTANCE.DRAW expects instanced model handle")
	}
	io, err := heap.Cast[*instancedModelObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("INSTANCE.DRAW: %w", err)
	}
	if shadowDeferActive() && InCamera3D() {
		draw3dMu.Lock()
		pendingInstDraw = append(pendingInstDraw, instancedDrawRec{instH: heap.Handle(args[0].IVal), lodMeshH: 0, lodDist: 0})
		draw3dMu.Unlock()
		return value.Nil, nil
	}
	drawInstancedRaster(io, nil, 0, false)
	return value.Nil, nil
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
		if o.hidden {
			return value.Nil, nil
		}
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
		n := io.count
		if n <= 0 || len(io.transforms) < n {
			return value.Nil, fmt.Errorf("MODEL.DRAW: instanced model not ready")
		}
		if shadowDeferActive() && InCamera3D() {
			draw3dMu.Lock()
			pendingInstDraw = append(pendingInstDraw, instancedDrawRec{instH: h, lodMeshH: 0, lodDist: 0})
			draw3dMu.Unlock()
			return value.Nil, nil
		}
		drawInstancedRaster(io, nil, 0, false)
		return value.Nil, nil
	}

	return value.Nil, fmt.Errorf("MODEL.DRAW: handle is not a Model, LODModel, or InstancedModel")
}
