//go:build cgo || (windows && !cgo)

package mbmatrix

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerVec2Extras(reg runtime.Registrar) {
	reg.Register("VEC2.TRANSFORMMAT4", "vec2", runtime.AdaptLegacy(m.vec2TransformMat4))
}

func (m *Module) vec2TransformMat4(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("VEC2.TRANSFORMMAT4 expects (vec2, mat4)")
	}
	v, err := m.vec2FromArgs(args, 0, "VEC2.TRANSFORMMAT4")
	if err != nil {
		return value.Nil, err
	}
	mat, err := m.matrixFromArgs(args, 1, "VEC2.TRANSFORMMAT4")
	if err != nil {
		return value.Nil, err
	}
	return m.allocVec2(fromV2(rl.Vector2Transform(toV2(v), toM(mat))))
}
