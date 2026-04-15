/*
 * Jolt Physics C Wrapper - Constraints/Joints
 */

#ifndef JOLT_WRAPPER_JOINT_H
#define JOLT_WRAPPER_JOINT_H

#include "physics.h"
#include "body.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef void* JoltConstraint;

// Create a hinge joint (revolute joint)
// px, py, pz: pivot point in world space
// ax, ay, az: rotation axis in world space
JoltConstraint JoltCreateHingeJoint(JoltPhysicsSystem system,
                                   JoltBodyID body1,
                                   JoltBodyID body2,
                                   float px, float py, float pz,
                                   float ax, float ay, float az);

// Destroy a constraint
void JoltDestroyConstraint(JoltPhysicsSystem system, JoltConstraint constraint);

#ifdef __cplusplus
}
#endif

#endif // JOLT_WRAPPER_JOINT_H
