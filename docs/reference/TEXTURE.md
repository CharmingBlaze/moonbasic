# Texture (`TEXTURE.*`)

GPU **texture** handles: images uploaded for fast `Draw.Texture*` calls. For **CPU** pixel work first, see [IMAGE.md](IMAGE.md), then **`TEXTURE.FROMIMAGE`** or **`TEXTURE.UPDATE`**.

**Registry keys** use uppercase `TEXTURE.` (e.g. `TEXTURE.LOAD`). In source, the **`Texture`** namespace maps to the same commands (`Texture.Load` → `TEXTURE.LOAD`).

**Threading:** Raylib GL calls belong on the **main thread** (see [ARCHITECTURE.md](../../ARCHITECTURE.md)).

---

### `Texture.Load(path)`
Loads a GPU texture from disk. Returns a **texture handle**.

### `Texture.FromImage(id)`
Creates a GPU texture from an in-memory `Image` handle.

### `Texture.Free(handle)`
Unloads GPU data and releases the handle from memory and its heap slot.

---

### `Texture.Width(handle)` / `Texture.Height(handle)`
Returns the integer pixel dimensions of the texture.

### `Texture.SetFilter(handle, filter)`
Sets the sampling filter (e.g., `FILTER_POINT`, `FILTER_BILINEAR`, `FILTER_TRILINEAR`).

---

### `RenderTarget.Make(w, h)`
Creates an off-screen render target (FBO). Returns a **handle** (`RenderTexture`).

### `RenderTarget.Begin(handle)`
Starts drawing into the specified render target. Subsequent `Draw.*` calls will target this FBO.

### `RenderTarget.End()`
Ends drawing into the current target and returns to the default framebuffer.

### `RenderTarget.Free(handle)`
Frees the render target and its associated color texture from memory.

The color attachment is often **Y-flipped** vs screen space; use **`Draw.TexturePro`** / **`Draw.TextureRec`** with a negative source height, or draw helpers that account for UV orientation, when compositing to the screen.

---

## Drawing

Use **`Draw.Texture`**, **`Draw.TextureRec`**, **`Draw.TexturePro`**, etc. (see [DRAW2D.md](DRAW2D.md)). Manifest coverage may list only a subset; the runtime exposes the full Raylib-backed draw family where CGO is enabled.

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

**Billboards:** **`ENTITY.CREATESPRITE`** accepts **`(textureHandle, width, height [, parent])`** so a loaded atlas applies to a 3D-facing quad; combine with **`TEXTURE.TICKALL`** and/or **`TEXTURE.SETFRAME`**.

**Meshes:** **`ENTITY.SCROLLMATERIAL`** adds **(du, dv)** to material 0’s scroll (same idea as **`MODEL.SCROLLTEXTURE`**). **`ENTITY.SETDETAILTEXTURE`** binds a second texture as **normal/detail** for the same material.

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
