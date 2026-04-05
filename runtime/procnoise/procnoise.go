// Package procnoise implements deterministic 2D/3D procedural noise used by runtime/mbgame
// (legacy PERLIN/SIMPLEX/VORONOI) and runtime/noisemod (Noise.* handles). Pure Go, no allocations
// in sampling paths.
package procnoise

import (
	"math"
)

func hashInt2(x, y int32) int32 {
	n := x*374761393 + y*668265263
	n = (n ^ (n >> 13)) * 1274126177
	return n ^ (n >> 16)
}

func hashInt3(x, y, z int32) int32 {
	n := x*374761393 + y*668265263 + z*1442695041
	n = (n ^ (n >> 13)) * 1274126177
	return n ^ (n >> 16)
}

// HashInt2 is exported for HASHINT / mbgame helpers (same algorithm as legacy pure_noise).
func HashInt2(x, y int32) int32 { return hashInt2(x, y) }

// HashInt3 is exported for HASHINT (3-arg) parity.
func HashInt3(x, y, z int32) int32 { return hashInt3(x, y, z) }

func smooth5(t float64) float64 { return t * t * t * (t*(t*6-15) + 10) }

// ValueNoise2 returns approximately [-1,1] smooth value noise at (x,y), keyed by seed.
func ValueNoise2(x, y float64, seed int32) float64 {
	x0 := math.Floor(x)
	y0 := math.Floor(y)
	tx := x - x0
	ty := y - y0
	sx := int32(seed)*7919 + int32(x0)
	sy := int32(seed)*31337 + int32(y0)

	n00 := float64(hashInt2(sx, sy)%1024)/512 - 1
	n10 := float64(hashInt2(sx+1, sy)%1024)/512 - 1
	n01 := float64(hashInt2(sx, sy+1)%1024)/512 - 1
	n11 := float64(hashInt2(sx+1, sy+1)%1024)/512 - 1

	ax := smooth5(tx)
	ay := smooth5(ty)
	ix0 := n00 + ax*(n10-n00)
	ix1 := n01 + ax*(n11-n01)
	return ix0 + ay*(ix1-ix0)
}

// Perlin2 is gradient-style value noise (legacy name); matches previous mbgame perlin2 at seed 0.
func Perlin2(x, y float64, seed int32) float64 {
	return ValueNoise2(x, y, seed)
}

// Perlin3 blends three 2D planes for a cheap 3D-ish field (~[-1,1]).
func Perlin3(x, y, z float64, seed int32) float64 {
	return (ValueNoise2(x, y, seed) + ValueNoise2(y, z, seed+1) + ValueNoise2(z, x, seed+2)) / 3
}

// FBM2 fractal Brownian motion in 2D.
func FBM2(x, y float64, octaves int, lacunarity, gain float64, seed int32) float64 {
	if octaves < 1 {
		octaves = 1
	}
	if lacunarity <= 0 {
		lacunarity = 2
	}
	if gain <= 0 {
		gain = 0.5
	}
	amp := 0.5
	f := 1.0
	sum := 0.0
	norm := 0.0
	for i := 0; i < octaves; i++ {
		sum += amp * ValueNoise2(x*f, y*f, seed+int32(i)*17)
		norm += amp
		amp *= gain
		f *= lacunarity
	}
	if norm > 0 {
		return sum / norm
	}
	return 0
}

// FBM2Weighted applies octaves with optional high-octave emphasis (0 = standard).
func FBM2Weighted(x, y float64, octaves int, lacunarity, gain, weighted float64, seed int32) float64 {
	if octaves < 1 {
		octaves = 1
	}
	if lacunarity <= 0 {
		lacunarity = 2
	}
	if gain <= 0 {
		gain = 0.5
	}
	if weighted < 0 {
		weighted = 0
	}
	if weighted > 1 {
		weighted = 1
	}
	amp := 0.5
	f := 1.0
	sum := 0.0
	norm := 0.0
	for i := 0; i < octaves; i++ {
		v := ValueNoise2(x*f, y*f, seed+int32(i)*17)
		if i > 0 && weighted > 0 {
			v *= 1 + weighted*float64(i)/float64(octaves)
		}
		sum += amp * v
		norm += amp
		amp *= gain
		f *= lacunarity
	}
	if norm > 0 {
		return sum / norm
	}
	return 0
}

// RidgedMulti2 is a ridged multifractal built from value noise (~[-1,1]).
func RidgedMulti2(x, y float64, octaves int, lacunarity, gain float64, seed int32) float64 {
	if octaves < 1 {
		octaves = 1
	}
	if lacunarity <= 0 {
		lacunarity = 2
	}
	if gain <= 0 {
		gain = 0.5
	}
	sum := 0.0
	amp := 0.5
	f := 1.0
	weight := 1.0
	for i := 0; i < octaves; i++ {
		v := ValueNoise2(x*f, y*f, seed+int32(i)*19)
		v = 1.0 - math.Abs(v)
		v *= v
		v *= weight
		weight = math.Max(0.001, v)
		sum += v * amp
		amp *= gain
		f *= lacunarity
	}
	return math.Tanh(sum * 4)
}

// PingPongFBM2 applies a ping-pong shaping to normalized FBM output (~[-1,1]).
func PingPongFBM2(x, y float64, octaves int, lacunarity, gain, pingStrength float64, seed int32) float64 {
	t := FBM2Weighted(x, y, octaves, lacunarity, gain, 0, seed)
	u := (t*0.5 + 0.5) * pingStrength
	u = u - math.Floor(u)
	p := 2 * u
	if p > 1 {
		p = 2 - p
	}
	return p*2 - 1
}

// Simplex2 legacy stand-in: rotated value-noise domain (matches old mbgame).
func Simplex2(x, y float64, seed int32) float64 {
	return ValueNoise2(x*1.1-y*0.3, y*1.1+x*0.2, seed)
}

// Simplex3 legacy stand-in.
func Simplex3(x, y, z float64, seed int32) float64 {
	return (Simplex2(x, y, seed) + Simplex2(y, z, seed+3) + Simplex2(z, x, seed+5)) / 3
}

// VoronoiDistance returns distance to nearest hashed cell centre (unbounded); use Cellular2 for [-1,1].
func VoronoiDistance(x, y float64, seed int32) float64 {
	cx := math.Floor(x)
	cy := math.Floor(y)
	minD := math.MaxFloat64
	for dy := -1.0; dy <= 1; dy++ {
		for dx := -1.0; dx <= 1; dx++ {
			gx := cx + dx
			gy := cy + dy
			h := hashInt2(int32(gx)+seed*131, int32(gy)+seed*171)
			px := gx + float64(h%1000)/1000
			py := gy + float64((h/1000)%1000)/1000
			d := math.Hypot(x-px, y-py)
			if d < minD {
				minD = d
			}
		}
	}
	return minD
}

// Cellular2 maps Voronoi distance to ~[-1,1] (distance-style cellular).
func Cellular2(x, y float64, seed int32) float64 {
	d := VoronoiDistance(x, y, seed)
	return math.Tanh(d * 3)
}

// Warp2 returns (wx, wy) domain warp offsets from low-frequency noise.
func Warp2(x, y, amp float64, seed int32) (float64, float64) {
	wx := Simplex2(x*0.05+13, y*0.05+7, seed) * amp
	wy := Simplex2(x*0.05+101, y*0.05+17, seed+9) * amp
	return wx, wy
}

// Tileable2 maps (x,y) through a torus parameterization then samples fn (approximate seamless tiling).
func Tileable2(x, y, w, h float64, seed int32, fn func(x, y float64, seed int32) float64) float64 {
	if w <= 1e-9 || h <= 1e-9 {
		return fn(x, y, seed)
	}
	twopi := math.Pi * 2
	sx := math.Sin(x/w*twopi) * 512
	sy := math.Sin(y/h*twopi) * 512
	cx := math.Cos(x/w*twopi) * 512
	cy := math.Cos(y/h*twopi) * 512
	return fn(sx+cy, sy+cx, seed)
}
