//go:build cgo || (windows && !cgo)

package mbcollision

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) requireHeap() error {
	if m.h == nil {
		return runtime.Errorf("collision builtins: heap not bound")
	}
	return nil
}

func argF(v value.Value) (float32, bool) {
	if f, ok := v.ToFloat(); ok {
		return float32(f), true
	}
	if i, ok := v.ToInt(); ok {
		return float32(i), true
	}
	return 0, false
}

func (m *Module) rayFromArgs(args []value.Value, ix int, op string) (rl.Ray, error) {
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return rl.Ray{}, fmt.Errorf("%s: argument %d must be ray handle", op, ix+1)
	}
	o, err := heap.Cast[*rayObj](m.h, heap.Handle(args[ix].IVal))
	if err != nil {
		return rl.Ray{}, fmt.Errorf("%s: %w", op, err)
	}
	return o.r, nil
}

// BBoxFromArgs retrieves a bounding box from a heap handle.
func BBoxFromArgs(s *heap.Store, h heap.Handle) (rl.BoundingBox, error) {
	o, err := heap.Cast[*bboxObj](s, h)
	if err != nil {
		return rl.BoundingBox{}, err
	}
	return o.box, nil
}

func rayCollisionScalars(reg runtime.Registrar, prefix string, nargs int, compute func([]value.Value) (rl.RayCollision, error)) {
	reg.Register(prefix+"_HIT", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		c, err := compute(args)
		if err != nil {
			return value.Nil, err
		}
		return value.FromBool(c.Hit), nil
	}))
	reg.Register(prefix+"_DISTANCE", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		c, err := compute(args)
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(float64(c.Distance)), nil
	}))
	reg.Register(prefix+"_POINTX", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		c, err := compute(args)
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(float64(c.Point.X)), nil
	}))
	reg.Register(prefix+"_POINTY", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		c, err := compute(args)
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(float64(c.Point.Y)), nil
	}))
	reg.Register(prefix+"_POINTZ", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		c, err := compute(args)
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(float64(c.Point.Z)), nil
	}))
	reg.Register(prefix+"_NORMALX", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		c, err := compute(args)
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(float64(c.Normal.X)), nil
	}))
	reg.Register(prefix+"_NORMALY", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		c, err := compute(args)
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(float64(c.Normal.Y)), nil
	}))
	reg.Register(prefix+"_NORMALZ", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		c, err := compute(args)
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(float64(c.Normal.Z)), nil
	}))
	_ = nargs // manifest documents arity; runtime checks inside compute
}

// rayPlaneCollision: plane n·p = d with n normalized; ray r(t) = o + t*dir.
func rayPlaneCollision(ray rl.Ray, nx, ny, nz, d float32) rl.RayCollision {
	n := rl.Vector3{X: nx, Y: ny, Z: nz}
	n = rl.Vector3Normalize(n)
	dir := rl.Vector3Normalize(ray.Direction)
	denom := rl.Vector3DotProduct(n, dir)
	if math.Abs(float64(denom)) < 1e-6 {
		return rl.RayCollision{Hit: false}
	}
	t := (d - rl.Vector3DotProduct(n, ray.Position)) / denom
	if t < 0 {
		return rl.RayCollision{Hit: false}
	}
	pt := rl.Vector3Add(ray.Position, rl.Vector3Scale(dir, t))
	return rl.RayCollision{Hit: true, Distance: t, Point: pt, Normal: n}
}

func rayModelCollision(ray rl.Ray, model rl.Model) rl.RayCollision {
	best := rl.RayCollision{Hit: false, Distance: float32(math.MaxFloat32)}
	tform := model.Transform
	meshes := model.GetMeshes()
	for i := range meshes {
		col := rl.GetRayCollisionMesh(ray, meshes[i], tform)
		if col.Hit && col.Distance < best.Distance {
			best = col
		}
	}
	if !best.Hit {
		return rl.RayCollision{Hit: false}
	}
	return best
}
