# Action & Gamepad Commands

Abstract input mapping that decouples game logic from hardware buttons.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Map logical actions to physical inputs with `ACTION.MAPKEY`, `ACTION.MAPJOY`, `ACTION.MAPMOUSE`, or `ACTION.MAPAXIS`.
2. Query action state each frame with `ACTION.DOWN`, `ACTION.PRESSED`, `ACTION.RELEASED`, or `ACTION.VALUE`.
3. Reset all mappings with `ACTION.RESET` when switching input profiles.

For raw gamepad access without the action layer, use `GAMEPAD.AXIS` and `GAMEPAD.BUTTON` directly or `AXIS.DPADY` for D-pad vertical input.

---

### `ACTION.MAPKEY(name, key)` / `MAPJOY` / `MAPMOUSE` / `MAPAXIS`
Binds a physical input to a named logical action.

- **Arguments**:
    - `name`: (String) Logical action name (e.g., "jump").
    - `key / button / axis`: (Integer) Hardware input constant.
- **Returns**: (None)
- **Example**:
    ```basic
    ACTION.MAPKEY("jump", KEY_SPACE)
    ACTION.MAPJOY("jump", 0, 0)
    ```

---

### `ACTION.DOWN(name)` / `PRESSED` / `RELEASED`
Polls the binary state of a named action.

- **Returns**: (Boolean)

---

### `ACTION.VALUE(name)`
Returns the analog value of an action (0.0 to 1.0).

- **Returns**: (Float)

---

### `ACTION.RESET()`
Clears all input mappings for all actions.

---

### `GAMEPAD.AXIS(idx, axis)`
Returns the raw float value of a gamepad axis (-1.0 to 1.0).

- **Returns**: (Float)

---

### `GAMEPAD.BUTTON(idx, btn)`
Returns `TRUE` if the gamepad button is currently held.

- **Returns**: (Boolean)

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
