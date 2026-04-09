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

// Frustum holds the world→clip matrix (Raylib: MatrixMultiply(view, proj)) plus legacy plane forms.
// Sphere/AABB visibility uses clip-space homogeneous tests against clipM so CPU culling matches
// the GPU frustum (row-sum plane extraction can false-negative for moving cameras/objects).
type Frustum struct {
	planes [6]Plane
	clipM  rl.Matrix
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

// ExtractFrustum builds frustum planes from the world→clip matrix used with Raylib’s column-vector
// multiply (cx = M0*x+M4*y+M8*z+M12, …). That matrix is MatrixMultiply(view, proj), which composes
// as the same linear map as proj·view in standard math — not MatrixMultiply(proj, view).
// Planes are extracted from the **rows** of that matrix (Gribb–Hartmann / clip-space convention).
//
// Use projectionMatrixForFrustum (RL cull near/far) — not rl.GetCameraProjectionMatrix, which is
// hardcoded 0.01/1000 in raylib-go and does not match rlgl when clip planes differ.
func ExtractFrustum(cam rl.Camera3D, aspectRatio float32) Frustum {
	view := rl.GetCameraMatrix(cam)
	proj := projectionMatrixForFrustum(cam, aspectRatio)
	pv := rl.MatrixMultiply(view, proj)

	// Row i in column-major storage: (M[i], M[i+4], M[i+8], M[i+12])
	r0 := [4]float32{pv.M0, pv.M4, pv.M8, pv.M12}
	r1 := [4]float32{pv.M1, pv.M5, pv.M9, pv.M13}
	r2 := [4]float32{pv.M2, pv.M6, pv.M10, pv.M14}
	r3 := [4]float32{pv.M3, pv.M7, pv.M11, pv.M15}

	var f Frustum
	f.planes[0] = normalisePlane(planeAdd(r3, r0)) // left
	f.planes[1] = normalisePlane(planeSub(r3, r0)) // right
	f.planes[2] = normalisePlane(planeAdd(r3, r1)) // bottom
	// Top and far half-spaces: negate vs raw r3-r1 / r3-r2 for rl.MatrixPerspective + row extraction; see tests.
	f.planes[3] = normalisePlane(negatePlane(planeSub(r3, r1))) // top
	f.planes[4] = normalisePlane(planeAdd(r3, r2))              // near
	f.planes[5] = normalisePlane(negatePlane(planeSub(r3, r2))) // far
	f.clipM = pv
	return f
}

// clipPointInsideNDC returns true if (x,y,z,1) projects inside the default NDC cube [-1,1] on all
// axes after perspective divide — matches OpenGL clip + perspective behavior used with rl.BeginMode3D.
func clipPointInsideNDC(pv rl.Matrix, x, y, z float32) bool {
	cx := pv.M0*x + pv.M4*y + pv.M8*z + pv.M12
	cy := pv.M1*x + pv.M5*y + pv.M9*z + pv.M13
	cz := pv.M2*x + pv.M6*y + pv.M10*z + pv.M14
	cw := pv.M3*x + pv.M7*y + pv.M11*z + pv.M15
	const eps = float32(8e-4)
	aw := float32(math.Abs(float64(cw)))
	if aw < 1e-6 {
		return false
	}
	ndx := cx / cw
	ndy := cy / cw
	ndz := cz / cw
	if ndx < -1-eps || ndx > 1+eps {
		return false
	}
	if ndy < -1-eps || ndy > 1+eps {
		return false
	}
	if ndz < -1-eps || ndz > 1+eps {
		return false
	}
	return true
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
// Uses clip-space tests on the sphere's axis-aligned bounding box corners vs clipM.
func (f Frustum) SphereVisible(cx, cy, cz, r float32) bool {
	if r < 0 {
		r = 0
	}
	x0, x1 := cx-r, cx+r
	y0, y1 := cy-r, cy+r
	z0, z1 := cz-r, cz+r
	for _, x := range []float32{x0, x1} {
		for _, y := range []float32{y0, y1} {
			for _, z := range []float32{z0, z1} {
				if clipPointInsideNDC(f.clipM, x, y, z) {
					return true
				}
			}
		}
	}
	return false
}

// AABBVisible tests an axis-aligned box against the frustum (8 clip-space corners).
func (f Frustum) AABBVisible(minX, minY, minZ, maxX, maxY, maxZ float32) bool {
	for _, x := range []float32{minX, maxX} {
		for _, y := range []float32{minY, maxY} {
			for _, z := range []float32{minZ, maxZ} {
				if clipPointInsideNDC(f.clipM, x, y, z) {
					return true
				}
			}
		}
	}
	return false
}

// PointVisible returns true if the world point projects inside the clip volume.
func (f Frustum) PointVisible(x, y, z float32) bool {
	return clipPointInsideNDC(f.clipM, x, y, z)
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
