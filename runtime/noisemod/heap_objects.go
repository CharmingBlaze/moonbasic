package noisemod

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

// noiseObj owns a configured procedural noise generator (pure Go; no C allocations).
// Coordinates with runtime/procnoise; stateful sampling shares the same core as legacy PERLIN/SIMPLEX in mbgame.
type noiseObj struct {
	noiseType string
	seed      int32

	frequency  float64
	octaves    int
	lacunarity float64
	gain       float64

	weightedStrength float64
	pingPongStrength float64

	cellularType string
	cellularDist string
	cellularJitter float64

	warpType string
	warpAmp  float64

	used  bool
	freed bool
}

func newNoiseObj() *noiseObj {
	return &noiseObj{
		noiseType:      "perlin",
		seed:           1337,
		frequency:      0.01,
		octaves:        3,
		lacunarity:     2,
		gain:           0.5,
		pingPongStrength: 2,
		cellularType:   "distance",
		cellularDist:   "euclidean",
		cellularJitter: 1,
		warpType:       "opensimplex2",
		warpAmp:        1,
	}
}

func (n *noiseObj) TypeName() string { return "Noise" }

func (n *noiseObj) TypeTag() uint16 { return heap.TagNoise }

func (n *noiseObj) Free() {
	if n.freed {
		return
	}
	n.freed = true
}

func (n *noiseObj) assertLive() error {
	if n.freed {
		return runtime.Errorf("Noise: use after free")
	}
	return nil
}

func (n *noiseObj) touch() {
	n.used = true
}

func (n *noiseObj) ensureMutable() error {
	if err := n.assertLive(); err != nil {
		return err
	}
	if n.used {
		return runtime.Errorf("Noise: call configuration (Set*) before the first Get/Fill*")
	}
	return nil
}
