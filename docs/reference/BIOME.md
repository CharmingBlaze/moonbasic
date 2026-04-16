# Biome Commands

Lightweight **temperature / humidity** state (`TagBiome`) for driving ambience or future terrain/weather blending. **CGO** required.

Biomes **do not** replace **`SCATTER.APPLY`**, **`TERRAIN.FILLPERLIN`**, or **`WEATHER.*`** — combine them in your loop or data layer.

Page shape: [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

Create a biome with **`BIOME.CREATE(name)`**, tune **`BIOME.SETTEMP`** and **`BIOME.SETHUMIDITY`**, and read the values where your gameplay or ambience system needs them.

---

### `BIOME.CREATE(name)`

Creates a biome descriptor. **`name`** is a string label. **`BIOME.MAKE`** is a deprecated alias. Returns a **handle**.

---

### `BIOME.FREE(biome)`

Frees the biome handle.

---

### `BIOME.SETTEMP(biome, t)`

Sets temperature (see runtime: manifest describes **celsius** scale).

---

### `BIOME.SETHUMIDITY(biome, h)`

Sets humidity **0–1** (normalized).

---

## Full Example

```basic
biome = BIOME.CREATE("forest")
BIOME.SETTEMP(biome, 18.0)
BIOME.SETHUMIDITY(biome, 0.65)
; ... drive audio, foliage, or weather from biome ...
BIOME.FREE(biome)
```

---

## See also

- [WEATHER.md](WEATHER.md)
- [TERRAIN.md](TERRAIN.md)
