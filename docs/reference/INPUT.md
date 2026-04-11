# Input Commands

Commands for handling keyboard, mouse, and gamepad input.

---

### `Input.KeyDown(key)`
Returns `TRUE` if the key is currently held down.

### `Input.KeyPressed(key)`
Returns `TRUE` on the frame the key was first pressed.

### `Input.KeyUp(key)`
Returns `TRUE` on the frame the key was released.

---

### `Input.MouseX()` / `Input.MouseY()`
Returns the current mouse pixel coordinates.

### `Input.MouseButtonDown(button)`
Returns `TRUE` if the mouse button is held down.

### `Input.MouseWheelMove()`
Returns the mouse wheel movement delta.

---

### `Input.Axis(negKey, posKey)`
Returns a value from -1.0 to 1.0 based on two keys.

### `Input.Movement2D(up, down, left, right)`
Returns a 2-element array handle `[forward, strafe]` based on four keys.

### `Input.MouseDelta()`
Returns a 2-element array handle `[dx, dy]` representing mouse movement since last frame.

---

## Keyboard Constants

### Letters
`KEY_A` to `KEY_Z`

### Numbers
`KEY_ZERO` to `KEY_NINE`

### Function Keys
`KEY_F1` to `KEY_F12`

### Arrow Keys
`KEY_UP`, `KEY_DOWN`, `KEY_LEFT`, `KEY_RIGHT`

### Special Keys
`KEY_SPACE`, `KEY_ESCAPE`, `KEY_ENTER`, `KEY_TAB`, `KEY_BACKSPACE`
`KEY_LEFT_SHIFT`, `KEY_LEFT_CONTROL`, `KEY_LEFT_ALT`
`KEY_RIGHT_SHIFT`, `KEY_RIGHT_CONTROL`, `KEY_RIGHT_ALT`

---

## Mouse Constants
`MOUSE_BUTTON_LEFT`, `MOUSE_RIGHT_BUTTON`, `MOUSE_MIDDLE_BUTTON`

---

## Cursor

### `Cursor.Hide()` / `Cursor.Show()`
Hides or shows the OS mouse cursor while over the window.

### `Cursor.Disable()` / `Cursor.Enable()`
Disables the cursor and switches to relative mouse mode (centered virtual cursor).

---

## Action Mapping

The action mapping system allows mapping physical inputs to abstract action names.

### `Input.MapKey(action, key)`
Maps a keyboard key to an action.

### `Input.ActionDown(action)`
Returns `TRUE` if any mapped input for the action is held down.

### `Input.ActionPressed(action)`
Returns `TRUE` only on the frame the action was first triggered.

### `Input.ActionAxis(action)`
Returns the analog axis value (-1.0 to 1.0) for the action.

---

## Keyboard Constants

### Letters
`KEY_A` to `KEY_Z`

### Numbers
`KEY_ZERO` to `KEY_NINE`

### Function Keys
`KEY_F1` to `KEY_F12`

### Arrow Keys
`KEY_UP`, `KEY_DOWN`, `KEY_LEFT`, `KEY_RIGHT`

### Special Keys
`KEY_SPACE`, `KEY_ESCAPE`, `KEY_ENTER`, `KEY_TAB`, `KEY_BACKSPACE`
`KEY_LEFT_SHIFT`, `KEY_LEFT_CONTROL`, `KEY_LEFT_ALT`
`KEY_RIGHT_SHIFT`, `KEY_RIGHT_CONTROL`, `KEY_RIGHT_ALT`

---

## Mouse Constants
`MOUSE_BUTTON_LEFT`, `MOUSE_RIGHT_BUTTON`, `MOUSE_MIDDLE_BUTTON`
