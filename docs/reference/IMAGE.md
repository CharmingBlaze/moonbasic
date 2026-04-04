# Image (CPU) — `Image.*` / `IMAGE.*`

**`Image.*`** commands work with **CPU-side** Raylib images: width × height pixel buffers you can read, write, and process in memory. They are **not** the same as **GPU textures** (`Texture.*`). Typical flow: build or load an image → **`Texture.FromImage`** → draw with **`Draw.Texture`**.

**Requires CGO** and a Raylib build (same as `Draw.*`, `Texture.*`).

---

## Load, create, export, free

| Command | Arguments | Returns | Notes |
|--------|-----------|---------|--------|
| `Image.Load(path$)` | file path | handle | `LoadImage` |
| `Image.LoadRaw(path$, w, h, format, headerSize)` | raw pixel file | handle | Raylib `PixelFormat` as `format` int |
| `Image.MakeBlank(w, h, r, g, b, a)` | size + clear color | handle | Filled rectangle |
| `Image.MakeCopy(img)` | image handle | handle | Deep copy |
| `Image.MakeText(text$, fontSize, r, g, b, a)` | string + style | handle | Transparent canvas; draws text with default font (fixed layout heuristic) |
| `Image.Export(img, path$)` | handle + path | bool | Writes file (format from extension) |
| `Image.Width(img)` / `Image.Height(img)` | handle | int | |
| `Image.Free(img)` | handle | — | Releases heap object and Raylib image |

---

## Crop, resize, flip, rotate

All of these **mutate** the image in place (Raylib-style).

| Command | Arguments | Notes |
|--------|-----------|--------|
| `Image.Crop(img, x, y, w, h)` | int rect | |
| `Image.Resize(img, newW, newH)` | int | Smooth filter |
| `Image.ResizeNN(img, newW, newH)` | int | Nearest-neighbor |
| `Image.FlipH(img)` / `Image.FlipV(img)` | handle | |
| `Image.Rotate(img, degrees)` | handle, int | Custom angle |
| `Image.RotateCW(img)` / `Image.RotateCCW(img)` | handle | 90° steps |

---

## Color adjustments

| Command | Arguments | Notes |
|--------|-----------|--------|
| `Image.ColorTint(img, r, g, b, a)` | handle + RGBA 0–255 | Multiply tint |
| `Image.ColorInvert(img)` | handle | |
| `Image.ColorGrayscale(img)` | handle | |
| `Image.ColorContrast(img, contrast#)` | handle, float | |
| `Image.ColorBrightness(img, brightness)` | handle, int | |
| `Image.ColorReplace(img, r1,g1,b1,a1, r2,g2,b2,a2)` | 8 ints | Replace one exact color with another |
| `Image.ClearBackground(img, r, g, b, a)` | handle + RGBA | Fills entire image |

---

## Drawing **onto** an image

Pixel coordinates are **integers**; colors are **RGBA 0–255**.

| Command | Arguments |
|--------|-----------|
| `Image.DrawPixel(img, x, y, r, g, b, a)` | |
| `Image.DrawRect(img, x, y, w, h, r, g, b, a)` | filled rect |
| `Image.DrawLine(img, x1, y1, x2, y2, r, g, b, a)` | |
| `Image.DrawCircle(img, cx, cy, radius, r, g, b, a)` | filled |
| `Image.DrawText(img, x, y, text$, fontSize, r, g, b, a)` | default font |
| `Image.DrawRectLines(img, x#, y#, w#, h#, thick, r, g, b, a)` | stroked rect (`float` position/size) |

---

## Blit, dither, mipmaps, format, alpha tools

| Command | Arguments | Notes |
|--------|-----------|--------|
| `Image.DrawImage(dest, src, sx, sy, sw, sh, dx, dy, dw, dh, r, g, b, a)` | source rect → dest rect + tint | `ImageDraw`; both handles must be images |
| `Image.Dither(img, rBpp, gBpp, bBpp, aBpp)` | handle + 4 ints | Floyd–Steinberg |
| `Image.Mipmaps(img)` | handle | Builds mip chain in CPU memory |
| `Image.Format(img, pixelFormatInt)` | handle, int | `PixelFormat` enum value |
| `Image.AlphaCrop(img, threshold#)` | handle, float | Crop to opaque bounds |
| `Image.AlphaClear(img, r, g, b, a, threshold#)` | | Pixels below alpha threshold → solid color |

---

## Read pixels & alpha bounding box

| Command | Arguments | Returns |
|--------|-----------|---------|
| `Image.GetColorR` / `GetColorG` / `GetColorB` / `GetColorA` | `(img, x, y)` | int 0–255 per channel |
| `Image.GetBBOXX` / `GetBBOXY` / `GetBBOXW` / `GetBBOXH` | `(img, alphaThreshold#)` | int tight bbox (`GetBBOXX` ends in **two X**s — matches `IMAGE.GETBBOXX`) |

---

## Clipboard → image

| Command | Arguments | Returns |
|--------|-----------|---------|
| `Clipboard.GetImage()` | (none) | handle |

Returns a **new** image copied from the OS clipboard, or fails if no image is available. Free with **`Image.Free`** when done.

---

## Example: composite, upload, draw

```basic
IF NOT Window.Open(640, 480, "Image → Texture") THEN END
ENDIF
Window.SetFPS(60)

a = Image.MakeBlank(128, 128, 40, 40, 50, 255)
b = Image.MakeBlank(32, 32, 200, 80, 80, 255)
Image.DrawImage(a, b, 0, 0, 32, 32, 48, 48, 32, 32, 255, 255, 255, 255)
Image.Free(b)

tex = Texture.FromImage(a)
Image.Free(a)

WHILE NOT Window.ShouldClose()
    Render.Clear(20, 24, 32)
    Draw.Texture(tex, 200, 160)
    Render.Frame()
WEND

Texture.Free(tex)
Window.Close()
```

---

## See also

- [TEXTURE.md](TEXTURE.md) — GPU textures, **`Texture.FromImage`**
- [DRAW2D.md](DRAW2D.md) — screen drawing
- [FONT.md](FONT.md) — TTF on-screen text (separate from **`Image.DrawText`**)
