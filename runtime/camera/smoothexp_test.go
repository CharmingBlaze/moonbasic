package mbcamera

import (
	"math"
	"testing"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func TestCamSmoothExpConverges(t *testing.T) {
	var m Module
	rt := &runtime.Runtime{}
	cur := value.FromFloat(0.0)
	tgt := value.FromFloat(10.0)
	hz := value.FromFloat(30.0)
	dt := value.FromFloat(1.0 / 60.0)
	var last float64
	for i := 0; i < 200; i++ {
		v, err := m.camSmoothExp(rt, cur, tgt, hz, dt)
		if err != nil {
			t.Fatal(err)
		}
		f, _ := v.ToFloat()
		last = f
		cur = v
	}
	if math.Abs(last-10.0) > 0.01 {
		t.Fatalf("expected ~10, got %g", last)
	}
}

func TestCamSmoothExpDtZero(t *testing.T) {
	var m Module
	rt := &runtime.Runtime{}
	v, err := m.camSmoothExp(rt, value.FromFloat(3.0), value.FromFloat(9.0), value.FromFloat(20.0), value.FromFloat(0.0))
	if err != nil {
		t.Fatal(err)
	}
	f, _ := v.ToFloat()
	if math.Abs(f-3.0) > 1e-9 {
		t.Fatalf("dt<=0 should keep current, got %g", f)
	}
}
