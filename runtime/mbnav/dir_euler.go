package mbnav

import "math"

// EulerFromWorldDirection returns pitch, yaw, roll (radians) suitable for ENTITY/BODY3D-style
// [pitch,yaw,roll] arrays, from a world-space direction (e.g. velocity or waypoint tangent).
// A zero-length vector returns (0,0,0).
func EulerFromWorldDirection(dx, dy, dz float64) (pitch, yaw, roll float64) {
	l := math.Sqrt(dx*dx + dy*dy + dz*dz)
	if l < 1e-10 {
		return 0, 0, 0
	}
	fx := dx / l
	fy := dy / l
	fz := dz / l
	pitch = math.Asin(-fy)
	yaw = math.Atan2(fx, fz)
	return pitch, yaw, 0
}
