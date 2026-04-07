package jolt

// #include "wrapper/character.h"
import "C"

// BackFaceMode controls how the character collides with back faces
type BackFaceMode int

const (
	// BackFaceModeIgnore - Ignore all back facing surfaces
	BackFaceModeIgnore BackFaceMode = 0
	// BackFaceModeCollide - Collide with back facing surfaces
	BackFaceModeCollide BackFaceMode = 1
)

// CharacterContact represents a collision contact for a virtual character
type CharacterContact struct {
	// Position where the character makes contact
	Position Vec3
	// LinearVelocity is the velocity of the contact point
	LinearVelocity Vec3
	// ContactNormal is the contact normal, pointing towards the character
	ContactNormal Vec3
	// SurfaceNormal is the surface normal of the contact
	SurfaceNormal Vec3
	// Distance to the contact (<= 0 means actual contact, > 0 means predictive)
	Distance float32
	// Fraction along the path where this contact takes place
	Fraction float32
	// BodyB is the ID of the body we're colliding with (nil if invalid)
	BodyB *BodyID
	// UserData is the user data of the body
	UserData uint64
	// IsSensorB indicates if the body is a sensor
	IsSensorB bool
	// HadCollision indicates if the character actually collided with the contact
	HadCollision bool
	// WasDiscarded indicates if the contact was discarded
	WasDiscarded bool
	// CanPushCharacter indicates if the velocity of the contact point can push the character
	CanPushCharacter bool
}

// CharacterVirtualSettings configures a virtual character
type CharacterVirtualSettings struct {
	// Shape is the collision shape for the character (required)
	Shape *Shape

	// Up is the character's up direction (default: {0, 1, 0})
	Up Vec3

	// MaxSlopeAngle is the maximum slope angle in radians that the character can walk on (default: 50 degrees)
	MaxSlopeAngle float32

	// Mass is the character mass in kg, used to push down objects (default: 70.0)
	Mass float32

	// MaxStrength is the maximum force the character can push other bodies with in Newtons (default: 100.0)
	MaxStrength float32

	// ShapeOffset is an extra offset applied to the shape in local space (default: {0, 0, 0})
	ShapeOffset Vec3

	// BackFaceMode controls collision with back faces (default: BackFaceModeCollide)
	BackFaceMode BackFaceMode

	// PredictiveContactDistance is how far to scan outside the shape for contacts.
	// 0 will cause the character to get stuck. Too high causes ghost collisions. (default: 0.1)
	PredictiveContactDistance float32

	// MaxCollisionIterations is the max amount of collision loops (default: 5)
	MaxCollisionIterations uint32

	// MaxConstraintIterations is how often to try stepping in constraint solving (default: 15)
	MaxConstraintIterations uint32

	// MinTimeRemaining is the early out condition - if this much time is left, we're done (default: 1.0e-4)
	MinTimeRemaining float32

	// CollisionTolerance is how far we're willing to penetrate geometry (default: 1.0e-3)
	CollisionTolerance float32

	// CharacterPadding is how far we try to stay away from geometry (default: 0.02)
	CharacterPadding float32

	// MaxNumHits is the max number of hits to collect to avoid excess contact points (default: 256)
	MaxNumHits uint32

	// HitReductionCosMaxAngle is cos(angle) for merging similar contact normals.
	// Default is ~2.5 degrees (0.999). Set to -1 to turn off. (default: 0.999)
	HitReductionCosMaxAngle float32

	// PenetrationRecoverySpeed governs how fast penetration is resolved.
	// 0 = nothing resolved, 1 = everything in one update (default: 1.0)
	PenetrationRecoverySpeed float32

	// EnhancedInternalEdgeRemoval removes ghost contacts with internal mesh edges.
	// More expensive but smoother movement over convex edges. (default: false)
	EnhancedInternalEdgeRemoval bool
}

// NewCharacterVirtualSettings creates settings with Jolt's default values
func NewCharacterVirtualSettings(shape *Shape) *CharacterVirtualSettings {
	return &CharacterVirtualSettings{
		Shape:                       shape,
		Up:                          Vec3{X: 0, Y: 1, Z: 0},
		MaxSlopeAngle:               DegreesToRadians(50.0),
		Mass:                        70.0,
		MaxStrength:                 100.0,
		ShapeOffset:                 Vec3{X: 0, Y: 0, Z: 0},
		BackFaceMode:                BackFaceModeCollide,
		PredictiveContactDistance:   0.1,
		MaxCollisionIterations:      5,
		MaxConstraintIterations:     15,
		MinTimeRemaining:            1.0e-4,
		CollisionTolerance:          1.0e-3,
		CharacterPadding:            0.02,
		MaxNumHits:                  256,
		HitReductionCosMaxAngle:     0.999,
		PenetrationRecoverySpeed:    1.0,
		EnhancedInternalEdgeRemoval: false,
	}
}

// CharacterVirtual represents a virtual character in the physics world
type CharacterVirtual struct {
	handle C.JoltCharacterVirtual
	ps     *PhysicsSystem
}

// GroundState indicates the ground contact state of a CharacterVirtual
type GroundState int

const (
	// GroundStateOnGround - Character is on the ground and can move freely
	GroundStateOnGround GroundState = 0
	// GroundStateOnSteepGround - Character is on a slope too steep to climb
	GroundStateOnSteepGround GroundState = 1
	// GroundStateNotSupported - Character is touching an object but not supported (should fall)
	GroundStateNotSupported GroundState = 2
	// GroundStateInAir - Character is in the air and not touching anything
	GroundStateInAir GroundState = 3
)

func (gs GroundState) String() string {
	switch gs {
	case GroundStateOnGround:
		return "OnGround"
	case GroundStateOnSteepGround:
		return "OnSteepGround"
	case GroundStateNotSupported:
		return "NotSupported"
	case GroundStateInAir:
		return "InAir"
	default:
		return "Unknown"
	}
}

// CreateCharacterVirtual creates a virtual character with the specified settings at the initial position
func (ps *PhysicsSystem) CreateCharacterVirtual(settings *CharacterVirtualSettings, position Vec3) *CharacterVirtual {
	// Convert Go settings to C settings
	cSettings := C.JoltCharacterVirtualSettings{
		shape:                       settings.Shape.handle,
		upX:                         C.float(settings.Up.X),
		upY:                         C.float(settings.Up.Y),
		upZ:                         C.float(settings.Up.Z),
		maxSlopeAngle:               C.float(settings.MaxSlopeAngle),
		mass:                        C.float(settings.Mass),
		maxStrength:                 C.float(settings.MaxStrength),
		shapeOffsetX:                C.float(settings.ShapeOffset.X),
		shapeOffsetY:                C.float(settings.ShapeOffset.Y),
		shapeOffsetZ:                C.float(settings.ShapeOffset.Z),
		backFaceMode:                C.JoltBackFaceMode(settings.BackFaceMode),
		predictiveContactDistance:   C.float(settings.PredictiveContactDistance),
		maxCollisionIterations:      C.uint(settings.MaxCollisionIterations),
		maxConstraintIterations:     C.uint(settings.MaxConstraintIterations),
		minTimeRemaining:            C.float(settings.MinTimeRemaining),
		collisionTolerance:          C.float(settings.CollisionTolerance),
		characterPadding:            C.float(settings.CharacterPadding),
		maxNumHits:                  C.uint(settings.MaxNumHits),
		hitReductionCosMaxAngle:     C.float(settings.HitReductionCosMaxAngle),
		penetrationRecoverySpeed:    C.float(settings.PenetrationRecoverySpeed),
		enhancedInternalEdgeRemoval: 0,
	}
	if settings.EnhancedInternalEdgeRemoval {
		cSettings.enhancedInternalEdgeRemoval = 1
	}

	handle := C.JoltCreateCharacterVirtual(
		ps.handle,
		&cSettings,
		C.float(position.X),
		C.float(position.Y),
		C.float(position.Z),
	)
	return &CharacterVirtual{handle: handle, ps: ps}
}

// Update advances the character simulation using the current velocity
// This is the basic update that moves the character according to its velocity and handles collisions.
// Note: You must apply gravity to the velocity yourself before calling this.
// deltaTime: duration of simulation step in seconds
// gravity: acceleration vector (e.g., Vec3{0, -9.81, 0} for Earth gravity) - applied when standing on objects
func (cv *CharacterVirtual) Update(deltaTime float32, gravity Vec3) {
	C.JoltCharacterVirtualUpdate(
		cv.handle,
		cv.ps.handle,
		C.float(deltaTime),
		C.float(gravity.X),
		C.float(gravity.Y),
		C.float(gravity.Z),
	)
}

// ExtendedUpdate advances the character simulation with combined movement logic
// Combines Update, StickToFloor, and WalkStairs into a unified operation
// deltaTime: duration of simulation step in seconds
// gravity: acceleration vector (e.g., Vec3{0, -9.81, 0} for Earth gravity)
func (cv *CharacterVirtual) ExtendedUpdate(deltaTime float32, gravity Vec3) {
	C.JoltCharacterVirtualExtendedUpdate(
		cv.handle,
		cv.ps.handle,
		C.float(deltaTime),
		C.float(gravity.X),
		C.float(gravity.Y),
		C.float(gravity.Z),
	)
}

// SetLinearVelocity sets the character's linear velocity
func (cv *CharacterVirtual) SetLinearVelocity(velocity Vec3) {
	C.JoltCharacterVirtualSetLinearVelocity(
		cv.handle,
		C.float(velocity.X),
		C.float(velocity.Y),
		C.float(velocity.Z),
	)
}

// GetLinearVelocity returns the current linear velocity of the character
func (cv *CharacterVirtual) GetLinearVelocity() Vec3 {
	var x, y, z C.float
	C.JoltCharacterVirtualGetLinearVelocity(cv.handle, &x, &y, &z)
	return Vec3{
		X: float32(x),
		Y: float32(y),
		Z: float32(z),
	}
}

// GetGroundVelocity returns the velocity clamped to the ground plane
func (cv *CharacterVirtual) GetGroundVelocity() Vec3 {
	var x, y, z C.float
	C.JoltCharacterVirtualGetGroundVelocity(cv.handle, &x, &y, &z)
	return Vec3{
		X: float32(x),
		Y: float32(y),
		Z: float32(z),
	}
}

// SetPosition sets the character's position in the world
func (cv *CharacterVirtual) SetPosition(position Vec3) {
	C.JoltCharacterVirtualSetPosition(
		cv.handle,
		C.float(position.X),
		C.float(position.Y),
		C.float(position.Z),
	)
}

// GetPosition returns the current position of the character
func (cv *CharacterVirtual) GetPosition() Vec3 {
	var x, y, z C.float
	C.JoltCharacterVirtualGetPosition(cv.handle, &x, &y, &z)
	return Vec3{
		X: float32(x),
		Y: float32(y),
		Z: float32(z),
	}
}

// Destroy frees the character resources
func (cv *CharacterVirtual) Destroy() {
	C.JoltDestroyCharacterVirtual(cv.handle)
}

// GetGroundState returns the current ground contact state
func (cv *CharacterVirtual) GetGroundState() GroundState {
	state := C.JoltCharacterVirtualGetGroundState(cv.handle)
	return GroundState(state)
}

// IsSupported returns true if the character is on ground or steep ground (not falling)
func (cv *CharacterVirtual) IsSupported() bool {
	result := C.JoltCharacterVirtualIsSupported(cv.handle)
	return result != 0
}

// SetShape changes the collision shape of the character
// shape: new collision shape for the character
// maxPenetrationDepth: maximum allowed penetration depth (typically 0.1)
func (cv *CharacterVirtual) SetShape(shape *Shape, maxPenetrationDepth float32) {
	C.JoltCharacterVirtualSetShape(
		cv.handle,
		shape.handle,
		C.float(maxPenetrationDepth),
		cv.ps.handle,
	)
}

// GetShape retrieves the current collision shape of the character
func (cv *CharacterVirtual) GetShape() *Shape {
	handle := C.JoltCharacterVirtualGetShape(cv.handle)
	return &Shape{handle: handle}
}

// PhysicsSystem returns the physics system that this character belongs to
func (cv *CharacterVirtual) PhysicsSystem() *PhysicsSystem {
	return cv.ps
}

// GetGroundNormal returns the normal vector of the ground surface
func (cv *CharacterVirtual) GetGroundNormal() Vec3 {
	var x, y, z C.float
	C.JoltCharacterVirtualGetGroundNormal(cv.handle, &x, &y, &z)
	return Vec3{
		X: float32(x),
		Y: float32(y),
		Z: float32(z),
	}
}

// GetGroundPosition returns the world position of the ground contact point
func (cv *CharacterVirtual) GetGroundPosition() Vec3 {
	var x, y, z C.float
	C.JoltCharacterVirtualGetGroundPosition(cv.handle, &x, &y, &z)
	return Vec3{
		X: float32(x),
		Y: float32(y),
		Z: float32(z),
	}
}

// GetActiveContacts returns the list of active contacts for the character
// maxContacts specifies the maximum number of contacts to retrieve (typically 256)
func (cv *CharacterVirtual) GetActiveContacts(maxContacts int) []CharacterContact {
	// Allocate C array for contacts
	cContacts := make([]C.JoltCharacterContact, maxContacts)

	// Call C function
	numContacts := int(C.JoltCharacterVirtualGetActiveContacts(
		cv.handle,
		&cContacts[0],
		C.int(maxContacts),
	))

	// Convert C contacts to Go contacts
	contacts := make([]CharacterContact, numContacts)
	for i := 0; i < numContacts; i++ {
		c := &cContacts[i]

		var bodyB *BodyID
		if c.bodyB != nil {
			bodyB = &BodyID{handle: c.bodyB}
		}

		contacts[i] = CharacterContact{
			Position: Vec3{
				X: float32(c.positionX),
				Y: float32(c.positionY),
				Z: float32(c.positionZ),
			},
			LinearVelocity: Vec3{
				X: float32(c.linearVelocityX),
				Y: float32(c.linearVelocityY),
				Z: float32(c.linearVelocityZ),
			},
			ContactNormal: Vec3{
				X: float32(c.contactNormalX),
				Y: float32(c.contactNormalY),
				Z: float32(c.contactNormalZ),
			},
			SurfaceNormal: Vec3{
				X: float32(c.surfaceNormalX),
				Y: float32(c.surfaceNormalY),
				Z: float32(c.surfaceNormalZ),
			},
			Distance:         float32(c.distance),
			Fraction:         float32(c.fraction),
			BodyB:            bodyB,
			UserData:         uint64(c.userData),
			IsSensorB:        c.isSensorB != 0,
			HadCollision:     c.hadCollision != 0,
			WasDiscarded:     c.wasDiscarded != 0,
			CanPushCharacter: c.canPushCharacter != 0,
		}
	}

	return contacts
}
