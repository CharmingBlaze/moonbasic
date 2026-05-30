# Controller Commands

Lightweight capsule character controller (non-Jolt). Simpler alternative to `CHARCONTROLLER` for basic grounded movement. Suitable for 2.5D and simple 3D games that don't need the full Jolt KCC.

For the full Jolt-backed controller see [CHARCONTROLLER.md](CHARCONTROLLER.md).

## Core Workflow

1. `CONTROLLER.CREATE(x, y, z, radius, height)` — create a capsule at position.
2. `CONTROLLER.MOVE(ctrl, dx, dy, dz)` each frame.
3. `CONTROLLER.GROUNDED(ctrl)` — check floor contact.
4. `CONTROLLER.JUMP(ctrl, impulse)` — jump.
5. `CONTROLLER.FREE(ctrl)` when done.

---

## Commands

### `CONTROLLER.CREATE(x, y, z, radius, height)` 

Creates a capsule controller at world position `(x, y, z)` with `radius` and `height`. Returns a **controller handle**.

---

### `CONTROLLER.MOVE(ctrl, dx, dy, dz)` 

Moves the controller by world delta `(dx, dy, dz)` this frame. Resolves against environment collision.

---

### `CONTROLLER.GROUNDED(ctrl)` 

Returns `TRUE` if the controller is resting on a surface.

---

### `CONTROLLER.JUMP(ctrl, impulse)` 

Applies an upward jump impulse.

---

### `CONTROLLER.FREE(ctrl)` 

Frees the controller handle.

---

## Full Example

Basic 3D character movement.

```basic
WINDOW.OPEN(960, 540, "Controller Demo")
WINDOW.SETFPS(60)

PHYSICS3D.START()
PHYSICS3D.SETGRAVITY(0, -20, 0)

cam  = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 6, -10)
CAMERA.SETTARGET(cam, 0, 0, 0)

ctrl = CONTROLLER.CREATE(0, 2, 0, 0.4, 1.8)
vy   = 0.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    dx = 0.0 : dz = 0.0

    IF INPUT.KEYDOWN(KEY_D) THEN dx =  5 * dt
    IF INPUT.KEYDOWN(KEY_A) THEN dx = -5 * dt
    IF INPUT.KEYDOWN(KEY_W) THEN dz = -5 * dt
    IF INPUT.KEYDOWN(KEY_S) THEN dz =  5 * dt

    IF CONTROLLER.GROUNDED(ctrl) THEN
        vy = 0
        IF INPUT.KEYPRESSED(KEY_SPACE) THEN vy = 8
    ELSE
        vy = vy - 20 * dt
    END IF

    CONTROLLER.MOVE(ctrl, dx, vy * dt, dz)
    PHYSICS3D.UPDATE()

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        DRAW.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

CONTROLLER.FREE(ctrl)
PHYSICS3D.STOP()
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `CONTROLLER.MAKE(entity)` | Deprecated alias of `CONTROLLER.CREATE`. |

---

## See also

- [CHARCONTROLLER.md](CHARCONTROLLER.md) — Jolt-backed full KCC
- [CHARACTERREF.md](CHARACTERREF.md) — handle-method Jolt KCC
- [PHYSICS3D.md](PHYSICS3D.md) — world setup
