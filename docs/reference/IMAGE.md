# Image (CPU) — `Image.*` / `IMAGE.*`

**CPU-side** pixel buffers (Raylib `Image`): read, mutate, save. **Not** GPU textures (`Texture.*`). Typical pipeline: **`IMAGE.CREATE`** (canonical) or deprecated **`IMAGE.MAKE`** / **`IMAGE.LOAD`** → optional edits → **`TEXTURE.FROMIMAGE`** → **`DRAW.TEXTURE`** on the main framebuffer → free when done.

**Requires CGO** and Raylib (same as `Draw.*`, `Texture.*`).

Registry keys use **dots and uppercase** (e.g. **`IMAGE.CREATE`**). PascalCase names below match docs/spec style.

---

### `Image.Load(path)`
Loads an image from disk (PNG, JPG, BMP, etc.). Returns a **handle**.

### `Image.Make(w, h [, r, g, b, a])`
Creates a new CPU image. If RGBA components are provided, fills the image with that color (0-255).

### `Image.Free(handle)`
Releases the heap slot and unloads the image memory.

---

### `Image.Width(handle)` / `Image.Height(handle)`
Returns the integer pixel dimensions of the image.

### `Image.Resize(handle, w, h)`
Resizes the image in memory using bilinear scaling.

### `Image.Export(handle, path)`
Saves the image to a file. The format is determined by the file extension.

---

### `Image.DrawPixel(handle, x, y, r, g, b, a)`
Draws a single pixel on the image at `(x, y)`.

### `Image.DrawRect(handle, x, y, w, h, r, g, b, a)`
Draws a filled rectangle on the image.

---

## Example: composite → texture → draw

```basic
Window.Open(640, 480, "Image to texture")
Window.SetFPS(60)

a = IMAGE.CREATE(128, 128)
IMAGE.CLEAR(a, 40, 40, 50, 255)
b = IMAGE.CREATE(32, 32, 200, 80, 80, 255)
IMAGE.DRAWIMAGE(a, b, 0, 0, 32, 32, 48, 48, 32, 32, 255, 255, 255, 255)
IMAGE.FREE(b)

tex = TEXTURE.FROMIMAGE(a)
IMAGE.FREE(a)

WHILE NOT Window.ShouldClose()
    Render.Clear(20, 24, 32)
    DRAW.TEXTURE(tex, 200, 160, 255, 255, 255, 255)
    Render.Frame()
WEND

TEXTURE.FREE(tex)
Window.Close()
```

---

## Common mistakes

- **`IMAGE.*` vs GPU** — To display pixels, use **`TEXTURE.FROMIMAGE`** then **`DRAW.TEXTURE`** (or equivalent), not **`IMAGE.*`** alone.
- **Unpaired `IMAGE.FREE`** — Each load/create should be freed when done.
- **`IMAGE.CREATE(w,h)`** — Transparent, not black-opaque, until you clear or paint.

---

## See also

- [TEXTURE.md](TEXTURE.md) — **`TEXTURE.FROMIMAGE`**, render targets
- [DRAW2D.md](DRAW2D.md) — screen-space drawing
- [FONT.md](FONT.md) — TTF on screen (separate from **`IMAGE.DRAWTEXT`**)
