/*
 * Shared object / broadphase layer IDs for physics.cpp and body.cpp.
 * Broadphase stays at 2 layers; multiple object layers map to MOVING broadphase.
 */
#ifndef JOLT_WRAPPER_PHYSICS_LAYERS_H
#define JOLT_WRAPPER_PHYSICS_LAYERS_H

#include <Jolt/Jolt.h>
#include <Jolt/Physics/Collision/ObjectLayer.h>
#include <Jolt/Physics/Collision/BroadPhase/BroadPhaseLayer.h>

namespace MBPHYS_Layers
{
	static constexpr JPH::ObjectLayer NON_MOVING = 0;
	static constexpr JPH::ObjectLayer MOVING = 1;
	static constexpr JPH::ObjectLayer CHARACTER = 2;
	static constexpr JPH::ObjectLayer SENSOR = 3;
	static constexpr JPH::ObjectLayer ONE_WAY = 4;
	static constexpr JPH::ObjectLayer NUM_LAYERS = 5;
};

namespace MBPHYS_BroadPhaseLayers
{
	static constexpr JPH::BroadPhaseLayer NON_MOVING(0);
	static constexpr JPH::BroadPhaseLayer MOVING(1);
	static constexpr JPH::uint NUM_LAYERS(2);
};

#endif
