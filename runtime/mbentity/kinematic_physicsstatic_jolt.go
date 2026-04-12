//go:build linux && cgo

package mbentity

// Jolt-linked physics3d uses a different body3dObj (no stub Pos/Shape). Host kinematic AABB helpers
// iterate scene entities only here; use Jolt Character / picks for physics-static geometry.

func foreachPhysicsStaticBox(resolveBox func(bx, by, bz, hw, hh, hd float32)) {
	_ = resolveBox
}

func foreachPhysicsStaticFloor(checkFloor func(bx, by, bz, bw, bh, bd float64)) {
	_ = checkFloor
}
