package procnoise

import (
	"testing"
)

func TestValueNoise2Deterministic(t *testing.T) {
	a := ValueNoise2(1.25, -3.5, 42)
	b := ValueNoise2(1.25, -3.5, 42)
	if a != b {
		t.Fatalf("same coords+seed: got %v vs %v", a, b)
	}
	c := ValueNoise2(1.25, -3.5, 43)
	if c == a {
		t.Fatalf("different seed should change value")
	}
}

func TestPerlin2Range(t *testing.T) {
	for i := 0; i < 50; i++ {
		v := Perlin2(float64(i)*0.1, float64(i)*-0.07, 7)
		if v < -1.01 || v > 1.01 {
			t.Fatalf("out of expected range: %v", v)
		}
	}
}
