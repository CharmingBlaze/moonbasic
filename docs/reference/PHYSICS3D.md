# 3D Physics Commands

Rigid-body 3D physics simulation using **Jolt Physics**. Requires **Windows or Linux + CGO + Jolt** for full simulation; other builds register the same keys against stubs that return explicit errors.

**Availability:** **`PHYSICS3D.*`** / **`BODY3D.*`** keys require **CGO + Jolt** ([jolt-go](https://github.com/bbitechnologies/jolt-go)). Registry map: [moonbasic-command-set/physics-3d.md](moonbasic-command-set/physics-3d.md).

**Terrain / heightfields:** vendored binding exposes **box / sphere / capsule / convex hull / mesh** shapes only. Jolt `HeightFieldShape` is not yet wired; align physics with generated meshes or separate bodies.

## Core Workflow

1. **`PHYSICS3D.START()`** — initialise the world once before any body creation.
2. **`BODY3D.CREATE(type)`** → add shapes → **`BODY3D.COMMIT(def, x, y, z)`** — build each body.
3. Each frame: **`PHYSICS3D.UPDATE()`** — advance simulation.
4. Read **`BODY3D.GETPOS`** / **`BODY3D.GETROT`** or use **`ENTITY.LINKPHYSBUFFER`** to sync visuals.
5. **`PHYSICS3D.STOP()`** — tear down when done.

### Method chaining 

All **`BODY3D.*`** mutating builtins return the body handle on success, so setters chain: `body.setPos(x,y,z).activate()`. See [UNIVERSAL_HANDLE_METHODS.md](UNIVERSAL_HANDLE_METHODS.md).

---

## World Management

### `PHYSICS3D.START()` 

Initialises the physics world. Must be called before any body creation or world queries.

---

### `PHYSICS3D.STOP()` 

Shuts down the simulation and frees all internal resources.

---

### `PHYSICS3D.UPDATE()` / `PHYSICS3D.STEP()` 

Advances the simulation one step. Call once per frame. Both keys call the same handler; **`UPDATE`** is preferred in new scripts (see [API_STANDARDIZATION_DIRECTIVE.md](../API_STANDARDIZATION_DIRECTIVE.md)).

---

### `PHYSICS3D.SETGRAVITY(x, y, z)` 

Sets the global gravity vector. Typical: `PHYSICS3D.SETGRAVITY(0, -10, 0)`.

---

## Body Creation

### `BODY3D.CREATE(type)` 

Creates a body *definition* (not yet in the world). `type`: `"static"`, `"dynamic"`, or `"kinematic"`. Returns a **definition handle**. Alias: `BODY3D.MAKE` (deprecated).

---

### `BODY3D.ADDBOX(def, w, h, d)` 

Adds a box collision shape (half-extents w, h, d) to the definition.

---

### `BODY3D.ADDSPHERE(def, radius)` 

Adds a sphere collision shape to the definition.

---

### `BODY3D.ADDCAPSULE(def, radius, height)` 

Adds a capsule collision shape to the definition.

---

### `BODY3D.ADDCONVEX(def, ...)` 

Adds a convex hull shape. See `commands.json` for vertex arity overloads.

---

### `BODY3D.ADDMESH(def, meshHandle)` 

Adds a mesh collision shape from an existing mesh handle (static bodies only).

---

### `BODY3D.COMMIT(def, x, y, z)` 

Finalises the body definition and inserts it into the world at `(x, y, z)`. Returns a **body handle**. The definition handle is consumed.

- *Handle shortcut*: `def.commit(x, y, z)`

---

### `BODY3D.SETMASS(body, mass)` 

Sets the mass of a dynamic body (kg). Must be called after `COMMIT`.

---

## Body Transform

### `BODY3D.SETPOS(body, x, y, z)` 

Teleports the body to a new world position. Does not wake a sleeping body automatically — call `BODY3D.ACTIVATE` after.

- *Handle shortcut*: `body.setPos(x, y, z)`

---

### `BODY3D.GETPOS(body)` 

Returns the current world position as three floats `x, y, z`.

- *Handle shortcut*: `body.getPos()`

---

### `BODY3D.GETROT(body)` 

Returns the current rotation as Euler angles `pitch, yaw, roll` (radians).

- *Handle shortcut*: `body.getRot()`

---

### `BODY3D.GETPOS` / `BODY3D.GETROT` — visual sync note 

There is no **`BODY3D.GETMATRIX`** builtin. For **`MESH.DRAW`**, combine **`TRANSFORM.TRANSLATION`** with **`TRANSFORM.ROTATION`** via **`TRANSFORM.MULTIPLY`**, or link an entity with **`ENTITY.LINKPHYSBUFFER`** and draw through **`ENTITY.DRAWALL`** / **`ENTITY.DRAW`**.

---

### `BODY3D.BUFFERINDEX(body)` 

Returns the shared matrix buffer index used by **`ENTITY.LINKPHYSBUFFER`**.

- *Handle shortcut*: `body.bufferIndex()`

---

## Body Velocity & Forces

### `BODY3D.SETLINEARVEL(body, vx, vy, vz)` 

Directly sets the linear velocity of a dynamic body. Alias: `BODY3D.SETVELOCITY`.

- *Handle shortcut*: `body.setLinearVel(vx, vy, vz)`

---

### `BODY3D.GETLINEARVEL(body)` 

Returns the current linear velocity as three floats. Alias: `BODY3D.GETVELOCITY`.

- *Handle shortcut*: `body.getLinearVel()`

---

### `BODY3D.APPLYFORCE(body, fx, fy, fz)` 

Applies a continuous force vector (mass-scaled). Alias: `BODY3D.ADDFORCE`.

- *Handle shortcut*: `body.applyForce(fx, fy, fz)`

---

### `BODY3D.APPLYIMPULSE(body, ix, iy, iz)` 

Applies an instant velocity-change impulse. Alias: `BODY3D.ADDIMPULSE`.

- *Handle shortcut*: `body.applyImpulse(ix, iy, iz)`

---

## Body State

### `BODY3D.ACTIVATE(body)` 

Force-wakes a sleeping body. Call after teleporting or after modifying velocity on a static/sleeping body.

- *Handle shortcut*: `body.activate()`

---

### `BODY3D.DEACTIVATE(body)` 

Puts a body to sleep immediately.

- *Handle shortcut*: `body.deactivate()`

---

### `BODY3D.FREE(body)` 

Removes the body from the simulation and frees its heap slot.

- *Handle shortcut*: `body.free()`

---

## Body Properties & Constraints

### `BODY3D.SETBOUNCE(body, restitution)` 

Sets the restitution (bounciness) coefficient (0..1).

- *Handle shortcut*: `body.setBounce(r)`

---

### `BODY3D.SETFRICTION(body, friction)` 

Sets the friction coefficient.

- *Handle shortcut*: `body.setFriction(f)`

---

### `BODY3D.SETDAMPING(body, linear, angular)` 

Sets air-resistance damping (0..1). **Not implemented** in the vendored binding — returns a runtime error.

- *Handle shortcut*: `body.setDamping(lin, ang)`

---

### `BODY3D.LOCKAXIS(body, flags)` 

Locks specific translation/rotation axes via a bitmask. **Not implemented** — returns a runtime error.

- *Handle shortcut*: `body.lockAxis(flags)`

---

### `BODY3D.SETGRAVITYFACTOR(body, factor)` 

Scales gravity for this body (e.g. `0` = weightless). **Not implemented** — returns a runtime error.

- *Handle shortcut*: `body.setGravityFactor(factor)`

---

### `BODY3D.SETCCD(body, toggle)` 

Enables/disables Continuous Collision Detection. **Not implemented** — returns a runtime error.

---

## Per-Body Collision Queries

Call these after **`PHYSICS3D.UPDATE`** each frame.

### `BODY3D.COLLIDED(body)` 

Returns `TRUE` if this body overlapped another body this step.

- *Handle shortcut*: `body.collided()`

---

### `BODY3D.COLLISIONOTHER(body)` 

Returns the **other** body handle from the last collision this step.

- *Handle shortcut*: `body.collisionOther()`

---

### `BODY3D.COLLISIONPOINT(body)` 

Returns the world contact point as a `[x, y, z]` array handle.

- *Handle shortcut*: `body.collisionPoint()`

---

### `BODY3D.COLLISIONNORMAL(body)` 

Returns the contact surface normal as a `[nx, ny, nz]` array handle.

- *Handle shortcut*: `body.collisionNormal()`

---

## Entity ↔ Body Bridge

### `ENTITY.LINKPHYSBUFFER(entity, bufferIndex)` 

Links a visual entity to a Jolt body's matrix buffer slot (index from `BODY3D.BUFFERINDEX`). The engine syncs the entity transform each frame after `PHYSICS3D.UPDATE` and registers the pair for frame collision queries.

---

### Entity-level collision helpers 

After `PHYSICS3D.STEP`, these globals query entity-linked bodies:

| Command | Returns | Notes |
|--------|---------|-------|
| **`EntityCollided(a, b)`** | `bool` | Both entities must be linked via `LINKPHYSBUFFER`. |
| **`CollisionNX/NY/NZ`** | `float` | Approximate contact normal (center-to-center). |
| **`CollisionPX/PY/PZ`** | `float` | Contact point; `CollisionY` aliases `CollisionPY`. |
| **`CollisionForce`** | `float` | Penetration-depth proxy — not a true post-solve impulse. |
| **`CountCollisions(e)`** | `int` | Number of distinct overlapping other entities this frame. |

**Ordering:** collect events **after** `PHYSICS3D.UPDATE`. `ENTITY.FREE` / `ENTITY.CLEARPHYSBUFFER` unregister the bridge.

---

## World Queries

### `PHYSICS3D.RAYCAST(ox, oy, oz, dx, dy, dz, maxDist)` 

Casts a ray from `(ox,oy,oz)` in direction `(dx,dy,dz)` scaled to `maxDist`. Returns a **1D float array handle** `[hit, nx, ny, nz, fraction, 0]`:

- `[0]` — `1` if hit, `0` if miss
- `[1..3]` — hit surface normal, or `0` on miss
- `[4]` — hit fraction along the ray (0–1)
- `[5]` — reserved (future body id)

Free the array handle when done.

---

## `PICK.*` — Screen / World Picking

Stage a ray, cast, then read results. Requires **CGO + Jolt**.

### `PICK.ORIGIN(x, y, z)` 

Sets the ray start point.

---

### `PICK.DIRECTION(dx, dy, dz)` 

Sets the ray direction. Vector length is used as max travel unless `PICK.MAXDIST` is set.

---

### `PICK.MAXDIST(d)` 

If set, normalises the direction and scales to this distance.

---

### `PICK.LAYERMASK(mask)` 

Bit `i` = accept hits on `ENTITY.COLLISIONLAYER` `i`. `0` = accept all.

---

### `PICK.CAST()` 

Fires the staged ray; fills result registers. Returns the hit entity or `0`. Entity must be linked via `ENTITY.LINKPHYSBUFFER`.

---

### `PICK.FROMCAMERA(cam, sx, sy)` 

Builds a ray from Raylib screen position `(sx, sy)`. Sets a default `MAXDIST` if unset.

---

### `PICK.SCREENCAST(cam, sx, sy)` 

`PICK.FROMCAMERA` + `PICK.CAST` in one call. Returns the hit entity or `0`.

---

### Pick result registers 

| Key | Value |
|-----|-------|
| **`PICK.X` / `PICK.Y` / `PICK.Z`** | Last hit world position |
| **`PICK.NX` / `PICK.NY` / `PICK.NZ`** | Last hit surface normal |
| **`PICK.ENTITY`** | Last hit entity |
| **`PICK.DIST`** | Distance along ray |
| **`PICK.HIT`** | `TRUE` if the last cast hit |

---

## Full Example: Falling Cube

```basic
WINDOW.OPEN(960, 540, "3D Physics")
WINDOW.SETFPS(60)

PHYSICS3D.START()
PHYSICS3D.SETGRAVITY(0, -10, 0)

cam = CAMERA.CREATE()
cam.setPos(0, 10, 20)
cam.setTarget(0, 0, 0)

; Floor
floorDef = BODY3D.CREATE("STATIC")
BODY3D.ADDBOX(floorDef, 25, 0.5, 25)
floorBody = BODY3D.COMMIT(floorDef, 0, -0.5, 0)

; Dynamic cube
cubeDef = BODY3D.CREATE("DYNAMIC")
BODY3D.ADDBOX(cubeDef, 1, 1, 1)
cubeBody = BODY3D.COMMIT(cubeDef, 0, 15, 0)
BODY3D.SETMASS(cubeBody, 1.0)

floorMesh = MESH.CREATECUBE(50, 1, 50)
cubeMesh  = MESH.CREATECUBE(2, 2, 2)
mat       = MATERIAL.CREATEDEFAULT()

WHILE NOT WINDOW.SHOULDCLOSE()
    PHYSICS3D.UPDATE()

    cx, cy, cz    = BODY3D.GETPOS(cubeBody)
    cp, cyaw, cr  = BODY3D.GETROT(cubeBody)
    xform = TRANSFORM.MULTIPLY(TRANSFORM.TRANSLATION(cx, cy, cz), TRANSFORM.ROTATION(cp, cyaw, cr))

    RENDER.CLEAR(10, 20, 40)
    RENDER.BEGIN3D(cam)
        MESH.DRAW(floorMesh, mat, TRANSFORM.TRANSLATION(0, -0.5, 0))
        MESH.DRAW(cubeMesh,  mat, xform)
        DRAW.GRID(100, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

BODY3D.FREE(floorBody)
BODY3D.FREE(cubeBody)
PHYSICS3D.STOP()
WINDOW.CLOSE()
```

---

## See also

- [PHYSICS_ADVANCED.md](PHYSICS_ADVANCED.md) — joints, constraints, vehicles
- [PHYSICS2D.md](PHYSICS2D.md) — Box2D 2D simulation
- [CHARCONTROLLER.md](CHARCONTROLLER.md) — capsule character controller
- [ENTITY.md](ENTITY.md) — entity system and `LINKPHYSBUFFER`
- [RAYCAST.md](RAYCAST.md) — non-physics ray queries
