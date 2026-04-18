# Mover Commands

Platform-style character movement helpers for entities: step-based XZ movement with slope/ceiling detection, landing, and relative movement. Simpler than `CHARCONTROLLER` — suited to grid or arcade movement.

## Core Workflow

1. `MOVER.MOVEXZ(entity, dx, dz, speed, gravity, stepHeight)` — attempt XZ movement with stepping. Returns the mover handle.
2. `MOVER.LAND(handle)` — snap the entity to the ground.
3. `MOVER.MOVEREL(handle, dx, dy, dz)` — apply a relative delta.
4. Read step results via `MOVER.MOVESTEPX` / `MOVER.MOVESTEPZ`.
5. `MOVER.FREE(handle)` when done.

---

## Movement

### `MOVER.MOVEXZ(entity, dx, dz, speed, gravity, stepHeight)` 

Attempts to move `entity` by `(dx, dz)` this frame at `speed`, applying `gravity`, and stepping up ledges up to `stepHeight`. Returns a **mover handle** for reading results.

---

### `MOVER.MOVESTEPX(entity, dx, dz, speed, gravity, stepHeight)` 

Returns the actual X displacement achieved (after collision resolution). Same arguments as `MOVEXZ`.

---

### `MOVER.MOVESTEPZ(entity, dx, dz, speed, gravity, stepHeight)` 

Returns the actual Z displacement achieved.

---

### `MOVER.MOVEREL(handle, dx, dy, dz)` 

Applies a relative displacement to the entity attached to `handle`.

---

### `MOVER.LAND(handle)` 

Snaps the entity to the ground surface below it (ground probe). Stops vertical velocity.

---

## Lifetime

### `MOVER.FREE(handle)` 

Frees the mover handle.

---

## Full Example

Arcade-style grid character with step-up on slopes.

```basic
WINDOW.OPEN(960, 540, "Mover Demo")
WINDOW.SETFPS(60)

cam  = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 8, -12)
CAMERA.SETTARGET(cam, 0, 0, 0)

player = ENTITY.CREATECUBE(1.0)
ENTITY.SETPOS(player, 0, 1, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()

    dx = 0.0
    dz = 0.0
    IF INPUT.KEYDOWN(KEY_RIGHT) THEN dx =  1.0
    IF INPUT.KEYDOWN(KEY_LEFT)  THEN dx = -1.0
    IF INPUT.KEYDOWN(KEY_UP)    THEN dz = -1.0
    IF INPUT.KEYDOWN(KEY_DOWN)  THEN dz =  1.0

    mv = MOVER.MOVEXZ(player, dx, dz, 5.0, -9.8, 0.4)
    MOVER.LAND(mv)
    MOVER.FREE(mv)

    ENTITY.UPDATE(dt)

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW3D.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

ENTITY.FREE(player)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [CHARCONTROLLER.md](CHARCONTROLLER.md) — full capsule character controller
- [CHARACTER.md](CHARACTER.md) — entity-bound character with physics
- [PLAYER.md](PLAYER.md) — high-level player controller
- [ENTITY.md](ENTITY.md) — entity system
