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

Initialises the physics world. Must be called once before any body creation or world queries.

- **Arguments**: (none)
- **Returns**: (none)

---

### `PHYSICS3D.STOP()` 

Shuts down the simulation and frees all internal resources.

- **Arguments**: (none)
- **Returns**: (none)

---

### `PHYSICS3D.UPDATE()` 

Advances the simulation one step. Call once per frame in your main loop.

- **Arguments**: (none)
- **Returns**: (none)
- **Alias**: `PHYSICS3D.STEP`

- **Example**:
  ```basic
  WHILE NOT WINDOW.SHOULDCLOSE()
      PHYSICS3D.UPDATE()
      RENDER.FRAME()
  WEND
  ```

---

### `PHYSICS3D.SETGRAVITY(x, y, z)` 

Sets the global gravity vector for the world.

- **Arguments**:
  - `x` (float): Gravity on the X axis.
  - `y` (float): Gravity on the Y axis (vertical).
  - `z` (float): Gravity on the Z axis.
- **Returns**: (none)

- **Example**:
  ```basic
  ; Earth-like gravity
  PHYSICS3D.SETGRAVITY(0, -9.81, 0)
  ```

---

## Body Creation

### `BODY3D.CREATE(type)` 

Creates a new body definition. This describes a body's properties before it is committed to the world.

- **Arguments**:
  - `type` (string): The motion type. Can be `"STATIC"`, `"DYNAMIC"`, or `"KINEMATIC"`.
- **Returns**: (handle) A handle to the new body definition.
- **Alias**: `BODY3D.MAKE` (deprecated)

- **Example**:
  ```basic
  def = BODY3D.CREATE("DYNAMIC")
  ```

---

### `BODY3D.ADDBOX(def, w, h, d)` / `ADDSPHERE` / `ADDCAPSULE`
Adds a primitive collision shape to a body definition.

- **Arguments**:
    - `def`: (Handle) The body definition.
    - `w, h, d`: (Float) Half-extents or dimensions.
- **Returns**: (Handle) The definition handle (for chaining).

---

### `BODY3D.COMMIT(def, x, y, z)`
Finalizes the definition and spawns the body in the world.

- **Arguments**:
    - `def`: (Handle) The body definition (consumed).
    - `x, y, z`: (Float) Initial world position.
- **Returns**: (Handle) The active physics body handle.
- **Example**:
    ```basic
    body = def.commit(0, 10, 0)
    ```

---

### `BODY3D.SETPOS(body, x, y, z)` / `SETROT`
Teleports the body or sets its Euler rotation (radians).

- **Returns**: (Handle) The body handle (for chaining).

---

### `BODY3D.GETPOS(body)` / `GETROT`
Returns the world position or Euler rotation.

- **Returns**: (Float, Float, Float) Destructured values.

---

---

### `BODY3D.ADDSPHERE(def, radius)` 

Adds a sphere collision shape to a body definition.

- **Arguments**:
  - `def` (handle): The body definition handle.
  - `radius` (float): The radius of the sphere.
- **Returns**: (handle) Returns the definition handle for chaining.

---

### `BODY3D.ADDCAPSULE(def, radius, height)` 

Adds a capsule collision shape to a body definition.

- **Arguments**:
  - `def` (handle): The body definition handle.
  - `radius` (float): The horizontal radius.
  - `height` (float): The total vertical height.
- **Returns**: (handle) Returns the definition handle for chaining.

---

### `BODY3D.ADDMESH(def, meshHandle)` 

Adds a static mesh collision shape. Only valid for `STATIC` bodies.

- **Arguments**:
  - `def` (handle): The body definition handle.
  - `meshHandle` (handle): The handle to a mesh (e.g. from `MESH.LOAD`).
- **Returns**: (handle) Returns the definition handle for chaining.

---

## Body Properties

### `BODY3D.SETMASS(body, mass)` 

Sets the mass of a dynamic body.

- *Handle shortcut*: `body.setMass(mass)`
- **Arguments**:
  - `body` (handle): The active physics body.
  - `mass` (float): Mass in kilograms.
- **Returns**: (handle) Returns the body handle for chaining.

---

### `BODY3D.SETBOUNCE(body, restitution)` 

Sets the "bounciness" of a body.

- *Handle shortcut*: `body.setBounce(restitution)`
- **Arguments**:
  - `body` (handle): The active physics body.
  - `restitution` (float): Range 0.0 (no bounce) to 1.0 (perfectly elastic).
- **Returns**: (handle) Returns the body handle for chaining.

---

### `BODY3D.SETFRICTION(body, friction)` 

Sets the surface friction of a body.

- *Handle shortcut*: `body.setFriction(friction)`
- **Arguments**:
  - `body` (handle): The active physics body.
  - `friction` (float): Higher values result in more friction.
- **Returns**: (handle) Returns the body handle for chaining.

---

## Motion & Forces

### `BODY3D.APPLYFORCE(body, x, y, z)` 

Applies a continuous force to a body.

- *Handle shortcut*: `body.applyForce(x, y, z)`
- **Arguments**:
  - `body` (handle): The active physics body.
  - `x, y, z` (float): The force vector.
- **Returns**: (none)

---

### `BODY3D.APPLYIMPULSE(body, x, y, z)` 

Applies an immediate momentum change (like a hit).

- *Handle shortcut*: `body.applyImpulse(x, y, z)`
- **Arguments**:
  - `body` (handle): The active physics body.
  - `x, y, z` (float): The impulse vector.
- **Returns**: (none)

---

### `BODY3D.SETLINEARVEL(body, x, y, z)` 

Directly sets the velocity of a dynamic body.

- *Handle shortcut*: `body.setLinearVel(x, y, z)`
- **Arguments**:
  - `body` (handle): The active physics body.
  - `x, y, z` (float): Velocity in units per second.
- **Returns**: (none)
- **Alias**: `BODY3D.SETVELOCITY`

---

## World Queries

### `PHYSICS3D.RAYCAST(ox, oy, oz, dx, dy, dz, maxDist)` 

Casts a ray into the physics world and returns the first hit.

- **Arguments**:
  - `ox, oy, oz` (float): Origin of the ray.
  - `dx, dy, dz` (float): Direction vector.
  - `maxDist` (float): Maximum search distance.
- **Returns**: (array) A 1D float array handle `[hit, nx, ny, nz, fraction]`.
  - `hit`: 1.0 if something was hit, 0.0 otherwise.
  - `nx, ny, nz`: Surface normal at hit point.
  - `fraction`: Distance percentage (0.0 to 1.0) along `maxDist`.

---

## Entity Bridge

### `ENTITY.LINKPHYSBUFFER(entity, bufferIndex)` 

Links a visual entity to a physics body for automatic synchronization.

- **Arguments**:
  - `entity` (handle): The visual entity.
  - `bufferIndex` (int): The index obtained from `BODY3D.BUFFERINDEX(body)`.
- **Returns**: (none)

- **Example**:
  ```basic
  body = def.commit(0, 10, 0)
  ent  = ENTITY.CREATE(model)
  ENTITY.LINKPHYSBUFFER(ent, body.bufferIndex())
  ```

---

## Full Example: Falling Cube

```basic
WINDOW.OPEN(960, 540, "3D Physics")
PHYSICS3D.START()
PHYSICS3D.SETGRAVITY(0, -10, 0)

cam = CAMERA.CREATE().pos(0, 10, 20).look(0, 0, 0)

; Floor
floorDef = BODY3D.CREATE("STATIC")
BODY3D.ADDBOX(floorDef, 25, 0.5, 25)
floorBody = BODY3D.COMMIT(floorDef, 0, -0.5, 0)

; Dynamic cube
cubeDef = BODY3D.CREATE("DYNAMIC")
BODY3D.ADDBOX(cubeDef, 1, 1, 1)
cubeBody = BODY3D.COMMIT(cubeDef, 0, 15, 0)

floorMesh = MESH.CREATECUBE(50, 1, 50)
cubeMesh  = MESH.CREATECUBE(2, 2, 2)
mat       = MATERIAL.CREATEDEFAULT()

WHILE NOT WINDOW.SHOULDCLOSE()
    PHYSICS3D.UPDATE()

    ; Read position/rotation from physics body
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

- [CHARACTER_PHYSICS.md](CHARACTER_PHYSICS.md) — Capsule character controller (KCC).
- [PHYSICS_ADVANCED.md](PHYSICS_ADVANCED.md) — Joints, constraints, and vehicles.
- [ENTITY.md](ENTITY.md) — Entity system and `LINKPHYSBUFFER`.
- [docs/PHYSICS.md](../PHYSICS.md) — Beginner's guide to 3D physics.
