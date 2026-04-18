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
Clears the color and depth buffers.

- **Arguments**:
    - `r, g, b`: (Float/Integer) Color components (0-255).
    - `a`: (Float, Optional) Alpha component.
- **Returns**: (None)

---

### `RENDER.FRAME()`
Ends the frame and presents the result to the window.

---

### `RENDER.WIDTH()` / `HEIGHT()`
Returns the physical framebuffer dimensions in pixels.

- **Returns**: (Float)

---

### `RENDER.SETAMBIENT(r, g, b [, a])`
Sets the 3D PBR hemispheric ambient tint.

- **Returns**: (None)

---

### `RENDER.SETSHADOWMAPSIZE(size)`
Sets the shadow map resolution in pixels.

- **Returns**: (None)

---

### `RENDER.SETBLEND(mode)`
Sets the active blend mode.

- **Arguments**:
    - `mode`: (Integer) Blend mode constant (e.g., `BLEND_ALPHA`).

---

### `RENDER.SCREENSHOT(path)`
Saves the current framebuffer to a PNG file.

