package mbgame

import (
	"math"
	"strings"
)

func easeInQuad(t float64) float64  { return t * t }
func easeOutQuad(t float64) float64 { return t * (2 - t) }
func easeInOutQuad(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return -1 + (4-2*t)*t
}

func easeInCubic(t float64) float64  { return t * t * t }
func easeOutCubic(t float64) float64 {
	t = 1 - t
	return 1 - t*t*t
}
func easeInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	return 1 - math.Pow(-2*t+2, 3)/2
}

func easeInSine(t float64) float64  { return 1 - math.Cos(t*math.Pi/2) }
func easeOutSine(t float64) float64  { return math.Sin(t * math.Pi / 2) }
func easeInOutSine(t float64) float64 { return 0.5 * (1 - math.Cos(math.Pi*t)) }

func easeInBack(t float64) float64 {
	const c1 = 1.70158
	const c3 = c1 + 1
	return c3*t*t*t - c1*t*t
}

func easeOutBack(t float64) float64 {
	const c1 = 1.70158
	const c3 = c1 + 1
	t--
	return 1 + c3*math.Pow(t, 3) + c1*math.Pow(t, 2)
}

func easeInBounce(t float64) float64 { return 1 - easeOutBounce(1-t) }

func easeOutBounce(t float64) float64 {
	const n1 = 7.5625
	const d1 = 2.75
	if t < 1/d1 {
		return n1 * t * t
	}
	if t < 2/d1 {
		t -= 1.5 / d1
		return n1*t*t + 0.75
	}
	if t < 2.5/d1 {
		t -= 2.25 / d1
		return n1*t*t + 0.9375
	}
	t -= 2.625 / d1
	return n1*t*t + 0.984375
}

func easeInElastic(t float64) float64 {
	if t == 0 || t == 1 {
		return t
	}
	const c4 = (2 * math.Pi) / 3
	return -math.Pow(2, 10*t-10) * math.Sin((t*10-10.75)*c4)
}

func easeOutElastic(t float64) float64 {
	if t == 0 || t == 1 {
		return t
	}
	const c4 = (2 * math.Pi) / 3
	return math.Pow(2, -10*t)*math.Sin((t*10-0.75)*c4) + 1
}

// easeLerp applies named easing to t in [0,1], then lerps a..b.
func easeLerp(a, b, t float64, name string) float64 {
	if t <= 0 {
		return a
	}
	if t >= 1 {
		return b
	}
	var u float64
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "easein", "in":
		u = easeInQuad(t)
	case "easeout", "out":
		u = easeOutQuad(t)
	case "easeinout", "inout":
		u = easeInOutQuad(t)
	case "easein3", "in3":
		u = easeInCubic(t)
	case "easeout3", "out3":
		u = easeOutCubic(t)
	case "easeinout3", "inout3":
		u = easeInOutCubic(t)
	case "easeinsine", "insine":
		u = easeInSine(t)
	case "easeoutsine", "outsine":
		u = easeOutSine(t)
	case "easeinoutsine", "inoutsine":
		u = easeInOutSine(t)
	case "easeinback", "inback":
		u = easeInBack(t)
	case "easeoutback", "outback":
		u = easeOutBack(t)
	case "easeinbounce", "inbounce":
		u = easeInBounce(t)
	case "easeoutbounce", "outbounce":
		u = easeOutBounce(t)
	case "easeinelastic", "inelastic":
		u = easeInElastic(t)
	case "easeoutelastic", "outelastic":
		u = easeOutElastic(t)
	default:
		u = easeOutQuad(t)
	}
	return a + (b-a)*u
}
