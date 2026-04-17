//go:build cgo || (windows && !cgo)

package mbmodel3d

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerModelTransform(m *Module, reg runtime.Registrar) {
	reg.Register("MODEL.SETPOS", "model", runtime.AdaptLegacy(m.modelSetPos))
	reg.Register("MODEL.SETPOSITION", "model", runtime.AdaptLegacy(m.modelSetPos))
}

// modelSetPos sets the model's root transform to a translation matrix.
// This replaces the full transform (rotation/scale from the previous matrix are cleared).
func (m *Module) modelSetPos(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("MODEL.SETPOS expects (model, x, y, z)")
	}
	x, ok1 := argFloat(args[1])
	y, ok2 := argFloat(args[2])
	z, ok3 := argFloat(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("MODEL.SETPOS: x, y, z must be numeric")
	}
	if mo, err := m.getModel(args, 0, "MODEL.SETPOS"); err == nil {
		mo.model.Transform = rl.MatrixTranslate(x, y, z)
		return args[0], nil
	}
	if lo, err := m.getLODModel(args, 0, "MODEL.SETPOS"); err == nil {
		lo.transform = rl.MatrixTranslate(x, y, z)
		return args[0], nil
	}
	return value.Nil, fmt.Errorf("MODEL.SETPOS: handle must be a Model or LODModel")
}
