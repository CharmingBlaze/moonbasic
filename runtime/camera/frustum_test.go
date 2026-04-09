//go:build cgo || (windows && !cgo)

package mbcamera

import (
	"math"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ExtractFrustum plane extraction must stay numerically stable; normals are normalized.
func TestExtractFrustumPlanesNormalized(t *testing.T) {
	cam := rl.Camera3D{
		Position:   rl.Vector3{X: 0, Y: 3, Z: 10},
		Target:     rl.Vector3{X: 0, Y: 1, Z: 0},
		Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
		Fovy:       55,
		Projection: rl.CameraPerspective,
	}
	f := ExtractFrustum(cam, 16.0/9.0)
	for i := range f.planes {
		p := f.planes[i]
		mag := math.Sqrt(float64(p.a*p.a + p.b*p.b + p.c*p.c))
		if mag < 0.9 || mag > 1.1 {
			t.Fatalf("plane %d: expected unit normal, got length %g", i, mag)
		}
	}
}

func TestProjectionMatrixUsesRLCullDistances(t *testing.T) {
	near := float32(rl.GetCullDistanceNear())
	far := float32(rl.GetCullDistanceFar())
	if near <= 0 || far <= near {
		t.Fatalf("unexpected RL clip range near=%g far=%g", near, far)
	}
	cam := rl.Camera3D{
		Position:   rl.Vector3{X: 0, Y: 2, Z: 8},
		Target:     rl.Vector3{X: 0, Y: 0, Z: 0},
		Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
		Fovy:       45,
		Projection: rl.CameraPerspective,
	}
	p := projectionMatrixForFrustum(cam, 16.0/9.0)
	// Perspective depth scale relates to near/far; zero near would collapse M10/M14.
	if math.Abs(float64(p.M14)) < 1e-6 {
		t.Fatal("projection matrix depth coupling looks degenerate")
	}
}

// Regression: look-at target must lie inside the frustum (entity CPU culling uses this).
func TestExtractFrustumLookAtTargetInside(t *testing.T) {
	cam := rl.Camera3D{
		Position:   rl.Vector3{X: 0, Y: 3, Z: 10},
		Target:     rl.Vector3{X: 0, Y: 1, Z: 0},
		Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
		Fovy:       55,
		Projection: rl.CameraPerspective,
	}
	f := ExtractFrustum(cam, 16.0/9.0)
	tx, ty, tz := cam.Target.X, cam.Target.Y, cam.Target.Z
	if !f.PointVisible(tx, ty, tz) {
		t.Fatalf("look-at target (%g,%g,%g) should be inside frustum", tx, ty, tz)
	}
}
