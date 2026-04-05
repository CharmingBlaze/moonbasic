package noisemod

import (
	"fmt"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) getNoise(rt *runtime.Runtime, args []value.Value, ix int) (*noiseObj, error) {
	hs, err := fixAllocHeap(rt, m)
	if err != nil {
		return nil, err
	}
	h, err := rt.ArgHandle(args, ix)
	if err != nil {
		return nil, err
	}
	o, err := heap.Cast[*noiseObj](hs, heap.Handle(h))
	if err != nil {
		return nil, err
	}
	return o, nil
}

func fixAllocHeap(rt *runtime.Runtime, m *Module) (*heap.Store, error) {
	if rt != nil && rt.Heap != nil {
		return rt.Heap, nil
	}
	if m.h != nil {
		return m.h, nil
	}
	return nil, runtime.Errorf("noise: heap not bound")
}

func (m *Module) noiseMake(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("NOISE.MAKE expects 0 arguments")
	}
	hs, err := fixAllocHeap(rt, m)
	if err != nil {
		return value.Nil, err
	}
	n := newNoiseObj()
	id, err := hs.Alloc(n)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) noiseFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	hs, err := fixAllocHeap(rt, m)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NOISE.FREE expects noise handle")
	}
	_ = hs.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func (m *Module) noiseSetType(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NOISE.SETTYPE expects (noise, type$)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.ensureMutable(); err != nil {
		return value.Nil, err
	}
	s, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	n.noiseType = strings.TrimSpace(s)
	return value.Nil, nil
}

func (m *Module) noiseSetSeed(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NOISE.SETSEED expects (noise, seed)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.ensureMutable(); err != nil {
		return value.Nil, err
	}
	seed, err := rt.ArgInt(args, 1)
	if err != nil {
		return value.Nil, err
	}
	n.seed = int32(seed)
	return value.Nil, nil
}

func (m *Module) noiseSetFrequency(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NOISE.SETFREQUENCY expects (noise, freq#)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.ensureMutable(); err != nil {
		return value.Nil, err
	}
	f, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	n.frequency = f
	return value.Nil, nil
}

func (m *Module) noiseSetOctaves(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NOISE.SETOCTAVES expects (noise, count)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.ensureMutable(); err != nil {
		return value.Nil, err
	}
	o, err := rt.ArgInt(args, 1)
	if err != nil {
		return value.Nil, err
	}
	n.octaves = int(o)
	if n.octaves < 1 {
		n.octaves = 1
	}
	return value.Nil, nil
}

func (m *Module) noiseSetLacunarity(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NOISE.SETLACUNARITY expects (noise, lac#)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.ensureMutable(); err != nil {
		return value.Nil, err
	}
	f, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	n.lacunarity = f
	return value.Nil, nil
}

func (m *Module) noiseSetGain(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NOISE.SETGAIN expects (noise, gain#)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.ensureMutable(); err != nil {
		return value.Nil, err
	}
	f, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	n.gain = f
	return value.Nil, nil
}

func (m *Module) noiseSetWeightedStrength(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NOISE.SETWEIGHTEDSTRENGTH expects (noise, strength#)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.ensureMutable(); err != nil {
		return value.Nil, err
	}
	f, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	n.weightedStrength = f
	return value.Nil, nil
}

func (m *Module) noiseSetPingPongStrength(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NOISE.SETPINGPONGSTRENGTH expects (noise, strength#)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.ensureMutable(); err != nil {
		return value.Nil, err
	}
	f, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	n.pingPongStrength = f
	return value.Nil, nil
}

func (m *Module) noiseSetCellularType(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NOISE.SETCELLULARTYPE expects (noise, type$)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.ensureMutable(); err != nil {
		return value.Nil, err
	}
	s, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	n.cellularType = strings.TrimSpace(s)
	return value.Nil, nil
}

func (m *Module) noiseSetCellularDistance(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NOISE.SETCELLULARDISTANCE expects (noise, func$)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.ensureMutable(); err != nil {
		return value.Nil, err
	}
	s, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	n.cellularDist = strings.TrimSpace(s)
	return value.Nil, nil
}

func (m *Module) noiseSetCellularJitter(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NOISE.SETCELLULARJITTER expects (noise, jitter#)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.ensureMutable(); err != nil {
		return value.Nil, err
	}
	f, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	n.cellularJitter = f
	return value.Nil, nil
}

func (m *Module) noiseSetDomainWarpType(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NOISE.SETDOMAINWARPTYPE expects (noise, type$)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.ensureMutable(); err != nil {
		return value.Nil, err
	}
	s, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	n.warpType = strings.TrimSpace(s)
	return value.Nil, nil
}

func (m *Module) noiseSetDomainWarpAmplitude(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NOISE.SETDOMAINWARPAMPLITUDE expects (noise, amp#)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.ensureMutable(); err != nil {
		return value.Nil, err
	}
	f, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	n.warpAmp = f
	return value.Nil, nil
}

func (m *Module) noiseGet(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("NOISE.GET expects (noise, x#, y#)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.assertLive(); err != nil {
		return value.Nil, err
	}
	x, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	y, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(n.Sample2D(x, y)), nil
}

func (m *Module) noiseGet3D(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("NOISE.GET3D expects (noise, x#, y#, z#)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.assertLive(); err != nil {
		return value.Nil, err
	}
	x, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	y, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	z, err := rt.ArgFloat(args, 3)
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(n.Sample3D(x, y, z)), nil
}

func (m *Module) noiseGetDomainWarped(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("NOISE.GETDOMAINWARPED expects (noise, x#, y#)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.assertLive(); err != nil {
		return value.Nil, err
	}
	x, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	y, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(n.SampleDomainWarped(x, y)), nil
}

func (m *Module) noiseGetNorm(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("NOISE.GETNORM expects (noise, x#, y#)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.assertLive(); err != nil {
		return value.Nil, err
	}
	x, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	y, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	v := normNoise(n.Sample2D(x, y))
	return value.FromFloat(v), nil
}

func (m *Module) noiseGetTileable(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("NOISE.GETTILEABLE expects (noise, x#, y#, w#, h#)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.assertLive(); err != nil {
		return value.Nil, err
	}
	x, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	y, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	w, err := rt.ArgFloat(args, 3)
	if err != nil {
		return value.Nil, err
	}
	h, err := rt.ArgFloat(args, 4)
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(n.SampleTileable(x, y, w, h)), nil
}

func (m *Module) noiseFillArray(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("NOISE.FILLARRAY expects (noise, arr, width, height, offsetX#, offsetY#)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.assertLive(); err != nil {
		return value.Nil, err
	}
	hs, err := fixAllocHeap(rt, m)
	if err != nil {
		return value.Nil, err
	}
	arr, err := heap.Cast[*heap.Array](hs, heap.Handle(args[1].IVal))
	if err != nil {
		return value.Nil, err
	}
	if arr.Kind != heap.ArrayKindFloat {
		return value.Nil, fmt.Errorf("NOISE.FILLARRAY: numeric array required")
	}
	wi, err := rt.ArgInt(args, 2)
	if err != nil {
		return value.Nil, err
	}
	hi, err := rt.ArgInt(args, 3)
	if err != nil {
		return value.Nil, err
	}
	offx, err := rt.ArgFloat(args, 4)
	if err != nil {
		return value.Nil, err
	}
	offy, err := rt.ArgFloat(args, 5)
	if err != nil {
		return value.Nil, err
	}
	width := int(wi)
	height := int(hi)
	if width < 1 || height < 1 {
		return value.Nil, fmt.Errorf("NOISE.FILLARRAY: width and height must be >= 1")
	}
	if width*height > arr.TotalElements() {
		return value.Nil, fmt.Errorf("NOISE.FILLARRAY: array smaller than width*height")
	}
	for yi := 0; yi < height; yi++ {
		for xi := 0; xi < width; xi++ {
			i := yi*width + xi
			v := n.Sample2D(float64(xi)+offx, float64(yi)+offy)
			if err := arr.Set([]int64{int64(i)}, v); err != nil {
				return value.Nil, err
			}
		}
	}
	return value.Nil, nil
}

func (m *Module) noiseFillArrayNorm(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("NOISE.FILLARRAYNORM expects (noise, arr, width, height, offsetX#, offsetY#)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.assertLive(); err != nil {
		return value.Nil, err
	}
	hs, err := fixAllocHeap(rt, m)
	if err != nil {
		return value.Nil, err
	}
	arr, err := heap.Cast[*heap.Array](hs, heap.Handle(args[1].IVal))
	if err != nil {
		return value.Nil, err
	}
	if arr.Kind != heap.ArrayKindFloat {
		return value.Nil, fmt.Errorf("NOISE.FILLARRAYNORM: numeric array required")
	}
	wi, err := rt.ArgInt(args, 2)
	if err != nil {
		return value.Nil, err
	}
	hi, err := rt.ArgInt(args, 3)
	if err != nil {
		return value.Nil, err
	}
	offx, err := rt.ArgFloat(args, 4)
	if err != nil {
		return value.Nil, err
	}
	offy, err := rt.ArgFloat(args, 5)
	if err != nil {
		return value.Nil, err
	}
	width := int(wi)
	height := int(hi)
	if width < 1 || height < 1 {
		return value.Nil, fmt.Errorf("NOISE.FILLARRAYNORM: width and height must be >= 1")
	}
	if width*height > arr.TotalElements() {
		return value.Nil, fmt.Errorf("NOISE.FILLARRAYNORM: array smaller than width*height")
	}
	for yi := 0; yi < height; yi++ {
		for xi := 0; xi < width; xi++ {
			i := yi*width + xi
			v := normNoise(n.Sample2D(float64(xi)+offx, float64(yi)+offy))
			if err := arr.Set([]int64{int64(i)}, v); err != nil {
				return value.Nil, err
			}
		}
	}
	return value.Nil, nil
}

func (m *Module) noiseMakePerlin(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NOISE.MAKEPERLIN expects (seed, freq#)")
	}
	seed, err := rt.ArgInt(args, 0)
	if err != nil {
		return value.Nil, err
	}
	f, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	hs, err := fixAllocHeap(rt, m)
	if err != nil {
		return value.Nil, err
	}
	n := newNoiseObj()
	n.noiseType = "perlin"
	n.seed = int32(seed)
	n.frequency = f
	id, err := hs.Alloc(n)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) noiseMakeSimplex(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NOISE.MAKESIMPLEX expects (seed, freq#)")
	}
	seed, err := rt.ArgInt(args, 0)
	if err != nil {
		return value.Nil, err
	}
	f, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	hs, err := fixAllocHeap(rt, m)
	if err != nil {
		return value.Nil, err
	}
	n := newNoiseObj()
	n.noiseType = "simplex"
	n.seed = int32(seed)
	n.frequency = f
	id, err := hs.Alloc(n)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) noiseMakeFractal(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("NOISE.MAKEFRACTAL expects (seed, freq#, octaves, type$)")
	}
	seed, err := rt.ArgInt(args, 0)
	if err != nil {
		return value.Nil, err
	}
	f, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	oc, err := rt.ArgInt(args, 2)
	if err != nil {
		return value.Nil, err
	}
	typ, err := rt.ArgString(args, 3)
	if err != nil {
		return value.Nil, err
	}
	hs, err := fixAllocHeap(rt, m)
	if err != nil {
		return value.Nil, err
	}
	n := newNoiseObj()
	n.seed = int32(seed)
	n.frequency = f
	n.octaves = int(oc)
	if n.octaves < 1 {
		n.octaves = 1
	}
	switch strings.ToUpper(strings.TrimSpace(typ)) {
	case "FBM", "FRACTAL_FBM":
		n.noiseType = "fractal_fbm"
	case "RIDGED", "FRACTAL_RIDGED":
		n.noiseType = "fractal_ridged"
	case "PINGPONG", "FRACTAL_PINGPONG":
		n.noiseType = "fractal_pingpong"
	default:
		n.noiseType = "fractal_fbm"
	}
	id, err := hs.Alloc(n)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) noiseMakeCellular(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("NOISE.MAKECELLULAR expects (seed, freq#, celltype$)")
	}
	seed, err := rt.ArgInt(args, 0)
	if err != nil {
		return value.Nil, err
	}
	f, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	ct, err := rt.ArgString(args, 2)
	if err != nil {
		return value.Nil, err
	}
	hs, err := fixAllocHeap(rt, m)
	if err != nil {
		return value.Nil, err
	}
	n := newNoiseObj()
	n.noiseType = "cellular"
	n.seed = int32(seed)
	n.frequency = f
	n.cellularType = strings.TrimSpace(ct)
	id, err := hs.Alloc(n)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) noiseMakeDomainWarp(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("NOISE.MAKEDOMAINWARP expects (seed, freq#, amp#)")
	}
	seed, err := rt.ArgInt(args, 0)
	if err != nil {
		return value.Nil, err
	}
	f, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	a, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	hs, err := fixAllocHeap(rt, m)
	if err != nil {
		return value.Nil, err
	}
	n := newNoiseObj()
	n.noiseType = "domain_warp"
	n.seed = int32(seed)
	n.frequency = f
	n.warpAmp = a
	id, err := hs.Alloc(n)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}
