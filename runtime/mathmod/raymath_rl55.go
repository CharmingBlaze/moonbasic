package mathmod

import "math"

// Raylib 5.5 raymath (float64) — transcribed from github.com/gen2brain/raylib-go/raylib raymath.go
// (Raylib 5.5 rlClamp, rlLerp, rlWrap).
//
// We do not import the raylib package here: rl's init() calls runtime.LockOSThread() for the
// graphics runtime; moonBASIC math builtins must run without that side effect in CLI-only use.

func clampRL55(value, min, max float64) float64 {
	var res float64
	if value < min {
		res = min
	} else {
		res = value
	}
	if res > max {
		return max
	}
	return res
}

func lerpRL55(start, end, amount float64) float64 {
	return start + amount*(end-start)
}

// wrapRL55 matches rlWrap. If max == min, returns min (raylib would divide by zero; BASIC guards).
func wrapRL55(value, min, max float64) float64 {
	if max == min {
		return min
	}
	return value - (max-min)*math.Floor((value-min)/(max-min))
}
