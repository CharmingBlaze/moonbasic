//go:build cgo || (windows && !cgo)

package mbgame

import (
	"testing"

	"moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func TestCollisionVecParity2D(t *testing.T) {
	h := heap.New()
	m := NewModule()
	m.BindHeap(h)

	pa, err := mbmatrix.AllocVec2Value(h, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	sa, err := mbmatrix.AllocVec2Value(h, 10, 10)
	if err != nil {
		t.Fatal(err)
	}
	pb, err := mbmatrix.AllocVec2Value(h, 5, 5)
	if err != nil {
		t.Fatal(err)
	}
	sb, err := mbmatrix.AllocVec2Value(h, 10, 10)
	if err != nil {
		t.Fatal(err)
	}
	v, err := m.collisionBoxOverlap2D([]value.Value{pa, sa, pb, sb})
	if err != nil {
		t.Fatal(err)
	}
	if v.Kind != value.KindBool {
		t.Fatalf("expected bool, got kind %v", v.Kind)
	}
	want := boxCollide2D(0, 0, 10, 10, 5, 5, 10, 10)
	got := v.IVal != 0
	if got != want {
		t.Fatalf("BOXOVERLAP2D got %v want %v", got, want)
	}
}

func TestCollisionVecParity3D(t *testing.T) {
	h := heap.New()
	m := NewModule()
	m.BindHeap(h)

	c1, err := mbmatrix.AllocVec3Value(h, 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	c2, err := mbmatrix.AllocVec3Value(h, 3, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	v, err := m.collisionSphereOverlap3D([]value.Value{c1, value.FromFloat(2), c2, value.FromFloat(2)})
	if err != nil {
		t.Fatal(err)
	}
	if v.Kind != value.KindBool {
		t.Fatalf("expected bool, got kind %v", v.Kind)
	}
	want := sphereCollide3D(0, 0, 0, 2, 3, 0, 0, 2)
	got := v.IVal != 0
	if got != want {
		t.Fatalf("SPHEREOVERLAP3D got %v want %v", got, want)
	}
}
