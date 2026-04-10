//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"
	"math"
	"unsafe"

	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// meshBuilderObj holds CPU-side vertices (xyz…) and triangle indices for procedural meshes.
type meshBuilderObj struct {
	release heap.ReleaseOnce
	verts   []float32
	idx     []uint16
	brushH  heap.Handle // optional TagBrush for GetSurfaceBrush / PaintSurface
}

func (o *meshBuilderObj) TypeName() string { return "MeshBuilder" }

func (o *meshBuilderObj) TypeTag() uint16 { return heap.TagMeshBuilder }

func (o *meshBuilderObj) Free() {
	o.release.Do(func() {
		o.verts = nil
		o.idx = nil
	})
}

func castMeshBuilder(h *heap.Store, hid heap.Handle) (*meshBuilderObj, error) {
	return heap.Cast[*meshBuilderObj](h, hid)
}

func (m *Module) entCreateSurface(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CreateSurface expects (entity#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("CreateSurface: invalid entity")
	}
	e := m.store().ents[id]
	ext := e.getExt()
	if e == nil || ext.procMeshH == 0 {
		return value.Nil, fmt.Errorf("CreateSurface: entity has no procedural mesh (use CreateMesh)")
	}
	return value.FromHandle(ext.procMeshH), nil
}

func (m *Module) entAddVertex(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("AddVertex expects (surface, x#, y#, z#)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("AddVertex: surface must be mesh builder handle")
	}
	b, err := castMeshBuilder(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argF32(args[1])
	y, ok2 := argF32(args[2])
	z, ok3 := argF32(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("AddVertex: coordinates must be numeric")
	}
	b.verts = append(b.verts, x, y, z)
	vi := int64(len(b.verts)/3 - 1)
	return value.FromInt(vi), nil
}

func (m *Module) entAddTriangle(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("AddTriangle expects (surface, v1#, v2#, v3#)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("AddTriangle: surface must be mesh builder handle")
	}
	b, err := castMeshBuilder(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	v0, ok0 := args[1].ToInt()
	v1, ok1 := args[2].ToInt()
	v2, ok2 := args[3].ToInt()
	if !ok0 || !ok1 || !ok2 || v0 < 0 || v1 < 0 || v2 < 0 {
		return value.Nil, fmt.Errorf("AddTriangle: vertex indices must be non-negative integers")
	}
	vc := len(b.verts) / 3
	if int(v0) >= vc || int(v1) >= vc || int(v2) >= vc {
		return value.Nil, fmt.Errorf("AddTriangle: vertex index out of range")
	}
	b.idx = append(b.idx, uint16(v0), uint16(v1), uint16(v2))
	return value.Nil, nil
}

func (m *Module) entVertexComponent(args []value.Value, axis int) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VertexX/Y/Z expects (surface, index#)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("surface must be mesh builder handle")
	}
	b, err := castMeshBuilder(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	ix, ok := args[1].ToInt()
	if !ok || ix < 0 {
		return value.Nil, fmt.Errorf("invalid vertex index")
	}
	base := int(ix) * 3
	if base+2 >= len(b.verts) {
		return value.Nil, fmt.Errorf("vertex index out of range")
	}
	return value.FromFloat(float64(b.verts[base+axis])), nil
}

func (m *Module) entVertexX(args []value.Value) (value.Value, error) { return m.entVertexComponent(args, 0) }
func (m *Module) entVertexY(args []value.Value) (value.Value, error) { return m.entVertexComponent(args, 1) }
func (m *Module) entVertexZ(args []value.Value) (value.Value, error) { return m.entVertexComponent(args, 2) }

func (m *Module) entUpdateMesh(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("UpdateMesh expects (entity#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("UpdateMesh: invalid entity")
	}
	e := m.store().ents[id]
	ext := e.getExt()
	if e == nil || ext.procMeshH == 0 {
		return value.Nil, fmt.Errorf("UpdateMesh: not a procedural mesh entity")
	}
	b, err := castMeshBuilder(m.h, ext.procMeshH)
	if err != nil {
		return value.Nil, err
	}
	if len(b.idx) < 3 || len(b.idx)%3 != 0 {
		return value.Nil, fmt.Errorf("UpdateMesh: need complete triangles (AddTriangle)")
	}
	vc := len(b.verts) / 3
	if vc < 3 {
		return value.Nil, fmt.Errorf("UpdateMesh: need at least 3 vertices")
	}

	norms := make([]float32, vc*3)
	for i := 0; i < len(b.idx); i += 3 {
		i0, i1, i2 := int(b.idx[i]), int(b.idx[i+1]), int(b.idx[i+2])
		if i0 >= vc || i1 >= vc || i2 >= vc {
			return value.Nil, fmt.Errorf("UpdateMesh: index out of range")
		}
		p0 := rl.Vector3{X: b.verts[i0*3], Y: b.verts[i0*3+1], Z: b.verts[i0*3+2]}
		p1 := rl.Vector3{X: b.verts[i1*3], Y: b.verts[i1*3+1], Z: b.verts[i1*3+2]}
		p2 := rl.Vector3{X: b.verts[i2*3], Y: b.verts[i2*3+1], Z: b.verts[i2*3+2]}
		e0 := rl.Vector3Subtract(p1, p0)
		e1 := rl.Vector3Subtract(p2, p0)
		fn := rl.Vector3CrossProduct(e0, e1)
		fn = rl.Vector3Normalize(fn)
		for _, vi := range []int{i0, i1, i2} {
			norms[vi*3] += fn.X
			norms[vi*3+1] += fn.Y
			norms[vi*3+2] += fn.Z
		}
	}
	for i := 0; i < vc; i++ {
		nx := norms[i*3]
		ny := norms[i*3+1]
		nz := norms[i*3+2]
		l := float32(math.Hypot(float64(nx), math.Hypot(float64(ny), float64(nz))))
		if l > 1e-8 {
			norms[i*3] /= l
			norms[i*3+1] /= l
			norms[i*3+2] /= l
		} else {
			norms[i*3], norms[i*3+1], norms[i*3+2] = 0, 1, 0
		}
	}

	pinV := append([]float32(nil), b.verts...)
	pinN := norms
	pinUV := make([]float32, vc*2)
	pinI := append([]uint16(nil), b.idx...)

	rm := rl.Mesh{}
	rm.VertexCount = int32(vc)
	rm.TriangleCount = int32(len(pinI) / 3)
	rm.Vertices = unsafe.SliceData(pinV)
	rm.Normals = unsafe.SliceData(pinN)
	rm.Texcoords = unsafe.SliceData(pinUV)
	rm.Indices = unsafe.SliceData(pinI)

	if e.hasRLModel {
		rl.UnloadModel(e.rlModel)
		e.hasRLModel = false
	}
	mod := rl.LoadModelFromMesh(rm)
	rl.UnloadMesh(&rm)
	if mod.MeshCount <= 0 {
		rl.UnloadModel(mod)
		return value.Nil, fmt.Errorf("UpdateMesh: GPU upload failed")
	}
	e.rlModel = mod
	e.hasRLModel = true
	e.kind = entKindMesh
	e.hidden = false
	return value.Nil, nil
}
