# Procedural Generation Commands

Index page for procedural noise, random number generation, and easing functions.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Where to find commands

Procedural generation is covered across several focused reference pages:

- **Noise** (`NOISE.*`, `PERLIN`, `FBMNOISE`, …) — see [NOISE.md](NOISE.md)
- **Random** (`RAND.*`, `RND`, `RNDRANGE`, `RANDOMIZE`, …) — see [MATH.md](MATH.md)
- **Easing** (`EASEIN`, `EASEOUT`, `EASELERP`, …) — see [EASING.md](EASING.md)
- **Biomes** (`BIOME.*`) — see [BIOME.md](BIOME.md)
- **Scatter** (`SCATTER.*`) — see [SCATTER_PROP_SPAWNER.md](SCATTER_PROP_SPAWNER.md)

For deterministic runs, seed the RNG via `RANDOMIZE` or `RAND.SEED` so procedural results are reproducible.
