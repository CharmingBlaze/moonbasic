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

## Extended Command Reference

### Creation aliases

| Command | Description |
|--------|-------------|
| `IMAGE.MAKE(w, h)` / `IMAGE.MAKEBLANK(w, h)` | Create a blank RGBA image. Aliases of `IMAGE.CREATEBLANK`. |
| `IMAGE.MAKECOPY(img)` | Alias of `IMAGE.CREATECOPY`. |
| `IMAGE.MAKETEXT(text, font, size, r, g, b)` | Alias of `IMAGE.CREATETEXT`. |
| `IMAGE.CREATEBLANK(w, h)` | Create blank RGBA image. |
| `IMAGE.CREATECOPY(img)` | Deep-copy an image handle. |
| `IMAGE.CREATETEXT(text, font, size, r, g, b)` | Render text into an image. |
| `IMAGE.COPY(img)` | Alias of `IMAGE.CREATECOPY`. |
| `IMAGE.LOADGIF(path)` | Load an animated GIF; returns first frame. |
| `IMAGE.LOADRAW(path, w, h, fmt)` | Load raw pixel data from file. |
| `IMAGE.LOADSEQUENCE(pattern, count)` | Load numbered image sequence (`frame_001.png`, …). |
| `IMAGE.TOTEXTURE(img)` | Upload image to GPU and return a texture handle. |

### Size queries

| Command | Description |
|--------|-------------|
| `IMAGE.GETWIDTH(img)` | Width in pixels. |
| `IMAGE.GETHEIGHT(img)` | Height in pixels. |
| `IMAGE.GETSIZE(img)` | Returns `[w, h]` array. |
| `IMAGE.GETBBOXX(img)` / `GETBBOXY(img)` | Bounding box top-left after alpha crop. |
| `IMAGE.GETBBOXW(img)` / `GETBBOXH(img)` | Bounding box width/height after alpha crop. |

### Pixel access

| Command | Description |
|--------|-------------|
| `IMAGE.GETPIXEL(img, x, y)` | Returns `[r,g,b,a]` of pixel at `(x,y)`. |
| `IMAGE.GETCOLORR(img, x, y)` | Red channel 0–255. |
| `IMAGE.GETCOLORG(img, x, y)` | Green channel 0–255. |
| `IMAGE.GETCOLORB(img, x, y)` | Blue channel 0–255. |
| `IMAGE.GETCOLORA(img, x, y)` | Alpha channel 0–255. |

### Transforms

| Command | Description |
|--------|-------------|
| `IMAGE.FLIPH(img)` | Flip horizontally in place. |
| `IMAGE.FLIPV(img)` | Flip vertically in place. |
| `IMAGE.ROTATE(img, degrees)` | Rotate by arbitrary angle (new image). |
| `IMAGE.ROTATECW(img)` | Rotate 90° clockwise. |
| `IMAGE.ROTATECCW(img)` | Rotate 90° counter-clockwise. |
| `IMAGE.CROP(img, x, y, w, h)` | Crop to rectangle in place. |
| `IMAGE.ALPHACROP(img)` | Crop to non-transparent bounding box. |
| `IMAGE.RESIZENN(img, w, h)` | Resize with nearest-neighbour filtering. |
| `IMAGE.MIPMAPS(img)` | Generate mipmap chain in image. |
| `IMAGE.DITHER(img)` | Apply dithering. |
| `IMAGE.SETFILTER(img, mode)` | Set scaling filter (`POINT`, `BILINEAR`). |

### Color adjustments

| Command | Description |
|--------|-------------|
| `IMAGE.COLORBRIGHTNESS(img, factor)` | Adjust brightness. |
| `IMAGE.COLORCONTRAST(img, factor)` | Adjust contrast. |
| `IMAGE.COLORGRAYSCALE(img)` | Convert to grayscale. |
| `IMAGE.COLORINVERT(img)` | Invert colors. |
| `IMAGE.COLORTINT(img, r, g, b, a)` | Multiply tint each pixel. |
| `IMAGE.COLORREPLACE(img, sr,sg,sb,sa, dr,dg,db,da)` | Replace one color with another. |
| `IMAGE.ALPHACLEAR(img, threshold)` | Zero alpha for pixels below threshold. |
| `IMAGE.CLEARBACKGROUND(img, r, g, b, a)` | Fill image with solid color. |

### Drawing into image

| Command | Description |
|--------|-------------|
| `IMAGE.DRAWLINE(img, x0, y0, x1, y1, r, g, b, a)` | Draw a line. |
| `IMAGE.DRAWCIRCLE(img, cx, cy, radius, r, g, b, a)` | Draw a filled circle. |
| `IMAGE.DRAWRECTLINES(img, x, y, w, h, r, g, b, a)` | Draw a rectangle outline. |

---

## See also

- [TEXTURE.md](TEXTURE.md) — **`TEXTURE.FROMIMAGE`**, render targets
- [DRAW2D.md](DRAW2D.md) — screen-space drawing
- [FONT.md](FONT.md) — TTF on screen (separate from **`IMAGE.DRAWTEXT`**)
