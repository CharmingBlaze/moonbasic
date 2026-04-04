package mbtween

import (
	"math"
	"strings"
)

func clamp01(t float64) float64 {
	if t < 0 {
		return 0
	}
	if t > 1 {
		return 1
	}
	return t
}

// ease applies named easing to t in [0,1].
func ease(t float64, name string) float64 {
	t = clamp01(t)
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "linear", "":
		return t
	case "easein":
		return easeInQuad(t)
	case "easeout":
		return easeOutQuad(t)
	case "easeinout":
		return easeInOutQuad(t)
	case "bounce":
		return easeOutBounce(t)
	case "elastic":
		return easeOutElastic(t)
	case "back":
		return easeOutBack(t)
	case "circ":
		return easeOutCirc(t)
	case "expo":
		return easeOutExpo(t)
	case "sine":
		return easeInOutSine(t)
	default:
		return t
	}
}

func easeInQuad(t float64) float64  { return t * t }
func easeOutQuad(t float64) float64 { return 1 - (1-t)*(1-t) }
func easeInOutQuad(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return 1 - math.Pow(-2*t+2, 2)/2
}

func easeInOutSine(t float64) float64 {
	return -(math.Cos(math.Pi*t) - 1) / 2
}

func easeOutCirc(t float64) float64 {
	return math.Sqrt(1 - math.Pow(t-1, 2))
}

func easeOutExpo(t float64) float64 {
	if t >= 1 {
		return 1
	}
	return 1 - math.Pow(2, -10*t)
}

func easeOutBack(t float64) float64 {
	const c1 = 1.70158
	const c3 = c1 + 1
	return 1 + c3*math.Pow(t-1, 3) + c1*math.Pow(t-1, 2)
}

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

func easeOutElastic(t float64) float64 {
	if t == 0 || t == 1 {
		return t
	}
	const c4 = (2 * math.Pi) / 3
	return math.Pow(2, -10*t)*math.Sin((t*10-0.75)*c4) + 1
}
