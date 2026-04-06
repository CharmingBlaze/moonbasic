//go:build cgo || (windows && !cgo)

package mbmodel3d

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/mbimage"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerMeshGen(m *Module, reg runtime.Registrar) {
	reg.Register("MESH.MAKEPOLY", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("MESH.MAKEPOLY expects sides, radius")
		}
		sides, ok1 := argInt(args[0])
		rad, ok2 := argFloat(args[1])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("MESH.MAKEPOLY: numeric sides and radius")
		}
		mesh := rl.GenMeshPoly(int(sides), rad)
		return m.allocMesh(mesh, "MESH.MAKEPOLY")
	}))

	reg.Register("MESH.MAKEPLANE", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("MESH.MAKEPLANE expects width, length, resX, resZ")
		}
		w, ok1 := argFloat(args[0])
		l, ok2 := argFloat(args[1])
		rx, ok3 := argInt(args[2])
		rz, ok4 := argInt(args[3])
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return value.Nil, fmt.Errorf("MESH.MAKEPLANE: numeric arguments")
		}
		mesh := rl.GenMeshPlane(w, l, int(rx), int(rz))
		return m.allocMesh(mesh, "MESH.MAKEPLANE")
	}))

	reg.Register("MESH.MAKECUBE", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("MESH.MAKECUBE expects width, height, length")
		}
		w, ok1 := argFloat(args[0])
		h, ok2 := argFloat(args[1])
		l, ok3 := argFloat(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("MESH.MAKECUBE: numeric dimensions")
		}
		mesh := rl.GenMeshCube(w, h, l)
		return m.allocMesh(mesh, "MESH.MAKECUBE")
	}))

	reg.Register("MESH.MAKESPHERE", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("MESH.MAKESPHERE expects radius, rings, slices")
		}
		rad, ok1 := argFloat(args[0])
		rings, ok2 := argInt(args[1])
		slices, ok3 := argInt(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("MESH.MAKESPHERE: numeric arguments")
		}
		mesh := rl.GenMeshSphere(rad, int(rings), int(slices))
		return m.allocMesh(mesh, "MESH.MAKESPHERE")
	}))

	reg.Register("MESH.MAKECYLINDER", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("MESH.MAKECYLINDER expects radius, height, slices")
		}
		rad, ok1 := argFloat(args[0])
		h, ok2 := argFloat(args[1])
		sl, ok3 := argInt(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("MESH.MAKECYLINDER: numeric arguments")
		}
		mesh := rl.GenMeshCylinder(rad, h, int(sl))
		return m.allocMesh(mesh, "MESH.MAKECYLINDER")
	}))

	reg.Register("MESH.MAKECONE", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("MESH.MAKECONE expects radius, height, slices")
		}
		rad, ok1 := argFloat(args[0])
		h, ok2 := argFloat(args[1])
		sl, ok3 := argInt(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("MESH.MAKECONE: numeric arguments")
		}
		mesh := rl.GenMeshCone(rad, h, int(sl))
		return m.allocMesh(mesh, "MESH.MAKECONE")
	}))

	reg.Register("MESH.MAKETORUS", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("MESH.MAKETORUS expects radius, size, radSeg, sides")
		}
		rad, ok1 := argFloat(args[0])
		sz, ok2 := argFloat(args[1])
		rseg, ok3 := argInt(args[2])
		sides, ok4 := argInt(args[3])
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return value.Nil, fmt.Errorf("MESH.MAKETORUS: numeric arguments")
		}
		mesh := rl.GenMeshTorus(rad, sz, int(rseg), int(sides))
		return m.allocMesh(mesh, "MESH.MAKETORUS")
	}))

	reg.Register("MESH.MAKEKNOT", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("MESH.MAKEKNOT expects radius, size, radSeg, sides")
		}
		rad, ok1 := argFloat(args[0])
		sz, ok2 := argFloat(args[1])
		rseg, ok3 := argInt(args[2])
		sides, ok4 := argInt(args[3])
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return value.Nil, fmt.Errorf("MESH.MAKEKNOT: numeric arguments")
		}
		mesh := rl.GenMeshKnot(rad, sz, int(rseg), int(sides))
		return m.allocMesh(mesh, "MESH.MAKEKNOT")
	}))

	reg.Register("MESH.MAKEHEIGHTMAP", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 4 || args[0].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("MESH.MAKEHEIGHTMAP expects image handle, sizeX, sizeY, sizeZ")
		}
		img, err := mbimage.RayImageForTexture(m.h, heap.Handle(args[0].IVal))
		if err != nil {
			return value.Nil, fmt.Errorf("MESH.MAKEHEIGHTMAP: %w", err)
		}
		sx, ok1 := argFloat(args[1])
		sy, ok2 := argFloat(args[2])
		sz, ok3 := argFloat(args[3])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("MESH.MAKEHEIGHTMAP: size must be numeric")
		}
		mesh := rl.GenMeshHeightmap(*img, rl.Vector3{X: sx, Y: sy, Z: sz})
		return m.allocMesh(mesh, "MESH.MAKEHEIGHTMAP")
	}))

	reg.Register("MESH.MAKECUBICMAP", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 4 || args[0].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("MESH.MAKECUBICMAP expects image handle, cubeX, cubeY, cubeZ")
		}
		img, err := mbimage.RayImageForTexture(m.h, heap.Handle(args[0].IVal))
		if err != nil {
			return value.Nil, fmt.Errorf("MESH.MAKECUBICMAP: %w", err)
		}
		cx, ok1 := argFloat(args[1])
		cy, ok2 := argFloat(args[2])
		cz, ok3 := argFloat(args[3])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("MESH.MAKECUBICMAP: cube size must be numeric")
		}
		mesh := rl.GenMeshCubicmap(*img, rl.Vector3{X: cx, Y: cy, Z: cz})
		return m.allocMesh(mesh, "MESH.MAKECUBICMAP")
	}))

	// Legacy manifest names (aliases)
	reg.Register("MESH.CUBE", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return m.allocMesh(rl.GenMeshCube(func() float32 { v, _ := argFloat(args[0]); return v }(), func() float32 { v, _ := argFloat(args[1]); return v }(), func() float32 { v, _ := argFloat(args[2]); return v }()), "MESH.CUBE")
	}))
	reg.Register("MESH.SPHERE", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return m.allocMesh(rl.GenMeshSphere(func() float32 { v, _ := argFloat(args[0]); return v }(), func() int { v, _ := argInt(args[1]); return int(v) }(), func() int { v, _ := argInt(args[2]); return int(v) }()), "MESH.SPHERE")
	}))
	reg.Register("MESH.PLANE", "mesh", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return m.allocMesh(rl.GenMeshPlane(func() float32 { v, _ := argFloat(args[0]); return v }(), func() float32 { v, _ := argFloat(args[1]); return v }(), func() int { v, _ := argInt(args[2]); return int(v) }(), func() int { v, _ := argInt(args[3]); return int(v) }()), "MESH.PLANE")
	}))
}
