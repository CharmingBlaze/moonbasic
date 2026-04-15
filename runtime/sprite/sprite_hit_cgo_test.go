//go:build cgo || (windows && !cgo)

package mbsprite

import (
	"math"
	"testing"
)

func TestSpritePointHit_unrotatedMatchesAABBWithOrigin(t *testing.T) {
	// dest (10,20), origin (0,0), 32x32 — same as old POINTHIT for this case
	if !spritePointHit(10, 20, 32, 32, 0, 0, 0, 15, 25) {
		t.Fatal("expected center sample inside")
	}
	if spritePointHit(10, 20, 32, 32, 0, 0, 0, 9, 25) {
		t.Fatal("expected left of drawable rect outside")
	}
	// origin shifts drawable top-left: dest (10,20), origin (5,5) -> rect at (5,15) size 32x32
	if !spritePointHit(10, 20, 32, 32, 5, 5, 0, 10, 20) {
		t.Fatal("expected point on top-left of drawn quad")
	}
}

func TestSpriteHitOverlap_axisAligned(t *testing.T) {
	a := &spriteObj{frameW: 32, frameH: 32, x: 0, y: 0, scaleX: 1, scaleY: 1}
	b := &spriteObj{frameW: 32, frameH: 32, x: 16, y: 0, scaleX: 1, scaleY: 1}
	if !spriteHitOverlap(a, b) {
		t.Fatal("expected overlap")
	}
	b.x = 64
	if spriteHitOverlap(a, b) {
		t.Fatal("expected separation")
	}
}

func TestSpriteHitOverlap_rotatedSquaresOverlap(t *testing.T) {
	// Two same-size squares, one rotated 45°, overlapping centers
	a := &spriteObj{frameW: 40, frameH: 40, x: 100, y: 100, scaleX: 1, scaleY: 1, originX: 20, originY: 20}
	b := &spriteObj{frameW: 40, frameH: 40, x: 100, y: 100, scaleX: 1, scaleY: 1, originX: 20, originY: 20}
	a.rotRad = float32(math.Pi / 4)
	b.rotRad = 0
	if !spriteHitOverlap(a, b) {
		t.Fatal("expected overlap at same dest with rotation")
	}
}

func TestSpriteQuadCorners_matchesRaylibRotation0(t *testing.T) {
	tl, bl, br, tr := spriteQuadCorners(10, 20, 32, 48, 5, 6, 0)
	if tl.x != 5 || tl.y != 14 {
		t.Fatalf("tl got (%v,%v)", tl.x, tl.y)
	}
	if tr.x != 37 || tr.y != 14 {
		t.Fatalf("tr got (%v,%v)", tr.x, tr.y)
	}
	if bl.x != 5 || bl.y != 62 {
		t.Fatalf("bl got (%v,%v)", bl.x, bl.y)
	}
	if br.x != 37 || br.y != 62 {
		t.Fatalf("br got (%v,%v)", br.x, br.y)
	}
}
