# Image Commands (CPU images)

`IMAGE.*` works with **CPU-side** Raylib images (pixel buffers). Use them for procedural art, blitting, then `Texture.FromImage` to upload to the GPU.

**Requires CGO** and a working Raylib build.

---

## Load / query / lifecycle

See manifest entries: `IMAGE.LOAD`, `IMAGE.WIDTH`, `IMAGE.HEIGHT`, `IMAGE.FREE`, `IMAGE.MAKEBLANK`, `IMAGE.MAKECOPY`, `IMAGE.EXPORT`, etc.

---

## Drawing into an image

- `IMAGE.DRAWPIXEL`, `DRAWRECT`, `DRAWLINE`, `DRAWCIRCLE`, `DRAWTEXT`

---

## Compositing and processing

### `IMAGE.DrawImage(dest, src, sx#, sy#, sw#, sh#, dx#, dy#, dw#, dh#, r, g, b, a)`

Blits a rectangle from `src` into `dest` with a tint (Raylib `ImageDraw`). Both arguments are **image** handles.

### `IMAGE.Dither(img, rBpp, gBpp, bBpp, aBpp)`

Floyd–Steinberg dithering to lower bit depth.

### `IMAGE.Mipmaps(img)`

Generates mip levels in memory for the image.

### `IMAGE.Format(img, pixelFormatInt)`

Converts pixel format (Raylib `PixelFormat` enum as integer).

### `IMAGE.DrawRectLines(img, x#, y#, w#, h#, thick, r, g, b, a)`

Stroked rectangle on the image.

### `IMAGE.AlphaCrop(img, threshold#)`

Crops to opaque bounds using alpha threshold.

### `IMAGE.AlphaClear(img, r, g, b, a, threshold#)`

Clears pixels below alpha threshold to a solid color.

```basic
a = IMAGE.MakeBlank(64, 64, 0, 0, 0, 255)
b = IMAGE.MakeBlank(32, 32, 255, 0, 0, 255)
IMAGE.DrawImage(a, b, 0, 0, 32, 32, 16, 16, 32, 32, 255, 255, 255, 255)
tex = Texture.FromImage(a)
IMAGE.Free(b)
IMAGE.Free(a)
```
