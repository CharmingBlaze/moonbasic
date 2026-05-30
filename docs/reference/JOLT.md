# Jolt Commands

Low-level Jolt Physics primitives. Prefer `BODY3D.*`, `PHYSICS3D.*`, and `SHAPE.*` for normal use — `JOLT.*` exposes the raw engine layer for advanced scenarios or when you need direct control.

Requires **CGO + Jolt** (Windows/Linux fullruntime).

## Core Workflow

1. `JOLT.INIT()` — initialise Jolt (usually done automatically by `PHYSICS3D.START`).
2. Create shapes: `JOLT.SHAPEBOX`, `JOLT.SHAPESPHERE`, `JOLT.SHAPECAPSULE`, `JOLT.SHAPECYLINDER`, `JOLT.SHAPEMESH`.
3. Create bodies: `JOLT.BODYCREATESTATIC` / `JOLT.BODYCREATEDYNAMIC` / `JOLT.BODYCREATEKINEMATIC`.
4. Add constraints: `JOLT.CONSTRAINTHINGE`, `JOLT.CONSTRAINTPOINT`, etc.
5. `JOLT.STEP(dt)` — advance the simulation.
6. `JOLT.SHUTDOWN()` — tear down Jolt.

---

## Lifecycle

### `JOLT.INIT()` 

Initialises the Jolt Physics engine. Called automatically by `PHYSICS3D.START`; only call manually if bypassing that.

---

### `JOLT.SHUTDOWN()` 

Shuts down Jolt and frees all physics resources.

---

### `JOLT.STEP(dt)` 

Steps the simulation by `dt` seconds.

---

### `JOLT.SETGRAVITY(x, y, z)` 

Sets the global gravity vector.

---

## Shapes

### `JOLT.SHAPEBOX(halfX, halfY, halfZ)` 

Creates a box collision shape (half-extents). Returns a shape handle.

---

### `JOLT.SHAPESPHERE(radius)` 

Creates a sphere shape. Returns a shape handle.

---

### `JOLT.SHAPECAPSULE(radius, halfHeight)` 

Creates a capsule shape. Returns a shape handle.

---

### `JOLT.SHAPECYLINDER(radius, halfHeight)` 

Creates a cylinder shape. Returns a shape handle.

---

### `JOLT.SHAPEMESH(meshHandle)` 

Creates a triangle mesh shape from a mesh handle (static bodies only). Returns a shape handle.

---

## Bodies

### `JOLT.BODYCREATESTATIC()` 

Creates a static (immovable) body. Returns a body handle.

---

### `JOLT.BODYCREATEDYNAMIC()` 

Creates a dynamic (physics-simulated) body. Returns a body handle.

---

### `JOLT.BODYCREATEKINEMATIC()` 

Creates a kinematic body (moved by code, affects dynamic bodies). Returns a body handle.

---

## Constraints

### `JOLT.CONSTRAINTHINGE(bodyA, bodyB)` 

Creates a hinge constraint between two bodies.

---

### `JOLT.CONSTRAINTPOINT(bodyA, bodyB)` 

Creates a point (ball-and-socket) constraint.

---

### `JOLT.CONSTRAINTFIXED(bodyA, bodyB)` 

Creates a fixed (weld) constraint — no relative motion.

---

### `JOLT.CONSTRAINTSLIDER(bodyA, bodyB)` 

Creates a slider (prismatic) constraint.

---

### `JOLT.CONSTRAINTDISTANCE(bodyA, bodyB)` 

Creates a distance constraint — maintains fixed separation.

---

## Queries

### `JOLT.RAYCAST(ox, oy, oz, dx, dy, dz)` 

Performs a raycast from origin `(ox, oy, oz)` in direction `(dx, dy, dz)`. Returns hit result data readable via `PICK.*`.

---

### `JOLT.COLLISIONQUERY(bodyHandle)` 

Runs a collision query for the given body. Returns overlapping contacts.

---

## Full Example

Direct Jolt usage bypassing PHYSICS3D.

```basic
; Advanced use: raw Jolt init + step loop
WINDOW.OPEN(960, 540, "Jolt Demo")
WINDOW.SETFPS(60)

JOLT.INIT()
JOLT.SETGRAVITY(0, -9.81, 0)

floor  = JOLT.BODYCREATESTATIC()
sphere = JOLT.BODYCREATEDYNAMIC()

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 5, -10)
CAMERA.SETTARGET(cam, 0, 0, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    JOLT.STEP(TIME.DELTA())

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        DRAW3D.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

JOLT.SHUTDOWN()
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [PHYSICS3D.md](PHYSICS3D.md) — high-level world setup (preferred)
- [BODY3D.md](BODY3D.md) — high-level body API
- [SHAPE.md](SHAPE.md) — reusable shape handles
- [JOINT3D.md](JOINT3D.md) — high-level constraints
