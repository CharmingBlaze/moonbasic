//go:build cgo

package mbmodel3d

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerModelLoad(m *Module, reg runtime.Registrar) {
	reg.Register("MODEL.LOAD", "model", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 || args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("MODEL.LOAD expects 1 string path")
		}
		path, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		mod := rl.LoadModel(path)
		id, err := m.h.Alloc(&modelObj{model: mod, loadedPath: path})
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	})

	reg.Register("MODEL.FREE", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 || args[0].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("MODEL.FREE expects model handle")
		}
		if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
			return value.Nil, err
		}
		return value.Nil, nil
	}))

	reg.Register("MODEL.GETMATERIALCOUNT", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MODEL.GETMATERIALCOUNT expects model handle")
		}
		o, err := m.getModel(args, 0, "MODEL.GETMATERIALCOUNT")
		if err != nil {
			return value.Nil, err
		}
		return value.FromInt(int64(o.model.MaterialCount)), nil
	}))
}
