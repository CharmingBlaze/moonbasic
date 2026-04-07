/*
 * Jolt Physics C Wrapper - Character Virtual Implementation
 */

#include "character.h"
#include "physics.h"
#include "core.h"
#include <Jolt/Jolt.h>
#include <Jolt/Core/TempAllocator.h>
#include <Jolt/Physics/Collision/Shape/CapsuleShape.h>
#include <Jolt/Physics/Character/CharacterVirtual.h>
#include <memory>

using namespace JPH;

// Collision layers (defined in physics.cpp)
namespace Layers
{
	static constexpr ObjectLayer MOVING = 1;
};

// Adapter: converts ObjectVsBroadPhaseLayerFilter to BroadPhaseLayerFilter for character collision
class BroadPhaseLayerFilterAdapter : public BroadPhaseLayerFilter
{
public:
	BroadPhaseLayerFilterAdapter(const ObjectVsBroadPhaseLayerFilter* filter, ObjectLayer layer)
		: m_filter(filter), m_object_layer(layer) {}

	virtual bool ShouldCollide(BroadPhaseLayer inLayer) const override
	{
		return m_filter->ShouldCollide(m_object_layer, inLayer);
	}

private:
	const ObjectVsBroadPhaseLayerFilter* m_filter;
	ObjectLayer m_object_layer;
};

// Adapter: converts ObjectLayerPairFilter to ObjectLayerFilter for character collision
class ObjectLayerFilterAdapter : public ObjectLayerFilter
{
public:
	ObjectLayerFilterAdapter(const ObjectLayerPairFilter* filter, ObjectLayer layer)
		: m_filter(filter), m_object_layer(layer) {}

	virtual bool ShouldCollide(ObjectLayer inLayer) const override
	{
		return m_filter->ShouldCollide(m_object_layer, inLayer);
	}

private:
	const ObjectLayerPairFilter* m_filter;
	ObjectLayer m_object_layer;
};

JoltCharacterVirtual JoltCreateCharacterVirtual(JoltPhysicsSystem system,
											 const JoltCharacterVirtualSettings* goSettings,
											 float x, float y, float z)
{
	PhysicsSystemWrapper *wrapper = static_cast<PhysicsSystemWrapper *>(system);
	const Shape* s = static_cast<const Shape*>(goSettings->shape);

	CharacterVirtualSettings settings;
	settings.mShape = s;
	settings.mUp = Vec3(goSettings->upX, goSettings->upY, goSettings->upZ);
	settings.mMaxSlopeAngle = goSettings->maxSlopeAngle;
	settings.mMass = goSettings->mass;
	settings.mMaxStrength = goSettings->maxStrength;
	settings.mShapeOffset = Vec3(goSettings->shapeOffsetX, goSettings->shapeOffsetY, goSettings->shapeOffsetZ);
	settings.mBackFaceMode = static_cast<EBackFaceMode>(goSettings->backFaceMode);
	settings.mPredictiveContactDistance = goSettings->predictiveContactDistance;
	settings.mMaxCollisionIterations = goSettings->maxCollisionIterations;
	settings.mMaxConstraintIterations = goSettings->maxConstraintIterations;
	settings.mMinTimeRemaining = goSettings->minTimeRemaining;
	settings.mCollisionTolerance = goSettings->collisionTolerance;
	settings.mCharacterPadding = goSettings->characterPadding;
	settings.mMaxNumHits = goSettings->maxNumHits;
	settings.mHitReductionCosMaxAngle = goSettings->hitReductionCosMaxAngle;
	settings.mPenetrationRecoverySpeed = goSettings->penetrationRecoverySpeed;
	settings.mEnhancedInternalEdgeRemoval = goSettings->enhancedInternalEdgeRemoval != 0;

	// Create at specified position using smart pointer for exception safety
	auto character = std::make_unique<CharacterVirtual>(&settings, RVec3(x, y, z), Quat::sIdentity(), GetPhysicsSystem(wrapper));
	return static_cast<JoltCharacterVirtual>(character.release());
}

void JoltDestroyCharacterVirtual(JoltCharacterVirtual character)
{
	CharacterVirtual* cv = static_cast<CharacterVirtual*>(character);
	delete cv;
}

void JoltCharacterVirtualUpdate(JoltCharacterVirtual character,
								JoltPhysicsSystem system,
								float deltaTime,
								float gravityX, float gravityY, float gravityZ)
{
	CharacterVirtual* cv = static_cast<CharacterVirtual*>(character);
	PhysicsSystemWrapper* wrapper = static_cast<PhysicsSystemWrapper*>(system);

	// Use MOVING layer for character (same as dynamic bodies)
	BroadPhaseLayerFilterAdapter broad_phase_filter(GetObjectVsBroadPhaseLayerFilter(wrapper), Layers::MOVING);
	ObjectLayerFilterAdapter object_layer_filter(GetObjectLayerPairFilter(wrapper), Layers::MOVING);

	// Call basic Update with gravity vector and layer filters
	cv->Update(
		deltaTime,
		Vec3(gravityX, gravityY, gravityZ),
		broad_phase_filter,
		object_layer_filter,
		{}, // Empty BodyFilter (collides with all bodies)
		{}, // Empty ShapeFilter (collides with all shapes)
		*gTempAllocator.get()
	);
}

void JoltCharacterVirtualExtendedUpdate(JoltCharacterVirtual character,
										JoltPhysicsSystem system,
										float deltaTime,
										float gravityX, float gravityY, float gravityZ)
{
	CharacterVirtual* cv = static_cast<CharacterVirtual*>(character);
	PhysicsSystemWrapper* wrapper = static_cast<PhysicsSystemWrapper*>(system);

	// Use default extended update settings
	CharacterVirtual::ExtendedUpdateSettings settings;

	// Use MOVING layer for character (same as dynamic bodies)
	BroadPhaseLayerFilterAdapter broad_phase_filter(GetObjectVsBroadPhaseLayerFilter(wrapper), Layers::MOVING);
	ObjectLayerFilterAdapter object_layer_filter(GetObjectLayerPairFilter(wrapper), Layers::MOVING);

	// Call ExtendedUpdate with gravity vector and layer filters
	cv->ExtendedUpdate(
		deltaTime,
		Vec3(gravityX, gravityY, gravityZ),
		settings,
		broad_phase_filter,
		object_layer_filter,
		{}, // Empty BodyFilter (collides with all bodies)
		{}, // Empty ShapeFilter (collides with all shapes)
		*gTempAllocator.get()
	);
}

void JoltCharacterVirtualSetLinearVelocity(JoltCharacterVirtual character,
										   float x, float y, float z)
{
	CharacterVirtual* cv = static_cast<CharacterVirtual*>(character);
	cv->SetLinearVelocity(Vec3(x, y, z));
}

void JoltCharacterVirtualGetLinearVelocity(const JoltCharacterVirtual character,
										   float* x, float* y, float* z)
{
	const CharacterVirtual* cv = static_cast<const CharacterVirtual*>(character);
	Vec3 vel = cv->GetLinearVelocity();
	*x = vel.GetX();
	*y = vel.GetY();
	*z = vel.GetZ();
}

void JoltCharacterVirtualGetGroundVelocity(const JoltCharacterVirtual character,
											float* x, float* y, float* z)
{
	const CharacterVirtual* cv = static_cast<const CharacterVirtual*>(character);
	Vec3 vel = cv->GetGroundVelocity();
	*x = vel.GetX();
	*y = vel.GetY();
	*z = vel.GetZ();
}

void JoltCharacterVirtualSetPosition(JoltCharacterVirtual character,
									 float x, float y, float z)
{
	CharacterVirtual* cv = static_cast<CharacterVirtual*>(character);
	cv->SetPosition(RVec3(x, y, z));
}

void JoltCharacterVirtualGetPosition(const JoltCharacterVirtual character,
									 float* x, float* y, float* z)
{
	const CharacterVirtual* cv = static_cast<const CharacterVirtual*>(character);
	RVec3 pos = cv->GetPosition();
	*x = static_cast<float>(pos.GetX());
	*y = static_cast<float>(pos.GetY());
	*z = static_cast<float>(pos.GetZ());
}

JoltGroundState JoltCharacterVirtualGetGroundState(const JoltCharacterVirtual character)
{
	const CharacterVirtual* cv = static_cast<const CharacterVirtual*>(character);
	CharacterBase::EGroundState state = cv->GetGroundState();
	return static_cast<JoltGroundState>(state);
}

int JoltCharacterVirtualIsSupported(const JoltCharacterVirtual character)
{
	const CharacterVirtual* cv = static_cast<const CharacterVirtual*>(character);
	return cv->IsSupported() ? 1 : 0;
}

void JoltCharacterVirtualSetShape(JoltCharacterVirtual character,
								  JoltShape shape,
								  float maxPenetrationDepth,
								  JoltPhysicsSystem system)
{
	CharacterVirtual* cv = static_cast<CharacterVirtual*>(character);
	const Shape* s = static_cast<const Shape*>(shape);
	PhysicsSystemWrapper* wrapper = static_cast<PhysicsSystemWrapper*>(system);

	// Use MOVING layer for character (same as dynamic bodies)
	BroadPhaseLayerFilterAdapter broad_phase_filter(GetObjectVsBroadPhaseLayerFilter(wrapper), Layers::MOVING);
	ObjectLayerFilterAdapter object_layer_filter(GetObjectLayerPairFilter(wrapper), Layers::MOVING);

	// Call SetShape with required filters
	cv->SetShape(
		s,
		maxPenetrationDepth,
		broad_phase_filter,
		object_layer_filter,
		{}, // Empty BodyFilter (collides with all bodies)
		{}, // Empty ShapeFilter (collides with all shapes)
		*gTempAllocator.get()
	);
}

// Get the shape of a virtual character
JoltShape JoltCharacterVirtualGetShape(const JoltCharacterVirtual character)
{
	const CharacterVirtual* cv = static_cast<const CharacterVirtual*>(character);
	return const_cast<Shape*>(cv->GetShape());
}

// Get the normal of the ground surface the character is standing on
void JoltCharacterVirtualGetGroundNormal(const JoltCharacterVirtual character,
										 float* x, float* y, float* z)
{
	const CharacterVirtual* cv = static_cast<const CharacterVirtual*>(character);
	Vec3 normal = cv->GetGroundNormal();
	*x = normal.GetX();
	*y = normal.GetY();
	*z = normal.GetZ();
}

// Get the position of the ground contact point
void JoltCharacterVirtualGetGroundPosition(const JoltCharacterVirtual character,
										   float* x, float* y, float* z)
{
	const CharacterVirtual* cv = static_cast<const CharacterVirtual*>(character);
	RVec3 pos = cv->GetGroundPosition();
	*x = static_cast<float>(pos.GetX());
	*y = static_cast<float>(pos.GetY());
	*z = static_cast<float>(pos.GetZ());
}

// Get the active contacts for the character
int JoltCharacterVirtualGetActiveContacts(const JoltCharacterVirtual character,
										  JoltCharacterContact* contacts,
										  int maxContacts)
{
	const CharacterVirtual* cv = static_cast<const CharacterVirtual*>(character);
	const CharacterVirtual::ContactList& activeContacts = cv->GetActiveContacts();

	int numContacts = static_cast<int>(activeContacts.size());
	int numToReturn = numContacts < maxContacts ? numContacts : maxContacts;

	for (int i = 0; i < numToReturn; i++)
	{
		const CharacterVirtual::Contact& c = activeContacts[i];

		// Copy position
		contacts[i].positionX = static_cast<float>(c.mPosition.GetX());
		contacts[i].positionY = static_cast<float>(c.mPosition.GetY());
		contacts[i].positionZ = static_cast<float>(c.mPosition.GetZ());

		// Copy linear velocity
		contacts[i].linearVelocityX = c.mLinearVelocity.GetX();
		contacts[i].linearVelocityY = c.mLinearVelocity.GetY();
		contacts[i].linearVelocityZ = c.mLinearVelocity.GetZ();

		// Copy contact normal
		contacts[i].contactNormalX = c.mContactNormal.GetX();
		contacts[i].contactNormalY = c.mContactNormal.GetY();
		contacts[i].contactNormalZ = c.mContactNormal.GetZ();

		// Copy surface normal
		contacts[i].surfaceNormalX = c.mSurfaceNormal.GetX();
		contacts[i].surfaceNormalY = c.mSurfaceNormal.GetY();
		contacts[i].surfaceNormalZ = c.mSurfaceNormal.GetZ();

		// Copy scalar fields
		contacts[i].distance = c.mDistance;
		contacts[i].fraction = c.mFraction;

		// Create a copy of the BodyID for the Go layer
		if (c.mBodyB.IsInvalid())
		{
			contacts[i].bodyB = nullptr;
		}
		else
		{
			contacts[i].bodyB = new BodyID(c.mBodyB);
		}

		contacts[i].userData = c.mUserData;

		// Copy bool fields (as int)
		contacts[i].isSensorB = c.mIsSensorB ? 1 : 0;
		contacts[i].hadCollision = c.mHadCollision ? 1 : 0;
		contacts[i].wasDiscarded = c.mWasDiscarded ? 1 : 0;
		contacts[i].canPushCharacter = c.mCanPushCharacter ? 1 : 0;
	}

	return numToReturn;
}
