/*
 * Jolt Physics C Wrapper - Physics System
 *
 * Handles physics world creation, management, and collision layers.
 */

#ifndef JOLT_WRAPPER_PHYSICS_H
#define JOLT_WRAPPER_PHYSICS_H

#ifdef __cplusplus
extern "C" {
#endif

// Opaque pointer types
typedef void* JoltPhysicsSystem;

// Create a new physics world
JoltPhysicsSystem JoltCreatePhysicsSystem();

// Destroy a physics world
void JoltDestroyPhysicsSystem(JoltPhysicsSystem system);

// Step the physics simulation by deltaTime seconds
void JoltPhysicsSystemUpdate(JoltPhysicsSystem system, float deltaTime);

#ifdef __cplusplus
}

// C++ only: Accessor functions for wrapper internals (used by character.cpp)
namespace JPH {
    class PhysicsSystem;
    class ObjectVsBroadPhaseLayerFilter;
    class ObjectLayerPairFilter;
}

struct PhysicsSystemWrapper;  // Opaque forward declaration

// Accessor functions
JPH::PhysicsSystem* GetPhysicsSystem(PhysicsSystemWrapper* wrapper);
const JPH::ObjectVsBroadPhaseLayerFilter* GetObjectVsBroadPhaseLayerFilter(PhysicsSystemWrapper* wrapper);
const JPH::ObjectLayerPairFilter* GetObjectLayerPairFilter(PhysicsSystemWrapper* wrapper);

#endif

#endif // JOLT_WRAPPER_PHYSICS_H
