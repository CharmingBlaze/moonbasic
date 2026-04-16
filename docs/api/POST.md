# Post-Processing & Effect Commands

Commands for adding screen-space visual effects to the final rendered image. Post-processing effects run after all 3D and 2D rendering and modify the entire framebuffer. Effects include bloom, SSAO, motion blur, depth of field, tonemapping, vignette, chromatic aberration, and more.

## Core Concepts

- **Post-processing pipeline** — A chain of full-screen shader passes applied to the framebuffer after all scene rendering.
- **POST.*** — Commands for managing the post-processing chain (add, remove, configure shaders).
- **EFFECT.*** — High-level shortcut commands for common effects (one call = fully configured).
- Effects are applied in the order they are added.

---

## Post-Processing Pipeline

### `Post.Add(shaderPath)`

Adds a custom post-processing shader to the pipeline.

- `shaderPath` (string) — Path to a GLSL fragment shader.

**Returns:** `handle`

```basic
vignette = Post.Add("shaders/vignette.glsl")
```

---

### `Post.AddShader(shaderPath)`

Alias for `Post.Add`. Adds a shader to the post-processing chain.

---

### `Post.Remove(handle)`

Removes a post-processing shader from the pipeline.

---

### `Post.SetParam(handle, paramName, value)`

Sets a uniform parameter on a post-processing shader.

- `handle` (handle) — Post shader handle.
- `paramName` (string) — GLSL uniform name.
- `value` (float/int) — Parameter value.

```basic
Post.SetParam(vignette, "strength", 0.4)
Post.SetParam(vignette, "softness", 0.8)
```

---

### `Post.SetTonemap(mode)` / `Render.SetTonemapping(mode)`

Sets the HDR tonemapping algorithm.

- `mode` (int) — 0 = None, 1 = Reinhard, 2 = ACES Filmic.

```basic
Post.SetTonemap(2)   ; ACES Filmic
```

---

### `Render.SetPostProcess(enabled)`

Enables or disables the entire post-processing pipeline.

- `enabled` (bool) — `TRUE` to enable.

---

## Quick Effect Shorthands

These add pre-configured effects with a single call.

### `Post.Bloom(intensity)`

Adds a bloom (glow) effect.

- `intensity` (float) — Bloom strength.

```basic
Post.Bloom(0.5)
```

---

### `Post.Vignette(strength)`

Adds a vignette (darkened edges) effect.

- `strength` (float) — Vignette strength.

```basic
Post.Vignette(0.3)
```

---

### `Post.Chromatic(offset)`

Adds chromatic aberration (color fringing at edges).

- `offset` (float) — Aberration offset.

```basic
Post.Chromatic(0.005)
```

---

## High-Level Effects

### `Effect.SSAO(radius, intensity)`

Enables Screen Space Ambient Occlusion — darkens creases and corners for added depth.

- `radius` (float) — Sample radius.
- `intensity` (float) — Effect strength.

```basic
Effect.SSAO(0.5, 1.0)
```

---

### `Effect.SSR(maxSteps, stepSize)`

Enables Screen Space Reflections.

- `maxSteps` (int) — Maximum ray march steps.
- `stepSize` (float) — Step size per march.

---

### `Effect.MotionBlur(intensity)`

Enables motion blur based on camera movement.

- `intensity` (float) — Blur strength.

```basic
Effect.MotionBlur(0.5)
```

---

### `Effect.DepthOfField(focusDistance, aperture)`

Enables depth of field (bokeh blur on out-of-focus areas).

- `focusDistance` (float) — Distance to the focus plane.
- `aperture` (float) — Aperture size (larger = more blur).

```basic
Effect.DepthOfField(10.0, 0.05)
```

---

### `Effect.Bloom(threshold, intensity)`

Enables bloom with configurable threshold and intensity.

- `threshold` (float) — Brightness threshold for bloom.
- `intensity` (float) — Bloom strength.

---

### `Effect.Tonemapping(mode)`

Sets tonemapping. Same as `Post.SetTonemap`.

---

### `Effect.Sharpen(intensity)`

Adds a sharpening post-process.

- `intensity` (float) — Sharpen strength.

---

### `Effect.Grain(intensity)`

Adds film grain noise.

- `intensity` (float) — Grain strength.

```basic
Effect.Grain(0.05)   ; Subtle film grain
```

---

### `Effect.Vignette(strength)`

Adds vignette. Same as `Post.Vignette`.

---

### `Effect.ChromaticAberration(offset)`

Adds chromatic aberration. Same as `Post.Chromatic`.

---

### `Effect.FXAA(enabled)`

Enables Fast Approximate Anti-Aliasing as a post-process.

- `enabled` (bool) — `TRUE` to enable.

```basic
Effect.FXAA(TRUE)
```

---

## Full Example

A 3D scene with multiple post-processing effects.

```basic
Window.Open(1280, 720, "Post-Processing Demo")
Window.SetFPS(60)

cam = Camera.Create()
cam.pos(0, 8, 15)
cam.look(0, 2, 0)
cam.fov(60)

; Enable post-processing effects
Render.SetAmbient(60, 60, 80)
Effect.Bloom(0.8, 0.5)
Effect.SSAO(0.5, 1.0)
Effect.FXAA(TRUE)
Effect.Vignette(0.25)
Post.SetTonemap(2)   ; ACES Filmic

; Create some bright objects for bloom
WHILE NOT Window.ShouldClose()
    dt = Time.Delta()

    Render.Clear(10, 10, 20)
    Camera.Begin(cam)
        Draw.Grid(20, 1.0)

        ; Bright cube (will glow with bloom)
        Draw.Cube(0, 2, 0, 2, 2, 2, 255, 200, 100, 255)

        ; Emissive spheres
        Draw.Sphere(-4, 1, 0, 1.0, 100, 255, 200, 255)
        Draw.Sphere(4, 1, 0, 1.0, 255, 100, 200, 255)
    Camera.End(cam)

    Draw.Text("SSAO + Bloom + FXAA + Vignette + ACES", 10, 10, 18, 255, 255, 255, 255)
    Render.Frame()
WEND

Camera.Free(cam)
Window.Close()
```

---

## See Also

- [RENDER](RENDER.md) — Render pipeline and blend modes
- [CAMERA](CAMERA.md) — Camera affects DOF focus and motion blur
- [WINDOW](WINDOW.md) — MSAA (hardware AA) alternative to FXAA
