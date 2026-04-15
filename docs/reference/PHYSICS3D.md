# 3D Physics Commands

Commands for creating and controlling a 3D physics simulation using Jolt Physics.

**Availability:** **`PHYSICS3D.*`** / **`BODY3D.*`** require **Linux + CGO** with [jolt-go](https://github.com/bbitechnologies/jolt-go); other builds return a stub error. **Registry map:** [moonbasic-command-set/physics-3d.md](moonbasic-command-set/physics-3d.md). Commands that are not wired by the vendored binding now return explicit **not implemented** errors instead of silently no-oping.

**Terrain / heightfields:** the vendored binding exposes **box / sphere / capsule / convex hull / mesh** shapes only. A **Jolt `HeightFieldShape`** for [`TERRAIN.*`](TERRAIN.md) heightmaps is not wired yet; align physics with generated meshes or separate bodies when needed.

## Core Workflow

1.  **Initialize**: Start the physics world with `Physics3D.Start()`.
2.  **Create Bodies**: Define and create physics bodies (`Body3D.Make`, `Body3D.AddShape`, `Body3D.Commit`).
3.  **Update**: Advance the simulation each frame with `Physics3D.Step()`.
4.  **Synchronize**: Use `Body3D.GetMatrix()` to update the transforms of your visual objects.
5.  **Cleanup**: Shut down the world with `Physics3D.Stop()` when done.

---

## World Management

### `Physics3D.Start()`
Initializes the 3D physics world. This must be called before any other physics commands.

### `Physics3D.Stop()`
Shuts down the physics simulation and frees all associated resources.

---

### `Physics3D.Step()`
Advances the physics simulation by one step. Call once per frame.

### `Physics3D.SetGravity(x, y, z)`
Sets the global gravity vector for the physics world.

---

### `Body3D.Make(type)`
Creates a body *definition*. `type`: `"static"`, `"dynamic"`, or `"kinematic"`.

### `Body3D.AddBox(def, w, h, d)`
Adds a box collision shape to the body definition.

### `Body3D.Commit(def, x, y, z)`
Finalizes the body and adds it to the world at the specified position. Returns a **body handle**.

---

### `Body3D.SetPos(handle, x, y, z)`
Teleports a physics body to a new world position.

### `Body3D.SetLinearVel(handle, vx, vy, vz)`
Directly sets the linear velocity of a dynamic body.

### `Body3D.GetMatrix(handle)`
Returns a handle to the body's transformation matrix for visual synchronization.

### `Body3D.Free(handle)`
Removes a body from the simulation and frees its memory.

---

## Body Properties & Constraints

For complex mechanical systems, joints, and advanced body control, see the **[Advanced Physics Guide](PHYSICS_ADVANCED.md)**.

### `Body3D.SetDamping(handle, linear, angular)`
Sets air resistance (0..1).  
Current native backend status: **not implemented** (returns runtime error).
- *Handle Shortcut*: `e.SetDamping(lin, ang)`

### `Body3D.LockAxis(handle, flags)`
Locks specific translation/rotation axes.
Current native backend status: **not implemented** (returns runtime error).
- *Handle Shortcut*: `e.LockAxis(flags)`

### `Body3D.SetGravityFactor(handle, factor)`
Scales gravity for this body (e.g., `0` = weightless).
Current native backend status: **not implemented** (returns runtime error).
- *Handle Shortcut*: `e.SetGravityFactor(factor)`

### `Body3D.SetCCD(handle, toggle)`
Enables Continuous Collision Detection.
Current native backend status: **not implemented** (returns runtime error).

---

## Queries (Linux + CGO + Jolt)

### `Physics3D.Raycast(ox, oy, oz, dx, dy, dz, maxDist)`

Casts a ray from `(ox,oy,oz)` in direction `(dx,dy,dz)`. The direction is **scaled** so its length does not exceed `maxDist`.

Returns a **new 1D float array handle** with 6 elements:

| Index | Value |
|------|--------|
| `0` | `1` if something was hit, `0` if miss |
| `1`–`3` | Hit normal (floats), or `0` on miss |
| `4` | Hit **fraction** along the clipped ray (0–1) |
| `5` | Reserved (`0`); a future version may supply a body id |

Free the array when finished if your program retains handles.

---

## `PICK.*` — world picking (Linux + CGO + Jolt)

Short **dot-notation** commands (no long argument lists). Stage a ray, then **`PICK.CAST`**, then read **`PICK.X`** … **`PICK.ENTITY`**.

| Command | Role |
|--------|------|
| **`PICK.ORIGIN(x, y, z)`** | Ray start |
| **`PICK.DIRECTION(dx, dy, dz)`** | Ray direction; **vector length** is used as max travel unless **`PICK.MAXDIST`** is set |
| **`PICK.MAXDIST(d)`** | If set, **normalize** direction and scale to this length |
| **`PICK.LAYERMASK(m)`** | Bit `i` = accept hits on **`ENTITY.COLLISIONLAYER`** `i`; **`0`** = accept all |
| **`PICK.RADIUS(r)`** | Reserved; non-zero is rejected until sphere pick exists |
| **`PICK.CAST()`** | Closest Jolt hit → fills registry; **returns entity** or `0` (entity must be linked via **`LINKPHYSBUFFER`**) |
| **`PICK.FROMCAMERA(cam, sx, sy)`** | Builds ray from Raylib screen position; sets default **`MAXDIST`** if unset |
| **`PICK.SCREENCAST(cam, sx, sy)`** | **`FROMCAMERA`** + **`CAST`** (returns entity) |
| **`PICK.X` … `PICK.Z`** | Last hit world position |
| **`PICK.NX` … `PICK.NZ`** | Last hit surface normal |
| **`PICK.ENTITY`** | Last hit entity |
| **`PICK.DIST`** | Distance along ray |
| **`PICK.HIT`** | Whether the last cast hit |

Register **`ENTITY.COLLISIONLAYER`** for **`PICK.LAYERMASK`** filtering (lookup wired from the entity module).

---

## Entity ↔ Jolt collision bridge (Linux + CGO)

After you **commit** rigid bodies with `BODY3D.COMMIT`, each body gets a shared **matrix buffer index** (`BODY3D.BUFFERINDEX(body)`). Link entities to those slots with **`ENTITY.LINKPHYSBUFFER(entity, bufferIndex)`**. That call also registers the pair for **frame collision queries** (implemented with post-step Jolt shape overlap, not a C++ contact listener—see note below).

Run **`PHYSICS3D.STEP`** each frame, then:

| Command | Meaning |
|--------|---------|
| **`EntityCollided(a, b)`** | `True` if bodies for `a` and `b` overlapped this step (both must be linked). |
| **`CollisionNX` / `NY` / `NZ`** | Approximate world normal from last successful **`EntityCollided`** (center-to-center). |
| **`CollisionPX` / `PY` / `PZ`** | Contact point from the shape query; **`CollisionY`** aliases **`CollisionPY`**. |
| **`CollisionForce`** | Uses penetration depth as a cheap impact proxy (not a true post-solve impulse). |
| **`CountCollisions(e)`** | Number of distinct overlapping **other** entities for `e` this frame (separate from legacy **`COUNTCOLLISIONS`**). |

**Ordering:** Collision events are collected **at the end of** `PHYSICS3D.STEP` (after matrix sync). Call **`EntityCollided`** in your game loop **after** stepping physics.

**Ghost entities:** **`ENTITY.FREE`** / **`ENTITY.CLEARPHYSBUFFER`** unregister the bridge so freed ids are not reported.

**`ENTITY.COLLISIONLAYER(e, layer)`** stores `0..31` for **`PICK.LAYERMASK`** filtering (see **`PICK.*`** above). The Jolt simulation still uses the vendored two-layer broadphase until `third_party` is extended.

**Note:** The repository’s `third_party/jolt-go` C++ layer is not modified here; the bridge uses **`CollideShapeGetHits`** overlap data so it stays buildable without rebuilding static Jolt libs.

---

## Full Example: Falling Cube

```basic
Window.Open(960, 540, "3D Physics Example")
Window.SetFPS(60)

; 1. Initialize Physics World
Physics3D.Start()
Physics3D.SetGravity(0, -10, 0)

; Setup camera
cam = CreateCamera()
cam.SetPos(0, 10, 20)
cam.SetTarget(0, 0, 0)

; 2. Create a static floor body
floor_def = Body3D.Make("static")
Body3D.AddBox(floor_def, 50, 1, 50)
floor_body = Body3D.Commit(floor_def, 0, -1, 0)

; 3. Create a dynamic cube body
cube_def = Body3D.Make("dynamic")
Body3D.AddBox(cube_def, 2, 2, 2)
cube_body = Body3D.Commit(cube_def, 0, 15, 0)
Body3D.SetMass(cube_body, 1.0)

; Create visual meshes to match the physics shapes
floor_mesh = Mesh.MakeCube(50, 1, 50)
cube_mesh = Mesh.MakeCube(2, 2, 2)
default_mat = Material.MakeDefault()

WHILE NOT Window.ShouldClose()
    ; 4. Update the physics simulation
    Physics3D.Step()

    Render.Clear(10, 20, 40)
    cam.Begin()
        ; 5. Synchronize visuals with physics
        floor_matrix = Body3D.GetMatrix(floor_body)
        cube_matrix = Body3D.GetMatrix(cube_body)

        Mesh.Draw(floor_mesh, default_mat, floor_matrix)
        Mesh.Draw(cube_mesh, default_mat, cube_matrix)

        Draw.Grid(100, 1.0)
    cam.End()
    Render.Frame()
WEND

; 6. Cleanup
Physics3D.Stop()
Window.Close()
```
