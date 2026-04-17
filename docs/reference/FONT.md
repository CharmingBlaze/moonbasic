# Font (`FONT.*`, `DRAW.TEXTFONT`)

**Conventions:** [STYLE_GUIDE.md](../../STYLE_GUIDE.md), [API_CONVENTIONS.md](API_CONVENTIONS.md) — reference pages use uppercase **`NAMESPACE.ACTION`**; Easy Mode (`Font.Load`, …) is [compatibility only](../../STYLE_GUIDE.md#easy-mode-compatibility-layer).

**Page shape:** [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) — see [WAVE.md](WAVE.md) (registry-first headings, **Full Example** at the end).

## Core Workflow

1. **Load:** **`FONT.LOAD(path)`** — `.ttf` / `.otf`; store the returned handle.
2. **Draw:** **`DRAW.TEXTFONT(handle, text, x, y, size, spacing, r, g, b, a)`** inside **`CAMERA2D.BEGIN`** / **`CAMERA2D.END`** (or your active 2D camera bracket).
3. **Free:** **`FONT.FREE(handle)`** when done.

---

### `FONT.LOAD(path)` 
Loads a `.ttf` or `.otf` font file from disk. Returns a **font handle**.

---

### `FONT.FREE(handle)` 
Unloads a font from memory and releases its heap slot.

---

### `DRAW.TEXTFONT(handle, text, x, y, size, spacing, r, g, b, a)` 
Draws text using a specific font handle. Call within **`CAMERA2D.BEGIN`** / **`CAMERA2D.END`** (or your active 2D camera bracket).
- `handle`: The handle of the loaded font.
- `text`: The string to draw.
- `x`, `y`: Screen position.
- `size`: Font size in pixels.
- `spacing`: Extra spacing between characters.
- `r, g, b, a`: Color components (0-255).

---

## Full Example

This example assumes you have a font file named `my_font.ttf` in the same directory as your script.

```basic
WINDOW.OPEN(800, 600, "Font Example")
WINDOW.SETFPS(60)

myFont = FONT.LOAD("my_font.ttf")

IF myFont = 0 THEN
    PRINT("Error: Could not load my_font.ttf")
    WINDOW.CLOSE()
    SYSTEM.EXIT()
ENDIF

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(50, 60, 70)

    CAMERA2D.BEGIN()
        DRAW.TEXTFONT(myFont, "Hello, moonBASIC!", 100, 200, 48, 2, 255, 200, 100, 255)
        DRAW.TEXT("This is the default system font.", 100, 300, 20, 200, 200, 200, 255)
    CAMERA2D.END()

    RENDER.FRAME()
WEND

FONT.FREE(myFont)
WINDOW.CLOSE()
```
