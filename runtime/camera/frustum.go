//go:build cgo || (windows && !cgo)

// View-frustum planes from combined projection*view (same as Raylib/RHI clip transform).
// Plane order: Left, Right, Bottom, Top, Near, Far. Stack-only types on extraction/tests — no heap allocs.
package mbcamera

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Plane is a normalised half-space ax+by+cz+d=0; visible side has ax+by+cz+d > 0.
type Plane struct {
	a, b, c, d float32
}

func (p Plane) distanceToPoint(x, y, z float32) float32 {
	return p.a*x + p.b*y + p.c*z + p.d
}

// Frustum holds six normalised planes: Left, Right, Bottom, Top, Near, Far.
type Frustum struct {
	planes [6]Plane
}

// projectionMatrixForFrustum matches rl.BeginMode3D + GetCameraProjectionMatrix math but uses
// rlGetCullDistanceNear/Far (see RL_CULL_DISTANCE_*) — not the hardcoded 0.01/1000 in
// GetCameraProjectionMatrix. Mismatched near/far made CPU culling reject geometry the GPU still drew.
func projectionMatrixForFrustum(cam rl.Camera3D, aspectRatio float32) rl.Matrix {
	near := float32(rl.GetCullDistanceNear())
	far := float32(rl.GetCullDistanceFar())
	if cam.Projection == rl.CameraOrthographic {
		top := cam.Fovy / 2.0
		right := top * aspectRatio
		return rl.MatrixOrtho(-right, right, -top, top, near, far)
	}
	fovRad := cam.Fovy * float32(math.Pi) / 180.0
	return rl.MatrixPerspective(fovRad, aspectRatio, near, far)
}

// ExtractFrustum builds frustum planes from PV = projection * view (column vector: clip = PV * v).
// Planes are extracted from the **rows** of PV (Gribb–Hartmann / clip-space convention), not from
// columns — column-based c0..c3 produced inverted tests and culled entire scenes.
func ExtractFrustum(cam rl.Camera3D, aspectRatio float32) Frustum {
	view := rl.GetCameraMatrix(cam)
	proj := projectionMatrixForFrustum(cam, aspectRatio)
	pv := rl.MatrixMultiply(proj, view)

	// Row i in column-major storage: (M[i], M[i+4], M[i+8], M[i+12])
	r0 := [4]float32{pv.M0, pv.M4, pv.M8, pv.M12}
	r1 := [4]float32{pv.M1, pv.M5, pv.M9, pv.M13}
	r2 := [4]float32{pv.M2, pv.M6, pv.M10, pv.M14}
	r3 := [4]float32{pv.M3, pv.M7, pv.M11, pv.M15}

	var f Frustum
	f.planes[0] = normalisePlane(planeAdd(r3, r0)) // left
	f.planes[1] = normalisePlane(planeSub(r3, r0)) // right
	f.planes[2] = normalisePlane(planeAdd(r3, r1)) // bottom
	// Top and far half-spaces are inverted vs the usual row-sum recipe for rl.MatrixPerspective;
	// without negation, center-of-frustum points fail PointVisible and entity culling drops all draws.
	f.planes[3] = normalisePlane(negatePlane(planeSub(r3, r1))) // top
	f.planes[4] = normalisePlane(planeAdd(r3, r2))              // near
	f.planes[5] = normalisePlane(negatePlane(planeSub(r3, r2))) // far
	return f
}

func planeAdd(a, b [4]float32) Plane {
	return Plane{a[0] + b[0], a[1] + b[1], a[2] + b[2], a[3] + b[3]}
}

func planeSub(a, b [4]float32) Plane {
	return Plane{a[0] - b[0], a[1] - b[1], a[2] - b[2], a[3] - b[3]}
}

func negatePlane(p Plane) Plane {
	return Plane{-p.a, -p.b, -p.c, -p.d}
}

func normalisePlane(p Plane) Plane {
	mag := float32(math.Sqrt(float64(p.a*p.a + p.b*p.b + p.c*p.c)))
	if mag < 1e-8 {
		return p
	}
	inv := 1.0 / mag
	return Plane{p.a * inv, p.b * inv, p.c * inv, p.d * inv}
}

// SphereVisible reports whether a sphere intersects the frustum (conservative).
func (f Frustum) SphereVisible(cx, cy, cz, r float32) bool {
	for i := 0; i < 6; i++ {
		p := &f.planes[i]
		if p.a*cx+p.b*cy+p.c*cz+p.d < -r {
			return false
		}
	}
	return true
}

// AABBVisible tests an axis-aligned box against the frustum (positive-vertex method).
func (f Frustum) AABBVisible(minX, minY, minZ, maxX, maxY, maxZ float32) bool {
	for i := 0; i < 6; i++ {
		p := &f.planes[i]
		var px, py, pz float32
		if p.a >= 0 {
			px = maxX
		} else {
			px = minX
		}
		if p.b >= 0 {
			py = maxY
		} else {
			py = minY
		}
		if p.c >= 0 {
			pz = maxZ
		} else {
			pz = minZ
		}
		if p.a*px+p.b*py+p.c*pz+p.d < 0 {
			return false
		}
	}
	return true
}

// PointVisible returns true if the point is on the visible side of all planes.
func (f Frustum) PointVisible(x, y, z float32) bool {
	for i := 0; i < 6; i++ {
		p := &f.planes[i]
		if p.a*x+p.b*y+p.c*z+p.d < 0 {
			return false
		}
	}
	return true
}

// WithinDistance is true if (cx,cy,cz) is within maxDist of (camX,camY,camZ) — squared distance only.
func WithinDistance(cx, cy, cz, camX, camY, camZ, maxDist float32) bool {
	dx := cx - camX
	dy := cy - camY
	dz := cz - camZ
	return dx*dx+dy*dy+dz*dz <= maxDist*maxDist
}

// BehindHorizon returns true if the top of the terrain feature at (cx,cz) is entirely below the camera's bottom-of-view angle.
func BehindHorizon(camX, camY, camZ, maxY, cx, cz, camPitchDeg, fovYDeg float32) bool {
	dx := cx - camX
	dz := cz - camZ
	horizDist := float32(math.Sqrt(float64(dx*dx + dz*dz)))
	if horizDist < 1.0 {
		return false
	}
	dy := maxY - camY
	angleToTop := float32(math.Atan2(float64(dy), float64(horizDist))) * (180.0 / math.Pi)
	bottomAngle := camPitchDeg - fovYDeg*0.5
	return angleToTop < bottomAngle
}

// --- Module-level state (updated at CAMERA.BEGIN / cleared at CAMERA.END) ---

var (
	activeFrustum Frustum
	frustumValid  bool
	activeCamPos  [3]float32
	activeCamPitch float32
	activeCamFOV   float32
	globalMaxDist  float32 = 1000.0
)

func setActiveFrustum(cam rl.Camera3D, aspect float32) {
	activeFrustum = ExtractFrustum(cam, aspect)
	frustumValid = true
	activeCamPos[0] = cam.Position.X
	activeCamPos[1] = cam.Position.Y
	activeCamPos[2] = cam.Position.Z
	activeCamFOV = cam.Fovy
	fwd := rl.Vector3Subtract(cam.Target, cam.Position)
	fwd = rl.Vector3Normalize(fwd)
	activeCamPitch = float32(math.Asin(float64(fwd.Y))) * (180.0 / math.Pi)
}

func clearActiveFrustum() { frustumValid = false }

// SetGlobalMaxDistance sets default max draw distance for CULL.INRANGE (3-arg) and frustum+distance tests.
func SetGlobalMaxDistance(d float32) { globalMaxDist = d }

// GlobalMaxDistance returns the current default max draw distance.
func GlobalMaxDistance() float32 { return globalMaxDist }

// AABBVisibleActive uses the frustum from the current CAMERA.BEGIN (if any).
func AABBVisibleActive(minX, minY, minZ, maxX, maxY, maxZ float32) bool {
	if !frustumValid {
		return true
	}
	return activeFrustum.AABBVisible(minX, minY, minZ, maxX, maxY, maxZ)
}

// SphereVisibleActive uses the active frustum from CAMERA.BEGIN.
func SphereVisibleActive(cx, cy, cz, r float32) bool {
	if !frustumValid {
		return true
	}
	return activeFrustum.SphereVisible(cx, cy, cz, r)
}

// WithinDistanceActive compares to active camera position and GlobalMaxDistance.
func WithinDistanceActive(cx, cy, cz float32) bool {
	return WithinDistance(cx, cy, cz, activeCamPos[0], activeCamPos[1], activeCamPos[2], globalMaxDist)
}

// BehindHorizonActive uses pitch/FOV captured at CAMERA.BEGIN.
func BehindHorizonActive(maxY, cx, cz float32) bool {
	if !frustumValid {
		return false
	}
	return BehindHorizon(
		activeCamPos[0], activeCamPos[1], activeCamPos[2],
		maxY, cx, cz,
		activeCamPitch, activeCamFOV,
	)
}
