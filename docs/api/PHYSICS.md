# Physics Commands

Commands for 3D and 2D rigid-body physics simulation. moonBASIC uses **Jolt Physics** for 3D simulation and a Box2D-style solver for 2D. The physics engine handles gravity, collisions, forces, impulses, joints, raycasting, and shape queries.

## Core Concepts

- **Body** — A rigid body in the physics world. Can be **dynamic** (moves), **static** (immovable), or **kinematic** (script-driven movement).
- **Shape** — The collision geometry (box, sphere, capsule, cylinder, mesh).
- **Step** — The physics simulation advances in discrete time steps. Call `Physics3D.Update(dt)` every frame.
- **Sync** — When entities have physics bodies, their positions are synced automatically during `Entity.Update(dt)`.
- **Jolt backend** — Desktop builds (Windows/Linux + CGO) use real Jolt physics. Stub builds degrade gracefully with error messages.

---

## 3D Physics Lifecycle

### `Physics3D.Init()`

Initializes the Jolt physics world. Called automatically when the first physics body is created.

---

### `Physics3D.Update(dt)` / `Physics3D.Step(dt)`

Advances the physics simulation by `dt` seconds. Call once per frame.

- `dt` (float) — Delta time.

**How it works:** Runs the Jolt collision detection and solver for one substep, then syncs body transforms back to entities via the `syncEntitiesFromPhysics` callback. Ground probes fire to detect entity grounding state.

```basic
WHILE NOT Window.ShouldClose()
    dt = Time.Delta()
    Entity.Update(dt)
    Physics3D.Update(dt)
    ; ... render ...
WEND
```

---

### `Physics3D.SetGravity(gx, gy, gz)` / `World.Gravity(gx, gy, gz)`

Sets the global gravity vector.

- `gx`, `gy`, `gz` (float) — Gravity components (typical: 0, -9.81, 0).

---

### `Physics3D.GetGravity()`

Returns the current gravity vector.

---

### `Physics3D.Destroy()`

Destroys the physics world and all bodies. Use for scene cleanup.

---

## 3D Body Creation

### `Body3D.Create(shapeHandle, motionType, mass)`

Creates a physics body with a shape, motion type, and mass.

- `shapeHandle` (handle) — Collision shape (from `Shape.Create*`).
- `motionType` (int) — 0 = Static, 1 = Dynamic, 2 = Kinematic.
- `mass` (float) — Mass in kg (ignored for static bodies).

**Returns:** `handle`

```basic
boxShape = Shape.CreateBox(1, 1, 1)
body = Body3D.Create(boxShape, 1, 10.0)  ; Dynamic, 10kg
Body3D.SetPos(body, 0, 5, 0)
```

---

### `Body3D.CreateStatic(shapeHandle)`

Creates a static (immovable) physics body.

---

### `Body3D.CreateKinematic(shapeHandle)`

Creates a kinematic body (moved by script, pushes dynamic bodies).

---

### `Body3D.Free(bodyHandle)`

Removes a body from the physics world and frees its handle.

---

## Body Properties

### `Body3D.SetPos(bodyHandle, x, y, z)` / `body.pos(x, y, z)`

Teleports a body to a new position.

- `x`, `y`, `z` (float) — World position.

**How it works:** For static/kinematic bodies, directly sets the position. For dynamic bodies, deactivates then repositions (teleport, not physical movement).

---

### `Body3D.GetPos(bodyHandle)`

Returns the body's current position.

---

### `Body3D.SetRot(bodyHandle, pitch, yaw, roll)` / `body.rot(p, y, r)`

Sets the body's rotation in Euler degrees.

---

### `Body3D.GetRot(bodyHandle)`

Returns the body's rotation.

---

### `Body3D.SetVel(bodyHandle, vx, vy, vz)` / `body.vel(vx, vy, vz)`

Sets the linear velocity of a dynamic body.

---

### `Body3D.GetVel(bodyHandle)`

Returns the body's linear velocity.

---

### `Body3D.SetAngVel(bodyHandle, wx, wy, wz)`

Sets the angular velocity.

---

### `Body3D.GetAngVel(bodyHandle)`

Returns the angular velocity.

---

### `Body3D.SetMass(bodyHandle, mass)` / `body.mass(m)`

Sets the body mass.

- `mass` (float) — Mass in kilograms.

---

### `Body3D.GetMass(bodyHandle)`

Returns the body's mass.

---

### `Body3D.SetFriction(bodyHandle, friction)` / `body.friction(f)`

Sets the friction coefficient (0.0 = ice, 1.0 = rough).

---

### `Body3D.SetRestitution(bodyHandle, restitution)` / `body.restitution(r)`

Sets the bounciness (0.0 = no bounce, 1.0 = perfect bounce).

**How it works:** Updates the Jolt body's restitution via `SetRestitutionToIndex`.

---

### `Body3D.SetGravityScale(bodyHandle, scale)`

Scales the gravity applied to this specific body. 0 = weightless, 1 = normal, 2 = double gravity.

---

### `Body3D.Activate(bodyHandle)` / `Body3D.Deactivate(bodyHandle)`

Wakes up or puts to sleep a dynamic body. Sleeping bodies don't consume CPU.

---

### `Body3D.IsActive(bodyHandle)`

Returns `TRUE` if the body is awake.

---

## Forces & Impulses

### `Body3D.ApplyForce(bodyHandle, fx, fy, fz)` / `body.force(fx, fy, fz)`

Applies a continuous force to the body's center of mass. Force is in Newtons and is applied over the frame.

**How it works:** Internally computed as `impulse × dt` to match Jolt's impulse-based API.

```basic
; Thruster pushing upward
Body3D.ApplyForce(rocket, 0, 1000, 0)
```

---

### `Body3D.ApplyImpulse(bodyHandle, ix, iy, iz)` / `body.impulse(ix, iy, iz)`

Applies an instantaneous impulse (momentum change) to the body.

```basic
; Explosion knockback
Body3D.ApplyImpulse(enemy, dirX * 50, 20, dirZ * 50)
```

---

### `Body3D.ApplyTorque(bodyHandle, tx, ty, tz)` / `body.torque(tx, ty, tz)`

Applies a rotational force (torque).

---

### `Body3D.ApplyForceAtPoint(bodyHandle, fx, fy, fz, px, py, pz)`

Applies a force at a specific world point, creating both linear and rotational effects.

---

## Collision Shapes

### `Shape.CreateBox(halfX, halfY, halfZ)`

Creates a box collision shape.

- `halfX`, `halfY`, `halfZ` (float) — Half-extents (the box extends ± this much from center).

**Returns:** `handle`

```basic
boxShape = Shape.CreateBox(0.5, 0.5, 0.5)  ; 1x1x1 cube
```

---

### `Shape.CreateSphere(radius)`

Creates a sphere collision shape.

- `radius` (float) — Sphere radius.

**Returns:** `handle`

---

### `Shape.CreateCapsule(radius, height)`

Creates a capsule collision shape (cylinder with hemispherical caps). Common for character controllers.

- `radius` (float) — Capsule radius.
- `height` (float) — Total height including caps.

**Returns:** `handle`

---

### `Shape.CreateCylinder(radius, height)`

Creates a cylinder collision shape.

**Returns:** `handle`

---

### `Shape.GetType(shapeHandle)` / `Shape.GetWidth(shapeHandle)` / `Shape.GetHeight(shapeHandle)` / `Shape.GetDepth(shapeHandle)` / `Shape.GetRadius(shapeHandle)`

Query shape properties.

---

## Raycasting

### `Physics3D.Raycast(originX, originY, originZ, dirX, dirY, dirZ, maxDist)`

Casts a ray into the physics world and returns the first hit.

**Returns:** Hit information (distance, normal, body handle).

**How it works:** Delegates to Jolt's ray cast API. Only tests against bodies that have collision shapes.

```basic
hit = Physics3D.Raycast(camX, camY, camZ, fwdX, fwdY, fwdZ, 100)
IF hit THEN
    ; Process hit
ENDIF
```

---

### `Physics3D.RaycastAll(originX, originY, originZ, dirX, dirY, dirZ, maxDist)`

Returns all bodies hit by a ray (not just the first).

---

### `Pick.RayToWorld(cameraHandle, screenX, screenY)`

Converts screen coordinates to a world ray for mouse picking.

---

### `Pick.ClosestBody(rayHandle)`

Returns the closest physics body intersected by a ray.

---

## Joints

### `Joint3D.CreateFixed(bodyA, bodyB)`

Creates a fixed joint locking two bodies together.

---

### `Joint3D.CreateHinge(bodyA, bodyB, anchorX, anchorY, anchorZ, axisX, axisY, axisZ)`

Creates a hinge (revolute) joint.

---

### `Joint3D.CreateSlider(bodyA, bodyB, axisX, axisY, axisZ)`

Creates a slider (prismatic) joint along an axis.

---

### `Joint3D.CreateDistance(bodyA, bodyB, minDist, maxDist)`

Creates a distance constraint.

---

### `Joint3D.Free(jointHandle)`

Removes a joint.

---

## 2D Physics

### `Physics2D.Step()`

Advances the 2D physics simulation by one frame.

---

### `Body2D.Create(shapeType, x, y, motionType)`

Creates a 2D physics body.

---

### `Body2D.SetPos(bodyHandle, x, y)` / `Body2D.GetPos(bodyHandle)`

Position manipulation.

---

### `Body2D.SetVel(bodyHandle, vx, vy)` / `Body2D.GetVel(bodyHandle)`

Velocity manipulation.

---

### `Body2D.ApplyForce(bodyHandle, fx, fy)` / `Body2D.ApplyImpulse(bodyHandle, ix, iy)`

Force/impulse application.

---

## Easy Mode Shortcuts

| Shortcut | Maps To |
|----------|---------|
| `UPDATEPHYSICS` | `Entity.Update(dt) + Physics3D.Update(dt)` |

---

## Full Example

A physics sandbox with falling objects and raycasting.

```basic
Window.Open(1280, 720, "Physics Demo")
Window.SetFPS(60)

cam = Camera.Create()
cam.pos(0, 10, 20)
cam.look(0, 2, 0)
cam.fov(60)

World.Gravity(0, -9.81, 0)

; Floor
floor = Entity.CreateBox(20, 0.5, 20, 100, 130, 100)
Entity.SetPos(floor, 0, -0.25, 0)
Entity.AddPhysics(floor, "box", 0)

; Spawn crates on click
crateCount = 0

WHILE NOT Window.ShouldClose()
    dt = Time.Delta()

    ; Spawn crate on left click
    IF Input.MousePressed(0) AND crateCount < 50 THEN
        crate = Entity.CreateBox(1, 1, 1, 180, 120, 60)
        Entity.SetPos(crate, RND(-5, 5), 8, RND(-5, 5))
        Entity.AddPhysics(crate, "box", 5)
        Entity.SetBounce(crate, 0.4)
        Entity.SetFriction(crate, 0.6)
        crateCount = crateCount + 1
    ENDIF

    ; Reset scene on R
    IF Input.KeyPressed(KEY_R) THEN
        Entity.FreeEntities()
        crateCount = 0
        ; Recreate floor
        floor = Entity.CreateBox(20, 0.5, 20, 100, 130, 100)
        Entity.SetPos(floor, 0, -0.25, 0)
        Entity.AddPhysics(floor, "box", 0)
    ENDIF

    ; Update
    Entity.Update(dt)
    Physics3D.Update(dt)

    ; Render
    Render.Clear(30, 30, 50)
    Camera.Begin(cam)
        Draw.Grid(20, 1.0)
        Entity.DrawAll()
    Camera.End(cam)

    Draw.Text("Click = Spawn Crate | R = Reset", 10, 10, 18, 255, 255, 255, 255)
    Draw.Text("Crates: " + STR(crateCount), 10, 32, 18, 200, 200, 200, 255)
    Render.Frame()
WEND

Entity.FreeEntities()
Camera.Free(cam)
Window.Close()
```

---

## See Also

- [ENTITY](ENTITY.md) — Entities with physics integration
- [WORLD](WORLD.md) — Gravity and world update
- [PLAYER](PLAYER.md) — Character controller built on Jolt KCC
- [CAMERA](CAMERA.md) — Raycasting from camera
