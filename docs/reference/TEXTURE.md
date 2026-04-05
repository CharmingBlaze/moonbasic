# Texture (`TEXTURE.*`)

GPU **texture** handles: images uploaded for fast `Draw.Texture*` calls. For **CPU** pixel work first, see [IMAGE.md](IMAGE.md), then **`TEXTURE.FROMIMAGE`** or **`TEXTURE.UPDATE`**.

**Registry keys** use uppercase `TEXTURE.` (e.g. `TEXTURE.LOAD`). In source, the **`Texture`** namespace maps to the same commands (`Texture.Load` → `TEXTURE.LOAD`).

**Threading:** Raylib GL calls belong on the **main thread** (see [ARCHITECTURE.md](../../ARCHITECTURE.md)).

---

## Loading and lifetime

| Command | Arguments | Returns | Notes |
|--------|-----------|---------|--------|
| **`TEXTURE.LOAD`** | `path$` | handle | `LoadTexture` from disk. |
| **`TEXTURE.FROMIMAGE`** | `image` | handle | From a heap **Image** (see [IMAGE.md](IMAGE.md)). |
| **`TEXTURE.FREE`** | `texture` | — | Unloads GPU data (`UnloadTexture`) unless the handle is **borrowed** (see render targets below). |

---

## Size, sampling, CPU upload

| Command | Arguments | Returns | Notes |
|--------|-----------|---------|--------|
| **`TEXTURE.WIDTH`** / **`TEXTURE.HEIGHT`** | `texture` | int | Pixel dimensions. |
| **`TEXTURE.SETFILTER`** | `texture`, `filter#` | — | Use globals **`FILTER_POINT`**, **`FILTER_BILINEAR`**, … (see below). |
| **`TEXTURE.SETWRAP`** | `texture`, `wrap#` | — | Use **`WRAP_REPEAT`**, **`WRAP_CLAMP`**, … |
| **`TEXTURE.UPDATE`** | `texture`, `image` | — | `UpdateTexture` from CPU **Image** pixels (format/size must match usage). |

Filter/wrap enum values are installed as globals by the runtime (same numeric values as Raylib):

- **Filter:** `FILTER_POINT`, `FILTER_BILINEAR`, `FILTER_TRILINEAR`, `FILTER_ANISOTROPIC_4X`, `FILTER_ANISOTROPIC_8X`, `FILTER_ANISOTROPIC_16X`
- **Wrap:** `WRAP_REPEAT`, `WRAP_CLAMP`, `WRAP_MIRROR_REPEAT`, `WRAP_MIRROR_CLAMP`

---

## Procedural textures

| Command | Arguments | Returns | Notes |
|--------|-----------|---------|--------|
| **`TEXTURE.GENWHITENOISE`** | `w`, `h` **or** `w`, `h`, `factor#` | handle | `GenImageWhiteNoise` → GPU texture. Default **factor** = **1**. |
| **`TEXTURE.GENCHECKED`** | `w`, `h`, `tileW`, `tileH`, `color1`, `color2` | handle | **Colors** are **COLOR** handles. |
| **`TEXTURE.GENGRADIENTV`** | `w`, `h`, `topColor`, `bottomColor` | handle | Vertical gradient. |
| **`TEXTURE.GENGRADIENTH`** | `w`, `h`, `leftColor`, `rightColor` | handle | Horizontal gradient (90° linear). |
| **`TEXTURE.GENCOLOR`** | `w`, `h`, `r`, `g`, `b`, `a` | handle | Flat RGBA fill. |

---

## Render targets (`RENDERTARGET.*`)

Off-screen rendering into a texture (FBO). Use for post-processing, minimaps, or multi-pass effects.

| Command | Arguments | Returns | Notes |
|--------|-----------|---------|--------|
| **`RENDERTARGET.MAKE`** | `width`, `height` | handle | `LoadRenderTexture`. Type name on handles: **`RenderTexture`**. |
| **`RENDERTARGET.FREE`** | `rt` | — | `UnloadRenderTexture`. |
| **`RENDERTARGET.BEGIN`** | `rt` | — | `BeginTextureMode` — draw and **`Render.Clear`** target this FBO. |
| **`RENDERTARGET.END`** | — | — | `EndTextureMode` — return to default framebuffer. |
| **`RENDERTARGET.TEXTURE`** | `rt` | handle | **Borrowed** **Texture** view of the **color** attachment. **`TEXTURE.FREE`** on this handle does **not** unload GPU memory; free the render target with **`RENDERTARGET.FREE`** (after dropping uses of the borrowed texture). |

Handle methods (on a **RenderTexture** handle): **`Begin`**, **`End`**, **`Free`**, **`Texture`** — same behaviour as the `RENDERTARGET.*` commands.

The color attachment is often **Y-flipped** vs screen space; use **`Draw.TexturePro`** / **`Draw.TextureRec`** with a negative source height, or draw helpers that account for UV orientation, when compositing to the screen.

---

## Drawing

Use **`Draw.Texture`**, **`Draw.TextureRec`**, **`Draw.TexturePro`**, etc. (see [DRAW2D.md](DRAW2D.md)). Manifest coverage may list only a subset; the runtime exposes the full Raylib-backed draw family where CGO is enabled.

---

## Atlas

Sprite sheets as a single GPU texture are documented in **[ATLAS.md](ATLAS.md)** (`ATLAS.*`).

---

## Example (load → draw → free)

```text
tex = TEXTURE.LOAD("assets/ui/panel.png")
WHILE NOT Window.ShouldClose()
    Render.Clear(30, 30, 40)
    Draw.Texture(tex, 10, 10, 255, 255, 255, 255)
    Render.Frame()
WEND
TEXTURE.FREE(tex)
```

---

## See also

- [IMAGE.md](IMAGE.md) — `IMAGE.MAKE`, `IMAGE.COPY`, export.
- [DRAW2D.md](DRAW2D.md) — `Draw.Texture*`, rectangles.
- [RENDER.md](RENDER.md) — `Render.Clear`, `Render.Frame`.
