# Biomes (`BIOME.*`)

Lightweight **temperature / humidity** state (`TagBiome`) for driving ambience or future terrain/weather blending. **CGO** required.

---

## `Biome.Make()` → handle

Creates a biome descriptor.

---

## `Biome.Free(biome)`

Frees the handle.

---

## `Biome.SetTemp(biome, t#)` / `Biome.SetHumidity(biome, h#)`

Sets normalized or arbitrary scales (see runtime) for **temperature** and **humidity**.

---

## Notes

Biomes **do not** duplicate **`Scatter.Apply`** or **`Terrain.FillPerlin`** — combine them in your game loop or data layer.

---

## See also

- [WEATHER.md](WEATHER.md)
- [TERRAIN.md](TERRAIN.md)
