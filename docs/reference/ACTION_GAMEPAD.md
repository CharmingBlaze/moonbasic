# Action & Gamepad Commands

Abstract input mapping that decouples game logic from hardware buttons.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Map logical actions to physical inputs with `ACTION.MAPKEY`, `ACTION.MAPJOY`, `ACTION.MAPMOUSE`, or `ACTION.MAPAXIS`.
2. Query action state each frame with `ACTION.DOWN`, `ACTION.PRESSED`, `ACTION.RELEASED`, or `ACTION.VALUE`.
3. Reset all mappings with `ACTION.RESET` when switching input profiles.

For raw gamepad access without the action layer, use `GAMEPAD.AXIS` and `GAMEPAD.BUTTON` directly or `AXIS.DPADY` for D-pad vertical input.

---

### `ACTION.MAPKEY(actionName, keyCode)` 

Binds a keyboard key to a named action.

- `actionName`: Logical action name (e.g. `"jump"`).
- `keyCode`: Keyboard key constant.

---

### `ACTION.MAPJOY(actionName, gamepadIndex, buttonIndex)` 

Binds a gamepad button to a named action.

- `actionName`: Logical action name.
- `gamepadIndex`: Gamepad slot (0-based).
- `buttonIndex`: Button index on the gamepad.

---

### `ACTION.MAPMOUSE(actionName, mouseButton)` 

Binds a mouse button to a named action.

- `actionName`: Logical action name.
- `mouseButton`: Mouse button index (0 = left, 1 = right, 2 = middle).

---

### `ACTION.MAPAXIS(actionName, gamepadIndex, axisIndex)` 

Binds a gamepad axis to a named action for analog queries via `ACTION.VALUE`.

- `actionName`: Logical action name.
- `gamepadIndex`: Gamepad slot (0-based).
- `axisIndex`: Axis index on the gamepad.

---

### `ACTION.DOWN(actionName)` 

Returns `TRUE` while the action is held down (any mapped input is active).

---

### `ACTION.PRESSED(actionName)` 

Returns `TRUE` only on the frame the action was first pressed.

---

### `ACTION.RELEASED(actionName)` 

Returns `TRUE` only on the frame the action was released.

---

### `ACTION.VALUE(actionName)` 

Returns the analog value of the action (0.0–1.0 for axes, 0/1 for digital).

---

### `ACTION.RESET()` 

Clears all action-to-input mappings.

---

### `GAMEPAD.AXIS(gamepadIndex, axisIndex)` 

Returns the raw float value of a gamepad axis (−1.0 to 1.0).

- `gamepadIndex`: Gamepad slot (0-based).
- `axisIndex`: Axis index.

---

### `GAMEPAD.BUTTON(gamepadIndex, buttonIndex)` 

Returns `TRUE` if the specified gamepad button is currently held.

- `gamepadIndex`: Gamepad slot (0-based).
- `buttonIndex`: Button index.

---

### `AXIS.DPADY(gamepadIndex)` 

Returns the vertical D-pad value for the given gamepad as a float (−1.0 up, 1.0 down, 0.0 neutral).

---

## Full Example

This example sets up jump and move actions, then polls them in a game loop.

```basic
; Map "jump" to Space key and gamepad button 0
ACTION.MAPKEY("jump", KEY_SPACE)
ACTION.MAPJOY("jump", 0, 0)

; Map "move_x" to left stick X axis
ACTION.MAPAXIS("move_x", 0, 0)

; Game loop
WHILE NOT WINDOW.SHOULDCLOSE()
    IF ACTION.PRESSED("jump")
        PRINT "Jump!"
    END IF

    move = ACTION.VALUE("move_x")
    IF ABS(move) > 0.1
        PRINT "Moving: " + STR(move)
    END IF

    ; Raw gamepad check
    IF GAMEPAD.BUTTON(0, 1)
        PRINT "Gamepad button 1 held"
    END IF

    RENDER.BEGINFRAME()
    RENDER.ENDFRAME()
WEND

ACTION.RESET()
```
