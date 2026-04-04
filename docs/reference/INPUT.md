# Input Commands

Commands for handling keyboard, mouse, and gamepad input.

---

## Direct Input

These commands directly query the state of a specific key or button.

### `Input.KeyDown(key)`

Returns `TRUE` if the specified key is currently being held down.

- `key`: The key code to check (e.g., `KEY_SPACE`, `KEY_W`).

### `Input.KeyPressed(key)`

Returns `TRUE` only on the frame the specified key was first pressed. Useful for single-trigger events like jumping.

### `Input.KeyUp(key)`

Returns `TRUE` only on the frame the specified key was released.

---

### `Input.MouseDown(button)`

Returns `TRUE` if the specified mouse button is currently being held down.

- `button`: The mouse button to check (e.g., `MOUSE_LEFT_BUTTON`).

### `Input.MouseX()` / `Input.MouseY()`

Returns the current X or Y coordinate of the mouse cursor.

---

## Action Mapping

The action mapping system is a powerful way to handle input. Instead of checking for specific keys, you define abstract "actions" and then check the state of those actions. This makes it easy to support multiple input devices and allow for user-configurable controls.

### 1. Define Mappings

First, map physical inputs to action names. This is typically done once at the start of your program.

- `Input.MapKey(action$, key)`: Maps a keyboard key to an action.
- `Input.MapGamepadButton(action$, gamepad, button)`: Maps a gamepad button.
- `Input.MapGamepadAxis(action$, gamepad, axis, direction)`: Maps a gamepad axis direction (e.g., left stick up) to an action.

```basic
; Define a "move_right" action for multiple inputs
Input.MapKey("move_right", KEY_D)
Input.MapKey("move_right", KEY_RIGHT) ; Also map the right arrow key
Input.MapGamepadButton("move_right", 0, GAMEPAD_BUTTON_RIGHT_FACE_RIGHT) ; D-pad right
```

### 2. Check Actions in the Game Loop

In your main loop, check the state of the action by its name, not the key.

- `Input.ActionDown(action$)`: Returns `TRUE` if any mapped input is currently active (held down).
- `Input.ActionPressed(action$)`: Returns `TRUE` only on the frame the action was first triggered.
- `Input.ActionReleased(action$)`: Returns `TRUE` on the frame the action was released.

```basic
; In the main loop, check the abstract action
IF Input.ActionDown("move_right") THEN
    player_x# = player_x# + speed# * Time.Delta()
ENDIF
```

### `Input.ActionAxis(action$)`

For analog controls like a joystick, this returns the axis value from -1.0 to 1.0.

```basic
; Map joystick axes
Input.MapGamepadAxis("move_horizontal", 0, GAMEPAD_AXIS_LEFT_X)

; Get the analog value
move_x# = Input.ActionAxis("move_horizontal")
player_x# = player_x# + move_x# * speed# * Time.Delta()
```
