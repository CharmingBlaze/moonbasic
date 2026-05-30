# Mouse Commands

Read mouse position, button state, scroll wheel, and control cursor visibility. For full input see [INPUT.md](INPUT.md).

## Core Workflow

1. Call `MOUSE.X()` / `MOUSE.Y()` each frame for cursor position.
2. Use `MOUSE.DOWN(btn)` / `MOUSE.PRESSED(btn)` / `MOUSE.RELEASED(btn)` for button state.
3. Read `MOUSE.WHEEL()` for scroll delta.
4. `MOUSE.HIDE()` / `MOUSE.SHOW()` to toggle cursor visibility.

Button constants: `MOUSE_LEFT` = 0, `MOUSE_RIGHT` = 1, `MOUSE_MIDDLE` = 2.

---

## Position

### `MOUSE.X()` 

Returns the current mouse cursor X position in screen pixels.

---

### `MOUSE.Y()` 

Returns the current mouse cursor Y position in screen pixels.

---

### `MOUSE.POSX()` 

Alias of `MOUSE.X()`.

---

### `MOUSE.POSY()` 

Alias of `MOUSE.Y()`.

---

### `MOUSE.DX()` 

Returns the horizontal delta (movement) since the last frame in pixels. Useful for camera rotation when the cursor is locked/hidden.

---

### `MOUSE.DY()` 

Returns the vertical delta since the last frame in pixels.

---

### `MOUSE.SETPOS(x, y)` 

Moves the cursor to screen position `(x, y)`.

---

## Buttons

### `MOUSE.DOWN(button)` 

Returns `TRUE` while `button` is held. `button`: `0` left, `1` right, `2` middle.

---

### `MOUSE.PRESSED(button)` 

Returns `TRUE` on the **first frame** the button is pressed (one-shot, not held).

---

### `MOUSE.RELEASED(button)` 

Returns `TRUE` on the **first frame** the button is released.

---

## Scroll Wheel

### `MOUSE.WHEEL()` 

Returns the scroll wheel delta for this frame. Positive = scroll up, negative = scroll down.

---

## Cursor Visibility

### `MOUSE.HIDE()` 

Hides the system cursor. The cursor still moves but is invisible.

---

### `MOUSE.SHOW()` 

Makes the system cursor visible again.

---

### `MOUSE.ENABLE()` 

Enables cursor input (default state).

---

### `MOUSE.DISABLE()` 

Disables cursor movement tracking (locks the cursor to the window for FPS-style control).

---

## Full Example

FPS-style camera look using mouse delta with left-click to shoot.

```basic
WINDOW.OPEN(960, 540, "Mouse Demo")
WINDOW.SETFPS(60)

MOUSE.HIDE()
MOUSE.DISABLE()

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 2, 0)
CAMERA.SETTARGET(cam, 0, 2, -1)

yaw   = 0.0
pitch = 0.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()

    yaw   = yaw   - MOUSE.DX() * 0.2
    pitch = pitch - MOUSE.DY() * 0.2
    pitch = MAX(-80, MIN(80, pitch))

    CAMERA.SETORBIT(cam, 0, 2, 0, yaw, pitch, 0.01)

    IF MOUSE.PRESSED(0) THEN PRINT "Shoot!"

    ; zoom with scroll
    fov = CAMERA.GETFOV(cam) - MOUSE.WHEEL() * 2
    CAMERA.SETFOV(cam, MAX(20, MIN(100, fov)))

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        DRAW3D.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

MOUSE.SHOW()
MOUSE.ENABLE()
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `MOUSE.SETPOSITION(x, y)` | Warp the OS cursor to screen position `(x, y)`. |

---

## See also

- [INPUT.md](INPUT.md) — unified keyboard, mouse, and gamepad API
- [CURSOR.md](CURSOR.md) — `CURSOR.SET` for custom cursor image
- [CAMERA.md](CAMERA.md) — camera orbit driven by mouse delta
