package value

import "testing"

func TestAddInt(t *testing.T) {
	v, err := Add(Int(2), Int(3))
	if err != nil || v.IVal != 5 {
		t.Fatal(v, err)
	}
}

func TestAddFloatCoerce(t *testing.T) {
	v, err := Add(Int(1), Float(2.5))
	if err != nil || v.Kind != KindFloat || v.FVal != 3.5 {
		t.Fatal(v, err)
	}
}
