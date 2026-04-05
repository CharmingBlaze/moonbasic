package noisemod

import (
	"math"
	"strings"

	"moonbasic/runtime/procnoise"
)

func normNoise(v float64) float64 {
	u := (v + 1) * 0.5
	if u < 0 {
		return 0
	}
	if u > 1 {
		return 1
	}
	return u
}

// evalAtNoiseSpace evaluates in scaled noise space (coordinates already multiplied by frequency).
func (n *noiseObj) evalAtNoiseSpace(xf, yf float64) float64 {
	switch strings.ToUpper(strings.TrimSpace(n.noiseType)) {
	case "DOMAIN_WARP":
		wx, wy := procnoise.Warp2(xf, yf, n.warpAmp, n.seed)
		return procnoise.ValueNoise2(wx, wy, n.seed)
	default:
		return n.rawAt(xf, yf)
	}
}

// Sample2D returns noise in ~[-1,1] at world coordinates.
func (n *noiseObj) Sample2D(x, y float64) float64 {
	n.touch()
	return n.evalAtNoiseSpace(x*n.frequency, y*n.frequency)
}

// SampleDomainWarped warps (x,y) in noise space then evaluates the current type (adds organic turbulence).
func (n *noiseObj) SampleDomainWarped(x, y float64) float64 {
	n.touch()
	xf, yf := x*n.frequency, y*n.frequency
	wx, wy := procnoise.Warp2(xf, yf, n.warpAmp, n.seed)
	return n.rawAt(wx, wy)
}

// Sample3D maps types to a cheap 3D field (~[-1,1]).
func (n *noiseObj) Sample3D(x, y, z float64) float64 {
	n.touch()
	xf, yf, zf := x*n.frequency, y*n.frequency, z*n.frequency
	switch strings.ToUpper(strings.TrimSpace(n.noiseType)) {
	case "DOMAIN_WARP":
		wx, wy := procnoise.Warp2(xf, yf, n.warpAmp, n.seed)
		wz := procnoise.Simplex2(yf, zf, n.seed+11) * n.warpAmp
		return procnoise.Perlin3(wx, wy, wz, n.seed)
	case "CELLULAR", "CELL":
		return (procnoise.Cellular2(xf, yf, n.seed) + procnoise.Cellular2(yf, zf, n.seed+3) + procnoise.Cellular2(zf, xf, n.seed+5)) / 3
	default:
		return procnoise.Perlin3(xf, yf, zf, n.seed)
	}
}

// SampleTileable approximate seamless tiling using a torus parameterization (see procnoise.Tileable2).
func (n *noiseObj) SampleTileable(x, y, w, h float64) float64 {
	n.touch()
	return procnoise.Tileable2(x, y, w, h, n.seed, func(px, py float64, _ int32) float64 {
		return n.evalAtNoiseSpace(px*n.frequency, py*n.frequency)
	})
}

func (n *noiseObj) rawAt(xf, yf float64) float64 {
	typ := strings.ToUpper(strings.TrimSpace(n.noiseType))
	switch typ {
	case "PERLIN":
		return procnoise.Perlin2(xf, yf, n.seed)
	case "SIMPLEX", "SIMPLEX_SMOOTH":
		v := procnoise.Simplex2(xf, yf, n.seed)
		if typ == "SIMPLEX_SMOOTH" {
			v = (v + procnoise.Simplex2(xf+0.1, yf+0.1, n.seed) + procnoise.Simplex2(xf-0.1, yf-0.1, n.seed)) / 3
		}
		return v
	case "VALUE", "VALUE_CUBIC":
		return procnoise.ValueNoise2(xf, yf, n.seed)
	case "CELLULAR", "CELL":
		return n.sampleCellular(xf, yf)
	case "FRACTAL_FBM":
		if n.weightedStrength > 0 {
			return procnoise.FBM2Weighted(xf, yf, n.octaves, n.lacunarity, n.gain, n.weightedStrength, n.seed)
		}
		return procnoise.FBM2(xf, yf, n.octaves, n.lacunarity, n.gain, n.seed)
	case "FRACTAL_RIDGED":
		return procnoise.RidgedMulti2(xf, yf, n.octaves, n.lacunarity, n.gain, n.seed)
	case "FRACTAL_PINGPONG", "PINGPONG":
		ps := n.pingPongStrength
		if ps <= 0 {
			ps = 2
		}
		return procnoise.PingPongFBM2(xf, yf, n.octaves, n.lacunarity, n.gain, ps, n.seed)
	default:
		return procnoise.Perlin2(xf, yf, n.seed)
	}
}

func (n *noiseObj) sampleCellular(xf, yf float64) float64 {
	ct := strings.ToUpper(strings.TrimSpace(n.cellularType))
	switch ct {
	case "CELL_VALUE", "CELLVALUE":
		cx := int32(math.Floor(xf))
		cy := int32(math.Floor(yf))
		return float64(procnoise.HashInt2(cx+n.seed*13, cy+n.seed*17)%1024)/512 - 1
	default:
		if strings.ToUpper(n.cellularDist) == "MANHATTAN" {
			// crude: scale coords for Manhattan-ish Voronoi
			return procnoise.Cellular2(xf*1.2, yf*1.2, n.seed)
		}
		return procnoise.Cellular2(xf, yf, n.seed)
	}
}
