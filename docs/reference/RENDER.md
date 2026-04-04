# Render Commands

Frame lifecycle and render-state helpers. Drawing itself uses **`Draw.*`**, **`Draw3D.*`**, **`Camera.Begin` / `Camera.End`**, and **`Camera2D.Begin` / `Camera2D.End`** — there are **no** `Render.BeginMode2D` / `Render.BeginMode3D` APIs.

## Frame loop

Each frame:

1. **`Render.Clear(...)`** — clear color / depth as configured.
2. **Optional modes:** `Camera2D.Begin` / `End` for transformed 2D; `Camera.Begin` / `Camera.End` for 3D.
3. **Draw** — `Draw.*`, `Draw3D.*`, models, etc.
4. **`Render.Frame()`** — present (swap buffers).

```basic
WHILE NOT Window.ShouldClose()
    Render.Clear(10, 20, 30)

    Camera2D.Begin()
        Draw.Rectangle(100, 100, 50, 50, 255, 255, 255, 255)
    Camera2D.End()

    Render.Frame()
WEND
```

---

## Core

### `Render.Clear([...])`

Overloads (see runtime):

- **0 args** — default clear.
- **1 arg** — color **handle**.
- **3 args** — `r, g, b` (0–255).
- **4 args** — `r, g, b, a`.

Requires an open window (`Window.Open`).

### `Render.Frame()`

Presents the frame. Must follow a `Render.Clear` for that frame (runtime enforces an active frame).

---

## Overlay / capture

### `Render.DrawFPS(x, y)`

Draws the FPS counter at integer pixel coordinates (registered with the render module).

### `Render.Screenshot(path$)`

Writes a PNG screenshot to `path$` (`TakeScreenshot`).

---

## Blending and depth

### `Render.SetBlend(mode)` / `Render.SetBlendMode(mode)`

Alias pair. `mode` is a numeric **Raylib blend mode** (e.g. `BLEND_ALPHA`, `BLEND_ADDITIVE` from key globals).

### `Render.SetDepthWrite(on)` / `Render.SetDepthMask(on)`

Alias pair. Enables or disables depth **mask** (writing to the depth buffer).

### `Render.SetDepthTest(on)`

Enables or disables depth **testing**.

---

## Raster state

### `Render.SetScissor(x, y, w, h)`

Enables scissor test and sets the rectangle (**integer** components).

### `Render.ClearScissor()`

Disables scissor test.

### `Render.SetWireframe(on)`

Enables or disables wireframe mode (Raylib wire mode).

---

## Window / pipeline hints

### `Render.SetMSAA(on)`

Toggles the **4× MSAA window hint** (affects setup; may require appropriate window flags).

### `Render.SetMode(mode$)`

`"forward"` or `"deferred"` — switches internal 3D render pipeline mode where supported.

### `Render.SetShadowMapSize(size)`

Sets shadow map resolution used by the 3D runtime (numeric size in pixels).

---

## 2D lighting ambient

When the light2d module is enabled (CGO build), ambient tint for the 2D lighting path:

### `Render.Set2DAmbient(r, g, b, a)`

**Four** integer components 0–255.

---

## Shader note

Global **`Render.BeginShader` / `Render.EndShader`** are **not** registered in the current runtime. Use **material/shader** APIs under **`Model.*`** / **`Shader.*`** (see [SHADER.md](SHADER.md)) or post-process paths where implemented.

---

## See also

- [CAMERA.md](CAMERA.md) — 2D/3D camera begin/end.
- [DRAW2D.md](DRAW2D.md), [DRAW3D.md](DRAW3D.md) — drawing commands.
