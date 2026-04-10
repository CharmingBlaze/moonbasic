# Image (CPU) — `Image.*` / `IMAGE.*`

**CPU-side** pixel buffers (Raylib `Image`): read, mutate, save. **Not** GPU textures (`Texture.*`). Typical pipeline: **`IMAGE.MAKE`** / **`IMAGE.LOAD`** → optional edits → **`TEXTURE.FROMIMAGE`** → **`DRAW.TEXTURE`** on the main framebuffer → free when done.

**Requires CGO** and Raylib (same as `Draw.*`, `Texture.*`).

Registry keys use **dots and uppercase** (e.g. `IMAGE.MAKE`). PascalCase names below match docs/spec style.

---

### Image.Make

```basic
h = IMAGE.MAKE(w, h)
h = IMAGE.MAKE(w, h, r, g, b, a)
```

**Two-argument** form creates an **RGBA image filled with transparent black** `(0,0,0,0)`. **Six-argument** form fills with the given **RGBA** (each channel **0–255**).

**Parameters**

| Name | Type | Description |
|---|---|---|
| w, h | int | Width and height in pixels. |
| r, g, b, a | int | Optional. Fill colour; required together with six-arg form. |

**Returns** — handle.

> **Common mistake:** Expecting **`IMAGE.MAKE(w,h)`** to be opaque black — it is **fully transparent** until you **`IMAGE.CLEAR`** or draw.

**Example**

```basic
img = IMAGE.MAKE(128, 128)
IMAGE.CLEAR(img, 40, 44, 52, 255)
```

**See also:** `IMAGE.MAKEBLANK`, `IMAGE.CLEAR`, `TEXTURE.FROMIMAGE`

---

### Image.MakeBlank

```basic
h = IMAGE.MAKEBLANK(w, h)
h = IMAGE.MAKEBLANK(w, h, r, g, b, a)
```

Same behaviour as **`IMAGE.MAKE`** (alias pair). Use either name.

---

### Image.Load

```basic
h = IMAGE.LOAD(path$)
```

Loads **PNG, JPG, BMP, TGA, GIF, HDR**, etc. from disk (Raylib). File is read and closed; you receive a **new** image handle.

**Parameters**

| Name | Type | Description |
|---|---|---|
| path | string | File path. |

**Returns** — handle.

---

### Image.LoadRaw

```basic
h = IMAGE.LOADRAW(path$, w, h, format, headerSize)
```

Loads raw pixel data. **`format`** is a Raylib **`PixelFormat`** integer; **`headerSize`** skips bytes at the start of the file.

**See also:** Raylib `PixelFormat` constants.

---

### Image.Copy / Image.MakeCopy

```basic
h2 = IMAGE.COPY(h)
h2 = IMAGE.MAKECOPY(h)
```

**Deep copy** — new handle, duplicated pixels. **`IMAGE.COPY`** and **`IMAGE.MAKECOPY`** are equivalent.

**Parameters**

| Name | Type | Description |
|---|---|---|
| h | handle | Source image. |

**Returns** — handle.

---

### Image.Free

```basic
IMAGE.FREE(h)
```

Releases the heap slot and Raylib image. A **second** **`IMAGE.FREE`** with the same handle value fails (stale handle), same as other heap objects.

---

### Image.Clear / Image.ClearBackground

```basic
IMAGE.CLEAR(h, r, g, b, a)
IMAGE.CLEARBACKGROUND(h, r, g, b, a)
```

Fills **every pixel** with **RGBA** (0–255 per channel). **`IMAGE.CLEAR`** and **`IMAGE.CLEARBACKGROUND`** are equivalent.

---

### Image.Width / Image.Height

```basic
w = IMAGE.WIDTH(h)
h = IMAGE.HEIGHT(h)
```

**Returns** — integer dimensions.

---

### Image.Export

```basic
ok = IMAGE.EXPORT(h, path$)
```

Writes an image file; format from **extension** (e.g. `.png`, `.jpg`, `.bmp`, `.tga`). **Returns** boolean success.

---

## Pixel reads (components)

Use separate channel queries (no array value in core **`IMAGE.*`**):

| Registry | Arguments | Returns |
|---|---|---|
| `IMAGE.GETCOLORR` / `G` / `B` / `A` | `(h, x, y)` | int 0–255 |

**Coordinates** are **integer** pixel indices.

---

## Drawing on an image

| Command | Purpose |
|---|---|
| `IMAGE.DRAWPIXEL` | Single pixel |
| `IMAGE.DRAWRECT` | Filled rectangle |
| `IMAGE.DRAWLINE` | Line |
| `IMAGE.DRAWCIRCLE` | Filled circle |
| `IMAGE.DRAWRECTLINES` | Rectangle outline (`float` geometry + thickness) |
| `IMAGE.DRAWTEXT` | Default font text: `(h, x, y, text, fontSize, r, g, b, a)` |
| `IMAGE.DRAWIMAGE` | Blit: source rect → dest rect + tint |

---

## Crop, resize, flip, rotate

All **mutate** in place.

| Command | Notes |
|---|---|
| `IMAGE.CROP` | Rectangle crop |
| `IMAGE.RESIZE` | Bilinear |
| `IMAGE.RESIZENN` | Nearest-neighbour |
| `IMAGE.FLIPH` / `IMAGE.FLIPV` | |
| `IMAGE.ROTATE` | Arbitrary degrees |
| `IMAGE.ROTATECW` / `IMAGE.ROTATECCW` | 90° steps |

---

## Colour adjustments

| Command | Notes |
|---|---|
| `IMAGE.COLORTINT` | Multiply tint RGBA |
| `IMAGE.COLORINVERT` | |
| `IMAGE.COLORGRAYSCALE` | |
| `IMAGE.COLORCONTRAST` | Numeric contrast |
| `IMAGE.COLORBRIGHTNESS` | Offset brightness |
| `IMAGE.COLORREPLACE` | Eight ints: from RGBA → to RGBA |

---

## Blit, mipmaps, format, alpha tools

| Command | Notes |
|---|---|
| `IMAGE.DITHER` | Floyd–Steinberg; four bpp values |
| `IMAGE.MIPMAPS` | CPU mipmap chain |
| `IMAGE.FORMAT` | Convert pixel format (Raylib enum int) |
| `IMAGE.ALPHACROP` | Crop to non-transparent bounds |
| `IMAGE.ALPHACLEAR` | Below alpha threshold → solid colour |

---

## Alpha bounding box

| Command | Returns |
|---|---|
| `IMAGE.GETBBOXX`, `GETBBOXY`, `GETBBOXW`, `GETBBOXH` | `(h, alphaThreshold)` → int |

---

## Clipboard

| Command | Returns |
|---|---|
| `CLIPBOARD.GETIMAGE` | New image handle from OS clipboard (if available) |

---

## Example: composite → texture → draw

```basic
Window.Open(640, 480, "Image to texture")
Window.SetFPS(60)

a = IMAGE.MAKE(128, 128)
IMAGE.CLEAR(a, 40, 40, 50, 255)
b = IMAGE.MAKE(32, 32, 200, 80, 80, 255)
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
- **`IMAGE.MAKE(w,h)`** — Transparent, not black-opaque, until you clear or paint.

---

## See also

- [TEXTURE.md](TEXTURE.md) — **`TEXTURE.FROMIMAGE`**, render targets
- [DRAW2D.md](DRAW2D.md) — screen-space drawing
- [FONT.md](FONT.md) — TTF on screen (separate from **`IMAGE.DRAWTEXT`**)
