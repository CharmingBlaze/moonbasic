# MoonBASIC Physics Ergonomics (Jolt)

This document describes the modern, simplified physics API for MoonBASIC. These commands are designed to be "Entity-First," abstracting away the complexity of the underlying Jolt physics system.

## 1. Quick Setup (`ENTITY.PHYSICS`)

The fastest way to add physics to an entity.

```basic
; One-line dynamic box setup (auto-sizes to model)
ENTITY.PHYSICS(cube, "BOX", 1.0)

; Static sphere setup
ENTITY.PHYSICS(ball, "SPHERE", 0.0)
```

**Parameters:**
- `id`: The entity handle.
- `type`: `"BOX"`, `"SPHERE"`, `"CAPSULE"`, or `"MESH"`.
- `mass`: `1.0` for dynamic, `0.0` for static.

---

## 2. Breaking Up Long Commands (Modular Setup)

If you need fine-tuned control without long argument lists, use the `PHYSICS.*` builder suite.

### Old Way (Low-Level)
```basic
b = BODY3D.MAKE("dynamic")
BODY3D.ADDBOX(b, 1.0, 0.5, 1.0)
bh = BODY3D.COMMIT(b, 0, 10, 0)
ENTITY.LINKPHYSBUFFER(myEnt, BODY3D.BUFFERINDEX(bh))
```

### New Way (Ergonomic)
```basic
PHYSICS.SHAPE(myEnt, "BOX")
PHYSICS.SIZE(myEnt, 1.0, 0.5, 1.0)
PHYSICS.FRICTION(myEnt, 0.8)
PHYSICS.BOUNCE(myEnt, 0.2)
PHYSICS.BUILD(myEnt, 1.0)
```

---

## 3. Interaction Helpers

You can now interact with physics entities directly using their entity ID. No more tracking body handles!

- **PHYSICS.IMPULSE(id, x, y, z)**: Apply an instant physical force.
- **PHYSICS.VELOCITY(id, x, y, z)**: Set the linear velocity.
- **PHYSICS.WAKE(id)**: Wake up an entity if it has fallen asleep (deactivated) to save CPU.

---

## 4. Automatic Dimension Calculation

When using `ENTITY.PHYSICS(id, "BOX")`, the engine automatically looks at the entity's `w, h, d` properties or model bounding box to determine the collision shape size. This eliminates the need to manually pass dimensions for every object in your scene.
