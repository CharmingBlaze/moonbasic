//go:build cgo || (windows && !cgo)

package mbcollision

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/runtime/mbmodel3d"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerRayBuiltins(reg runtime.Registrar) {
	reg.Register("RAY.MAKE", "collision", runtime.AdaptLegacy(m.rayMake))
	reg.Register("RAY.CREATE", "collision", runtime.AdaptLegacy(m.rayCreate))
	reg.Register("RAY.FREE", "collision", runtime.AdaptLegacy(m.rayFree))
	reg.Register("RAY.SETPOS", "collision", runtime.AdaptLegacy(m.raySetPos))
	reg.Register("RAY.SETPOSITION", "collision", runtime.AdaptLegacy(m.raySetPos))
	reg.Register("RAY.SETDIR", "collision", runtime.AdaptLegacy(m.raySetDir))
	reg.Register("RAY.GETPOS", "collision", runtime.AdaptLegacy(m.rayGetPos))
	reg.Register("RAY.GETDIR", "collision", runtime.AdaptLegacy(m.rayGetDir))

	rayCollisionScalars(reg, "RAY.HITSPHERE", 5, func(args []value.Value) (rl.RayCollision, error) {
		if err := m.requireHeap(); err != nil {
			return rl.RayCollision{}, err
		}
		if len(args) != 5 {
			return rl.RayCollision{}, fmt.Errorf("RAY.HITSPHERE_* expects 5 arguments (ray, cx, cy, cz, r)")
		}
		ray, err := m.rayFromArgs(args, 0, "RAY.HITSPHERE")
		if err != nil {
			return rl.RayCollision{}, err
		}
		cx, ok1 := argF(args[1])
		cy, ok2 := argF(args[2])
		cz, ok3 := argF(args[3])
		rad, ok4 := argF(args[4])
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return rl.RayCollision{}, fmt.Errorf("RAY.HITSPHERE_*: sphere parameters must be numeric")
		}
		center := rl.Vector3{X: cx, Y: cy, Z: cz}
		return rl.GetRayCollisionSphere(ray, center, rad), nil
	})

	rayCollisionScalars(reg, "RAY.HITBOX", 7, func(args []value.Value) (rl.RayCollision, error) {
		if err := m.requireHeap(); err != nil {
			return rl.RayCollision{}, err
		}
		if len(args) != 7 {
			return rl.RayCollision{}, fmt.Errorf("RAY.HITBOX_* expects 7 arguments (ray, minX, minY, minZ, maxX, maxY, maxZ)")
		}
		ray, err := m.rayFromArgs(args, 0, "RAY.HITBOX")
		if err != nil {
			return rl.RayCollision{}, err
		}
		var min, max rl.Vector3
		for i, tgt := range []*rl.Vector3{&min, &max} {
			base := 1 + i*3
			x, okx := argF(args[base])
			y, oky := argF(args[base+1])
			z, okz := argF(args[base+2])
			if !okx || !oky || !okz {
				return rl.RayCollision{}, fmt.Errorf("RAY.HITBOX_*: box corner must be numeric")
			}
			tgt.X, tgt.Y, tgt.Z = x, y, z
		}
		box := rl.BoundingBox{Min: min, Max: max}
		return rl.GetRayCollisionBox(ray, box), nil
	})

	rayCollisionScalars(reg, "RAY.HITPLANE", 5, func(args []value.Value) (rl.RayCollision, error) {
		if err := m.requireHeap(); err != nil {
			return rl.RayCollision{}, err
		}
		if len(args) != 5 {
			return rl.RayCollision{}, fmt.Errorf("RAY.HITPLANE_* expects 5 arguments (ray, nx, ny, nz, d)")
		}
		ray, err := m.rayFromArgs(args, 0, "RAY.HITPLANE")
		if err != nil {
			return rl.RayCollision{}, err
		}
		nx, ok1 := argF(args[1])
		ny, ok2 := argF(args[2])
		nz, ok3 := argF(args[3])
		d, ok4 := argF(args[4])
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return rl.RayCollision{}, fmt.Errorf("RAY.HITPLANE_*: plane parameters must be numeric")
		}
		return rayPlaneCollision(ray, nx, ny, nz, d), nil
	})

	rayCollisionScalars(reg, "RAY.HITTRIANGLE", 10, func(args []value.Value) (rl.RayCollision, error) {
		if err := m.requireHeap(); err != nil {
			return rl.RayCollision{}, err
		}
		if len(args) != 10 {
			return rl.RayCollision{}, fmt.Errorf("RAY.HITTRIANGLE_* expects 10 arguments (ray, p1x,p1y,p1z, p2x,p2y,p2z, p3x,p3y,p3z)")
		}
		ray, err := m.rayFromArgs(args, 0, "RAY.HITTRIANGLE")
		if err != nil {
			return rl.RayCollision{}, err
		}
		var v1, v2, v3 rl.Vector3
		for i, tgt := range []*rl.Vector3{&v1, &v2, &v3} {
			base := 1 + i*3
			x, okx := argF(args[base])
			y, oky := argF(args[base+1])
			z, okz := argF(args[base+2])
			if !okx || !oky || !okz {
				return rl.RayCollision{}, fmt.Errorf("RAY.HITTRIANGLE_*: vertex must be numeric")
			}
			tgt.X, tgt.Y, tgt.Z = x, y, z
		}
		return rl.GetRayCollisionTriangle(ray, v1, v2, v3), nil
	})

	rayCollisionScalars(reg, "RAY.HITMESH", 3, func(args []value.Value) (rl.RayCollision, error) {
		if err := m.requireHeap(); err != nil {
			return rl.RayCollision{}, err
		}
		if len(args) != 3 {
			return rl.RayCollision{}, fmt.Errorf("RAY.HITMESH_* expects 3 arguments (ray, mesh, matrixHandle)")
		}
		ray, err := m.rayFromArgs(args, 0, "RAY.HITMESH")
		if err != nil {
			return rl.RayCollision{}, err
		}
		if args[1].Kind != value.KindHandle {
			return rl.RayCollision{}, fmt.Errorf("RAY.HITMESH_*: mesh must be handle")
		}
		matH := heap.Handle(0)
		if args[2].Kind == value.KindHandle {
			matH = heap.Handle(args[2].IVal)
		} else if i, ok := args[2].ToInt(); ok && i == 0 {
			matH = 0
		} else if f, ok := argF(args[2]); ok && f == 0 {
			matH = 0
		} else {
			return rl.RayCollision{}, fmt.Errorf("RAY.HITMESH_*: matrix must be handle or 0")
		}
		mesh, err := mbmodel3d.MeshRaylib(m.h, heap.Handle(args[1].IVal))
		if err != nil {
			return rl.RayCollision{}, fmt.Errorf("RAY.HITMESH_*: %w", err)
		}
		mat, err := mbmatrix.MatrixRaylib(m.h, matH)
		if err != nil {
			return rl.RayCollision{}, fmt.Errorf("RAY.HITMESH_*: matrix: %w", err)
		}
		return rl.GetRayCollisionMesh(ray, mesh, mat), nil
	})

	hitModelRay := func(op string) func([]value.Value) (rl.RayCollision, error) {
		return func(args []value.Value) (rl.RayCollision, error) {
			if err := m.requireHeap(); err != nil {
				return rl.RayCollision{}, err
			}
			if len(args) != 2 {
				return rl.RayCollision{}, fmt.Errorf("%s_* expects 2 arguments (ray, model)", op)
			}
			ray, err := m.rayFromArgs(args, 0, op)
			if err != nil {
				return rl.RayCollision{}, err
			}
			if args[1].Kind != value.KindHandle {
				return rl.RayCollision{}, fmt.Errorf("%s_*: model must be handle", op)
			}
			model, err := mbmodel3d.ModelRaylib(m.h, heap.Handle(args[1].IVal))
			if err != nil {
				return rl.RayCollision{}, fmt.Errorf("%s_*: %w", op, err)
			}
			return rayModelCollision(ray, model), nil
		}
	}
	rayCollisionScalars(reg, "RAY.HITMODEL", 2, hitModelRay("RAY.HITMODEL"))
	rayCollisionScalars(reg, "RAY.INTERSECTSMODEL", 2, hitModelRay("RAY.INTERSECTSMODEL"))
}

func (m *Module) rayMake(args []value.Value) (value.Value, error) {
	return m.allocRay(args, "RAY.MAKE")
}

func (m *Module) rayCreate(args []value.Value) (value.Value, error) {
	return m.allocRay(args, "RAY.CREATE")
}

func (m *Module) allocRay(args []value.Value, op string) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("%s expects 6 arguments (ox, oy, oz, dx, dy, dz)", op)
	}
	var pos, dir rl.Vector3
	for i, tgt := range []*rl.Vector3{&pos, &dir} {
		base := i * 3
		x, okx := argF(args[base])
		y, oky := argF(args[base+1])
		z, okz := argF(args[base+2])
		if !okx || !oky || !okz {
			return value.Nil, fmt.Errorf("%s: arguments must be numeric", op)
		}
		tgt.X, tgt.Y, tgt.Z = x, y, z
	}
	o := &rayObj{r: rl.Ray{Position: pos, Direction: dir}}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) rayFree(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("RAY.FREE expects ray handle")
	}
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func (m *Module) raySetPos(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("RAY.SETPOS expects (handle, x, y, z)")
	}
	o, err := heap.Cast[*rayObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("RAY.SETPOS: %w", err)
	}
	x, ok1 := argF(args[1])
	y, ok2 := argF(args[2])
	z, ok3 := argF(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("RAY.SETPOS: components must be numeric")
	}
	o.r.Position = rl.Vector3{X: x, Y: y, Z: z}
	return args[0], nil
}

func (m *Module) raySetDir(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("RAY.SETDIR expects (handle, x, y, z)")
	}
	o, err := heap.Cast[*rayObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("RAY.SETDIR: %w", err)
	}
	x, ok1 := argF(args[1])
	y, ok2 := argF(args[2])
	z, ok3 := argF(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("RAY.SETDIR: components must be numeric")
	}
	o.r.Direction = rl.Vector3{X: x, Y: y, Z: z}
	return args[0], nil
}

func (m *Module) rayGetPos(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("RAY.GETPOS expects ray handle")
	}
	o, err := heap.Cast[*rayObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("RAY.GETPOS: %w", err)
	}
	return mbmatrix.AllocVec3Value(m.h, o.r.Position.X, o.r.Position.Y, o.r.Position.Z)
}

func (m *Module) rayGetDir(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("RAY.GETDIR expects ray handle")
	}
	o, err := heap.Cast[*rayObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("RAY.GETDIR: %w", err)
	}
	return mbmatrix.AllocVec3Value(m.h, o.r.Direction.X, o.r.Direction.Y, o.r.Direction.Z)
}
