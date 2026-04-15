/*
 * Jolt Physics C Wrapper - Constraints/Joints Implementation
 */

#include "joint.h"
#include <Jolt/Jolt.h>
#include <Jolt/Physics/PhysicsSystem.h>
#include <Jolt/Physics/Constraints/HingeConstraint.h>
#include <Jolt/Physics/Body/Body.h>
#include <Jolt/Physics/Body/BodyLockMulti.h>
#include <cmath>

using namespace JPH;

// Defined in physics.cpp (opaque here so this TU does not depend on PhysicsSystemWrapper layout)
struct PhysicsSystemWrapper;
PhysicsSystem* GetPhysicsSystem(PhysicsSystemWrapper* wrapper);

JoltConstraint JoltCreateHingeJoint(JoltPhysicsSystem system,
                                   JoltBodyID body1,
                                   JoltBodyID body2,
                                   float px, float py, float pz,
                                   float ax, float ay, float az)
{
    PhysicsSystemWrapper *wrapper = static_cast<PhysicsSystemWrapper *>(system);
    PhysicsSystem* ps = GetPhysicsSystem(wrapper);

    const BodyID *bid1 = static_cast<const BodyID *>(body1);
    const BodyID *bid2 = static_cast<const BodyID *>(body2);

    HingeConstraintSettings settings;
    settings.mSpace = EConstraintSpace::WorldSpace;
    settings.mPoint1 = settings.mPoint2 = RVec3(px, py, pz);
    settings.mHingeAxis1 = settings.mHingeAxis2 = Vec3(ax, ay, az);

    // Normal axis must be perpendicular to hinge axis
    Vec3 normal = Vec3::sAxisX();
    if (std::fabs(ax) > 0.9f) {
        normal = Vec3::sAxisY();
    }
    settings.mNormalAxis1 = settings.mNormalAxis2 =
        settings.mHingeAxis1.GetCrossProduct(normal).NormalizedOr(Vec3::sAxisZ());

    TwoBodyConstraint* constraint = nullptr;
    const BodyID ids[2] = { *bid1, *bid2 };
    {
        BodyLockMultiWrite lock(ps->GetBodyLockInterface(), ids, 2);
        Body* b0 = lock.GetBody(0);
        Body* b1 = lock.GetBody(1);
        if (b0 != nullptr && b1 != nullptr) {
            constraint = settings.Create(*b0, *b1);
            if (constraint != nullptr) {
                ps->AddConstraint(constraint);
            }
        }
    }

    return static_cast<JoltConstraint>(constraint);
}

void JoltDestroyConstraint(JoltPhysicsSystem system, JoltConstraint constraint)
{
    PhysicsSystemWrapper *wrapper = static_cast<PhysicsSystemWrapper *>(system);
    PhysicsSystem* ps = GetPhysicsSystem(wrapper);
    Constraint* c = static_cast<Constraint*>(constraint);

    if (c) {
        ps->RemoveConstraint(c);
    }
}
