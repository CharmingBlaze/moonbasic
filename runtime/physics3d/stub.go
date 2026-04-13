//go:build !linux || !cgo

package mbphysics3d

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// stubHint explains why native Jolt is unavailable on this build (see jolt-go cgo_* files: Linux/Darwin only today).
const stubHint = "native Jolt is not linked on this build (need Linux/macOS amd64/arm64 + CGO; github.com/bbitechnologies/jolt-go has no Windows CGO libs yet). Use Linux CI, WSL2 for dev parity, or contribute Windows static libs to jolt-go."

type Vec3 struct {
	X, Y, Z float32
}

type ShapeObj struct {
	Kind int
	F1, F2, F3 float32
}
type body3dObj struct {
	ID          int
	Pos, Rot    Vec3
	Shape       *ShapeObj
	Layer       int
	Collision   bool
}

var (
	staticBodies = make(map[heap.Handle]*body3dObj)
	nextBodyID   = 1
)

func registerPhysics3DCommands(m *Module, reg runtime.Registrar) {
	reg.Register("PHYSICS3D.GETSCRATCHFLOAT", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phGetScratchFloat(m, a) }))
	
	noop := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return value.Nil, nil }
	
	// Shape API
	reg.Register("SHAPE.CREATEBOX", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		hx, _ := args[0].ToFloat(); hy, _ := args[1].ToFloat(); hz, _ := args[2].ToFloat()
		id, _ := m.h.Alloc(&ShapeObj{1, float32(hx), float32(hy), float32(hz)})
		return value.FromHandle(id), nil
	})
	reg.Register("SHAPE.CREATESPHERE", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		r, _ := args[0].ToFloat()
		id, _ := m.h.Alloc(&ShapeObj{2, float32(r), 0, 0})
		return value.FromHandle(id), nil
	})
	reg.Register("SHAPE.CREATECAPSULE", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		r, _ := args[0].ToFloat(); h, _ := args[1].ToFloat()
		id, _ := m.h.Alloc(&ShapeObj{3, float32(r), float32(h), 0})
		return value.FromHandle(id), nil
	})
	reg.Register("SHAPE.CREATECYLINDER", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		r, _ := args[0].ToFloat(); h, _ := args[1].ToFloat()
		id, _ := m.h.Alloc(&ShapeObj{4, float32(r), float32(h), 0})
		return value.FromHandle(id), nil
	})
	reg.Register("SHAPEREF.FREE", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		m.h.Free(heap.Handle(args[0].IVal))
		return value.Nil, nil
	})

	// Body API
	reg.Register("STATIC.CREATE", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		sh, _ := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
		body := &body3dObj{ID: nextBodyID, Shape: sh, Collision: true}
		nextBodyID++
		id, _ := m.h.Alloc(body)
		staticBodies[id] = body
		return value.FromHandle(id), nil
	})
	reg.Register("KINEMATIC.CREATE", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		sh, _ := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
		body := &body3dObj{ID: nextBodyID, Shape: sh, Collision: true}
		nextBodyID++
		id, _ := m.h.Alloc(body)
		return value.FromHandle(id), nil
	})
	reg.Register("BODYREF.SETPOSITION", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		bo, _ := hGetBody(m, args[0])
		x, _ := args[1].ToFloat(); y, _ := args[2].ToFloat(); z, _ := args[3].ToFloat()
		bo.Pos = Vec3{X: float32(x), Y: float32(y), Z: float32(z)}
		return value.Nil, nil
	})
	reg.Register("BODYREF.SETROTATION", "physics3d", noop)
	reg.Register("BODYREF.FREE", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		bh := heap.Handle(args[0].IVal)
		delete(staticBodies, bh)
		m.h.Free(bh)
		return value.Nil, nil
	})

	reg.Register("PHYSICS3D.GETGRAVITYX", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		_ = args
		return value.FromFloat(0), nil
	})
	reg.Register("PHYSICS3D.GETGRAVITYY", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		_ = args
		return value.FromFloat(-9.81), nil
	})
	reg.Register("PHYSICS3D.GETGRAVITYZ", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		_ = args
		return value.FromFloat(0), nil
	})
	reg.Register("PHYSICS.GETGRAVITYX", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return value.FromFloat(0), nil
	})
	reg.Register("PHYSICS.GETGRAVITYY", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return value.FromFloat(-9.81), nil
	})
	reg.Register("PHYSICS.GETGRAVITYZ", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return value.FromFloat(0), nil
	})

	reg.Register("SHAPE.GETTYPE", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		sh, err := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
		if err != nil {
			return value.Nil, err
		}
		return value.FromInt(int64(sh.Kind)), nil
	})
	reg.Register("SHAPE.GETWIDTH", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		sh, err := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(float64(sh.F1)), nil
	})
	reg.Register("SHAPE.GETHEIGHT", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		sh, err := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(float64(sh.F2)), nil
	})
	reg.Register("SHAPE.GETDEPTH", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		sh, err := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(float64(sh.F3)), nil
	})
	reg.Register("SHAPE.GETRADIUS", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		sh, err := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(float64(sh.F1)), nil
	})
	reg.Register("SHAPE.GETSIZEX", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		sh, err := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(float64(sh.F1)), nil
	})
	reg.Register("SHAPE.GETSIZEY", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		sh, err := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(float64(sh.F2)), nil
	})
	reg.Register("SHAPE.GETSIZEZ", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		sh, err := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(float64(sh.F3)), nil
	})

	// Legacy / Other
	noopKeys := []string{
		"WORLD.SETGRAVITY", "PHYSICS3D.START", "PHYSICS3D.STOP", "PHYSICS3D.STEP", "PHYSICS3D.SETSUBSTEPS",
		"PHYSICS3D.SETTIMESTEP", "PHYSICS3D.SETGRAVITY",
		"PHYSICS.START", "PHYSICS.STOP", "PHYSICS.SETGRAVITY", "PHYSICS.STEP", "PHYSICS.SETSUBSTEPS",
		"KINEMATICREF.SETVELOCITY", "KINEMATICREF.UPDATE", "BODYREF.SETLAYER", "BODYREF.ENABLECOLLISION",
	}
	for _, k := range noopKeys {
		reg.Register(k, "physics3d", noop)
	}

	// World API (Easy Mode)
	reg.Register("WORLD.SETUP", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return PHWorldSetup(m, args)
	})
	reg.Register("LEVEL.STATIC", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return phLevelStatic(m, args)
	})
}

func hGetBody(m *Module, v value.Value) (*body3dObj, error) {
	return heap.Cast[*body3dObj](m.h, heap.Handle(v.IVal))
}

func (s *ShapeObj) TypeName() string { return "Shape" }
func (s *ShapeObj) TypeTag() uint16  { return heap.TagShape }
func (s *ShapeObj) Free()            {}

func (b *body3dObj) TypeName() string { return "Body3D" }
func (b *body3dObj) TypeTag() uint16 {
	return heap.TagPhysicsBody // Fallback sharing
}
func (b *body3dObj) Free() {}

func shutdownPhysics3D(m *Module) { _ = m }

// Exported for charcontroller/stub.go
func GetStaticBodyRegistry() map[heap.Handle]*body3dObj { return staticBodies }

type body3dObjExport struct {
	Pos   Vec3
	Shape *ShapeObj
}

func ApplyImpulseToIndex(idx int, x, y, z float32)     {}
func GetLinearVelocityToIndex(idx int) (x, y, z float32) { return 0, 0, 0 }
func GetBodyQuaternionForBufferIndex(idx int) (x, y, z, w float32, ok bool) { return 0, 0, 0, 1, false }
func SetVelocityToIndex(idx int, x, y, z float32)      {}
func SetPositionToIndex(idx int, x, y, z float32)      {}
func SetFrictionToIndex(idx int, x float32)            {}
func SetRestitutionToIndex(idx int, x float32)         {}
func WakeIndex(idx int)                                {}
func ApplyForceToIndex(idx int, x, y, z float32)       {}
func RotateToIndex(idx int, p, y, r float32)           {}
func SetGravityFactorToIndex(idx int, x float32)      {}

func PHWorldSetup(m *Module, args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func BDAddMesh(h *heap.Store, args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func BDCommit(h *heap.Store, args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func BDBufferIndex(h *heap.Store, args []value.Value) (value.Value, error) {
	return value.FromInt(0), nil
}

func phLevelStatic(m *Module, args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

type BuilderObj struct {
	Motion   int
	Friction float32
	Shape    *ShapeObj
}

func (b *BuilderObj) TypeName() string { return "Body3DBuilder" }
func (b *BuilderObj) TypeTag() uint16  { return heap.TagPhysicsBuilder }
func (b *BuilderObj) Free()            {}
