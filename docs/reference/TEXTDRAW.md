# TextDraw / TextExObj Commands

Retained-mode text draw objects. Configure once, update text or position each frame, and call `DRAW` once — avoids rebuilding text draw calls from scratch every frame.

Two handle types:
- **`TEXTDRAW`** — default font, pixel-integer position, integer size.
- **`TEXTEXOBJ`** — custom font handle, float position, float size and spacing.

## Core Workflow

1. Create via `TEXT.MAKE(text, x, y, size, r, g, b, a)` → returns a **TEXTDRAW handle** — or via the appropriate constructor.
2. Update each frame with `TEXTDRAW.SETTEXT` / `TEXTDRAW.POS` / `TEXTDRAW.COLOR`.
3. `TEXTDRAW.DRAW(handle)` to render.
4. `TEXTDRAW.FREE(handle)` when done.

---

## TEXTDRAW Commands

### `TEXTDRAW.POS(handle, x, y)` 

Sets the screen position in pixels.

---

### `TEXTDRAW.SIZE(handle, size)` 

Sets the font size in pixels.

---

### `TEXTDRAW.COLOR(handle, r, g, b, a)` / `TEXTDRAW.COL(handle, r, g, b, a)` 

Sets the text color (0–255).

---

### `TEXTDRAW.SETTEXT(handle, text)` 

Updates the text string.

---

### `TEXTDRAW.DRAW(handle)` 

Renders the text this frame.

---

### `TEXTDRAW.FREE(handle)` 

Frees the handle.

---

## TEXTEXOBJ Commands

### `TEXTEXOBJ.POS(handle, x, y)` 

Sets position as floats.

---

### `TEXTEXOBJ.SIZE(handle, size)` 

Sets font size as float.

---

### `TEXTEXOBJ.SPACING(handle, spacing)` 

Sets character spacing.

---

### `TEXTEXOBJ.COLOR(handle, r, g, b, a)` 

Sets text color.

---

### `TEXTEXOBJ.SETTEXT(handle, text)` 

Updates the text string.

---

### `TEXTEXOBJ.DRAW(handle)` 

Renders this frame.

---

### `TEXTEXOBJ.FREE(handle)` 

Frees the handle.

---

## Full Example

Score display that updates each frame without rebuilding.

```basic
WINDOW.OPEN(800, 450, "TextDraw Demo")
WINDOW.SETFPS(60)

score = 0
lbl = TEXT.MAKE("Score: 0", 10, 10, 24, 255, 255, 255, 255)

WHILE NOT WINDOW.SHOULDCLOSE()
    IF INPUT.KEYPRESSED(KEY_SPACE) THEN
        score = score + 10
        TEXTDRAW.SETTEXT(lbl, "Score: " + STR(score))
    END IF

    RENDER.CLEAR(20, 20, 40)
    TEXTDRAW.DRAW(lbl)
    RENDER.FRAME()
WEND

TEXTDRAW.FREE(lbl)
WINDOW.CLOSE()
```

---

## See also

- [DRAW.md](DRAW.md) — immediate-mode `DRAW.TEXT`
- [FONT.md](FONT.md) — `FONT.LOAD` for custom fonts with `TEXTEXOBJ`
