# Image (`IMAGE.*`)

**CPU-side** pixel buffers (Raylib `Image`): read, mutate, save. **Not** GPU textures (`TEXTURE.*`). Typical pipeline: **`IMAGE.CREATE`** (canonical) or deprecated **`IMAGE.MAKE`** / **`IMAGE.LOAD`** → optional edits → **`TEXTURE.FROMIMAGE`** → **`DRAW.TEXTURE`** on the main framebuffer → free when done.

**Conventions:** [STYLE_GUIDE.md](../../STYLE_GUIDE.md), [API_CONVENTIONS.md](API_CONVENTIONS.md) — reference pages use uppercase **`NAMESPACE.ACTION`**; Easy Mode (`Image.Load`, …) is [compatibility only](../../STYLE_GUIDE.md#easy-mode-compatibility-layer).

**Page shape:** [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) — see [WAVE.md](WAVE.md) (registry-first headings, **Full Example** at the end).

**Requires CGO** and Raylib (same as **`DRAW.*`**, **`TEXTURE.*`**).

Registry keys use **dots and uppercase** (e.g. **`IMAGE.CREATE`**). In source, the **`Image`** namespace maps to the same commands (`Image.Load` → `IMAGE.LOAD`).

---

### `IMAGE.LOAD(path)`
Loads an image from disk (PNG, JPG, BMP, etc.). Returns a **handle**.

### `IMAGE.CREATE(w, h [, r, g, b, a])`
Creates a new CPU image. If RGBA components are provided, fills the image with that color (0-255). **`IMAGE.MAKE`** is a **deprecated** alias of **`IMAGE.CREATE`**.

### `IMAGE.FREE(handle)`
Releases the heap slot and unloads the image memory.

---

### `IMAGE.WIDTH(handle)` / `IMAGE.HEIGHT(handle)`
Returns the integer pixel dimensions of the image.

### `IMAGE.RESIZE(handle, w, h)`
Resizes the image in memory using bilinear scaling.

### `IMAGE.EXPORT(handle, path)`
Saves the image to a file. The format is determined by the file extension.

---

### `IMAGE.DRAWPIXEL(handle, x, y, r, g, b, a)`
Draws a single pixel on the image at `(x, y)`.

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
