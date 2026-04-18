# Light Commands

3D light handles (directional, point, spot) with shadow casting and ambient control.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Create a light with `LIGHT.CREATE` (type: `"directional"`, `"point"`, or `"spot"`).
2. Position with `LIGHT.SETPOS` / `LIGHT.SETDIR`, color with `LIGHT.SETCOLOR`.
3. Enable shadows with `LIGHT.SETSHADOW` and set map size with `RENDER.SETSHADOWMAPSIZE`.
4. Set ambient with `RENDER.SETAMBIENT`.
5. Free with `LIGHT.FREE`.

For 2D lighting see [SPRITE.md](SPRITE.md) (`LIGHT2D.*`).

---

### `LIGHT.CREATE(type)`
Creates a new light source of the specified type.

- **Arguments**:
    - `type`: (String) `"directional"`, `"point"`, or `"spot"`.
- **Returns**: (Handle) The new light handle.
- **Example**:
    ```basic
    sun = LIGHT.CREATE("directional")
    ```

---

### `LIGHT.FREE(handle)`
Unloads the light and frees its resources.

---

### `LIGHT.SETPOS(handle, x, y, z)` / `SETDIR`
Sets the position or direction of the light.

- **Arguments**:
    - `handle`: (Handle) The light to modify.
    - `x, y, z`: (Float) World coordinates or direction vector.
- **Returns**: (Handle) The light handle (for chaining).

---

### `LIGHT.SETCOLOR(handle, r, g, b [, a])`
Sets the color and intensity of the light.

- **Arguments**:
    - `r, g, b`: (Float/Integer) Color components (0-255).
    - `a`: (Float, Optional) Multiplier for light strength.
- **Returns**: (Handle) The light handle (for chaining).

---

### `LIGHT.SETRANGE(handle, range)`
Sets the maximum distance for point/spot lights.

- **Returns**: (Handle) The light handle (for chaining).

---

### `LIGHT.SETSHADOW(handle, toggle)`
Enables or disables shadow casting.

- **Arguments**:
    - `toggle`: (Boolean) `TRUE` to enable shadows.
- **Returns**: (Handle) The light handle (for chaining).

---

### `LIGHT.SETTARGET(handle, x, y, z)`
Sets the world point for shadow camera focus.

- **Returns**: (Handle) The light handle (for chaining).

---

### `RENDER.SETAMBIENT` 

```basic
RENDER.SETAMBIENT(r, g, b)
RENDER.SETAMBIENT(r, g, b, a)
```

**3D PBR** hemispheric ambient tint (per-channel multiplier on albedo). Components may be **0.0–1.0** or **0–255** (values > 1 are normalized as 8-bit). With **four** arguments, **`a`** scales **all three** RGB channels together (useful for global ambient strength).

**Parameters**

| Name | Type | Description |
|---|---|---|
| r, g, b | float | Ambient tint per channel. |
| a | float | Optional. Scales r, g, b together (default `1` when using the 3-argument form). |

**Example**

```basic
RENDER.SETAMBIENT(0.05, 0.06, 0.08)
RENDER.SETAMBIENT(13, 15, 20, 255)
```

**See also:** **`RENDER.SETSHADOWMAPSIZE`**, [MODEL.md](MODEL.md)

---

### `RENDER.SETSHADOWMAPSIZE` 

```basic
RENDER.SETSHADOWMAPSIZE(size)
```

Sets shadow map resolution in pixels. Larger = sharper shadows, more VRAM for the depth target. Clamped by the engine.

**Parameters**

| Name | Type | Description |
|---|---|---|
| size | int | Edge length in pixels (prefer powers of two, e.g. 512–4096). |

**Example**

```basic
RENDER.SETSHADOWMAPSIZE(2048)
```

---

## PBR and shadows

- Materials use the engine PBR path so the fragment shader receives `lightDir`, `lightColor`, `ambientColor`, `lightVP`, `shadowBiasK`, and the shadow depth map.
- If no shadow caster is registered, shadows are skipped for that frame.

---

## Full Example

A scene with a directional shadow-casting light and a point fill light.

```basic
WINDOW.OPEN(960, 540, "Lighting Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 8, -12)
CAMERA.SETTARGET(cam, 0, 0, 0)

; Shadow-casting sun
sun = LIGHT.CREATE("directional")
LIGHT.SETDIR(sun, -0.5, -1.0, -0.3)
LIGHT.SETCOLOR(sun, 255, 240, 200, 255)
LIGHT.SETSHADOW(sun, TRUE)
RENDER.SETSHADOWMAPSIZE(2048)

; Warm fill point light
fill = LIGHT.CREATE("point")
LIGHT.SETPOS(fill, 3, 4, -2)
LIGHT.SETCOLOR(fill, 255, 180, 100, 180)
LIGHT.SETRANGE(fill, 12.0)

RENDER.SETAMBIENT(30, 30, 50, 255)

cube = ENTITY.CREATECUBE(2.0)

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(10, 15, 25)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW3D.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

LIGHT.FREE(sun)
LIGHT.FREE(fill)
ENTITY.FREE(cube)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [MODEL.md](MODEL.md) — models and materials
- [RENDER.md](RENDER.md) — `RENDER.SETAMBIENT`, shadow map size
- [SHADER.md](SHADER.md) — custom PBR shaders
