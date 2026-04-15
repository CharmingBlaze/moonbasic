/*
 * Jolt Physics C Wrapper - Character Virtual Implementation
 */

#include "character.h"
#include "physics.h"
#include "physics_bridge.h"
#include "physics_layers.h"
#include "core.h"
#include <Jolt/Jolt.h>
#include <Jolt/Physics/PhysicsSystem.h>
#include <Jolt/Physics/Body/BodyLock.h>
#include <Jolt/Core/TempAllocator.h>
#include <Jolt/Physics/Collision/Shape/CapsuleShape.h>
#include <Jolt/Physics/Character/CharacterVirtual.h>
#include <memory>
#include <mutex>
#include <vector>
#include <atomic>

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

// Character contact listener that pushes events into a queue
class CharacterContactListenerImpl : public CharacterContactListener
{
public:
    explicit CharacterContactListenerImpl(PhysicsSystem *inPhysicsSystem)
        : mPhysicsSystem(inPhysicsSystem), m_enabled(false) {}

    void SetEnabled(bool enabled) { m_enabled = enabled; }
    bool IsEnabled() const { return m_enabled; }

    virtual void OnContactAdded(const CharacterVirtual *inCharacter, const BodyID &inBodyID2, const SubShapeID &inSubShapeID2, RVec3Arg inContactPosition, Vec3Arg inContactNormal, CharacterContactSettings &ioSettings) override
    {
        // One-way platforms: when hitting ONE_WAY layer from below, disable contact response (pass-through).
        if (mPhysicsSystem != nullptr)
        {
            BodyLockRead lock(mPhysicsSystem->GetBodyLockInterface(), inBodyID2);
            if (lock.Succeeded())
            {
                const Body &body = lock.GetBody();
                if (body.GetObjectLayer() == MBPHYS_Layers::ONE_WAY)
                {
                    Vec3 up = inCharacter->GetUp();
                    if (inContactNormal.Dot(up) > 0.05f)
                    {
                        ioSettings.mCanPushCharacter = false;
                        ioSettings.mCanReceiveImpulses = false;
                    }
                }
            }
        }

        if (!m_enabled) return;

        std::lock_guard<std::mutex> lock(m_mutex);
        if (m_events.size() < 1024) { // Avoid unbounded growth
            JoltCharacterContactEvent event;
            event.bodyB = inBodyID2.GetIndexAndSequenceNumber();
            event.positionX = (float)inContactPosition.GetX();
            event.positionY = (float)inContactPosition.GetY();
            event.positionZ = (float)inContactPosition.GetZ();
            event.normalX = inContactNormal.GetX();
            event.normalY = inContactNormal.GetY();
            event.normalZ = inContactNormal.GetZ();
            event.distance = 0.0f; // OnContactAdded is actual contact
            m_events.push_back(event);
        }
    }

    int Drain(JoltCharacterContactEvent* outEvents, int maxEvents)
    {
        std::lock_guard<std::mutex> lock(m_mutex);
        int count = (int)m_events.size();
        if (count == 0) return 0;
        if (count > maxEvents) count = maxEvents;
        
        std::copy(m_events.begin(), m_events.begin() + count, outEvents);
        
        m_events.erase(m_events.begin(), m_events.begin() + count);
        return count;
    }

private:
    PhysicsSystem *mPhysicsSystem;
    std::mutex m_mutex;
    std::vector<JoltCharacterContactEvent> m_events;
    std::atomic<bool> m_enabled;
};

// Internal wrapper structure to keep track of character and its listener
struct CharacterVirtualWrapper
{
    CharacterVirtual* mCharacter;
    CharacterContactListenerImpl mListener;

    CharacterVirtualWrapper(const CharacterVirtualSettings* inSettings, RVec3Arg inPosition, QuatArg inRotation, PhysicsSystem* inSystem)
        : mListener(inSystem)
    {
        mCharacter = new CharacterVirtual(inSettings, inPosition, inRotation, inSystem);
        mCharacter->SetListener(&mListener);
    }

    ~CharacterVirtualWrapper()
    {
        delete mCharacter;
    }
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

	// Create wrapper which manages character and listener
	auto wrapper_obj = std::make_unique<CharacterVirtualWrapper>(&settings, RVec3(x, y, z), Quat::sIdentity(), GetPhysicsSystem(wrapper));
	return static_cast<JoltCharacterVirtual>(wrapper_obj.release());
}

void JoltDestroyCharacterVirtual(JoltCharacterVirtual character)
{
	CharacterVirtualWrapper* wrapper = static_cast<CharacterVirtualWrapper*>(character);
	delete wrapper;
}

void JoltCharacterVirtualUpdate(JoltCharacterVirtual character,
								JoltPhysicsSystem system,
								float deltaTime,
								float gravityX, float gravityY, float gravityZ)
{
	CharacterVirtualWrapper* wrapper_obj = static_cast<CharacterVirtualWrapper*>(character);
	CharacterVirtual* cv = wrapper_obj->mCharacter;
	PhysicsSystemWrapper* physics_wrapper = static_cast<PhysicsSystemWrapper*>(system);

	// Use MOVING layer for character (same as dynamic bodies)
	BroadPhaseLayerFilterAdapter broad_phase_filter(GetObjectVsBroadPhaseLayerFilter(physics_wrapper), Layers::MOVING);
	ObjectLayerFilterAdapter object_layer_filter(GetObjectLayerPairFilter(physics_wrapper), Layers::MOVING);

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
										float gravityX, float gravityY, float gravityZ,
										const JoltCharacterExtendedUpdateSettings* extendedSettingsOrNull)
{
	CharacterVirtualWrapper* wrapper_obj = static_cast<CharacterVirtualWrapper*>(character);
	CharacterVirtual* cv = wrapper_obj->mCharacter;
	PhysicsSystemWrapper* physics_wrapper = static_cast<PhysicsSystemWrapper*>(system);

	CharacterVirtual::ExtendedUpdateSettings settings;
	if (extendedSettingsOrNull != nullptr)
	{
		settings.mStickToFloorStepDown =
			Vec3(extendedSettingsOrNull->stickToFloorStepDownX,
				 extendedSettingsOrNull->stickToFloorStepDownY,
				 extendedSettingsOrNull->stickToFloorStepDownZ);
		settings.mWalkStairsStepUp =
			Vec3(extendedSettingsOrNull->walkStairsStepUpX,
				 extendedSettingsOrNull->walkStairsStepUpY,
				 extendedSettingsOrNull->walkStairsStepUpZ);
		settings.mWalkStairsMinStepForward = extendedSettingsOrNull->walkStairsMinStepForward;
		settings.mWalkStairsStepForwardTest = extendedSettingsOrNull->walkStairsStepForwardTest;
		settings.mWalkStairsCosAngleForwardContact = extendedSettingsOrNull->walkStairsCosAngleForwardContact;
		settings.mWalkStairsStepDownExtra =
			Vec3(extendedSettingsOrNull->walkStairsStepDownExtraX,
				 extendedSettingsOrNull->walkStairsStepDownExtraY,
				 extendedSettingsOrNull->walkStairsStepDownExtraZ);
	}

	// Use MOVING layer for character (same as dynamic bodies)
	BroadPhaseLayerFilterAdapter broad_phase_filter(GetObjectVsBroadPhaseLayerFilter(physics_wrapper), Layers::MOVING);
	ObjectLayerFilterAdapter object_layer_filter(GetObjectLayerPairFilter(physics_wrapper), Layers::MOVING);

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
	CharacterVirtualWrapper* wrapper_obj = static_cast<CharacterVirtualWrapper*>(character);
	CharacterVirtual* cv = wrapper_obj->mCharacter;
	cv->SetLinearVelocity(Vec3(x, y, z));
}

void JoltCharacterVirtualGetLinearVelocity(const JoltCharacterVirtual character,
										   float* x, float* y, float* z)
{
	CharacterVirtualWrapper* wrapper_obj = const_cast<CharacterVirtualWrapper*>(static_cast<const CharacterVirtualWrapper*>(character));
	CharacterVirtual* cv = wrapper_obj->mCharacter;
	Vec3 vel = cv->GetLinearVelocity();
	*x = vel.GetX();
	*y = vel.GetY();
	*z = vel.GetZ();
}

void JoltCharacterVirtualGetGroundVelocity(const JoltCharacterVirtual character,
											float* x, float* y, float* z)
{
	CharacterVirtualWrapper* wrapper_obj = const_cast<CharacterVirtualWrapper*>(static_cast<const CharacterVirtualWrapper*>(character));
	CharacterVirtual* cv = wrapper_obj->mCharacter;
	Vec3 vel = cv->GetGroundVelocity();
	*x = vel.GetX();
	*y = vel.GetY();
	*z = vel.GetZ();
}

void JoltCharacterVirtualSetPosition(JoltCharacterVirtual character,
									 float x, float y, float z)
{
	CharacterVirtualWrapper* wrapper_obj = static_cast<CharacterVirtualWrapper*>(character);
	CharacterVirtual* cv = wrapper_obj->mCharacter;
	cv->SetPosition(RVec3(x, y, z));
}

void JoltCharacterVirtualGetPosition(const JoltCharacterVirtual character,
									 float* x, float* y, float* z)
{
	CharacterVirtualWrapper* wrapper_obj = const_cast<CharacterVirtualWrapper*>(static_cast<const CharacterVirtualWrapper*>(character));
	CharacterVirtual* cv = wrapper_obj->mCharacter;
	RVec3 pos = cv->GetPosition();
	*x = static_cast<float>(pos.GetX());
	*y = static_cast<float>(pos.GetY());
	*z = static_cast<float>(pos.GetZ());
}

JoltGroundState JoltCharacterVirtualGetGroundState(const JoltCharacterVirtual character)
{
	CharacterVirtualWrapper* wrapper_obj = const_cast<CharacterVirtualWrapper*>(static_cast<const CharacterVirtualWrapper*>(character));
	CharacterVirtual* cv = wrapper_obj->mCharacter;
	CharacterBase::EGroundState state = cv->GetGroundState();
	return static_cast<JoltGroundState>(state);
}

int JoltCharacterVirtualIsSupported(const JoltCharacterVirtual character)
{
	CharacterVirtualWrapper* wrapper_obj = const_cast<CharacterVirtualWrapper*>(static_cast<const CharacterVirtualWrapper*>(character));
	CharacterVirtual* cv = wrapper_obj->mCharacter;
	return cv->IsSupported() ? 1 : 0;
}

void JoltCharacterVirtualSetShape(JoltCharacterVirtual character,
								  JoltShape shape,
								  float maxPenetrationDepth,
								  JoltPhysicsSystem system)
{
	CharacterVirtualWrapper* wrapper_obj = static_cast<CharacterVirtualWrapper*>(character);
	CharacterVirtual* cv = wrapper_obj->mCharacter;
	const Shape* s = static_cast<const Shape*>(shape);
	PhysicsSystemWrapper* physics_wrapper = static_cast<PhysicsSystemWrapper*>(system);

	// Use MOVING layer for character (same as dynamic bodies)
	BroadPhaseLayerFilterAdapter broad_phase_filter(GetObjectVsBroadPhaseLayerFilter(physics_wrapper), Layers::MOVING);
	ObjectLayerFilterAdapter object_layer_filter(GetObjectLayerPairFilter(physics_wrapper), Layers::MOVING);

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
	CharacterVirtualWrapper* wrapper_obj = const_cast<CharacterVirtualWrapper*>(static_cast<const CharacterVirtualWrapper*>(character));
	CharacterVirtual* cv = wrapper_obj->mCharacter;
	return const_cast<Shape*>(cv->GetShape());
}

// Get the normal of the ground surface the character is standing on
void JoltCharacterVirtualGetGroundNormal(const JoltCharacterVirtual character,
										 float* x, float* y, float* z)
{
	CharacterVirtualWrapper* wrapper_obj = const_cast<CharacterVirtualWrapper*>(static_cast<const CharacterVirtualWrapper*>(character));
	CharacterVirtual* cv = wrapper_obj->mCharacter;
	Vec3 normal = cv->GetGroundNormal();
	*x = normal.GetX();
	*y = normal.GetY();
	*z = normal.GetZ();
}

// Get the position of the ground contact point
void JoltCharacterVirtualGetGroundPosition(const JoltCharacterVirtual character,
										   float* x, float* y, float* z)
{
	CharacterVirtualWrapper* wrapper_obj = const_cast<CharacterVirtualWrapper*>(static_cast<const CharacterVirtualWrapper*>(character));
	CharacterVirtual* cv = wrapper_obj->mCharacter;
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
	CharacterVirtualWrapper* wrapper_obj = const_cast<CharacterVirtualWrapper*>(static_cast<const CharacterVirtualWrapper*>(character));
	CharacterVirtual* cv = wrapper_obj->mCharacter;
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

void JoltCharacterVirtualSetContactListenerEnabled(JoltCharacterVirtual character, int enabled)
{
    CharacterVirtualWrapper* wrapper_obj = static_cast<CharacterVirtualWrapper*>(character);
    wrapper_obj->mListener.SetEnabled(enabled != 0);
}

int JoltCharacterVirtualDrainContactQueue(JoltCharacterVirtual character,
                                         JoltCharacterContactEvent* events,
                                         int maxEvents)
{
    CharacterVirtualWrapper* wrapper_obj = static_cast<CharacterVirtualWrapper*>(character);
    return wrapper_obj->mListener.Drain(events, maxEvents);
}
