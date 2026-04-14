//go:build (!linux && !windows) || !cgo
//
// Non-Jolt stub: same API as pick_cgo.go; raycasts return no hit. Used when (!linux && !windows) || !cgo.
// fullruntime builds working. See AGENTS.md “Physics sync & Jolt”.
//
package mbphysics3d

// SetPickLayerLookup is a no-op without Jolt.
func SetPickLayerLookup(fn func(int64) (uint8, bool)) { _ = fn }

func registerPickCommands(m *Module, reg runtime.Registrar) {
	reg.Register("PICK.ORIGIN", "physics3d", runtime.AdaptLegacy(m.pickOrigin))
	reg.Register("PICK.DIRECTION", "physics3d", runtime.AdaptLegacy(m.pickDirection))
	reg.Register("PICK.MAXDIST", "physics3d", runtime.AdaptLegacy(m.pickMaxDistSet))
	reg.Register("PICK.LAYERMASK", "physics3d", runtime.AdaptLegacy(m.pickLayerMaskSet))
	reg.Register("PICK.RADIUS", "physics3d", runtime.AdaptLegacy(m.pickRadiusSet))
	reg.Register("PICK.CAST", "physics3d", runtime.AdaptLegacy(m.pickCast))
	reg.Register("PICK.FROMCAMERA", "physics3d", runtime.AdaptLegacy(m.pickFromCamera))
	reg.Register("PICK.SCREENCAST", "physics3d", runtime.AdaptLegacy(m.pickScreenCast))
	reg.Register("PICK.X", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return m.pickGet(a, func() float64 { return 0 }) }))
	reg.Register("PICK.Y", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return m.pickGet(a, func() float64 { return 0 }) }))
	reg.Register("PICK.Z", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return m.pickGet(a, func() float64 { return 0 }) }))
	reg.Register("PICK.NX", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return m.pickGet(a, func() float64 { return 0 }) }))
	reg.Register("PICK.NY", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return m.pickGet(a, func() float64 { return 0 }) }))
	reg.Register("PICK.NZ", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return m.pickGet(a, func() float64 { return 0 }) }))
	reg.Register("PICK.ENTITY", "physics3d", runtime.AdaptLegacy(m.pickEntityGet))
	reg.Register("PICK.DIST", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return m.pickGet(a, func() float64 { return 0 }) }))
	reg.Register("PICK.HIT", "physics3d", runtime.AdaptLegacy(m.pickHitGet))
	reg.Register("PHYSICS3D.MOUSEHIT", "physics3d", runtime.AdaptLegacy(m.phMouseHit))
	reg.Register("WORLD.MOUSETOENTITY", "physics3d", runtime.AdaptLegacy(m.camRaycastMouseEntity))
	reg.Register("WORLD.MOUSEPICK", "physics3d", runtime.AdaptLegacy(m.camRaycastMouseEntity))
	reg.Register("CAMERA.RAYCASTMOUSE", "camera", runtime.AdaptLegacy(m.camRaycastMouseEntity))
}

func (m *Module) pickOrigin(args []value.Value) (value.Value, error)      { return value.Nil, nil }
func (m *Module) pickDirection(args []value.Value) (value.Value, error)   { return value.Nil, nil }
func (m *Module) pickMaxDistSet(args []value.Value) (value.Value, error)  { return value.Nil, nil }
func (m *Module) pickLayerMaskSet(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) pickRadiusSet(args []value.Value) (value.Value, error)    { return value.Nil, nil }
func (m *Module) pickFromCamera(args []value.Value) (value.Value, error)   { return value.Nil, nil }
func (m *Module) pickScreenCast(args []value.Value) (value.Value, error)   { return value.Nil, nil }
func (m *Module) camRaycastMouseEntity(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) phMouseHit(args []value.Value) (value.Value, error)       { return value.Nil, nil }
func (m *Module) pickCast(args []value.Value) (value.Value, error)         { return value.FromInt(0), nil }
func (m *Module) pickGet(args []value.Value, f func() float64) (value.Value, error) {
	return value.FromFloat(0), nil
}
func (m *Module) pickEntityGet(args []value.Value) (value.Value, error) { return value.FromInt(0), nil }
func (m *Module) pickHitGet(args []value.Value) (value.Value, error)    { return value.False, nil }

func resetPickState() {}

// PickCastEntityID is unavailable without Jolt.
func PickCastEntityID(ox, oy, oz, dx, dy, dz, maxDist float64) int64 {
	_, _, _, _, _, _, _ = ox, oy, oz, dx, dy, dz, maxDist
	return 0
}

var groundRaycastHook func(ox, oy, oz, maxDown float64) (nx, ny, nz, hitY float64, ok bool)

func SetRaycastHook(fn func(ox, oy, oz, maxDown float64) (nx, ny, nz, hitY float64, ok bool)) {
	groundRaycastHook = fn
}

// RaycastDownGroundProbe is unavailable without Jolt (unless hook is wired).
func RaycastDownGroundProbe(ox, oy, oz, maxDown float64) (nx, ny, nz, hitY float64, ok bool) {
	if groundRaycastHook != nil {
		return groundRaycastHook(ox, oy, oz, maxDown)
	}
	return 0, 1, 0, 0, false
}

// RaycastDownNormal is unavailable without Jolt.
func RaycastDownNormal(ox, oy, oz, maxDown float64) (nx, ny, nz float64, ok bool) {
	nx, ny, nz, _, ok = RaycastDownGroundProbe(ox, oy, oz, maxDown)
	return nx, ny, nz, ok
}
func (m *Module) pickRadiusSet(args []value.Value) (value.Value, error)    { return value.Nil, nil }
func (m *Module) pickFromCamera(args []value.Value) (value.Value, error)   { return value.Nil, nil }
func (m *Module) pickScreenCast(args []value.Value) (value.Value, error)   { return value.Nil, nil }
func (m *Module) camRaycastMouseEntity(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) phMouseHit(args []value.Value) (value.Value, error)       { return value.Nil, nil }
func (m *Module) pickCast(args []value.Value) (value.Value, error)         { return value.FromInt(0), nil }
func (m *Module) pickGet(args []value.Value, f func() float64) (value.Value, error) {
	return value.FromFloat(0), nil
}
func (m *Module) pickEntityGet(args []value.Value) (value.Value, error) { return value.FromInt(0), nil }
func (m *Module) pickHitGet(args []value.Value) (value.Value, error)    { return value.False, nil }
