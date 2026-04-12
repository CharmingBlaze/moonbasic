//go:build !linux || !cgo

package mbentity

import mbphysics3d "moonbasic/runtime/physics3d"

// foreachPhysicsStaticBox applies resolveBox to heap static bodies (physics3d stub layout: Pos + Shape).
// Not used on linux+cgo — Jolt body3dObj is a different type (see kinematic_physicsstatic_jolt.go).
func foreachPhysicsStaticBox(resolveBox func(bx, by, bz, hw, hh, hd float32)) {
	for _, b := range mbphysics3d.GetStaticBodyRegistry() {
		if b == nil || b.Shape == nil {
			continue
		}
		if b.Shape.Kind == 1 {
			resolveBox(b.Pos.X, b.Pos.Y, b.Pos.Z, b.Shape.F1, b.Shape.F2, b.Shape.F3)
		}
	}
}

func foreachPhysicsStaticFloor(checkFloor func(bx, by, bz, bw, bh, bd float64)) {
	for _, b := range mbphysics3d.GetStaticBodyRegistry() {
		if b == nil || b.Shape == nil {
			continue
		}
		if b.Shape.Kind == 1 {
			checkFloor(float64(b.Pos.X), float64(b.Pos.Y), float64(b.Pos.Z), float64(b.Shape.F1*2), float64(b.Shape.F2*2), float64(b.Shape.F3*2))
		}
	}
}
