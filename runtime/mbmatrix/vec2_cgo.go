//go:build cgo || (windows && !cgo)

package mbmatrix

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

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
	reg.Register("VEC2.LERP", "vec2", runtime.AdaptLegacy(m.vec2Lerp))
	reg.Register("VEC2.DISTANCE", "vec2", runtime.AdaptLegacy(m.vec2Distance))
	reg.Register("VEC2.ANGLE", "vec2", runtime.AdaptLegacy(m.vec2Angle))
	reg.Register("VEC2.ROTATE", "vec2", runtime.AdaptLegacy(m.vec2Rotate))
}

func (m *Module) allocVec2(v rl.Vector2) (value.Value, error) {
	id, err := m.h.Alloc(&vec2Obj{v: v})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
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
	return m.allocVec2(rl.Vector2{X: x, Y: y})
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
	o.v = rl.Vector2{X: x, Y: y}
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
	return m.allocVec2(rl.Vector2Add(a, b))
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
	return m.allocVec2(rl.Vector2Subtract(a, b))
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
	return m.allocVec2(rl.Vector2Scale(v, s))
}

func (m *Module) vec2Length(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("VEC2.LENGTH expects vec2 handle")
	}
	v, err := m.vec2FromArgs(args, 0, "VEC2.LENGTH")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(rl.Vector2Length(v))), nil
}

func (m *Module) vec2Normalize(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("VEC2.NORMALIZE expects vec2 handle")
	}
	v, err := m.vec2FromArgs(args, 0, "VEC2.NORMALIZE")
	if err != nil {
		return value.Nil, err
	}
	return m.allocVec2(rl.Vector2Normalize(v))
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
	return m.allocVec2(rl.Vector2Lerp(a, b, t))
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
	return value.FromFloat(float64(rl.Vector2Distance(a, b))), nil
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
	return value.FromFloat(float64(rl.Vector2Angle(a, b))), nil
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
	return m.allocVec2(rl.Vector2Rotate(v, ang))
}
