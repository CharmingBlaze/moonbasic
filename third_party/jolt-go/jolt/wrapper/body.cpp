/*
 * Jolt Physics C Wrapper - Body Operations Implementation
 */

#include "body.h"
#include "physics.h"
#include "physics_layers.h"
#include <Jolt/Jolt.h>
#include <Jolt/Physics/PhysicsSystem.h>
#include <Jolt/Physics/Body/Body.h>
#include <Jolt/Physics/Body/AllowedDOFs.h>
#include <Jolt/Physics/Body/BodyCreationSettings.h>
#include <Jolt/Physics/Body/BodyInterface.h>
#include <Jolt/Physics/Body/BodyLockInterface.h>
#include <Jolt/Physics/Body/MotionProperties.h>
#include <memory>

using namespace JPH;

namespace Layers = MBPHYS_Layers;

JoltBodyInterface JoltPhysicsSystemGetBodyInterface(JoltPhysicsSystem system)
{
	PhysicsSystemWrapper *wrapper = static_cast<PhysicsSystemWrapper *>(system);
	PhysicsSystem* ps = GetPhysicsSystem(wrapper);
	BodyInterface* bi = &ps->GetBodyInterface();

	return static_cast<JoltBodyInterface>(bi);
}

void JoltGetBodyPosition(const JoltBodyInterface bodyInterface,
						 const JoltBodyID bodyID,
						 float *x, float *y, float *z)
{
	const BodyInterface *bi = static_cast<const BodyInterface *>(bodyInterface);
	const BodyID *bid = static_cast<const BodyID *>(bodyID);

	RVec3 pos = bi->GetPosition(*bid);
	*x = static_cast<float>(pos.GetX());
	*y = static_cast<float>(pos.GetY());
	*z = static_cast<float>(pos.GetZ());
}

void JoltGetBodyRotation(const JoltBodyInterface bodyInterface,
						 const JoltBodyID bodyID,
						 float *x, float *y, float *z, float *w)
{
	const BodyInterface *bi = static_cast<const BodyInterface *>(bodyInterface);
	const BodyID *bid = static_cast<const BodyID *>(bodyID);

	Quat q = bi->GetRotation(*bid);
	*x = q.GetX();
	*y = q.GetY();
	*z = q.GetZ();
	*w = q.GetW();
}

void JoltSetBodyPosition(JoltBodyInterface bodyInterface,
						 JoltBodyID bodyID,
						 float x, float y, float z)
{
	BodyInterface *bi = static_cast<BodyInterface *>(bodyInterface);
	const BodyID *bid = static_cast<const BodyID *>(bodyID);

	bi->SetPosition(*bid, RVec3(x, y, z), EActivation::DontActivate);
}

void JoltGetBodyLinearVelocity(const JoltBodyInterface bodyInterface,
							   const JoltBodyID bodyID,
							   float *x, float *y, float *z)
{
	const BodyInterface *bi = static_cast<const BodyInterface *>(bodyInterface);
	const BodyID *bid = static_cast<const BodyID *>(bodyID);

	Vec3 vel = bi->GetLinearVelocity(*bid);
	*x = vel.GetX();
	*y = vel.GetY();
	*z = vel.GetZ();
}

void JoltSetBodyLinearVelocity(JoltBodyInterface bodyInterface,
							   JoltBodyID bodyID,
							   float x, float y, float z)
{
	BodyInterface *bi = static_cast<BodyInterface *>(bodyInterface);
	const BodyID *bid = static_cast<const BodyID *>(bodyID);

	bi->SetLinearVelocity(*bid, Vec3(x, y, z));
}

void JoltAddBodyImpulse(JoltBodyInterface bodyInterface,
						JoltBodyID bodyID,
						float x, float y, float z)
{
	BodyInterface *bi = static_cast<BodyInterface *>(bodyInterface);
	const BodyID *bid = static_cast<const BodyID *>(bodyID);

	bi->AddImpulse(*bid, Vec3(x, y, z));
}

void JoltSetBodyDamping(JoltPhysicsSystem system, JoltBodyID bodyID, float linear, float angular)
{
	PhysicsSystemWrapper *wrapper = static_cast<PhysicsSystemWrapper *>(system);
	PhysicsSystem* ps = GetPhysicsSystem(wrapper);
	const BodyID *bid = static_cast<const BodyID *>(bodyID);

	BodyLockWrite lock(ps->GetBodyLockInterface(), *bid);
	if (lock.Succeeded())
	{
		Body &body = lock.GetBody();
		if (body.GetMotionProperties())
		{
			body.GetMotionProperties()->SetLinearDamping(linear);
			body.GetMotionProperties()->SetAngularDamping(angular);
		}
	}
}

void JoltSetBodyMotionQuality(JoltBodyInterface bodyInterface,
                             JoltBodyID bodyID,
                             JoltMotionQuality quality)
{
    BodyInterface *bi = static_cast<BodyInterface *>(bodyInterface);
    const BodyID *bid = static_cast<const BodyID *>(bodyID);
    bi->SetMotionQuality(*bid, static_cast<EMotionQuality>(quality));
}

void JoltSetBodyAllowedDOFs(JoltPhysicsSystem system,
                           JoltBodyID bodyID,
                           int allowedDOFs)
{
    PhysicsSystemWrapper *wrapper = static_cast<PhysicsSystemWrapper *>(system);
    PhysicsSystem* ps = GetPhysicsSystem(wrapper);
    const BodyID *bid = static_cast<const BodyID *>(bodyID);

    BodyLockWrite lock(ps->GetBodyLockInterface(), *bid);
    if (lock.Succeeded())
    {
        Body &body = lock.GetBody();
        body.SetAllowedDOFs(static_cast<EAllowedDOFs>(static_cast<uint8_t>(allowedDOFs)));
    }
}

void JoltSetBodyGravityFactor(JoltBodyInterface bodyInterface,
                             JoltBodyID bodyID,
                             float gravityFactor)
{
    BodyInterface *bi = static_cast<BodyInterface *>(bodyInterface);
    const BodyID *bid = static_cast<const BodyID *>(bodyID);
    bi->SetGravityFactor(*bid, gravityFactor);
}

void JoltSetBodyIsSensor(JoltBodyInterface bodyInterface,
                        JoltBodyID bodyID,
                        int isSensor)
{
    BodyInterface *bi = static_cast<BodyInterface *>(bodyInterface);
    const BodyID *bid = static_cast<const BodyID *>(bodyID);
    bi->SetIsSensor(*bid, isSensor != 0);
}

void JoltSetBodyFriction(JoltBodyInterface bodyInterface,
						 JoltBodyID bodyID,
						 float friction)
{
	BodyInterface *bi = static_cast<BodyInterface *>(bodyInterface);
	const BodyID *bid = static_cast<const BodyID *>(bodyID);
	bi->SetFriction(*bid, friction);
}

void JoltSetBodyRestitution(JoltBodyInterface bodyInterface,
							JoltBodyID bodyID,
							float restitution)
{
	BodyInterface *bi = static_cast<BodyInterface *>(bodyInterface);
	const BodyID *bid = static_cast<const BodyID *>(bodyID);
	bi->SetRestitution(*bid, restitution);
}

void JoltBatchGetBodyTransforms(JoltBodyInterface bodyInterface,
								const JoltBodyID *bodyIDs,
								int count,
								float *out)
{
	const BodyInterface *bi = static_cast<const BodyInterface *>(bodyInterface);
	float *w = out;
	for (int i = 0; i < count; i++)
	{
		const BodyID *bid = static_cast<const BodyID *>(bodyIDs[i]);
		RVec3 pos = bi->GetPosition(*bid);
		Quat q = bi->GetRotation(*bid);
		w[0] = static_cast<float>(pos.GetX());
		w[1] = static_cast<float>(pos.GetY());
		w[2] = static_cast<float>(pos.GetZ());
		w[3] = q.GetX();
		w[4] = q.GetY();
		w[5] = q.GetZ();
		w[6] = q.GetW();
		w[7] = bi->IsActive(*bid) ? 1.0f : 0.0f;
		w += 8;
	}
}

void JoltBatchApplyGravityDelta(JoltBodyInterface bodyInterface,
								const JoltBodyID *bodyIDs,
								int count,
								float dvx,
								float dvy,
								float dvz)
{
	BodyInterface *bi = static_cast<BodyInterface *>(bodyInterface);
	for (int i = 0; i < count; i++)
	{
		const BodyID *bid = static_cast<const BodyID *>(bodyIDs[i]);
		Vec3 v = bi->GetLinearVelocity(*bid);
		bi->SetLinearVelocity(*bid, Vec3(v.GetX() + dvx, v.GetY() + dvy, v.GetZ() + dvz));
		bi->ActivateBody(*bid);
	}
}

JoltBodyID JoltCreateBody(JoltBodyInterface bodyInterface,
						  JoltShape shape,
						  float x, float y, float z,
						  JoltMotionType motionType,
						  int isSensor,
						  float friction,
						  float restitution,
						  int allowedDOFsOrZero,
						  int objectLayerOrMinusOne)
{
	BodyInterface *bi = static_cast<BodyInterface *>(bodyInterface);
	const Shape *s = static_cast<const Shape *>(shape);

	EMotionType joltMotionType;
	ObjectLayer layer;

	if (objectLayerOrMinusOne >= 0 && objectLayerOrMinusOne < (int)Layers::NUM_LAYERS)
	{
		layer = static_cast<ObjectLayer>(objectLayerOrMinusOne);
		switch (motionType)
		{
		case JoltMotionTypeStatic:
			joltMotionType = EMotionType::Static;
			break;
		case JoltMotionTypeKinematic:
			joltMotionType = EMotionType::Kinematic;
			break;
		case JoltMotionTypeDynamic:
			joltMotionType = EMotionType::Dynamic;
			break;
		default:
			joltMotionType = EMotionType::Static;
			break;
		}
	}
	else
	{
		switch (motionType)
		{
		case JoltMotionTypeStatic:
			joltMotionType = EMotionType::Static;
			layer = Layers::NON_MOVING;
			break;
		case JoltMotionTypeKinematic:
			joltMotionType = EMotionType::Kinematic;
			layer = isSensor ? Layers::SENSOR : Layers::MOVING;
			break;
		case JoltMotionTypeDynamic:
			joltMotionType = EMotionType::Dynamic;
			layer = isSensor ? Layers::SENSOR : Layers::MOVING;
			break;
		default:
			joltMotionType = EMotionType::Static;
			layer = Layers::NON_MOVING;
			break;
		}
	}

	BodyCreationSettings body_settings(
		s,
		RVec3(x, y, z),
		Quat::sIdentity(),
		joltMotionType,
		layer);

	body_settings.mIsSensor = (isSensor != 0);
	body_settings.mFriction = friction;
	body_settings.mRestitution = restitution;

    // Set Tri-Tier defaults if appropriate
    if (joltMotionType == EMotionType::Static) {
        body_settings.mFriction = 0.5f; // Architect's Stage default
    } else if (joltMotionType == EMotionType::Dynamic) {
        body_settings.mLinearDamping = 0.05f; // Architect's Prop default
        body_settings.mAngularDamping = 0.05f;
    }

	if (allowedDOFsOrZero != 0)
	{
		body_settings.mAllowedDOFs = static_cast<EAllowedDOFs>(static_cast<uint8_t>(allowedDOFsOrZero));
	}

	Body *body = bi->CreateBody(body_settings);
	if (!body)
	{
		return nullptr;
	}

	bi->AddBody(body->GetID(), EActivation::DontActivate);

	auto bodyIDPtr = std::make_unique<BodyID>(body->GetID());
	return static_cast<JoltBodyID>(bodyIDPtr.release());
}

void JoltActivateBody(JoltBodyInterface bodyInterface, JoltBodyID bodyID)
{
	BodyInterface *bi = static_cast<BodyInterface *>(bodyInterface);
	const BodyID *bid = static_cast<const BodyID *>(bodyID);

	bi->ActivateBody(*bid);
}

void JoltDeactivateBody(JoltBodyInterface bodyInterface, JoltBodyID bodyID)
{
	BodyInterface *bi = static_cast<BodyInterface *>(bodyInterface);
	const BodyID *bid = static_cast<const BodyID *>(bodyID);

	bi->DeactivateBody(*bid);
}

void JoltSetBodyShape(JoltBodyInterface bodyInterface,
					 JoltBodyID bodyID,
                     JoltShape shape,
					 int updateMassProperties)
{
	BodyInterface *bi = static_cast<BodyInterface *>(bodyInterface);
	const BodyID *bid = static_cast<const BodyID *>(bodyID);
	const Shape *s = static_cast<const Shape *>(shape);

	bi->SetShape(*bid, s, updateMassProperties != 0, EActivation::Activate);
}

unsigned int JoltBodyIDGetIndexAndSequenceNumber(JoltBodyID bodyID)
{
	const BodyID *bid = static_cast<const BodyID *>(bodyID);
	return bid->GetIndexAndSequenceNumber();
}

void JoltDestroyBodyID(JoltBodyID bodyID)
{
	BodyID *bid = static_cast<BodyID *>(bodyID);
	delete bid;
}
