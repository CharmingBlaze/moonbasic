# Advanced Physics Commands

Powerful commands for high-level physics machinery, constraints, and automated world interactions.

## Core Workflow

Advanced physics typically involves connecting dynamic bodies with joints to create complex mechanisms.

1. **Create Bodies**: Define your rigid parts using `BODY3D.CREATE()` or `ENTITY.ADDPHYSICS()`.
2. **Connect Joints**: Link parts using `JOINT.CREATEHINGE()` or `JOINT.CREATEPOINT()`.
3. **Configure World**: Enable automated behaviors like `WATER.AUTOPHYSICS()`.

```basic
hero = ENTITY.LOADMESH("hero.glb")
hero.AddPhysics("DYNAMIC", "CAPSULE")
hero.LockAxis(48) // Lock X/Z rotation (keep upright)

door = ENTITY.CREATECUBE()
door.Scale(0.1, 2, 1)
door.AddPhysics("DYNAMIC", "BOX")

// Create a hinge at the door edge
JOINT.CREATEHINGE(door, 0, 0, 0, 0, 1, 0)
```

---

## Joints & Constraints

### `JOINT.CREATEHINGE(b1, b2, px, py, pz, ax, ay, az)`
Creates a hinge joint between two bodies.
- `b1`, `b2`: The two bodies to connect (handles).
- `px, py, pz`: Pivot point in world space.
- `ax, ay, az`: Rotation axis (e.g., `0, 1, 0` for a vertical door hinge).

### `JOINT.CREATEPOINT(b1, b2, px, py, pz)`
Creates a point-to-point (ball and socket) joint.
- `px, py, pz`: Pivot point where the two bodies meet.

### `JOINT.FREE(joint)`
Destroys a joint. The connected bodies remain but are no longer constrained.

---

## Advanced Body Control

### `BODY3D.LOCKAXIS(body, flags)`
Locks specific axes for translation or rotation.
- `flags`: Sum of 1 (X), 2 (Y), 4 (Z) for linear; 8 (X), 16 (Y), 32 (Z) for angular.
- *Handle Shortcut*: `e.LockAxis(flags)`

### `BODY3D.SETDAMPING(body, linear, angular)`
Sets air resistance. High values make objects feel "heavy" in air or honey.
- *Handle Shortcut*: `e.SetDamping(lin, ang)`

### `BODY3D.SETGRAVITYFACTOR(body, factor)`
Scales gravity for a single body. `0` makes it weightless; `2.0` makes it twice as heavy.
- *Handle Shortcut*: `e.SetGravityFactor(factor)`

### `BODY3D.SETCCD(body, toggle)`
Enables Continuous Collision Detection. Prevents fast-moving bullets from phasing through walls.
- *Handle Shortcut*: `e.SetCCD(1)`

---

## Automated Systems

### `WATER.AUTOPHYSICS(toggle)`
When enabled, all physics-driven entities automatically receive buoyancy forces when they enter a `WATER` volume.
- Force is proportional to the submerged volume.
- Provides realistic floating for crates, barrels, and boats.

---

## Platform Parity

Desktop **Windows and Linux** with **`CGO_ENABLED=1`** and linked Jolt use the **same** native physics path. **`!cgo`** or missing Jolt libs use **stubs** (limited or error-returning APIs).

| Feature | With native Jolt (desktop CGO) | Stub / no Jolt |
| :--- | :--- | :--- |
| **Joints** | Placeholder handles (wrapper growth pending) | No-op / errors |
| **LockAxis** | Parsed; some setters not implemented | Stub |
| **Buoyancy** | Volume-based where wired | Scripted / limited |
| **CCD** | Supported where exposed in wrapper | Not available |
