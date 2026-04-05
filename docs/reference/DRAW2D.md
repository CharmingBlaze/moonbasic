# 2D Drawing Commands

Raylib-backed 2D drawing. Typical frame flow: `Render.Clear` → (optional) `Camera2D.Begin` / `Camera2D.End` → `Draw.*` → `Render.Frame`.

- **No 2D camera:** coordinates are screen pixels (top-left origin).
- **With `Camera2D`:** world coordinates are transformed by offset, target, zoom, and rotation.

Color components `r, g, b, a` are **0–255** unless noted.

---

## Shapes (filled / strokes)

| Command | Arguments (summary) |
|---------|---------------------|
| `Draw.Rectangle` | `x, y, w, h, r, g, b, a` |
| `Draw.RectangleRounded` | `x, y, w, h, radius, r, g, b, a` |
| `Draw.RectLines` | `x, y, w, h, thick, r, g, b, a` |
| `Draw.RectPro` | `x, y, w, h, ox, oy, rot, r, g, b, a` — rotated rectangle; `rot` in degrees |
| `Draw.RectGradV` | `x, y, w, h, top rgba, bottom rgba` (12 args) |
| `Draw.RectGradH` | `x, y, w, h, left rgba, right rgba` (12 args) |
| `Draw.RectGrad` | `x, y, w, h, four corners × rgba` (20 args) |
| `Draw.Circle` | `cx, cy, radius, r, g, b, a` |
| `Draw.CircleLines` | same as `Draw.Circle` |
| `Draw.CircleSector` | `cx, cy, radius#, start#, end#, segments, r, g, b, a` — angles in degrees |
| `Draw.CircleGradient` | `cx, cy, radius#, inner r,g,b,a, outer r,g,b,a` |
| `Draw.Ellipse` | `cx, cy, rx, ry, r, g, b, a` |
| `Draw.EllipseLines` | same as `Draw.Ellipse` |
| `Draw.Ring` | `cx, cy, innerR#, outerR#, start#, end#, segments, r, g, b, a` |
| `Draw.RingLines` | same arity as `Draw.Ring` |
| `Draw.Triangle` | `x1,y1, x2,y2, x3,y3, r, g, b, a` |
| `Draw.TriangleLines` | same arity as `Draw.Triangle` |
| `Draw.Poly` | `cx, cy, sides, radius#, rotation#, r, g, b, a` |
| `Draw.PolyLines` | `cx, cy, sides, radius#, rotation#, thick, r, g, b, a` |

### UI-style helpers

| Command | Arguments (summary) |
|---------|---------------------|
| `Draw.ProgressBar` | `x, y, w, h, t#, r, g, b, a` — **`t`** in **0..1**; gray track + filled bar |
| `Draw.HealthBar` | `x, y, w, h, current#, max#, r, g, b, a` — fill ratio **`current/max`** |
| `Draw.CenterText` | `text$, y, size, r, g, b, a` — horizontally centered on the screen |
| `Draw.RightText` | `text$, marginRight, y, size, r, g, b, a` — right-aligned with **`marginRight`** from the right edge |
| `Draw.ShadowText` | `text$, x, y, size, r, g, b, a` — offset shadow |
| `Draw.OutlineText` | `text$, x, y, size, r, g, b, a` — 8-neighbor outline |
| `Draw.Crosshair` | `cx, cy, radius, r, g, b, a` |
| `Draw.RectGrid` | `x, y, cellW, cellH, cols, rows, line rgba, fill rgba` (14 args) |

---

## Lines and splines

| Command | Arguments |
|---------|-----------|
| `Draw.Line` | `x1, y1, x2, y2, r, g, b, a` |
| `Draw.LineEx` | `x1, y1, x2, y2, thick#, r, g, b, a` |
| `Draw.LineBezier` | `x1,y1, x2,y2, thick#, r, g, b, a` |
| `Draw.LineBezierQuad` | `x1,y1, cx,cy, x2,y2, thick#, r, g, b, a` |
| `Draw.LineBezierCubic` | `x1,y1, c1x,c1y, c2x,c2y, x2,y2, thick#, r, g, b, a` |
| `Draw.SplineLinear` | `(pointsArray, thick#, r, g, b, a)` |
| `Draw.SplineBasis` | same |
| `Draw.SplineCatmullRom` | same |
| `Draw.SplineBezierQuad` | same |
| `Draw.SplineBezierCubic` | same |

**`pointsArray`:** heap array handle with an **even** length of floats: `x0,y0,x1,y1,...`.

---

## Textures

All texture commands take a **texture handle** from `Texture.Load` / `Texture.FromImage`, etc.

| Command | Arguments (summary) |
|---------|---------------------|
| `Draw.Texture` | `tex, x, y, r, g, b, a` — integer pixel position |
| `Draw.TextureV` | `tex, x, y, r, g, b, a` — float position |
| `Draw.TextureEx` | `tex, x, y, rot#, scale#, r, g, b, a` |
| `Draw.TextureRec` | `tex, srcX#, srcY#, srcW#, srcH#, x, y, r, g, b, a` |
| `Draw.TexturePro` | `tex, src rect (4 floats), dest rect (4 floats), ox#, oy#, rot#, r, g, b, a` |
| `Draw.TextureFull` | `tex` — stretches full texture to the screen |
| `Draw.TextureFlipped` | `tex` — draws render-target texture flipped (Y) full screen |
| `Draw.TextureTiled` | `tex, src rect (4 floats), dest rect (4 floats), ox#, oy#, rot#, scale#, r, g, b, a` (17 args). Tiles with `Draw.TexturePro` internally. For **`rot# = 0` and `ox# = oy# = 0`** behavior matches tiled fills; non-zero rotation/origin uses a **single** `Draw.TexturePro` over the destination (not per-tile). |
| `Draw.TextureNPatch` | `tex, L, T, R, B, x, y, w, h, r, g, b, a` — border widths then destination rect |

---

## Text and measurement

| Command | Arguments |
|---------|-----------|
| `Draw.Text` | `text$, x, y, size, r, g, b, a` — default font |
| `Draw.TextEx` | `font, text$, x#, y#, size#, spacing#, r, g, b, a` |
| `Draw.TextFont` | **Alias of `Draw.TextEx`** (same 10 arguments) |
| `Draw.TextPro` | `font, text$, x#, y#, ox#, oy#, rot#, size#, spacing#, r, g, b, a` |
| `Draw.TextWidth` | `text$, size` → width in pixels |
| `Draw.TextFontWidth` | `font, text$, size#, spacing#` → `[w#, h#]` array |
| `MeasureText` | same as `Draw.TextWidth` |
| `MeasureTextEx` | `font, text$, size#, spacing#` → `[w#, h#]` array |
| `GetTextCodepointCount` | `text$` → integer |

---

## Pixels, arc, grid

| Command | Arguments |
|---------|-----------|
| `Draw.Pixel` / `Draw.PixelV` | `x, y, r, g, b, a` (int vs float coordinates) |
| `Draw.SetPixelColor` | Alias of `Draw.Pixel` |
| `Draw.GetPixelColor` | `x, y` → **handle** to a 4-element float array `[r, g, b, a]` (0–255) sampled from the screen |
| `Draw.Dot` | `x, y, size#, r, g, b, a` |
| `Draw.Arc` | `cx, cy, radius#, start#, end#, thick#, r, g, b, a` |
| `Draw.Grid2D` | `spacing#, r, g, b, a` — 2D cell grid in screen space |

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
player_x# = 375
player_y# = 275

WHILE NOT Window.ShouldClose()
    dt# = Time.Delta()
    IF Input.KeyDown(KEY_RIGHT) THEN player_x# = player_x# + 200 * dt#
    IF Input.KeyDown(KEY_LEFT) THEN player_x# = player_x# - 200 * dt#

    Render.Clear(14, 22, 33)
    Draw.Rectangle(0, 500, 800, 100, 40, 50, 60, 255)
    Draw.Texture(tex, INT(player_x#), INT(player_y#), 255, 255, 255, 255)
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
