//go:build cgo

package mbmodel3d

import (
	"fmt"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	mbmatrix "moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// Raylib mesh VBO indices (rlgl default shader attribute locations).
const (
	meshBufPosition = 0
	meshBufTexcoord = 1
	meshBufNormal   = 2
)

func registerMeshOps(m *Module, reg runtime.Registrar) {
	reg.Register("MESH.UPLOAD", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("MESH.UPLOAD expects mesh handle, dynamic")
		}
		o, err := m.getMesh(args, 0, "MESH.UPLOAD")
		if err != nil {
			return value.Nil, err
		}
		dyn, ok := argBool(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("MESH.UPLOAD: dynamic must be bool or number")
		}
		rl.UploadMesh(&o.m, dyn)
		return value.Nil, nil
	}))

	reg.Register("MESH.DRAW", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("MESH.DRAW expects mesh handle, material handle, matrix handle (or 0 for identity)")
		}
		mo, err := m.getMesh(args, 0, "MESH.DRAW")
		if err != nil {
			return value.Nil, err
		}
		mato, err := m.getMaterial(args, 1, "MESH.DRAW")
		if err != nil {
			return value.Nil, err
		}
		var transform rl.Matrix
		switch args[2].Kind {
		case value.KindHandle:
			transform, err = mbmatrix.MatrixRaylib(m.h, heap.Handle(args[2].IVal))
			if err != nil {
				return value.Nil, err
			}
		default:
			if i, ok := args[2].ToInt(); ok && i == 0 {
				transform = rl.MatrixIdentity()
			} else if f, ok := args[2].ToFloat(); ok && f == 0 {
				transform = rl.MatrixIdentity()
			} else {
				return value.Nil, fmt.Errorf("MESH.DRAW: third argument must be matrix handle or 0")
			}
		}
		if shadowDeferActive() && InCamera3D() {
			draw3dMu.Lock()
			deferredMeshes = append(deferredMeshes, deferredMeshRec{
				meshH: heap.Handle(args[0].IVal),
				matH:  heap.Handle(args[1].IVal),
				mtx:   transform,
			})
			draw3dMu.Unlock()
			return value.Nil, nil
		}
		son := shadowDeferActive()
		if mato.pbr && son {
			bindPBRDrawState(&mato.mat, true)
		}
		rl.DrawMesh(mo.m, mato.mat, transform)
		if mato.pbr && son {
			clearShadowMapSlot(&mato.mat)
		}
		return value.Nil, nil
	}))

	reg.Register("MESH.DRAWROTATED", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("MESH.DRAWROTATED expects mesh handle, material handle, rx, ry, rz (radians)")
		}
		mo, err := m.getMesh(args, 0, "MESH.DRAWROTATED")
		if err != nil {
			return value.Nil, err
		}
		mato, err := m.getMaterial(args, 1, "MESH.DRAWROTATED")
		if err != nil {
			return value.Nil, err
		}
		rx, ok1 := argFloat(args[2])
		ry, ok2 := argFloat(args[3])
		rz, ok3 := argFloat(args[4])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("MESH.DRAWROTATED: rotation angles must be numeric")
		}
		transform := rl.MatrixRotateXYZ(rl.Vector3{X: rx, Y: ry, Z: rz})
		if shadowDeferActive() && InCamera3D() {
			draw3dMu.Lock()
			deferredMeshes = append(deferredMeshes, deferredMeshRec{
				meshH: heap.Handle(args[0].IVal),
				matH:  heap.Handle(args[1].IVal),
				mtx:   transform,
			})
			draw3dMu.Unlock()
			return value.Nil, nil
		}
		son := shadowDeferActive()
		if mato.pbr && son {
			bindPBRDrawState(&mato.mat, true)
		}
		rl.DrawMesh(mo.m, mato.mat, transform)
		if mato.pbr && son {
			clearShadowMapSlot(&mato.mat)
		}
		return value.Nil, nil
	}))

	reg.Register("MESH.UPDATEVERTEX", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 10 {
			return value.Nil, fmt.Errorf("MESH.UPDATEVERTEX expects handle, idx, x,y,z, nx,ny,nz, u,v")
		}
		o, err := m.getMesh(args, 0, "MESH.UPDATEVERTEX")
		if err != nil {
			return value.Nil, err
		}
		idx, ok := argInt(args[1])
		if !ok || idx < 0 || idx >= o.m.VertexCount {
			return value.Nil, fmt.Errorf("MESH.UPDATEVERTEX: invalid vertex index")
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
			return value.Nil, fmt.Errorf("MESH.UPDATEVERTEX: components must be numeric")
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
	}))

	bboxAxis := func(axis int) func([]value.Value) (value.Value, error) {
		return func(args []value.Value) (value.Value, error) {
			if err := m.requireHeap(); err != nil {
				return value.Nil, err
			}
			if len(args) != 1 {
				return value.Nil, fmt.Errorf("MESH.GETBBOX*: expects mesh handle")
			}
			o, err := m.getMesh(args, 0, "MESH.GETBBOX")
			if err != nil {
				return value.Nil, err
			}
			b := rl.GetMeshBoundingBox(o.m)
			switch axis {
			case 0:
				return value.FromFloat(float64(b.Min.X)), nil
			case 1:
				return value.FromFloat(float64(b.Min.Y)), nil
			case 2:
				return value.FromFloat(float64(b.Min.Z)), nil
			case 3:
				return value.FromFloat(float64(b.Max.X)), nil
			case 4:
				return value.FromFloat(float64(b.Max.Y)), nil
			default:
				return value.FromFloat(float64(b.Max.Z)), nil
			}
		}
	}

	reg.Register("MESH.GETBBOXMINX", "mesh", runtime.AdaptLegacy(bboxAxis(0)))
	reg.Register("MESH.GETBBOXMINY", "mesh", runtime.AdaptLegacy(bboxAxis(1)))
	reg.Register("MESH.GETBBOXMINZ", "mesh", runtime.AdaptLegacy(bboxAxis(2)))
	reg.Register("MESH.GETBBOXMAXX", "mesh", runtime.AdaptLegacy(bboxAxis(3)))
	reg.Register("MESH.GETBBOXMAXY", "mesh", runtime.AdaptLegacy(bboxAxis(4)))
	reg.Register("MESH.GETBBOXMAXZ", "mesh", runtime.AdaptLegacy(bboxAxis(5)))

	reg.Register("MESH.GENTANGENTS", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MESH.GENTANGENTS expects mesh handle")
		}
		if _, err := m.getMesh(args, 0, "MESH.GENTANGENTS"); err != nil {
			return value.Nil, err
		}
		// raylib-go exposes GenMeshTangents only in the purego (!cgo) build, not alongside rmodels.cgo.
		return value.Nil, runtime.Errorf("MESH.GENTANGENTS: not available when moonbasic is built with CGO; use CGO_ENABLED=0 (purego raylib-go) if you need tangent generation")
	}))

	reg.Register("MESH.FREE", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 || args[0].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("MESH.FREE expects mesh handle")
		}
		if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
			return value.Nil, err
		}
		return value.Nil, nil
	}))
}
