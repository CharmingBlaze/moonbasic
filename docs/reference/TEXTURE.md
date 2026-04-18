# Texture Commands

GPU texture handles for loading, drawing, atlas animation, and render targets.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Load with `TEXTURE.LOAD` or create from a CPU image with `TEXTURE.FROMIMAGE`.
2. Draw with `DRAW.TEXTURE` / `DRAW.TEXTUREREC` / `DRAW.TEXTUREPRO` (see [DRAW2D.md](DRAW2D.md)).
3. For sprite-sheet animation use `TEXTURE.LOADANIM` + `TEXTURE.PLAY` + `TEXTURE.TICKALL`.
4. Free with `TEXTURE.FREE`.

For CPU pixel buffers see [IMAGE.md](IMAGE.md).

---

### `TEXTURE.LOAD(path)`
Loads a GPU texture from disk. Returns a **texture handle**.

- **Arguments**:
    - `path`: (String) File path relative to working directory.
- **Returns**: (Handle) The new texture handle.
- **Example**:
    ```basic
    tex = TEXTURE.LOAD("assets/grass.png")
    ```

---

### `TEXTURE.FROMIMAGE(imgHandle)`
Creates a GPU texture from an in-memory `Image` handle.

- **Arguments**:
    - `imgHandle`: (Handle) The CPU image source.
- **Returns**: (Handle) A new GPU texture.

---

### `TEXTURE.FREE(handle)`
Unloads GPU data and releases the handle from memory and its heap slot.

---

### `TEXTURE.WIDTH(handle)` / `TEXTURE.HEIGHT(handle)`
Returns the integer pixel dimensions of the texture.

---

### `TEXTURE.SETFILTER(handle, filter)`
Sets the sampling filter.

- **Arguments**:
    - `handle`: (Handle) The texture to modify.
    - `filter`: (Integer) Filter mode constant (e.g., `FILTER_POINT`).
- **Returns**: (Handle) The texture handle (for chaining).

---

### `RENDERTARGET.CREATE(w, h)`
Creates an off-screen render target (FBO). Returns a **handle**.

- **Arguments**:
    - `w, h`: (Integer) Dimensions in pixels.
- **Returns**: (Handle) The new render target handle.

---

### `RENDERTARGET.BEGIN(handle)`
Starts drawing into the specified render target.

---

### `RENDERTARGET.END()`
Ends drawing into the current target and returns to the default framebuffer.

---

### `RENDERTARGET.FREE(handle)`
Frees the render target and its associated color texture from memory.

The color attachment is often **Y-flipped** vs screen space; use **`DRAW.TEXTUREPRO`** / **`DRAW.TEXTUREREC`** with a negative source height, or draw helpers that account for UV orientation, when compositing to the screen.

---

## Drawing

Use **`DRAW.TEXTURE`**, **`DRAW.TEXTUREREC`**, **`DRAW.TEXTUREPRO`**, etc. (see [DRAW2D.md](DRAW2D.md)). Manifest coverage may list only a subset; the runtime exposes the full Raylib-backed draw family where CGO is enabled.

---

## Atlas

Sprite sheets as a single GPU texture are documented in **[ATLAS.md](ATLAS.md)** (`ATLAS.*` — JSON-packed rectangles).

### Uniform grid animation (`TEXTURE.SETGRID`, `TEXTURE.*`) 

For **equal-sized frames** laid out in a regular **columns × rows** grid on one texture (water ripples, fire strips, etc.):

| Command | Purpose |
|--------|---------|
| **`TEXTURE.SETGRID`** | `(texture, columns, rows)` — frame layout |
| **`TEXTURE.SETFRAME`** | `(texture, frameIndex)` — pick a cell (0-based) |
| **`TEXTURE.LOADANIM`** | `(path, columns, rows)` — load + set grid in one step |
| **`TEXTURE.PLAY`** | `(texture, fps, loop)` — auto-advance frames |
| **`TEXTURE.STOPANIM`** | Stop auto-advance |
| **`TEXTURE.TICKALL`** | Call **once per frame** (optional `dt`) so **`TEXTURE.PLAY`** advances |
| **`TEXTURE.SETUVSCROLL`** | `(texture, speedU, speedV)` — scroll source rect (for “infinite” flow) |
| **`TEXTURE.SETDISTORTION`** | `(texture, amount)` — hint for shader-side distortion |

**Billboards:** **`ENTITY.CREATESPRITE`** accepts **`(textureHandle, width, height [, parent])`** so a loaded grid/atlas applies to a 3D-facing quad; combine with **`TEXTURE.TICKALL`** and/or **`TEXTURE.SETFRAME`**. Full workflow (maps, animation modes): [**`SPRITE3D.md`**](SPRITE3D.md).

**Meshes:** **`ENTITY.SCROLLMATERIAL`** adds **(du, dv)** to material 0’s scroll (same idea as **`MODEL.SCROLLTEXTURE`**). **`ENTITY.SETDETAILTEXTURE`** binds a second texture as **normal/detail** for the same material.

---

## Full Example (load → draw → free)

```basic
WINDOW.OPEN(800, 600, "Texture draw")
WINDOW.SETFPS(60)
tex = TEXTURE.LOAD("assets/ui/panel.png")
WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(30, 30, 40)
    DRAW.TEXTURE(tex, 10, 10, 255, 255, 255, 255)
    RENDER.FRAME()
WEND
TEXTURE.FREE(tex)
WINDOW.CLOSE()
```

---

## Extended Command Reference

### Procedural generation

| Command | Description |
|--------|-------------|
| `TEXTURE.GENCOLOR(w, h, r, g, b, a)` | Generate a solid-color texture. |
| `TEXTURE.GENCHECKED(w, h, cx, cy, c1, c2)` | Generate a checkerboard texture. |
| `TEXTURE.GENGRADIENTH(w, h, left, right)` | Horizontal gradient texture. |
| `TEXTURE.GENGRADIENTV(w, h, top, bottom)` | Vertical gradient texture. |
| `TEXTURE.GENWHITENOISE(w, h, factor)` | White noise texture. |

### Queries

| Command | Description |
|--------|-------------|
| `TEXTURE.GETWIDTH(tex)` | Width in pixels. |
| `TEXTURE.GETHEIGHT(tex)` | Height in pixels. |
| `TEXTURE.GETSIZE(tex)` | Returns `[w, h]` array. |
| `TEXTURE.ISLOADED(tex)` | Returns `TRUE` if async load completed. |

### Settings

| Command | Description |
|--------|-------------|
| `TEXTURE.SETWRAP(tex, mode)` | Set UV wrap mode (`REPEAT`, `CLAMP`, `MIRROR`). |
| `TEXTURE.SETDEFAULTFILTER(mode)` | Set global default filter (`POINT`, `BILINEAR`, `TRILINEAR`). |
| `TEXTURE.UPDATE(tex, img)` | Upload updated CPU image to existing GPU texture. |
| `TEXTURE.RELOAD(tex, path)` | Reload texture from file in-place. |
| `TEXTURE.LOADASYNC(path, callback)` | Non-blocking load; fires `callback(handle)` when ready. |

---

## See also

- [IMAGE.md](IMAGE.md) — **`IMAGE.CREATE`** / deprecated **`IMAGE.MAKE`**, `IMAGE.COPY`, export.
- [DRAW2D.md](DRAW2D.md) — **`DRAW.TEXTURE*`**, rectangles.
- [RENDER.md](RENDER.md) — **`RENDER.CLEAR`**, **`RENDER.FRAME`**.
