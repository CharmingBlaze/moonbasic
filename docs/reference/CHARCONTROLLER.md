# CharController Commands

Jolt `CharacterVirtual`-backed capsule character controller (integer entity-id API). Create one per player entity, then drive it each frame with move/jump, and poll ground state.

Requires **CGO + Jolt** (Windows/Linux fullruntime).

## Core Workflow

1. `CHARCONTROLLER.CREATE(entity, radius, height, stepHeight, slopeLimit, snapDist)` — attach a KCC.
2. Each frame:
   - `CHARCONTROLLER.MOVE(entity, dx, dz)` — set horizontal movement.
   - `CHARCONTROLLER.APPLYGRAVITY(entity, dt)` — apply gravity.
   - `CHARCONTROLLER.UPDATE(entity, dt)` — resolve.
3. Read state: `CHARCONTROLLER.GROUNDED(entity)`, `CHARCONTROLLER.X/Y/Z(entity)`.
4. `CHARCONTROLLER.FREE(entity)` when done.

---

## Creation

### `CHARCONTROLLER.CREATE(entity, radius, height, stepHeight, slopeLimit, snapDist)` 

Attaches a Jolt CharacterVirtual KCC to `entity`. `radius`/`height` define the capsule. `stepHeight` = stair climb height. `slopeLimit` = max walkable slope degrees. `snapDist` = floor snap distance.

- *Handle shortcut*: `CHARACTER.CREATE(entity, radius, height)`

---

## Movement

### `CHARCONTROLLER.MOVE(entity, dx, dz)` 

Sets the horizontal movement delta this frame in world units.

---

### `CHARCONTROLLER.SETVELOCITY(entity, vx, vy, vz)` 

Overrides velocity directly.

---

### `CHARCONTROLLER.ADDVELOCITY(entity, vx, vy, vz)` 

Adds to current velocity (use for jump / knockback).

---

### `CHARCONTROLLER.APPLYGRAVITY(entity, dt)` 

Applies per-entity gravity this step. Call before `UPDATE`.

---

### `CHARCONTROLLER.UPDATE(entity, dt)` 

Resolves the KCC for one frame — collision detection, slope handling, step-up. Call once per frame.

---

### `CHARCONTROLLER.JUMP(entity, impulse)` 

Applies an upward jump impulse.

---

## Position

### `CHARCONTROLLER.X(entity)` / `CHARCONTROLLER.Y(entity)` / `CHARCONTROLLER.Z(entity)` 

Returns the current world position component.

---

### `CHARCONTROLLER.SETPOS(entity, x, y, z)` 

Teleports the character.

---

## Ground State

### `CHARCONTROLLER.GROUNDED(entity)` 

Returns `TRUE` when standing on a supported surface.

---

### `CHARCONTROLLER.GETGROUNDSTATE(entity)` 

Returns Jolt `EGroundState` int: `0`=OnGround, `1`=OnSteepGround, `2`=NotSupported, `3`=InAir.

---

### `CHARCONTROLLER.GETGROUNDVELOCITY(entity)` 

Returns platform velocity `[vx, vy, vz]` array (for moving platforms).

---

## Velocity Queries

### `CHARCONTROLLER.GETVELOCITY(entity)` 

Returns `[vx, vy, vz]` velocity array.

---

### `CHARCONTROLLER.ISSLIDING(entity)` 

Returns `TRUE` if sliding on steep ground.

---

### `CHARCONTROLLER.ISCEILING(entity)` 

Returns `TRUE` if head contact detected.

---

## Configuration

### `CHARCONTROLLER.SETGRAVITY(entity, scale)` 

Sets per-entity gravity scale.

---

### `CHARCONTROLLER.SETMAXSLOPE(entity, degrees)` 

Sets max walkable slope.

---

### `CHARCONTROLLER.SETSTEPHEIGHT(entity, height)` 

Sets stair step-up height.

---

### `CHARCONTROLLER.SETSNAPDIST(entity, dist)` 

Sets floor snap distance.

---

## Aliases & additional queries

| Command | Description |
|--------|-------------|
| **`CHARCONTROLLER.ISGROUNDED(entity)`** | Alias of `CHARCONTROLLER.GROUNDED`. Returns `TRUE` when on a supported surface. |
| **`CHARCONTROLLER.MAKE(...)`** | Deprecated alias of `CHARCONTROLLER.CREATE`. Use `CREATE`. |
| **`CHARCONTROLLER.SETPOSITION(entity, x, y, z)`** | Alias of `CHARCONTROLLER.SETPOS`. Teleports the character. |
| **`CHARCONTROLLER.GETPOS(entity)`** | Returns `[x, y, z]` world position array. |
| **`CHARCONTROLLER.GETLINEARVEL(entity)`** | Returns world linear velocity `[vx, vy, vz]` from Jolt. |
| **`CHARCONTROLLER.GETGROUNDNORMAL(entity)`** | Returns ground contact normal `[nx, ny, nz]`; returns up vector when airborne. |

---

## Lifetime

### `CHARCONTROLLER.FREE(entity)` 

Destroys the KCC for the entity.

---

## Full Example

```basic
WINDOW.OPEN(960, 540, "CharController Demo")
WINDOW.SETFPS(60)

PHYSICS3D.START()
PHYSICS3D.SETGRAVITY(0, -20, 0)

cam    = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 6, -10)
CAMERA.SETTARGET(cam, 0, 0, 0)

player = ENTITY.CREATECAPSULE(0.4, 1.8)
ENTITY.SETPOS(player, 0, 2, 0)
CHARCONTROLLER.CREATE(player, 0.4, 1.8, 0.4, 45, 0.2)
CHARCONTROLLER.SETGRAVITY(player, 1.0)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    dx = 0.0 : dz = 0.0
    IF INPUT.KEYDOWN(KEY_D) THEN dx =  5 * dt
    IF INPUT.KEYDOWN(KEY_A) THEN dx = -5 * dt
    IF INPUT.KEYDOWN(KEY_W) THEN dz = -5 * dt
    IF INPUT.KEYDOWN(KEY_S) THEN dz =  5 * dt

    CHARCONTROLLER.MOVE(player, dx, dz)
    IF INPUT.KEYPRESSED(KEY_SPACE) AND CHARCONTROLLER.GROUNDED(player) THEN
        CHARCONTROLLER.JUMP(player, 8)
    END IF
    CHARCONTROLLER.APPLYGRAVITY(player, dt)
    CHARCONTROLLER.UPDATE(player, dt)

    PHYSICS3D.UPDATE()
    ENTITY.UPDATE(dt)

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

CHARCONTROLLER.FREE(player)
ENTITY.FREE(player)
PHYSICS3D.STOP()
WINDOW.CLOSE()
```

---

## See also

- [CHARACTERREF.md](CHARACTERREF.md) — handle-based KCC (same system, different API)
- [CHAR.md](CHAR.md) — `CHAR.*` alias namespace
- [PLAYER.md](PLAYER.md) — high-level player controller
- [PHYSICS3D.md](PHYSICS3D.md) — world setup
