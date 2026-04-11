# Render Commands

Frame lifecycle and render-state helpers.

## Core Workflow

1. **Clear**: Call `Render.Clear()` at the start of each frame.
2. **Draw**: Issue drawing commands (`Draw.*`, `Draw3D.*`, etc.).
3. **Present**: Call `Render.Frame()` to display the result.

```basic
WHILE NOT Window.ShouldClose()
    Render.Clear(10, 20, 30)
    
    ; Drawing goes here
    Draw.Rectangle(100, 100, 50, 50, 255, 255, 255, 255)
    
    Render.Frame()
WEND
```

---

## Frame Lifecycle

### `Render.Clear(r, g, b [, a])`

Clears the color and depth buffers. RGBA components may be **0.0–1.0** or **0–255** (normalized internally).

### `Render.Frame()`

Ends the current frame, flushes any pending draw commands, and presents the result to the screen (swap buffers). Call exactly once at the end of your main loop.

---

## Screen Metrics & FPS

### `Render.Width()` / `Render.Height()`

Returns physical framebuffer dimensions in pixels (float64).

### `Render.DrawFPS(x, y)`

Draws the current frame rate overlay at the specified pixel coordinates.

---

## Lighting & Shadow Mapping

### `Render.SetAmbient(r, g, b [, a])`

**3D PBR** hemispheric ambient tint (per-channel multiplier on albedo). Components may be **0.0–1.0** or **0–255**. With four arguments, **`a`** scales all three RGB channels together for global ambient strength.

### `Render.SetShadowMapSize(size)`

Sets the resolution of the depth texture used for shadows. Typical values: 512, 1024, 2048, 4096.

---

## Blending & Depth State

### `Render.SetBlend(mode)`

Sets the Raylib blend mode (integer constant). Use globals **`BLEND_ALPHA`**, **`BLEND_ADDITIVE`**, etc.

### `Render.SetDepthWrite(on)`

Toggles whether drawing writes to the depth buffer.

### `Render.SetDepthTest(on)`

Toggles whether drawing performs depth testing.

---

## Capture

### `Render.Screenshot(path)`

Writes the current framebuffer to a PNG file at the specified path.

