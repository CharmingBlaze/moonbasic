# Cursor Commands

Control the system cursor: show, hide, enable/disable tracking, check visibility, and set a custom cursor image.

See [MOUSE.md](MOUSE.md) for mouse position and button state.

## Core Workflow

- Call `CURSOR.HIDE()` at game start for FPS / action games.
- Use `CURSOR.SET(id)` to switch between system cursor shapes.
- Re-show with `CURSOR.SHOW()` on exit or in menus.

---

## Visibility

### `CURSOR.SHOW()` 

Makes the system cursor visible.

---

### `CURSOR.HIDE()` 

Hides the system cursor. The cursor still moves internally.

---

### `CURSOR.ISHIDDEN()` 

Returns `TRUE` if the cursor is currently hidden.

---

### `CURSOR.ISONSCREEN()` 

Returns `TRUE` if the cursor is within the window bounds.

---

## Enable / Disable

### `CURSOR.ENABLE()` 

Enables cursor input (default). Cursor moves freely.

---

### `CURSOR.DISABLE()` 

Disables cursor movement — locks the cursor to the window. Use for FPS camera control (read delta via `MOUSE.DX` / `MOUSE.DY`).

---

### `CURSOR.ISENABLED()` 

Returns `TRUE` if the cursor is currently enabled.

---

## Appearance

### `CURSOR.SET(cursorId)` 

Sets the cursor to a system cursor shape. `cursorId` constants (Raylib `MouseCursor`): `0` = default arrow, `1` = ibeam (text), `2` = crosshair, `3` = pointing hand, `4` = resize EW, `5` = resize NS, `6` = resize NWSE, `7` = resize NESW, `8` = resize all, `9` = not allowed.

---

## Full Example

FPS game start: hide cursor, re-show on menu.

```basic
WINDOW.OPEN(960, 540, "Cursor Demo")
WINDOW.SETFPS(60)

inMenu = TRUE

WHILE NOT WINDOW.SHOULDCLOSE()
    IF inMenu THEN
        CURSOR.SHOW()
        CURSOR.ENABLE()
        CURSOR.SET(0)   ; arrow

        IF MOUSE.PRESSED(0)
            inMenu = FALSE
            CURSOR.HIDE()
            CURSOR.DISABLE()
        END IF

        RENDER.CLEAR(30, 30, 50)
        DRAW.TEXT("Click to start", 320, 240, 24, 255, 255, 255, 255)
    ELSE
        ; FPS mode - read mouse delta for look
        dx = MOUSE.DX()
        dy = MOUSE.DY()

        RENDER.CLEAR(20, 25, 35)
        DRAW.TEXT("ESC to menu", 10, 10, 18, 200, 200, 200, 255)

        IF INPUT.KEYDOWN(KEY_ESCAPE)
            inMenu = TRUE
        END IF
    END IF

    RENDER.FRAME()
WEND

CURSOR.SHOW()
CURSOR.ENABLE()
WINDOW.CLOSE()
```

---

## See also

- [MOUSE.md](MOUSE.md) — mouse position, buttons, delta
- [INPUT.md](INPUT.md) — unified input system
- [WINDOW.md](WINDOW.md) — window focus and close
