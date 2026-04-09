//go:build cgo || (windows && !cgo)

package mbmodel3d

import (
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/mbjobs"
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
		obj := &modelObj{model: mod, loadedPath: path, animSpeed: 1}
		obj.setFinalizer()
		id, err := m.h.Alloc(obj)
		if err != nil {
			return value.Nil, err
		}
		obj.loaded = true // Synchronous load is ready immediately
		return value.FromHandle(id), nil
	})

	reg.Register("MODEL.LOADASYNC", "model", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 || args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("MODEL.LOADASYNC expects 1 string path")
		}
		path, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}

		obj := &modelObj{loadedPath: path, animSpeed: 1, isLoading: true}
		obj.setFinalizer()
		id, err := m.h.Alloc(obj)
		if err != nil {
			return value.Nil, err
		}

		// Background Task: File Check & Hand-off to Main Thread
		mbjobs.EnqueueJob(func() {
			// Simulate disk activity or check for file existence
			if _, err := os.Stat(path); err != nil {
				obj.mu.Lock()
				obj.loadError = err.Error()
				obj.isLoading = false
				obj.mu.Unlock()
				return
			}

			// Hand off to Main Thread for OpenGL calls
			enqueueOnMainThread(func() {
				mod := rl.LoadModel(path)
				obj.mu.Lock()
				obj.model = mod
				obj.loaded = true
				obj.isLoading = false
				obj.mu.Unlock()
			})
		})

		return value.FromHandle(id), nil
	})

	// MODEL.MAKE wraps LoadModelFromMesh — copies mesh data into a new Model (independent of the mesh handle).
	reg.Register("MODEL.MAKE", "model", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 || args[0].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("MODEL.MAKE expects mesh handle")
		}
		o, err := m.getMesh(args, 0, "MODEL.MAKE")
		if err != nil {
			return value.Nil, err
		}
		mod := rl.LoadModelFromMesh(o.m)
		o.consumedByModel = true
		obj := &modelObj{model: mod, loadedPath: "", animSpeed: 1}
		obj.setFinalizer()
		id, err := m.h.Alloc(obj)
		if err != nil {
			o.consumedByModel = false
			rl.UnloadModel(mod)
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	})

	freeModelOrInst := func(args []value.Value) (value.Value, error) {
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
	}
	reg.Register("MODEL.FREE", "model", runtime.AdaptLegacy(freeModelOrInst))
	reg.Register("INSTANCE.FREE", "model", runtime.AdaptLegacy(freeModelOrInst))

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

	reg.Register("MODEL.ISLOADED", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MODEL.ISLOADED expects model handle")
		}
		o, err := m.getModel(args, 0, "MODEL.ISLOADED")
		if err != nil {
			return value.Nil, err
		}
		o.mu.RLock()
		defer o.mu.RUnlock()
		return value.FromBool(o.loaded), nil
	}))
}
