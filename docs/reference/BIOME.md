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

### `BIOME.MAKE(name)` 

Alias for `BIOME.CREATE`.

---

### `BIOME.SETTEMP(biomeHandle, temperature)` 

Sets the temperature value for the biome (typically 0.0–1.0).

---

### `BIOME.SETHUMIDITY(biomeHandle, humidity)` 

Sets the humidity value for the biome (typically 0.0–1.0).

---

### `BIOME.FREE(biomeHandle)` 

Frees the biome handle from memory.

---

## Full Example

This example creates two biomes and assigns properties for a procedural terrain system.

```basic
desert = BIOME.CREATE("desert")
BIOME.SETTEMP(desert, 0.9)
BIOME.SETHUMIDITY(desert, 0.1)

jungle = BIOME.CREATE("jungle")
BIOME.SETTEMP(jungle, 0.7)
BIOME.SETHUMIDITY(jungle, 0.9)

PRINT "Biomes configured."

; Use biome handles with terrain chunk placement ...

BIOME.FREE(desert)
BIOME.FREE(jungle)
```

---

## See also

- [WEATHER.md](WEATHER.md)
- [TERRAIN.md](TERRAIN.md)
