//go:build cgo || (windows && !cgo)
//
// Same mbentity registration on Windows and Linux; Jolt contact data is populated only on Linux+CGO.
// See AGENTS.md “Physics sync & Jolt”.
//
package mbentity

import (
	"fmt"

	"moonbasic/runtime"
	mbphysics3d "moonbasic/runtime/physics3d"
	"moonbasic/vm/value"
)

func registerJoltEntityCollisionAPI(m *Module, r runtime.Registrar) {
	r.Register("ENTITY.COLLISIONLAYER", "entity", runtime.AdaptLegacy(m.entCollisionLayer))
	r.Register("EntityCollisionLayer", "entity", runtime.AdaptLegacy(m.entCollisionLayer))
	r.Register("EntityCollided", "entity", runtime.AdaptLegacy(m.entEntityCollidedPair))
	r.Register("ENTITYPHYSICSTOUCH", "entity", runtime.AdaptLegacy(m.entEntityCollidedPair))
	// Last-contact globals from Jolt (0 arguments). Use ENTITY.COLLISIONX(entity [, index]) for rule-based hits.
	r.Register("PhysicsCollisionNX", "entity", runtime.AdaptLegacy(m.entBridgeCollisionNX))
	r.Register("PhysicsCollisionNY", "entity", runtime.AdaptLegacy(m.entBridgeCollisionNY))
	r.Register("PhysicsCollisionNZ", "entity", runtime.AdaptLegacy(m.entBridgeCollisionNZ))
	r.Register("PhysicsCollisionPX", "entity", runtime.AdaptLegacy(m.entBridgeCollisionPX))
	r.Register("PhysicsCollisionPY", "entity", runtime.AdaptLegacy(m.entBridgeCollisionPY))
	r.Register("PhysicsCollisionPZ", "entity", runtime.AdaptLegacy(m.entBridgeCollisionPZ))
	r.Register("PhysicsCollisionY", "entity", runtime.AdaptLegacy(m.entBridgeCollisionPY))
	r.Register("PhysicsCollisionForce", "entity", runtime.AdaptLegacy(m.entBridgeCollisionForce))
	r.Register("PhysicsContactCount", "entity", runtime.AdaptLegacy(m.entCountPhysicsContacts))
	// Jolt contact count per entity (same as PhysicsContactCount).
	r.Register("CountCollisions", "entity", runtime.AdaptLegacy(m.entCountPhysicsContacts))
}

// entCollisionLayer stores a 0–31 layer id for future Jolt object-layer filtering (simulation filter not wired yet).
func (m *Module) entCollisionLayer(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.COLLISIONLAYER expects (entity, layer)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.COLLISIONLAYER: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.COLLISIONLAYER: unknown entity")
	}
	ly, ok := args[1].ToInt()
	if !ok || ly < 0 || ly > 31 {
		return value.Nil, fmt.Errorf("ENTITY.COLLISIONLAYER: layer must be 0..31")
	}
	e.collisionLayer = uint8(ly)
	return value.Nil, nil
}

// entEntityCollidedPair returns true if two entities had a Jolt-backed contact in the last PHYSICS3D.STEP (Linux+CGO).
func (m *Module) entEntityCollidedPair(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("EntityCollided expects (entity, entity)")
	}
	a, ok1 := m.entID(args[0])
	b, ok2 := m.entID(args[1])
	if !ok1 || !ok2 || a < 1 || b < 1 {
		return value.Nil, fmt.Errorf("EntityCollided: invalid entity")
	}
	if m.store().ents[a] == nil || m.store().ents[b] == nil {
		return value.Nil, fmt.Errorf("EntityCollided: unknown entity")
	}
	_, ok := mbphysics3d.PairCollidedThisFrame(a, b)
	return value.FromBool(ok), nil
}

func (m *Module) entBridgeCollisionNX(args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CollisionNX expects 0 arguments")
	}
	return value.FromFloat(mbphysics3d.LastCollisionData().NX), nil
}

func (m *Module) entBridgeCollisionNY(args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CollisionNY expects 0 arguments")
	}
	return value.FromFloat(mbphysics3d.LastCollisionData().NY), nil
}

func (m *Module) entBridgeCollisionNZ(args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CollisionNZ expects 0 arguments")
	}
	return value.FromFloat(mbphysics3d.LastCollisionData().NZ), nil
}

func (m *Module) entBridgeCollisionPX(args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CollisionPX expects 0 arguments")
	}
	return value.FromFloat(mbphysics3d.LastCollisionData().PX), nil
}

func (m *Module) entBridgeCollisionPY(args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CollisionPY expects 0 arguments")
	}
	return value.FromFloat(mbphysics3d.LastCollisionData().PY), nil
}

func (m *Module) entBridgeCollisionPZ(args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CollisionPZ expects 0 arguments")
	}
	return value.FromFloat(mbphysics3d.LastCollisionData().PZ), nil
}

func (m *Module) entBridgeCollisionForce(args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CollisionForce expects 0 arguments")
	}
	return value.FromFloat(mbphysics3d.LastCollisionData().Force), nil
}

func (m *Module) entCountPhysicsContacts(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CountCollisions expects entity")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("CountCollisions: invalid entity")
	}
	if m.store().ents[id] == nil {
		return value.Nil, fmt.Errorf("CountCollisions: unknown entity")
	}
	n := mbphysics3d.CountCollisionsForEntity(id)
	return value.FromInt(int64(n)), nil
}
