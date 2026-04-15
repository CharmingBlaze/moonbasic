/*
 * Forward declarations for symbols defined in physics.cpp.
 * Every wrapper .cpp compiles as its own translation unit (CI builds a static
 * archive from *.cpp); they cannot rely on MSVC-specific loose typing.
 */

#ifndef JOLT_WRAPPER_PHYSICS_BRIDGE_H
#define JOLT_WRAPPER_PHYSICS_BRIDGE_H

#include <Jolt/Jolt.h>

namespace JPH {
class PhysicsSystem;
class ObjectVsBroadPhaseLayerFilter;
class ObjectLayerPairFilter;
}

struct PhysicsSystemWrapper;

JPH::PhysicsSystem *GetPhysicsSystem(PhysicsSystemWrapper *wrapper);
const JPH::ObjectVsBroadPhaseLayerFilter *GetObjectVsBroadPhaseLayerFilter(PhysicsSystemWrapper *wrapper);
const JPH::ObjectLayerPairFilter *GetObjectLayerPairFilter(PhysicsSystemWrapper *wrapper);

#endif
