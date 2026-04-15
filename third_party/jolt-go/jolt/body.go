package jolt

// #include "wrapper/body.h"
import "C"
import "unsafe"

// MotionType determines how a body responds to forces
type MotionType int

const (
	MotionTypeStatic    MotionType = C.JoltMotionTypeStatic    // Immovable, zero velocity
	MotionTypeKinematic MotionType = C.JoltMotionTypeKinematic // Movable by user, doesn't respond to forces
	MotionTypeDynamic   MotionType = C.JoltMotionTypeDynamic   // Affected by forces
)

// MotionQuality determines the depth of collision detection
type MotionQuality int

const (
	MotionQualityDiscrete   MotionQuality = C.JoltMotionQualityDiscrete   // Discrete collision detection
	MotionQualityLinearCast MotionQuality = C.JoltMotionQualityLinearCast // Continuous collision detection (CCD)
)

// AllowedDOFs passes bits to Jolt to lock rotation/translation axes.
const (
	AllowedDOFsAll        int = 0
	AllowedDOFsTranslationX int = 1
	AllowedDOFsTranslationY int = 2
	AllowedDOFsTranslationZ int = 4
	AllowedDOFsRotationX    int = 8
	AllowedDOFsRotationY    int = 16
	AllowedDOFsRotationZ    int = 32
	AllowedDOFsPlaneXZ      int = 1 | 4 | 16 // Transl X/Z, Rot Y (2D Platformer)
)

// Object layer IDs (must match physics_layers.h / wrapper). Use ObjectLayerAuto to derive from motion + sensor.
const (
	ObjectLayerAuto      int = -1
	ObjectLayerNonMoving int = 0
	ObjectLayerMoving    int = 1
	ObjectLayerCharacter int = 2
	ObjectLayerSensor    int = 3
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

// IndexAndSequenceNumber returns Jolt's packed body id (matches CharacterContactEvent.BodyB).
func (b *BodyID) IndexAndSequenceNumber() uint32 {
	if b == nil {
		return 0
	}
	return uint32(C.JoltBodyIDGetIndexAndSequenceNumber(b.handle))
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
func (bi *BodyInterface) CreateBody(shape *Shape, position Vec3, motionType MotionType, isSensor bool, friction, restitution float32, allowedDOFsOrZero int, objectLayerOrMinusOne ...int) *BodyID {
	sensor := C.int(0)
	if isSensor {
		sensor = C.int(1)
	}
	oly := C.int(-1)
	if len(objectLayerOrMinusOne) > 0 {
		oly = C.int(objectLayerOrMinusOne[0])
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
		oly,
	)

	return &BodyID{handle: handle}
}

// BatchGetBodyTransforms fills out with 8 floats per body: px,py,pz, qx,qy,qz,qw, active(1/0).
func (bi *BodyInterface) BatchGetBodyTransforms(ids []*BodyID, out []float32) {
	if len(ids) == 0 {
		return
	}
	if len(out) < len(ids)*8 {
		return
	}
	cIDs := make([]C.JoltBodyID, len(ids))
	for i := range ids {
		cIDs[i] = ids[i].handle
	}
	C.JoltBatchGetBodyTransforms(bi.handle, (*C.JoltBodyID)(unsafe.Pointer(&cIDs[0])), C.int(len(ids)), (*C.float)(unsafe.Pointer(&out[0])))
}

// BatchApplyGravityDelta adds (dvx,dvy,dvz) to each body's linear velocity in one CGO call.
func (bi *BodyInterface) BatchApplyGravityDelta(ids []*BodyID, dvx, dvy, dvz float32) {
	if len(ids) == 0 {
		return
	}
	cIDs := make([]C.JoltBodyID, len(ids))
	for i := range ids {
		cIDs[i] = ids[i].handle
	}
	C.JoltBatchApplyGravityDelta(bi.handle, (*C.JoltBodyID)(unsafe.Pointer(&cIDs[0])), C.int(len(ids)), C.float(dvx), C.float(dvy), C.float(dvz))
}

// SetFriction sets the friction coefficient on an existing body (0 = no friction).
func (bi *BodyInterface) SetFriction(bodyID *BodyID, friction float32) {
	C.JoltSetBodyFriction(bi.handle, bodyID.handle, C.float(friction))
}

// SetRestitution sets restitution on an existing body (0 = no bounce).
func (bi *BodyInterface) SetRestitution(bodyID *BodyID, restitution float32) {
	C.JoltSetBodyRestitution(bi.handle, bodyID.handle, C.float(restitution))
}

// SetBodyDamping sets linear and angular damping on an existing body.
func (ps *PhysicsSystem) SetBodyDamping(bodyID *BodyID, linear, angular float32) {
	C.JoltSetBodyDamping(ps.handle, bodyID.handle, C.float(linear), C.float(angular))
}

// SetMotionQuality sets motion quality for a body (e.g., enable CCD).
func (bi *BodyInterface) SetMotionQuality(bodyID *BodyID, quality MotionQuality) {
    C.JoltSetBodyMotionQuality(bi.handle, bodyID.handle, C.JoltMotionQuality(quality))
}

// SetAllowedDOFs sets allowed degrees of freedom (Lock Axis). Requires the system for locking.
func (ps *PhysicsSystem) SetAllowedDOFs(bodyID *BodyID, allowedDOFs int) {
    C.JoltSetBodyAllowedDOFs(ps.handle, bodyID.handle, C.int(allowedDOFs))
}

// SetGravityFactor sets the gravity multiplier for a body.
func (bi *BodyInterface) SetGravityFactor(bodyID *BodyID, factor float32) {
    C.JoltSetBodyGravityFactor(bi.handle, bodyID.handle, C.float(factor))
}

// SetIsSensor toggles the sensor flag on a body.
func (bi *BodyInterface) SetIsSensor(bodyID *BodyID, isSensor bool) {
    sensor := C.int(0)
    if isSensor { sensor = C.int(1) }
    C.JoltSetBodyIsSensor(bi.handle, bodyID.handle, sensor)
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
func (bi *BodyInterface) SetShape(bodyID *BodyID, shape *Shape, updateMassProperties bool) {
	update := C.int(0)
	if updateMassProperties {
		update = C.int(1)
	}
	C.JoltSetBodyShape(bi.handle, bodyID.handle, shape.handle, update)
}
