# Character Controller Commands (`CHARCONTROLLER.*`)

Low-level **heap handle** API around a **capsule character controller**. On **Linux + CGO** this is Jolt’s **`CharacterVirtual`** (full sliding, stairs, ground queries). On **other fullruntime builds** the same **`CHARCONTROLLER.*`** keys are registered against a **lightweight AABB + static-body stub** so scripts compile and run; behavior is simpler than Jolt.

**Documentation order:** [Platform priority](../DEVELOPER.md#platform-priority-windows-then-linux) — when OSes differ, Windows-first notes apply.

For **`CHARACTER.CREATE`** / **`CHARACTERREF.*`** (entity-bound or host KCC), see [CHARACTER.md](CHARACTER.md). For **`PLAYER.*`** / **`CHAR.*`**, see [PLAYER.md](PLAYER.md) and [KCC.md](KCC.md).

## Core workflow

1. **`PHYSICS3D.START()`** (or **`WORLD.SETUP()`**) and set world gravity as needed.
2. **`CHARCONTROLLER.MAKE(radius#, height#, x#, y#, z#)`** → controller **handle**.
3. Each frame: input → **`CHARCONTROLLER.MOVE(handle, dx#, dy#, dz#)`** (or velocity-driven workflows via **`CHARACTERREF.*`** on Linux when bound to the same backing capsule).
4. Sync visuals: **`CHARCONTROLLER.GETPOS`**, **`CHARCONTROLLER.X` / `.Y` / `.Z`**, or **`CHARCONTROLLER.GETLINEARVEL`** / ground helpers below.
5. **`CHARCONTROLLER.FREE(handle)`** when done.

---

## Creation and lifetime

| Command | Notes |
|--------|--------|
| **`CHARCONTROLLER.MAKE(radius#, height#, x#, y#, z#)`** | Capsule at world position; returns **handle**. Linux+Jolt requires an active **`PHYSICS3D`** session (`CHARCONTROLLER: PHYSICS3D not started` otherwise). |
| **`CHARCONTROLLER.FREE(handle)`** | Releases heap slot; Linux tears down Jolt **`CharacterVirtual`** safely before physics shutdown. |

---

## Pose and motion

| Command | Notes |
|--------|--------|
| **`CHARCONTROLLER.SETPOS(handle, x#, y#, z#)`** / **`CHARCONTROLLER.SETPOSITION(...)`** | Alias pair: set world position (then internal update). |
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
Window.Open(960, 540, "Character Controller")
Window.SetFPS(60)

; 1. Start Physics
Physics3D.Start()
Physics3D.SetGravity(0, -10, 0)

; Setup camera and floor
cam = Camera.Make()
cam.SetTarget(0, 5, 0)
floor_def = Body3D.Make("static")
Body3D.AddBox(floor_def, 100, 1, 100)
floor_body = Body3D.Commit(floor_def, 0, 0, 0)
floor_mesh = Mesh.MakeCube(100, 1, 100)
mat = Material.MakeDefault()

; 2. Create Controller
player = CHARCONTROLLER.MAKE(0.5, 2.0, 0, 5, 0)
player_mesh = Mesh.MakeCapsule(0.5, 2.0, 16, 16)

WHILE NOT Window.ShouldClose()
    Physics3D.Step()

    ; 3. Update controller from input
    speed = 5.0 * Time.Delta()
    dx = 0
    dz = 0
    IF Input.KeyDown(KEY_W) THEN dz = -speed
    IF Input.KeyDown(KEY_S) THEN dz = speed
    IF Input.KeyDown(KEY_A) THEN dx = -speed
    IF Input.KeyDown(KEY_D) THEN dx = speed
    CHARCONTROLLER.MOVE(player, dx, 0, dz)

    ; 4. Synchronize visuals
    player_x = CHARCONTROLLER.X(player)
    player_y = CHARCONTROLLER.Y(player)
    player_z = CHARCONTROLLER.Z(player)
    cam.SetPos(player_x, player_y + 10, player_z + 15)
    cam.SetTarget(player_x, player_y, player_z)

    player_transform = Transform.Translation(player_x, player_y, player_z)

    Render.Clear(20, 30, 40)
    cam.Begin()
        Mesh.Draw(floor_mesh, mat, Body3D.GetMatrix(floor_body))
        Mesh.Draw(player_mesh, mat, player_transform)
        Draw.Grid(100, 1.0)
    cam.End()
    Render.Frame()
WEND

CHARCONTROLLER.FREE(player)
Physics3D.Stop()
Window.Close()
```
