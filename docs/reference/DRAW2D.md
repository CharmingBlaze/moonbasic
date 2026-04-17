# 2D Drawing Commands

Immediate-mode 2D shapes, text, and texture drawing on the screen framebuffer.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. `RENDER.CLEAR` to start the frame.
2. Draw with `DRAW.RECTANGLE`, `DRAW.CIRCLE`, `DRAW.LINE`, `DRAW.TEXT`, `DRAW.TEXTURE`, etc.
3. `RENDER.FRAME` to present.

Color components `r, g, b, a` are **0–255**. For 2D camera transforms see `CAMERA2D.BEGIN` / `CAMERA2D.END` in [CAMERA.md](CAMERA.md).

---

### `DRAW.RECTANGLE(x, y, w, h, r, g, b, a)` 
Draws a filled rectangle at the specified screen coordinates.
- `x, y`: Top-left corner.
- `w, h`: Dimensions.
- `r, g, b, a`: Color components (0-255).

---

### `DRAW.RECTLINES(x, y, w, h, thick, r, g, b, a)` 
Draws a rectangle outline with a specific thickness.

---

### `DRAW.CIRCLE(cx, cy, radius, r, g, b, a)` 
Draws a filled circle.
- `cx, cy`: Center position.
- `radius`: Circle radius.

---

### `DRAW.LINE(x1, y1, x2, y2, r, g, b, a)` 
Draws a line between two points.

---

### `DRAW.TEXT(text, x, y, size, r, g, b, a)` 
Draws text using the default font.
- `text`: The string to display.
- `size`: Font size in pixels.

---

### `DRAW.TEXTURE(id, x, y, r, g, b, a)` 
Draws a texture handle at the specified position with a tint color. Use `255, 255, 255, 255` for no tint.

---

### `DEBUG.PRINT(template, v0 [, v1 …])` 

Quick **debug HUD** lines: **`template`** uses placeholders **`{0}`** … **`{9}`**, filled from the following values. Draws with the default font at a fixed top-left column, **stacking downward** each frame; the vertical cursor **resets** when the render **frame** advances (same timing as **`RENDER.FRAME`** / runtime frame counter). For positioned or styled HUD text, use **`DRAW.TEXT`** instead.

---

## Pixels, arc, grid

| Command | Arguments |
|---------|-----------|
| `DRAW.PIXEL` / `DRAW.PIXELV` | `x, y, r, g, b, a` — integer vs float screen coordinates |
| `DRAW.SETPIXELCOLOR` | Alias of **`DRAW.PIXEL`** |
| `DRAW.GETPIXELCOLOR` | `x, y` → **handle** to a 4-element float array `[r, g, b, a]` (0–255) sampled from the screen |
| `DRAW.DOT` | `x, y, size, r, g, b, a` |
| `DRAW.ARC` | `cx, cy, radius, start, end, thick, r, g, b, a` |
| `DRAW.GRID2D` | `spacing, r, g, b, a` — 2D cell grid in screen space |

---

## Overlay

| Command | Arguments |
|---------|-----------|
| `RENDER.DRAWFPS` | `x, y` — draws FPS counter (registered on the render module) |

---

## Full Example

```basic
WINDOW.OPEN(800, 600, "2D Drawing")
WINDOW.SETFPS(60)
tex = TEXTURE.LOAD("player.png")
player_x = 375
player_y = 275

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    IF INPUT.KEYDOWN(KEY_RIGHT) THEN player_x = player_x + 200 * dt
    IF INPUT.KEYDOWN(KEY_LEFT) THEN player_x = player_x - 200 * dt

    RENDER.CLEAR(14, 22, 33)
    DRAW.RECTANGLE(0, 500, 800, 100, 40, 50, 60, 255)
    DRAW.TEXTURE(tex, INT(player_x), INT(player_y), 255, 255, 255, 255)
    DRAW.RECTANGLE(10, 10, 200, 20, 200, 0, 0, 255)
    DRAW.TEXT("HEALTH", 15, 12, 16, 255, 255, 255, 255)
    RENDER.FRAME()
WEND

TEXTURE.FREE(tex)
WINDOW.CLOSE()
```

---

## See also

- [DRAW3D.md](DRAW3D.md) — 3D primitives (`Draw3D.*` / `DRAW.*` aliases).
- [CAMERA.md](CAMERA.md) — **`CAMERA2D.*`** transforms.
- [RENDER.md](RENDER.md) — clear / frame / render state.
- [TEXTURE.md](TEXTURE.md) — loading and freeing textures.
