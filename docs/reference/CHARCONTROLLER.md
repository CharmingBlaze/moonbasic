# Character Controller Commands (`CHARCONTROLLER.*`)

Low-level **heap handle** API around a **capsule character controller**. On **Linux + CGO** this is Jolt’s **`CharacterVirtual`** (full sliding, stairs, ground queries). On **other fullruntime builds** the same **`CHARCONTROLLER.*`** keys are registered against a **lightweight AABB + static-body stub** so scripts compile and run; behavior is simpler than Jolt.

**Documentation order:** [Platform priority](../DEVELOPER.md#platform-priority-windows-then-linux) — when OSes differ, Windows-first notes apply.

For **`CHARACTER.CREATE`** / **`CHARACTERREF.*`** (entity-bound, Jolt on desktop CGO), see [CHARACTER.md](CHARACTER.md). For **`PLAYER.*`** / **`CHAR.*`**, see [PLAYER.md](PLAYER.md) and [KCC.md](KCC.md).

## Core Workflow

1. **`PHYSICS3D.START()`** (or **`WORLD.SETUP()`**) and set world gravity as needed.
2. **`CHARCONTROLLER.CREATE(radius#, height#, x#, y#, z#)`** → controller **handle**.
3. Each frame: input → **`CHARCONTROLLER.MOVE(handle, dx#, dy#, dz#)`** (or velocity-driven workflows via **`CHARACTERREF.*`** on Linux when bound to the same backing capsule).
4. Sync visuals: **`CHARCONTROLLER.GETPOS`**, **`CHARCONTROLLER.X` / `.Y` / `.Z`**, or **`CHARCONTROLLER.GETLINEARVEL`** / ground helpers below.
5. **`CHARCONTROLLER.FREE(handle)`** when done.

---

## Creation and Lifetime

### `CHARCONTROLLER.CREATE(radius, height, x, y, z)` 

Creates a capsule controller at world position `(x, y, z)`. Returns a **handle**. Requires an active `PHYSICS3D` session on Linux+Jolt (`CHARCONTROLLER: PHYSICS3D not started` otherwise). Alias: `CHARCONTROLLER.MAKE` (deprecated).

- *Handle shortcut*: n/a — use returned handle directly

---

### `CHARCONTROLLER.FREE(handle)` 

Releases the heap slot and tears down the Jolt `CharacterVirtual` safely before physics shutdown.

- *Handle shortcut*: `ctrl.free()`

---

## Position

### `CHARCONTROLLER.SETPOS(handle, x, y, z)` 

Sets world position and triggers an internal update. Alias: `CHARCONTROLLER.SETPOSITION` (deprecated).

- *Handle shortcut*: `ctrl.setPos(x, y, z)`

---

### `CHARCONTROLLER.GETPOS(handle)` 

Returns current world position as a `[x, y, z]` array handle.

- *Handle shortcut*: `ctrl.getPos()`

---

### `CHARCONTROLLER.X(handle)` / `CHARCONTROLLER.Y(handle)` / `CHARCONTROLLER.Z(handle)` 

Returns the individual X, Y, or Z world coordinate as a scalar float.

---

### `CHARCONTROLLER.TELEPORT(handle, x, y, z)` 

Snaps to position and **clears linear velocity**. Useful for spawn points and cutscene transitions.

- *Handle shortcut*: `ctrl.teleport(x, y, z)`

---

## Motion

### `CHARCONTROLLER.MOVE(handle, dx, dy, dz)` 

Applies a displacement this frame. Collisions resolved via Jolt extended update (Linux+CGO) or stub slide (other builds).

- *Handle shortcut*: `ctrl.move(dx, dy, dz)`

---

## Ground & Velocity

These map to `CharacterVirtual` on **Windows/Linux + CGO**. On stub builds, values are approximated for portable gameplay code.

### `CHARCONTROLLER.ISGROUNDED(handle)` 

Returns `TRUE` when the capsule has a supported floor beneath it.

- *Handle shortcut*: `ctrl.isGrounded()`

---

### `CHARCONTROLLER.GROUNDSTATE(handle)` 

Returns Jolt `EGroundState`: `0` OnGround, `1` OnSteepGround, `2` NotSupported, `3` InAir. Stub: `0` or `3` only.

- *Handle shortcut*: `ctrl.groundState()`

---

### `CHARCONTROLLER.GETLINEARVEL(handle)` 

Returns world linear velocity as a `[vx, vy, vz]` array handle.

- *Handle shortcut*: `ctrl.getLinearVel()`

---

### `CHARCONTROLLER.GETGROUNDVELOCITY(handle)` 

Returns velocity of the ground surface under the capsule as a `[vx, vy, vz]` array handle. Stub: horizontal components when grounded.

- *Handle shortcut*: `ctrl.getGroundVelocity()`

---

### `CHARCONTROLLER.GETGROUNDNORMAL(handle)` 

Returns the contact normal under the capsule as a `[nx, ny, nz]` array handle. Stub: up-vector when grounded, zero otherwise.

- *Handle shortcut*: `ctrl.getGroundNormal()`

---

## Full Example

```basic
WINDOW.OPEN(960, 540, "Character Controller")
WINDOW.SETFPS(60)

PHYSICS3D.START()
PHYSICS3D.SETGRAVITY(0, -10, 0)

cam = CAMERA.CREATE()
cam.SETTARGET(0, 5, 0)

floorDef = BODY3D.CREATE("STATIC")
BODY3D.ADDBOX(floorDef, 50, 0.5, 50)
floorBody = BODY3D.COMMIT(floorDef, 0, 0, 0)
floorMesh = MESH.CREATECUBE(100, 1, 100)
mat = MATERIAL.CREATEDEFAULT()

player = CHARCONTROLLER.CREATE(0.5, 2.0, 0, 5, 0)
playerMesh = MESH.CREATECAPSULE(0.5, 2.0, 16, 16)

WHILE NOT WINDOW.SHOULDCLOSE()
    PHYSICS3D.UPDATE()

    speed = 5.0 * TIME.DELTA()
    dx = 0
    dz = 0
    IF INPUT.KEYDOWN(KEY_W) THEN dz = -speed
    IF INPUT.KEYDOWN(KEY_S) THEN dz = speed
    IF INPUT.KEYDOWN(KEY_A) THEN dx = -speed
    IF INPUT.KEYDOWN(KEY_D) THEN dx = speed
    CHARCONTROLLER.MOVE(player, dx, 0, dz)

    player_x = CHARCONTROLLER.X(player)
    player_y = CHARCONTROLLER.Y(player)
    player_z = CHARCONTROLLER.Z(player)
    cam.SETPOS(player_x, player_y + 10, player_z + 15)
    cam.SETTARGET(player_x, player_y, player_z)

    playerTx = TRANSFORM.TRANSLATION(player_x, player_y, player_z)

    RENDER.CLEAR(20, 30, 40)
    RENDER.Begin3D(cam)
        MESH.DRAW(floorMesh, mat, TRANSFORM.TRANSLATION(0, 0, 0))
        MESH.DRAW(playerMesh, mat, playerTx)
        DRAW.GRID(100, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

CHARCONTROLLER.FREE(player)
BODY3D.FREE(floorBody)
PHYSICS3D.STOP()
WINDOW.CLOSE()
```

---

## See also

- [CHARACTER.md](CHARACTER.md) — `CHARACTER.CREATE` entity-bound controller (`CHARACTERREF.*`)
- [PHYSICS3D.md](PHYSICS3D.md) — world setup, gravity, `BODY3D.*`
- [KCC.md](KCC.md) — `CHAR.*` / `PLAYER.*` gameplay layer
- [PLAYER.md](PLAYER.md) — `PLAYER.CREATE`, swim, nav
