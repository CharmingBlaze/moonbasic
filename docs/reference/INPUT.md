# Input Commands

Keyboard, mouse, cursor, and action-mapping input queries.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

Poll keys with `INPUT.KEYDOWN` / `INPUT.KEYPRESSED`, read mouse with `INPUT.MOUSEX` / `INPUT.MOUSEY`, and use `INPUT.AXIS` / `INPUT.MOVEMENT2D` for gameplay movement. Map named actions with `INPUT.MAPKEY`. Hide/lock the cursor with `CURSOR.*`.

---

## Keyboard and mouse

### `INPUT.KEYDOWN(key)` / `KEYPRESSED` / `KEYUP`
Keyboard state queries.

- **Arguments**:
    - `key`: (Integer) Key constant (e.g., `KEY_A`).
- **Returns**: (Boolean)

---

### `INPUT.MOUSEX()` / `MOUSEY()`
Returns the current mouse pixel coordinates.

- **Returns**: (Float)

---

### `INPUT.MOUSEDOWN(button)`
Returns `TRUE` if the mouse button is held down.

- **Arguments**:
    - `button`: (Integer) Mouse button constant.
- **Returns**: (Boolean)

---

### `INPUT.MOUSEWHEELMOVE()`
Returns the mouse wheel movement delta.

- **Returns**: (Float)

---

## Axes and movement helpers

### `INPUT.AXIS(negKey, posKey)`
Returns a value from -1.0 to 1.0 based on two keys.

- **Arguments**:
    - `negKey`: (Integer) Key for negative direction (-1.0).
    - `posKey`: (Integer) Key for positive direction (1.0).
- **Returns**: (Float)

---

### `INPUT.AXISDEG(negKey, posKey, degPerSec, dt)`
Keyboard orbit / yaw-style delta in degrees per second.

- **Arguments**:
    - `degPerSec`: (Float) Speed of rotation.
    - `dt`: (Float) Delta time.
- **Returns**: (Float)

---

### `INPUT.MOVEMENT2D(up, down, left, right)`
Returns a handle to a 2-element array `[forward, strafe]`.

- **Returns**: (Handle) A 2-float array.

---

### `INPUT.MOUSEDELTA()`
Returns a handle to a 2-element array `[dx, dy]` for relative movement.

- **Returns**: (Handle) A 2-float array.

---

## Cursor

### `CURSOR.HIDE()` / `CURSOR.SHOW()` 
Hides or shows the OS mouse cursor while over the window.

---

### `CURSOR.DISABLE()` / `CURSOR.ENABLE()` 
Disables the cursor and switches to relative mouse mode (centered virtual cursor).

---

## Action mapping

### `INPUT.MAPKEY(action, key)` 
Maps a keyboard key to an action.

---

### `INPUT.ACTIONDOWN(action)` / `INPUT.ACTIONPRESSED(action)` / `INPUT.ACTIONRELEASED(action)` 
Return whether mapped actions are held, newly pressed, or released.

---

### `INPUT.ACTIONAXIS(action)` 
Returns the analog axis value (-1.0 to 1.0) for the action.

---

## Keyboard constants

### Letters 
`KEY_A` to `KEY_Z`

---

### Numbers 
`KEY_ZERO` to `KEY_NINE`

---

### Function keys 
`KEY_F1` to `KEY_F12`

---

### Arrow keys 
`KEY_UP`, `KEY_DOWN`, `KEY_LEFT`, `KEY_RIGHT`

---

### Special keys 
`KEY_SPACE`, `KEY_ESCAPE`, `KEY_ENTER`, `KEY_TAB`, `KEY_BACKSPACE`  
`KEY_LEFT_SHIFT`, `KEY_LEFT_CONTROL`, `KEY_LEFT_ALT`  
`KEY_RIGHT_SHIFT`, `KEY_RIGHT_CONTROL`, `KEY_RIGHT_ALT`

---

## Mouse constants

Raylib-style names such as **`MOUSE_BUTTON_LEFT`** / **`MOUSE_BUTTON_RIGHT`** / **`MOUSE_BUTTON_MIDDLE`** appear in some samples; others use **`MOUSE_LEFT_BUTTON`**, **`MOUSE_RIGHT_BUTTON`**, **`MOUSE_MIDDLE_BUTTON`**. Use the identifiers your build binds to the active Raylib layer.

---

## Easy Mode name map (compatibility only)

| Dotted facade | Registry |
|---------------|----------|
| `Input.KeyDown` | `INPUT.KEYDOWN` |
| `Input.KeyPressed` | `INPUT.KEYPRESSED` |
| `Input.KeyUp` | `INPUT.KEYUP` |
| `Input.MouseX` / `MouseY` | `INPUT.MOUSEX` / `INPUT.MOUSEY` |
| `Input.MouseButtonDown` | `INPUT.MOUSEDOWN` |
| `Input.MouseWheelMove` | `INPUT.MOUSEWHEELMOVE` |
| `Input.Axis` | `INPUT.AXIS` |
| `Input.AxisDeg` | `INPUT.AXISDEG` |
| `Input.Movement2D` | `INPUT.MOVEMENT2D` |
| `Input.MouseDelta` | `INPUT.MOUSEDELTA` |
| `Input.MapKey` | `INPUT.MAPKEY` |
| `Input.ActionDown` / `ActionPressed` / `ActionAxis` | `INPUT.ACTIONDOWN` / `ACTIONPRESSED` / `ACTIONAXIS` |
| `Cursor.Hide` / `Show` / `Disable` / `Enable` | `CURSOR.HIDE` / `SHOW` / `DISABLE` / `ENABLE` |

---

## Full Example

```basic
WINDOW.OPEN(640, 480, "Input Demo")
WINDOW.SETFPS(60)

x = 320
y = 240

WHILE NOT WINDOW.SHOULDCLOSE()
    IF INPUT.KEYDOWN(KEY_RIGHT) THEN x = x + 3
    IF INPUT.KEYDOWN(KEY_LEFT)  THEN x = x - 3
    IF INPUT.KEYDOWN(KEY_DOWN)  THEN y = y + 3
    IF INPUT.KEYDOWN(KEY_UP)    THEN y = y - 3

    IF INPUT.KEYPRESSED(KEY_SPACE) THEN PRINT "Jump!"

    mx = INPUT.MOUSEX()
    my = INPUT.MOUSEY()

    RENDER.CLEAR(20, 20, 30)
    DRAW.CIRCLE(x, y, 16, 255, 100, 100, 255)
    DRAW.CIRCLE(mx, my, 8, 100, 255, 100, 255)
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## Extended Command Reference

### Mouse deltas & speed

| Command | Description |
|--------|-------------|
| `INPUT.MOUSEDELTAX()` / `INPUT.MOUSEDX()` | Mouse X movement this frame in pixels. |
| `INPUT.MOUSEDELTAY()` / `INPUT.MOUSEDY()` | Mouse Y movement this frame in pixels. |
| `INPUT.MOUSEXSPEED()` | Mouse X movement speed (px/s). |
| `INPUT.MOUSEYSPEED()` | Mouse Y movement speed (px/s). |
| `INPUT.SETMOUSEPOS(x, y)` | Warp mouse cursor to screen position. |
| `INPUT.SETMOUSEOFFSET(dx, dy)` | Offset applied to reported mouse position. |
| `INPUT.SETMOUSESCALE(sx, sy)` | Scale applied to mouse delta values. |
| `INPUT.LOCKMOUSE(bool)` | Lock/unlock cursor to window centre. |
| `INPUT.GETMOUSEWORLDPOS(cam)` | Returns `[x,y,z]` world position under mouse cursor. |

### Additional key helpers

| Command | Description |
|--------|-------------|
| `INPUT.KEYHIT(key)` | Alias of `INPUT.KEYPRESSED`. |
| `INPUT.MOUSEHIT(btn)` | Returns `TRUE` on the first frame mouse `btn` is pressed. |
| `INPUT.MOUSEPRESSED(btn)` | Alias of `INPUT.MOUSEHIT`. |
| `INPUT.MOUSERELEASED(btn)` | Returns `TRUE` on the first frame `btn` is released. |
| `INPUT.CHARPRESSED()` | Returns the Unicode code point of the last typed character. |
| `INPUT.GETKEYNAME(key)` | Returns the human-readable name string for a key constant. |
| `INPUT.GETINACTIVITY()` | Returns seconds since any input was received. |

### Gamepad

| Command | Description |
|--------|-------------|
| `INPUT.ISGAMEPADAVAILABLE(id)` | Returns `TRUE` if gamepad `id` is connected. |
| `INPUT.JOYX(id)` / `INPUT.JOYY(id)` | Left stick X/Y for gamepad `id`. |
| `INPUT.JOYBUTTON(id, btn)` | Returns `TRUE` if button `btn` is held on gamepad `id`. |
| `INPUT.JOYDOWN(id, btn)` | Alias of `INPUT.JOYBUTTON`. |
| `INPUT.GAMEPADAXISCOUNT(id)` | Number of axes on gamepad `id`. |
| `INPUT.GAMEPADBUTTONCOUNT(id)` | Number of buttons on gamepad `id`. |
| `INPUT.GETGAMEPADAXISVALUE(id, axis)` | Raw axis value -1.0..1.0. |
| `INPUT.MAPGAMEPADAXIS(id, axis, action, deadzone)` | Map axis to named action. |
| `INPUT.MAPGAMEPADBUTTON(id, btn, action)` | Map button to named action. |
| `INPUT.SETGAMEPADMAPPINGS(csv)` | Set SDL-format gamepad mapping string. |
| `INPUT.LOADMAPPINGS(path)` | Load action-to-input mappings from a file. |
| `INPUT.SAVEMAPPINGS(path)` | Save current mappings to a file. |

### Touch

| Command | Description |
|--------|-------------|
| `INPUT.TOUCHCOUNT()` | Number of active touch points. |
| `INPUT.TOUCHX(index)` / `INPUT.TOUCHY(index)` | Position of touch point `index`. |
| `INPUT.TOUCHPRESSED(index)` | Returns `TRUE` on the first frame touch `index` is down. |
| `INPUT.GETTOUCHPOINTID(index)` | Returns the OS touch point id. |

### Camera-relative movement

| Command | Description |
|--------|-------------|
| `INPUT.MOVEDIR(camYaw)` | Returns `[dx, dz]` camera-relative WASD direction vector (normalized). |

## See also

- [GAMEPAD.md](GAMEPAD.md) — `GAMEPAD.*` namespace
- [CURSOR.md](CURSOR.md) — cursor visibility and lock
- [KEY.md](KEY.md) — `KEY_*` constants
