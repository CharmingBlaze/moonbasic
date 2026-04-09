//go:build cgo || (windows && !cgo)

package mbcamera

import (
	"math"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// orbitCam matches ApplySetOrbit / Camera.OrbitEntity (examples/mario64/main_entities.mb defaults).
func orbitCam(tx, ty, tz, yaw, pitch, dist float32) rl.Camera3D {
	if dist < 0.15 {
		dist = 0.15
	}
	maxPitch := float32(1.45)
	if pitch > maxPitch {
		pitch = maxPitch
	}
	if pitch < -maxPitch {
		pitch = -maxPitch
	}
	sy, cy := math.Sin(float64(yaw)), math.Cos(float64(yaw))
	sp, cp := math.Sin(float64(pitch)), math.Cos(float64(pitch))
	hdist := float64(dist) * cp
	px := float32(float64(tx) + sy*hdist)
	py := float32(float64(ty) + float64(dist)*sp)
	pz := float32(float64(tz) + cy*hdist)
	return rl.Camera3D{
		Position:   rl.Vector3{X: px, Y: py, Z: pz},
		Target:     rl.Vector3{X: tx, Y: ty, Z: tz},
		Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
		Fovy:       55,
		Projection: rl.CameraPerspective,
	}
}

// Regression: mario64 entity sample — floor, platforms, player + hat must not be CPU-culled as invisible.
func TestExtractFrustumMarioEntitySceneVisible(t *testing.T) {
	aspect := float32(960.0 / 540.0)
	cam := orbitCam(0, 1, 0, 0, 0.22, 7.5)
	f := ExtractFrustum(cam, aspect)

	tx, ty, tz := cam.Target.X, cam.Target.Y, cam.Target.Z
	if !f.PointVisible(tx, ty, tz) {
		t.Fatalf("camera target (player) should be visible, got cam Pos=(%v,%v,%v) T=(%v,%v,%v)",
			cam.Position.X, cam.Position.Y, cam.Position.Z, tx, ty, tz)
	}

	// Floor: 48×0.5×48 centered at (0,-0.25,0)
	floorMinX, floorMinY, floorMinZ := float32(-24.0), float32(-0.5), float32(-24.0)
	floorMaxX, floorMaxY, floorMaxZ := float32(24.0), float32(0.0), float32(24.0)
	if !f.AABBVisible(floorMinX, floorMinY, floorMinZ, floorMaxX, floorMaxY, floorMaxZ) {
		t.Fatalf("floor AABB should intersect frustum (CPU culling was hiding the whole ground)")
	}

	// Player sphere center (0,1,0) r=0.45
	if !f.SphereVisible(0, 1, 0, 0.45) {
		t.Fatalf("player sphere should be visible")
	}

	// Hat (child), approximate world center ~ (0, 1.85, 0), r ~ 0.16
	if !f.SphereVisible(0, 1.85, 0, 0.2) {
		t.Fatalf("hat sphere should be visible")
	}
}
