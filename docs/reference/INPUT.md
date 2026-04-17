# Input Commands

Keyboard, mouse, cursor, and action-mapping input queries.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

Poll keys with `INPUT.KEYDOWN` / `INPUT.KEYPRESSED`, read mouse with `INPUT.MOUSEX` / `INPUT.MOUSEY`, and use `INPUT.AXIS` / `INPUT.MOVEMENT2D` for gameplay movement. Map named actions with `INPUT.MAPKEY`. Hide/lock the cursor with `CURSOR.*`.

---

## Keyboard and mouse

### `INPUT.KEYDOWN(key)` 
Returns `TRUE` if the key is currently held down.

---

### `INPUT.KEYPRESSED(key)` 
Returns `TRUE` on the frame the key was first pressed.

---

### `INPUT.KEYUP(key)` 
Returns `TRUE` on the frame the key was released.

---

### `INPUT.MOUSEX()` / `INPUT.MOUSEY()` 
Returns the current mouse pixel coordinates.

---

### `INPUT.MOUSEDOWN(button)` 
Returns `TRUE` if the mouse button is held down.

---

### `INPUT.MOUSEWHEELMOVE()` 
Returns the mouse wheel movement delta.

---

## Axes and movement helpers

### `INPUT.AXIS(negKey, posKey)` 
Returns a value from -1.0 to 1.0 based on two keys.

---

### `INPUT.AXISDEG(negKey, posKey, degPerSec, dt)` 
Keyboard orbit / yaw-style delta in degrees per second, scaled by **`dt`** (see [CAMERA.md](CAMERA.md)).

---

### `INPUT.MOVEMENT2D(up, down, left, right)` 
Returns a 2-element array handle `[forward, strafe]` based on four keys.

---

### `INPUT.MOUSEDELTA()` 
Returns a 2-element array handle `[dx, dy]` for movement since last frame (see also **`INPUT.MOUSEDELTAX`** / **`INPUT.MOUSEDELTAY`**).

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
