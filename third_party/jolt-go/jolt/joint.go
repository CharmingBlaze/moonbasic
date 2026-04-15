package jolt

// #include "wrapper/joint.h"
import "C"

// Constraint represents a physics joint
type Constraint struct {
    handle C.JoltConstraint
}

// Destroy removes the constraint from the physics system
func (ps *PhysicsSystem) DestroyConstraint(c *Constraint) {
    C.JoltDestroyConstraint(ps.handle, c.handle)
}

// CreateHingeJoint creates a revolute joint between two bodies.
func (ps *PhysicsSystem) CreateHingeJoint(body1, body2 *BodyID, pivot, axis Vec3) *Constraint {
    handle := C.JoltCreateHingeJoint(
        ps.handle,
        body1.handle,
        body2.handle,
        C.float(pivot.X), C.float(pivot.Y), C.float(pivot.Z),
        C.float(axis.X), C.float(axis.Y), C.float(axis.Z),
    )
    if handle == nil {
        return nil
    }
    return &Constraint{handle: handle}
}
