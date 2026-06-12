# Texture Commands

Commands for loading, managing, and manipulating 2D textures. Textures are GPU-resident images used for sprites, billboards, model materials, UI elements, and render targets.

## Core Concepts

- **Texture** — An image loaded onto the GPU for fast rendering. Created from files (`.png`, `.jpg`, `.bmp`, `.tga`, etc.) or from in-memory images.
- **Handle** — Every texture returns a handle that must be freed when no longer needed to avoid GPU memory leaks.
- **Filters** — Textures support filter presets applied at load time (nearest-neighbor for pixel art, bilinear for smooth scaling).
- **Async loading** — Textures can be loaded in the background while showing a loading screen.
- Textures have a **finalizer** set at load time — if you forget to free a texture, the Go garbage collector will release the GPU resource eventually, but explicit `Free` is preferred.

---

## Loading

### `Texture.Load(filePath)`

Loads a texture from an image file. Supports `.png`, `.jpg`, `.bmp`, `.tga`, `.gif`, `.hdr`, `.psd`, `.dds`, `.ktx`, `.astc`.

- `filePath` (string) — Path to the image file.

**Returns:** `handle`

**How it works:** Reads the file, decodes it, uploads pixel data to the GPU, applies default texture filtering, and returns a heap handle. The `TextureObject` struct stores the Raylib texture, source path, filter flags, and UV scale.

```basic
playerTex = Texture.Load("assets/player.png")
ASSERT(playerTex <> 0, "Failed to load player texture")
```

---

### `Texture.LoadAsync(filePath)`

Begins loading a texture in the background. Use `Texture.IsLoaded` to poll completion.

- `filePath` (string) — Path to the image file.

**Returns:** `handle` — A pending texture handle.

**How it works:** Queues the file for background decoding. The handle is valid immediately but the texture data isn't usable until `Texture.IsLoaded` returns `TRUE`. Useful for loading screens.

```basic
tex = Texture.LoadAsync("assets/large_world.png")
WHILE NOT Texture.IsLoaded(tex)
    Render.Clear(0, 0, 0)
    Draw.Text("Loading...", 10, 10, 20, 255, 255, 255, 255)
    Render.Frame()
WEND
```

---

### `Texture.IsLoaded(textureHandle)`

Returns `TRUE` if an asynchronously loaded texture has finished loading and is ready to use.

- `textureHandle` (handle) — Texture handle.

**Returns:** `bool`

---

### `Texture.FromImage(imageHandle)` / `Image.ToTexture(imageHandle)`

Creates a GPU texture from a CPU-side image handle. Use this after manipulating an image in memory.

- `imageHandle` (handle) — Image handle (from `IMAGE.LOAD` or generated).

**Returns:** `handle`

---

### `Texture.Reload(textureHandle)`

Reloads a texture from its original source file. Useful for hot-reloading assets during development.

- `textureHandle` (handle) — Texture to reload.

---

### `Texture.Free(textureHandle)`

Frees a texture from GPU memory.

- `textureHandle` (handle) — Texture to free.

**How it works:** Calls `UnloadTexture` on the underlying Raylib texture and removes the handle from the heap.

```basic
Texture.Free(playerTex)
```

---

## Texture Drawing

Textures can be drawn using the `Draw.Billboard` family (3D) or through Sprite objects (2D). For direct 2D texture rendering, use `DrawTex2`:

```basic
; Create a texture draw object
tex = Texture.Load("assets/icon.png")
icon = DrawTex2(tex, 100, 100, 255, 255, 255, 255)

; Each frame
DrawTex2.Draw(icon)

; Cleanup
DrawTex2.Free(icon)
Texture.Free(tex)
```

---

## Advanced Texture Drawing

Textures loaded with `Texture.Load` can also be drawn with advanced wrappers that support rotation, scaling, source rectangles (for texture atlases), and tinting. These commands are registered in the `draw/texture_adv_wrappers_cgo.go` module.

---

## Easy Mode Shortcuts

| Shortcut | Maps To |
|----------|---------|
| `LoadTexture(path)` | `Texture.Load(path)` |
| `LOADTEXTURE(path)` | `Texture.Load(path)` |
| `FreeTexture(h)` | `Texture.Free(h)` |
| `FREETEXTURE(h)` | `Texture.Free(h)` |

---

## Full Example

Loading textures, drawing them as billboards in 3D, and cleaning up properly.

```basic
Window.Open(1280, 720, "Texture Demo")
Window.SetFPS(60)

cam = Camera.Create()
cam.pos(0, 5, 10)
cam.look(0, 1, 0)
cam.fov(60)

; Load textures
treeTex = Texture.Load("assets/tree.png")
grassTex = Texture.Load("assets/grass.png")

WHILE NOT Window.ShouldClose()
    Render.Clear(100, 180, 255)

    Camera.Begin(cam)
        Draw.Grid(20, 1.0)

        ; Draw trees as billboards
        FOR i = 0 TO 9
            x = (i - 5) * 3
            Draw.Billboard(cam, treeTex, x, 2, -5, 4.0, 255, 255, 255, 255)
        NEXT

    Camera.End(cam)

    Draw.Text("Texture Billboards", 10, 10, 24, 255, 255, 255, 255)
    Render.Frame()
WEND

; Always free textures
Texture.Free(treeTex)
Texture.Free(grassTex)
Camera.Free(cam)
Window.Close()
```

---

## See Also

- [DRAW](DRAW.md) — Billboard and 2D texture drawing
- [SPRITE](SPRITE.md) — 2D sprite system built on textures
- [MODEL](MODEL.md) — Textures applied to 3D model materials
- [MATERIAL](MATERIAL.md) — Material texture assignment
