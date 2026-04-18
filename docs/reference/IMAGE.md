# Image Commands

CPU-side pixel buffers: create, load, mutate, export, then upload to GPU via `TEXTURE.FROMIMAGE`.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Load with `IMAGE.LOAD` or create with `IMAGE.CREATE`.
2. Optionally draw on the image with `IMAGE.DRAWPIXEL`, `IMAGE.DRAWRECT`, etc.
3. Upload to GPU with `TEXTURE.FROMIMAGE` (see [TEXTURE.md](TEXTURE.md)).
4. Free CPU image with `IMAGE.FREE`.

---

### `IMAGE.LOAD(path)`
Loads an image from disk (PNG, JPG, BMP, etc.). Returns a **handle**.

- **Arguments**:
    - `path`: (String) File path relative to working directory.
- **Returns**: (Handle) The new image handle.
- **Example**:
    ```basic
    img = IMAGE.LOAD("hero.png")
    ```

---

### `IMAGE.CREATE(w, h [, r, g, b, a])`
Creates a new CPU image.

- **Arguments**:
    - `w, h`: (Integer) Dimensions in pixels.
    - `r, g, b, a`: (Integer, optional) Initial fill color (0-255).
- **Returns**: (Handle) The new image handle.
- **Example**:
    ```basic
    a = IMAGE.CREATE(128, 128, 255, 0, 0, 255) ; 128x128 Red box
    ```

---

### `IMAGE.FREE(handle)`
Releases the heap slot and unloads the image memory.

---

### `IMAGE.WIDTH(handle)` / `IMAGE.HEIGHT(handle)`
Returns the integer pixel dimensions of the image.

---

### `IMAGE.RESIZE(handle, w, h)`
Resizes the image in memory using bilinear scaling.

- **Returns**: (Handle) The modified image handle (for chaining).

---

### `IMAGE.EXPORT(handle, path)`
Saves the image to a file. The format is determined by the file extension.

---

### `IMAGE.DRAWPIXEL(handle, x, y, r, g, b, a)`
Draws a single pixel on the image.

---

### `IMAGE.DRAWRECT(handle, x, y, w, h, r, g, b, a)`
Draws a filled rectangle on the image.

---

## Full Example (composite → texture → draw)

```basic
WINDOW.OPEN(640, 480, "Image to texture")
WINDOW.SETFPS(60)

a = IMAGE.CREATE(128, 128)
IMAGE.CLEAR(a, 40, 40, 50, 255)
b = IMAGE.CREATE(32, 32, 200, 80, 80, 255)
IMAGE.DRAWIMAGE(a, b, 0, 0, 32, 32, 48, 48, 32, 32, 255, 255, 255, 255)
IMAGE.FREE(b)

tex = TEXTURE.FROMIMAGE(a)
IMAGE.FREE(a)

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(20, 24, 32)
    DRAW.TEXTURE(tex, 200, 160, 255, 255, 255, 255)
    RENDER.FRAME()
WEND

TEXTURE.FREE(tex)
WINDOW.CLOSE()
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
