package value

import (
	"testing"
)

type testStringMap map[int32]string

func (m testStringMap) GetString(i int32) (string, bool) {
	s, ok := m[i]
	return s, ok
}

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

func TestStringAtPoolThenHeap(t *testing.T) {
	pool := []string{"hello"}
	heap := testStringMap{7: "runtime"}
	if s := StringAt(FromStringIndex(0), pool, nil); s != "hello" {
		t.Fatalf("pool: got %q", s)
	}
	if s := StringAt(FromStringIndex(7), pool, heap); s != "runtime" {
		t.Fatalf("heap: got %q", s)
	}
}

func TestEqualStringValueCrossTable(t *testing.T) {
	pool := []string{"a", "b"}
	heap := testStringMap{5: "a"}
	a := FromStringIndex(0) // pool "a"
	b := FromStringIndex(5) // heap "a"
	if !EqualStringValue(a, b, pool, heap) {
		t.Fatal("expected text-equal strings")
	}
	if Equal(a, b) {
		t.Fatal("index identity should differ")
	}
}

func TestLessStringCrossTable(t *testing.T) {
	pool := []string{"m"}
	heap := testStringMap{9: "z"}
	x := FromStringIndex(0)
	y := FromStringIndex(9)
	lt, err := Less(x, y, pool, heap)
	if err != nil || !lt {
		t.Fatalf("Less m < z: lt=%v err=%v", lt, err)
	}
}
