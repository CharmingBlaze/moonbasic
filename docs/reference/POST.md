# Post Commands

Shader-based post-process pipeline: add named effects or custom shader passes, set parameters, and control tonemapping.

For toggle-based effect shortcuts see [EFFECT.md](EFFECT.md).

## Core Workflow

1. `POST.ADD(effectName)` to enable a named effect, or `POST.ADDSHADER(shaderHandle)` for custom passes.
2. `POST.SETPARAM(effectName, paramName, value)` to tune parameters.
3. Effects are applied in order each `RENDER.FRAME`.
4. `POST.REMOVE(effectName)` to disable.

---

## Adding Effects

### `POST.ADD(effectName)` 

Enables a named post-process effect. Effect names: `"bloom"`, `"vignette"`, `"chromatic"`, `"ssao"`, `"ssr"`, `"motionblur"`, `"dof"`, `"fxaa"`, `"grain"`, `"sharpen"`.

---

### `POST.ADDSHADER(shaderHandle)` 

Adds a custom GLSL shader as a full-screen post-process pass. `shaderHandle` from `SHADER.LOAD`.

---

### `POST.REMOVE(effectName)` 

Removes a previously added effect by name.

---

## Parameters

### `POST.SETPARAM(effectName, paramName, value)` 

Sets a float parameter on a named effect. Example: `POST.SETPARAM("bloom", "threshold", 0.7)`.

---

### `POST.SETTONEMAP(mode)` 

Sets the tonemapping mode as an integer constant. `0` = none, `1` = Reinhard, `2` = ACES, `3` = filmic.

---

## Shortcuts

### `POST.BLOOM(intensity)` 

Quick-set bloom intensity (enables if not already added).

---

### `POST.CHROMATIC(offset)` 

Quick-set chromatic aberration offset.

---

### `POST.VIGNETTE(strength, softness)` 

Quick-set vignette strength and softness.

---

## Full Example

Custom shader post-process pass combined with built-in bloom.

```basic
WINDOW.OPEN(960, 540, "Post Demo")
WINDOW.SETFPS(60)

POST.ADD("bloom")
POST.SETPARAM("bloom", "threshold", 0.6)
POST.SETPARAM("bloom", "intensity", 1.8)
POST.ADD("vignette")
POST.SETPARAM("vignette", "strength", 0.35)
POST.SETTONEMAP(2)   ; ACES

cam  = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 4, -8)
CAMERA.SETTARGET(cam, 0, 0, 0)
cube = ENTITY.CREATECUBE(2.0)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    ENTITY.TURN(cube, 0, 30 * dt, 0)
    ENTITY.UPDATE(dt)
    RENDER.CLEAR(5, 5, 10)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
    RENDER.END3D()
    RENDER.FRAME()
WEND

POST.REMOVE("bloom")
POST.REMOVE("vignette")
ENTITY.FREE(cube)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [EFFECT.md](EFFECT.md) — toggle-style effect shortcuts
- [SHADER.md](SHADER.md) — writing custom shader passes
- [RENDER.md](RENDER.md) — render pipeline overview
