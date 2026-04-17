//go:build cgo || (windows && !cgo)

package mbmodel3d

import (
	"fmt"
	"math"
	"strings"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerMeshComplete(m *Module, reg runtime.Registrar) {
	reg.Register("MESH.LOAD", "mesh", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 || args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("MESH.LOAD expects path")
		}
		path, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		path = strings.TrimSpace(path)
		if path == "" {
			return value.Nil, fmt.Errorf("MESH.LOAD: path required")
		}
		mod := rl.LoadModel(path)
		meshes := mod.GetMeshes()
		if len(meshes) == 0 {
			rl.UnloadModel(mod)
			return value.Nil, fmt.Errorf("MESH.LOAD: no meshes in %q", path)
		}
		mesh := meshes[0]
		obj := &meshObj{m: mesh, backingModel: mod}
		obj.setFinalizer()
		id, err := m.h.Alloc(obj)
		if err != nil {
			rl.UnloadModel(mod)
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	})

	reg.Register("MESH.EXPORT", "mesh", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
			return value.Nil, fmt.Errorf("MESH.EXPORT expects (mesh, path)")
		}
		o, err := m.getMesh(args, 0, "MESH.EXPORT")
		if err != nil {
			return value.Nil, err
		}
		path, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Nil, err
		}
		path = strings.TrimSpace(path)
		if path == "" {
			return value.Nil, fmt.Errorf("MESH.EXPORT: path required")
		}
		rl.ExportMesh(o.m, path)
		return value.Nil, nil
	})

	reg.Register("MESH.GETBOUNDS", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MESH.GETBOUNDS expects mesh handle")
		}
		o, err := m.getMesh(args, 0, "MESH.GETBOUNDS")
		if err != nil {
			return value.Nil, err
		}
		b := rl.GetMeshBoundingBox(o.m)
		a, err := heap.NewArray([]int64{6})
		if err != nil {
			return value.Nil, err
		}
		_ = a.Set([]int64{0}, float64(b.Min.X))
		_ = a.Set([]int64{1}, float64(b.Min.Y))
		_ = a.Set([]int64{2}, float64(b.Min.Z))
		_ = a.Set([]int64{3}, float64(b.Max.X))
		_ = a.Set([]int64{4}, float64(b.Max.Y))
		_ = a.Set([]int64{5}, float64(b.Max.Z))
		id, err := m.h.Alloc(a)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}))

	reg.Register("MESH.GENERATEBOUNDS", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MESH.GENERATEBOUNDS expects mesh handle")
		}
		if _, err := m.getMesh(args, 0, "MESH.GENERATEBOUNDS"); err != nil {
			return value.Nil, err
		}
		// Raylib has no separate "commit bounds" — geometry already defines bounds; GetMeshBoundingBox recomputes.
		return value.Nil, nil
	}))

	reg.Register("MESH.GENERATENORMALS", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MESH.GENERATENORMALS expects mesh handle")
		}
		if _, err := m.getMesh(args, 0, "MESH.GENERATENORMALS"); err != nil {
			return value.Nil, err
		}
		return value.Nil, fmt.Errorf("MESH.GENERATENORMALS: not available — use procedural MESH.MAKE* meshes or precompute normals offline")
	}))

	reg.Register("MESH.UPDATEVERTICES", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("MESH.UPDATEVERTICES expects (mesh, verts_array)")
		}
		o, err := m.getMesh(args, 0, "MESH.UPDATEVERTICES")
		if err != nil {
			return value.Nil, err
		}
		if args[1].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("MESH.UPDATEVERTICES: verts must be float array handle")
		}
		ah := heap.Handle(args[1].IVal)
		n := m.h.ArrayFlatLen(ah)
		if n < 0 {
			return value.Nil, fmt.Errorf("MESH.UPDATEVERTICES: invalid array")
		}
		per := 8 // x,y,z,nx,ny,nz,u,v
		if n%per != 0 {
			return value.Nil, fmt.Errorf("MESH.UPDATEVERTICES: vertex float count must be multiple of 8")
		}
		vc := n / per
		if int32(vc) != o.m.VertexCount {
			return value.Nil, fmt.Errorf("MESH.UPDATEVERTICES: vertex count mismatch (mesh has %d vertices)", o.m.VertexCount)
		}
		for i := 0; i < vc; i++ {
			base := int64(i * per)
			vv := make([]value.Value, 10)
			vv[0] = args[0]
			for j := 0; j < 9; j++ {
				f, ok := m.h.ArrayGetFloat(ah, base+int64(j))
				if !ok {
					return value.Nil, fmt.Errorf("MESH.UPDATEVERTICES: bad array data")
				}
				vv[j+1] = value.FromFloat(f)
			}
			if _, err := meshUpdateVertexOne(m, vv); err != nil {
				return value.Nil, err
			}
		}
		return value.Nil, nil
	}))

	reg.Register("MESH.DRAWAT", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("MESH.DRAWAT expects (mesh, material, x, y, z)")
		}
		mo, err := m.getMesh(args, 0, "MESH.DRAWAT")
		if err != nil {
			return value.Nil, err
		}
		mato, err := m.getMaterial(args, 1, "MESH.DRAWAT")
		if err != nil {
			return value.Nil, err
		}
		x, ok1 := argFloat(args[2])
		y, ok2 := argFloat(args[3])
		z, ok3 := argFloat(args[4])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("MESH.DRAWAT: position must be numeric")
		}
		t := rl.MatrixTranslate(x, y, z)
		if shadowDeferActive() && InCamera3D() {
			draw3dMu.Lock()
			deferredMeshes = append(deferredMeshes, deferredMeshRec{
				meshH: heap.Handle(args[0].IVal),
				matH:  heap.Handle(args[1].IVal),
				mtx:   t,
			})
			draw3dMu.Unlock()
			return value.Nil, nil
		}
		son := shadowDeferActive()
		if mato.pbr && son {
			bindPBRDrawState(mato, true)
		}
		rl.DrawMesh(mo.m, mato.mat, t)
		if mato.pbr && son {
			clearShadowMapSlot(&mato.mat)
		}
		return value.Nil, nil
	}))

	reg.Register("MESH.DRAWINSTANCED", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("MESH.DRAWINSTANCED expects (mesh, material, transforms_array, count)")
		}
		mo, err := m.getMesh(args, 0, "MESH.DRAWINSTANCED")
		if err != nil {
			return value.Nil, err
		}
		mato, err := m.getMaterial(args, 1, "MESH.DRAWINSTANCED")
		if err != nil {
			return value.Nil, err
		}
		if args[2].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("MESH.DRAWINSTANCED: transforms must be float array")
		}
		ah := heap.Handle(args[2].IVal)
		cnt, ok := argInt(args[3])
		if !ok || cnt < 1 {
			return value.Nil, fmt.Errorf("MESH.DRAWINSTANCED: count must be >= 1")
		}
		flat := m.h.ArrayFlatLen(ah)
		if flat != int(cnt)*16 {
			return value.Nil, fmt.Errorf("MESH.DRAWINSTANCED: expected %d floats (16 per instance), got %d", cnt*16, flat)
		}
		mats := make([]rl.Matrix, int(cnt))
		for i := 0; i < int(cnt); i++ {
			var m16 [16]float32
			for j := 0; j < 16; j++ {
				f, ok := m.h.ArrayGetFloat(ah, int64(i*16+j))
				if !ok {
					return value.Nil, fmt.Errorf("MESH.DRAWINSTANCED: bad matrix data")
				}
				m16[j] = float32(f)
			}
			mats[i] = rl.Matrix{
				M0: m16[0], M4: m16[4], M8: m16[8], M12: m16[12],
				M1: m16[1], M5: m16[5], M9: m16[9], M13: m16[13],
				M2: m16[2], M6: m16[6], M10: m16[10], M14: m16[14],
				M3: m16[3], M7: m16[7], M11: m16[11], M15: m16[15],
			}
		}
		if shadowDeferActive() && InCamera3D() {
			// Instanced deferred path not implemented — draw immediately.
		}
		son := shadowDeferActive()
		if mato.pbr && son {
			bindPBRDrawState(mato, true)
		}
		drawMeshInstancedMO(mo, mato, mats, cnt)
		if mato.pbr && son {
			clearShadowMapSlot(&mato.mat)
		}
		return value.Nil, nil
	}))

	reg.Register("MESH.MAKECUSTOM", "mesh", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("MESH.MAKECUSTOM expects (verts_array, indices_array)")
		}
		_ = rt
		if args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("MESH.MAKECUSTOM: both arguments must be array handles")
		}
		vh := heap.Handle(args[0].IVal)
		ih := heap.Handle(args[1].IVal)
		vn := m.h.ArrayFlatLen(vh)
		in := m.h.ArrayFlatLen(ih)
		if vn < 8 || vn%8 != 0 {
			return value.Nil, fmt.Errorf("MESH.MAKECUSTOM: verts must have length multiple of 8 (x,y,z,nx,ny,nz,u,v)")
		}
		if in < 3 || in%3 != 0 {
			return value.Nil, fmt.Errorf("MESH.MAKECUSTOM: indices must have length multiple of 3")
		}
		vc := vn / 8
		tc := in / 3
		verts := make([]float32, vc*3)
		norms := make([]float32, vc*3)
		uvs := make([]float32, vc*2)
		for i := 0; i < vc; i++ {
			base := int64(i * 8)
			for j := 0; j < 3; j++ {
				f, ok := m.h.ArrayGetFloat(vh, base+int64(j))
				if !ok {
					return value.Nil, fmt.Errorf("MESH.MAKECUSTOM: bad vertex data")
				}
				verts[i*3+j] = float32(f)
			}
			for j := 0; j < 3; j++ {
				f, ok := m.h.ArrayGetFloat(vh, base+3+int64(j))
				if !ok {
					return value.Nil, fmt.Errorf("MESH.MAKECUSTOM: bad normal data")
				}
				norms[i*3+j] = float32(f)
			}
			for j := 0; j < 2; j++ {
				f, ok := m.h.ArrayGetFloat(vh, base+6+int64(j))
				if !ok {
					return value.Nil, fmt.Errorf("MESH.MAKECUSTOM: bad texcoord data")
				}
				uvs[i*2+j] = float32(f)
			}
		}
		idx := make([]uint16, in)
		for i := 0; i < in; i++ {
			f, ok := m.h.ArrayGetFloat(ih, int64(i))
			if !ok {
				return value.Nil, fmt.Errorf("MESH.MAKECUSTOM: bad index data")
			}
			if f < 0 || f >= float64(vc) {
				return value.Nil, fmt.Errorf("MESH.MAKECUSTOM: index out of range")
			}
			idx[i] = uint16(math.Round(f))
		}
		rm := rl.Mesh{}
		rm.VertexCount = int32(vc)
		rm.TriangleCount = int32(tc)
		rm.Vertices = unsafe.SliceData(verts)
		rm.Normals = unsafe.SliceData(norms)
		rm.Texcoords = unsafe.SliceData(uvs)
		rm.Indices = unsafe.SliceData(idx)
		obj := &meshObj{
			m: rm, pinVerts: verts, pinNorms: norms, pinUVs: uvs, pinIdx: idx,
		}
		obj.setFinalizer()
		id, err := m.h.Alloc(obj)
		if err != nil {
			return value.Nil, err
		}
		o, err := heap.Cast[*meshObj](m.h, id)
		if err != nil {
			return value.Nil, err
		}
		rl.UploadMesh(&o.m, false)
		return value.FromHandle(id), nil
	})
	reg.Register("MESH.CREATECUSTOM", "mesh", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("MESH.CREATECUSTOM expects (verts_array, indices_array)")
		}
		_ = rt
		if args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("MESH.CREATECUSTOM: both arguments must be array handles")
		}
		vh := heap.Handle(args[0].IVal)
		ih := heap.Handle(args[1].IVal)
		vn := m.h.ArrayFlatLen(vh)
		in := m.h.ArrayFlatLen(ih)
		if vn < 8 || vn%8 != 0 {
			return value.Nil, fmt.Errorf("MESH.CREATECUSTOM: verts must have length multiple of 8 (x,y,z,nx,ny,nz,u,v)")
		}
		if in < 3 || in%3 != 0 {
			return value.Nil, fmt.Errorf("MESH.CREATECUSTOM: indices must have length multiple of 3")
		}
		vc := vn / 8
		tc := in / 3
		verts := make([]float32, vc*3)
		norms := make([]float32, vc*3)
		uvs := make([]float32, vc*2)
		for i := 0; i < vc; i++ {
			base := int64(i * 8)
			for j := 0; j < 3; j++ {
				f, ok := m.h.ArrayGetFloat(vh, base+int64(j))
				if !ok {
					return value.Nil, fmt.Errorf("MESH.CREATECUSTOM: bad vertex data")
				}
				verts[i*3+j] = float32(f)
			}
			for j := 0; j < 3; j++ {
				f, ok := m.h.ArrayGetFloat(vh, base+3+int64(j))
				if !ok {
					return value.Nil, fmt.Errorf("MESH.CREATECUSTOM: bad normal data")
				}
				norms[i*3+j] = float32(f)
			}
			for j := 0; j < 2; j++ {
				f, ok := m.h.ArrayGetFloat(vh, base+6+int64(j))
				if !ok {
					return value.Nil, fmt.Errorf("MESH.CREATECUSTOM: bad texcoord data")
				}
				uvs[i*2+j] = float32(f)
			}
		}
		idx := make([]uint16, in)
		for i := 0; i < in; i++ {
			f, ok := m.h.ArrayGetFloat(ih, int64(i))
			if !ok {
				return value.Nil, fmt.Errorf("MESH.CREATECUSTOM: bad index data")
			}
			if f < 0 || f >= float64(vc) {
				return value.Nil, fmt.Errorf("MESH.CREATECUSTOM: index out of range")
			}
			idx[i] = uint16(math.Round(f))
		}
		rm := rl.Mesh{}
		rm.VertexCount = int32(vc)
		rm.TriangleCount = int32(tc)
		rm.Vertices = unsafe.SliceData(verts)
		rm.Normals = unsafe.SliceData(norms)
		rm.Texcoords = unsafe.SliceData(uvs)
		rm.Indices = unsafe.SliceData(idx)
		obj := &meshObj{
			m: rm, pinVerts: verts, pinNorms: norms, pinUVs: uvs, pinIdx: idx,
		}
		obj.setFinalizer()
		id, err := m.h.Alloc(obj)
		if err != nil {
			return value.Nil, err
		}
		o, err := heap.Cast[*meshObj](m.h, id)
		if err != nil {
			return value.Nil, err
		}
		rl.UploadMesh(&o.m, false)
		return value.FromHandle(id), nil
	})

	reg.Register("MESH.MAKECAPSULE", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		_ = m
		return value.Nil, fmt.Errorf("MESH.MAKECAPSULE: not available — this raylib build has no GenMeshCapsule; use MESH.MAKECYLINDER and spheres or DRAW3D.CAPSULE")
	}))
	reg.Register("MESH.CREATECAPSULE", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		_ = m
		return value.Nil, fmt.Errorf("MESH.CREATECAPSULE: not available — this raylib build has no GenMeshCapsule; use MESH.MAKECYLINDER and spheres or DRAW3D.CAPSULE")
	}))

	stubOpt := func(name string) func([]value.Value) (value.Value, error) {
		return func(args []value.Value) (value.Value, error) {
			_ = args
			return value.Nil, fmt.Errorf("%s: meshoptimizer not linked in moonbasic — use external tooling to optimise meshes", name)
		}
	}
	reg.Register("MESH.OPTIMISEALL", "mesh", runtime.AdaptLegacy(stubOpt("MESH.OPTIMISEALL")))
	reg.Register("MESH.OPTIMIZEALL", "mesh", runtime.AdaptLegacy(stubOpt("MESH.OPTIMIZEALL")))
	reg.Register("MESH.OPTIMISEVERTEXCACHE", "mesh", runtime.AdaptLegacy(stubOpt("MESH.OPTIMISEVERTEXCACHE")))
	reg.Register("MESH.OPTIMIZEVERTEXCACHE", "mesh", runtime.AdaptLegacy(stubOpt("MESH.OPTIMIZEVERTEXCACHE")))
	reg.Register("MESH.OPTIMISEOVERDRAW", "mesh", runtime.AdaptLegacy(stubOpt("MESH.OPTIMISEOVERDRAW")))
	reg.Register("MESH.OPTIMIZEOVERDRAW", "mesh", runtime.AdaptLegacy(stubOpt("MESH.OPTIMIZEOVERDRAW")))
	reg.Register("MESH.OPTIMISEFETCH", "mesh", runtime.AdaptLegacy(stubOpt("MESH.OPTIMISEFETCH")))
	reg.Register("MESH.OPTIMIZEFETCH", "mesh", runtime.AdaptLegacy(stubOpt("MESH.OPTIMIZEFETCH")))
	reg.Register("MESH.GENERATELOD", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		_ = args
		return value.Nil, fmt.Errorf("MESH.GENERATELOD: not implemented — use MODEL.LOADLOD for file-based LOD")
	}))
	reg.Register("MESH.GENERATELODCHAIN", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		_ = args
		return value.Nil, fmt.Errorf("MESH.GENERATELODCHAIN: not implemented")
	}))
}

func meshUpdateVertexOne(m *Module, args []value.Value) (value.Value, error) {
	// Reuse logic: delegate to same code path as MESH.UPDATEVERTEX registration
	if len(args) != 10 {
		return value.Nil, fmt.Errorf("internal")
	}
	o, err := m.getMesh(args, 0, "MESH.UPDATEVERTICES")
	if err != nil {
		return value.Nil, err
	}
	idx, ok := argInt(args[1])
	if !ok || idx < 0 || idx >= o.m.VertexCount {
		return value.Nil, fmt.Errorf("MESH.UPDATEVERTICES: invalid vertex index")
	}
	vi := int(idx)
	x, ok1 := argFloat(args[2])
	y, ok2 := argFloat(args[3])
	z, ok3 := argFloat(args[4])
	nx, ok4 := argFloat(args[5])
	ny, ok5 := argFloat(args[6])
	nz, ok6 := argFloat(args[7])
	u, ok7 := argFloat(args[8])
	vCoord, ok8 := argFloat(args[9])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 {
		return value.Nil, fmt.Errorf("MESH.UPDATEVERTICES: components must be numeric")
	}
	if o.m.Vertices != nil {
		verts := unsafe.Slice(o.m.Vertices, int(o.m.VertexCount)*3)
		verts[vi*3+0] = x
		verts[vi*3+1] = y
		verts[vi*3+2] = z
	}
	if o.m.Normals != nil {
		norms := unsafe.Slice(o.m.Normals, int(o.m.VertexCount)*3)
		norms[vi*3+0] = nx
		norms[vi*3+1] = ny
		norms[vi*3+2] = nz
	}
	if o.m.Texcoords != nil {
		uvs := unsafe.Slice(o.m.Texcoords, int(o.m.VertexCount)*2)
		uvs[vi*2+0] = u
		uvs[vi*2+1] = vCoord
	}
	if o.m.VaoID != 0 {
		if o.m.Vertices != nil {
			verts := unsafe.Slice(o.m.Vertices, int(o.m.VertexCount)*3)
			off := vi * 12
			buf := unsafe.Slice((*byte)(unsafe.Pointer(&verts[vi*3])), 12)
			rl.UpdateMeshBuffer(o.m, meshBufPosition, buf, off)
		}
		if o.m.Texcoords != nil {
			uvs := unsafe.Slice(o.m.Texcoords, int(o.m.VertexCount)*2)
			off := vi * 8
			buf := unsafe.Slice((*byte)(unsafe.Pointer(&uvs[vi*2])), 8)
			rl.UpdateMeshBuffer(o.m, meshBufTexcoord, buf, off)
		}
		if o.m.Normals != nil {
			norms := unsafe.Slice(o.m.Normals, int(o.m.VertexCount)*3)
			off := vi * 12
			buf := unsafe.Slice((*byte)(unsafe.Pointer(&norms[vi*3])), 12)
			rl.UpdateMeshBuffer(o.m, meshBufNormal, buf, off)
		}
	}
	return value.Nil, nil
}
