# Input Commands

Commands for reading keyboard, mouse, gamepad, touch, and gesture input. moonBASIC provides three tiers of input API: **low-level polling** (`Input.*`), **facade objects** (`Key.*`, `Mouse.*`, `Gamepad.*`), and **action mapping** (`Action.*`) for abstract game controls.

## Core Concepts

- **Down** — the key/button is currently held this frame.
- **Pressed** — the key/button was first pressed this frame (edge-triggered, fires once).
- **Released / Up** — the key/button was released this frame (edge-triggered).
- **Key codes** — integer constants like `KEY_W`, `KEY_SPACE`, `KEY_ESCAPE`. See Raylib key constants.
- **Mouse buttons** — 0 = Left, 1 = Right, 2 = Middle.

---

## Keyboard

### `Input.KeyDown(keyCode)`

Returns `TRUE` while the key is held down. Fires every frame the key is pressed.

- `keyCode` (int) — Raylib key constant (e.g., `KEY_W`, `KEY_SPACE`).

**Returns:** `bool`

```basic
IF Input.KeyDown(KEY_W) THEN
    playerY = playerY - speed * dt
ENDIF
```

---

### `Input.KeyPressed(keyCode)`

Returns `TRUE` only on the frame the key was first pressed. Does not repeat.

- `keyCode` (int) — Key constant.

**Returns:** `bool`

```basic
IF Input.KeyPressed(KEY_SPACE) THEN
    playerVelY = jumpForce
ENDIF
```

---

### `Input.KeyUp(keyCode)`

Returns `TRUE` only on the frame the key was released.

- `keyCode` (int) — Key constant.

**Returns:** `bool`

---

### `Input.KeyHit(keyCode)`

Alias for `Input.KeyPressed`. Returns `TRUE` on the frame the key was first pressed.

---

### `Input.GetKeyName(keyCode)`

Returns the human-readable name of a key code (e.g., `"W"`, `"Space"`, `"Escape"`).

- `keyCode` (int) — Key constant.

**Returns:** `string`

---

### `Input.CharPressed()`

Returns the Unicode character code of the last character typed this frame. Returns 0 if no character was typed. Useful for text input fields.

**Returns:** `int`

---

### `Input.GetInactivity()`

Returns the number of seconds since the last keyboard or mouse input event.

**Returns:** `float`

---

## Mouse

### `Input.MouseX()` / `Mouse.X()`

Returns the current mouse X position in window pixels.

**Returns:** `int`

---

### `Input.MouseY()` / `Mouse.Y()`

Returns the current mouse Y position in window pixels.

**Returns:** `int`

---

### `Input.MouseDown(button)` / `Mouse.Down(button)`

Returns `TRUE` while a mouse button is held down.

- `button` (int) — 0 = Left, 1 = Right, 2 = Middle.

**Returns:** `bool`

```basic
IF Mouse.Down(0) THEN
    ; Left mouse button is held
    Draw.Circle(Mouse.X(), Mouse.Y(), 5, 255, 0, 0, 255)
ENDIF
```

---

### `Input.MousePressed(button)` / `Mouse.Pressed(button)`

Returns `TRUE` on the frame the mouse button was first pressed.

- `button` (int) — Mouse button index.

**Returns:** `bool`

---

### `Input.MouseReleased(button)` / `Mouse.Released(button)`

Returns `TRUE` on the frame the mouse button was released.

- `button` (int) — Mouse button index.

**Returns:** `bool`

---

### `Input.MouseHit(button)`

Returns `TRUE` on the frame the mouse button was clicked. Same as `MousePressed`.

---

### `Input.MouseDeltaX()` / `Mouse.DX()` / `MOUSEDX`

Returns the horizontal mouse movement in pixels since the last frame. Positive = moved right.

**Returns:** `float`

---

### `Input.MouseDeltaY()` / `Mouse.DY()` / `MOUSEDY`

Returns the vertical mouse movement in pixels since the last frame. Positive = moved down.

**Returns:** `float`

```basic
; FPS-style camera rotation
yaw = yaw + Mouse.DX() * sensitivity
pitch = pitch + Mouse.DY() * sensitivity
```

---

### `Input.MouseWheelMove()` / `Mouse.Wheel()` / `MOUSEWHEEL`

Returns the mouse wheel scroll amount this frame. Positive = scrolled up.

**Returns:** `float`

---

### `Input.MouseDelta()`

Returns both delta X and Y as a combined value. Use `MouseDeltaX`/`MouseDeltaY` for individual axes.

---

### `Input.SetMousePos(x, y)` / `Mouse.SetPos(x, y)`

Warps the mouse cursor to the given screen position.

- `x` (int) — X position.
- `y` (int) — Y position.

---

### `Input.SetMouseScale(scaleX, scaleY)`

Sets a scale factor for mouse coordinates. Useful for HiDPI displays or custom coordinate systems.

- `scaleX` (float) — Horizontal scale.
- `scaleY` (float) — Vertical scale.

---

### `Input.SetMouseOffset(offsetX, offsetY)`

Sets an offset for mouse coordinates.

- `offsetX` (int) — X offset.
- `offsetY` (int) — Y offset.

---

### `Input.GetMouseWorldPos()`

Returns the mouse position transformed to world coordinates (requires an active camera).

**Returns:** `float` (use with camera unproject for full 3D world picking)

---

### `Input.LockMouse(enabled)`

Locks or unlocks the mouse cursor. When locked, the cursor is hidden and confined to the window — ideal for FPS games.

- `enabled` (bool) — `TRUE` to lock.

```basic
; Lock mouse for FPS controls
Input.LockMouse(TRUE)
```

---

## Cursor Visibility

### `Mouse.Show()` / `Cursor.Show()`

Shows the mouse cursor.

---

### `Mouse.Hide()` / `Cursor.Hide()`

Hides the mouse cursor.

---

### `Cursor.IsHidden()`

Returns `TRUE` if the cursor is currently hidden.

**Returns:** `bool`

---

### `Cursor.IsOnScreen()`

Returns `TRUE` if the cursor is inside the window bounds.

**Returns:** `bool`

---

### `Mouse.Enable()` / `Cursor.Enable()`

Enables cursor input (re-enables after `Disable`).

---

### `Mouse.Disable()` / `Cursor.Disable()`

Disables cursor input entirely.

---

### `Cursor.IsEnabled()`

Returns `TRUE` if cursor input is enabled.

**Returns:** `bool`

---

### `Cursor.Set(cursorType)`

Changes the cursor icon.

- `cursorType` (int) — Cursor type constant (arrow, hand, crosshair, etc.).

---

## Movement & Axis Helpers

### `Input.Axis()`

Returns a value from -1.0 to 1.0 representing the current horizontal movement input. Reads WASD and arrow keys. Useful for platformer movement.

**Returns:** `float`

```basic
moveX = Input.Axis() * speed * dt
```

---

### `Input.Movement2D()`

Returns a normalized 2D movement vector from WASD/arrow keys as a combined value.

---

### `Input.MoveDir()`

Returns the movement direction angle in radians based on current WASD input.

**Returns:** `float`

---

### `Input.AxisDeg()` / `Input.Orbit()`

Returns the axis input as a degree angle. Useful for orbit camera rotation.

**Returns:** `float`

---

## Gamepad

### `Input.IsGamepadAvailable(gamepadId)`

Returns `TRUE` if a gamepad is connected at the given index.

- `gamepadId` (int) — Gamepad index (0-based).

**Returns:** `bool`

---

### `Input.GetGamepadAxisValue(gamepadId, axis)`

Returns the raw axis value (-1.0 to 1.0) for a gamepad axis.

- `gamepadId` (int) — Gamepad index.
- `axis` (int) — Axis index (0 = Left X, 1 = Left Y, 2 = Right X, 3 = Right Y, 4 = Left Trigger, 5 = Right Trigger).

**Returns:** `float`

---

### `Gamepad.Axis(gamepadId, axis)`

Facade for reading a gamepad axis value.

---

### `Gamepad.Button(gamepadId, button)`

Returns `TRUE` if a gamepad button is pressed.

- `gamepadId` (int) — Gamepad index.
- `button` (int) — Button index.

**Returns:** `bool`

---

### `Input.GamepadButtonCount(gamepadId)`

Returns the number of buttons on a gamepad.

**Returns:** `int`

---

### `Input.GamepadAxisCount(gamepadId)`

Returns the number of axes on a gamepad.

**Returns:** `int`

---

### `Input.SetGamepadMappings(mappingString)`

Loads a custom gamepad mapping string (SDL format). Useful for supporting non-standard controllers.

- `mappingString` (string) — SDL gamepad mapping string.

---

## Touch Input

### `Input.TouchCount()`

Returns the number of active touch points.

**Returns:** `int`

---

### `Input.TouchX(index)`

Returns the X position of a touch point.

- `index` (int) — Touch point index.

**Returns:** `int`

---

### `Input.TouchY(index)`

Returns the Y position of a touch point.

- `index` (int) — Touch point index.

**Returns:** `int`

---

### `Input.TouchPressed(index)`

Returns `TRUE` if a touch point was pressed this frame.

**Returns:** `bool`

---

### `Input.GetTouchPointID(index)`

Returns the platform ID for a touch point.

**Returns:** `int`

---

## Gestures

### `Gesture.Enable(flags)`

Enables specific gesture detection. Pass a bitmask of gesture types.

- `flags` (int) — Gesture type bitmask.

---

### `Gesture.IsDetected(gestureType)`

Returns `TRUE` if a specific gesture was detected this frame.

**Returns:** `bool`

---

### `Gesture.GetDetected()`

Returns the ID of the last detected gesture.

**Returns:** `int`

---

### `Gesture.GetHoldDuration()`

Returns how long a hold gesture has been active, in seconds.

**Returns:** `float`

---

### `Gesture.GetDragVectorX()` / `Gesture.GetDragVectorY()`

Returns the drag vector components.

**Returns:** `float`

---

### `Gesture.GetDragAngle()`

Returns the angle of the drag gesture.

**Returns:** `float`

---

### `Gesture.GetPinchVectorX()` / `Gesture.GetPinchVectorY()`

Returns the pinch vector components.

**Returns:** `float`

---

### `Gesture.GetPinchAngle()`

Returns the pinch gesture angle.

**Returns:** `float`

---

## Action Mapping

Action mapping lets you define abstract game actions and bind them to multiple input devices. This decouples your game logic from specific keys/buttons.

### `Action.MapKey(actionName, keyCode)`

Binds a keyboard key to a named action.

- `actionName` (string) — Action identifier (e.g., `"jump"`, `"fire"`).
- `keyCode` (int) — Key constant.

```basic
Action.MapKey("jump", KEY_SPACE)
Action.MapKey("jump", KEY_W)     ; Multiple keys for same action
Action.MapMouse("fire", 0)       ; Left click = fire
Action.MapJoy("jump", 0)         ; Gamepad button 0 = jump
```

---

### `Action.MapMouse(actionName, button)`

Binds a mouse button to a named action.

---

### `Action.MapJoy(actionName, button)`

Binds a gamepad button to a named action.

---

### `Action.MapAxis(actionName, gamepadAxis)`

Binds a gamepad axis to a named action.

---

### `Action.Down(actionName)`

Returns `TRUE` while the action is held.

- `actionName` (string) — Action name.

**Returns:** `bool`

```basic
IF Action.Down("fire") THEN
    FireBullet()
ENDIF
```

---

### `Action.Pressed(actionName)`

Returns `TRUE` on the frame the action was first triggered.

**Returns:** `bool`

---

### `Action.Released(actionName)`

Returns `TRUE` on the frame the action was released.

**Returns:** `bool`

---

### `Action.Value(actionName)`

Returns the analog value of an action (0.0–1.0 for buttons, -1.0–1.0 for axes).

**Returns:** `float`

---

### `Action.Reset()`

Clears all action bindings.

---

### `Input.SaveMappings(filePath)`

Saves action mappings to a file for persistence.

- `filePath` (string) — Output file path.

---

### `Input.LoadMappings(filePath)`

Loads action mappings from a file.

- `filePath` (string) — Mapping file path.

---

## Easy Mode Shortcuts

| Shortcut | Maps To |
|----------|---------|
| `KeyDown(code)` | `Input.KeyDown(code)` |
| `KeyHit(code)` | `Input.KeyPressed(code)` |
| `KEYDOWN(code)` | `Input.KeyDown(code)` |
| `KEYHIT(code)` | `Input.KeyPressed(code)` |
| `KEYUP(code)` | `Input.KeyUp(code)` |
| `MouseDown(btn)` | `Input.MouseDown(btn)` |
| `MouseHit(btn)` | `Input.MouseHit(btn)` |
| `MOUSEX` | `Input.MouseX()` |
| `MOUSEY` | `Input.MouseY()` |
| `MOUSEDX` | `Input.MouseDeltaX()` |
| `MOUSEDY` | `Input.MouseDeltaY()` |
| `MOUSEWHEEL` | `Input.MouseWheelMove()` |
| `AXIS()` | `Input.Axis()` |
| `JoyX(id)` | `Input.JoyX(id)` |
| `JoyY(id)` | `Input.JoyY(id)` |
| `JoyDown(id, btn)` | `Input.JoyButton(id, btn)` |

---

## Full Example

A complete input demo showing keyboard, mouse, and action mapping.

```basic
Window.Open(1280, 720, "Input Demo")
Window.SetFPS(60)

; Set up action mapping
Action.MapKey("move_left", KEY_A)
Action.MapKey("move_left", KEY_LEFT)
Action.MapKey("move_right", KEY_D)
Action.MapKey("move_right", KEY_RIGHT)
Action.MapKey("jump", KEY_SPACE)
Action.MapMouse("shoot", 0)

playerX = 640
playerY = 500
playerVelY = 0
gravity = 800
speed = 300

WHILE NOT Window.ShouldClose()
    dt = Time.Delta()

    ; Movement with action mapping
    IF Action.Down("move_left") THEN
        playerX = playerX - speed * dt
    ENDIF
    IF Action.Down("move_right") THEN
        playerX = playerX + speed * dt
    ENDIF

    ; Jump
    IF Action.Pressed("jump") AND playerY >= 500 THEN
        playerVelY = -400
    ENDIF

    ; Gravity
    playerVelY = playerVelY + gravity * dt
    playerY = playerY + playerVelY * dt
    IF playerY > 500 THEN
        playerY = 500
        playerVelY = 0
    ENDIF

    ; Shoot on click
    IF Action.Pressed("shoot") THEN
        PRINT "Shot at " + STR(Mouse.X()) + ", " + STR(Mouse.Y())
    ENDIF

    ; Render
    Render.Clear(20, 20, 30)
    Draw.Rectangle(playerX - 16, playerY - 32, 32, 32, 100, 200, 255, 255)
    Draw.Text("WASD/Arrows = Move | Space = Jump | Click = Shoot", 10, 10, 16, 200, 200, 200, 255)
    Draw.Text("Mouse: " + STR(Mouse.X()) + ", " + STR(Mouse.Y()), 10, 30, 16, 150, 150, 150, 255)
    Render.Frame()
WEND

Window.Close()
```

---

## See Also

- [WINDOW](WINDOW.md) — Window creation
- [CAMERA](CAMERA.md) — FPS camera mode uses mouse input
- [ACTION](ACTION.md) — Detailed action mapping reference
