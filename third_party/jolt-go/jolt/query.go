package jolt

// #include "wrapper/query.h"
import "C"

// CollisionHit contains information about a single collision detected during a shape query
type CollisionHit struct {
	BodyID           *BodyID // The body that was hit
	ContactPoint     Vec3    // The contact point in world space
	PenetrationDepth float32 // How deep the shapes overlap (negative if separated)
}

// RaycastHit contains information about a single raycast hit
type RaycastHit struct {
	BodyID   *BodyID // The body that was hit (nil if no hit)
	HitPoint Vec3    // The position where the ray hit the surface
	Normal   Vec3    // The surface normal at the hit point
	Fraction float32 // The fraction along the ray where the hit occurred [0, 1]
}

// CollideShape checks if a shape at the given position collides with any bodies in the physics system.
// This performs a static overlap test - the shape itself is not added to the physics system.
//
// Parameters:
//   - shape: The collision shape to test
//   - position: Position in world space to place the shape
//   - penetrationTolerance: Distance threshold for collision detection in meters
//
// Returns true if any collision is detected, false otherwise.
//
// Example usage:
//
//	sphere := jolt.CreateSphere(1.0)
//	defer sphere.Destroy()
//
//	if ps.CollideShape(sphere, jolt.Vec3{X: 0, Y: 5, Z: 0}, 0) {
//	    fmt.Println("Collision detected!")
//	}
func (ps *PhysicsSystem) CollideShape(shape *Shape, position Vec3, penetrationTolerance float32) bool {
	result := C.JoltCollideShape(
		ps.handle,
		shape.handle,
		C.float(position.X),
		C.float(position.Y),
		C.float(position.Z),
		C.float(penetrationTolerance),
	)
	return result != 0
}

// CollideShapeGetHits performs a shape collision query and returns detailed information about all hits.
// This is useful when you need to know which specific bodies were hit and where.
//
// Parameters:
//   - shape: The collision shape to test
//   - position: Position in world space to place the shape
//   - maxHits: Maximum number of hits to return (limits memory allocation)
//   - penetrationTolerance: A factor that determines the accuracy of the penetration depth calculation.
//     If the change of the squared distance is less than tolerance * current_penetration_depth^2 the algorithm will terminate. (unit: dimensionless)
//
// Returns a slice of CollisionHit containing information about each collision.
//
// Example usage:
//
//	sphere := jolt.CreateSphere(1.0)
//	defer sphere.Destroy()
//
//	hits := ps.CollideShapeGetHits(sphere, jolt.Vec3{X: 0, Y: 5, Z: 0}, 10, 0)
//	for _, hit := range hits {
//	    fmt.Printf("Hit body at %.2f, %.2f, %.2f (depth: %.2f)\n",
//	        hit.ContactPoint.X, hit.ContactPoint.Y, hit.ContactPoint.Z,
//	        hit.PenetrationDepth)
//	}
func (ps *PhysicsSystem) CollideShapeGetHits(shape *Shape, position Vec3, maxHits int, penetrationTolerance float32) []CollisionHit {
	if maxHits <= 0 {
		return []CollisionHit{}
	}

	// Allocate C array for results
	cHits := make([]C.JoltCollisionHit, maxHits)

	numHits := C.JoltCollideShapeGetHits(
		ps.handle,
		shape.handle,
		C.float(position.X),
		C.float(position.Y),
		C.float(position.Z),
		&cHits[0],
		C.int(maxHits),
		C.float(penetrationTolerance),
	)

	// Convert C results to Go
	hits := make([]CollisionHit, int(numHits))
	for i := 0; i < int(numHits); i++ {
		cHit := cHits[i]
		hits[i] = CollisionHit{
			BodyID: &BodyID{handle: cHit.bodyID},
			ContactPoint: Vec3{
				X: float32(cHit.contactPointX),
				Y: float32(cHit.contactPointY),
				Z: float32(cHit.contactPointZ),
			},
			PenetrationDepth: float32(cHit.penetrationDepth),
		}
	}

	return hits
}

// CastRay performs a raycast from origin in the specified direction and returns the closest hit.
// The direction vector does not need to be normalized - its length determines the maximum ray distance.
//
// Parameters:
//   - origin: Starting position of the ray in world space
//   - direction: Direction and length of the ray (ray goes from origin to origin + direction)
//
// Returns:
//   - hit: Information about the closest hit (HitPoint, Normal, Fraction, BodyID)
//   - hasHit: true if the ray hit something, false otherwise
//
// Example usage:
//
//	// Cast a ray from (0, 10, 0) downward for 20 units
//	origin := jolt.Vec3{X: 0, Y: 10, Z: 0}
//	direction := jolt.Vec3{X: 0, Y: -20, Z: 0}
//	hit, hasHit := ps.CastRay(origin, direction)
//	if hasHit {
//	    fmt.Printf("Hit at %.2f, %.2f, %.2f (fraction: %.2f)\n",
//	        hit.HitPoint.X, hit.HitPoint.Y, hit.HitPoint.Z, hit.Fraction)
//	}
func (ps *PhysicsSystem) CastRay(origin, direction Vec3) (RaycastHit, bool) {
	var cHit C.JoltRaycastHit

	result := C.JoltCastRay(
		ps.handle,
		C.float(origin.X),
		C.float(origin.Y),
		C.float(origin.Z),
		C.float(direction.X),
		C.float(direction.Y),
		C.float(direction.Z),
		&cHit,
	)

	if result == 0 {
		return RaycastHit{}, false
	}

	hit := RaycastHit{
		BodyID: &BodyID{handle: cHit.bodyID},
		HitPoint: Vec3{
			X: float32(cHit.hitPointX),
			Y: float32(cHit.hitPointY),
			Z: float32(cHit.hitPointZ),
		},
		Normal: Vec3{
			X: float32(cHit.normalX),
			Y: float32(cHit.normalY),
			Z: float32(cHit.normalZ),
		},
		Fraction: float32(cHit.fraction),
	}

	return hit, true
}

// CastRayGetHits performs a raycast and returns all hits along the ray, sorted by distance.
// The direction vector does not need to be normalized - its length determines the maximum ray distance.
//
// Parameters:
//   - origin: Starting position of the ray in world space
//   - direction: Direction and length of the ray (ray goes from origin to origin + direction)
//   - maxHits: Maximum number of hits to return (limits memory allocation)
//
// Returns a slice of RaycastHit containing all hits sorted by distance (closest first).
//
// Example usage:
//
//	// Cast a ray and get all hits
//	origin := jolt.Vec3{X: 0, Y: 10, Z: 0}
//	direction := jolt.Vec3{X: 0, Y: -20, Z: 0}
//	hits := ps.CastRayGetHits(origin, direction, 10)
//	for i, hit := range hits {
//	    fmt.Printf("Hit %d: body at %.2f, %.2f, %.2f (fraction: %.2f)\n",
//	        i, hit.HitPoint.X, hit.HitPoint.Y, hit.HitPoint.Z, hit.Fraction)
//	}
func (ps *PhysicsSystem) CastRayGetHits(origin, direction Vec3, maxHits int) []RaycastHit {
	if maxHits <= 0 {
		return []RaycastHit{}
	}

	// Allocate C array for results
	cHits := make([]C.JoltRaycastHit, maxHits)

	numHits := C.JoltCastRayGetHits(
		ps.handle,
		C.float(origin.X),
		C.float(origin.Y),
		C.float(origin.Z),
		C.float(direction.X),
		C.float(direction.Y),
		C.float(direction.Z),
		&cHits[0],
		C.int(maxHits),
	)

	// Convert C results to Go
	hits := make([]RaycastHit, int(numHits))
	for i := 0; i < int(numHits); i++ {
		cHit := cHits[i]
		hits[i] = RaycastHit{
			BodyID: &BodyID{handle: cHit.bodyID},
			HitPoint: Vec3{
				X: float32(cHit.hitPointX),
				Y: float32(cHit.hitPointY),
				Z: float32(cHit.hitPointZ),
			},
			Normal: Vec3{
				X: float32(cHit.normalX),
				Y: float32(cHit.normalY),
				Z: float32(cHit.normalZ),
			},
			Fraction: float32(cHit.fraction),
		}
	}

	return hits
}
