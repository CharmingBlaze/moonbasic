package mathmod

import (
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// registerGamePlaneHelpers registers 2D/3D gameplay-oriented math (no Raylib).
func (m *Module) registerGamePlaneHelpers(r runtime.Registrar) {
	regFlat := func(short, long string, fn runtime.BuiltinFn) {
		r.Register(short, "math", fn)
		r.Register(long, "math", fn)
	}

	twoFloat := func(f func(float64, float64) float64) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			if len(args) != 2 {
				return value.Nil, errNArgs(2, len(args))
			}
			a, _ := args[0].ToFloat()
			b, _ := args[1].ToFloat()
			return value.FromFloat(f(a, b)), nil
		}
	}
	threeFloat := func(f func(float64, float64, float64) float64) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			if len(args) != 3 {
				return value.Nil, errNArgs(3, len(args))
			}
			a, _ := args[0].ToFloat()
			b, _ := args[1].ToFloat()
			c, _ := args[2].ToFloat()
			return value.FromFloat(f(a, b, c)), nil
		}
	}
	fourFloat := func(f func(float64, float64, float64, float64) float64) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			if len(args) != 4 {
				return value.Nil, errNArgs(4, len(args))
			}
			a, _ := args[0].ToFloat()
			b, _ := args[1].ToFloat()
			c, _ := args[2].ToFloat()
			d, _ := args[3].ToFloat()
			return value.FromFloat(f(a, b, c, d)), nil
		}
	}

	// Ken Perlin's smootherstep — useful for cameras, fades, procedural feel.
	regFlat("SMOOTHERSTEP", "MATH.SMOOTHERSTEP", threeFloat(smootherstepRL))

	// Horizontal distance on the XZ plane (ignores Y). Common for 3D chase / aggro radii.
	regFlat("HDIST", "MATH.HDIST", fourFloat(hdistXZ))
	regFlat("HDISTSQ", "MATH.HDISTSQ", fourFloat(hdistsqXZ))

	// Yaw (radians) that aligns +forward in XZ with (dx, dz), matching INPUT.MOVEDIR / MOVEX-MOVEZ convention.
	regFlat("YAWFROMXZ", "MATH.YAWFROMXZ", twoFloat(yawFromXZ))

	// Heading in degrees [0,360) on the XZ plane from (x1,z1) toward (x2,z2), 0° = +Z, increasing clockwise when +X is right (matches YAWFROMXZ / MOVEX convention).
	regFlat("ANGLETO", "MATH.ANGLETO", fourFloat(angleToDegXZ))

	// Shortest signed angle difference in radians (for AI steering, compare to LERPANGLE which interpolates).
	regFlat("ANGLEDIFFRAD", "MATH.ANGLEDIFFRAD", twoFloat(angleDiffRad))

	// 2D distance / squared — same as DISTANCE2D / DISTANCESQ2D in mbgame; exposed under MATH for discovery.
	regFlat("DIST2D", "MATH.DIST2D", fourFloat(dist2D))
	regFlat("DISTSQ2D", "MATH.DISTSQ2D", fourFloat(distSq2D))
}

func smootherstepRL(edge0, edge1, x float64) float64 {
	if edge1 == edge0 {
		return 0
	}
	t := (x - edge0) / (edge1 - edge0)
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	return t * t * t * (t*(t*6-15) + 10)
}

func hdistXZ(x1, z1, x2, z2 float64) float64 {
	dx := x2 - x1
	dz := z2 - z1
	return math.Hypot(dx, dz)
}

func hdistsqXZ(x1, z1, x2, z2 float64) float64 {
	dx := x2 - x1
	dz := z2 - z1
	return dx*dx + dz*dz
}

// yawFromXZ returns atan2(dx, dz): when (dx,dz) is unit forward, matches sin/cos split used by MOVEX/MOVEZ at that yaw.
func yawFromXZ(dx, dz float64) float64 {
	return math.Atan2(dx, dz)
}

func angleToDegXZ(x1, z1, x2, z2 float64) float64 {
	dx := x2 - x1
	dz := z2 - z1
	rad := math.Atan2(dx, dz)
	deg := rad * 180 / math.Pi
	if deg < 0 {
		deg += 360
	}
	return deg
}

func angleDiffRad(a, b float64) float64 {
	return math.Atan2(math.Sin(b-a), math.Cos(b-a))
}

func dist2D(x1, y1, x2, y2 float64) float64 {
	return math.Hypot(x2-x1, y2-y1)
}

func distSq2D(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return dx*dx + dy*dy
}
