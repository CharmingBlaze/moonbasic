# Light and shadow (3D)

moonBASIC exposes **3D** lights as small **CPU-side handles** (`TypeName` `Light`). The built-in **PBR + shadow** path uses a **single directional** light for diffuse shading and **one shadow-casting light** at a time. Registry keys use **dots and uppercase** (e.g. `LIGHT.MAKE`); the tables below show both **PascalCase** names (as in specs) and the **canonical keys**.

**Engine-style constructors (`LIGHT.CREATEPOINT`, `CREATEDIRECTIONAL`, `CREATESPOT`, …):** [CAMERA_LIGHT_RENDER.md](CAMERA_LIGHT_RENDER.md).

**GPU / memory**

- **`RENDER.SETSHADOWMAPSIZE`** controls one depth render target; larger sizes use more VRAM and bandwidth.
- Call **`LIGHT.FREE`** when a light is no longer needed. The underlying **`Heap.Free`** invalidates the handle after the first successful free; calling **`LIGHT.FREE`** again with the same numeric handle returns a **runtime error** (stale handle). The light object’s own teardown is safe once invoked.

For **2D** lighting, use **`LIGHT2D.*`** and **`RENDER.SET2DAMBIENT`** — see [SPRITE.md](SPRITE.md).

---

### Light.Make

```basic
sun = Light.Make()
sun = Light.Make("directional")
```

Creates a light handle. Optional `kind`: `"directional"`, `"point"`, or `"spot"` (stored for your logic; the stock PBR path is built around a **directional** sun).

**Parameters**

| Name | Type | Description |
|---|---|---|
| kind | string | Optional. Light kind label. |

**Returns** — handle (integer).

**Notes** — Default colour is white, intensity `1`, direction roughly toward the origin, shadow target `(0, 2, 0)`.

> **Common mistake:** Expecting multiple shadow casters. Only **one** light should have **`LIGHT.SETSHADOW`** enabled at a time; the last one wins.

**Example**

```basic
sun = Light.Make("directional")
LIGHT.SETDIR(sun, -0.4, -0.8, -0.4)
```

**See also:** `LIGHT.SETSHADOW`, `RENDER.SETSHADOWMAPSIZE`

---

### Light.Free

```basic
Light.Free(sun)
```

Releases the light handle via the heap. Clears **shadow caster** registration if this light was casting shadows.

**Parameters**

| Name | Type | Description |
|---|---|---|
| light | handle | Light returned from `Light.Make`. |

> **Common mistake:** Calling **`LIGHT.FREE`** twice with the same handle value after the first successful free — the heap reports a stale handle error.

**See also:** `LIGHT.SETSHADOW`

---

### Light.SetDir

```basic
LIGHT.SETDIR(light, x#, y#, z#)
```

Sets the **light travel direction** (into the scene), normalized internally. Used for the directional sun and shadow frustum alignment.

**Parameters**

| Name | Type | Description |
|---|---|---|
| light | handle | Light handle. |
| x, y, z | float | Direction components (need not be pre-normalized). |

---

### Light.SetShadow

```basic
LIGHT.SETSHADOW(light, enabled?)
```

Enables shadow mapping for this light. **Only one** light may cast shadows; the most recently enabled wins.

**Parameters**

| Name | Type | Description |
|---|---|---|
| light | handle | Light handle. |
| enabled? | bool or int | Non-zero enables. |

---

### Light.SetColor

```basic
LIGHT.SETCOLOR(light, r, g, b)
LIGHT.SETCOLOR(light, r, g, b, a)
```

Sets RGB **and optional alpha**. RGB may be **0.0–1.0** or **0–255** (any channel &gt; 1 scales all three channels as 8-bit). The **fourth** argument `a` multiplies the diffuse contribution (after intensity); use **0.0–1.0** or **0–255**.

**Parameters**

| Name | Type | Description |
|---|---|---|
| light | handle | Light handle. |
| r, g, b | float / int | Red, green, blue. |
| a | float / int | Optional. Overall colour scale (default `1`). |

> **Common mistake:** Forgetting that **`a`** applies on top of **`LIGHT.SETINTENSITY`** — both scale the final diffuse RGB.

**Example**

```basic
LIGHT.SETCOLOR(sun, 255, 240, 220, 255)
LIGHT.SETINTENSITY(sun, 1.2)
```

**See also:** `LIGHT.SETINTENSITY`

---

### Light.SetIntensity

```basic
LIGHT.SETINTENSITY(light, amount#)
```

Non-negative scale for diffuse RGB (applied together with **`LIGHT.SETCOLOR`**).

**Parameters**

| Name | Type | Description |
|---|---|---|
| light | handle | Light handle. |
| amount | float | Intensity; negative values are clamped to `0`. |

---

### Light.SetPosition / Light.SetPos

```basic
LIGHT.SETPOSITION(light, x#, y#, z#)
LIGHT.SETPOS(light, x#, y#, z#)
```

World position for **point** / **spot** workflows (stored for API completeness; extend custom shaders if you need full point-light shading in the default PBR path).

**Parameters**

| Name | Type | Description |
|---|---|---|
| light | handle | Light handle. |
| x, y, z | float | World position. |

---

### Light.SetTarget

```basic
LIGHT.SETTARGET(light, x#, y#, z#)
```

World point the **orthographic shadow camera** looks at (default `0, 2, 0`). Adjust so your scene sits in the shadow frustum.

**Parameters**

| Name | Type | Description |
|---|---|---|
| light | handle | Light handle. |
| x, y, z | float | Look-at point in world space. |

> **Common mistake:** Moving only **`LIGHT.SETDIR`** and wondering why shadows slide — pair direction with a sensible **`LIGHT.SETTARGET`** for the shadow volume.

---

### Light.SetShadowBias

```basic
LIGHT.SETSHADOWBIAS(light, bias#)
```

Multiplier for **depth bias** in shadow sampling (typical **0.5–2.0**, clamped internally). Higher reduces **acne**; lower reduces **peter-panning**.

**Parameters**

| Name | Type | Description |
|---|---|---|
| light | handle | Light handle. |
| bias | float | Bias multiplier (default `1`). |

---

### Light.SetInnerCone / Light.SetOuterCone

```basic
LIGHT.SETINNERCONE(light, angle#)
LIGHT.SETOUTERCONE(light, angle#)
```

Spotlight cone angles in **degrees** (stored for API completeness).

**Parameters**

| Name | Type | Description |
|---|---|---|
| light | handle | Light handle. |
| angle | float | Half-cone style angle in degrees (engine defaults: inner 25°, outer 35° at creation). |

---

### Light.SetRange

```basic
LIGHT.SETRANGE(light, range#)
```

Attenuation range for point/spot lights (stored).

**Parameters**

| Name | Type | Description |
|---|---|---|
| light | handle | Light handle. |
| range | float | Range distance. |

---

### Light.Enable

```basic
LIGHT.ENABLE(light, enabled?)
```

Master on/off. When disabled, diffuse contribution is zero and shadow caster registration is cleared if this light was the caster.

**Parameters**

| Name | Type | Description |
|---|---|---|
| light | handle | Light handle. |
| enabled? | bool or int | Non-zero enables. |

---

### Light.IsEnabled

```basic
ok = LIGHT.ISENABLED(light)
```

Returns **1** if the light exists and is enabled, **0** otherwise.

**Returns** — integer (`0` or `1`).

---

### Render.SetAmbient

```basic
RENDER.SETAMBIENT(r#, g#, b#)
RENDER.SETAMBIENT(r#, g#, b#, a#)
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
