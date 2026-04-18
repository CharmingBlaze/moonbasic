# Effect Commands

Toggle and tune post-process rendering effects: bloom, SSAO, SSR, motion blur, depth-of-field, tonemapping, FXAA, vignette, grain, sharpen, and chromatic aberration.

Requires **fullruntime** (CGO + Raylib). See also [POST.md](POST.md) for shader-based post-process pipeline.

## Core Workflow

1. Call any `EFFECT.*` command to enable/disable effects globally.
2. Optional second/third arguments fine-tune intensity or thresholds.
3. Effects apply to the next `RENDER.FRAME`.

---

## Ambient Occlusion

### `EFFECT.SSAO(enabled [, radius [, bias]])` 

Enables/disables Screen Space Ambient Occlusion. `radius` controls the sample hemisphere size (default `0.5`). `bias` reduces self-occlusion artifacts (default `0.025`).

---

### `EFFECT.SSR(enabled [, maxDist [, resolution]])` 

Enables/disables Screen Space Reflections. `maxDist` limits ray travel distance. `resolution` scales the reflection buffer (0–1).

---

## Motion & Depth

### `EFFECT.MOTIONBLUR(enabled [, strength])` 

Enables/disables per-object motion blur. `strength` scales the blur amount (default `1.0`).

---

### `EFFECT.DEPTHOFFIELD(enabled [, focusDist [, focusRange]])` 

Enables/disables depth of field (bokeh blur). `focusDist` is the in-focus world distance. `focusRange` is the falloff distance.

---

## Bloom

### `EFFECT.BLOOM(enabled [, threshold [, intensity]])` 

Enables/disables bloom (light bleed). `threshold` is the luminance cutoff. `intensity` scales the bloom contribution.

---

## Tonemapping

### `EFFECT.TONEMAPPING(mode)` 

Sets the tonemapping operator. `mode` is a string: `"none"`, `"reinhard"`, `"aces"`, `"filmic"`, `"uncharted2"`.

---

## Sharpening

### `EFFECT.SHARPEN(enabled [, strength])` 

Enables/disables image sharpening. `strength` controls the sharpen amount.

---

## Film Effects

### `EFFECT.GRAIN(enabled [, intensity])` 

Adds film grain noise. `intensity` 0–1.

---

### `EFFECT.VIGNETTE(enabled [, strength])` 

Darkens screen edges. `strength` 0–1.

---

### `EFFECT.CHROMATICABERRATION(enabled [, offset])` 

Splits RGB channels slightly. `offset` controls the pixel separation.

---

## Anti-aliasing

### `EFFECT.FXAA(enabled)` 

Enables/disables FXAA anti-aliasing.

---

## Full Example

Cinematic look with bloom, vignette, and ACES tonemapping.

```basic
WINDOW.OPEN(1280, 720, "Effects Demo")
WINDOW.SETFPS(60)

EFFECT.BLOOM(TRUE, 0.8, 1.5)
EFFECT.TONEMAPPING("aces")
EFFECT.VIGNETTE(TRUE, 0.4)
EFFECT.FXAA(TRUE)
EFFECT.SSAO(TRUE, 0.4, 0.02)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 5, -10)
CAMERA.SETTARGET(cam, 0, 0, 0)

sun = LIGHT.CREATE("directional")
LIGHT.SETDIR(sun, -0.5, -1, -0.3)
LIGHT.SETCOLOR(sun, 255, 240, 200, 255)

cube = ENTITY.CREATECUBE(2.0)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    ENTITY.TURN(cube, 0, 45 * dt, 0)
    ENTITY.UPDATE(dt)

    RENDER.CLEAR(5, 5, 10)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
    RENDER.END3D()
    RENDER.FRAME()
WEND

EFFECT.BLOOM(FALSE)
EFFECT.SSAO(FALSE)
ENTITY.FREE(cube)
LIGHT.FREE(sun)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [POST.md](POST.md) — shader-based post-process pipeline
- [RENDER.md](RENDER.md) — `RENDER.SETAMBIENT`, shadow map
- [LIGHT.md](LIGHT.md) — lighting for PBR + bloom interaction
- [SHADER.md](SHADER.md) — custom shaders
