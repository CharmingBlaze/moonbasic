/*
 * Jolt Physics C Wrapper - Physics System Implementation
 */

#include "physics.h"
#include "core.h"
#include "physics_layers.h"
#include <Jolt/Jolt.h>
#include <Jolt/Core/TempAllocator.h>
#include <Jolt/Core/JobSystemThreadPool.h>
#include <Jolt/Physics/PhysicsSettings.h>
#include <Jolt/Physics/PhysicsSystem.h>
#include <Jolt/Physics/Collision/ContactListener.h>
#include <Jolt/Physics/Collision/RayCast.h>
#include <Jolt/Physics/Collision/CastResult.h>
#include <vector>
#include <mutex>
#include <memory>

using namespace JPH;

namespace Layers = MBPHYS_Layers;
namespace BroadPhaseLayers = MBPHYS_BroadPhaseLayers;

// Maps object layers to broad phase layers
class BPLayerInterfaceImpl final : public BroadPhaseLayerInterface
{
public:
	BPLayerInterfaceImpl()
	{
		mObjectToBroadPhase[Layers::NON_MOVING] = BroadPhaseLayers::NON_MOVING;
		mObjectToBroadPhase[Layers::MOVING] = BroadPhaseLayers::MOVING;
		mObjectToBroadPhase[Layers::CHARACTER] = BroadPhaseLayers::MOVING;
		mObjectToBroadPhase[Layers::SENSOR] = BroadPhaseLayers::MOVING;
		mObjectToBroadPhase[Layers::ONE_WAY] = BroadPhaseLayers::MOVING;
	}

	virtual uint GetNumBroadPhaseLayers() const override
	{
		return BroadPhaseLayers::NUM_LAYERS;
	}

	virtual BroadPhaseLayer GetBroadPhaseLayer(ObjectLayer inLayer) const override
	{
		JPH_ASSERT(inLayer < Layers::NUM_LAYERS);
		return mObjectToBroadPhase[inLayer];
	}

#if defined(JPH_EXTERNAL_PROFILE) || defined(JPH_PROFILE_ENABLED)
	virtual const char* GetBroadPhaseLayerName(BroadPhaseLayer inLayer) const override
	{
		switch ((BroadPhaseLayer::Type)inLayer)
		{
		case (BroadPhaseLayer::Type)BroadPhaseLayers::NON_MOVING:	return "NON_MOVING";
		case (BroadPhaseLayer::Type)BroadPhaseLayers::MOVING:		return "MOVING";
		default:													return "INVALID";
		}
	}
#endif // JPH_EXTERNAL_PROFILE || JPH_PROFILE_ENABLED

private:
	BroadPhaseLayer mObjectToBroadPhase[Layers::NUM_LAYERS];
};

// Filters which broad phase layers can collide
class ObjectVsBroadPhaseLayerFilterImpl : public ObjectVsBroadPhaseLayerFilter
{
public:
	virtual bool ShouldCollide(ObjectLayer inLayer1, BroadPhaseLayer inLayer2) const override
	{
		switch (inLayer1)
		{
		case Layers::NON_MOVING:
			return inLayer2 == BroadPhaseLayers::MOVING;
		case Layers::MOVING:
		case Layers::CHARACTER:
		case Layers::SENSOR:
		case Layers::ONE_WAY:
			return true;
		default:
			JPH_ASSERT(false);
			return false;
		}
	}
};

// Filters which object layers can collide with each other (symmetric)
class ObjectLayerPairFilterImpl : public ObjectLayerPairFilter
{
public:
	virtual bool ShouldCollide(ObjectLayer inObject1, ObjectLayer inObject2) const override
	{
		if (inObject1 == Layers::NON_MOVING && inObject2 == Layers::NON_MOVING)
			return false;
		if (inObject1 == Layers::SENSOR && inObject2 == Layers::SENSOR)
			return false;
		if (inObject1 == Layers::ONE_WAY && inObject2 == Layers::ONE_WAY)
			return false;
		return true;
	}
};

// Internal collision event queue for rigid bodies (sensors/triggers)
struct InternalCollisionEvent {
    BodyID body1;
    BodyID body2;
};

// Global contact listener implementation
class ContactListenerImpl : public ContactListener
{
public:
    std::mutex mQueueMutex;
    std::vector<InternalCollisionEvent> mEventQueue;

    virtual ValidateResult OnContactValidate(const Body &inBody1, const Body &inBody2, RVec3Arg inBaseOffset, const CollideShapeResult &inCollisionResult) override
    {
        return ValidateResult::AcceptAllContactsForThisBodyPair;
    }

    virtual void OnContactAdded(const Body &inBody1, const Body &inBody2, const ContactManifold &inManifold, ContactSettings &ioSettings) override
    {
        std::lock_guard<std::mutex> lock(mQueueMutex);
        mEventQueue.push_back({ inBody1.GetID(), inBody2.GetID() });
    }

    virtual void OnContactPersisted(const Body &inBody1, const Body &inBody2, const ContactManifold &inManifold, ContactSettings &ioSettings) override
    {
        // For sensors, we only care about the entry (Added), but if we want 'Stay' logic, we could add it here.
    }

    virtual void OnContactRemoved(const SubShapeIDPair &inSubShapePair) override
    {
        // Entry/Exit logic could be expanded here.
    }
};

// Wrapper to keep layer interfaces alive (PhysicsSystem stores references to them)
struct PhysicsSystemWrapper
{
	std::unique_ptr<PhysicsSystem> system;
	std::unique_ptr<BPLayerInterfaceImpl> broad_phase_layer_interface;
	std::unique_ptr<ObjectVsBroadPhaseLayerFilterImpl> object_vs_broadphase_layer_filter;
	std::unique_ptr<ObjectLayerPairFilterImpl> object_vs_object_layer_filter;
    std::unique_ptr<ContactListenerImpl> contact_listener;

	~PhysicsSystemWrapper() = default;
};

JoltPhysicsSystem JoltCreatePhysicsSystem()
{
	const uint cMaxBodies = 10240;
	const uint cNumBodyMutexes = 0;
	const uint cMaxBodyPairs = 65536;
	const uint cMaxContactConstraints = 20480;

	auto wrapper = std::make_unique<PhysicsSystemWrapper>();

	wrapper->broad_phase_layer_interface = std::make_unique<BPLayerInterfaceImpl>();
	wrapper->object_vs_broadphase_layer_filter = std::make_unique<ObjectVsBroadPhaseLayerFilterImpl>();
	wrapper->object_vs_object_layer_filter = std::make_unique<ObjectLayerPairFilterImpl>();
    wrapper->contact_listener = std::make_unique<ContactListenerImpl>();

	wrapper->system = std::make_unique<PhysicsSystem>();
	wrapper->system->Init(cMaxBodies, cNumBodyMutexes, cMaxBodyPairs, cMaxContactConstraints,
						  *wrapper->broad_phase_layer_interface,
						  *wrapper->object_vs_broadphase_layer_filter,
						  *wrapper->object_vs_object_layer_filter);

    wrapper->system->SetContactListener(wrapper->contact_listener.get());

	return static_cast<JoltPhysicsSystem>(wrapper.release());
}

void JoltDestroyPhysicsSystem(JoltPhysicsSystem system)
{
	PhysicsSystemWrapper *wrapper = static_cast<PhysicsSystemWrapper *>(system);
	delete wrapper;
}

void JoltPhysicsSystemUpdate(JoltPhysicsSystem system, float deltaTime)
{
	PhysicsSystemWrapper *wrapper = static_cast<PhysicsSystemWrapper *>(system);
	wrapper->system->Update(deltaTime, 1, gTempAllocator.get(), gJobSystem.get());
}

int JoltPhysicsSystemDrainContactQueue(JoltPhysicsSystem system, JoltCollisionEvent* out, int maxCount)
{
    PhysicsSystemWrapper *wrapper = static_cast<PhysicsSystemWrapper *>(system);
    auto& listener = wrapper->contact_listener;
    
    std::lock_guard<std::mutex> lock(listener->mQueueMutex);
    int count = static_cast<int>(listener->mEventQueue.size());
    if (count > maxCount) count = maxCount;

    for (int i = 0; i < count; i++) {
        // We need to return heap-allocated BodyID pointers that Go can manage, 
        // to match the pattern in JoltCreateBody.
        // Actually, for events, it's easier to just pass the IDs and have Go look them up.
        // But the C struct JoltCollisionEvent expects void*.
        
        auto b1 = std::make_unique<BodyID>(listener->mEventQueue[i].body1);
        auto b2 = std::make_unique<BodyID>(listener->mEventQueue[i].body2);
        
        out[i].body1 = static_cast<void*>(b1.release());
        out[i].body2 = static_cast<void*>(b2.release());
    }

    listener->mEventQueue.erase(listener->mEventQueue.begin(), listener->mEventQueue.begin() + count);
    return (int)count;
}

PhysicsSystem* GetPhysicsSystem(PhysicsSystemWrapper* wrapper)
{
	return wrapper->system.get();
}

const ObjectVsBroadPhaseLayerFilter* GetObjectVsBroadPhaseLayerFilter(PhysicsSystemWrapper* wrapper)
{
	return wrapper->object_vs_broadphase_layer_filter.get();
}

const ObjectLayerPairFilter* GetObjectLayerPairFilter(PhysicsSystemWrapper* wrapper)
{
	return wrapper->object_vs_object_layer_filter.get();
}
