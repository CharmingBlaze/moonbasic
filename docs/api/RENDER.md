# Render Commands

Commands that control the frame lifecycle, screen clearing, 3D camera passes, and advanced rendering features like blend modes, post-processing, wireframe, and screenshots.

## Core Workflow

Every frame in moonBASIC follows this pattern:

```basic
Render.Clear(r, g, b)     ; Begin frame + clear screen
; ... draw commands ...
Render.Frame()             ; Present frame to screen
```

`RENDER.CLEAR` implicitly calls `BeginDrawing()` if a frame isn't already active. `RENDER.FRAME` calls `EndDrawing()`, presents the back buffer, and advances the frame counter. All draw commands must happen between these two calls.

For 3D rendering, wrap your 3D draw calls in a `RENDER.BEGIN3D` / `RENDER.END3D` pair (or use `CAMERA.BEGIN` / `CAMERA.END`).

---

## Frame Lifecycle

### `Render.Clear(r, g, b)` / `Render.Clear(r, g, b, a)` / `Render.Clear(colorHandle)` / `Render.Clear()`

Begins a new frame and clears the screen to the given color. This is the **first call** of every frame.

**Overloads:**
- `Render.Clear()` — Clear to black (0, 0, 0).
- `Render.Clear(r, g, b)` — Clear to an RGB color (0–255 each).
- `Render.Clear(r, g, b, a)` — Clear to an RGBA color.
- `Render.Clear(colorHandle)` — Clear using a `Color` handle created with `COLOR.CREATE`.

**How it works:** If no frame is currently active, calls `BeginDrawing()` first, then `ClearBackground()` with the requested color. This means you can safely call `Render.Clear` multiple times per frame (it won't start a new frame each time — it just re-clears).

**Errors:**
- `"RENDER.CLEAR: window is not open (call WINDOW.OPEN first)"`
- `"RENDER.CLEAR: r, g, b must be numeric"`

```basic
; Simple clear to dark blue
Render.Clear(20, 20, 60)

; Clear to black (shorthand)
Render.Clear()

; Clear with alpha
Render.Clear(0, 0, 0, 128)
```

---

### `Render.Frame()`

Ends the current frame, presents it to the screen, and advances the frame counter. This is the **last call** of every frame.

**How it works:**
1. Flushes any active render targets.
2. Runs frame-draw hooks (overlay systems, debug displays).
3. Runs the frame-end hook (if registered).
4. Increments `rt.FrameCount`.
5. If a script error occurred, draws it as red text in the top-left corner.
6. Draws the world-flash overlay if active (e.g., from `WORLD.FLASH`).
7. Calls `EndDrawing()` to present the back buffer.
8. Calls `Gosched()` to yield the Go scheduler.
9. Optionally logs FPS to diagnostics every 60 frames.

**Errors:**
- `"RENDER.FRAME: window is not open"`
- `"RENDER.FRAME: no active frame (call RENDER.CLEAR first)"`

```basic
Render.Clear(0, 0, 0)
Draw.Text("Hello!", 10, 10, 20, 255, 255, 255, 255)
Render.Frame()
```

---

## 3D Camera Pass

### `Render.Begin3D(cameraHandle)`

Begins a 3D rendering pass with the given camera. All 3D draw commands (meshes, entities, grid, lines, etc.) must happen between `Begin3D` and `End3D`.

- `cameraHandle` (handle) — A camera created with `CAMERA.CREATE`.

**How it works:** Calls Raylib's `BeginMode3D` with the camera's internal 3D configuration. The depth buffer is active during this pass.

```basic
cam = Camera.Create()
cam.pos(0, 10, 20)
cam.look(0, 0, 0)

Render.Clear(40, 40, 60)
Render.Begin3D(cam)
    Draw.Grid(20, 1.0)
    Draw.Cube(0, 1, 0, 2, 2, 2, 255, 0, 0, 255)
Render.End3D()
Render.Frame()
```

---

### `Render.End3D()`

Ends the current 3D rendering pass. Returns to 2D screen-space drawing.

---

## Render Dimensions

### `Render.Width()`

Returns the current render target width. This may differ from `Window.Width()` if using render-to-texture.

**Returns:** `int`

---

### `Render.Height()`

Returns the current render target height.

**Returns:** `int`

---

## Blend & Depth

### `Render.SetBlend(mode)` / `Render.SetBlendMode(mode)`

Sets the GPU blend mode for subsequent draw calls. Common modes: 0 = Alpha (default), 1 = Additive, 2 = Multiplied.

- `mode` (int) — Raylib blend mode constant.

**How it works:** Calls `BeginBlendMode` with the given mode. Affects all subsequent 2D/3D draw commands until changed again.

```basic
; Draw an additive glow
Render.SetBlend(1)
Draw.Circle(400, 300, 50, 255, 200, 100, 128)
Render.SetBlend(0) ; Back to normal alpha
```

---

### `Render.SetDepthWrite(enabled)` / `Render.SetDepthMask(enabled)`

Enables or disables writing to the depth buffer. Useful for transparent objects that should not occlude geometry behind them.

- `enabled` (bool) — `TRUE` to enable depth writing, `FALSE` to disable.

---

### `Render.SetDepthTest(enabled)`

Enables or disables depth testing. When disabled, objects are drawn in call order regardless of Z distance.

- `enabled` (bool) — `TRUE` to enable depth testing.

---

## Scissor (Clipping)

### `Render.SetScissor(x, y, width, height)`

Sets a rectangular clipping region. Only pixels inside this rectangle will be drawn.

- `x` (int) — Left edge.
- `y` (int) — Top edge.
- `width` (int) — Clip width.
- `height` (int) — Clip height.

```basic
; Only draw inside a box
Render.SetScissor(100, 100, 400, 300)
Draw.Text("This is clipped!", 110, 110, 20, 255, 255, 255, 255)
Render.ClearScissor()
```

---

### `Render.ClearScissor()`

Removes the scissor clipping region, restoring full-screen drawing.

---

## Wireframe

### `Render.SetWireframe(enabled)`

Toggles wireframe rendering mode for all subsequent 3D draw calls.

- `enabled` (bool) — `TRUE` for wireframe, `FALSE` for solid.

```basic
Render.SetWireframe(TRUE)
; All 3D meshes will render as wireframes
Draw.Cube(0, 1, 0, 2, 2, 2, 255, 255, 255, 255)
Render.SetWireframe(FALSE)
```

---

## Screenshots

### `Render.Screenshot(filePath)`

Saves the current frame as a PNG image file.

- `filePath` (string) — Output file path (e.g., `"screenshot.png"`).

```basic
IF Input.KeyPressed(KEY_F12) THEN
    Render.Screenshot("screenshot_" + STR(Time.Millisecs()) + ".png")
ENDIF
```

---

## Anti-Aliasing

### `Render.SetMSAA(samples)`

Sets the MSAA sample count for the current render context. See also `Window.SetMSAA` which must be called before window creation for full-context MSAA.

- `samples` (int) — Number of samples (2, 4, 8).

---

## Shadow Maps

### `Render.SetShadowMapSize(size)`

Sets the resolution of the shadow map texture. Larger values give sharper shadows but use more GPU memory.

- `size` (int) — Shadow map edge size in pixels (e.g., 1024, 2048, 4096).

---

## Ambient & Fog

### `Render.SetAmbient(r, g, b)`

Sets the global ambient light color for 3D PBR rendering. This is the base light that illuminates all objects regardless of light placement.

- `r` (int) — Red (0–255).
- `g` (int) — Green (0–255).
- `b` (int) — Blue (0–255).

```basic
; Warm sunset ambient
Render.SetAmbient(200, 150, 100)
```

---

### `Render.SetFog(mode, density, start, end, r, g, b)`

Configures distance fog for 3D scenes. Fog makes distant objects fade toward a color, adding depth.

- `mode` (int) — Fog mode (0 = off, 1 = linear, 2 = exponential).
- Additional parameters vary by mode.

---

## Post-Processing

### `Render.SetBloom(intensity)`

Enables or adjusts bloom post-processing (glow effect on bright areas).

- `intensity` (float) — Bloom intensity. 0 disables.

---

### `Render.SetPostProcess(enabled)`

Enables or disables the post-processing pipeline globally.

- `enabled` (bool) — `TRUE` to enable post-processing.

---

### `Render.SetMode(mode)`

Sets the render mode (forward, deferred, etc.).

- `mode` (int) — Render mode constant.

---

### `Render.SetTonemapping(mode)`

Sets the HDR tonemapping algorithm. Common modes: 0 = none, 1 = Reinhard, 2 = ACES.

- `mode` (int) — Tonemapping algorithm.

---

## Easy Mode Shortcuts

| Shortcut | Maps To |
|----------|---------|
| `SKYCOLOR(r, g, b)` | `Render.Clear(r, g, b)` |
| `AMBIENTLIGHT(r, g, b)` | `Render.SetAmbient(r, g, b)` |

---

## Full Example

A complete frame loop with 3D rendering, fog, and a screenshot hotkey.

```basic
Window.Open(1280, 720, "Render Demo")
Window.SetFPS(60)

cam = Camera.Create()
cam.pos(0, 8, 15)
cam.look(0, 0, 0)
cam.fov(60)

; Set up scene lighting
Render.SetAmbient(80, 80, 100)

WHILE NOT Window.ShouldClose()
    Render.Clear(30, 30, 50)

    ; 3D pass
    Render.Begin3D(cam)
        Draw.Grid(20, 1.0)
        Draw.Cube(0, 1, 0, 2, 2, 2, 200, 50, 50, 255)
        Draw.Sphere(4, 1.5, 0, 1.0, 50, 200, 50, 255)
    Render.End3D()

    ; 2D overlay
    Draw.Text("Render Demo", 10, 10, 24, 255, 255, 255, 255)
    Draw.Text("FPS: " + STR(Window.GetFPS()), 10, 40, 18, 100, 255, 100, 255)
    Draw.Text("F12 = Screenshot", 10, 65, 16, 180, 180, 180, 255)

    ; Screenshot hotkey
    IF Input.KeyPressed(KEY_F12) THEN
        Render.Screenshot("screenshot.png")
    ENDIF

    Render.Frame()
WEND

Camera.Free(cam)
Window.Close()
```

---

## See Also

- [WINDOW](WINDOW.md) — Window creation and management
- [CAMERA](CAMERA.md) — 3D camera setup for `Render.Begin3D`
- [DRAW](DRAW.md) — 2D and 3D drawing commands
- [POST](POST.md) — Post-processing effects (bloom, vignette, etc.)
- [EFFECT](EFFECT.md) — Visual effects (SSAO, motion blur, DOF)
