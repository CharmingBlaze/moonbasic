# RenderTarget Commands

Off-screen render targets (framebuffers): render a scene to a texture, then use that texture on a mesh or as a sprite for mirrors, portals, mini-maps, and UI panels.

## Core Workflow

1. `RENDERTARGET.CREATE(width, height)` — allocate a framebuffer.
2. `RENDERTARGET.BEGIN(rt)` — redirect drawing into the framebuffer.
3. Draw your scene normally.
4. `RENDERTARGET.END()` — restore default framebuffer.
5. `RENDERTARGET.TEXTURE(rt)` — get the texture handle, use with `MATERIAL.SETTEXTURE` or `DRAW.TEXTURE`.
6. `RENDERTARGET.FREE(rt)` when done.

---

## Creation

### `RENDERTARGET.CREATE(width, height)` 

Allocates an off-screen framebuffer of `width × height` pixels. Returns a **rendertarget handle**.

---

## Render Pair

### `RENDERTARGET.BEGIN(rt)` 

Redirects all subsequent draw calls into the render target's framebuffer.

---

### `RENDERTARGET.END()` 

Ends render-target drawing and restores the default window framebuffer. Always pair with `RENDERTARGET.BEGIN`.

---

## Texture Access

### `RENDERTARGET.TEXTURE(rt)` 

Returns the render target's color buffer as a **texture handle**. Use with `MATERIAL.SETTEXTURE`, `DRAW.TEXTURE`, or `SPRITE.*`.

---

## Lifetime

### `RENDERTARGET.FREE(rt)` 

Frees the render target and its GPU texture.

---

## Full Example

Mini-map rendered to a texture displayed in the corner of the screen.

```basic
WINDOW.OPEN(960, 540, "RenderTarget Demo")
WINDOW.SETFPS(60)

; main camera
cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 5, -10)
CAMERA.SETTARGET(cam, 0, 0, 0)

; top-down mini-map camera
mmCam = CAMERA.CREATE()
CAMERA.SETPOS(mmCam, 0, 30, 0)
CAMERA.SETTARGET(mmCam, 0, 0, 0)

; off-screen target for the mini-map
rt = RENDERTARGET.CREATE(256, 256)

cube = ENTITY.CREATECUBE(2.0)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    ENTITY.TURN(cube, 0, 30 * dt, 0)
    ENTITY.UPDATE(dt)

    ; render mini-map into rt
    RENDERTARGET.BEGIN(rt)
        RENDER.CLEAR(10, 20, 10)
        RENDER.BEGIN3D(mmCam)
            ENTITY.DRAWALL()
            DRAW3D.GRID(20, 1.0)
        RENDER.END3D()
    RENDERTARGET.END()

    ; main scene
    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW3D.GRID(20, 1.0)
    RENDER.END3D()

    ; draw mini-map overlay
    mmTex = RENDERTARGET.TEXTURE(rt)
    DRAW.TEXTURE(mmTex, 700, 10, 250, 250, 255, 255, 255, 220)

    RENDER.FRAME()
WEND

RENDERTARGET.FREE(rt)
ENTITY.FREE(cube)
CAMERA.FREE(cam)
CAMERA.FREE(mmCam)
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `RENDERTARGET.MAKE(w, h)` | Deprecated alias of `RENDERTARGET.CREATE`. |

---

## See also

- [MATERIAL.md](MATERIAL.md) — apply rt texture to a mesh
- [TEXTURE.md](TEXTURE.md) — regular texture loading
- [CAMERA.md](CAMERA.md) — camera setup
- [SHADER.md](SHADER.md) — custom shader using rt texture
