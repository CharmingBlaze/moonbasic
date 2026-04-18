# Body3D Commands

3D rigid body creation and simulation using Jolt Physics. Build body definitions, add collision shapes, commit to the world, then read state each frame.

Requires **`CGO_ENABLED=1`** and linked Jolt on Windows/Linux. See [PHYSICS3D.md](PHYSICS3D.md) and [PHYSICS_ADVANCED.md](PHYSICS_ADVANCED.md).

## Core Workflow

1. `PHYSICS3D.START()` — initialise the 3D world.
2. `BODY3D.CREATE(type)` — begin a definition (`"DYNAMIC"`, `"STATIC"`, `"KINEMATIC"`).
3. `BODY3D.ADDBOX` / `BODY3D.ADDSPHERE` / `BODY3D.ADDCAPSULE` / `BODY3D.ADDMESH` — attach shape.
4. `BODY3D.COMMIT(def, x, y, z)` — finalise; returns a **live body handle**.
5. Each frame: `PHYSICS3D.UPDATE()` → read `BODY3D.X/Y/Z` / `BODY3D.GETPOS` → sync visuals.
6. `BODY3D.FREE(handle)` when done.

---

## Creation

### `BODY3D.CREATE(type)` 

Begins a body definition. `type` is `"DYNAMIC"`, `"STATIC"`, or `"KINEMATIC"`. Returns a **definition handle**. Call shape-add commands then `BODY3D.COMMIT`.

---

### `BODY3D.ADDBOX(def, halfX, halfY, halfZ)` 

Attaches a box collision shape to the definition. Values are **half-extents** in world units.

---

### `BODY3D.ADDSPHERE(def, radius)` 

Attaches a sphere shape with the given `radius`.

---

### `BODY3D.ADDCAPSULE(def, radius, halfHeight)` 

Attaches a capsule (cylinder with hemispherical caps). `halfHeight` is half the cylindrical section.

---

### `BODY3D.ADDMESH(def, meshHandle)` 

Attaches a triangle mesh shape (static bodies only). `meshHandle` is a loaded mesh handle from `MESH.*`.

---

### `BODY3D.COMMIT(def, x, y, z)` 

Finalises and places the body at world position `(x, y, z)`. Returns the **live body handle**.

---

## Position & Rotation

### `BODY3D.SETPOS(body, x, y, z)` 

Teleports the body to world coordinates. Wakes sleeping bodies.

- *Handle shortcut*: `body.setPos(x, y, z)`

---

### `BODY3D.GETPOS(body)` 

Returns current position as a 3-element array handle `[x, y, z]`.

- *Handle shortcut*: `body.getPos()`

---

### `BODY3D.POS(body)` 

Alias of `BODY3D.GETPOS`. Returns `[x, y, z]` array handle.

---

### `BODY3D.X(body)` / `BODY3D.Y(body)` / `BODY3D.Z(body)` 

Returns a single world coordinate component as a float scalar.

---

### `BODY3D.SETROT(body, pitch, yaw, roll)` 

Sets the body orientation in **degrees** (Euler angles). Returns the body handle for chaining.

- *Handle shortcut*: `body.setRot(pitch, yaw, roll)`

---

### `BODY3D.GETROT(body)` 

Returns orientation as a 3-element array handle `[pitch, yaw, roll]` in degrees.

- *Handle shortcut*: `body.getRot()`

---

### `BODY3D.ROT(body)` 

Alias of `BODY3D.GETROT`.

---

## Scale

### `BODY3D.SETSCALE(body, sx, sy, sz)` 

Scales the collision shape for primitive bodies (box/sphere/capsule). Not supported for mesh bodies.

- *Handle shortcut*: `body.setScale(sx, sy, sz)`

---

### `BODY3D.GETSCALE(body)` 

Returns `[sx, sy, sz]` scale factors. Mesh bodies always return `[1, 1, 1]`.

---

### `BODY3D.SCALE(body)` 

Alias of `BODY3D.GETSCALE`.

---

## Mass & Material

### `BODY3D.SETMASS(body, mass)` 

Sets the body mass in kg. Returns the body handle.

---

### `BODY3D.GETMASS(body)` / `BODY3D.MASS(body)` 

Returns mass in kg.

---

### `BODY3D.SETFRICTION(body, friction)` 

Sets the friction coefficient (0–1). Returns the body handle.

---

### `BODY3D.GETFRICTION(body)` / `BODY3D.FRICTION(body)` 

Returns the current friction coefficient.

---

### `BODY3D.SETRESTITUTION(body, value)` 

Sets the restitution (bounciness) 0–1. Returns the body handle.

- Alias: `BODY3D.BOUNCE` (getter)

---

### `BODY3D.GETRESTITUTION(body)` / `BODY3D.RESTITUTION(body)` / `BODY3D.BOUNCE(body)` 

Returns the restitution value.

---

## Velocity & Forces

### `BODY3D.SETLINEARVEL(body, vx, vy, vz)` 

Sets the linear velocity directly. Returns the body handle.

- *Handle shortcut*: `body.setLinearVel(vx, vy, vz)`

---

### `BODY3D.GETLINEARVEL(body)` / `BODY3D.GETVELOCITY(body)` / `BODY3D.VELOCITY(body)` / `BODY3D.VEL(body)` 

Returns linear velocity as a 3-element array handle `[vx, vy, vz]`.

---

### `BODY3D.SETANGULARVEL(body, ax, ay, az)` 

Sets the angular velocity in radians per second. Returns the body handle.

---

### `BODY3D.GETANGULARVEL(body)` / `BODY3D.ANGULARVEL(body)` 

Returns angular velocity as a 3-element array handle `[ax, ay, az]`.

---

### `BODY3D.APPLYFORCE(body, fx, fy, fz)` 

Applies a continuous force (Newtons) at the body centre this frame.

---

### `BODY3D.APPLYIMPULSE(body, ix, iy, iz)` 

Applies an instantaneous impulse (kg·m/s) at the body centre.

---

### `BODY3D.APPLYTORQUE(body, tx, ty, tz)` 

Applies a rotational torque this frame.

---

## Advanced Properties

### `BODY3D.SETGRAVITYFACTOR(body, factor)` 

Scales gravity for this body only. `0.0` = weightless, `2.0` = double gravity.

- *Handle shortcut*: `body.setGravityFactor(factor)`

---

### `BODY3D.GETGRAVITYFACTOR(body)` 

Returns the gravity scale factor.

---

### `BODY3D.SETDAMPING(body, linear, angular)` 

Sets linear and angular damping (air resistance). Returns the body handle.

- *Handle shortcut*: `body.setDamping(linear, angular)`

---

### `BODY3D.GETDAMPING(body)` 

Returns `[linearDamp, angularDamp]` array handle.

---

### `BODY3D.SETCCD(body, enabled)` 

Enables (`TRUE`) or disables Continuous Collision Detection — prevents fast objects tunnelling through thin geometry.

- *Handle shortcut*: `body.setCCD(1)`

---

### `BODY3D.GETCCD(body)` 

Returns `TRUE` if CCD is currently enabled.

---

### `BODY3D.LOCKAXIS(body, flags)` 

Locks specific motion or rotation axes. `flags` is a bitmask: `1`=X, `2`=Y, `4`=Z (linear); `8`=RotX, `16`=RotY, `32`=RotZ. Combine with `+`. E.g. `48` locks X and Y rotation (keeps upright).

---

### `BODY3D.ACTIVATE(body)` / `BODY3D.DEACTIVATE(body)` 

Manually wake or sleep a body.

---

### `BODY3D.BUFFERINDEX(body)` 

Returns the internal physics buffer index for this body. Useful for advanced sync with the entity bridge.

---

## Collision Queries

### `BODY3D.COLLIDED(body)` 

Returns `1` if the body had a contact during the last physics step.

---

### `BODY3D.COLLISIONOTHER(body)` 

Returns the handle of the other body from the most recent contact.

---

### `BODY3D.COLLISIONPOINT(body)` 

Returns the contact point as a 3-element array handle `[x, y, z]`.

---

### `BODY3D.COLLISIONNORMAL(body)` 

Returns the contact normal as a 3-element array handle `[nx, ny, nz]`.

---

## Lifetime

### `BODY3D.FREE(body)` 

Removes the body from the physics world and frees its resources.

- *Handle shortcut*: `body.free()`

---

## Full Example

A dynamic crate falling onto a static floor with CCD enabled.

```basic
WINDOW.OPEN(960, 540, "Body3D Demo")
WINDOW.SETFPS(60)

PHYSICS3D.START()
PHYSICS3D.SETGRAVITY(0, -10, 0)

; Static floor
floorDef = BODY3D.CREATE("STATIC")
BODY3D.ADDBOX(floorDef, 20, 0.5, 20)
floor = BODY3D.COMMIT(floorDef, 0, -0.5, 0)

; Dynamic crate with CCD
crateDef = BODY3D.CREATE("DYNAMIC")
BODY3D.ADDBOX(crateDef, 0.5, 0.5, 0.5)
BODY3D.SETRESTITUTION(crateDef, 0.4)
crate = BODY3D.COMMIT(crateDef, 0, 10, 0)
BODY3D.SETCCD(crate, TRUE)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 5, -12)
CAMERA.SETTARGET(cam, 0, 0, 0)

mesh = MESH.CREATECUBE(1, 1, 1)
mat  = MATERIAL.CREATEDEFAULT()

WHILE NOT WINDOW.SHOULDCLOSE()
    PHYSICS3D.UPDATE()

    cx = BODY3D.X(crate)
    cy = BODY3D.Y(crate)
    cz = BODY3D.Z(crate)
    t  = TRANSFORM.TRANSLATION(cx, cy, cz)

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        MESH.DRAW(mesh, mat, t)
        DRAW3D.GRID(20, 1.0)
    RENDER.END3D()
    TRANSFORM.FREE(t)
    RENDER.FRAME()
WEND

BODY3D.FREE(crate)
BODY3D.FREE(floor)
MESH.UNLOAD(mesh)
PHYSICS3D.STOP()
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `BODY3D.MAKE(shape, x,y,z)` | Deprecated alias of `BODY3D.CREATE`. |
| `BODY3D.SETPOSITION(body, x,y,z)` | Teleport body to world position (wakes it). |
| `BODY3D.SETVELOCITY(body, vx,vy,vz)` | Set linear velocity directly. |

---

## See also

- [PHYSICS3D.md](PHYSICS3D.md) — world init, gravity, step
- [PHYSICS_ADVANCED.md](PHYSICS_ADVANCED.md) — joints, CCD, lock-axis, damping
- [SHAPE.md](SHAPE.md) — reusable shape handles
- [JOINT3D.md](JOINT3D.md) — hinge, cone, slider constraints
- [ENTITY.md](ENTITY.md) — entity bridge (`ENTITY.ADDPHYSICS`)
