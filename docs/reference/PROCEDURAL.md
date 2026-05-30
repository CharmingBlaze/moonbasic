# Procedural Generation Commands

Index page for procedural noise, random number generation, and easing functions.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

- **Noise:** `n = NOISE.CREATE()` → configure type/frequency → `NOISE.SAMPLE2D(n, x, z)` → `NOISE.FREE(n)`.
- **Seeded RNG:** `RANDOMIZE(seed)` for global, or `r = RAND.CREATE(seed)` for independent streams.
- **Easing:** Wrap any 0–1 progress `t` with `SMOOTHSTEP(0, 1, t)` or `EASEIN(t, power)`.

---

## Where to find commands

Procedural generation is covered across several focused reference pages:

- **Noise** (`NOISE.*`, `PERLIN`, `FBMNOISE`, …) — see [NOISE.md](NOISE.md)
- **Random** (`RAND.*`, `RND`, `RNDRANGE`, `RANDOMIZE`, …) — see [MATH.md](MATH.md)
- **Easing** (`EASEIN`, `EASEOUT`, `EASELERP`, …) — see [EASING.md](EASING.md)
- **Biomes** (`BIOME.*`) — see [BIOME.md](BIOME.md)
- **Scatter** (`SCATTER.*`) — see [SCATTER_PROP_SPAWNER.md](SCATTER_PROP_SPAWNER.md)

For deterministic runs, seed the RNG via `RANDOMIZE` or `RAND.SEED` so procedural results are reproducible.

---

## Full Example

Terrain height map using Perlin noise with a seeded RNG for prop scatter.

```basic
WINDOW.OPEN(960, 540, "Procedural Demo")
WINDOW.SETFPS(60)

RANDOMIZE(42)   ; reproducible run

n = NOISE.CREATE()
NOISE.SETTYPE(n, NOISE_PERLIN)
NOISE.SETFREQUENCY(n, 0.06)
NOISE.SETOCTAVES(n, 4)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 20, -5)
CAMERA.SETTARGET(cam, 0, 0, 10)

RENDER.CLEAR(80, 120, 160)
RENDER.BEGIN3D(cam)
    FOR z = -16 TO 16
        FOR x = -16 TO 16
            h = NOISE.SAMPLE2D(n, x, z) * 4.0
            DRAW.CUBE(x, h * 0.5, z, 1.0, h + 0.01, 1.0, 0, 60, 140, 60, 255)
        NEXT x
    NEXT z
    DRAW.GRID(32, 1.0)
RENDER.END3D()
RENDER.FRAME()

NOISE.FREE(n)
WINDOW.CLOSE()
```

---

## See also

- [NOISE.md](NOISE.md) — `NOISE.*` procedural noise
- [RAND.md](RAND.md) — seeded RNG handles
- [MATH.md](MATH.md) — `RND`, `RNDRANGE`, `RANDOMIZE`
- [EASING.md](EASING.md) — easing functions
