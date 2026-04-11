# Light and shadow (3D)

moonBASIC exposes **3D** lights as small **CPU-side handles** (`TypeName` `Light`). The built-in **PBR + shadow** path uses a **single directional** light for diffuse shading and **one shadow-casting light** at a time. Registry keys use **dots and uppercase** (e.g. `LIGHT.MAKE`); the tables below show both **PascalCase** names (as in specs) and the **canonical keys**.

**Engine-style constructors (`LIGHT.CREATEPOINT`, `CREATEDIRECTIONAL`, `CREATESPOT`, …):** [CAMERA_LIGHT_RENDER.md](CAMERA_LIGHT_RENDER.md).

**GPU / memory**

- **`RENDER.SETSHADOWMAPSIZE`** controls one depth render target; larger sizes use more VRAM and bandwidth.
- Call **`LIGHT.FREE`** when a light is no longer needed. The underlying **`Heap.Free`** invalidates the handle after the first successful free; calling **`LIGHT.FREE`** again with the same numeric handle returns a **runtime error** (stale handle). The light object’s own teardown is safe once invoked.

For **2D** lighting, use **`LIGHT2D.*`** and **`RENDER.SET2DAMBIENT`** — see [SPRITE.md](SPRITE.md).

---

### `Light.Make(type)`
Creates a new light source of the specified type. Returns a **handle** to the light.
- `type`: The type of light: `"directional"`, `"point"`, or `"spot"`.

### `Light.Free(handle)`
Unloads the light and frees its resources.

---

### `Light.SetPos(handle, x, y, z)`
Sets the world position of a point or spot light.

### `Light.SetDir(handle, x, y, z)`
Sets the direction vector for a directional or spot light.

---

### `Light.SetColor(handle, r, g, b [, a])`
Sets the color and intensity of the light (0-255). The optional alpha component multiplies the overall light strength.

### `Light.SetRange(handle, range)`
Sets the maximum distance at which the light has an effect (for point and spot lights).

---

### `Light.SetShadow(handle, toggle)`
Enables or disables shadow casting for the light. Only **one** shadow-casting light is supported at a time.

### `Light.SetInnerCone(handle, degrees)` / `Light.SetOuterCone(handle, degrees)`
Sets the inner and outer half-cone angles for spotlights in degrees.

### `Light.SetTarget(handle, x, y, z)`
Sets the world point the **orthographic shadow camera** looks at. Correctly framing your scene in this volume is required for shadows.

---

### Render.SetAmbient

```basic
RENDER.SETAMBIENT(r, g, b)
RENDER.SETAMBIENT(r, g, b, a)
```

**3D PBR** hemispheric ambient tint (per-channel multiplier on albedo). Components may be **0.0–1.0** or **0–255** (values &gt; 1 are normalized as 8-bit). With **four** arguments, **`a`** scales **all three** RGB channels together (useful for global ambient strength).

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

**See also:** `RENDER.SETSHADOWMAPSIZE`, [MODEL.md](MODEL.md)

---

### Render.SetShadowMapSize

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

**See also:** [MODEL.md](MODEL.md), [MEMORY.md](../MEMORY.md), [ARCHITECTURE.md](../../ARCHITECTURE.md)
