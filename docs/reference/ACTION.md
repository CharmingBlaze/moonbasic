# Action Commands

Named input action mapping system. Bind named actions to keyboard keys, mouse buttons, gamepad axes, and joystick buttons. Query actions by name instead of hardcoded key codes — makes remapping trivial.

## Core Workflow

1. `ACTION.MAPKEY(name, keyCode)` — bind an action name to a key.
2. `ACTION.MAPMOUSE(name, button)` / `ACTION.MAPJOY(name, pad, button)` — bind mouse or gamepad.
3. Each frame: `ACTION.DOWN(name)`, `ACTION.PRESSED(name)`, `ACTION.RELEASED(name)`, `ACTION.VALUE(name)`.
4. `ACTION.RESET()` to clear all bindings.

---

## Binding

### `ACTION.MAPKEY(actionName, keyCode)` 

Binds `actionName` to a keyboard key constant (e.g. `KEY_SPACE`, `KEY_W`).

---

### `ACTION.MAPMOUSE(actionName, button)` 

Binds `actionName` to a mouse button: `0` = left, `1` = right, `2` = middle.

---

### `ACTION.MAPJOY(actionName, gamepadId, button)` 

Binds `actionName` to a gamepad button on device `gamepadId`.

---

### `ACTION.MAPAXIS(actionName, gamepadId, axis)` 

Binds `actionName` to a gamepad axis. `ACTION.VALUE` returns the axis float.

---

## Querying

### `ACTION.DOWN(actionName)` 

Returns `TRUE` while the bound input is held.

---

### `ACTION.PRESSED(actionName)` 

Returns `TRUE` on the first frame the bound input is pressed.

---

### `ACTION.RELEASED(actionName)` 

Returns `TRUE` on the first frame the bound input is released.

---

### `ACTION.VALUE(actionName)` 

Returns the float value of the bound input — `0.0`/`1.0` for buttons, `-1.0` to `1.0` for axes.

---

## Reset

### `ACTION.RESET()` 

Clears all action bindings.

---

## Full Example

Rebindable controls using ACTION.

```basic
WINDOW.OPEN(800, 450, "Action Demo")
WINDOW.SETFPS(60)

; default bindings
ACTION.MAPKEY("jump",    KEY_SPACE)
ACTION.MAPKEY("moveR",   KEY_D)
ACTION.MAPKEY("moveL",   KEY_A)
ACTION.MAPJOY("jump",    0, GAMEPAD_BUTTON_RIGHT_FACE_DOWN)
ACTION.MAPAXIS("moveR",  0, GAMEPAD_AXIS_LEFT_X)

px = 400.0
py = 300.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()

    px = px + ACTION.VALUE("moveR") * 200 * dt
    IF ACTION.DOWN("moveL") THEN px = px - 200 * dt
    IF ACTION.PRESSED("jump") THEN py = py - 120

    RENDER.CLEAR(20, 20, 40)
    DRAW.CIRCLE(INT(px), INT(py), 16, 80, 200, 255, 255)
    DRAW.TEXT("SPACE/GamepadA = jump | A/D = move", 10, 10, 18, 200, 200, 200, 255)
    RENDER.FRAME()
WEND

ACTION.RESET()
WINDOW.CLOSE()
```

---

## See also

- [INPUT.md](INPUT.md) — raw key/mouse/gamepad polling
- [GAMEPAD.md](GAMEPAD.md) — `GAMEPAD.AXIS` / `GAMEPAD.BUTTON`
- [MOUSE.md](MOUSE.md) — mouse buttons
