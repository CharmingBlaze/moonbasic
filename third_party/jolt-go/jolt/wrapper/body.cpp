/*
 * Jolt Physics C Wrapper - Body Operations Implementation
 */

#include "body.h"
#include "physics.h"
#include <Jolt/Jolt.h>
#include <Jolt/Physics/PhysicsSystem.h>
#include <Jolt/Physics/Body/Body.h>
#include <Jolt/Physics/Body/BodyCreationSettings.h>
#include <Jolt/Physics/Body/BodyInterface.h>
#include <memory>

using namespace JPH;

// Collision layers (defined in physics.cpp)
namespace Layers
{
	static constexpr ObjectLayer NON_MOVING = 0;
	static constexpr ObjectLayer MOVING = 1;
};

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

void JoltSetBodyPosition(JoltBodyInterface bodyInterface,
						 JoltBodyID bodyID,
						 float x, float y, float z)
{
	BodyInterface *bi = static_cast<BodyInterface *>(bodyInterface);
	const BodyID *bid = static_cast<const BodyID *>(bodyID);

	bi->SetPosition(*bid, RVec3(x, y, z), EActivation::DontActivate);
}

JoltBodyID JoltCreateBody(JoltBodyInterface bodyInterface,
						  JoltShape shape,
						  float x, float y, float z,
						  JoltMotionType motionType,
						  int isSensor)
{
	BodyInterface *bi = static_cast<BodyInterface *>(bodyInterface);
	const Shape *s = static_cast<const Shape *>(shape);

	// Convert motion type
	EMotionType joltMotionType;
	ObjectLayer layer;
	switch (motionType)
	{
	case JoltMotionTypeStatic:
		joltMotionType = EMotionType::Static;
		layer = Layers::NON_MOVING;
		break;
	case JoltMotionTypeKinematic:
		joltMotionType = EMotionType::Kinematic;
		layer = Layers::MOVING;
		break;
	case JoltMotionTypeDynamic:
		joltMotionType = EMotionType::Dynamic;
		layer = Layers::MOVING;
		break;
	default:
		joltMotionType = EMotionType::Static;
		layer = Layers::NON_MOVING;
		break;
	}

	BodyCreationSettings body_settings(
		s,
		RVec3(x, y, z),
		Quat::sIdentity(),
		joltMotionType,
		layer);

	// Set sensor flag if requested
	body_settings.mIsSensor = (isSensor != 0);

	Body *body = bi->CreateBody(body_settings);
	if (!body)
	{
		return nullptr;
	}

	// Don't activate yet - caller will activate when ready
	bi->AddBody(body->GetID(), EActivation::DontActivate);

	// Use smart pointer for exception safety, then release to caller
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

void JoltDestroyBodyID(JoltBodyID bodyID)
{
	BodyID *bid = static_cast<BodyID *>(bodyID);
	delete bid;
}
