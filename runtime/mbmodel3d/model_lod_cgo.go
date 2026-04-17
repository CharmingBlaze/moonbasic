//go:build cgo || (windows && !cgo)

package mbmodel3d

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerModelLOD(m *Module, reg runtime.Registrar) {
	reg.Register("MODEL.LOADLOD", "model", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("MODEL.LOADLOD expects 3 string paths (high$, med$, low$)")
		}
		var paths [3]string
		for i := range paths {
			if args[i].Kind != value.KindString {
				return value.Nil, fmt.Errorf("MODEL.LOADLOD: argument %d must be string path", i+1)
			}
			p, err := rt.ArgString(args, i)
			if err != nil {
				return value.Nil, err
			}
			paths[i] = p
		}
		lo := &lodModelObj{transform: rl.MatrixIdentity()}
		for i := range paths {
			mod := rl.LoadModel(paths[i])
			if mod.MeshCount == 0 {
				for j := 0; j < i; j++ {
					rl.UnloadModel(lo.models[j])
					lo.models[j] = rl.Model{}
				}
				return value.Nil, fmt.Errorf("MODEL.LOADLOD: failed to load %q", paths[i])
			}
			lo.models[i] = mod
		}
		lo.setFinalizer()
		id, err := m.h.Alloc(lo)
		if err != nil {
			lo.Free()
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	})

	reg.Register("MODEL.SETLODDISTANCES", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("MODEL.SETLODDISTANCES expects (lodModel, band0, band1, band2)")
		}
		lo, err := m.getLODModel(args, 0, "MODEL.SETLODDISTANCES")
		if err != nil {
			return value.Nil, err
		}
		b0, ok1 := argFloat(args[1])
		b1, ok2 := argFloat(args[2])
		b2, ok3 := argFloat(args[3])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("MODEL.SETLODDISTANCES: distances must be numeric")
		}
		if b0 <= 0 || b1 <= b0 || b2 <= b1 {
			return value.Nil, fmt.Errorf("MODEL.SETLODDISTANCES: require 0 < band0 < band1 < band2")
		}
		lo.band0, lo.band1, lo.band2 = b0, b1, b2
		lo.configured = true
		return args[0], nil
	}))
}
