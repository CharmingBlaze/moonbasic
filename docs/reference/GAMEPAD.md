# Gamepad Commands

Direct gamepad axis and button polling by device ID. For named action mapping see [ACTION.md](ACTION.md).

## Core Workflow

1. Check `GAME.ISGAMEPADAVAILABLE(id)` before reading.
2. `GAMEPAD.AXIS(gamepadId, axisId)` — returns -1.0 to 1.0.
3. `GAMEPAD.BUTTON(gamepadId, button)` — returns `TRUE` while held.

Axis constants: `GAMEPAD_AXIS_LEFT_X` = 0, `GAMEPAD_AXIS_LEFT_Y` = 1, `GAMEPAD_AXIS_RIGHT_X` = 2, `GAMEPAD_AXIS_RIGHT_Y` = 3, `GAMEPAD_AXIS_LEFT_TRIGGER` = 4, `GAMEPAD_AXIS_RIGHT_TRIGGER` = 5.

---

## Commands

### `GAMEPAD.AXIS(gamepadId, axisId)` 

Returns the current float value of `axisId` on gamepad `gamepadId`. Range: `-1.0` to `1.0` (triggers: `0.0` to `1.0`).

---

### `GAMEPAD.BUTTON(gamepadId, button)` 

Returns `TRUE` while `button` is held on `gamepadId`.

---

## Full Example

Two-stick movement and camera.

```basic
WINDOW.OPEN(960, 540, "Gamepad Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 6, -10)
CAMERA.SETTARGET(cam, 0, 0, 0)

px = 0.0  pz = 0.0
yaw = 0.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()

    IF GAME.ISGAMEPADAVAILABLE(0) THEN
        lx = GAMEPAD.AXIS(0, GAMEPAD_AXIS_LEFT_X)
        lz = GAMEPAD.AXIS(0, GAMEPAD_AXIS_LEFT_Y)
        rx = GAMEPAD.AXIS(0, GAMEPAD_AXIS_RIGHT_X)
        px = px + lx * 5 * dt
        pz = pz + lz * 5 * dt
        yaw = yaw + rx * 90 * dt
    END IF

    CAMERA.SETORBIT(cam, px, 0, pz, yaw, 25, 12)

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        DRAW.SPHERE(px, 0.5, pz, 0.5, 80, 200, 255, 255)
        DRAW.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [ACTION.md](ACTION.md) — named action bindings (preferred)
- [INPUT.md](INPUT.md) — unified input including gamepad
- [GAME.md](GAME.md) — `GAME.JOYX/Y`, `GAME.JOYBUTTON`
