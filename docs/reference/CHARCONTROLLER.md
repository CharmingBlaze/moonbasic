# Character Controller Commands (`CHARCONTROLLER.*`)

Low-level **heap handle** API around a **capsule character controller**. On **Linux + CGO** this is Jolt’s **`CharacterVirtual`** (full sliding, stairs, ground queries). On **other fullruntime builds** the same **`CHARCONTROLLER.*`** keys are registered against a **lightweight AABB + static-body stub** so scripts compile and run; behavior is simpler than Jolt.

**Documentation order:** [Platform priority](../DEVELOPER.md#platform-priority-windows-then-linux) — when OSes differ, Windows-first notes apply.

For **`CHARACTER.CREATE`** / **`CHARACTERREF.*`** (entity-bound, Jolt on desktop CGO), see [CHARACTER.md](CHARACTER.md). For **`PLAYER.*`** / **`CHAR.*`**, see [PLAYER.md](PLAYER.md) and [KCC.md](KCC.md).

## Core workflow

1. **`PHYSICS3D.START()`** (or **`WORLD.SETUP()`**) and set world gravity as needed.
2. **`CHARCONTROLLER.CREATE(radius#, height#, x#, y#, z#)`** → controller **handle**.
3. Each frame: input → **`CHARCONTROLLER.MOVE(handle, dx#, dy#, dz#)`** (or velocity-driven workflows via **`CHARACTERREF.*`** on Linux when bound to the same backing capsule).
4. Sync visuals: **`CHARCONTROLLER.GETPOS`**, **`CHARCONTROLLER.X` / `.Y` / `.Z`**, or **`CHARCONTROLLER.GETLINEARVEL`** / ground helpers below.
5. **`CHARCONTROLLER.FREE(handle)`** when done.

---

## Creation and lifetime

| Command | Notes |
|--------|--------|
| **`CHARCONTROLLER.CREATE(radius#, height#, x#, y#, z#)`** | Capsule at world position; returns **handle**. Linux+Jolt requires an active **`PHYSICS3D`** session (`CHARCONTROLLER: PHYSICS3D not started` otherwise). |
| **`CHARCONTROLLER.FREE(handle)`** | Releases heap slot; Linux tears down Jolt **`CharacterVirtual`** safely before physics shutdown. |

---

## Pose and motion

| Command | Notes |
|--------|--------|
| **`CHARCONTROLLER.SETPOS(handle, x#, y#, z#)`** (canonical); deprecated **`CHARCONTROLLER.SETPOSITION`** | Set world position (then internal update). |
| **`CHARCONTROLLER.MOVE(handle, dx#, dy#, dz#)`** | Apply displacement; collisions resolved via Jolt extended update (Linux) or stub slide (other OS). |
| **`CHARCONTROLLER.TELEPORT(handle, x#, y#, z#)`** | Snap to position and **clear linear velocity** (useful for spawn / cutscenes). |
| **`CHARCONTROLLER.GETPOS(handle)`** | **Array handle** `[x, y, z]`. |
| **`CHARCONTROLLER.X(handle)`** / **`.Y`** / **`.Z`** | Scalar world components. |

---

## Ground and velocity (Jolt-rich)

These map closely to **`CharacterVirtual`** on **Linux + CGO**. On the **stub** path, values are approximated so gameplay code can stay portable.

| Command | Returns | Meaning |
|--------|---------|---------|
| **`CHARCONTROLLER.ISGROUNDED(handle)`** | `bool` | Supported floor under the capsule. |
| **`CHARCONTROLLER.GROUNDSTATE(handle)`** | `int` | Jolt **`EGroundState`**: **0** OnGround, **1** OnSteepGround, **2** NotSupported, **3** InAir. Stub: **0** or **3** only. |
| **`CHARCONTROLLER.GETLINEARVEL(handle)`** | `[vx, vy, vz]` array | World linear velocity. |
| **`CHARCONTROLLER.GETGROUNDVELOCITY(handle)`** | `[vx, vy, vz]` array | **`GetGroundVelocity()`** — velocity clamped to the ground plane (Jolt). Stub: horizontal components when grounded. |
| **`CHARCONTROLLER.GETGROUNDNORMAL(handle)`** | `[nx, ny, nz]` array | Contact normal under the capsule; stub returns **up** when grounded else **zero**. |

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
