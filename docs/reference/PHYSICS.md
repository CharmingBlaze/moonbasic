# Physics Commands

Entity-level physics shorthand namespace. `PHYSICS.*` operates directly on entity IDs (ints) rather than `BODY3D` handles — the classic "auto-physics" workflow where entities self-manage their physics body.

For full handle-based physics see [PHYSICS3D.md](PHYSICS3D.md), [BODY3D.md](BODY3D.md).

## Core Workflow

1. `PHYSICS.START()` — initialise the world.
2. `PHYSICS.AUTO(entityId, type, mass)` — attach a physics body to an entity automatically.
3. Apply forces/velocity: `PHYSICS.IMPULSE`, `PHYSICS.VELOCITY`, `PHYSICS.FORCE`.
4. `PHYSICS.STEP(dt)` each frame.
5. `PHYSICS.STOP()` on exit.

---

## Lifecycle

### `PHYSICS.START()` 

Initialises the physics world. Alias of `PHYSICS3D.START`.

---

### `PHYSICS.STOP()` 

Shuts down the physics world.

---

### `PHYSICS.STEP(dt)` 

Advances the simulation by `dt` seconds.

---

### `PHYSICS.SETGRAVITY(x, y, z)` 

Sets the global gravity vector.

---

### `PHYSICS.GETGRAVITYX()` / `PHYSICS.GETGRAVITYY()` / `PHYSICS.GETGRAVITYZ()` 

Returns individual gravity components.

---

### `PHYSICS.SETSUBSTEPS(n)` 

Sets the number of simulation sub-steps per frame for higher accuracy.

---

## Auto-Attach

### `PHYSICS.AUTO(entityId, type, mass)` 

Automatically creates and attaches a physics body to `entityId`. `type` is a string: `"DYNAMIC"`, `"STATIC"`, or `"KINEMATIC"`. `mass` in kg.

---

### `PHYSICS.AUTOCREATE(entityId)` 

Creates a default dynamic physics body for the entity from its mesh bounds.

---

### `PHYSICS.SHAPE(entityId, shapeName, ...)` 

Sets the collision shape type for the entity's physics body.

---

### `PHYSICS.SIZE(entityId, sx, sy, sz)` 

Sets the collision shape dimensions.

---

### `PHYSICS.BUILD(entityId, mass)` 

Finalises the physics body for an entity.

---

## Forces & Velocity

### `PHYSICS.IMPULSE(entityId, ix, iy, iz)` 

Applies an instantaneous impulse.

---

### `PHYSICS.VELOCITY(entityId, vx, vy, vz)` 

Sets the linear velocity directly.

---

### `PHYSICS.FORCE(entityId, fx, fy, fz)` 

Applies a continuous force this frame.

---

### `PHYSICS.TORQUE(bodyHandle, tx, ty, tz)` 

Applies torque to a body handle.

---

## Material

### `PHYSICS.FRICTION(entityId, value)` 

Sets the friction coefficient.

---

### `PHYSICS.BOUNCE(entityId, value)` 

Sets the restitution (bounciness).

---

### `PHYSICS.GRAVITY(entityId, scale)` 

Sets per-entity gravity scale.

---

## CCD & Rotation

### `PHYSICS.CCD(entityId, enabled)` 

Enables/disables Continuous Collision Detection for this entity.

---

### `PHYSICS.SETROT(entityId, pitch, yaw, roll)` 

Sets physics body rotation in degrees.

---

## Wake

### `PHYSICS.WAKE(entityId)` 

Wakes a sleeping body.

---

## Queries

### `PHYSICS.RAYCAST(ox, oy, oz, dx, dy, dz, maxDist)` 

Performs a raycast and returns a hit result handle. Read with `PICK.*`.

---

### `PHYSICS.SPHERECAST(...)` / `PHYSICS.BOXCAST(...)` 

Sphere and box sweep queries.

---

### `PHYSICS.ENABLE(entityId)` / `PHYSICS.DISABLE(entityId)` 

Enable or disable physics on an entity.

---

## Buoyancy

### `PHYSICS.SETBUOYANCY(entityId, density)` 

Sets a buoyancy water-density override for the entity.

---

### `PHYSICS.GETBUOYANCY(entityId)` 

Returns the buoyancy density.

---

## Effects

### `PHYSICS.EXPLOSION(x, y, z, radius, force)` 

Applies a radial explosion impulse to all dynamic bodies within `radius`.

---

## Full Example

Crates falling with entity-level physics.

```basic
WINDOW.OPEN(960, 540, "Physics Demo")
WINDOW.SETFPS(60)

PHYSICS.START()
PHYSICS.SETGRAVITY(0, -10, 0)

floor = ENTITY.CREATEPLANE(20, 20)
ENTITY.SETPOS(floor, 0, 0, 0)
PHYSICS.AUTO(floor, "STATIC", 0)
PHYSICS.BUILD(floor, 0)

FOR i = 1 TO 5
    box = ENTITY.CREATECUBE(1.0)
    ENTITY.SETPOS(box, RNDF(-4, 4), i * 2 + 2, RNDF(-4, 4))
    PHYSICS.AUTO(box, "DYNAMIC", 1.0)
    PHYSICS.BUILD(box, 1.0)
NEXT i

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 10, -15)
CAMERA.SETTARGET(cam, 0, 4, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    PHYSICS.STEP(TIME.DELTA())
    ENTITY.UPDATE(TIME.DELTA())

    IF INPUT.KEYPRESSED(KEY_SPACE)
        PHYSICS.EXPLOSION(0, 1, 0, 6, 12)
    END IF

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

PHYSICS.STOP()
WINDOW.CLOSE()
```

---

## See also

- [PHYSICS3D.md](PHYSICS3D.md) — preferred full 3D physics API
- [BODY3D.md](BODY3D.md) — handle-based body creation
- [ENTITY.md](ENTITY.md) — `ENTITY.ADDPHYSICS` bridge
