package terrain

import (
	"math"
)

func smooth5(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

func hash01(perm *[512]int, x, y int) float64 {
	n := perm[(perm[x&255]+y)&255]
	return float64(n) / 255.0
}

// valueNoise2D smooth interpolated value noise in [0,1).
func valueNoise2D(x, y float64, perm *[512]int) float64 {
	x0 := int(math.Floor(x))
	y0 := int(math.Floor(y))
	xf := x - float64(x0)
	yf := y - float64(y0)
	x0 &= 255
	y0 &= 255
	x1 := (x0 + 1) & 255
	y1 := (y0 + 1) & 255
	u := smooth5(xf)
	v := smooth5(yf)
	n00 := hash01(perm, x0, y0)
	n10 := hash01(perm, x1, y0)
	n01 := hash01(perm, x0, y1)
	n11 := hash01(perm, x1, y1)
	nx0 := n00 + u*(n10-n00)
	nx1 := n01 + u*(n11-n01)
	return nx0 + v*(nx1-nx0)
}

func fbm2D(x, y, scale float64, perm *[512]int, octaves int) float64 {
	amp := 1.0
	freq := scale
	sum := 0.0
	norm := 0.0
	for o := 0; o < octaves; o++ {
		sum += amp * valueNoise2D(x*freq, y*freq, perm)
		norm += amp
		amp *= 0.5
		freq *= 2
	}
	if norm > 0 {
		return sum / norm
	}
	return 0
}

func makePerm(seed int64) [512]int {
	var p [256]int
	for i := range p {
		p[i] = i
	}
	s := uint64(seed)
	if s == 0 {
		s = 1
	}
	for i := 255; i > 0; i-- {
		s = s*6364136223846793005 + 1
		j := int(s % uint64(i+1))
		p[i], p[j] = p[j], p[i]
	}
	var out [512]int
	copy(out[:256], p[:])
	copy(out[256:], p[:])
	return out
}

func fillPerlinHeights(heights []float32, w, h int, scale float64, amp float32, seed int64) {
	perm := makePerm(seed)
	for z := 0; z < h; z++ {
		for x := 0; x < w; x++ {
			nx := float64(x) * 0.1
			nz := float64(z) * 0.1
			v := fbm2D(nx, nz, scale, &perm, 4)
			if v < 0 {
				v = 0
			}
			if v > 1 {
				v = 1
			}
			heights[z*w+x] = float32(v) * amp
		}
	}
}
