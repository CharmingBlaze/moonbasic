/*
 * Jolt Physics C Wrapper - Physics System
 *
 * The primary entry point for the physics simulation.
 */

#ifndef JOLT_WRAPPER_PHYSICS_H
#define JOLT_WRAPPER_PHYSICS_H

#ifdef __cplusplus
extern "C" {
#endif

// Opaque pointer to the physics system wrapper
typedef void* JoltPhysicsSystem;

// Initialize global Jolt resources (allocators, job system)
// Returns 1 on success, 0 on failure.
int JoltInit();

// Shutdown global Jolt resources
void JoltShutdown();

// Create a new physics system
JoltPhysicsSystem JoltCreatePhysicsSystem();

// Destroy a physics system
void JoltDestroyPhysicsSystem(JoltPhysicsSystem system);

// Update the physics system (step the simulation)
void JoltPhysicsSystemUpdate(JoltPhysicsSystem system, float deltaTime);

// Collision Event Struct
typedef struct {
    void* body1;
    void* body2;
} JoltCollisionEvent;

// Drain the global contact queue of the physics system.
// Returns the number of events written to the 'out' buffer.
int JoltPhysicsSystemDrainContactQueue(JoltPhysicsSystem system, JoltCollisionEvent* out, int maxCount);

#ifdef __cplusplus
}
#endif

#endif // JOLT_WRAPPER_PHYSICS_H
