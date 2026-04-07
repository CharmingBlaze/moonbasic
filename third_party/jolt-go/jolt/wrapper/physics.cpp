/*
 * Jolt Physics C Wrapper - Physics System Implementation
 *
 * CRITICAL BUG FIX: PhysicsSystem::Init() takes references to layer interfaces,
 * so they must be heap-allocated and kept alive. We use PhysicsSystemWrapper
 * to manage their lifetime alongside the PhysicsSystem.
 */

#include "physics.h"
#include "core.h"
#include <Jolt/Jolt.h>
#include <Jolt/Core/TempAllocator.h>
#include <Jolt/Core/JobSystemThreadPool.h>
#include <Jolt/Physics/PhysicsSettings.h>
#include <Jolt/Physics/PhysicsSystem.h>
#include <memory>

using namespace JPH;

// Collision layers: NON_MOVING (static) and MOVING (dynamic)
namespace Layers
{
	static constexpr ObjectLayer NON_MOVING = 0;
	static constexpr ObjectLayer MOVING = 1;

	static constexpr ObjectLayer NUM_LAYERS = 2;
};

namespace BroadPhaseLayers
{
	static constexpr BroadPhaseLayer NON_MOVING(0);
	static constexpr BroadPhaseLayer MOVING(1);

	static constexpr uint NUM_LAYERS(2);
};

// Maps object layers to broad phase layers
class BPLayerInterfaceImpl final : public BroadPhaseLayerInterface
{
public:
	BPLayerInterfaceImpl()
	{
		mObjectToBroadPhase[Layers::NON_MOVING] = BroadPhaseLayers::NON_MOVING;
		mObjectToBroadPhase[Layers::MOVING] = BroadPhaseLayers::MOVING;
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
			return true;
		default:
			JPH_ASSERT(false);
			return false;
		}
	}
};

// Filters which object layers can collide with each other
class ObjectLayerPairFilterImpl : public ObjectLayerPairFilter
{
public:
	virtual bool ShouldCollide(ObjectLayer inObject1, ObjectLayer inObject2) const override
	{
		switch (inObject1)
		{
		case Layers::NON_MOVING:
			return inObject2 == Layers::MOVING;
		case Layers::MOVING:
			return true;
		default:
			JPH_ASSERT(false);
			return false;
		}
	}
};

// Wrapper to keep layer interfaces alive (PhysicsSystem stores references to them)
struct PhysicsSystemWrapper
{
	std::unique_ptr<PhysicsSystem> system;
	std::unique_ptr<BPLayerInterfaceImpl> broad_phase_layer_interface;
	std::unique_ptr<ObjectVsBroadPhaseLayerFilterImpl> object_vs_broadphase_layer_filter;
	std::unique_ptr<ObjectLayerPairFilterImpl> object_vs_object_layer_filter;

	~PhysicsSystemWrapper() = default;
};

JoltPhysicsSystem JoltCreatePhysicsSystem()
{
	// ref: https://github.com/godotengine/godot/blob/e47fb8b8989fd5589c65c4b0ac980de2e936c041/modules/jolt_physics/jolt_project_settings.cpp#L71
	const uint cMaxBodies = 10240;
	const uint cNumBodyMutexes = 0;
	const uint cMaxBodyPairs = 65536;
	const uint cMaxContactConstraints = 20480;

	// Create wrapper to hold PhysicsSystem and layer interfaces
	auto wrapper = std::make_unique<PhysicsSystemWrapper>();

	// Create layer interfaces using smart pointers
	wrapper->broad_phase_layer_interface = std::make_unique<BPLayerInterfaceImpl>();
	wrapper->object_vs_broadphase_layer_filter = std::make_unique<ObjectVsBroadPhaseLayerFilterImpl>();
	wrapper->object_vs_object_layer_filter = std::make_unique<ObjectLayerPairFilterImpl>();

	// Create physics system
	wrapper->system = std::make_unique<PhysicsSystem>();
	wrapper->system->Init(cMaxBodies, cNumBodyMutexes, cMaxBodyPairs, cMaxContactConstraints,
						  *wrapper->broad_phase_layer_interface,
						  *wrapper->object_vs_broadphase_layer_filter,
						  *wrapper->object_vs_object_layer_filter);

	// Release ownership to caller (Go will manage lifetime via JoltDestroyPhysicsSystem)
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

// C++ only: Accessor functions for wrapper internals
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
