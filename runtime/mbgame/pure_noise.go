package mbgame

import (
	"math"

	"moonbasic/runtime/procnoise"
)

func hashInt2(x, y int32) int32 {
	return procnoise.HashInt2(x, y)
}

func hashInt3(x, y, z int32) int32 {
	return procnoise.HashInt3(x, y, z)
}

func hashFloat2(x, y float64) float64 {
	ix := int32(math.Round(x * 1000))
	iy := int32(math.Round(y * 1000))
	u := uint32(procnoise.HashInt2(ix, iy))
	return float64(u%1000000) / 1000000
}

func perlin2(x, y float64) float64 {
	return procnoise.Perlin2(x, y, 0)
}

func perlin3(x, y, z float64) float64 {
	return procnoise.Perlin3(x, y, z, 0)
}

func fbm2(x, y float64, octaves int) float64 {
	return procnoise.FBM2(x, y, octaves, 2, 0.5, 0)
}

func voronoi2(x, y float64) float64 {
	return procnoise.VoronoiDistance(x, y, 0)
}

func simplex2(x, y float64) float64 {
	return procnoise.Simplex2(x, y, 0)
}

func simplex3(x, y, z float64) float64 {
	return procnoise.Simplex3(x, y, z, 0)
}
