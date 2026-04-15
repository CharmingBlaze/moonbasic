/*
 * Jolt Physics C Wrapper - Constraints/Joints Implementation
 */

#include "joint.h"
#include <Jolt/Jolt.h>
#include <Jolt/Physics/PhysicsSystem.h>
#include <Jolt/Physics/Constraints/HingeConstraint.h>
#include <Jolt/Physics/Body/BodyInterface.h>

using namespace JPH;

JoltConstraint JoltCreateHingeJoint(JoltPhysicsSystem system,
                                   JoltBodyID body1,
                                   JoltBodyID body2,
                                   float px, float py, float pz,
                                   float ax, float ay, float az)
{
    PhysicsSystemWrapper *wrapper = static_cast<PhysicsSystemWrapper *>(system);
    PhysicsSystem* ps = wrapper->system.get();
    BodyInterface& bi = ps->GetBodyInterface();

    const BodyID *bid1 = static_cast<const BodyID *>(body1);
    const BodyID *bid2 = static_cast<const BodyID *>(body2);

    HingeConstraintSettings settings;
    settings.mSpace = EConstraintSpace::WorldSpace;
    settings.mPoint1 = settings.mPoint2 = RVec3(px, py, pz);
    settings.mHingeAxis1 = settings.mHingeAxis2 = Vec3(ax, ay, az);
    
    // Normal axis must be perpendicular to hinge axis
    Vec3 normal = Vec3::sAxisX();
    if (std::abs(ax) > 0.9f) normal = Vec3::sAxisY(); 
    settings.mNormalAxis1 = settings.mNormalAxis2 = settings.mHingeAxis1.GetCrossProduct(normal).NormalizedOr(Vec3::sAxisZ());

    Constraint* constraint = settings.Create(*ps->GetBodyLockInterface().LockBody(*bid1), *ps->GetBodyLockInterface().LockBody(*bid2));
    if (constraint) {
        ps->AddConstraint(constraint);
    }

    return static_cast<JoltConstraint>(constraint);
}

void JoltDestroyConstraint(JoltPhysicsSystem system, JoltConstraint constraint)
{
    PhysicsSystemWrapper *wrapper = static_cast<PhysicsSystemWrapper *>(system);
    PhysicsSystem* ps = wrapper->system.get();
    Constraint* c = static_cast<Constraint*>(constraint);

    if (c) {
        ps->RemoveConstraint(c);
        // Constraint is a RefCounted object in Jolt, but we don't have a Ref<> here.
        // ps->RemoveConstraint will decrement ref count.
    }
}
