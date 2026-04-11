package jolt

// #include "wrapper/body.h"
import "C"

// MotionType determines how a body responds to forces
type MotionType int

const (
	MotionTypeStatic    MotionType = C.JoltMotionTypeStatic    // Immovable, zero velocity
	MotionTypeKinematic MotionType = C.JoltMotionTypeKinematic // Movable by user, doesn't respond to forces
	MotionTypeDynamic   MotionType = C.JoltMotionTypeDynamic   // Affected by forces
)

// AllowedDOFsOrZero passes 0 to Jolt body creation to keep EAllowedDOFs::All.
// AllowedDOFsPlatformer is world-space TranslationX|Y|Z + RotationY (upright capsule / character).
const (
	AllowedDOFsOrZero     int = 0
	AllowedDOFsPlatformer int = 0x17 // 1+2+4+16; see Jolt AllowedDOFs.h
)

// BodyInterface provides methods to create and manipulate physics bodies
type BodyInterface struct {
	handle C.JoltBodyInterface
}

// GetBodyInterface returns the interface for creating/manipulating bodies
func (ps *PhysicsSystem) GetBodyInterface() *BodyInterface {
	handle := C.JoltPhysicsSystemGetBodyInterface(ps.handle)
	return &BodyInterface{handle: handle}
}

// BodyID uniquely identifies a physics body
type BodyID struct {
	handle C.JoltBodyID
}

// Destroy frees the body ID
func (b *BodyID) Destroy() {
	C.JoltDestroyBodyID(b.handle)
}

// GetPosition returns the current position of a body
func (bi *BodyInterface) GetPosition(bodyID *BodyID) Vec3 {
	var x, y, z C.float
	C.JoltGetBodyPosition(bi.handle, bodyID.handle, &x, &y, &z)
	return Vec3{
		X: float32(x),
		Y: float32(y),
		Z: float32(z),
	}
}

// GetRotation returns the body's world-space orientation (Jolt quaternion: x, y, z, w).
func (bi *BodyInterface) GetRotation(bodyID *BodyID) Quat {
	var x, y, z, w C.float
	C.JoltGetBodyRotation(bi.handle, bodyID.handle, &x, &y, &z, &w)
	return Quat{
		X: float32(x),
		Y: float32(y),
		Z: float32(z),
		W: float32(w),
	}
}

// CreateBody creates a body with specific motion type and sensor flag.
//
// Parameters:
//   - shape: The collision shape
//   - position: Initial position
//   - motionType: MotionTypeStatic, MotionTypeKinematic, or MotionTypeDynamic
//   - isSensor: If true, body is detected by queries but doesn't generate contact forces
//
// Examples:
//
//	// Create static ground
//	box := jolt.CreateBox(jolt.Vec3{X: 10, Y: 0.5, Z: 10})
//	ground := bi.CreateBody(box, jolt.Vec3{X: 0, Y: 0, Z: 0}, jolt.MotionTypeStatic, false)
//
//	// Create dynamic sphere
//	sphere := jolt.CreateSphere(1.0)
//	ball := bi.CreateBody(sphere, jolt.Vec3{X: 0, Y: 10, Z: 0}, jolt.MotionTypeDynamic, false)
//	bi.ActivateBody(ball)
//
//	// Create kinematic sensor
//	capsule := jolt.CreateCapsule(0.5, 1.8)
//	sensor := bi.CreateBody(capsule, jolt.Vec3{X: 0, Y: 1, Z: 0}, jolt.MotionTypeKinematic, true, 0.2, 0, jolt.AllowedDOFsOrZero)
//	bi.ActivateBody(sensor)
func (bi *BodyInterface) CreateBody(shape *Shape, position Vec3, motionType MotionType, isSensor bool, friction, restitution float32, allowedDOFsOrZero int) *BodyID {
	sensor := C.int(0)
	if isSensor {
		sensor = C.int(1)
	}

	handle := C.JoltCreateBody(
		bi.handle,
		shape.handle,
		C.float(position.X),
		C.float(position.Y),
		C.float(position.Z),
		C.JoltMotionType(motionType),
		sensor,
		C.float(friction),
		C.float(restitution),
		C.int(allowedDOFsOrZero),
	)

	return &BodyID{handle: handle}
}

// SetFriction sets the friction coefficient on an existing body (0 = no friction).
func (bi *BodyInterface) SetFriction(bodyID *BodyID, friction float32) {
	C.JoltSetBodyFriction(bi.handle, bodyID.handle, C.float(friction))
}

// SetRestitution sets restitution on an existing body (0 = no bounce).
func (bi *BodyInterface) SetRestitution(bodyID *BodyID, restitution float32) {
	C.JoltSetBodyRestitution(bi.handle, bodyID.handle, C.float(restitution))
}

// SetPosition updates the position of a body
func (bi *BodyInterface) SetPosition(bodyID *BodyID, position Vec3) {
	C.JoltSetBodyPosition(
		bi.handle,
		bodyID.handle,
		C.float(position.X),
		C.float(position.Y),
		C.float(position.Z),
	)
}

// GetLinearVelocity returns the linear velocity of a body (world space).
func (bi *BodyInterface) GetLinearVelocity(bodyID *BodyID) Vec3 {
	var x, y, z C.float
	C.JoltGetBodyLinearVelocity(bi.handle, bodyID.handle, &x, &y, &z)
	return Vec3{X: float32(x), Y: float32(y), Z: float32(z)}
}

// SetLinearVelocity sets the linear velocity of a body (world space).
func (bi *BodyInterface) SetLinearVelocity(bodyID *BodyID, velocity Vec3) {
	C.JoltSetBodyLinearVelocity(
		bi.handle,
		bodyID.handle,
		C.float(velocity.X),
		C.float(velocity.Y),
		C.float(velocity.Z),
	)
}

// AddImpulse applies an impulse to a dynamic body (world space).
func (bi *BodyInterface) AddImpulse(bodyID *BodyID, impulse Vec3) {
	C.JoltAddBodyImpulse(
		bi.handle,
		bodyID.handle,
		C.float(impulse.X),
		C.float(impulse.Y),
		C.float(impulse.Z),
	)
}

// ActivateBody makes a body participate in the simulation
func (bi *BodyInterface) ActivateBody(bodyID *BodyID) {
	C.JoltActivateBody(bi.handle, bodyID.handle)
}

// DeactivateBody removes a body from active simulation
func (bi *BodyInterface) DeactivateBody(bodyID *BodyID) {
	C.JoltDeactivateBody(bi.handle, bodyID.handle)
}

// SetShape changes the collision shape of a body
//
// Parameters:
//   - bodyID: The body to modify
//   - shape: The new collision shape
//   - updateMassProperties: If true, recalculates mass/inertia from the new shape
//
// Note: This automatically activates the body
func (bi *BodyInterface) SetShape(bodyID *BodyID, shape *Shape, updateMassProperties bool) {
	update := C.int(0)
	if updateMassProperties {
		update = C.int(1)
	}
	C.JoltSetBodyShape(bi.handle, bodyID.handle, shape.handle, update)
}
