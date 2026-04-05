# Noise (`NOISE.*`)

moonBASIC exposes **stateful** procedural noise under the **`Noise.*`** namespace. The implementation is **pure Go** ([`runtime/procnoise`](../../runtime/procnoise)) shared with legacy globals **`PERLIN`**, **`SIMPLEX`**, **`VORONOI`**, **`FBMNOISE`** in [`runtime/mbgame`](../../runtime/mbgame) (same core at implicit seed **0**). Use **`NOISE.*`** when you need seeds, fractal settings, or fill helpers; use globals for one-off expressions.

**Memory:** `NoiseObject` holds only Go fields. **`NOISE.FREE`** is idempotent. Native **Raylib** work happens only for **`NOISE.FILLIMAGE`** (CPU image pixels). See [`MEMORY.md`](../MEMORY.md).

> **Common mistake:** moonBASIC is case-agnostic. A variable named `noise` becomes `NOISE` and **shadows** the `Noise.*` namespace — `Noise.MakeFractal` may be parsed as a **method call on your handle**. Use names like `ng`, `gen`, or `terrainNoise`.

**See also:** [`TERRAIN.FILLPERLIN`](../../compiler/builtinmanifest/commands.json) (heightfield fill), legacy **`PERLIN`** / **`FBMNOISE`** in mbgame.

---

### Noise.Make

```basic
ng = Noise.Make()
```

Creates a generator with defaults: type `"perlin"`, seed `1337`, frequency `0.01`, `3` octaves for fractal types.

**Returns** — noise handle (`int`).

**Example**

```basic
ng = Noise.Make()
Noise.SetSeed(ng, 99)
Noise.SetFrequency(ng, 0.004)
h# = Noise.Get(ng, 10, 20)
Noise.Free(ng)
```

---

### Noise.Free

```basic
Noise.Free(ng)
```

Releases generator state. Safe to call twice (second is a no-op on the handle table side after invalidation).

**Parameters**

| Name | Type | Description |
|------|------|-------------|
| ng | int | Noise handle |

---

### Noise.SetType

```basic
Noise.SetType(ng, type$)
```

Selects the algorithm **before the first sample**. `type$` examples: `"perlin"`, `"simplex"`, `"simplex_smooth"`, `"value"`, `"cellular"`, `"fractal_fbm"`, `"fractal_ridged"`, `"fractal_pingpong"`, `"domain_warp"`.

**Parameters**

| Name | Type | Description |
|------|------|-------------|
| ng | int | Noise handle |
| type$ | string | Algorithm name (case-insensitive) |

---

### Noise.SetSeed

```basic
Noise.SetSeed(ng, seed)
```

Integer seed for deterministic worlds (multiplayer / replay).

**Parameters**

| Name | Type | Description |
|------|------|-------------|
| ng | int | Noise handle |
| seed | int | Seed |

---

### Noise.SetFrequency

```basic
Noise.SetFrequency(ng, freq#)
```

Feature size: lower = smoother/larger hills; typical terrain `0.001`–`0.05`.

---

### Noise.SetOctaves / Noise.SetLacunarity / Noise.SetGain

```basic
Noise.SetOctaves(ng, count)
Noise.SetLacunarity(ng, lac#)
Noise.SetGain(ng, gain#)
```

Fractal controls (used by `fractal_*` types). Defaults: octaves `3`, lacunarity `2`, gain `0.5`.

---

### Noise.SetWeightedStrength

```basic
Noise.SetWeightedStrength(ng, strength#)
```

Emphasises higher octaves when using **`fractal_fbm`** (`0` = off, `1` = strong).

---

### Noise.SetPingPongStrength

```basic
Noise.SetPingPongStrength(ng, strength#)
```

Shapes **`fractal_pingpong`** output (default internal `2` if unset).

---

### Noise.SetCellularType / SetCellularDistance / SetCellularJitter

```basic
Noise.SetCellularType(ng, type$)
Noise.SetCellularDistance(ng, func$)
Noise.SetCellularJitter(ng, jitter#)
```

Cellular / Voronoi flavour. `type$` examples: `"distance"`, `"cell_value"`. Distance `"manhattan"` selects a different metric (approximate). Jitter is reserved for future fine-tuning.

---

### Noise.SetDomainWarpType / Noise.SetDomainWarpAmplitude

```basic
Noise.SetDomainWarpType(ng, type$)
Noise.SetDomainWarpAmplitude(ng, amp#)
```

`type$` is stored for compatibility; warp uses internal low-frequency **`Simplex2`** offsets. **`amp#`** scales warp strength (default `1`).

---

### Noise.Get

```basic
h# = Noise.Get(ng, x#, y#)
```

Samples **2D** noise ~`[-1,1]`. Locks configuration (no further **`Set*`**).

**Parameters**

| Name | Type | Description |
|------|------|-------------|
| ng | int | Handle |
| x#, y# | float | World coordinates |

---

### Noise.Get3D

```basic
h# = Noise.Get3D(ng, x#, y#, z#)
```

Cheap 3D field (blended planes / cellular blend). ~`[-1,1]`.

---

### Noise.GetDomainWarped

```basic
h# = Noise.GetDomainWarped(ng, x#, y#)
```

Applies domain warp, then evaluates the active type (turbulent coastlines, etc.).

---

### Noise.GetNorm

```basic
h# = Noise.GetNorm(ng, x#, y#)
```

Returns **`0..1`**: `(Get + 1) * 0.5` clamped.

---

### Noise.GetTileable

```basic
h# = Noise.GetTileable(ng, x#, y#, w#, h#)
```

Approximate **seamless** tiling using a torus parameterisation; `w#`, `h#` are tile size in the same units as `x`, `y`.

---

### Noise.FillArray

```basic
Noise.FillArray(ng, arr, width, height, offsetX#, offsetY#)
```

Writes **`width*height`** floats into **`arr`** (numeric array), row-major. Values ~`[-1,1]`.

> **Common mistake:** `arr` must have at least `width*height` elements.

---

### Noise.FillArrayNorm

Same as **`FillArray`**, but writes **`0..1`**.

---

### Noise.FillImage

```basic
Noise.FillImage(ng, img, offsetX#, offsetY#)
```

Fills a greyscale **`Image`** (CPU) for debugging / textures. **Requires CGO** (Raylib). On `!cgo` builds, returns an error.

---

### Noise.MakePerlin / MakeSimplex / MakeFractal / MakeCellular / MakeDomainWarp

```basic
ng = Noise.MakePerlin(seed, freq#)
ng = Noise.MakeSimplex(seed, freq#)
ng = Noise.MakeFractal(seed, freq#, octaves, type$)
ng = Noise.MakeCellular(seed, freq#, celltype$)
ng = Noise.MakeDomainWarp(seed, freq#, amp#)
```

Convenience constructors (pre-configured, no separate **`SetType`** needed).  
**`MakeFractal`**: `type$` is **`"fbm"`**, **`"ridged"`**, or **`"pingpong"`** (aliases accepted).

**Example**

```basic
ng = Noise.MakeFractal(42, 0.005, 6, "ridged")
h# = Noise.GetNorm(ng, x#, z#)
Noise.Free(ng)
```

---

## Choosing a noise type

| type$ | Typical use |
|--------|-------------|
| `simplex` | Smooth hills, general terrain |
| `fractal_fbm` | Rolling organic terrain |
| `fractal_ridged` | Ridges, cliffs |
| `cellular` + `distance` | Voronoi-style features |
| `domain_warp` | Warped domains, turbulent blending |

---

## Samples in-tree

| File | Purpose |
|------|---------|
| [`testdata/noise_test.mb`](../../testdata/noise_test.mb) | Headless checks |
| [`testdata/noise_terrain.mb`](../../testdata/noise_terrain.mb) | Windowed greyscale preview |

Run: `go run . --check testdata/noise_test.mb`
