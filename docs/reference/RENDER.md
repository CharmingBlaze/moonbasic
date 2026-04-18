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

---

## Extended Command Reference

### Frame loop aliases

| Command | Description |
|--------|-------------|
| `RENDER.BEGINFRAME()` | Alias of `RENDER.CLEAR` / frame start. |
| `RENDER.ENDFRAME()` | Alias of `RENDER.FRAME` / frame present. |
| `RENDER.SETFPS(fps)` | Set target frames per second (alias of `WINDOW.SETFPS`). |
| `RENDER.DRAWFPS(x, y)` | Draw built-in FPS counter overlay at `(x, y)`. |

### 3D mode aliases

| Command | Description |
|--------|-------------|
| `RENDER.BEGIN3D(cam)` | Begin 3D rendering with `cam` (canonical). |
| `RENDER.END3D()` | End 3D rendering (canonical). |
| `RENDER.BEGINMODE3D(cam)` | Alias of `RENDER.BEGIN3D`. |
| `RENDER.ENDMODE3D()` | Alias of `RENDER.END3D`. |

### 2D mode

| Command | Description |
|--------|-------------|
| `RENDER.BEGINMODE2D(cam2d)` | Begin 2D camera-transformed rendering. |
| `RENDER.ENDMODE2D()` | End 2D camera mode. |
| `RENDER.SET2DAMBIENT(r, g, b, a)` | Set 2D ambient light color for lit sprites. |

### Shader pass

| Command | Description |
|--------|-------------|
| `RENDER.BEGINSHADER(shader)` | Activate a shader for subsequent draws. |
| `RENDER.ENDSHADER()` | Deactivate the current shader. |

### GPU state

| Command | Description |
|--------|-------------|
| `RENDER.SETBLENDMODE(mode)` | Set blend mode (`NONE`, `ALPHA`, `ADD`, `MULTIPLY`, `SUBTRACT`). |
| `RENDER.SETDEPTHTEST(bool)` | Enable/disable depth test. |
| `RENDER.SETDEPTHWRITE(bool)` | Enable/disable depth buffer writes. |
| `RENDER.SETDEPTHMASK(bool)` | Alias of `SETDEPTHWRITE`. |
| `RENDER.SETCULLFACE(mode)` | Face culling (`NONE`, `FRONT`, `BACK`). |
| `RENDER.SETWIREFRAME(bool)` | Global wireframe toggle. |
| `RENDER.SETSCISSOR(x, y, w, h)` | Enable scissor rectangle. |
| `RENDER.CLEARSCISSOR()` | Disable scissor test. |
| `RENDER.CLEARCACHE()` | Flush internal render state cache. |
| `RENDER.SETMSAA(samples)` | Set MSAA sample count (requires restart). |
| `RENDER.SETMODE(mode)` | Switch between rendering modes (forward/deferred). |

### Post-processing & sky

| Command | Description |
|--------|-------------|
| `RENDER.SETPOSTPROCESS(shader)` | Assign a fullscreen post-process shader. |
| `RENDER.SETSKYBOX(texHandle)` | Set environment/skybox cubemap. |
| `RENDER.SETIBLINTENSITY(v)` | Set IBL (image-based lighting) intensity. |
| `RENDER.SETIBLSPLIT(v)` | Set IBL split factor. |
| `RENDER.SETTONEMAPPING(mode)` | Set tone-mapping operator (`LINEAR`, `ACES`, `FILMIC`). |
| `RENDER.SETBLOOM(threshold, intensity)` | Configure bloom post effect. |
| `RENDER.SETFOG(r, g, b, near, far)` | Set linear fog color and range. |

## See also

- [DRAW2D.md](DRAW2D.md) — 2D draw commands
- [DRAW3D.md](DRAW3D.md) — 3D draw commands
- [SHADER.md](SHADER.md) — custom shaders
- [WINDOW.md](WINDOW.md) — window and FPS settings

