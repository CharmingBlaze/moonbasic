# Noise Commands

Stateful procedural noise generators: Perlin, Simplex, cellular, fractal, and domain warp.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Create a generator with `NOISE.MAKE` or a convenience constructor (`NOISE.MAKEPERLIN`, `NOISE.MAKEFRACTAL`, etc.).
2. Configure with `NOISE.SETSEED`, `NOISE.SETFREQUENCY`, `NOISE.SETOCTAVES`, etc.
3. Sample with `NOISE.GET`, `NOISE.GET3D`, `NOISE.GETNORM`, or fill arrays/images.
4. Free with `NOISE.FREE`.

> **Naming tip:** Do not name your variable `noise` — it shadows the namespace. Use `ng`, `gen`, or `terrainNoise`.

---

### NOISE.MAKE 

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
h = Noise.Get(ng, 10, 20)
Noise.Free(ng)
```

---

### NOISE.FREE 

```basic
Noise.Free(ng)
```

Releases generator state. Safe to call twice (second is a no-op on the handle table side after invalidation).

**Parameters**

| Name | Type | Description |
|------|------|-------------|
| ng | int | Noise handle |

---

### NOISE.SETTYPE 

```basic
Noise.SetType(ng, type)
```

Selects the algorithm **before the first sample**. `type` examples: `"perlin"`, `"simplex"`, `"simplex_smooth"`, `"value"`, `"cellular"`, `"fractal_fbm"`, `"fractal_ridged"`, `"fractal_pingpong"`, `"domain_warp"`.

**Parameters**

| Name | Type | Description |
|------|------|-------------|
| ng | int | Noise handle |
| type | string | Algorithm name (case-insensitive) |

---

### NOISE.SETSEED 

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

### NOISE.SETFREQUENCY 

```basic
Noise.SetFrequency(ng, freq)
```

Feature size: lower = smoother/larger hills; typical terrain `0.001`–`0.05`.

---

### NOISE.SETOCTAVES / NOISE.SETLACUNARITY / NOISE.SETGAIN 

```basic
Noise.SetOctaves(ng, count)
Noise.SetLacunarity(ng, lac)
Noise.SetGain(ng, gain)
```

Fractal controls (used by `fractal_*` types). Defaults: octaves `3`, lacunarity `2`, gain `0.5`.

---

### NOISE.SETWEIGHTEDSTRENGTH 

```basic
Noise.SetWeightedStrength(ng, strength)
```

Emphasises higher octaves when using **`fractal_fbm`** (`0` = off, `1` = strong).

---

### NOISE.SETPINGPONGSTRENGTH 

```basic
Noise.SetPingPongStrength(ng, strength)
```

Shapes **`fractal_pingpong`** output (default internal `2` if unset).

---

### NOISE.SETCELLULARTYPE / NOISE.SETCELLULARDISTANCE / NOISE.SETCELLULARJITTER 

```basic
Noise.SetCellularType(ng, type)
Noise.SetCellularDistance(ng, func)
Noise.SetCellularJitter(ng, jitter)
```

Cellular / Voronoi flavour. `type` examples: `"distance"`, `"cell_value"`. Distance `"manhattan"` selects a different metric (approximate). Jitter is reserved for future fine-tuning.

---

### NOISE.SETDOMAINWARPTYPE / NOISE.SETDOMAINWARPAMPLITUDE 

```basic
Noise.SetDomainWarpType(ng, type)
Noise.SetDomainWarpAmplitude(ng, amp)
```

`type` is stored for compatibility; warp uses internal low-frequency **`Simplex2`** offsets. **`amp`** scales warp strength (default `1`).

---

### NOISE.GET 

```basic
h = Noise.Get(ng, x, y)
```

Samples **2D** noise ~`[-1,1]`. Locks configuration (no further **`Set*`**).

**Parameters**

| Name | Type | Description |
|------|------|-------------|
| ng | int | Handle |
| x, y | float | World coordinates |

---

### NOISE.GET3D 

```basic
h = Noise.Get3D(ng, x, y, z)
```

Cheap 3D field (blended planes / cellular blend). ~`[-1,1]`.

---

### NOISE.GETDOMAINWARPED 

```basic
h = Noise.GetDomainWarped(ng, x, y)
```

Applies domain warp, then evaluates the active type (turbulent coastlines, etc.).

---

### NOISE.GETNORM 

```basic
h = Noise.GetNorm(ng, x, y)
```

Returns **`0..1`**: `(Get + 1) * 0.5` clamped.

---

### NOISE.GETTILEABLE 

```basic
h = Noise.GetTileable(ng, x, y, w, h)
```

Approximate **seamless** tiling using a torus parameterisation; `w`, `h` are tile size in the same units as `x`, `y`.

---

### NOISE.FILLARRAY 

```basic
Noise.FillArray(ng, arr, width, height, offsetX, offsetY)
```

Writes **`width*height`** floats into **`arr`** (numeric array), row-major. Values ~`[-1,1]`.

> **Common mistake:** `arr` must have at least `width*height` elements.

---

### NOISE.FILLARRAYNORM 

Same as **`FillArray`**, but writes **`0..1`**.

---

### NOISE.FILLIMAGE 

```basic
Noise.FillImage(ng, img, offsetX, offsetY)
```

Fills a greyscale **`Image`** (CPU) for debugging / textures. **Requires CGO** (Raylib). On `!cgo` builds, returns an error.

---

### NOISE.MAKEPERLIN / NOISE.MAKESIMPLEX / NOISE.MAKEFRACTAL / NOISE.MAKECELLULAR / NOISE.MAKEDOMAINWARP 

```basic
ng = Noise.MakePerlin(seed, freq)
ng = Noise.MakeSimplex(seed, freq)
ng = Noise.MakeFractal(seed, freq, octaves, type)
ng = Noise.MakeCellular(seed, freq, celltype)
ng = Noise.MakeDomainWarp(seed, freq, amp)
```

Convenience constructors (pre-configured, no separate **`SetType`** needed).  
**`MakeFractal`**: `type` is **`"fbm"`**, **`"ridged"`**, or **`"pingpong"`** (aliases accepted).

**Example**

```basic
ng = Noise.MakeFractal(42, 0.005, 6, "ridged")
h = Noise.GetNorm(ng, x, z)
Noise.Free(ng)
```

---

## Choosing a noise type

| type | Typical use |
|--------|-------------|
| `simplex` | Smooth hills, general terrain |
| `fractal_fbm` | Rolling organic terrain |
| `fractal_ridged` | Ridges, cliffs |
| `cellular` + `distance` | Voronoi-style features |
| `domain_warp` | Warped domains, turbulent blending |

---

## Full Example

Fractal terrain heightmap sampled into a 2D display.

```basic
WINDOW.OPEN(512, 512, "Noise Demo")
WINDOW.SETFPS(60)

ng = NOISE.MAKE()
NOISE.SETTYPE(ng, "fractal_fbm")
NOISE.SETSEED(ng, 42)
NOISE.SETFREQUENCY(ng, 0.005)
NOISE.SETOCTAVES(ng, 5)

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(0, 0, 0)
    FOR y = 0 TO 511
        FOR x = 0 TO 511
            v = NOISE.GETNORM(ng, x, y)   ; 0.0..1.0
            c = INT(v * 255)
            DRAW.PIXEL(x, y, c, c, c, 255)
        NEXT x
    NEXT y
    RENDER.FRAME()
WEND

NOISE.FREE(ng)
WINDOW.CLOSE()
```

---

## Samples in-tree

| File | Purpose |
|------|---------|
| [`testdata/noise_test.mb`](../../testdata/noise_test.mb) | Headless checks |
| [`testdata/noise_terrain.mb`](../../testdata/noise_terrain.mb) | Windowed greyscale preview |

Run: `go run . --check testdata/noise_test.mb`
