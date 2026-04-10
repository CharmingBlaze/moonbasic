//go:build cgo || (windows && !cgo)

package mbmatrix

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerVec3(reg runtime.Registrar) {
	reg.Register("VEC3.MAKE", "vec3", runtime.AdaptLegacy(m.vec3Make))
	reg.Register("VEC3.FREE", "vec3", runtime.AdaptLegacy(m.vec3Free))
	reg.Register("VEC3.X", "vec3", runtime.AdaptLegacy(m.vec3X))
	reg.Register("VEC3.Y", "vec3", runtime.AdaptLegacy(m.vec3Y))
	reg.Register("VEC3.Z", "vec3", runtime.AdaptLegacy(m.vec3Z))
	reg.Register("VEC3.SET", "vec3", runtime.AdaptLegacy(m.vec3Set))
	reg.Register("VEC3.ADD", "vec3", runtime.AdaptLegacy(m.vec3Add))
	reg.Register("VEC3.SUB", "vec3", runtime.AdaptLegacy(m.vec3Sub))
	reg.Register("VEC3.MUL", "vec3", runtime.AdaptLegacy(m.vec3Mul))
	reg.Register("VEC3.DIV", "vec3", runtime.AdaptLegacy(m.vec3Div))
	reg.Register("VEC3.DOT", "vec3", runtime.AdaptLegacy(m.vec3Dot))
	reg.Register("VEC3.CROSS", "vec3", runtime.AdaptLegacy(m.vec3Cross))
	reg.Register("VEC3.LENGTH", "vec3", runtime.AdaptLegacy(m.vec3Length))
	reg.Register("VEC3.NORMALIZE", "vec3", runtime.AdaptLegacy(m.vec3Normalize))
	reg.Register("VEC3.LERP", "vec3", runtime.AdaptLegacy(m.vec3Lerp))
	reg.Register("VEC3.DISTANCE", "vec3", runtime.AdaptLegacy(m.vec3Distance))
	reg.Register("VEC3.DIST", "vec3", runtime.AdaptLegacy(m.vec3Dist))
	reg.Register("VEC3.DISTSQ", "vec3", runtime.AdaptLegacy(m.vec3DistSq))
	reg.Register("VEC3.REFLECT", "vec3", runtime.AdaptLegacy(m.vec3Reflect))
	reg.Register("VEC3.NEGATE", "vec3", runtime.AdaptLegacy(m.vec3Negate))
	reg.Register("VEC3.EQUALS", "vec3", runtime.AdaptLegacy(m.vec3Equals))

	reg.Register("VEC3.VEC3", "vec3", runtime.AdaptLegacy(m.vec3Make))
	reg.Register("VEC3.VECADD", "vec3", runtime.AdaptLegacy(m.vec3Add))
	reg.Register("VEC3.VECSUB", "vec3", runtime.AdaptLegacy(m.vec3Sub))
	reg.Register("VEC3.VECSCALE", "vec3", runtime.AdaptLegacy(m.vec3Mul))
	reg.Register("VEC3.VECNORMALIZE", "vec3", runtime.AdaptLegacy(m.vec3Normalize))
	reg.Register("VEC3.VECDOT", "vec3", runtime.AdaptLegacy(m.vec3Dot))
	reg.Register("VEC3.VECCROSS", "vec3", runtime.AdaptLegacy(m.vec3Cross))
	reg.Register("VEC3.VECLENGTH", "vec3", runtime.AdaptLegacy(m.vec3Length))
}

func (m *Module) allocVec3(v rl.Vector3) (value.Value, error) {
	id, err := m.h.Alloc(&vec3Obj{v: v})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) vec3Make(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("VEC3.MAKE expects 3 arguments (x, y, z)")
	}
	x, ok1 := argF(args[0])
	y, ok2 := argF(args[1])
	z, ok3 := argF(args[2])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("VEC3.MAKE: components must be numeric")
	}
	return m.allocVec3(rl.Vector3{X: x, Y: y, Z: z})
}

func (m *Module) vec3Free(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("VEC3.FREE expects vec3 handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) vec3X(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("VEC3.X expects vec3 handle")
	}
	v, err := m.vec3FromArgs(args, 0, "VEC3.X")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(v.X)), nil
}

func (m *Module) vec3Y(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("VEC3.Y expects vec3 handle")
	}
	v, err := m.vec3FromArgs(args, 0, "VEC3.Y")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(v.Y)), nil
}

func (m *Module) vec3Z(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("VEC3.Z expects vec3 handle")
	}
	v, err := m.vec3FromArgs(args, 0, "VEC3.Z")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(v.Z)), nil
}

func (m *Module) vec3Set(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("VEC3.SET expects (handle, x, y, z)")
	}
	o, err := heap.Cast[*vec3Obj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("VEC3.SET: %w", err)
	}
	x, ok1 := argF(args[1])
	y, ok2 := argF(args[2])
	z, ok3 := argF(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("VEC3.SET: components must be numeric")
	}
	o.v = rl.Vector3{X: x, Y: y, Z: z}
	return value.Nil, nil
}

func (m *Module) vec3Add(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC3.ADD expects two vec3 handles")
	}
	a, err := m.vec3FromArgs(args, 0, "VEC3.ADD")
	if err != nil {
		return value.Nil, err
	}
	b, err := m.vec3FromArgs(args, 1, "VEC3.ADD")
	if err != nil {
		return value.Nil, err
	}
	return m.allocVec3(rl.Vector3Add(a, b))
}

func (m *Module) vec3Sub(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC3.SUB expects two vec3 handles")
	}
	a, err := m.vec3FromArgs(args, 0, "VEC3.SUB")
	if err != nil {
		return value.Nil, err
	}
	b, err := m.vec3FromArgs(args, 1, "VEC3.SUB")
	if err != nil {
		return value.Nil, err
	}
	return m.allocVec3(rl.Vector3Subtract(a, b))
}

func (m *Module) vec3Mul(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC3.MUL expects (vec3, scalar)")
	}
	v, err := m.vec3FromArgs(args, 0, "VEC3.MUL")
	if err != nil {
		return value.Nil, err
	}
	s, ok := argF(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("VEC3.MUL: scalar must be numeric")
	}
	return m.allocVec3(rl.Vector3Scale(v, s))
}

func (m *Module) vec3Div(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC3.DIV expects (vec3, scalar)")
	}
	v, err := m.vec3FromArgs(args, 0, "VEC3.DIV")
	if err != nil {
		return value.Nil, err
	}
	s, ok := argF(args[1])
	if !ok || float64(s) == 0 {
		return value.Nil, fmt.Errorf("VEC3.DIV: non-zero scalar required")
	}
	return m.allocVec3(rl.Vector3Scale(v, 1/s))
}

func (m *Module) vec3Dot(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC3.DOT expects two vec3 handles")
	}
	a, err := m.vec3FromArgs(args, 0, "VEC3.DOT")
	if err != nil {
		return value.Nil, err
	}
	b, err := m.vec3FromArgs(args, 1, "VEC3.DOT")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(rl.Vector3DotProduct(a, b))), nil
}

func (m *Module) vec3Cross(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC3.CROSS expects two vec3 handles")
	}
	a, err := m.vec3FromArgs(args, 0, "VEC3.CROSS")
	if err != nil {
		return value.Nil, err
	}
	b, err := m.vec3FromArgs(args, 1, "VEC3.CROSS")
	if err != nil {
		return value.Nil, err
	}
	return m.allocVec3(rl.Vector3CrossProduct(a, b))
}

func (m *Module) vec3Length(args []value.Value) (value.Value, error) {
	if len(args) == 3 {
		x, ok1 := argF(args[0])
		y, ok2 := argF(args[1])
		z, ok3 := argF(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("VEC3.LENGTH: components must be numeric")
		}
		return value.FromFloat(math.Sqrt(float64(x*x + y*y + z*z))), nil
	}
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("VEC3.LENGTH expects vec3 handle or (x, y, z)")
	}
	v, err := m.vec3FromArgs(args, 0, "VEC3.LENGTH")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(rl.Vector3Length(v))), nil
}

func (m *Module) vec3Normalize(args []value.Value) (value.Value, error) {
	if len(args) == 3 {
		x, ok1 := argF(args[0])
		y, ok2 := argF(args[1])
		z, ok3 := argF(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("VEC3.NORMALIZE: components must be numeric")
		}
		mag := float32(math.Sqrt(float64(x*x + y*y + z*z)))
		if mag > 0 {
			x /= mag
			y /= mag
			z /= mag
		}
		return m.allocTuple3(x, y, z)
	}
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("VEC3.NORMALIZE expects vec3 handle or (x, y, z)")
	}
	v, err := m.vec3FromArgs(args, 0, "VEC3.NORMALIZE")
	if err != nil {
		return value.Nil, err
	}
	return m.allocVec3(rl.Vector3Normalize(v))
}

func (m *Module) allocTuple3(x, y, z float32) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	arr, err := heap.NewArrayOfKind([]int64{3}, heap.ArrayKindFloat, 0)
	if err != nil {
		return value.Nil, err
	}
	arr.Floats[0] = float64(x)
	arr.Floats[1] = float64(y)
	arr.Floats[2] = float64(z)
	h, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func (m *Module) vec3Lerp(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("VEC3.LERP expects (a, b, t)")
	}
	a, err := m.vec3FromArgs(args, 0, "VEC3.LERP")
	if err != nil {
		return value.Nil, err
	}
	b, err := m.vec3FromArgs(args, 1, "VEC3.LERP")
	if err != nil {
		return value.Nil, err
	}
	t, ok := argF(args[2])
	if !ok {
		return value.Nil, fmt.Errorf("VEC3.LERP: t must be numeric")
	}
	return m.allocVec3(rl.Vector3Lerp(a, b, t))
}

func (m *Module) vec3Dist(args []value.Value) (value.Value, error) {
	if len(args) == 6 {
		x1, ok1 := argF(args[0])
		y1, ok2 := argF(args[1])
		z1, ok3 := argF(args[2])
		x2, ok4 := argF(args[3])
		y2, ok5 := argF(args[4])
		z2, ok6 := argF(args[5])
		if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
			return value.Nil, fmt.Errorf("VEC3.DIST: components must be numeric")
		}
		dx := float64(x2 - x1)
		dy := float64(y2 - y1)
		dz := float64(z2 - z1)
		return value.FromFloat(math.Sqrt(dx*dx + dy*dy + dz*dz)), nil
	}
	return m.vec3Distance(args)
}

func (m *Module) vec3DistSq(args []value.Value) (value.Value, error) {
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("VEC3.DISTSQ expects (x1, y1, z1, x2, y2, z2)")
	}
	x1, ok1 := argF(args[0])
	y1, ok2 := argF(args[1])
	z1, ok3 := argF(args[2])
	x2, ok4 := argF(args[3])
	y2, ok5 := argF(args[4])
	z2, ok6 := argF(args[5])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("VEC3.DISTSQ: components must be numeric")
	}
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	dz := float64(z2 - z1)
	return value.FromFloat(dx*dx + dy*dy + dz*dz), nil
}

func (m *Module) vec3Distance(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC3.DISTANCE expects two vec3 handles")
	}
	a, err := m.vec3FromArgs(args, 0, "VEC3.DISTANCE")
	if err != nil {
		return value.Nil, err
	}
	b, err := m.vec3FromArgs(args, 1, "VEC3.DISTANCE")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(rl.Vector3Distance(a, b))), nil
}

func (m *Module) vec3Reflect(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC3.REFLECT expects (v, normal)")
	}
	v, err := m.vec3FromArgs(args, 0, "VEC3.REFLECT")
	if err != nil {
		return value.Nil, err
	}
	n, err := m.vec3FromArgs(args, 1, "VEC3.REFLECT")
	if err != nil {
		return value.Nil, err
	}
	return m.allocVec3(rl.Vector3Reflect(v, n))
}

func (m *Module) vec3Negate(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("VEC3.NEGATE expects vec3 handle")
	}
	v, err := m.vec3FromArgs(args, 0, "VEC3.NEGATE")
	if err != nil {
		return value.Nil, err
	}
	return m.allocVec3(rl.Vector3Negate(v))
}

func (m *Module) vec3Equals(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC3.EQUALS expects two vec3 handles")
	}
	a, err := m.vec3FromArgs(args, 0, "VEC3.EQUALS")
	if err != nil {
		return value.Nil, err
	}
	b, err := m.vec3FromArgs(args, 1, "VEC3.EQUALS")
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(rl.Vector3Equals(a, b)), nil
}
