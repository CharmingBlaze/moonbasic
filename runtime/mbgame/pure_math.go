package mbgame

import (
	"math"
	"strings"
)

func deg2rad(d float64) float64 { return d * math.Pi / 180 }

// NEWXVALUE / NEWYVALUE: 2D step from (x,y) at angle (degrees) by dist.
func newXValue(x, angleDeg, dist float64) float64 {
	a := deg2rad(angleDeg)
	return x + math.Cos(a)*dist
}

func newYValue(y, angleDeg, dist float64) float64 {
	a := deg2rad(angleDeg)
	return y + math.Sin(a)*dist
}

// NEWZVALUE: vertical component for spherical-ish move (anglex/y in degrees).
func newZValue(z, angleXDeg, angleYDeg, dist float64) float64 {
	ax := deg2rad(angleXDeg)
	ay := deg2rad(angleYDeg)
	return z + dist*math.Sin(ay)*math.Cos(ax)
}

func pointDir2D(x1, y1, x2, y2 float64) float64 {
	return math.Atan2(y2-y1, x2-x1) * 180 / math.Pi
}

// POINTDIR3D: angle (degrees) on a chosen axis plane projection.
func pointDir3D(x1, y1, z1, x2, y2, z2 float64, axis string) float64 {
	switch strings.ToLower(strings.TrimSpace(axis)) {
	case "x":
		return math.Atan2(z2-z1, y2-y1) * 180 / math.Pi
	case "y":
		return math.Atan2(z2-z1, x2-x1) * 180 / math.Pi
	case "z":
		return math.Atan2(y2-y1, x2-x1) * 180 / math.Pi
	default:
		return pointDir2D(x1, y1, x2, y2)
	}
}

// CURVEVALUE: move src toward dest by (dest-src)/speed per call (DBPro-style).
func curveValue(dest, src, speed float64) float64 {
	if speed == 0 {
		return src
	}
	return src + (dest-src)/speed
}

// CURVEANGLE: shortest path on circle (degrees).
func curveAngle(destDeg, srcDeg, speed float64) float64 {
	if speed == 0 {
		return srcDeg
	}
	diff := math.Remainder(destDeg-srcDeg, 360)
	return srcDeg + diff/speed
}

func oscillate(elapsedSec, speed, minV, maxV float64) float64 {
	mid := (minV + maxV) / 2
	amp := (maxV - minV) / 2
	return mid + amp*math.Sin(elapsedSec*speed)
}

func wrapValue(v, minV, maxV float64) float64 {
	if maxV <= minV {
		return v
	}
	r := maxV - minV
	return minV + math.Mod(math.Mod(v-minV, r)+r, r)
}

func approach(cur, target, step float64) float64 {
	if step < 0 {
		step = -step
	}
	if cur < target {
		if cur+step > target {
			return target
		}
		return cur + step
	}
	if cur > target {
		if cur-step < target {
			return target
		}
		return cur - step
	}
	return target
}
