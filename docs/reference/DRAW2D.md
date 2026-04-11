# 2D Drawing Commands

Raylib-backed 2D drawing. Typical frame flow: `Render.Clear` → (optional) `Camera2D.Begin` / `Camera2D.End` → `Draw.*` → `Render.Frame`.

**BlitzPlus-style names** (`Plot`, `Line`, `Rect`, …) map to these **`DRAW.*`** commands with different color/buffer rules — see [BLITZ_COMMAND_INDEX.md](BLITZ_COMMAND_INDEX.md). Aliases: **`DRAW.PLOT`** = **`DRAW.PIXEL`**; **`DRAW.OVAL`** / **`DRAW.OVALLINES`** = **`DRAW.ELLIPSE`** / **`DRAW.ELLIPSELINES`**.

- **No 2D camera:** coordinates are screen pixels (top-left origin).
- **With `Camera2D`:** world coordinates are transformed by offset, target, zoom, and rotation.

Color components `r, g, b, a` are **0–255** unless noted.

---

### `Draw.Rectangle(x, y, w, h, r, g, b, a)`
Draws a filled rectangle at the specified screen coordinates.
- `x, y`: Top-left corner.
- `w, h`: Dimensions.
- `r, g, b, a`: Color components (0-255).

### `Draw.RectLines(x, y, w, h, thick, r, g, b, a)`
Draws a rectangle outline with a specific thickness.

---

### `Draw.Circle(cx, cy, radius, r, g, b, a)`
Draws a filled circle.
- `cx, cy`: Center position.
- `radius`: Circle radius.

### `Draw.Line(x1, y1, x2, y2, r, g, b, a)`
Draws a line between two points.

---

### `Draw.Text(text, x, y, size, r, g, b, a)`
Draws text using the default font.
- `text`: The string to display.
- `size`: Font size in pixels.

### `Draw.Texture(id, x, y, r, g, b, a)`
Draws a texture handle at the specified position with a tint color. Use `255, 255, 255, 255` for no tint.

### `DEBUG.PRINT(template, v0 [, v1 …])`

Quick **debug HUD** lines: **`template`** uses placeholders **`{0}`** … **`{9}`**, filled from the following values. Draws with the default font at a fixed top-left column, **stacking downward** each frame; the vertical cursor **resets** when the render **frame** advances (same timing as **`RENDER.FRAME`** / runtime frame counter). For positioned or styled HUD text, use **`Draw.Text`** instead.

---

## Pixels, arc, grid

| Command | Arguments |
|---------|-----------|
| `Draw.Pixel` / `Draw.PixelV` | `x, y, r, g, b, a` (int vs float coordinates) |
| `Draw.SetPixelColor` | Alias of `Draw.Pixel` |
| `Draw.GetPixelColor` | `x, y` → **handle** to a 4-element float array `[r, g, b, a]` (0–255) sampled from the screen |
| `Draw.Dot` | `x, y, size, r, g, b, a` |
| `Draw.Arc` | `cx, cy, radius, start, end, thick, r, g, b, a` |
| `Draw.Grid2D` | `spacing, r, g, b, a` — 2D cell grid in screen space |

---

## Overlay

| Command | Arguments |
|---------|-----------|
| `Render.DrawFPS` | `x, y` — draws FPS counter (registered on the render module; same as `Draw` group in docs tables) |

---

## Full example

```basic
Window.Open(800, 600, "2D Drawing")
Window.SetFPS(60)
tex = Texture.Load("player.png")
player_x = 375
player_y = 275

WHILE NOT Window.ShouldClose()
    dt = Time.Delta()
    IF Input.KeyDown(KEY_RIGHT) THEN player_x = player_x + 200 * dt
    IF Input.KeyDown(KEY_LEFT) THEN player_x = player_x - 200 * dt

    Render.Clear(14, 22, 33)
    Draw.Rectangle(0, 500, 800, 100, 40, 50, 60, 255)
    Draw.Texture(tex, INT(player_x), INT(player_y), 255, 255, 255, 255)
    Draw.Rectangle(10, 10, 200, 20, 200, 0, 0, 255)
    Draw.Text("HEALTH", 15, 12, 16, 255, 255, 255, 255)
    Render.Frame()
WEND

Texture.Free(tex)
Window.Close()
```

---

## See also

- [DRAW3D.md](DRAW3D.md) — 3D primitives (`Draw3D.*` / `Draw.*` aliases).
- [CAMERA.md](CAMERA.md) — `Camera2D` transforms.
- [RENDER.md](RENDER.md) — clear / frame / render state.
- [TEXTURE.md](TEXTURE.md) — loading and freeing textures.
