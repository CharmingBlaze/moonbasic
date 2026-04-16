# Render Commands

Frame lifecycle and render-state helpers.

**Page shape:** [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern** for lifecycle + example below).

## Core Workflow

1. **Clear**: Call **`RENDER.CLEAR`** at the start of each frame.
2. **Draw**: Issue drawing commands (**`DRAW.*`**, **`DRAW3D.*`**, entity draws, …).
3. **Present**: Call **`RENDER.FRAME`** to display the result.

---

## Full Example

```basic
WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(10, 20, 30)
    DRAW.RECTANGLE(100, 100, 50, 50, 255, 255, 255, 255)
    RENDER.FRAME()
WEND
```

---

## Frame Lifecycle

### `RENDER.CLEAR(r, g, b [, a])`

Clears the color and depth buffers. RGBA components may be **0.0–1.0** or **0–255** (normalized internally).

### `RENDER.FRAME()`

Ends the current frame, flushes any pending draw commands, and presents the result to the screen (swap buffers). Call exactly once at the end of your main loop.

---

## Screen Metrics & FPS

### `RENDER.WIDTH()` / `RENDER.HEIGHT()`

Returns physical framebuffer dimensions in pixels (float64).

### `RENDER.DRAWFPS(x, y)`

Draws the current frame rate overlay at the specified pixel coordinates.

---

## Lighting & Shadow Mapping

### `RENDER.SETAMBIENT(r, g, b [, a])`

**3D PBR** hemispheric ambient tint (per-channel multiplier on albedo). Components may be **0.0–1.0** or **0–255**. With four arguments, **`a`** scales all three RGB channels together for global ambient strength.

### `RENDER.SETSHADOWMAPSIZE(size)`

Sets the resolution of the depth texture used for shadows. Typical values: 512, 1024, 2048, 4096.

---

## Blending & Depth State

### `RENDER.SETBLEND(mode)`

Sets the Raylib blend mode (integer constant). Use globals **`BLEND_ALPHA`**, **`BLEND_ADDITIVE`**, etc.

### `RENDER.SETDEPTHWRITE(on)`

Toggles whether drawing writes to the depth buffer.

### `RENDER.SETDEPTHTEST(on)`

Toggles whether drawing performs depth testing.

---

## Capture

### `RENDER.SCREENSHOT(path)`

Writes the current framebuffer to a PNG file at the specified path.

