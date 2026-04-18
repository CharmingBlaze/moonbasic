# Player Commands

High-level KCC helpers: create, move, jump, swim, query grounding, and spatial lookups.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. `PHYSICS3D.START()` and set gravity.
2. `PLAYER.CREATE(entity, radius, height)` on a positioned capsule entity.
3. Each frame: `PLAYER.MOVE` / `PLAYER.JUMP` / `PLAYER.ISGROUNDED`.
4. Query targets with `PLAYER.GETLOOKTARGET`, nearby entities with `PLAYER.GETNEARBY`.

For character controller tuning and guide see [CHARACTER_PHYSICS.md](CHARACTER_PHYSICS.md).

## Platform

Order: **Windows** first, **Linux** second ([DEVELOPER.md](../DEVELOPER.md#platform-priority-windows-then-linux)).

| Feature | `!cgo` / stub build | Windows + Linux, CGO + Jolt |
|--------|---------------------|------------------------------|
| **`PLAYER.CREATE` / `MOVE` / `JUMP` / `ISGROUNDED` / `SYNCANIM`** | Clear error (requires CGO + Jolt) | Supported |
| **`PLAYER.GETLOOKTARGET` / `GETNEARBY` / `SETSTATE`** | Stub / limited | Full when **`PHYSICS3D.START`** and entity pipeline are active |

Start the world with **`PHYSICS3D.START()`** before **`PLAYER.CREATE`**.

**Gameplay-oriented KCC guide (beginner → advanced):** **[CHARACTER_PHYSICS.md](CHARACTER_PHYSICS.md)** — **`CHAR.MOVE`**, **`CHAR.MOVEWITHCAMERA`**, **`NAVTO` / `NAVUPDATE`**, **`WORLD.MOUSEFLOOR`**, **`WORLD.MOUSEPICK`**, and RPG helpers.

**Entity argument:** Commands that take **`entity`** accept either a **numeric entity id** (`1`, `2`, …) or an **EntityRef handle** from **`MODEL.CREATECAPSULE`**, **`CUBE`**, **`SPHERE`**, etc. Do not pass the raw heap slot as a plain integer; use the handle returned by the constructor.

---

## Player-centric KCC getters (implicit subject)

After **`Player.Create(...)`** / **`Character.Create(...)`** / **`Char.Make(...)`**, the runtime remembers the **last KCC subject** (**implicit hero**). Most **`Player.Get*`** queries accept **either**:

- **`()`** — use the implicit subject (the capsule you created last in this session), or  
- **`(entity)`** — query a specific entity id / **EntityRef**.

If you call **`Player.GetPositionX()`** (or any zero-arg getter) **before** any KCC exists, the runtime reports an error (no implicit subject).

**Short names** (same handlers as the long forms; see [API_CONSISTENCY.md](../API_CONSISTENCY.md)):

| Short | Equivalent |
|-------|------------|
| **`Player.GetX()`** / **`GetY()`** / **`GetZ()`** | **`Player.GetPositionX()`** / **`GetPositionY()`** / **`GetPositionZ()`** |
| **`Player.GetPitch()`** / **`GetYaw()`** / **`GetRoll()`** | **`Player.GetRotationPitch()`** / **`GetRotationYaw()`** / **`GetRotationRoll()`** |
| **`Player.GetGrounded()`** | **`Player.IsGrounded()`** |
| **`Player.GetGravity()`** | **`Player.GetGravityScale()`** (per-character scale, not world gravity) |
| **`Player.GetCapsuleRadius()`** / **`GetCapsuleHeight()`** | **`Player.GetRadius()`** / **`Player.GetHeight()`** |
| **`Player.GetShapeType()`** | Returns **`"capsule"`** (CharacterVirtual) |

**World gravity** (global, not per-player): **`Physics.GetGravityX()`** / **`Physics.GetGravityY()`** / **`Physics.GetGravityZ()`** (and **`Physics3D.GetGravity*`** aliases) — see [PHYSICS3D.md](PHYSICS3D.md).

**Not exposed** as **`Ray.*`** / **`Sweep.*`** / **`Debug.*`** getters today — use **`PICK.*`**, **`PHYSICS3D`**, and engine logging as documented in [PHYSICS3D.md](PHYSICS3D.md). Physics-wide **body counts** / **collision counts** are not surfaced as **`Physics.GetBodyCount`**-style builtins yet.

---

### `PLAYER.CREATE(entity [, radius, height])`
Spawns a capsule character controller at the entity's world position.

- **Arguments**:
    - `entity`: (Handle) The entity to control.
    - `radius`: (Float, Optional) Capsule radius (default 0.4).
    - `height`: (Float, Optional) Total height (default 1.75).
- **Returns**: (Handle) The entity handle.
- **Example**:
    ```basic
    hero = ENTITY.LOAD("hero.iqm")
    PLAYER.CREATE(hero, 0.4, 1.8)
    ```

---

### `PLAYER.MOVE(entity, vx, vz)`
Moves the character with world-space horizontal velocity.

- **Arguments**:
    - `entity`: (Handle) The character entity.
    - `vx, vz`: (Float) Velocity in units per second.
- **Returns**: (Handle) The entity handle.

---

### `PLAYER.MOVEWITHCAMERA(entity, camera, fwd, side, speed)`
Moves the character relative to a camera's orientation.

- **Arguments**:
    - `entity`: (Handle) The character entity.
    - `camera`: (Handle) The reference camera.
    - `fwd`: (Float) Forward/Back input (-1 to 1).
    - `side`: (Float) Left/Right input (-1 to 1).
    - `speed`: (Float) Movement speed.
- **Returns**: (Handle) The entity handle.

---

### `PLAYER.JUMP(entity, impulseY)`
Applies an upward vertical impulse.

- **Returns**: (Handle) The entity handle.

---

### `PLAYER.ISGROUNDED(entity)`
Returns `TRUE` if the character reports ground support.

- **Returns**: (Boolean)

---

### `PLAYER.GETLOOKTARGET(entity, maxDist)`
Returns the entity ID being looked at by this character.

- **Arguments**:
    - `entity`: (Handle) The character entity.
    - `maxDist`: (Float) Maximum ray distance.
- **Returns**: (Integer) Entity ID, or 0 if nothing hit.

---

### `PLAYER.GETNEARBY(entity, radius, tag)`
Returns a list of entities within range matching a tag.

- **Arguments**:
    - `entity`: (Handle) The source entity.
    - `radius`: (Float) Search radius.
    - `tag`: (String) Case-insensitive glob (e.g., "Enemy*").
- **Returns**: (Handle) A numeric array handle of entity IDs.

---

### `PLAYER.TELEPORT(entity, x, y, z)`
Instantly snaps the character to a new position and clears velocity.

- **Returns**: (Handle) The entity handle.

### Tuning, queries, and gameplay helpers

| Command | Description |
|--------|-------------|
| **`PLAYER.SETGRAVITYSCALE(entity, scale)`** | Scales **gravity on Y** during **`CharacterMoveXZVelocity`** (**1** = default; values below **1** lighten gravity; above **1** strengthen it). |
| **`PLAYER.GETCROUCH(entity)`** / **`PLAYER.SETCROUCH(entity, bool)`** | Stored **crouch** flag for gameplay. **Capsule height** is not changed yet (Jolt wrapper limitation). |
| **`PLAYER.SWIM(entity, buoyancy, drag)`** | **Swim mode**: **buoyancy** (0–1) reduces downward gravity; **drag** damps horizontal velocity per second. Use **`(0, 0)`** to disable. |
| **`PLAYER.SETSTEPOFFSET(entity, height)`** | Alias of **`PLAYER.SETSTEPHEIGHT`** (reserved for future stair tuning). |
| **`PLAYER.GETSTANDNORMAL(entity)`** → **vec3 handle** | Ground/floor normal under the feet (**`GetGroundNormal`** or short downward ray). |
| **`PLAYER.PUSH(player, target, force)`** | Forward **horizontal** push on **target** via **`ENTITY.ADDFORCE`**-style integration; scaled by **`PLAYER.SETMASS`**. |
| **`PLAYER.GRAB(player, target)`** | Each **`PLAYER.MOVE`**, repositions **target** in front of the player (**`target 0`** releases). Not a Jolt **fixed constraint** yet. |
| **`PLAYER.SETMASS(entity, mass)`** | Stores **gameplay mass** (e.g. **`PLAYER.PUSH`**); Jolt **CharacterVirtual** mass is fixed at **`PLAYER.CREATE`**. |
| **`PLAYER.GETSURFACETYPE(entity)`** → **string** | Downward **Jolt** ray → hit entity → **`SurfaceMaterialHint`** from glTF **`material` / `footstep`** metadata or **Blender tag**; else **`Default`**. |
| **`PLAYER.SETFOVKICK` / `PLAYER.GETFOVKICK`** | Stores **extra FOV degrees** per entity; each frame do **`Camera.SetFOV(cam, base + Player.GetFovKick(hero))`** (or your own base). |
| **`PLAYER.ISMOVING(entity)`** → **bool** | **True** if horizontal **linear speed** is above ~**0.05** (for footsteps / sprint FX). |
| **`PLAYER.GETPOSITIONX` / `Y` / `Z`**, **`GETROTATIONPITCH` / `YAW` / `ROLL`**, **`GETVELOCITYX` / `Y` / `Z`**, **`GETSPEED`** | **float** — world pose / velocity from **CharacterVirtual** + entity bridge. |
| **`PLAYER.GETONSLOPE`**, **`GETONWALL`**, **`GETSLOPEANGLE`**, **`GETISJUMPING`**, **`GETISFALLING`** | **bool** / **float** — ground and motion hints (**`GETONSLOPE`** mirrors **`ISONSTEEPSLOPE`**; **`GETONWALL`** uses Jolt **NotSupported**). |
| **`PLAYER.GETMAXSLOPE`**, **`GETSTEPHEIGHT`**, **`GETGRAVITYSCALE`**, **`GETFRICTION`**, **`GETSNAPDISTANCE`**, **`GETHEIGHT`**, **`GETRADIUS`** | **float** — tuned capsule / stair / gravity / stick-down. |
| **`PLAYER.GETLAYER`**, **`GETMASK`** | **int** — reserved (**0**). |
| **`PLAYER.GETCOLLISIONENABLED`** | **bool** — reserved (**true**). |
| **`CHAR.GET*`** | Same signatures as **`PLAYER.GET*`** (aliases). |
| **`PLAYER.SnapToGround(entity, terrain, offset)`** | Sets **Y** from **`Terrain.GetHeight`** at the entity’s **XZ** plus **offset** (feet vs pivot). On **Linux + Jolt** after **`PLAYER.CREATE`**, also syncs the **CharacterVirtual** capsule. |
| **`PLAYER.ISSWIMMING(entity)`** → **bool** | **True** when the entity’s position lies inside a **`WATER`** volume column (between **bed** and the wavy surface). Use with **`PLAYER.SETGRAVITYSCALE`** for floatier movement. |

---

## Velocity & impulse

| Command | Description |
|--------|-------------|
| **`PLAYER.SETVELOCITY(entity, vx, vy, vz)`** | Override the character's full velocity vector directly. |
| **`PLAYER.ADDIMPULSE(entity, ix, iy, iz)`** | Add an instant impulse (world units/s) to the character velocity. |
| **`PLAYER.GETVELOCITYX(entity)`** / **`GETVELOCITYY`** / **`GETVELOCITYZ`** | Per-axis world velocity. Aliases: **`GETVX`** / **`GETVY`** / **`GETVZ`**. |
| **`PLAYER.GETGROUNDVELOCITYX(entity)`** / **`GETGROUNDVELOCITYY(entity)`** / **`GETGROUNDVELOCITYZ(entity)`** | Moving-platform velocity projected onto the ground plane per axis. |
| **`PLAYER.GETSUBMERGEDFACTOR(entity)`** | **0.0** = fully above water; **1.0** = fully submerged (for swimming blend). |

---

## Ground & ceiling state

| Command | Description |
|--------|-------------|
| **`PLAYER.GETGROUNDSTATE(entity)`** → **int** | Jolt `EGroundState`: 0=OnGround, 1=OnSteepGround, 2=NotSupported, 3=InAir. |
| **`PLAYER.GETCEILING(entity)`** → **bool** | **True** when head contact is detected (ceiling collision). |
| **`PLAYER.GETISSLIDING(entity)`** → **bool** | **True** when the character is sliding down a steep slope (`OnSteepGround`). |
| **`PLAYER.ISSUBMERGED(entity)`** → **bool** | **True** when entity center is below water surface. |

---

## Tuning setters

| Command | Description |
|--------|-------------|
| **`PLAYER.SETAIRCONTROL(entity, factor)`** | Scale (0–1) for XZ control while airborne (default 0.3). |
| **`PLAYER.SETGROUNDCONTROL(entity, factor)`** | Scale for XZ movement response on ground (default 1.0). |
| **`PLAYER.SETJUMPBUFFER(entity, seconds)`** | Coyote-time window after leaving ground where jump is still accepted. |
| **`PLAYER.SETPADDING(entity, dist)`** | Skin-width padding around the capsule for penetration recovery. |
| **`PLAYER.SETSLOPELIMIT(entity, degrees)`** | Maximum walkable slope angle; steeper ground = `OnSteepGround`. |
| **`PLAYER.SETSTICKFLOOR(entity, dist)`** | Downward snap distance to stick to floors on ramps and steps. |

---

## Events & movement helpers

| Command | Description |
|--------|-------------|
| **`PLAYER.ONTRIGGER(entity, callback)`** | Register a callback fired when the character overlaps a trigger volume. |
| **`PLAYER.MOVERELATIVE(camYaw, fwd, strafe, speed, dt)`** | Returns `[dx, dz]` camera-relative delta — same as `MOVESTEPX` + `MOVESTEPZ` combined. |

---

## Full Example

```moonbasic
WINDOW.OPEN(1280, 720, "Player")
WINDOW.SETFPS(60)
PHYSICS3D.START()
hero = ENTITY.LOAD("hero.iqm")
PLAYER.CREATE(hero)
cam = CAMERA.CREATE()

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    PLAYER.MOVE(hero, INPUT.AXIS(KEY_S, KEY_W) * 5.0, INPUT.AXIS(KEY_A, KEY_D) * 5.0)
    IF PLAYER.ISGROUNDED(hero) AND INPUT.KEYPRESSED(KEY_SPACE) THEN
        PLAYER.JUMP(hero, 6.0)
    ENDIF
    target = PLAYER.GETLOOKTARGET(hero, 3.0)
    IF target <> 0 AND INPUT.KEYPRESSED(KEY_E) THEN
        fn = LEVEL.MATCHSCRIPTBIND(EntityName(target))
        REM dispatch fn in your script...
    ENDIF
    PLAYER.SYNCANIM(hero, 0.12)
    RENDER.CLEAR(20, 24, 32)
    RENDER.Begin3D(cam)
        ENTITY.DRAW(hero)
    RENDER.END3D()
    RENDER.FRAME()
WEND
```

Naming: use **`LEVEL.LOAD`** / **`ENTITY.DRAW`** (or your project’s draw path), not **`SCENE.DRAW`**, so **`SCENE.*`** stays reserved for **mbscene** game scenes.

---

## See also

- [CHARACTER.md](CHARACTER.md) — **`CHARACTER.CREATE(entity, r, h)`**, **`CHARACTERREF.*`**, entity-bound Jolt KCC  
- [CHARACTER_PHYSICS.md](CHARACTER_PHYSICS.md) — character controller tutorial and full sample
- [LEVEL.md](LEVEL.md) — glTF, tags, **`LEVEL.BINDSCRIPT`**  
- [PHYSICS3D.md](PHYSICS3D.md) — Jolt world, **`PICK.*`**, rays  
- [API_CONSISTENCY.md](../API_CONSISTENCY.md) — machine-generated list of every **`PLAYER.*`** / **`CHAR.*`** name and arity (from `commands.json`; use after manifest changes)  
