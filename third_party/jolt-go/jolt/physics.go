package jolt

// #include "wrapper/physics.h"
// #include "wrapper/body.h"
import "C"

// PhysicsSystem represents a physics simulation world
type PhysicsSystem struct {
	handle C.JoltPhysicsSystem
}

// NewPhysicsSystem creates a new physics world
func NewPhysicsSystem() *PhysicsSystem {
	handle := C.JoltCreatePhysicsSystem()
	return &PhysicsSystem{handle: handle}
}

// Destroy frees the physics system
func (ps *PhysicsSystem) Destroy() {
	C.JoltDestroyPhysicsSystem(ps.handle)
}

// Update advances the simulation by deltaTime seconds
func (ps *PhysicsSystem) Update(deltaTime float32) {
	C.JoltPhysicsSystemUpdate(ps.handle, C.float(deltaTime))
}

// CollisionEvent represents a contact between two bodies
type CollisionEvent struct {
	Body1 *BodyID
	Body2 *BodyID
}

// DrainContactQueue fetches all pending collision events from the physics system.
func (ps *PhysicsSystem) DrainContactQueue(maxCount int) []CollisionEvent {
	cEvents := make([]C.JoltCollisionEvent, maxCount)
	count := int(C.JoltPhysicsSystemDrainContactQueue(ps.handle, &cEvents[0], C.int(maxCount)))
	if count <= 0 {
		return nil
	}

	res := make([]CollisionEvent, count)
	for i := 0; i < count; i++ {
		res[i] = CollisionEvent{
			Body1: &BodyID{handle: C.JoltBodyID(cEvents[i].body1)},
			Body2: &BodyID{handle: C.JoltBodyID(cEvents[i].body2)},
		}
	}
	return res
}
