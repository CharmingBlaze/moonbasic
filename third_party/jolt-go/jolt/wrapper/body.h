/*
 * Jolt Physics C Wrapper - Body Operations
 *
 * Handles rigid body creation and manipulation.
 */

#ifndef JOLT_WRAPPER_BODY_H
#define JOLT_WRAPPER_BODY_H

#include "physics.h"

#ifdef __cplusplus
extern "C" {
#endif

// Opaque pointer types
typedef void* JoltBodyInterface;
typedef void* JoltBodyID;
typedef void* JoltShape;

// Motion type enum (matches Jolt's EMotionType)
typedef enum {
    JoltMotionTypeStatic = 0,    // Immovable, zero velocity
    JoltMotionTypeKinematic = 1, // Movable by user, zero velocity response to forces
    JoltMotionTypeDynamic = 2    // Affected by forces
} JoltMotionType;

// Motion quality enum (matches Jolt's EMotionQuality)
typedef enum {
    JoltMotionQualityDiscrete = 0,    // Discrete collision detection
    JoltMotionQualityLinearCast = 1,  // Continuous collision detection (CCD)
} JoltMotionQuality;

// Get the body interface for creating/manipulating bodies
JoltBodyInterface JoltPhysicsSystemGetBodyInterface(JoltPhysicsSystem system);

// Get the position of a body
void JoltGetBodyPosition(const JoltBodyInterface bodyInterface,
                        const JoltBodyID bodyID,
                        float* x, float* y, float* z);

// Get the rotation of a body (world-space quaternion x, y, z, w - Jolt convention)
void JoltGetBodyRotation(const JoltBodyInterface bodyInterface,
                        const JoltBodyID bodyID,
                        float* x, float* y, float* z, float* w);

// Set the position of a body
void JoltSetBodyPosition(JoltBodyInterface bodyInterface,
                        JoltBodyID bodyID,
                        float x, float y, float z);

// Get linear velocity (world space, units/s)
void JoltGetBodyLinearVelocity(const JoltBodyInterface bodyInterface,
                               const JoltBodyID bodyID,
                               float* x, float* y, float* z);

// Set linear velocity (world space, units/s)
void JoltSetBodyLinearVelocity(JoltBodyInterface bodyInterface,
                               JoltBodyID bodyID,
                               float x, float y, float z);

// Apply impulse (world space, kg?m/s for dynamic bodies)
void JoltAddBodyImpulse(JoltBodyInterface bodyInterface,
                        JoltBodyID bodyID,
                        float x, float y, float z);

// Create a body with motion type, sensor flag, friction/restitution, and optional DOF lock.
// allowedDOFsOrZero: 0 = all DOFs (Jolt default). Non-zero = EAllowedDOFs bitmask (see Jolt AllowedDOFs.h).
// objectLayerOrMinusOne: -1 = derive from motion (static vs moving/character) and isSensor (sensor layer);
//   otherwise 0..3 = explicit object layer (NON_MOVING, MOVING, CHARACTER, SENSOR).
JoltBodyID JoltCreateBody(JoltBodyInterface bodyInterface,
                          JoltShape shape,
                          float x, float y, float z,
                          JoltMotionType motionType,
                          int isSensor,
                          float friction,
                          float restitution,
                          int allowedDOFsOrZero,
                          int objectLayerOrMinusOne);

// Batched pose read: one CGO crossing. Writes 8 floats per body:
// px, py, pz, qx, qy, qz, qw, active(1 = body active / not sleeping, 0 otherwise).
void JoltBatchGetBodyTransforms(JoltBodyInterface bodyInterface,
                                const JoltBodyID* bodyIDs,
                                int count,
                                float* out);

// Add (dvx,dvy,dvz) to linear velocity for each body in one call (e.g. per-step gravity).
void JoltBatchApplyGravityDelta(JoltBodyInterface bodyInterface,
                                const JoltBodyID* bodyIDs,
                                int count,
                                float dvx, float dvy, float dvz);

// Runtime material tweaks (BodyInterface)
void JoltSetBodyFriction(JoltBodyInterface bodyInterface,
                         JoltBodyID bodyID,
                         float friction);

void JoltSetBodyRestitution(JoltBodyInterface bodyInterface,
                            JoltBodyID bodyID,
                            float restitution);

// Set damping (linear and angular)
void JoltSetBodyDamping(JoltPhysicsSystem system,
                       JoltBodyID bodyID,
                       float linear,
                       float angular);

// Set motion quality (CCD toggle)
void JoltSetBodyMotionQuality(JoltBodyInterface bodyInterface,
                             JoltBodyID bodyID,
                             JoltMotionQuality quality);

// Set allowed degrees of freedom (Lock Axis)
// This requires a lock and potentially re-calculating motion properties.
void JoltSetBodyAllowedDOFs(JoltPhysicsSystem system,
                           JoltBodyID bodyID,
                           int allowedDOFs);

// Set gravity factor (0 = weightless, 1 = normal)
void JoltSetBodyGravityFactor(JoltBodyInterface bodyInterface,
                             JoltBodyID bodyID,
                             float gravityFactor);

// Set sensor flag
void JoltSetBodyIsSensor(JoltBodyInterface bodyInterface,
                        JoltBodyID bodyID,
                        int isSensor);

// Activate a body (makes it participate in simulation)
void JoltActivateBody(JoltBodyInterface bodyInterface, JoltBodyID bodyID);

// Deactivate a body (removes from active simulation)
void JoltDeactivateBody(JoltBodyInterface bodyInterface, JoltBodyID bodyID);

// Set the shape of a body
void JoltSetBodyShape(JoltBodyInterface bodyInterface,
                     JoltBodyID bodyID,
                     JoltShape shape,
                     int updateMassProperties);

// Index+sequence packing (matches JPH::BodyID::GetIndexAndSequenceNumber) for contact/event lookup.
unsigned int JoltBodyIDGetIndexAndSequenceNumber(JoltBodyID bodyID);

// Destroy a body ID
void JoltDestroyBodyID(JoltBodyID bodyID);

#ifdef __cplusplus
}
#endif

#endif // JOLT_WRAPPER_BODY_H
