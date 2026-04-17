# Movement Commands

Interpolation-based movement and camera-relative XZ stepping helpers.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

Use `MOVE.*` for simple value interpolation each frame, and `MOVER.*` for entity-level XZ movement relative to a camera yaw.

Combine with `DELTATIME()` or `DT` for frame-rate independent motion. For arcade collision math see [COLLISION.md](COLLISION.md); for input axes see [INPUT.md](INPUT.md).

---

### `MOVE.TOWARD(current, target, maxDelta)` 

Moves `current` toward `target` by at most `maxDelta`. Returns the new value. Clamps so it never overshoots.

---

### `MOVE.LERP(a, b, t)` 

Linear interpolation between `a` and `b` by factor `t` (0.0–1.0).

---

### `MOVER.MOVEXZ(entityHandle, yaw, forward, strafe, speed, dt)` 

Moves an entity on the XZ plane relative to camera `yaw`. Returns the mover handle.

- `forward`, `strafe`: Input axes (−1 to 1).
- `speed`: Movement speed.
- `dt`: Delta time.

---

### `MOVER.MOVESTEPX(entityHandle, yaw, forward, strafe, speed, dt)` 

Returns the X component of a camera-relative XZ step without applying it.

---

### `MOVER.MOVESTEPZ(entityHandle, yaw, forward, strafe, speed, dt)` 

Returns the Z component of a camera-relative XZ step without applying it.

---

### `MOVER.MOVEREL(entityHandle, dx, dy, dz)` 

Moves an entity by a relative offset in world space.

---

### `MOVER.LAND(entityHandle)` 

Snaps the entity to the ground (Y = terrain or collision surface).

---

### `MOVER.FREE(entityHandle)` 

Frees the mover state associated with the entity.

---

## Full Example

This example moves a character with camera-relative controls.

```basic
player = ENTITY.LOAD("player.glb")
cam = CAMERA.CREATE()
speed = 5.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = DELTATIME()
    yaw = CAMERA.GETYAW(cam)

    ; Get input
    fwd = 0.0
    stf = 0.0
    IF INPUT.KEYDOWN(KEY_W) THEN fwd = 1.0
    IF INPUT.KEYDOWN(KEY_S) THEN fwd = -1.0
    IF INPUT.KEYDOWN(KEY_A) THEN stf = -1.0
    IF INPUT.KEYDOWN(KEY_D) THEN stf = 1.0

    MOVER.MOVEXZ(player, yaw, fwd, stf, speed, dt)
    MOVER.LAND(player)

    RENDER.BEGINFRAME()
    RENDER.BEGINMODE3D(cam)
    ENTITY.DRAW(player)
    RENDER.ENDMODE3D()
    RENDER.ENDFRAME()
WEND
```
