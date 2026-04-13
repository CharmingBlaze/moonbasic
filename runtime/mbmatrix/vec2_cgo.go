//go:build cgo || (windows && !cgo)

package mbmatrix

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/hal"
	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerVec2(reg runtime.Registrar) {
	reg.Register("VEC2.MAKE", "vec2", runtime.AdaptLegacy(m.vec2Make))
	reg.Register("VEC2.FREE", "vec2", runtime.AdaptLegacy(m.vec2Free))
	reg.Register("VEC2.X", "vec2", runtime.AdaptLegacy(m.vec2X))
	reg.Register("VEC2.Y", "vec2", runtime.AdaptLegacy(m.vec2Y))
	reg.Register("VEC2.SET", "vec2", runtime.AdaptLegacy(m.vec2Set))
	reg.Register("VEC2.ADD", "vec2", runtime.AdaptLegacy(m.vec2Add))
	reg.Register("VEC2.SUB", "vec2", runtime.AdaptLegacy(m.vec2Sub))
	reg.Register("VEC2.MUL", "vec2", runtime.AdaptLegacy(m.vec2Mul))
	reg.Register("VEC2.LENGTH", "vec2", runtime.AdaptLegacy(m.vec2Length))
	reg.Register("VEC2.NORMALIZE", "vec2", runtime.AdaptLegacy(m.vec2Normalize))
	reg.Register("VEC2.MOVE_TOWARD", "vec2", runtime.AdaptLegacy(m.vec2MoveToward))
	reg.Register("VEC2.DIST", "vec2", runtime.AdaptLegacy(m.vec2Dist))
	reg.Register("VEC2.DISTSQ", "vec2", runtime.AdaptLegacy(m.vec2DistSq))
	reg.Register("VEC2.PUSHOUT", "vec2", runtime.AdaptLegacy(m.vec2PushOut))
	reg.Register("VEC2.LERP", "vec2", runtime.AdaptLegacy(m.vec2Lerp))
	reg.Register("VEC2.DISTANCE", "vec2", runtime.AdaptLegacy(m.vec2Distance))
	reg.Register("VEC2.ANGLE", "vec2", runtime.AdaptLegacy(m.vec2Angle))
	reg.Register("VEC2.ROTATE", "vec2", runtime.AdaptLegacy(m.vec2Rotate))
}

func (m *Module) vec2Make(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC2.MAKE expects 2 arguments (x, y)")
	}
	x, ok1 := argF(args[0])
	y, ok2 := argF(args[1])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("VEC2.MAKE: components must be numeric")
	}
	return m.allocVec2(hal.V2{X: x, Y: y})
}

func (m *Module) vec2Free(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("VEC2.FREE expects vec2 handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) vec2X(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("VEC2.X expects vec2 handle")
	}
	v, err := m.vec2FromArgs(args, 0, "VEC2.X")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(v.X)), nil
}

func (m *Module) vec2Y(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("VEC2.Y expects vec2 handle")
	}
	v, err := m.vec2FromArgs(args, 0, "VEC2.Y")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(v.Y)), nil
}

func (m *Module) vec2Set(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("VEC2.SET expects (handle, x, y)")
	}
	o, err := heap.Cast[*vec2Obj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("VEC2.SET: %w", err)
	}
	x, ok1 := argF(args[1])
	y, ok2 := argF(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("VEC2.SET: components must be numeric")
	}
	o.v = hal.V2{X: x, Y: y}
	return value.Nil, nil
}

func (m *Module) vec2Add(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC2.ADD expects two vec2 handles")
	}
	a, err := m.vec2FromArgs(args, 0, "VEC2.ADD")
	if err != nil {
		return value.Nil, err
	}
	b, err := m.vec2FromArgs(args, 1, "VEC2.ADD")
	if err != nil {
		return value.Nil, err
	}
	return m.allocVec2(fromV2(rl.Vector2Add(toV2(a), toV2(b))))
}

func (m *Module) vec2Sub(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC2.SUB expects two vec2 handles")
	}
	a, err := m.vec2FromArgs(args, 0, "VEC2.SUB")
	if err != nil {
		return value.Nil, err
	}
	b, err := m.vec2FromArgs(args, 1, "VEC2.SUB")
	if err != nil {
		return value.Nil, err
	}
	return m.allocVec2(fromV2(rl.Vector2Subtract(toV2(a), toV2(b))))
}

func (m *Module) vec2Mul(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC2.MUL expects (vec2, scalar)")
	}
	v, err := m.vec2FromArgs(args, 0, "VEC2.MUL")
	if err != nil {
		return value.Nil, err
	}
	s, ok := argF(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("VEC2.MUL: scalar must be numeric")
	}
	return m.allocVec2(fromV2(rl.Vector2Scale(toV2(v), s)))
}

func (m *Module) vec2Length(args []value.Value) (value.Value, error) {
	if len(args) == 2 {
		x, ok1 := argF(args[0])
		y, ok2 := argF(args[1])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("VEC2.LENGTH: components must be numeric")
		}
		return value.FromFloat(math.Hypot(float64(x), float64(y))), nil
	}
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("VEC2.LENGTH expects vec2 handle or (x, y)")
	}
	v, err := m.vec2FromArgs(args, 0, "VEC2.LENGTH")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(rl.Vector2Length(toV2(v)))), nil
}

func (m *Module) vec2Normalize(args []value.Value) (value.Value, error) {
	if len(args) == 2 {
		x, ok1 := argF(args[0])
		y, ok2 := argF(args[1])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("VEC2.NORMALIZE: components must be numeric")
		}
		mag := float32(math.Hypot(float64(x), float64(y)))
		if mag > 0 {
			x /= mag
			y /= mag
		}
		return m.allocTuple2(x, y)
	}
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("VEC2.NORMALIZE expects vec2 handle or (x, y)")
	}
	v, err := m.vec2FromArgs(args, 0, "VEC2.NORMALIZE")
	if err != nil {
		return value.Nil, err
	}
	return m.allocVec2(fromV2(rl.Vector2Normalize(toV2(v))))
}

func (m *Module) vec2MoveToward(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("VEC2.MOVE_TOWARD expects (fromX, fromY, toX, toY, maxDist)")
	}
	fx, ok1 := argF(args[0])
	fy, ok2 := argF(args[1])
	tx, ok3 := argF(args[2])
	ty, ok4 := argF(args[3])
	maxDist, ok5 := argF(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("VEC2.MOVE_TOWARD: arguments must be numeric")
	}
	if maxDist <= 0 {
		return m.allocTuple2(fx, fy)
	}
	dx := tx - fx
	dy := ty - fy
	dist := float32(math.Hypot(float64(dx), float64(dy)))
	if dist <= maxDist || dist <= 1e-6 {
		return m.allocTuple2(tx, ty)
	}
	t := maxDist / dist
	return m.allocTuple2(fx+dx*t, fy+dy*t)
}

func (m *Module) allocTuple2(x, y float32) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	arr, err := heap.NewArrayOfKind([]int64{2}, heap.ArrayKindFloat, 0)
	if err != nil {
		return value.Nil, err
	}
	arr.Floats[0] = float64(x)
	arr.Floats[1] = float64(y)
	h, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func (m *Module) vec2Lerp(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("VEC2.LERP expects (a, b, t)")
	}
	a, err := m.vec2FromArgs(args, 0, "VEC2.LERP")
	if err != nil {
		return value.Nil, err
	}
	b, err := m.vec2FromArgs(args, 1, "VEC2.LERP")
	if err != nil {
		return value.Nil, err
	}
	t, ok := argF(args[2])
	if !ok {
		return value.Nil, fmt.Errorf("VEC2.LERP: t must be numeric")
	}
	return m.allocVec2(fromV2(rl.Vector2Lerp(toV2(a), toV2(b), t)))
}

func (m *Module) vec2Dist(args []value.Value) (value.Value, error) {
	if len(args) == 4 {
		x1, ok1 := argF(args[0])
		y1, ok2 := argF(args[1])
		x2, ok3 := argF(args[2])
		y2, ok4 := argF(args[3])
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return value.Nil, fmt.Errorf("VEC2.DIST: components must be numeric")
		}
		dx := float64(x2 - x1)
		dy := float64(y2 - y1)
		return value.FromFloat(math.Hypot(dx, dy)), nil
	}
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC2.DIST expects (x1,y1,x2,y2) or two vec2 handles")
	}
	a, err := m.vec2FromArgs(args, 0, "VEC2.DIST")
	if err != nil {
		return value.Nil, err
	}
	b, err := m.vec2FromArgs(args, 1, "VEC2.DIST")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(rl.Vector2Distance(toV2(a), toV2(b)))), nil
}

func (m *Module) vec2DistSq(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("VEC2.DISTSQ expects (x1, y1, x2, y2)")
	}
	x1, ok1 := argF(args[0])
	y1, ok2 := argF(args[1])
	x2, ok3 := argF(args[2])
	y2, ok4 := argF(args[3])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("VEC2.DISTSQ: components must be numeric")
	}
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	return value.FromFloat(dx*dx + dy*dy), nil
}

func (m *Module) vec2PushOut(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("VEC2.PUSHOUT expects (x, z, cx, cz, minRadius)")
	}
	x, ok1 := argF(args[0])
	z, ok2 := argF(args[1])
	cx, ok3 := argF(args[2])
	cz, ok4 := argF(args[3])
	minR, ok5 := argF(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("VEC2.PUSHOUT: arguments must be numeric")
	}
	if minR <= 0 {
		return m.allocTuple2(x, z)
	}
	dx := x - cx
	dz := z - cz
	dist := float32(math.Hypot(float64(dx), float64(dz)))
	if dist >= minR {
		return m.allocTuple2(x, z)
	}
	if dist < 1e-6 {
		return m.allocTuple2(cx+minR, cz)
	}
	push := minR - dist
	t := push / dist
	return m.allocTuple2(x+dx*t, z+dz*t)
}

func (m *Module) vec2Distance(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC2.DISTANCE expects two vec2 handles")
	}
	a, err := m.vec2FromArgs(args, 0, "VEC2.DISTANCE")
	if err != nil {
		return value.Nil, err
	}
	b, err := m.vec2FromArgs(args, 1, "VEC2.DISTANCE")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(rl.Vector2Distance(toV2(a), toV2(b)))), nil
}

func (m *Module) vec2Angle(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC2.ANGLE expects two vec2 handles")
	}
	a, err := m.vec2FromArgs(args, 0, "VEC2.ANGLE")
	if err != nil {
		return value.Nil, err
	}
	b, err := m.vec2FromArgs(args, 1, "VEC2.ANGLE")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(rl.Vector2Angle(toV2(a), toV2(b)))), nil
}

func (m *Module) vec2Rotate(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC2.ROTATE expects (vec2, angleRadians)")
	}
	v, err := m.vec2FromArgs(args, 0, "VEC2.ROTATE")
	if err != nil {
		return value.Nil, err
	}
	ang, ok := argF(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("VEC2.ROTATE: angle must be numeric")
	}
	return m.allocVec2(fromV2(rl.Vector2Rotate(toV2(v), ang)))
}
